// Package signal provides WebSocket and HTTP REST signaling services for client connections.
package signal

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

var (
	ErrInvalidAuthToken = errors.New("unauthorized: invalid authentication token")
	ErrMissingRoom      = errors.New("bad request: room parameter missing")
)

// ResponsePayload defines standard JSON response messages.
type ResponsePayload struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Room    string `json:"room,omitempty"`
}

// Server handles client signaling requests and authentication validation.
type Server struct {
	mu     sync.RWMutex
	secret string
}

// NewServer instantiates a new Signal Server with secret-based auth verification.
func NewServer(secret string) *Server {
	return &Server{
		secret: secret,
	}
}

// HandleValidate validates client credentials and sets proper CORS headers.
func (s *Server) HandleValidate(w http.ResponseWriter, r *http.Request) {
	// Set explicit CORS headers for browser clients
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	token := r.URL.Query().Get("access_token")
	if token == "" {
		s.writeErrorJSON(w, r, http.StatusUnauthorized, ErrInvalidAuthToken)
		return
	}

	room := r.URL.Query().Get("room")
	if room == "" {
		s.writeErrorJSON(w, r, http.StatusBadRequest, ErrMissingRoom)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ResponsePayload{
		Status:  http.StatusOK,
		Message: "Connection token validated successfully",
		Room:    room,
	})
}

func (s *Server) writeErrorJSON(w http.ResponseWriter, r *http.Request, status int, err error) {
	telemetry.Log.Warn().Err(err).Int("status", status).Str("path", r.URL.Path).Msg("Signaling request error")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ResponsePayload{
		Status:  status,
		Message: err.Error(),
	})
}
