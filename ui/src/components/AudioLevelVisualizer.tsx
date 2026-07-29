import React from 'react';

interface AudioLevelVisualizerProps {
  isActive: boolean;
  levelDbov?: number; // e.g. -60 to 0
}

export const AudioLevelVisualizer: React.FC<AudioLevelVisualizerProps> = ({ isActive, levelDbov = -50 }) => {
  // Normalize dBov to 0..1 height percentage
  const normalized = Math.max(0, Math.min(100, (levelDbov + 60) * 1.66));

  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: '3px',
        height: '18px',
        padding: '0 4px',
      }}
    >
      {[0.8, 1.2, 0.6, 1.4, 0.9].map((scale, idx) => (
        <div
          key={idx}
          style={{
            width: '3px',
            height: isActive ? `${Math.max(4, normalized * scale)}px` : '4px',
            borderRadius: '2px',
            background: isActive
              ? 'linear-gradient(180deg, #00f2fe, #4facfe)'
              : 'rgba(255, 255, 255, 0.2)',
            transition: 'height 0.1s ease',
          }}
        />
      ))}
    </div>
  );
};
