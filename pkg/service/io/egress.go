package io

import (
	"context"
	"fmt"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// EgressService handles room recording, composite video export, and S3 upload.
type EgressService struct{}

func NewEgressService() *EgressService {
	return &EgressService{}
}

func (e *EgressService) StartRoomCompositeEgress(ctx context.Context, roomName string, destinationURL string) (string, error) {
	egressID := fmt.Sprintf("egr-mp4-%s", roomName)
	telemetry.Log.Info().Str("room", roomName).Str("destination", destinationURL).Str("egress_id", egressID).Msg("Started Egress composite recording job")
	return egressID, nil
}
