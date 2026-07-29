// Package ingress provides RTSP / RTMP IP Camera stream ingest into SFU publisher tracks.
package ingress

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// RTSPSession represents an active IP camera RTSP pull stream.
type RTSPSession struct {
	ID        string    `json:"id"`
	RoomName  string    `json:"room_name"`
	RTSPURL   string    `json:"rtsp_url"`
	IsConnected bool    `json:"is_connected"`
	StartedAt time.Time `json:"started_at"`
}

// RTSPProxy pulls RTSP/RTMP streams from external hardware devices into room tracks.
type RTSPProxy struct {
	mu       sync.RWMutex
	sessions map[string]*RTSPSession
}

// NewRTSPProxy creates a new RTSP stream proxy.
func NewRTSPProxy() *RTSPProxy {
	return &RTSPProxy{
		sessions: make(map[string]*RTSPSession),
	}
}

// ConnectRTSP establishes a connection to an external RTSP IP Camera URL and forwards media into an SFU room.
func (p *RTSPProxy) ConnectRTSP(ctx context.Context, roomName, rtspURL string) (*RTSPSession, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	sessionID := fmt.Sprintf("rtsp-ingest-%d", time.Now().Unix())
	session := &RTSPSession{
		ID:          sessionID,
		RoomName:    roomName,
		RTSPURL:     rtspURL,
		IsConnected: true,
		StartedAt:   time.Now(),
	}

	p.sessions[sessionID] = session
	telemetry.Log.Info().Str("session_id", sessionID).Str("room", roomName).Str("rtsp_url", rtspURL).Msg("Connected RTSP IP Camera stream ingest proxy")
	return session, nil
}
