import React, { useState } from 'react';
import { Activity, Bot, Cpu, Globe, Radio, ShieldCheck, Volume2, Sparkles } from 'lucide-react';
import { RoomList, Participant } from './components/RoomList';
import { TelemetryPanel } from './components/TelemetryPanel';
import { AgentConsole } from './components/AgentConsole';
import { ProducerStudio } from './components/ProducerStudio';
import { StreamChat, ChatMsg } from './components/StreamChat';
import { LanguageProvider, useTranslation } from './i18n/LanguageContext';
import { Language } from './i18n/translations';

function DashboardContent() {
  const { language, setLanguage, t } = useTranslation();
  const [roomName, setRoomName] = useState('enterprise-studio-room');
  const [isConnected, setIsConnected] = useState(false);
  const [agentDispatched, setAgentDispatched] = useState(false);
  const [activeModel, setActiveModel] = useState('gemini-3.6-flash');
  const [logs, setLogs] = useState<string[]>([
    '[SYSTEM] VibeCast Enterprise Media Server initialized',
    '[MCP] Model Context Protocol Server active on :8080',
    '[SERVICE] Ingress, Egress, and SIP adapters ready',
    '[ROUTER] Dual-stack Redis & In-Memory Cluster Router ready',
  ]);

  const [participants, setParticipants] = useState<Participant[]>([
    {
      id: 'pub-1',
      name: 'Vladimir Chopine (Host)',
      isMuted: false,
      role: 'publisher',
      qualityScore: 5,
      rttMs: 14,
      dbov: -18.4,
    },
    {
      id: 'agent-mcp-1',
      name: 'AI Realtime Voice Assistant',
      isMuted: false,
      role: 'agent',
      qualityScore: 5,
      rttMs: 8,
      dbov: -22.1,
    },
  ]);

  const [chatMessages, setChatMessages] = useState<ChatMsg[]>([
    {
      id: '1',
      senderName: 'Viewer_Alex',
      text: 'Great live stream quality! Loving the AI assistant.',
      timestamp: '11:45 AM',
      status: 'APPROVED',
      toxicity: 0.02,
    },
    {
      id: '2',
      senderName: 'Producer_Guard',
      text: 'VLM vision moderation active on video feed.',
      timestamp: '11:46 AM',
      status: 'APPROVED',
      toxicity: 0.05,
    },
  ]);

  const handleSendChatMessage = (text: string) => {
    const isToxic = text.toLowerCase().includes('spam') || text.toLowerCase().includes('hate');
    const newMsg: ChatMsg = {
      id: String(Date.now()),
      senderName: 'StreamProducer',
      text: text,
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      status: isToxic ? 'FLAGGED' : 'APPROVED',
      toxicity: isToxic ? 0.88 : 0.04,
    };
    setChatMessages((prev) => [...prev, newMsg]);
    addLog(`[CHAT AI MODERATOR] Evaluated chat message. Status: ${newMsg.status} (Toxicity: ${Math.round(newMsg.toxicity * 100)}%)`);
  };

  const handleSensitivityChange = (sens: string) => {
    addLog(`[PRODUCER STUDIO] Set VLM Guardrail Sensitivity to: ${sens}`);
  };

  const handleToggleOverlay = (overlay: string) => {
    addLog(`[PRODUCER STUDIO] Toggled Stream Overlay: ${overlay}`);
  };

  const handleConnect = () => {
    setIsConnected(!isConnected);
    addLog(isConnected ? 'Disconnected from WebRTC room' : `Connected to WebRTC room: ${roomName}`);
  };

  const handleDispatchAgent = (prompt: string, model: string, tool: string) => {
    setAgentDispatched(true);
    setActiveModel(model);
    addLog(`[MCP] Executed Tool '${tool}' with Model [${model}]`);
    addLog(`[MCP Agent] System Prompt: "${prompt.slice(0, 60)}..."`);
  };

  const handleToggleMute = (id: string) => {
    setParticipants((prev) =>
      prev.map((p) => (p.id === id ? { ...p, isMuted: !p.isMuted } : p))
    );
    addLog(`[ROUTER] Toggled audio mute state for participant: ${id}`);
  };

  const addLog = (msg: string) => {
    setLogs((prev) => [`[${new Date().toLocaleTimeString()}] ${msg}`, ...prev.slice(0, 15)]);
  };

  return (
    <div className="dashboard-container">
      {/* Header Banner */}
      <header className="header-banner">
        <div className="logo-group">
          <div className="logo-badge">V</div>
          <div>
            <h1 className="title-text">{t('title')}</h1>
            <p style={{ fontSize: '0.85rem', color: '#9ca3af' }}>
              {t('subtitle')}
            </p>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', background: 'rgba(255,255,255,0.05)', padding: '0.3rem 0.6rem', borderRadius: '8px', border: '1px solid rgba(255,255,255,0.1)' }}>
            <Globe size={14} color="#00f2fe" />
            <select
              value={language}
              onChange={(e) => setLanguage(e.target.value as Language)}
              style={{
                background: 'transparent',
                border: 'none',
                color: '#fff',
                fontFamily: 'inherit',
                fontSize: '0.8rem',
                cursor: 'pointer',
              }}
            >
              <option value="en" style={{ background: '#0a0d14' }}>🇺🇸 English</option>
              <option value="es" style={{ background: '#0a0d14' }}>🇪🇸 Español</option>
              <option value="fr" style={{ background: '#0a0d14' }}>🇫🇷 Français</option>
              <option value="de" style={{ background: '#0a0d14' }}>🇩🇪 Deutsch</option>
              <option value="ja" style={{ background: '#0a0d14' }}>🇯🇵 日本語</option>
              <option value="zh" style={{ background: '#0a0d14' }}>🇨🇳 中文</option>
              <option value="ru" style={{ background: '#0a0d14' }}>🇷🇺 Русский</option>
            </select>
          </div>

          <div className="badge-mcp">
            <div className="badge-pulse"></div>
            {t('mcpActive')}
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', fontSize: '0.8rem', color: '#4facfe' }}>
            <ShieldCheck size={16} /> {t('enterpriseNode')}
          </div>
        </div>
      </header>

      <header className="header-banner">
        <div className="logo-group">
          <div className="logo-badge">A</div>
          <div>
            <h1 className="title-text">AgenticSFU Enterprise</h1>
            <p style={{ fontSize: '0.85rem', color: '#9ca3af' }}>
              Geekatplay Studio (Vladimir Chopine) • Modern WebRTC SFU & MCP Agent Architecture
            </p>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
          <div className="badge-mcp">
            <div className="badge-pulse"></div>
            MCP Agent Protocol Active (:8080)
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', fontSize: '0.8rem', color: '#4facfe' }}>
            <ShieldCheck size={16} /> Enterprise Node
          </div>
        </div>
      </header>

      {/* Main Grid Layout */}
      <div className="grid-layout">
        {/* Left Column: Stage Preview & Controls */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <div className="card">
            <div className="card-header" style={{ justifyContent: 'space-between' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <Radio size={20} color="#00f2fe" />
                Real-Time Media Router & AI Stage
              </div>
              {isConnected && (
                <span style={{ fontSize: '0.8rem', color: '#00f2fe', display: 'flex', alignItems: 'center', gap: '0.3rem' }}>
                  <Sparkles size={14} /> Active WebRTC Transport
                </span>
              )}
            </div>

            <div className="stage-preview">
              {isConnected ? (
                <>
                  <div className="audio-waves">
                    <div className="wave-bar"></div>
                    <div className="wave-bar"></div>
                    <div className="wave-bar"></div>
                    <div className="wave-bar"></div>
                    <div className="wave-bar"></div>
                  </div>
                  <p style={{ color: '#00f2fe', fontWeight: 600 }}>
                    Active DownTrack Selective Forwarding Unit - Room [{roomName}]
                  </p>
                  {agentDispatched && (
                    <div
                      style={{
                        background: 'rgba(79, 172, 254, 0.15)',
                        border: '1px solid rgba(79, 172, 254, 0.4)',
                        padding: '0.5rem 1rem',
                        borderRadius: '8px',
                        color: '#4facfe',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        fontSize: '0.9rem',
                      }}
                    >
                      <Bot size={16} /> AI Multimodal Voice Agent Connected ({activeModel} • PCM 16kHz)
                    </div>
                  )}
                </>
              ) : (
                <>
                  <Volume2 size={48} color="#4b5563" />
                  <p style={{ color: '#9ca3af' }}>No Active WebRTC PeerConnection</p>
                  <p style={{ color: '#6b7280', fontSize: '0.8rem' }}>Enter a room name below and click 'Join Room' to publish & subscribe</p>
                </>
              )}
            </div>

            <div style={{ display: 'flex', gap: '1rem', marginTop: '1.5rem' }}>
              <input
                type="text"
                value={roomName}
                onChange={(e) => setRoomName(e.target.value)}
                placeholder="Enter Room Name"
                style={{
                  background: 'rgba(255,255,255,0.05)',
                  border: '1px solid rgba(255,255,255,0.1)',
                  padding: '0.75rem 1rem',
                  borderRadius: '10px',
                  color: '#fff',
                  flex: 1,
                  fontFamily: 'inherit',
                }}
              />
              <button className="btn" onClick={handleConnect}>
                {isConnected ? 'Disconnect' : 'Join Room'}
              </button>
            </div>
          </div>

          <ProducerStudio
            onSensitivityChange={handleSensitivityChange}
            onToggleOverlay={handleToggleOverlay}
          />

          <RoomList
            roomName={roomName}
            participants={isConnected ? participants : []}
            onToggleMute={handleToggleMute}
          />
        </div>

        {/* Right Column: Stream Chat, Telemetry, MCP Agent Console & Logs */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <StreamChat
            messages={chatMessages}
            onSendMessage={handleSendChatMessage}
          />
          <TelemetryPanel
            bytesIn={isConnected ? 1485760 : 0}
            bytesOut={isConnected ? 2897152 : 0}
            nackCount={isConnected ? 2 : 0}
            rttMs={14}
            packetLoss={0.1}
            activeConnections={isConnected ? participants.length : 0}
          />

          <AgentConsole onDispatch={handleDispatchAgent} />

          {/* Activity Log Card */}
          <div className="card" style={{ flex: 1 }}>
            <div className="card-header">
              <Cpu size={20} color="#4facfe" />
              Event Console & MCP Log Audit
            </div>
            <div
              style={{
                fontFamily: 'var(--font-mono)',
                fontSize: '0.8rem',
                color: '#9ca3af',
                display: 'flex',
                flexDirection: 'column',
                gap: '0.5rem',
                maxHeight: '200px',
                overflowY: 'auto',
              }}
            >
              {logs.map((log, idx) => (
                <div key={idx}>{log}</div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function App() {
  return (
    <LanguageProvider>
      <DashboardContent />
    </LanguageProvider>
  );
}


