// Command egress-server launches standalone Media Egress Recording and Streaming daemon.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/agenticsfu/agentic-sfu/pkg/egress"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

func main() {
	telemetry.InitLogger("info", "json")
	telemetry.Log.Info().Msg("Starting VibeCast Standalone Egress Server Daemon...")

	recorder, err := egress.NewTrackRecorder("./recordings")
	if err != nil {
		telemetry.Log.Error().Err(err).Msg("Failed to initialize track recorder")
		os.Exit(1)
	}

	pusher := egress.NewRTMPPusher()
	uploader := egress.NewCloudStorageUploader(egress.StorageS3, "agentic-sfu-recordings")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Simulate recording trigger
	job, err := recorder.StartRecording(ctx, "enterprise-studio-room", "mp4")
	if err == nil {
		telemetry.Log.Info().Str("job_id", job.ID).Msg("Egress Recorder session active")
	}

	_ = pusher
	_ = uploader

	telemetry.Log.Info().Msg("Egress Server running. Press Ctrl+C to terminate.")
	<-ctx.Done()

	if job != nil {
		_, _ = recorder.StopRecording(job.ID)
	}

	telemetry.Log.Info().Msg("Egress Server shut down gracefully.")
}
