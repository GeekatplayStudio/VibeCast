import React from 'react';
import { Wifi, AlertTriangle } from 'lucide-react';

interface ConnectionQualityBadgeProps {
  score: number; // 1 to 5
  packetLoss?: number;
  rtt?: number;
}

export const ConnectionQualityBadge: React.FC<ConnectionQualityBadgeProps> = ({ score, packetLoss = 0, rtt = 15 }) => {
  const getQualityColor = () => {
    if (score >= 4) return '#10b981'; // Green
    if (score >= 3) return '#3b82f6'; // Blue
    if (score >= 2) return '#f59e0b'; // Orange
    return '#ef4444'; // Red
  };

  const getQualityLabel = () => {
    if (score >= 4) return 'Excellent';
    if (score >= 3) return 'Good';
    if (score >= 2) return 'Fair';
    return 'Poor';
  };

  return (
    <div
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: '0.4rem',
        padding: '0.25rem 0.6rem',
        borderRadius: '6px',
        background: 'rgba(255, 255, 255, 0.05)',
        border: `1px solid ${getQualityColor()}`,
        fontSize: '0.75rem',
        fontWeight: 600,
        color: getQualityColor(),
      }}
      title={`Packet Loss: ${packetLoss}% | RTT: ${rtt}ms`}
    >
      {score <= 2 ? <AlertTriangle size={12} /> : <Wifi size={12} />}
      <span>{getQualityLabel()} ({rtt}ms)</span>
    </div>
  );
};
