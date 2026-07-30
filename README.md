# VibeCast

> **High-Performance WebRTC SFU, Native Model Context Protocol (MCP) AI Agent Engine & Producer Studio**
> Created by **Geekatplay Studio (Vladimir Chopine)** under the Apache 2.0 Open Source License.

---

## 🙏 Credits & Attribution

**VibeCast** was inspired by and derived from open WebRTC architecture patterns pioneeringly demonstrated by projects like [LiveKit](https://github.com/livekit/livekit). 

While building upon these foundational principles, **VibeCast** is tailored specifically for content creators, live streamers, and AI developers:
- 🛡️ **Security & Access Control**: Enforces strict HMAC-SHA256 JWT access token validation and structured concurrency across signaling and REST endpoints.
- 🎨 **Producer Studio UI**: Features a modern cyber-glassmorphic control panel with real-time VAD decibel wave meters, connection quality metrics, active participant controls, and multi-language (7 i18n) localization.
- 🤖 **Native Model Context Protocol (MCP) Integration**: Built-in MCP Server (`:8080`) & WebRTC SCTP DataChannel bridge for AI agent tool discovery and execution.
- 💬 **AI Stream Chat & VLM/LLM Moderation**: Built-in live stream chat with open-source Vision-Language Model (VLM) & LLM guardrails that monitor chat in real time, auto-enforce toxicity thresholds, and protect stream producers.

---

## 🏗️ Microservice System Architecture

```
                                ┌──────────────────────────────────┐
                                │          VibeCast Suite          │
                                └────────────────┬─────────────────┘
                                                 │
   ┌────────────────┬─────────────────┬─────────┴────────┬──────────────────┐
   │                │                 │                  │                  │
┌──▼─────────┐ ┌────▼──────────┐ ┌────▼─────────┐ ┌──────▼──────────┐ ┌─────▼──────────┐
│ sfu-server │ │agent-worker   │ │egress-server │ │ingress-server   │ │sip-server     │
│ (:7880)    │ │(Voice & Chat) │ │(MP4/RTMP/S3) │ │(WHIP :8088)     │ │(SIP UDP :5060)│
└──┬─────────┘ └────┬──────────┘ └────┬─────────┘ └──────┬──────────┘ └─────┬──────────┘
   │                │                 │                  │                  │
   └────────────────┴─────────────────┼──────────────────┴──────────────────┘
                                      │
                   ┌──────────────────▼───────────────────┐
                   │   Cyber-Glassmorphic Web UI          │
                   └──────────────────────────────────────┘
```

---

## 🌟 Key Features

- 🚀 **WebRTC Media Router**: Efficient RTP packet forwarding, sequence remapping, and real-time connection quality scoring.
- 🎬 **Producer Studio Suite**: Dedicated creator control panel with AI VLM video feed monitoring, stream chat moderation, guardrail sensitivity controls, and overlay triggers.
- 🤖 **Native Model Context Protocol (MCP) Server**: Dual-mode MCP agent tool bridge (HTTP REST `:8080` & WebRTC SCTP DataChannels).
- 💬 **AI Chat Moderator & VLM Guardrails**: Real-time open-source LLM/VLM chat monitoring that automatically flags toxicity, blocks spam, and enforces producer rules.
- 🎙️ **AI Voice Assistant Pipeline**: Direct PCM-16/Opus audio frame ingestion with VAD speech end-pointing.
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

## 🧪 Testing & Verification

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

## 📜 License & Legal Attribution

Licensed under the **Apache License, Version 2.0**. Free for open-source use, creators, and developers.
Original concepts and architectural inspiration credited to [LiveKit, Inc.](https://github.com/livekit/livekit).
