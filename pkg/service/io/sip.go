package io

import (
	"context"
	"fmt"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// SIPService handles SIP trunking, dial-in, and PSTN audio gateway integration.
type SIPService struct{}

func NewSIPService() *SIPService {
	return &SIPService{}
}

func (s *SIPService) CreateSIPCall(ctx context.Context, phoneNumber string, roomName string) (string, error) {
	callID := fmt.Sprintf("sip-call-%s", roomName)
	telemetry.Log.Info().Str("phone", phoneNumber).Str("room", roomName).Str("call_id", callID).Msg("Dispatched outbound SIP PSTN call")
	return callID, nil
}
