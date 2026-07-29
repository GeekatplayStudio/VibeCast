package rtc

import (
	"sync"
)

// SubscriptionRule defines filter criteria for track subscriptions.
type SubscriptionRule struct {
	AllowedTrackIDs []string `json:"allowed_track_ids"`
	PausedTrackIDs  []string `json:"paused_track_ids"`
	MaxQuality      string   `json:"max_quality"`
}

// SubscriptionManager manages dynamic track subscription state per participant.
type SubscriptionManager struct {
	mu    sync.RWMutex
	rules map[string]*SubscriptionRule
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		rules: make(map[string]*SubscriptionRule),
	}
}

func (s *SubscriptionManager) SetRule(participantID string, rule *SubscriptionRule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[participantID] = rule
}

func (s *SubscriptionManager) IsSubscribed(participantID string, trackID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rule, ok := s.rules[participantID]
	if !ok {
		return true // Default subscribe to all
	}

	for _, paused := range rule.PausedTrackIDs {
		if paused == trackID {
			return false
		}
	}
	return true
}
