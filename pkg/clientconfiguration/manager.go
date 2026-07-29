// Package clientconfiguration manages dynamic client SDK features, bitrate caps, and codec compatibility matrix.
package clientconfiguration

import (
	"sync"
)

// ClientConfiguration specifies adaptive client settings.
type ClientConfiguration struct {
	MaxBitrateBps   uint32   `json:"max_bitrate_bps"`
	SupportedCodecs []string `json:"supported_codecs"`
	EnableSimulcast bool     `json:"enable_simulcast"`
}

// Manager dynamically resolves optimal client configuration overrides.
type Manager struct {
	mu            sync.RWMutex
	defaultConfig ClientConfiguration
}

func NewManager() *Manager {
	return &Manager{
		defaultConfig: ClientConfiguration{
			MaxBitrateBps:   3_000_000,
			SupportedCodecs: []string{"opus", "vp8", "h264", "vp9", "av1"},
			EnableSimulcast: true,
		},
	}
}

func (m *Manager) GetConfigurationForClient(_ string) ClientConfiguration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.defaultConfig
}
