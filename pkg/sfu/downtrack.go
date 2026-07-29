// Package sfu implements high-performance media track forwarding and packet distribution.
package sfu

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

// DownTrack represents a subscriber media track consuming incoming RTP packets.
type DownTrack struct {
	mu           sync.RWMutex
	id           string
	trackLocal   *webrtc.TrackLocalStaticRTP
	packetsSent  uint64
	bytesSent    uint64
	lastSeqNum   uint16
	seqOffset    uint16
}

// NewDownTrack initializes a new local static RTP track for downstream distribution.
func NewDownTrack(id string, codec webrtc.RTPCodecCapability, streamID string) (*DownTrack, error) {
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(codec, id, streamID)
	if err != nil {
		return nil, fmt.Errorf("failed to create local static rtp track: %w", err)
	}

	return &DownTrack{
		id:         id,
		trackLocal: trackLocal,
	}, nil
}

// WriteRTP forwards an incoming RTP packet to the downstream subscriber track with sequence adjustment.
func (d *DownTrack) WriteRTP(packet *rtp.Packet) error {
	d.mu.Lock()
	atomic.AddUint64(&d.packetsSent, 1)
	atomic.AddUint64(&d.bytesSent, uint64(packet.MarshalSize()))
	
	// Remap sequence number to maintain continuous packet index for subscribers
	d.lastSeqNum++
	modifiedPacket := *packet
	modifiedPacket.SequenceNumber = d.lastSeqNum
	d.mu.Unlock()

	return d.trackLocal.WriteRTP(&modifiedPacket)
}

// TrackLocal returns the Pion TrackLocalStaticRTP handle.
func (d *DownTrack) TrackLocal() *webrtc.TrackLocalStaticRTP {
	return d.trackLocal
}

// Stats returns the total packets and bytes sent through this DownTrack.
func (d *DownTrack) Stats() (uint64, uint64) {
	return atomic.LoadUint64(&d.packetsSent), atomic.LoadUint64(&d.bytesSent)
}

