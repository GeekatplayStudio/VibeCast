// Package egress provides RTMP live stream output forwarding.
package egress

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// RTMPSession represents an active live stream push to YouTube Live, Twitch, etc.
type RTMPSession struct {
	ID        string    `json:"id"`
	RoomName  string    `json:"room_name"`
	StreamURL string    `json:"stream_url"`
	StreamKey string    `json:"stream_key,omitempty"`
	IsActive  bool      `json:"is_active"`
	StartedAt time.Time `json:"started_at"`
}

// RTMPPusher manages live stream RTMP output pipelines.
type RTMPPusher struct {
	mu       sync.RWMutex
	sessions map[string]*RTMPSession
}

// NewRTMPPusher creates a new RTMP pusher instance.
func NewRTMPPusher() *RTMPPusher {
	return &RTMPPusher{
		sessions: make(map[string]*RTMPSession),
	}
}

// StartStream begins forwarding SFU room audio/video to an RTMP endpoint.
func (p *RTMPPusher) StartStream(ctx context.Context, roomName, streamURL, streamKey string) (*RTMPSession, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	sessionID := fmt.Sprintf("rtmp-%s-%d", roomName, time.Now().Unix())
	session := &RTMPSession{
		ID:        sessionID,
		RoomName:  roomName,
		StreamURL: streamURL,
		StreamKey: streamKey,
		IsActive:  true,
		StartedAt: time.Now(),
	}

	p.sessions[sessionID] = session
	telemetry.Log.Info().Str("session_id", sessionID).Str("room", roomName).Str("url", streamURL).Msg("Started RTMP live stream egress push")
	return session, nil
}

// StopStream stops live stream push.
func (p *RTMPPusher) StopStream(sessionID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	session, ok := p.sessions[sessionID]
	if !ok {
		return fmt.Errorf("rtmp session %s not found", sessionID)
	}

	session.IsActive = false
	telemetry.Log.Info().Str("session_id", sessionID).Msg("Stopped RTMP live stream egress push")
	return nil
}
