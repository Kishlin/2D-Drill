package types_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/types"
)

func TestAABB_Intersects(t *testing.T) {
	tests := []struct {
		name     string
		a, b     types.AABB
		expected bool
	}{
		{
			name:     "Overlapping boxes",
			a:        types.NewAABB(0, 0, 64, 64),
			b:        types.NewAABB(32, 32, 64, 64),
			expected: true,
		},
		{
			name:     "Separated boxes",
			a:        types.NewAABB(0, 0, 64, 64),
			b:        types.NewAABB(100, 100, 64, 64),
			expected: false,
		},
		{
			name:     "Edge touching (no overlap)",
			a:        types.NewAABB(0, 0, 64, 64),
			b:        types.NewAABB(64, 0, 64, 64),
			expected: false,
		},
		{
			name:     "Fully contained",
			a:        types.NewAABB(10, 10, 20, 20),
			b:        types.NewAABB(0, 0, 64, 64),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Intersects(tt.b)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAABB_Penetration(t *testing.T) {
	tests := []struct {
		name       string
		a, b       types.AABB
		expectedDX float32
		expectedDY float32
	}{
		{
			name:       "Overlap from left",
			a:          types.NewAABB(0, 0, 64, 64),
			b:          types.NewAABB(32, 0, 64, 64),
			expectedDX: 32,  // a.X < b.X, so overlap X = 32
			expectedDY: -64, // a.Y == b.Y, goes to else: -overlapY
		},
		{
			name:       "Overlap from right",
			a:          types.NewAABB(64, 0, 64, 64),
			b:          types.NewAABB(32, 0, 64, 64),
			expectedDX: -32, // a.X > b.X, so -overlapX = -32
			expectedDY: -64, // a.Y == b.Y, goes to else: -overlapY
		},
		{
			name:       "Overlap from above",
			a:          types.NewAABB(0, 0, 64, 64),
			b:          types.NewAABB(0, 32, 64, 64),
			expectedDX: -64, // a.X == b.X, goes to else: -overlapX
			expectedDY: 32,  // a.Y < b.Y, so overlapY = 32
		},
		{
			name:       "Overlap from below",
			a:          types.NewAABB(0, 64, 64, 64),
			b:          types.NewAABB(0, 32, 64, 64),
			expectedDX: -64, // a.X == b.X, goes to else: -overlapX
			expectedDY: -32, // a.Y > b.Y, so -overlapY = -32
		},
		{
			name:       "No overlap",
			a:          types.NewAABB(0, 0, 64, 64),
			b:          types.NewAABB(100, 100, 64, 64),
			expectedDX: 0,
			expectedDY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dx, dy := tt.a.Penetration(tt.b)
			if dx != tt.expectedDX || dy != tt.expectedDY {
				t.Errorf("Expected (%f, %f), got (%f, %f)",
					tt.expectedDX, tt.expectedDY, dx, dy)
			}
		})
	}
}

func TestAABB_MinMax(t *testing.T) {
	aabb := types.NewAABB(10, 20, 30, 40)

	minVal := aabb.Min()
	if minVal.X != 10 || minVal.Y != 20 {
		t.Errorf("Min() expected (10, 20), got (%f, %f)", minVal.X, minVal.Y)
	}

	maxVal := aabb.Max()
	if maxVal.X != 40 || maxVal.Y != 60 {
		t.Errorf("Max() expected (40, 60), got (%f, %f)", maxVal.X, maxVal.Y)
	}
}
