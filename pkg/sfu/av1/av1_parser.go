// Package av1 provides AV1 OBU (Open Bitstream Unit) depacketization and temporal layer filtering.
package av1

import (
	"fmt"
	"sync"

	"github.com/pion/rtp"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// OBUType represents AV1 Open Bitstream Unit header type.
type OBUType uint8

const (
	OBUTypeSequenceHeader OBUType = 1
	OBUTypeTemporalDelimiter OBUType = 2
	OBUTypeFrameHeader    OBUType = 3
	OBUTypeTileGroup      OBUType = 4
	OBUTypeMetadata       OBUType = 5
	OBUTypeFrame          OBUType = 6
)

// AV1Parser depacketizes AV1 RTP payloads for selective forwarding.
type AV1Parser struct {
	mu           sync.RWMutex
	keyframeCount uint64
}

// NewAV1Parser creates a new AV1 codec packet parser.
func NewAV1Parser() *AV1Parser {
	return &AV1Parser{}
}

// ParseRTP inspects AV1 RTP packet headers and extracts OBU metadata.
func (p *AV1Parser) ParseRTP(packet *rtp.Packet) (isKeyframe bool, obuType OBUType, err error) {
	if len(packet.Payload) < 2 {
		return false, 0, fmt.Errorf("av1 payload too short")
	}

	// Extract OBU type from first byte (bits 3..6)
	headerByte := packet.Payload[0]
	obuType = OBUType((headerByte >> 3) & 0x0F)

	if obuType == OBUTypeSequenceHeader || obuType == OBUTypeFrame {
		p.mu.Lock()
		p.keyframeCount++
		p.mu.Unlock()
		isKeyframe = true
		telemetry.Log.Debug().Uint16("seq", packet.SequenceNumber).Msg("AV1Parser detected AV1 Keyframe OBU")
	}

	return isKeyframe, obuType, nil
}
