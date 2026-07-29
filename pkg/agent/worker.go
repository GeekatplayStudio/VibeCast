// Package agent provides worker execution harness for AI agent room participation.
package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// AgentWorker manages a single agent participant worker process.
type AgentWorker struct {
	mu          sync.RWMutex
	workerID    string
	currentJob  *AgentJob
	pipeline    *VoicePipeline
	cancelFunc  context.CancelFunc
	activeRooms map[string]bool
}

// NewAgentWorker initializes an AI agent worker.
func NewAgentWorker(workerID string) *AgentWorker {
	return &AgentWorker{
		workerID:    workerID,
		pipeline:    NewVoicePipeline(16000, 1, -38.0),
		activeRooms: make(map[string]bool),
	}
}

// ExecuteJob joins the room specified in the AgentJob and launches the VoicePipeline.
func (w *AgentWorker) ExecuteJob(ctx context.Context, job *AgentJob) error {
	w.mu.Lock()
	w.currentJob = job
	jobCtx, cancel := context.WithCancel(ctx)
	w.cancelFunc = cancel
	w.activeRooms[job.RoomName] = true
	w.mu.Unlock()

	telemetry.Log.Info().Str("worker_id", w.workerID).Str("job_id", job.ID).Str("room", job.RoomName).Str("model", job.Model).Msg("AgentWorker starting job execution")

	// Set up voice pipeline callbacks
	w.pipeline.OnSpeechStart(func() {
		telemetry.Log.Info().Str("job_id", job.ID).Msg("AgentWorker: User is speaking...")
	})

	w.pipeline.OnSpeechEnd(func() {
		telemetry.Log.Info().Str("job_id", job.ID).Msg("AgentWorker: User speech ended -> Generating LLM audio response")
	})

	go w.pipeline.Start(jobCtx)

	return nil
}

// Stop terminates the current running job.
func (w *AgentWorker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cancelFunc != nil {
		w.cancelFunc()
	}
	if w.currentJob != nil {
		w.currentJob.State = JobStateTerminated
		telemetry.Log.Info().Str("job_id", w.currentJob.ID).Msg("AgentWorker job stopped successfully")
	}
}

// Status returns current worker metadata.
func (w *AgentWorker) Status() map[string]interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var activeJobID string
	if w.currentJob != nil {
		activeJobID = w.currentJob.ID
	}

	return map[string]interface{}{
		"worker_id":     w.workerID,
		"active_job":    activeJobID,
		"active_rooms":  len(w.activeRooms),
		"uptime_seconds": time.Now().Unix(),
	}
}
