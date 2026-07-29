// Package sfu provides screen sharing composition and track prioritization.
package sfu

import (
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// ScreenShareTrack metadata for screen capture streams.
type ScreenShareTrack struct {
	TrackID    string `json:"track_id"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	MaxFPS     int    `json:"max_fps"`
	IsPriority bool   `json:"is_priority"`
}

// ScreenShareManager prioritizes screen sharing video tracks over camera streams.
type ScreenShareManager struct {
	mu           sync.RWMutex
	activeTrack  *ScreenShareTrack
	priorityBits uint64
}

// NewScreenShareManager creates a screen sharing track manager instance.
func NewScreenShareManager() *ScreenShareManager {
	return &ScreenShareManager{}
}

// RegisterScreenShare binds an active screen sharing track.
func (m *ScreenShareManager) RegisterScreenShare(trackID string, width, height, maxFPS int) *ScreenShareTrack {
	m.mu.Lock()
	defer m.mu.Unlock()

	track := &ScreenShareTrack{
		TrackID:    trackID,
		Width:      width,
		Height:     height,
		MaxFPS:     maxFPS,
		IsPriority: true,
	}

	m.activeTrack = track
	telemetry.Log.Info().Str("track_id", trackID).Int("width", width).Int("height", height).Int("fps", maxFPS).Msg("Registered active high-resolution screen share track")
	return track
}

// ActiveScreenShare returns the current active screen share track.
func (m *ScreenShareManager) ActiveScreenShare() (*ScreenShareTrack, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.activeTrack == nil {
		return nil, false
	}
	return m.activeTrack, true
}
