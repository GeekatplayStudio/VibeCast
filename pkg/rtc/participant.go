package rtc

import (
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/transport"
)

// ParticipantState represents the lifecycle status of a connected user or agent.
type ParticipantState int

const (
	StateJoining ParticipantState = iota
	StateActive
	StateDisconnected
)

// Participant represents a connected client or AI agent in a room.
type Participant struct {
	mu           sync.RWMutex
	identity     string
	name         string
	state        ParticipantState
	permissions  Permissions
	joinedAt     time.Time
	publisherPC  *transport.Transport
	subscriberPC *transport.Transport
	published    map[string]*MediaTrack
}

// NewParticipant initializes a room participant with assigned identity and permissions.
func NewParticipant(identity string, name string, perms Permissions) *Participant {
	return &Participant{
		identity:    identity,
		name:        name,
		state:       StateJoining,
		permissions: perms,
		joinedAt:    time.Now(),
		published:   make(map[string]*MediaTrack),
	}
}

func (p *Participant) Identity() string {
	return p.identity
}

func (p *Participant) Name() string {
	return p.name
}

func (p *Participant) SetState(state ParticipantState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.state = state
}

func (p *Participant) State() ParticipantState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

func (p *Participant) Permissions() Permissions {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.permissions
}

func (p *Participant) AddTrack(track *MediaTrack) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.published[track.ID()] = track
}

func (p *Participant) GetTrack(trackID string) *MediaTrack {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.published[trackID]
}

func (p *Participant) RemoveTrack(trackID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.published, trackID)
}

func (p *Participant) PublishedTracks() []*MediaTrack {
	p.mu.RLock()
	defer p.mu.RUnlock()
	tracks := make([]*MediaTrack, 0, len(p.published))
	for _, t := range p.published {
		tracks = append(tracks, t)
	}
	return tracks
}
