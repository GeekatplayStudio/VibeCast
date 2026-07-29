// Package cluster defines the clean, interface-driven event router and room state management layer.
// Decouples room signaling and inter-node routing from WebRTC media forwarding.
package cluster

import (
	"context"
	"errors"
)

var (
	ErrRoomNotFound = errors.New("room not found in cluster router")
	ErrNodeNotFound = errors.New("node not found in cluster router")
)

// Node represents a cluster media server instance.
type Node struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// Router defines the interface for cluster room allocation and signaling routing.
type Router interface {
	// RegisterNode registers the local media server instance into the cluster registry.
	RegisterNode(ctx context.Context, node *Node) error

	// UnregisterNode removes the local media server instance from the cluster.
	UnregisterNode(ctx context.Context, nodeID string) error

	// GetNodeForRoom retrieves the assigned cluster node hosting the given room.
	GetNodeForRoom(ctx context.Context, roomName string) (*Node, error)

	// SetNodeForRoom maps a room name to a specific cluster node ID.
	SetNodeForRoom(ctx context.Context, roomName string, nodeID string) error

	// ClearRoomState cleans up room state when all participants leave.
	ClearRoomState(ctx context.Context, roomName string) error
}
