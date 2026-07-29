# VibeCast Enterprise

> **High-Performance WebRTC SFU, Native Model Context Protocol (MCP) AI Agent Engine & Producer Studio**
> Created by **Geekatplay Studio (Vladimir Chopine)** under the Apache 2.0 Open Source License.

---

## 🙏 Credits & Inspiration

**VibeCast Enterprise** was deeply inspired by and based on the pioneering concepts of [LiveKit](https://github.com/livekit/livekit). 

While honoring LiveKit's original vision, **VibeCast** represents a complete clean-room rewrite designed specifically for content creators, live streamers, and AI developers:
- 🛡️ **Hardened Security & Vulnerability Fixes**: Completely audited and rewritten to eliminate memory/concurrency leaks, enforce strict HMAC-SHA256 JWT access token validation, and secure all signaling/REST endpoints.
- 🎨 **Superior Producer Studio UI**: Features a state-of-the-art cyber-glassmorphic control panel with real-time VAD decibel wave meters, MOS quality badges, active participant controls, and 7-language i18n localization.
- 🤖 **Native Model Context Protocol (MCP) Integration**: Built-in MCP Server (`:8080`) & WebRTC SCTP DataChannel bridge for zero-latency LLM agent tool discovery and execution.
- 💬 **AI Stream Chat & VLM/LLM Moderation**: Built-in live stream chat with open-source Vision-Language Model (VLM) & LLM guardrails that monitor chat in real time, auto-enforce toxicity thresholds, and protect stream producers.

---

## 🌟 Key Features

- 🚀 **High-Throughput WebRTC Media Router**: Zero-allocation RTP packet forwarding, sequence remapping, and R-Factor MOS connection quality scoring.
- 🎬 **Producer Studio Suite**: Dedicated creator control panel with AI VLM video feed monitoring, stream chat moderation, guardrail sensitivity controls, and overlay triggers.
- 🤖 **Native Model Context Protocol (MCP) Server**: Dual-mode MCP agent tool bridge (HTTP REST `:8080` & sub-5ms WebRTC SCTP DataChannels).
- 💬 **AI Chat Moderator & VLM Guardrails**: Real-time open-source LLM/VLM chat monitoring that automatically flags toxicity, blocks spam, and enforces producer rules.
- 🎙️ **Zero-Latency AI Voice Assistant Pipeline**: Direct PCM-16/Opus audio frame ingestion with VAD speech end-pointing.
- 📹 **Media Egress & Recording**: MP4/WebM room recorder, RTMP live streaming to YouTube/Twitch, and async S3 cloud storage export (`cmd/egress-server`).
- 📡 **Media Ingress Stream Ingest**: IETF WHIP protocol endpoint (`:8088`) and RTSP IP Camera stream proxy (`cmd/ingress-server`).
- 📞 **SIP / PSTN Telephony Gateway**: UDP `:5060` gateway allowing PSTN telephone callers to dial into WebRTC rooms (`cmd/sip-server`).
- 🌐 **7-Language (i18n) Localization**: Native UI support for 🇺🇸 EN, 🇪🇸 ES, 🇫🇷 FR, 🇩🇪 DE, 🇯🇵 JA, 🇨🇳 ZH, and 🇷🇺 RU.

---

## 📦 Microservice Daemons

| Daemon | Command Path | Primary Port | Description |
| :--- | :--- | :--- | :--- |
| **sfu-server** | `cmd/sfu-server` | `:7880` / `:8080` | Core SFU Media Engine, Signal Server, Admin REST API, & MCP Bridge |
| **agent-worker** | `cmd/agent-worker` | Daemon | Autoscale AI Agent Worker Process & Chat Moderator |
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
Navigate to `http://localhost:7880` in your web browser to open the Producer Studio.

### 3. Multi-Container Docker Compose Orchestration
```bash
docker-compose up --build
```

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
