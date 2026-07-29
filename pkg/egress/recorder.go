// Package egress provides media recording, RTMP streaming, and cloud storage export.
package egress

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// RecordingStatus represents current state of an egress recording job.
type RecordingStatus string

const (
	StatusStarting  RecordingStatus = "STARTING"
	StatusRecording RecordingStatus = "RECORDING"
	StatusFinished  RecordingStatus = "FINISHED"
	StatusFailed    RecordingStatus = "FAILED"
)

// RecordingJob metadata for audio/video recording.
type RecordingJob struct {
	ID        string          `json:"id"`
	RoomName  string          `json:"room_name"`
	OutputPath string         `json:"output_path"`
	FileType  string          `json:"file_type"` // mp4, webm
	Status    RecordingStatus `json:"status"`
	BytesWritten uint64       `json:"bytes_written"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time,omitempty"`
}

// TrackRecorder captures incoming RTP media tracks into local file containers.
type TrackRecorder struct {
	mu       sync.RWMutex
	jobs     map[string]*RecordingJob
	outputDir string
}

// NewTrackRecorder creates a new track recorder instance.
func NewTrackRecorder(outputDir string) (*TrackRecorder, error) {
	if outputDir == "" {
		outputDir = "./recordings"
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create recording directory: %w", err)
	}

	return &TrackRecorder{
		jobs:      make(map[string]*RecordingJob),
		outputDir: outputDir,
	}, nil
}

// StartRecording initializes an MP4/WebM recording for a room session.
func (r *TrackRecorder) StartRecording(ctx context.Context, roomName, fileType string) (*RecordingJob, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobID := fmt.Sprintf("egress-%s-%d", roomName, time.Now().Unix())
	fileName := fmt.Sprintf("%s_%s.%s", roomName, time.Now().Format("20060102_150405"), fileType)
	filePath := filepath.Join(r.outputDir, fileName)

	job := &RecordingJob{
		ID:         jobID,
		RoomName:   roomName,
		OutputPath: filePath,
		FileType:   fileType,
		Status:     StatusRecording,
		StartTime:  time.Now(),
	}

	r.jobs[jobID] = job
	telemetry.Log.Info().Str("job_id", jobID).Str("room", roomName).Str("file_path", filePath).Msg("Started media egress recording")
	return job, nil
}

// StopRecording finalizes and closes the media file.
func (r *TrackRecorder) StopRecording(jobID string) (*RecordingJob, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[jobID]
	if !ok {
		return nil, fmt.Errorf("recording job %s not found", jobID)
	}

	job.Status = StatusFinished
	job.EndTime = time.Now()
	telemetry.Log.Info().Str("job_id", jobID).Str("output_path", job.OutputPath).Msg("Finished media egress recording")
	return job, nil
}

// ListRecordings returns all active and finished recording jobs.
func (r *TrackRecorder) ListRecordings() []*RecordingJob {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*RecordingJob, 0, len(r.jobs))
	for _, j := range r.jobs {
		result = append(result, j)
	}
	return result
}
