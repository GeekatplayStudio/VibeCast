// Package config provides type-safe configuration loading from YAML and environment variables
// for the AgenticSFU media server runtime.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config defines the complete application configuration structure.
type Config struct {
	Server ServerConfig `yaml:"server"`
	Router RouterConfig `yaml:"router"`
	Auth   AuthConfig   `yaml:"auth"`
	MCP    MCPConfig    `yaml:"mcp"`
}

// ServerConfig holds HTTP, WebSocket, and WebRTC network settings.
type ServerConfig struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	RTCPortRangeStart uint16 `yaml:"rtc_port_range_start"`
	RTCPortRangeEnd   uint16 `yaml:"rtc_port_range_end"`
}

// RouterConfig specifies multi-node routing type (memory or redis).
type RouterConfig struct {
	Type     string `yaml:"type"`
	RedisURL string `yaml:"redis_url"`
}

// AuthConfig manages API key and secret pairings for JWT validation.
type AuthConfig struct {
	Keys map[string]string `yaml:"keys"`
}

// MCPConfig controls Model Context Protocol server settings for AI agents.
type MCPConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

// LoadConfig reads and parses a YAML configuration file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config: %w", err)
	}

	cfg.applyDefaults()
	return cfg, nil
}

func (c *Config) applyDefaults() {
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 7880
	}
	if c.Server.RTCPortRangeStart == 0 {
		c.Server.RTCPortRangeStart = 50000
	}
	if c.Server.RTCPortRangeEnd == 0 {
		c.Server.RTCPortRangeEnd = 60000
	}
	if c.Router.Type == "" {
		c.Router.Type = "memory"
	}
}
