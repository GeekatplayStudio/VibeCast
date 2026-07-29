import React from 'react';
import { Activity, ArrowDownRight, ArrowUpRight, Cpu, Zap } from 'lucide-react';

interface TelemetryPanelProps {
  bytesIn: number;
  bytesOut: number;
  nackCount: number;
  rttMs?: number;
  packetLoss?: number;
  activeConnections?: number;
}

export const TelemetryPanel: React.FC<TelemetryPanelProps> = ({
  bytesIn,
  bytesOut,
  nackCount,
  rttMs = 18,
  packetLoss = 0.2,
  activeConnections = 2,
}) => {
  return (
    <div className="card">
      <div className="card-header" style={{ justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <Activity size={18} color="#10b981" />
          SFU Real-Time Media Telemetry
        </div>
        <span style={{ fontSize: '0.75rem', background: 'rgba(16, 185, 129, 0.2)', color: '#10b981', padding: '2px 8px', borderRadius: '12px' }}>
          LIVE
        </span>
      </div>

      <div className="metrics-list">
        <div className="metric-item">
          <span style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', color: '#9ca3af' }}>
            <ArrowDownRight size={14} color="#10b981" /> Inbound Bandwidth:
          </span>
          <span style={{ color: '#10b981', fontWeight: 600 }}>{(bytesIn / 1024).toFixed(1)} KB/s</span>
        </div>

        <div className="metric-item">
          <span style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', color: '#9ca3af' }}>
            <ArrowUpRight size={14} color="#00f2fe" /> Outbound Bandwidth:
          </span>
          <span style={{ color: '#00f2fe', fontWeight: 600 }}>{(bytesOut / 1024).toFixed(1)} KB/s</span>
        </div>

        <div className="metric-item">
          <span style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', color: '#9ca3af' }}>
            <Zap size={14} color="#f59e0b" /> Latency (RTT) / Loss:
          </span>
          <span style={{ color: '#f59e0b', fontWeight: 600 }}>{rttMs}ms | {packetLoss}%</span>
        </div>

        <div className="metric-item">
          <span style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', color: '#9ca3af' }}>
            <Cpu size={14} color="#4facfe" /> Active PeerConnections:
          </span>
          <span style={{ color: '#4facfe', fontWeight: 600 }}>{activeConnections} session(s)</span>
        </div>
      </div>
    </div>
  );
};

