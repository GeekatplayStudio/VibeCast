package sfu

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/pion/webrtc/v3"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// Router coordinates publisher tracks and forwards media to subscriber DownTracks.
type Router struct {
	mu         sync.RWMutex
	roomName   string
	downTracks map[string][]*DownTrack
}

// NewRouter initializes an SFU Media Router for a room.
func NewRouter(roomName string) *Router {
	return &Router{
		roomName:   roomName,
		downTracks: make(map[string][]*DownTrack),
	}
}

// AddPublisherTrack binds an incoming remote track and spawns a packet read loop.
func (r *Router) AddPublisherTrack(ctx context.Context, track *webrtc.TrackRemote) {
	trackID := track.ID()
	telemetry.Log.Info().Str("room", r.roomName).Str("track_id", trackID).Str("kind", track.Kind().String()).Msg("Adding publisher track")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				packet, _, err := track.ReadRTP()
				if err == io.EOF {
					return
				}
				if err != nil {
					continue
				}

				r.mu.RLock()
				subscribers := r.downTracks[trackID]
				r.mu.RUnlock()

				for _, dt := range subscribers {
					_ = dt.WriteRTP(packet)
				}
			}
		}
	}()
}

// Subscribe creates a subscriber DownTrack bound to a publisher track ID.
func (r *Router) Subscribe(trackID string, codec webrtc.RTPCodecCapability, streamID string) (*DownTrack, error) {
	dt, err := NewDownTrack(fmt.sprintf("sub-%s", trackID), codec, streamID)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.downTracks[trackID] = append(r.downTracks[trackID], dt)
	r.mu.Unlock()

	return dt, nil
}
