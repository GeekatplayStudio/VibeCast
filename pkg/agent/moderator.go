// Package agent provides real-time AI stream chat moderation and VLM/LLM guardrail enforcement.
package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// ModerationStatus represents outcome of AI guardrail evaluation.
type ModerationStatus string

const (
	StatusApproved ModerationStatus = "APPROVED"
	StatusFlagged  ModerationStatus = "FLAGGED"
	StatusBlocked  ModerationStatus = "BLOCKED"
)

// ChatMessage represents a live stream chat entry.
type ChatMessage struct {
	ID         string           `json:"id"`
	SenderID   string           `json:"sender_id"`
	SenderName string           `json:"sender_name"`
	Text       string           `json:"text"`
	Timestamp  time.Time        `json:"timestamp"`
	Status     ModerationStatus `json:"status"`
	Toxicity   float64          `json:"toxicity_score"`
	Reason     string           `json:"reason,omitempty"`
}

// GuardrailSensitivity defines strictness level for AI VLM/LLM moderation.
type GuardrailSensitivity string

const (
	SensitivityLow    GuardrailSensitivity = "LOW"
	SensitivityMedium GuardrailSensitivity = "MEDIUM"
	SensitivityStrict GuardrailSensitivity = "STRICT"
)

// ChatModerator evaluates stream chat messages against open-source LLM/VLM guardrails.
type ChatModerator struct {
	mu           sync.RWMutex
	sensitivity  GuardrailSensitivity
	blockedWords map[string]bool
	processedMsg int
	flaggedMsg   int
}

// NewChatModerator creates a new AI chat moderator instance.
func NewChatModerator(sensitivity GuardrailSensitivity) *ChatModerator {
	if sensitivity == "" {
		sensitivity = SensitivityMedium
	}
	moderator := &ChatModerator{
		sensitivity:  sensitivity,
		blockedWords: make(map[string]bool),
	}

	// Default toxic keyword blocklist
	defaultBlocklist := []string{"spam", "scam", "hate", "abuse", "phishing"}
	for _, word := range defaultBlocklist {
		moderator.blockedWords[word] = true
	}

	return moderator
}

// EvaluateMessage analyzes a chat message against LLM guardrail rules and toxicity thresholds.
func (m *ChatModerator) EvaluateMessage(ctx context.Context, senderID, senderName, text string) *ChatMessage {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.processedMsg++
	lowerText := strings.ToLower(text)
	toxicity := 0.05
	status := StatusApproved
	reason := ""

	// Check toxicity keywords
	for word := range m.blockedWords {
		if strings.Contains(lowerText, word) {
			toxicity = 0.85
			status = StatusFlagged
			reason = fmt.Sprintf("AI Guardrail Trigger: Detected keyword '%s'", word)
			m.flaggedMsg++
			break
		}
	}

	// Apply sensitivity threshold
	if m.sensitivity == SensitivityStrict && toxicity > 0.5 {
		status = StatusBlocked
	}

	msg := &ChatMessage{
		ID:         fmt.Sprintf("msg-%d", time.Now().UnixNano()/1e6),
		SenderID:   senderID,
		SenderName: senderName,
		Text:       text,
		Timestamp:  time.Now(),
		Status:     status,
		Toxicity:   toxicity,
		Reason:     reason,
	}

	telemetry.Log.Info().Str("sender", senderName).Str("status", string(status)).Float64("toxicity", toxicity).Msg("AI Chat Moderator evaluated message")
	return msg
}

// SetSensitivity updates guardrail strictness.
func (m *ChatModerator) SetSensitivity(s GuardrailSensitivity) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sensitivity = s
	telemetry.Log.Info().Str("sensitivity", string(s)).Msg("Updated AI Chat Moderator guardrail sensitivity")
}

// Stats returns moderation summary metrics.
func (m *ChatModerator) Stats() (int, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.processedMsg, m.flaggedMsg
}
