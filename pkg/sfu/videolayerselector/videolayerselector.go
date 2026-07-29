// Package videolayerselector handles SVC (Scalable Video Coding) and Simulcast layer switching.
package videolayerselector

import (
	"sync"
)

type VideoSpatialLayer int

const (
	SpatialLow VideoSpatialLayer = iota
	SpatialMid
	SpatialHigh
)

// Selector manages layer switching decisions based on subscriber bandwidth.
type Selector struct {
	mu           sync.RWMutex
	currentLayer VideoSpatialLayer
}

func NewSelector() *Selector {
	return &Selector{
		currentLayer: SpatialHigh,
	}
}

func (s *Selector) SelectLayer(targetBitrate uint32) VideoSpatialLayer {
	s.mu.Lock()
	defer s.mu.Unlock()

	if targetBitrate < 300_000 {
		s.currentLayer = SpatialLow
	} else if targetBitrate < 1_000_000 {
		s.currentLayer = SpatialMid
	} else {
		s.currentLayer = SpatialHigh
	}
	return s.currentLayer
}

func (s *Selector) CurrentLayer() VideoSpatialLayer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentLayer
}
