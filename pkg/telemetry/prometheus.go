// Package telemetry provides Prometheus metrics endpoints for WebRTC media server monitoring.
package telemetry

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// MetricsCollector tracks aggregate WebRTC transport and MCP tool counters.
type MetricsCollector struct {
	activeRooms       uint64
	activeTrackCount  uint64
	bytesReceived     uint64
	bytesSent         uint64
	mcpToolExecutions uint64
}

// GlobalCollector instance.
var GlobalCollector = &MetricsCollector{}

// IncRoom increments active room count.
func (m *MetricsCollector) IncRoom() { atomic.AddUint64(&m.activeRooms, 1) }

// DecRoom decrements active room count.
func (m *MetricsCollector) DecRoom() { atomic.AddUint64(&m.activeRooms, ^uint64(0)) }

// IncMCPTool increments MCP tool call count.
func (m *MetricsCollector) IncMCPTool() { atomic.AddUint64(&m.mcpToolExecutions, 1) }

// MetricsHandler serves OpenMetrics / Prometheus plaintext payload.
func (m *MetricsCollector) MetricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP agentic_sfu_active_rooms Current active WebRTC room sessions\n")
	fmt.Fprintf(w, "# TYPE agentic_sfu_active_rooms gauge\n")
	fmt.Fprintf(w, "agentic_sfu_active_rooms %d\n", atomic.LoadUint64(&m.activeRooms))

	fmt.Fprintf(w, "# HELP agentic_sfu_mcp_tool_calls_total Total MCP tool executions\n")
	fmt.Fprintf(w, "# TYPE agentic_sfu_mcp_tool_calls_total counter\n")
	fmt.Fprintf(w, "agentic_sfu_mcp_tool_calls_total %d\n", atomic.LoadUint64(&m.mcpToolExecutions))
}
