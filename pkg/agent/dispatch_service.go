// Package agent manages real-time AI multimodal and voice agent participant dispatching.
package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// JobState represents the lifecycle status of an AI agent job.
type JobState string

const (
	JobStateSpawning   JobState = "SPAWNING"
	JobStateConnected  JobState = "CONNECTED"
	JobStateActive     JobState = "ACTIVE"
	JobStateTerminated JobState = "TERMINATED"
	JobStateFailed     JobState = "FAILED"
)

// AgentJob holds metadata and status for a dispatched agent worker task.
type AgentJob struct {
	ID        string    `json:"id"`
	AgentName string    `json:"agent_name"`
	RoomName  string    `json:"room_name"`
	Prompt    string    `json:"prompt"`
	Model     string    `json:"model"`
	State     JobState  `json:"state"`
	WorkerID  string    `json:"worker_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WorkerInfo tracks active worker node availability.
type WorkerInfo struct {
	ID          string    `json:"id"`
	ActiveJobs  int       `json:"active_jobs"`
	MaxJobs     int       `json:"max_jobs"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

// AgentDispatchService coordinates prompt-driven AI agent dispatching and worker process lifecycle.
type AgentDispatchService struct {
	mu      sync.RWMutex
	jobs    map[string]*AgentJob
	workers map[string]*WorkerInfo
}

// NewAgentDispatchService initializes the agent dispatcher.
func NewAgentDispatchService() *AgentDispatchService {
	return &AgentDispatchService{
		jobs:    make(map[string]*AgentJob),
		workers: make(map[string]*WorkerInfo),
	}
}

// RegisterWorker registers an active agent worker node.
func (s *AgentDispatchService) RegisterWorker(workerID string, maxJobs int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.workers[workerID] = &WorkerInfo{
		ID:            workerID,
		MaxJobs:       maxJobs,
		LastHeartbeat: time.Now(),
	}
	telemetry.Log.Info().Str("worker_id", workerID).Int("max_jobs", maxJobs).Msg("Registered AI Agent worker node")
}

// DispatchAgent creates and queues a new agent job for execution.
func (s *AgentDispatchService) DispatchAgent(ctx context.Context, agentName, roomName, prompt, model string) (*AgentJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobID := fmt.Sprintf("job-%s-%d", agentName, time.Now().UnixNano()/1e6)
	now := time.Now()

	job := &AgentJob{
		ID:        jobID,
		AgentName: agentName,
		RoomName:  roomName,
		Prompt:    prompt,
		Model:     model,
		State:     JobStateSpawning,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.jobs[jobID] = job
	telemetry.Log.Info().Str("job_id", jobID).Str("agent_name", agentName).Str("room", roomName).Msg("Dispatched AI Agent worker job")
	return job, nil
}

// UpdateJobState transitions an agent job's status.
func (s *AgentDispatchService) UpdateJobState(jobID string, state JobState, workerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return fmt.Errorf("job %s not found", jobID)
	}

	job.State = state
	if workerID != "" {
		job.WorkerID = workerID
	}
	job.UpdatedAt = time.Now()

	telemetry.Log.Info().Str("job_id", jobID).Str("state", string(state)).Str("worker_id", workerID).Msg("Updated agent job state")
	return nil
}

// ListJobs returns all active agent jobs.
func (s *AgentDispatchService) ListJobs() []*AgentJob {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*AgentJob, 0, len(s.jobs))
	for _, j := range s.jobs {
		result = append(result, j)
	}
	return result
}

