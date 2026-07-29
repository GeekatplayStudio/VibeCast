package cluster

import (
	"context"
	"sync"
)

// MemoryRouter implements an in-memory, thread-safe Router for standalone deployments.
type MemoryRouter struct {
	mu        sync.RWMutex
	nodes     map[string]*Node
	roomNodes map[string]string
}

// NewMemoryRouter instantiates a new thread-safe memory router.
func NewMemoryRouter() *MemoryRouter {
	return &MemoryRouter{
		nodes:     make(map[string]*Node),
		roomNodes: make(map[string]string),
	}
}

func (m *MemoryRouter) RegisterNode(_ context.Context, node *Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes[node.ID] = node
	return nil
}

func (m *MemoryRouter) UnregisterNode(_ context.Context, nodeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.nodes, nodeID)
	return nil
}

func (m *MemoryRouter) GetNodeForRoom(_ context.Context, roomName string) (*Node, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodeID, ok := m.roomNodes[roomName]
	if !ok {
		return nil, ErrRoomNotFound
	}

	node, ok := m.nodes[nodeID]
	if !ok {
		return nil, ErrNodeNotFound
	}

	return node, nil
}

func (m *MemoryRouter) SetNodeForRoom(_ context.Context, roomName string, nodeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.roomNodes[roomName] = nodeID
	return nil
}

func (m *MemoryRouter) ClearRoomState(_ context.Context, roomName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.roomNodes, roomName)
	return nil
}
