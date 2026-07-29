// Package connectionquality provides participant network quality scoring and MOS metrics calculation.
package connectionquality

import (
	"math"
	"sync"
	"time"
)

// QualityScore represents participant connection grade (1 to 5).
type QualityScore int

const (
	ScorePoor      QualityScore = 1
	ScoreFair      QualityScore = 2
	ScoreGood      QualityScore = 3
	ScoreVeryGood  QualityScore = 4
	ScoreExcellent QualityScore = 5
)

// Stats holds raw RTC transport metrics.
type Stats struct {
	PacketLossFraction float64       `json:"packet_loss_fraction"`
	JitterMs           float64       `json:"jitter_ms"`
	RttMs              float64       `json:"rtt_ms"`
	BitrateBps         uint64        `json:"bitrate_bps"`
	Score              QualityScore  `json:"score"`
	LastUpdated        time.Time     `json:"last_updated"`
}

// Manager tracks and calculates network quality scores across participants.
type Manager struct {
	mu    sync.RWMutex
	stats map[string]*Stats
}

// NewManager initializes a Connection Quality Manager.
func NewManager() *Manager {
	return &Manager{
		stats: make(map[string]*Stats),
	}
}

// UpdateStats calculates MOS (Mean Opinion Score) based on loss, jitter, and RTT.
func (m *Manager) UpdateStats(participantID string, lossFraction, jitterMs, rttMs float64, bitrateBps uint64) QualityScore {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Compute effective latency: RTT + Jitter * 2
	effectiveLatency := rttMs + (jitterMs * 2.0)
	
	// Compute R-factor score calculation
	rFactor := 93.2 - (effectiveLatency / 40.0)
	if effectiveLatency > 160.0 {
		rFactor = 93.2 - ((effectiveLatency - 120.0) / 10.0)
	}
	rFactor -= (lossFraction * 100.0 * 2.5)
	if rFactor < 0 {
		rFactor = 0
	}

	score := ScoreExcellent
	switch {
	case rFactor >= 80:
		score = ScoreExcellent
	case rFactor >= 70:
		score = ScoreVeryGood
	case rFactor >= 60:
		score = ScoreGood
	case rFactor >= 50:
		score = ScoreFair
	default:
		score = ScorePoor
	}

	m.stats[participantID] = &Stats{
		PacketLossFraction: math.Round(lossFraction*1000) / 1000,
		JitterMs:           math.Round(jitterMs*10) / 10,
		RttMs:              math.Round(rttMs*10) / 10,
		BitrateBps:         bitrateBps,
		Score:              score,
		LastUpdated:        time.Now(),
	}

	return score
}

// GetStats returns stats for a specific participant.
func (m *Manager) GetStats(participantID string) (*Stats, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	st, ok := m.stats[participantID]
	return st, ok
}
