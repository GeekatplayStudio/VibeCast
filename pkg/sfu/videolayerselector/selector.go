// Package videolayerselector provides spatial and temporal layer switching for VP8/H264/AV1 simulcast.
package videolayerselector

import (
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// SpatialLayer represents video resolution quality tier (0: Low/180p, 1: Medium/360p, 2: High/720p).
type SpatialLayer int

const (
	SpatialLow    SpatialLayer = 0
	SpatialMedium SpatialLayer = 1
	SpatialHigh   SpatialLayer = 2
)

// TemporalLayer represents frame rate quality tier (0: 15fps, 1: 30fps, 2: 60fps).
type TemporalLayer int

// Selector evaluates available subscriber bandwidth and dynamically selects optimal spatial/temporal video layer.
type Selector struct {
	mu             sync.RWMutex
	currentSpatial SpatialLayer
	currentTemporal TemporalLayer
	maxSpatial     SpatialLayer
	maxTemporal    TemporalLayer
}

// NewSelector creates a simulcast layer selector instance.
func NewSelector(maxSpatial SpatialLayer, maxTemporal TemporalLayer) *Selector {
	return &Selector{
		currentSpatial: maxSpatial,
		currentTemporal: maxTemporal,
		maxSpatial:     maxSpatial,
		maxTemporal:    maxTemporal,
	}
}

// SelectLayer selects the appropriate spatial/temporal layer based on subscriber estimated bandwidth in Bps.
func (s *Selector) SelectLayer(estimatedBitrateBps uint64) (SpatialLayer, TemporalLayer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	targetSpatial := SpatialHigh
	switch {
	case estimatedBitrateBps < 300000:
		targetSpatial = SpatialLow
	case estimatedBitrateBps < 800000:
		targetSpatial = SpatialMedium
	default:
		targetSpatial = SpatialHigh
	}

	if targetSpatial > s.maxSpatial {
		targetSpatial = s.maxSpatial
	}

	if targetSpatial != s.currentSpatial {
		telemetry.Log.Info().Int("old_layer", int(s.currentSpatial)).Int("new_layer", int(targetSpatial)).Uint64("bitrate", estimatedBitrateBps).Msg("Simulcast Selector switched video spatial layer")
		s.currentSpatial = targetSpatial
	}

	return s.currentSpatial, s.currentTemporal
}
