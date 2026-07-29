// Command agent-worker runs a standalone AI agent worker daemon.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/agent"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

func main() {
	telemetry.InitLogger("info", "json")
	telemetry.Log.Info().Msg("Starting VibeCast Standalone AI Agent Worker Process...")

	workerID := "agent-worker-node-1"
	worker := agent.NewAgentWorker(workerID)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Simulate receiving a room job
	job := &agent.AgentJob{
		ID:        "job-init-demo",
		AgentName: "VoiceAssistant",
		RoomName:  "enterprise-studio-room",
		Prompt:    "Act as real-time audio translator & conversation summary assistant",
		Model:     "gemini-3.6-flash",
		State:     agent.JobStateConnected,
		CreatedAt: time.Now(),
	}

	if err := worker.ExecuteJob(ctx, job); err != nil {
		telemetry.Log.Error().Err(err).Msg("Failed to execute agent job")
		os.Exit(1)
	}

	telemetry.Log.Info().Str("worker_id", workerID).Msg("Agent Worker Process running. Press Ctrl+C to terminate.")
	<-ctx.Done()

	worker.Stop()
	telemetry.Log.Info().Msg("Agent Worker Process shut down gracefully.")
}
