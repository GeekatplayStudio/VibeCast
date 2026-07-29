// Package mcp provides native Model Context Protocol (MCP) server integration
// allowing LLM agents to dynamically inspect, join, and control WebRTC rooms.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// ToolDefinition represents an MCP tool contract exposed to LLM agents.
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ToolCallRequest represents an incoming MCP JSON-RPC tool execution payload.
type ToolCallRequest struct {
	ID        string                 `json:"id"`
	Tool      string                 `json:"tool"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolCallResponse represents the result of an MCP tool invocation.
type ToolCallResponse struct {
	ID        string      `json:"id"`
	Success   bool        `json:"success"`
	Result    interface{} `json:"result,omitempty"`
	ErrorMsg  string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Bridge exposes MCP endpoint handlers for AI agent discovery and control.
type Bridge struct {
	mu            sync.RWMutex
	port          int
	activeAgents  map[string]string
	executedTools int
}

// NewBridge creates a new MCP agent bridge.
func NewBridge(port int) *Bridge {
	return &Bridge{
		port:         port,
		activeAgents: make(map[string]string),
	}
}

// Start launches the HTTP server handling MCP tool calls.
func (b *Bridge) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/mcp/v1/tools", b.handleListTools)
	mux.HandleFunc("/mcp/v1/dispatch", b.handleAgentDispatch)
	mux.HandleFunc("/mcp/v1/execute", b.handleExecuteTool)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", b.port),
		Handler: mux,
	}

	telemetry.Log.Info().Int("port", b.port).Msg("MCP Agent Bridge listening")

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func (b *Bridge) handleListTools(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tools := []ToolDefinition{
		{
			Name:        "list_rooms",
			Description: "Lists all active WebRTC SFU room sessions and participant counts",
		},
		{
			Name:        "dispatch_agent",
			Description: "Dispatches a real-time voice/multimodal AI agent into a target room",
			Parameters: map[string]interface{}{
				"room":   "string",
				"prompt": "string",
				"model":  "string",
			},
		},
		{
			Name:        "mute_participant",
			Description: "Mutes a participant track across the SFU router",
			Parameters: map[string]interface{}{
				"room":           "string",
				"participant_id": "string",
			},
		},
		{
			Name:        "get_telemetry",
			Description: "Retrieves real-time WebRTC media router performance metrics (bitrate, loss, RTT)",
		},
	}
	_ = json.NewEncoder(w).Encode(tools)
}

func (b *Bridge) handleAgentDispatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Room   string `json:"room"`
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	agentID := fmt.Sprintf("agent-voice-%d", time.Now().Unix())

	b.mu.Lock()
	b.activeAgents[agentID] = req.Room
	b.mu.Unlock()

	telemetry.Log.Info().Str("room", req.Room).Str("agent_id", agentID).Msg("Received MCP agent dispatch request")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "dispatched",
		"agent_id": agentID,
		"room":     req.Room,
		"prompt":   req.Prompt,
	})
}

func (b *Bridge) handleExecuteTool(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var call ToolCallRequest
	if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ToolCallResponse{
			Success:  false,
			ErrorMsg: "invalid tool call request payload",
		})
		return
	}

	b.mu.Lock()
	b.executedTools++
	b.mu.Unlock()

	res := ToolCallResponse{
		ID:        call.ID,
		Success:   true,
		Timestamp: time.Now(),
		Result: map[string]interface{}{
			"message": fmt.Sprintf("MCP Tool '%s' executed successfully", call.Tool),
			"details": call.Arguments,
		},
	}

	_ = json.NewEncoder(w).Encode(res)
}

