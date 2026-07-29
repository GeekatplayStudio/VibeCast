package telemetry

import (
	"sync/atomic"
)

// MetricsTracker records server-wide packet counts and active stream statistics.
type MetricsTracker struct {
	bytesIn       uint64
	bytesOut      uint64
	packetsIn     uint64
	packetsOut    uint64
	nackCount     uint64
	pliCount      uint64
}

// GlobalMetrics is the singleton metrics recorder.
var GlobalMetrics = &MetricsTracker{}

func (m *MetricsTracker) AddBytesIn(n uint64) {
	atomic.AddUint64(&m.bytesIn, n)
}

func (m *MetricsTracker) AddBytesOut(n uint64) {
	atomic.AddUint64(&m.bytesOut, n)
}

func (m *MetricsTracker) IncrementPacketsIn() {
	atomic.AddUint64(&m.packetsIn, 1)
}

func (m *MetricsTracker) IncrementPacketsOut() {
	atomic.AddUint64(&m.packetsOut, 1)
}

func (m *MetricsTracker) IncrementNACK() {
	atomic.AddUint64(&m.nackCount, 1)
}

func (m *MetricsTracker) IncrementPLI() {
	atomic.AddUint64(&m.pliCount, 1)
}

func (m *MetricsTracker) Snapshot() map[string]uint64 {
	return map[string]uint64{
		"bytes_in":    atomic.LoadUint64(&m.bytesIn),
		"bytes_out":   atomic.LoadUint64(&m.bytesOut),
		"packets_in":  atomic.LoadUint64(&m.packetsIn),
		"packets_out": atomic.LoadUint64(&m.packetsOut),
		"nack_count":  atomic.LoadUint64(&m.nackCount),
		"pli_count":   atomic.LoadUint64(&m.pliCount),
	}
}
