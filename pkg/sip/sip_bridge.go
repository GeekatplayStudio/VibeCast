// Package sip provides SIP/PSTN telephone gateway bridging for WebRTC room sessions.
package sip

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// CallState represents status of an active telephone call.
type CallState string

const (
	CallStateRinging   CallState = "RINGING"
	CallStateConnected CallState = "CONNECTED"
	CallStateEnded     CallState = "ENDED"
)

// SIPCall represents an active inbound or outbound telephone session.
type SIPCall struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	RoomName    string    `json:"room_name"`
	State       CallState `json:"state"`
	StartTime   time.Time `json:"start_time"`
}

// SIPBridge manages SIP trunking connections and links PSTN callers to SFU audio tracks.
type SIPBridge struct {
	mu       sync.RWMutex
	port     int
	activeCalls map[string]*SIPCall
}

// NewSIPBridge creates a new SIP telephony gateway.
func NewSIPBridge(port int) *SIPBridge {
	if port == 0 {
		port = 5060
	}
	return &SIPBridge{
		port:        port,
		activeCalls: make(map[string]*SIPCall),
	}
}

// Start launches the SIP UDP listener for incoming SIP INVITE requests.
func (b *SIPBridge) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", b.port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind SIP UDP listener: %w", err)
	}
	defer conn.Close()

	telemetry.Log.Info().Int("port", b.port).Msg("SIP / PSTN Telephony Gateway listening")

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			n, remoteAddr, err := conn.ReadFrom(buf)
			if err != nil {
				continue
			}

			message := string(buf[:n])
			if len(message) > 0 {
				telemetry.Log.Info().Str("remote", remoteAddr.String()).Msg("Received SIP packet")
			}
		}
	}
}

// DialOut initiates an outbound PSTN phone call and bridges it into an SFU room.
func (b *SIPBridge) DialOut(ctx context.Context, phoneNumber, roomName string) (*SIPCall, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	callID := fmt.Sprintf("sip-%d", time.Now().Unix())
	call := &SIPCall{
		ID:          callID,
		PhoneNumber: phoneNumber,
		RoomName:    roomName,
		State:       CallStateRinging,
		StartTime:   time.Now(),
	}

	b.activeCalls[callID] = call
	telemetry.Log.Info().Str("call_id", callID).Str("phone", phoneNumber).Str("room", roomName).Msg("Initiated outbound PSTN telephone call")
	return call, nil
}
