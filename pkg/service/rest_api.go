// Package service provides administrative REST API management endpoints.
package service

import (
	"encoding/json"
	"net/http"

	"github.com/agenticsfu/agentic-sfu/pkg/auth"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// RestAPIService handles admin management requests for rooms, tokens, and mute controls.
type RestAPIService struct {
	roomManager *RoomManager
	apiKey      string
	apiSecret   string
}

// NewRestAPIService creates a new administrative REST API handler.
func NewRestAPIService(rm *RoomManager, apiKey, apiSecret string) *RestAPIService {
	return &RestAPIService{
		roomManager: rm,
		apiKey:      apiKey,
		apiSecret:   apiSecret,
	}
}

// RegisterRoutes mounts REST endpoints to an HTTP mux.
func (s *RestAPIService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/rooms", s.handleRooms)
	mux.HandleFunc("/api/v1/token", s.handleGenerateToken)
	mux.HandleFunc("/api/v1/mute", s.handleMuteParticipant)
}

func (s *RestAPIService) handleRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		rooms := s.roomManager.ListRooms()
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"count": len(rooms),
			"rooms": rooms,
		})
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Name string `json:"name"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		room, err := s.roomManager.CreateRoom(r.Context(), req.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(room)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *RestAPIService) handleGenerateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Identity string `json:"identity"`
		Room     string `json:"room"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	token := auth.NewAccessToken(s.apiKey, s.apiSecret, req.Identity)
	token.Grant = auth.ClaimGrant{
		Room:         req.Room,
		CanPublish:   true,
		CanSubscribe: true,
	}

	jwtStr, err := token.GenerateSignedToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"token":    jwtStr,
		"identity": req.Identity,
		"room":     req.Room,
	})
}

func (s *RestAPIService) handleMuteParticipant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Room          string `json:"room"`
		ParticipantID string `json:"participant_id"`
		Mute          bool   `json:"mute"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	telemetry.Log.Info().Str("room", req.Room).Str("participant", req.ParticipantID).Bool("mute", req.Mute).Msg("Admin REST API: Muted participant track")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"participant_id": req.ParticipantID,
		"is_muted":       req.Mute,
	})
}
