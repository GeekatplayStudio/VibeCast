# VibeCast Enterprise

> **High-Performance WebRTC SFU & Native Model Context Protocol (MCP) AI Agent Engine**
> Created by **Geekatplay Studio (Vladimir Chopine)** under the Apache 2.0 Open Source License.

---

## 🌟 Key Features

- 🚀 **High-Throughput WebRTC Media Router**: Zero-allocation RTP packet forwarding, sequence remapping, and R-Factor MOS connection quality scoring.
- 🤖 **Native Model Context Protocol (MCP) Server**: Dual-mode MCP agent tool bridge (HTTP REST `:8080` & sub-5ms WebRTC SCTP DataChannels).
- 🎙️ **Zero-Latency AI Voice Assistant Pipeline**: Direct PCM-16/Opus audio frame ingestion with VAD speech end-pointing.
- 📹 **Media Egress & Recording**: MP4/WebM room recorder, RTMP live streaming to YouTube/Twitch, and async S3 cloud storage export (`cmd/egress-server`).
- 📡 **Media Ingress Stream Ingest**: IETF WHIP protocol endpoint (`:8088`) and RTSP IP Camera stream proxy (`cmd/ingress-server`).
- 📞 **SIP / PSTN Telephony Gateway**: UDP `:5060` gateway allowing PSTN telephone callers to dial into WebRTC rooms (`cmd/sip-server`).
- 🎨 **Superior Cyber-Glassmorphism Web UI**: Interactive dashboard with real-time VAD decibel wave meters, MOS quality badges, participant mute controls, MCP prompt console with LLM model selection (Gemini 3.6 Flash, Claude 3.5, GPT-4o), 7-language i18n localization, and live telemetry.
- 📊 **Prometheus & Webhook Engine**: OpenMetrics `/metrics` endpoint and signed HMAC-SHA256 HTTP Webhook event dispatching.

---

## 📦 Microservice Daemons

| Daemon | Command Path | Primary Port | Description |
| :--- | :--- | :--- | :--- |
| **sfu-server** | `cmd/sfu-server` | `:7880` / `:8080` | Core SFU Media Engine, Signal Server, Admin REST API, & MCP Bridge |
| **agent-worker** | `cmd/agent-worker` | Daemon | Autoscale AI Agent Worker Process Runner |
| **egress-server** | `cmd/egress-server` | Daemon | Room Media Recording, RTMP Stream Push, & S3 Uploader |
| **ingress-server** | `cmd/ingress-server` | `:8088` | WHIP WebRTC Ingest & RTSP IP Camera Proxy |
| **sip-server** | `cmd/sip-server` | `:5060` (UDP) | SIP / PSTN Telephone Gateway Bridge |

---

## 🛠️ Quick Start

### 1. Build All Microservices (Windows / PowerShell)
```powershell
.\scripts\build.ps1
```

### 2. Run Single-Binary Server (Includes React UI & MCP Bridge)
```powershell
.\bin\sfu-server.exe
```
Navigate to `http://localhost:7880` in your web browser.

### 3. Multi-Container Docker Compose Orchestration
```bash
docker-compose up --build
```

---

## 🌐 Multi-Language (i18n) Support

Supported in **7 Languages**:
- 🇺🇸 **English (`en`)**
- 🇪🇸 **Spanish (`es`)**
- 🇫🇷 **French (`fr`)**
- 🇩🇪 **German (`de`)**
- 🇯🇵 **Japanese (`ja`)**
- 🇨🇳 **Chinese (`zh`)**
- 🇷🇺 **Russian (`ru`)**

---

## 🧪 Testing

Run the full end-to-end integration test suite:
```powershell
go test ./pkg/... ./test/...
```

Build the React UI production bundle:
```powershell
cd ui
npm run build
```

---

## 📜 License

Licensed under the **Apache License, Version 2.0**. Free for open-source use, creators, and enterprise deployments.
