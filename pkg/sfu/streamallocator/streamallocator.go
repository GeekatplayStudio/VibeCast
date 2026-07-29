// Package streamallocator manages dynamic bandwidth allocation across multiple video streams.
package streamallocator

import (
	"sync"
)

// StreamAllocator dynamically distributes available bandwidth across published tracks.
type StreamAllocator struct {
	mu           sync.RWMutex
	totalBitrate uint32
}

func NewStreamAllocator() *StreamAllocator {
	return &StreamAllocator{
		totalBitrate: 5_000_000, // 5 Mbps default allocation
	}
}

func (s *StreamAllocator) AllocateBandwidth(trackCount int) uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if trackCount == 0 {
		return s.totalBitrate
	}
	return s.totalBitrate / uint32(trackCount)
}
