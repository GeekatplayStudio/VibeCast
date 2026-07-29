// Package agent provides real-time multimodal audio frame pipelines and voice activity end-pointing.
package agent

import (
	"context"
	"math"
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// FrameHandler processes raw PCM-16 audio frames from WebRTC subscriber tracks.
type FrameHandler func(samples []int16, dbov float64)

// VoicePipeline orchestrates real-time audio frame ingestion, VAD speech end-pointing, and LLM stream dispatch.
type VoicePipeline struct {
	mu             sync.RWMutex
	sampleRate     int
	channels       int
	vadThreshold   float64
	isSpeaking     bool
	silenceFrames  int
	speechFrames   int
	onSpeechStart  func()
	onSpeechEnd    func()
	onFrame        FrameHandler
}

// NewVoicePipeline creates a new real-time voice pipeline.
func NewVoicePipeline(sampleRate, channels int, vadThresholdDbov float64) *VoicePipeline {
	if vadThresholdDbov == 0 {
		vadThresholdDbov = -38.0
	}
	return &VoicePipeline{
		sampleRate:   sampleRate,
		channels:     channels,
		vadThreshold: vadThresholdDbov,
	}
}

// OnSpeechStart sets callback for when user begins speaking.
func (p *VoicePipeline) OnSpeechStart(fn func()) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onSpeechStart = fn
}

// OnSpeechEnd sets callback for when user finishes speaking (speech end-pointing).
func (p *VoicePipeline) OnSpeechEnd(fn func()) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onSpeechEnd = fn
}

// OnFrame sets callback for every ingested audio frame.
func (p *VoicePipeline) OnFrame(fn FrameHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onFrame = fn
}

// PushPCM16Frame ingests a chunk of 16-bit signed PCM audio samples.
func (p *VoicePipeline) PushPCM16Frame(samples []int16) {
	if len(samples) == 0 {
		return
	}

	// Calculate RMS & dBov
	var sumSquares float64
	for _, sample := range samples {
		normalized := float64(sample) / 32768.0
		sumSquares += normalized * normalized
	}
	rms := math.Sqrt(sumSquares / float64(len(samples)))
	dbov := -127.0
	if rms > 1e-6 {
		dbov = 20.0 * math.Log10(rms)
	}

	p.mu.Lock()
	if p.onFrame != nil {
		p.onFrame(samples, dbov)
	}

	// Evaluate VAD Speech End-pointing
	if dbov >= p.vadThreshold {
		p.speechFrames++
		p.silenceFrames = 0
		if !p.isSpeaking && p.speechFrames >= 3 {
			p.isSpeaking = true
			telemetry.Log.Debug().Float64("dbov", dbov).Msg("VoicePipeline: Speech started")
			if p.onSpeechStart != nil {
				go p.onSpeechStart()
			}
		}
	} else {
		p.silenceFrames++
		p.speechFrames = 0
		if p.isSpeaking && p.silenceFrames >= 15 { // ~300ms silence threshold
			p.isSpeaking = false
			telemetry.Log.Debug().Msg("VoicePipeline: Speech ended (end-pointing triggered)")
			if p.onSpeechEnd != nil {
				go p.onSpeechEnd()
			}
		}
	}
	p.mu.Unlock()
}

// Start begins processing context.
func (p *VoicePipeline) Start(ctx context.Context) {
	telemetry.Log.Info().Int("sample_rate", p.sampleRate).Msg("VoicePipeline audio engine active")
	<-ctx.Done()
}
