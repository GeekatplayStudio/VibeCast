import React from 'react';
import { Mic, MicOff, Shield, User, Video } from 'lucide-react';
import { ConnectionQualityBadge } from './ConnectionQualityBadge';
import { AudioLevelVisualizer } from './AudioLevelVisualizer';

export interface Participant {
  id: string;
  name: string;
  isMuted: boolean;
  role: 'publisher' | 'agent' | 'subscriber';
  qualityScore: number;
  rttMs: number;
  dbov: number;
}

interface RoomListProps {
  roomName: string;
  participants: Participant[];
  onToggleMute?: (id: string) => void;
}

export const RoomList: React.FC<RoomListProps> = ({ roomName, participants, onToggleMute }) => {
  return (
    <div className="card">
      <div className="card-header" style={{ justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <Video size={18} color="#00f2fe" />
          Active Room: {roomName}
        </div>
        <span style={{ fontSize: '0.8rem', color: '#9ca3af' }}>{participants.length} Participants</span>
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        {participants.length === 0 ? (
          <div style={{ fontSize: '0.85rem', color: '#9ca3af', textAlign: 'center', padding: '1rem' }}>
            No active participants in room
          </div>
        ) : (
          participants.map((p) => (
            <div
              key={p.id}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '0.75rem 1rem',
                background: 'rgba(255,255,255,0.03)',
                borderRadius: '10px',
                border: '1px solid rgba(255,255,255,0.05)',
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                <User size={16} color={p.role === 'agent' ? '#00f2fe' : '#9ca3af'} />
                <div>
                  <div style={{ fontSize: '0.9rem', color: '#f3f4f6', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
                    {p.name}
                    {p.role === 'agent' && (
                      <span style={{ fontSize: '0.65rem', background: 'rgba(0, 242, 254, 0.2)', color: '#00f2fe', padding: '1px 6px', borderRadius: '4px' }}>
                        MCP AGENT
                      </span>
                    )}
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.2rem' }}>
                    <ConnectionQualityBadge score={p.qualityScore} rtt={p.rttMs} />
                    <AudioLevelVisualizer isActive={!p.isMuted} levelDbov={p.dbov} />
                  </div>
                </div>
              </div>
              <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                <button
                  onClick={() => onToggleMute && onToggleMute(p.id)}
                  style={{
                    background: p.isMuted ? 'rgba(239, 68, 68, 0.2)' : 'rgba(16, 185, 129, 0.2)',
                    border: `1px solid ${p.isMuted ? '#ef4444' : '#10b981'}`,
                    color: p.isMuted ? '#ef4444' : '#10b981',
                    borderRadius: '6px',
                    padding: '0.3rem 0.6rem',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '0.3rem',
                    fontSize: '0.75rem',
                  }}
                >
                  {p.isMuted ? <MicOff size={12} /> : <Mic size={12} />}
                  {p.isMuted ? 'Muted' : 'Live'}
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

