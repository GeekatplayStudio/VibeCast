import React, { useState } from 'react';
import { MessageSquare, ShieldAlert, ShieldCheck, Send } from 'lucide-react';

export interface ChatMsg {
  id: string;
  senderName: string;
  text: string;
  timestamp: string;
  status: 'APPROVED' | 'FLAGGED' | 'BLOCKED';
  toxicity: number;
}

interface StreamChatProps {
  messages: ChatMsg[];
  onSendMessage: (text: string) => void;
}

export const StreamChat: React.FC<StreamChatProps> = ({ messages, onSendMessage }) => {
  const [inputText, setInputText] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (inputText.trim()) {
      onSendMessage(inputText);
      setInputText('');
    }
  };

  return (
    <div className="card" style={{ display: 'flex', flexDirection: 'column', height: '360px' }}>
      <div className="card-header" style={{ justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <MessageSquare size={18} color="#00f2fe" />
          Live Stream Chat
        </div>
        <span style={{ fontSize: '0.75rem', background: 'rgba(0, 242, 254, 0.15)', color: '#00f2fe', padding: '2px 8px', borderRadius: '10px' }}>
          AI VLM Moderated
        </span>
      </div>

      <div style={{ flex: 1, overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: '0.5rem', marginBottom: '0.75rem' }}>
        {messages.length === 0 ? (
          <div style={{ color: '#9ca3af', fontSize: '0.85rem', textAlign: 'center', margin: 'auto' }}>
            No chat messages yet. Type below to send.
          </div>
        ) : (
          messages.map((m) => (
            <div
              key={m.id}
              style={{
                background: m.status === 'FLAGGED' ? 'rgba(239, 68, 68, 0.1)' : 'rgba(255, 255, 255, 0.03)',
                border: `1px solid ${m.status === 'FLAGGED' ? 'rgba(239, 68, 68, 0.4)' : 'rgba(255, 255, 255, 0.05)'}`,
                borderRadius: '8px',
                padding: '0.5rem 0.75rem',
                fontSize: '0.85rem',
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.2rem' }}>
                <span style={{ fontWeight: 600, color: '#4facfe' }}>{m.senderName}</span>
                <span style={{ fontSize: '0.7rem', color: '#9ca3af', display: 'flex', alignItems: 'center', gap: '0.3rem' }}>
                  {m.status === 'FLAGGED' ? (
                    <span style={{ color: '#ef4444', display: 'flex', alignItems: 'center', gap: '0.2rem' }}>
                      <ShieldAlert size={12} /> Flagged ({Math.round(m.toxicity * 100)}%)
                    </span>
                  ) : (
                    <span style={{ color: '#10b981', display: 'flex', alignItems: 'center', gap: '0.2rem' }}>
                      <ShieldCheck size={12} /> Safe
                    </span>
                  )}
                  • {m.timestamp}
                </span>
              </div>
              <p style={{ color: m.status === 'FLAGGED' ? '#fca5a5' : '#f3f4f6' }}>{m.text}</p>
            </div>
          ))
        )}
      </div>

      <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '0.5rem' }}>
        <input
          type="text"
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          placeholder="Send a message to stream chat..."
          style={{
            flex: 1,
            background: 'rgba(255, 255, 255, 0.05)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '8px',
            padding: '0.5rem 0.75rem',
            color: '#fff',
            fontFamily: 'inherit',
            fontSize: '0.85rem',
          }}
        />
        <button type="submit" className="btn" style={{ padding: '0.5rem 1rem' }}>
          <Send size={14} />
        </button>
      </form>
    </div>
  );
};
