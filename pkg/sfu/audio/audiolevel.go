// Package audio provides active speaker detection, audio level parsing, and loudness calculations.
package audio

import (
	"math"
	"sync"
)

// ActiveSpeakerObserver monitors audio levels across participants to detect who is currently speaking.
type ActiveSpeakerObserver struct {
	mu            sync.RWMutex
	thresholdDbov  float64
	activeSpeaker  string
	participantDbs map[string]float64
}

func NewActiveSpeakerObserver(thresholdDbov float64) *ActiveSpeakerObserver {
	if thresholdDbov == 0 {
		thresholdDbov = -35.0 // dBov threshold for voice activity
	}
	return &ActiveSpeakerObserver{
		thresholdDbov:  thresholdDbov,
		participantDbs: make(map[string]float64),
	}
}

// ObserveLevel updates a participant's current audio level and evaluates active speaker status.
func (a *ActiveSpeakerObserver) ObserveLevel(participantID string, levelDbov float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.participantDbs[participantID] = levelDbov

	if levelDbov >= a.thresholdDbov {
		a.activeSpeaker = participantID
	}
}

// ActiveSpeaker returns the identity of the currently active speaker.
func (a *ActiveSpeakerObserver) ActiveSpeaker() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.activeSpeaker
}

// GetLevels returns current dBov levels for all active audio streams.
func (a *ActiveSpeakerObserver) GetLevels() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	result := make(map[string]float64, len(a.participantDbs))
	for k, v := range a.participantDbs {
		result[k] = v
	}
	return result
}

// ComputePCM16RMS calculates root-mean-square energy and converts to dBov for 16-bit PCM audio samples.
func ComputePCM16RMS(samples []int16) float64 {
	if len(samples) == 0 {
		return -127.0
	}
	var sumSquares float64
	for _, sample := range samples {
		val := float64(sample) / 32768.0
		sumSquares += val * val
	}
	rms := math.Sqrt(sumSquares / float64(len(samples)))
	if rms <= 1e-6 {
		return -127.0
	}
	return 20.0 * math.Log10(rms)
}

// LinearToDBov converts raw PCM audio amplitude (0 to 32767) into dBov scale.
func LinearToDBov(linear uint16) float64 {
	if linear == 0 {
		return -127.0
	}
	ratio := float64(linear) / 32767.0
	return 20.0 * math.Log10(ratio)
}

