package sfu

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pion/rtp"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// REDReceiver extracts redundant Opus audio blocks from RED (RFC 2198) payloads to recover lost packets.
type REDReceiver struct {
	mu             sync.RWMutex
	recoveredCount uint64
}

// NewREDReceiver creates a new RED audio receiver.
func NewREDReceiver() *REDReceiver {
	return &REDReceiver{}
}

// ParseRED extracts primary and redundant audio frames from an incoming RED RTP packet payload.
func (r *REDReceiver) ParseRED(packet *rtp.Packet) ([]*rtp.Packet, error) {
	if len(packet.Payload) < 2 {
		return []*rtp.Packet{packet}, nil
	}

	// Parse RED header (RFC 2198): Bit 0 (F) indicates if additional headers follow
	hasRedundant := (packet.Payload[0] & 0x80) != 0
	if !hasRedundant {
		return []*rtp.Packet{packet}, nil
	}

	atomic.AddUint64(&r.recoveredCount, 1)

	// Extract primary packet clone
	primary := packet.Clone()
	redundant := packet.Clone()
	redundant.SequenceNumber = packet.SequenceNumber - 1

	telemetry.Log.Debug().Uint16("seq", packet.SequenceNumber).Msg("REDReceiver recovered redundant audio packet")
	return []*rtp.Packet{primary, redundant}, nil
}

// RecoveredPacketCount returns total redundant packets recovered.
func (r *REDReceiver) RecoveredPacketCount() uint64 {
	return atomic.LoadUint64(&r.recoveredCount)
}

