package rtc

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

// MediaTrack represents a published audio or video track within a room session.
type MediaTrack struct {
	mu           sync.RWMutex
	id           string
	kind         webrtc.RTPCodecType
	publisherID  string
	simulcast    bool
	isMuted      bool
	subscriberCount int
}

// NewMediaTrack constructs a new published MediaTrack.
func NewMediaTrack(id string, kind webrtc.RTPCodecType, publisherID string) *MediaTrack {
	return &MediaTrack{
		id:          id,
		kind:        kind,
		publisherID: publisherID,
	}
}

func (m *MediaTrack) ID() string {
	return m.id
}

func (m *MediaTrack) Kind() webrtc.RTPCodecType {
	return m.kind
}

func (m *MediaTrack) PublisherID() string {
	return m.publisherID
}

func (m *MediaTrack) SetMuted(muted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isMuted = muted
}

func (m *MediaTrack) IsMuted() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isMuted
}

func (m *MediaTrack) IncrementSubscribers() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscriberCount++
	return m.subscriberCount
}

func (m *MediaTrack) DecrementSubscribers() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.subscriberCount > 0 {
		m.subscriberCount--
	}
	return m.subscriberCount
}
