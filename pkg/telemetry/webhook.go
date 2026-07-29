// Package telemetry provides outbound signed HTTP Webhook event dispatching.
package telemetry

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// WebhookEvent represents an outbound notification event payload.
type WebhookEvent struct {
	ID        string      `json:"id"`
	Event     string      `json:"event"` // room_started, participant_joined, agent_dispatched, recording_finished
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// WebhookNotifier manages async HTTP POST event webhooks to creator endpoints.
type WebhookNotifier struct {
	mu         sync.RWMutex
	urls       []string
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

// NewWebhookNotifier creates a new Webhook event dispatcher.
func NewWebhookNotifier(urls []string, apiKey, apiSecret string) *WebhookNotifier {
	return &WebhookNotifier{
		urls:      urls,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// EmitEvent dispatches a signed JSON webhook payload to configured URLs.
func (w *WebhookNotifier) EmitEvent(eventType string, data interface{}) {
	w.mu.RLock()
	urls := make([]string, len(w.urls))
	copy(urls, w.urls)
	w.mu.RUnlock()

	if len(urls) == 0 {
		return
	}

	event := WebhookEvent{
		ID:        fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		Event:     eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return
	}

	// Compute HMAC-SHA256 signature header
	h := hmac.New(sha256.New, []byte(w.apiSecret))
	h.Write(payloadBytes)
	signature := hex.EncodeToString(h.Sum(nil))

	for _, url := range urls {
		go func(targetURL string) {
			req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(payloadBytes))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))
			req.Header.Set("X-AgenticSFU-Signature", signature)

			resp, err := w.httpClient.Do(req)
			if err == nil && resp != nil {
				_ = resp.Body.Close()
			}
		}(url)
	}
}
