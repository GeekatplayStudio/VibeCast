// Command ingress-server runs standalone WHIP and RTSP Stream Ingest Daemon.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/agenticsfu/agentic-sfu/pkg/ingress"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

func main() {
	telemetry.InitLogger("info", "json")
	telemetry.Log.Info().Msg("Starting VibeCast Standalone Ingress Server Daemon...")

	whipServer := ingress.NewWHIPServer(8088)
	rtspProxy := ingress.NewRTSPProxy()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := whipServer.Start(ctx); err != nil {
			telemetry.Log.Error().Err(err).Msg("WHIP Server terminated")
		}
	}()

	_ = rtspProxy

	telemetry.Log.Info().Msg("Ingress Server running on :8088. Press Ctrl+C to terminate.")
	<-ctx.Done()

	telemetry.Log.Info().Msg("Ingress Server shut down gracefully.")
}
