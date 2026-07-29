// Package ingress provides WHIP WebRTC HTTP Ingestion Protocol & RTSP media stream ingest.
package ingress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// WHIPSession represents an incoming WebRTC stream publish request via WHIP HTTP POST.
type WHIPSession struct {
	ID        string    `json:"id"`
	RoomName  string    `json:"room_name"`
	StreamID  string    `json:"stream_id"`
	SDPOffer  string    `json:"sdp_offer"`
	SDPAnswer string    `json:"sdp_answer"`
	CreatedAt time.Time `json:"created_at"`
}

// WHIPServer handles IETF WHIP protocol standard endpoints.
type WHIPServer struct {
	mu       sync.RWMutex
	port     int
	sessions map[string]*WHIPSession
}

// NewWHIPServer creates a new WHIP ingestion server instance.
func NewWHIPServer(port int) *WHIPServer {
	return &WHIPServer{
		port:     port,
		sessions: make(map[string]*WHIPSession),
	}
}

// Start launches the WHIP HTTP listener for OBS Studio / WHIP clients.
func (s *WHIPServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/whip/v1/stream", s.handleWHIPPublish)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	telemetry.Log.Info().Int("port", s.port).Msg("WHIP Stream Ingress Server listening")

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func (s *WHIPServer) handleWHIPPublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "default-ingress-room"
	}

	sessionID := fmt.Sprintf("whip-%d", time.Now().UnixNano()/1e6)
	session := &WHIPSession{
		ID:        sessionID,
		RoomName:  roomName,
		StreamID:  fmt.Sprintf("stream-%s", sessionID),
		CreatedAt: time.Now(),
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/sdp")
	w.Header().Set("Location", fmt.Sprintf("/whip/v1/stream/%s", sessionID))
	w.WriteHeader(http.StatusCreated)

	telemetry.Log.Info().Str("session_id", sessionID).Str("room", roomName).Msg("Accepted WHIP stream publish request")
}
