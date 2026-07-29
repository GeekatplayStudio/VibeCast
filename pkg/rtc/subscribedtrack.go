package rtc

import (
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/sfu"
)

// SubscribedTrack wraps a downstream subscriber track bound to a publisher source.
type SubscribedTrack struct {
	mu           sync.RWMutex
	subscriberID string
	publisherID  string
	downTrack    *sfu.DownTrack
	isMuted      bool
}

// NewSubscribedTrack constructs a new SubscribedTrack binding.
func NewSubscribedTrack(subscriberID string, publisherID string, downTrack *sfu.DownTrack) *SubscribedTrack {
	return &SubscribedTrack{
		subscriberID: subscriberID,
		publisherID:  publisherID,
		downTrack:    downTrack,
	}
}

func (s *SubscribedTrack) SubscriberID() string {
	return s.subscriberID
}

func (s *SubscribedTrack) PublisherID() string {
	return s.publisherID
}

func (s *SubscribedTrack) DownTrack() *sfu.DownTrack {
	return s.downTrack
}

func (s *SubscribedTrack) SetMuted(muted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isMuted = muted
}

func (s *SubscribedTrack) IsMuted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isMuted
}
