package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/cluster"
	"github.com/agenticsfu/agentic-sfu/pkg/rtc"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// RoomManager coordinates room creation, storage, and cluster node assignment.
type RoomManager struct {
	mu     sync.RWMutex
	store  *LocalStore
	router cluster.Router
	nodeID string
}

func NewRoomManager(store *LocalStore, router cluster.Router, nodeID string) *RoomManager {
	return &RoomManager{
		store:  store,
		router: router,
		nodeID: nodeID,
	}
}

// GetOrCreateRoom loads an existing room session or creates and registers a new room.
func (m *RoomManager) GetOrCreateRoom(ctx context.Context, name string) (*rtc.Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, err := m.store.LoadRoom(ctx, name)
	if err == nil {
		return room, nil
	}

	sid := fmt.Sprintf("sid-%d", time.Now().UnixNano())
	room = rtc.NewRoom(name, sid)

	if err := m.store.StoreRoom(ctx, room); err != nil {
		return nil, err
	}

	if err := m.router.SetNodeForRoom(ctx, name, m.nodeID); err != nil {
		telemetry.Log.Warn().Err(err).Str("room", name).Msg("Failed to update router node map for room")
	}

	telemetry.Log.Info().Str("room", name).Str("sid", sid).Msg("Allocated new room session")
	return room, nil
}

// CloseRoom tears down an active room session and removes cluster state.
func (m *RoomManager) CloseRoom(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_ = m.store.DeleteRoom(ctx, name)
	_ = m.router.ClearRoomState(ctx, name)

	telemetry.Log.Info().Str("room", name).Msg("Closed room session")
	return nil
}
