// Package bwe implements Transport-Wide Congestion Control (TWCC) and real-time Bandwidth Estimation.
package bwe

import (
	"sync"
	"time"
)

// BandwidthEstimator computes adaptive target bitrates based on TWCC delay signals.
type BandwidthEstimator struct {
	mu            sync.RWMutex
	targetBitrate uint32
	lastUpdate    time.Time
}

func NewBandwidthEstimator() *BandwidthEstimator {
	return &BandwidthEstimator{
		targetBitrate: 2_500_000, // 2.5 Mbps initial estimate
		lastUpdate:    time.Now(),
	}
}

// OnTWCCFeedback processes incoming TWCC feedback packets and adjusts target bitrates.
func (b *BandwidthEstimator) OnTWCCFeedback(rttMs uint32, lossRate float64) uint32 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if lossRate > 0.10 {
		// Reduce target bitrate by 15% on high loss
		b.targetBitrate = uint32(float64(b.targetBitrate) * 0.85)
	} else if lossRate < 0.02 && rttMs < 100 {
		// Probe upward by 5% on stable low-latency connection
		b.targetBitrate = uint32(float64(b.targetBitrate) * 1.05)
	}

	b.lastUpdate = time.Now()
	return b.targetBitrate
}

func (b *BandwidthEstimator) TargetBitrate() uint32 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.targetBitrate
}
