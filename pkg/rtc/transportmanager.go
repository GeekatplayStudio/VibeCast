package rtc

import (
	"fmt"
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/transport"
	"github.com/pion/webrtc/v3"
)

// TransportManager manages dual Publisher and Subscriber WebRTC transports per participant.
type TransportManager struct {
	mu           sync.RWMutex
	publisher    *transport.Transport
	subscriber   *transport.Transport
	webrtcAPI    *webrtc.API
}

func NewTransportManager(api *webrtc.API) *TransportManager {
	return &TransportManager{
		webrtcAPI: api,
	}
}

func (tm *TransportManager) CreatePublisherTransport(participantID string) (*transport.Transport, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, err := transport.NewTransport(fmt.Sprintf("pub-%s", participantID), tm.webrtcAPI)
	if err != nil {
		return nil, err
	}
	tm.publisher = t
	return t, nil
}

func (tm *TransportManager) CreateSubscriberTransport(participantID string) (*transport.Transport, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, err := transport.NewTransport(fmt.Sprintf("sub-%s", participantID), tm.webrtcAPI)
	if err != nil {
		return nil, err
	}
	tm.subscriber = t
	return t, nil
}

func (tm *TransportManager) Publisher() *transport.Transport {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.publisher
}

func (tm *TransportManager) Subscriber() *transport.Transport {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.subscriber
}
