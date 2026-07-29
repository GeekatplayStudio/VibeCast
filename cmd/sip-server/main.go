// Command sip-server runs standalone SIP / PSTN Telephony Gateway Daemon.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/agenticsfu/agentic-sfu/pkg/sip"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

func main() {
	telemetry.InitLogger("info", "json")
	telemetry.Log.Info().Msg("Starting VibeCast Standalone SIP/PSTN Telephony Gateway...")

	bridge := sip.NewSIPBridge(5060)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := bridge.Start(ctx); err != nil {
			telemetry.Log.Error().Err(err).Msg("SIP Bridge listener error")
		}
	}()

	telemetry.Log.Info().Msg("SIP Gateway Server running on UDP :5060. Press Ctrl+C to terminate.")
	<-ctx.Done()

	telemetry.Log.Info().Msg("SIP Gateway Server shut down gracefully.")
}
