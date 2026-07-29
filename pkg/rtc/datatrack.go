package rtc

import (
	"sync"
)

// DataPacket represents a custom text or binary message sent over WebRTC DataChannel.
type DataPacket struct {
	Kind        string `json:"kind"` // "reliable" or "lossy"
	Destination []string `json:"destination,omitempty"`
	Payload     []byte `json:"payload"`
}

// DataTrackManager handles real-time DataChannel message distribution across participants.
type DataTrackManager struct {
	mu           sync.RWMutex
	subscribers  map[string]chan *DataPacket
}

func NewDataTrackManager() *DataTrackManager {
	return &DataTrackManager{
		subscribers: make(map[string]chan *DataPacket),
	}
}

func (d *DataTrackManager) RegisterSubscriber(participantID string) chan *DataPacket {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan *DataPacket, 100)
	d.subscribers[participantID] = ch
	return ch
}

func (d *DataTrackManager) UnregisterSubscriber(participantID string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if ch, ok := d.subscribers[participantID]; ok {
		close(ch)
		delete(d.subscribers, participantID)
	}
}

func (d *DataTrackManager) BroadcastData(packet *DataPacket) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, ch := range d.subscribers {
		select {
		case ch <- packet:
		default:
			// Non-blocking drop if buffer full
		}
	}
}
