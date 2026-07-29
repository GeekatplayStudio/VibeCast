package sfu

import (
	"sync"

	"github.com/pion/rtp"
)

// PacketBuffer implements a high-performance RTP packet ring-buffer for NACK retransmissions.
type PacketBuffer struct {
	mu       sync.RWMutex
	capacity uint16
	packets  map[uint16]*rtp.Packet
}

// NewPacketBuffer constructs a packet buffer with specified capacity.
func NewPacketBuffer(capacity uint16) *PacketBuffer {
	if capacity == 0 {
		capacity = 500
	}
	return &PacketBuffer{
		capacity: capacity,
		packets:  make(map[uint16]*rtp.Packet),
	}
}

// Push adds an RTP packet into the ring buffer.
func (b *PacketBuffer) Push(packet *rtp.Packet) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.packets[packet.SequenceNumber] = packet.Clone()
	if uint16(len(b.packets)) > b.capacity {
		// Evict oldest packet sequence number
		var oldestSeq uint16
		for seq := range b.packets {
			oldestSeq = seq
			break
		}
		delete(b.packets, oldestSeq)
	}
}

// GetPacket retrieves a buffered RTP packet by sequence number for NACK handling.
func (b *PacketBuffer) GetPacket(seq uint16) *rtp.Packet {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.packets[seq]
}
