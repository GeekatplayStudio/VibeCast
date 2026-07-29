// Package streamtracker monitors RTP packet arrival cadence to detect track stalls or muted inputs.
package streamtracker

import (
	"sync"
	"time"
)

type StreamState int

const (
	StateInactive StreamState = iota
	StateActive
)

// StreamTracker tracks packet timestamps to signal track activity changes.
type StreamTracker struct {
	mu           sync.RWMutex
	lastPacketAt time.Time
	state        StreamState
	timeout      time.Duration
}

func NewStreamTracker(timeout time.Duration) *StreamTracker {
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	return &StreamTracker{
		state:   StateInactive,
		timeout: timeout,
	}
}

// ObservePacket updates the last seen timestamp of a stream.
func (st *StreamTracker) ObservePacket() {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.lastPacketAt = time.Now()
	st.state = StateActive
}

// CheckStatus evaluates whether the stream has timed out and become inactive.
func (st *StreamTracker) CheckStatus() StreamState {
	st.mu.Lock()
	defer st.mu.Unlock()

	if time.Since(st.lastPacketAt) > st.timeout {
		st.state = StateInactive
	}
	return st.state
}
