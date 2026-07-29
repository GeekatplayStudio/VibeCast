import React, { useState } from 'react';
import { Eye, Layers, ShieldCheck, Sliders, Video } from 'lucide-react';

interface ProducerStudioProps {
  onSensitivityChange: (sensitivity: string) => void;
  onToggleOverlay: (overlay: string) => void;
}

export const ProducerStudio: React.FC<ProducerStudioProps> = ({
  onSensitivityChange,
  onToggleOverlay,
}) => {
  const [sensitivity, setSensitivity] = useState('MEDIUM');
  const [vlmStatus, setVlmStatus] = useState('ACTIVE');
  const [activeOverlays, setActiveOverlays] = useState<Record<string, boolean>>({
    lowerThird: true,
    aiAvatar: true,
    chatOverlay: false,
  });

  const handleSensitivityChange = (val: string) => {
    setSensitivity(val);
    onSensitivityChange(val);
  };

  const handleOverlayToggle = (key: string) => {
    const next = !activeOverlays[key];
    setActiveOverlays((prev) => ({ ...prev, [key]: next }));
    onToggleOverlay(key);
  };

  return (
    <div className="card">
      <div className="card-header" style={{ justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <Sliders size={18} color="#4facfe" />
          Producer Studio Control Suite
        </div>
        <span style={{ fontSize: '0.75rem', background: 'rgba(16, 185, 129, 0.2)', color: '#10b981', padding: '2px 8px', borderRadius: '10px' }}>
          VLM Stream Guardrails ON
        </span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
        {/* VLM Vision Analysis & Guardrails */}
        <div
          style={{
            background: 'rgba(255, 255, 255, 0.03)',
            borderRadius: '10px',
            padding: '0.75rem',
            border: '1px solid rgba(255, 255, 255, 0.05)',
          }}
        >
          <div style={{ fontSize: '0.85rem', fontWeight: 600, color: '#f3f4f6', marginBottom: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
            <Eye size={14} color="#00f2fe" /> Open-Source VLM Video Analysis
          </div>
          <div style={{ fontSize: '0.75rem', color: '#9ca3af', marginBottom: '0.75rem' }}>
            Monitors stream video frames for compliance & brand safety
          </div>
          <div style={{ display: 'flex', gap: '0.4rem' }}>
            {['LOW', 'MEDIUM', 'STRICT'].map((level) => (
              <button
                key={level}
                onClick={() => handleSensitivityChange(level)}
                style={{
                  flex: 1,
                  background: sensitivity === level ? 'linear-gradient(135deg, #00f2fe, #4facfe)' : 'rgba(255,255,255,0.05)',
                  color: sensitivity === level ? '#000' : '#9ca3af',
                  fontWeight: 600,
                  border: 'none',
                  borderRadius: '6px',
                  padding: '0.35rem 0.5rem',
                  fontSize: '0.75rem',
                  cursor: 'pointer',
                }}
              >
                {level}
              </button>
            ))}
          </div>
        </div>

        {/* Live Stream Overlay Triggers */}
        <div
          style={{
            background: 'rgba(255, 255, 255, 0.03)',
            borderRadius: '10px',
            padding: '0.75rem',
            border: '1px solid rgba(255, 255, 255, 0.05)',
          }}
        >
          <div style={{ fontSize: '0.85rem', fontWeight: 600, color: '#f3f4f6', marginBottom: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
            <Layers size={14} color="#10b981" /> Producer Stream Overlays
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.4rem' }}>
            {[
              { key: 'lowerThird', label: 'Lower Thirds Graphic' },
              { key: 'aiAvatar', label: 'AI Co-Host Video Overlay' },
              { key: 'chatOverlay', label: 'On-Screen Stream Chat' },
            ].map((item) => (
              <div
                key={item.key}
                onClick={() => handleOverlayToggle(item.key)}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '0.3rem 0.6rem',
                  background: activeOverlays[item.key] ? 'rgba(16, 185, 129, 0.15)' : 'rgba(255,255,255,0.03)',
                  border: `1px solid ${activeOverlays[item.key] ? '#10b981' : 'rgba(255,255,255,0.05)'}`,
                  borderRadius: '6px',
                  cursor: 'pointer',
                  fontSize: '0.75rem',
                  color: activeOverlays[item.key] ? '#10b981' : '#9ca3af',
                }}
              >
                <span>{item.label}</span>
                <span style={{ fontWeight: 600 }}>{activeOverlays[item.key] ? 'ON' : 'OFF'}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
