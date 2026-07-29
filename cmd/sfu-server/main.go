// Copyright 2026 Geekatplay Studio (Vladimir Chopine).
// Based on original concepts by LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/agenticsfu/agentic-sfu/pkg/cluster"
	"github.com/agenticsfu/agentic-sfu/pkg/config"
	"github.com/agenticsfu/agentic-sfu/pkg/mcp"
	"github.com/agenticsfu/agentic-sfu/pkg/service"
	sig "github.com/agenticsfu/agentic-sfu/pkg/signal"
	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

func main() {
	telemetry.Log.Info().Msg("Starting VibeCast Enterprise Media Server...")

	// 1. Load Configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		telemetry.Log.Warn().Err(err).Msg("Could not load config.yaml, using defaults")
		cfg = &config.Config{}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize Cluster Router
	router := cluster.NewMemoryRouter()
	localNode := &cluster.Node{
		ID:      "node-local-1",
		Address: cfg.Server.Host,
		Port:    cfg.Server.Port,
	}
	if err := router.RegisterNode(ctx, localNode); err != nil {
		telemetry.Log.Fatal().Err(err).Msg("Failed to register node")
	}

	// 3. Initialize Signal Server, Admin REST API & Prometheus Metrics
	signalServer := sig.NewServer("devkey")
	http.HandleFunc("/rtc/validate", signalServer.HandleValidate)
	http.HandleFunc("/metrics", telemetry.GlobalCollector.MetricsHandler)

	roomStore := service.NewLocalRoomStore()
	roomManager := service.NewRoomManager(roomStore)
	restAPI := service.NewRestAPIService(roomManager, "devkey", "secret123")
	restAPI.RegisterRoutes(http.DefaultServeMux)

	// 4. Mount React UI Static Dashboard Handler (Single-Binary Distribution)
	uiDir := "./ui/dist"
	if _, err := os.Stat(uiDir); err == nil {
		fs := http.FileServer(http.Dir(uiDir))
		http.Handle("/", fs)
		telemetry.Log.Info().Str("ui_dir", uiDir).Msg("Serving Enterprise Web UI Dashboard at http://localhost:7880")
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`<!DOCTYPE html><html><head><title>AgenticSFU Enterprise</title></head><body style="background:#0a0d14;color:#f3f4f6;font-family:sans-serif;text-align:center;padding:50px;"><h1>AgenticSFU Core Media Server Active</h1><p>MCP Agent Protocol active on :8080 | Metrics active on :7880/metrics</p></body></html>`))
		})
	}

	// 5. Initialize MCP Agent Bridge
	if cfg.MCP.Enabled {
		mcpBridge := mcp.NewBridge(8080)
		go func() {
			if err := mcpBridge.Start(ctx); err != nil {
				telemetry.Log.Error().Err(err).Msg("MCP Bridge error")
			}
		}()
	}

	// 6. Start Signal HTTP Server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	telemetry.Log.Info().Str("address", serverAddr).Msg("Signal server listening")

	go func() {
		if err := http.ListenAndServe(serverAddr, nil); err != nil {
			telemetry.Log.Error().Err(err).Msg("Signal server stopped")
		}
	}()

	// 7. Graceful Shutdown Signal Handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	telemetry.Log.Info().Msg("Shutting down AgenticSFU Media Server...")
}

