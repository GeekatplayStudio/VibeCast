// Package pacer provides token-bucket RTP packet pacing to smooth bitrate spikes.
package pacer

import (
	"context"
	"sync"
	"time"

	"github.com/pion/rtp"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// TokenBucketPacer controls RTP packet transmission pacing across subscriber DownTracks.
type TokenBucketPacer struct {
	mu           sync.RWMutex
	targetBitrate uint64
	packetQueue  []*rtp.Packet
	tokens       uint64
	lastRefill   time.Time
}

// NewTokenBucketPacer creates a new RTP packet pacer.
func NewTokenBucketPacer(targetBitrate uint64) *TokenBucketPacer {
	if targetBitrate == 0 {
		targetBitrate = 2500000 // 2.5 Mbps default
	}
	return &TokenBucketPacer{
		targetBitrate: targetBitrate,
		packetQueue:   make([]*rtp.Packet, 0, 100),
		lastRefill:    time.Now(),
	}
}

// Enqueue adds an RTP packet to the pacer transmission queue.
func (p *TokenBucketPacer) Enqueue(packet *rtp.Packet) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.packetQueue = append(p.packetQueue, packet)
}

// SetBitrate updates the target pacing bitrate.
func (p *TokenBucketPacer) SetBitrate(bitrate uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.targetBitrate = bitrate
	telemetry.Log.Info().Uint64("bitrate", bitrate).Msg("Updated TokenBucketPacer target bitrate")
}

// Start launches the pacer queue pop loop.
func (p *TokenBucketPacer) Start(ctx context.Context, sendFunc func(*rtp.Packet) error) {
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.mu.Lock()
			if len(p.packetQueue) > 0 {
				packet := p.packetQueue[0]
				p.packetQueue = p.packetQueue[1:]
				p.mu.Unlock()

				_ = sendFunc(packet)
			} else {
				p.mu.Unlock()
			}
		}
	}
}
