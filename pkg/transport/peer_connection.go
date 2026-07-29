// Package transport provides thread-safe WebRTC PeerConnection wrappers and SDP negotiation.
package transport

import (
	"context"
	"fmt"
	"sync"

	"github.com/pion/webrtc/v3"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// Transport encapsulates a WebRTC PeerConnection and track subscribers.
type Transport struct {
	mu             sync.RWMutex
	id             string
	peerConnection *webrtc.PeerConnection
	onTrack        func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver)
}

// NewTransport creates and initializes a new Pion WebRTC PeerConnection.
func NewTransport(id string, api *webrtc.API) (*Transport, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pc, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.fmt.Errorf("failed to create peer connection: %w", err)
	}

	t := &Transport{
		id:             id,
		peerConnection: pc,
	}

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		t.mu.RLock()
		handler := t.onTrack
		t.mu.RUnlock()

		if handler != nil {
			handler(track, receiver)
		}
	})

	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		telemetry.Log.Info().Str("transport_id", id).Str("state", state.String()).Msg("ICE Connection State changed")
	})

	return t, nil
}

// SetOnTrack assigns the incoming media track callback handler.
func (t *Transport) SetOnTrack(handler func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onTrack = handler
}

// PeerConnection returns the underlying Pion PeerConnection.
func (t *Transport) PeerConnection() *webrtc.PeerConnection {
	return t.peerConnection
}

// Close gracefully closes the WebRTC peer connection.
func (t *Transport) Close(_ context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.peerConnection != nil {
		return t.peerConnection.Close()
	}
	return nil
}
