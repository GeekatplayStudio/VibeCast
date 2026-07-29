// Package io manages external media streaming adapters: RTMP/WHIP ingress, composite egress recording, and SIP gateways.
package io

import (
	"context"
	"fmt"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// IngressService handles incoming external live streams (RTMP, WHIP, OBS Studio).
type IngressService struct{}

func NewIngressService() *IngressService {
	return &IngressService{}
}

func (s *IngressService) CreateIngress(ctx context.Context, roomName string, inputType string) (string, error) {
	streamID := fmt.Sprintf("ing-%s-%s", inputType, roomName)
	telemetry.Log.Info().Str("room", roomName).Str("type", inputType).Str("stream_id", streamID).Msg("Created stream Ingress endpoint")
	return streamID, nil
}
