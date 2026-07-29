package cluster_test

import (
	"context"
	"testing"

	"github.com/agenticsfu/agentic-sfu/pkg/cluster"
)

func TestMemoryRouter(t *testing.T) {
	ctx := context.Background()
	r := cluster.NewMemoryRouter()

	node := &cluster.Node{
		ID:      "node-1",
		Address: "127.0.0.1",
		Port:    7880,
	}

	// Test Node Registration
	err := r.RegisterNode(ctx, node)
	if err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Test Setting Room Node
	err = r.SetNodeForRoom(ctx, "room-alpha", "node-1")
	if err != nil {
		t.Fatalf("failed to set room node: %v", err)
	}

	// Test Getting Room Node
	foundNode, err := r.GetNodeForRoom(ctx, "room-alpha")
	if err != nil {
		t.Fatalf("failed to get node for room: %v", err)
	}
	if foundNode.ID != "node-1" {
		t.Errorf("expected node-1, got %s", foundNode.ID)
	}

	// Test Clearing Room State
	err = r.ClearRoomState(ctx, "room-alpha")
	if err != nil {
		t.Fatalf("failed to clear room state: %v", err)
	}

	_, err = r.GetNodeForRoom(ctx, "room-alpha")
	if err != cluster.ErrRoomNotFound {
		t.Errorf("expected ErrRoomNotFound, got %v", err)
	}
}
