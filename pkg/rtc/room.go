package rtc

import (
	"context"
	"errors"
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/sfu"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

var (
	ErrParticipantAlreadyExists = errors.New("participant identity already exists in room")
	ErrParticipantNotFound      = errors.New("participant identity not found in room")
)

// Room manages connected participants, audio/video track distribution, and SFU media routing.
type Room struct {
	mu           sync.RWMutex
	name         string
	sid          string
	participants map[string]*Participant
	sfuRouter    *sfu.Router
	emptyTimeout int // in seconds
}

// NewRoom creates a new active Room instance.
func NewRoom(name string, sid string) *Room {
	return &Room{
		name:         name,
		sid:          sid,
		participants: make(map[string]*Participant),
		sfuRouter:    sfu.NewRouter(name),
		emptyTimeout: 300,
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) SID() string {
	return r.sid
}

func (r *Room) SFURouter() *sfu.Router {
	return r.sfuRouter
}

// JoinParticipant adds a new participant to the room state.
func (r *Room) JoinParticipant(p *Participant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.participants[p.Identity()]; exists {
		return ErrParticipantAlreadyExists
	}

	r.participants[p.Identity()] = p
	p.SetState(StateActive)

	telemetry.Log.Info().Str("room", r.name).Str("identity", p.Identity()).Msg("Participant joined room")
	return nil
}

// LeaveParticipant removes a participant and cleans up their published media tracks.
func (r *Room) LeaveParticipant(_ context.Context, identity string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.participants[identity]
	if !exists {
		return ErrParticipantNotFound
	}

	p.SetState(StateDisconnected)
	delete(r.participants, identity)

	telemetry.Log.Info().Str("room", r.name).Str("identity", identity).Msg("Participant left room")
	return nil
}

func (r *Room) GetParticipant(identity string) *Participant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.participants[identity]
}

func (r *Room) ParticipantCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.participants)
}

func (r *Room) Participants() []*Participant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*Participant, 0, len(r.participants))
	for _, p := range r.participants {
		list = append(list, p)
	}
	return list
}
