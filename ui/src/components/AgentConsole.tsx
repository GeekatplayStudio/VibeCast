import React, { useState } from 'react';
import { Bot, Cpu, Send, Wrench } from 'lucide-react';

interface AgentConsoleProps {
  onDispatch: (prompt: string, model: string, tool: string) => void;
}

export const AgentConsole: React.FC<AgentConsoleProps> = ({ onDispatch }) => {
  const [prompt, setPrompt] = useState('Act as real-time audio translator & conversation summary assistant');
  const [selectedModel, setSelectedModel] = useState('gemini-3.6-flash');
  const [selectedTool, setSelectedTool] = useState('dispatch_agent');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (prompt.trim()) {
      onDispatch(prompt, selectedModel, selectedTool);
    }
  };

  return (
    <div className="card">
      <div className="card-header">
        <Bot size={18} color="#4facfe" />
        MCP Agent & Tool Console
      </div>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        <div style={{ display: 'flex', gap: '0.5rem' }}>
          <select
            value={selectedModel}
            onChange={(e) => setSelectedModel(e.target.value)}
            style={{
              flex: 1,
              background: 'rgba(255,255,255,0.05)',
              border: '1px solid rgba(255,255,255,0.1)',
              borderRadius: '8px',
              padding: '0.5rem',
              color: '#fff',
              fontFamily: 'inherit',
              fontSize: '0.85rem',
            }}
          >
            <option value="gemini-3.6-flash">Gemini 3.6 Flash (High Speed)</option>
            <option value="claude-3.5-sonnet">Claude 3.5 Sonnet</option>
            <option value="gpt-4o-realtime">GPT-4o Realtime Voice</option>
          </select>

          <select
            value={selectedTool}
            onChange={(e) => setSelectedTool(e.target.value)}
            style={{
              flex: 1,
              background: 'rgba(255,255,255,0.05)',
              border: '1px solid rgba(255,255,255,0.1)',
              borderRadius: '8px',
              padding: '0.5rem',
              color: '#fff',
              fontFamily: 'inherit',
              fontSize: '0.85rem',
            }}
          >
            <option value="dispatch_agent">MCP Tool: dispatch_agent</option>
            <option value="mute_participant">MCP Tool: mute_participant</option>
            <option value="get_telemetry">MCP Tool: get_telemetry</option>
          </select>
        </div>

        <textarea
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
          rows={3}
          style={{
            background: 'rgba(255,255,255,0.03)',
            border: '1px solid rgba(255,255,255,0.08)',
            borderRadius: '8px',
            padding: '0.75rem',
            color: '#fff',
            fontFamily: 'inherit',
            resize: 'none',
          }}
        />
        <button
          type="submit"
          className="btn"
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            gap: '0.5rem',
            background: 'linear-gradient(135deg, #4facfe, #00f2fe)',
          }}
        >
          <Send size={16} /> Execute MCP Tool Call
        </button>
      </form>
    </div>
  );
};

