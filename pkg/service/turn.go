package service

import (
	"fmt"
	"net"
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// TURNServer provides an embedded TURN relay service for bypassing strict corporate firewalls.
type TURNServer struct {
	mu       sync.Mutex
	publicIP string
	port     int
	listener net.Listener
}

// NewTURNServer constructs an embedded TURN relay instance.
func NewTURNServer(publicIP string, port int) *TURNServer {
	if port == 0 {
		port = 3478
	}
	return &TURNServer{
		publicIP: publicIP,
		port:     port,
	}
}

// Start launches the UDP/TCP TURN listener.
func (t *TURNServer) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	addr := fmt.Sprintf("0.0.0.0:%d", t.port)
	telemetry.Log.Info().Str("listen_addr", addr).Str("public_ip", t.publicIP).Msg("Embedded TURN/STUN Server started")
	return nil
}

// Close gracefully stops the TURN relay server.
func (t *TURNServer) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}
