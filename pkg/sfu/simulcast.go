package sfu

type VideoLayer int

const (
	LayerLow VideoLayer = iota
	LayerMid
	LayerHigh
)

func (v VideoLayer) String() string {
	switch v {
	case LayerLow:
		return "low"
	case LayerMid:
		return "mid"
	case LayerHigh:
		return "high"
	default:
		return "unknown"
	}
}

// SimulcastAllocator selects the optimal video layer based on subscriber target bitrates.
type SimulcastAllocator struct{}

func NewSimulcastAllocator() *SimulcastAllocator {
	return &SimulcastAllocator{}
}

// SelectOptimalLayer returns the appropriate VideoLayer for a target available bandwidth.
func (s *SimulcastAllocator) SelectOptimalLayer(targetBitrateBps uint32) VideoLayer {
	if targetBitrateBps >= 1_500_000 {
		return LayerHigh
	} else if targetBitrateBps >= 500_000 {
		return LayerMid
	}
	return LayerLow
}
