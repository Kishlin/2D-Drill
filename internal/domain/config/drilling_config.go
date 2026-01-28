package config

type DrillingConfig struct {
	MinDrillingDuration   float32 // seconds at ground level (base duration)
	MaxDrillingDuration   float32 // depth factor at max depth
	FloorDrillingDuration float32 // absolute minimum (safety clamp)
}
