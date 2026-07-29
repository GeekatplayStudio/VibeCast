// Package buffer manages RTP packet buffering, NACK queues, and packet retransmission.
package buffer

import (
	"sync"

	"github.com/pion/rtp"
)

// Buffer manages an RTP packet ring-buffer with NACK retransmission lookup.
type Buffer struct {
	mu       sync.RWMutex
	capacity uint16
	packets  map[uint16]*rtp.Packet
}

// NewBuffer constructs a packet buffer instance.
func NewBuffer(capacity uint16) *Buffer {
	if capacity == 0 {
		capacity = 1000
	}
	return &Buffer{
		capacity: capacity,
		packets:  make(map[uint16]*rtp.Packet),
	}
}

// Push adds an RTP packet to the ring buffer.
func (b *Buffer) Push(packet *rtp.Packet) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.packets[packet.SequenceNumber] = packet.Clone()
	if uint16(len(b.packets)) > b.capacity {
		var oldest uint16
		for seq := range b.packets {
			oldest = seq
			break
		}
		delete(b.packets, oldest)
	}
}

// GetPacket retrieves a packet by sequence number for NACK retransmission.
func (b *Buffer) GetPacket(seq uint16) *rtp.Packet {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.packets[seq]
}
