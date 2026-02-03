package rag

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
		tolerance float64
	}{
		{
			name:      "identical vectors",
			a:         []float32{1, 2, 3},
			b:         []float32{1, 2, 3},
			expected:  1.0,
			tolerance: 0.0001,
		},
		{
			name:      "opposite vectors",
			a:         []float32{1, 0, 0},
			b:         []float32{-1, 0, 0},
			expected:  -1.0,
			tolerance: 0.0001,
		},
		{
			name:      "orthogonal vectors",
			a:         []float32{1, 0, 0},
			b:         []float32{0, 1, 0},
			expected:  0.0,
			tolerance: 0.0001,
		},
		{
			name:      "different lengths",
			a:         []float32{1, 2},
			b:         []float32{1, 2, 3},
			expected:  0.0,
			tolerance: 0.0,
		},
		{
			name:      "zero vector",
			a:         []float32{0, 0, 0},
			b:         []float32{1, 2, 3},
			expected:  0.0,
			tolerance: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.a, tt.b)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("CosineSimilarity() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestDotProductSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
	}{
		{
			name:     "positive dot product",
			a:        []float32{1, 2, 3},
			b:        []float32{4, 5, 6},
			expected: 32.0, // 1*4 + 2*5 + 3*6 = 32
		},
		{
			name:     "zero dot product",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "different lengths",
			a:        []float32{1, 2},
			b:        []float32{1, 2, 3},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DotProductSimilarity(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("DotProductSimilarity() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestEuclideanDistance(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
		tolerance float64
	}{
		{
			name:      "identical vectors",
			a:         []float32{1, 2, 3},
			b:         []float32{1, 2, 3},
			expected:  0.0,
			tolerance: 0.0001,
		},
		{
			name:      "unit distance",
			a:         []float32{0, 0, 0},
			b:         []float32{1, 0, 0},
			expected:  1.0,
			tolerance: 0.0001,
		},
		{
			name:      "sqrt(3) distance",
			a:         []float32{0, 0, 0},
			b:         []float32{1, 1, 1},
			expected:  math.Sqrt(3),
			tolerance: 0.0001,
		},
		{
			name:     "different lengths",
			a:        []float32{1, 2},
			b:        []float32{1, 2, 3},
			expected: math.Inf(1),
			tolerance: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EuclideanDistance(tt.a, tt.b)
			if math.IsInf(tt.expected, 1) {
				if !math.IsInf(result, 1) {
					t.Errorf("EuclideanDistance() = %v, expected +Inf", result)
				}
			} else if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("EuclideanDistance() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateSimilarity(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{1, 2, 3}

	// Test Cosine
	cosine := CalculateSimilarity(a, b, Cosine)
	if math.Abs(cosine-1.0) > 0.0001 {
		t.Errorf("CalculateSimilarity(Cosine) = %v, expected 1.0", cosine)
	}

	// Test DotProduct
	dot := CalculateSimilarity(a, b, DotProduct)
	if math.Abs(dot-14.0) > 0.0001 {
		t.Errorf("CalculateSimilarity(DotProduct) = %v, expected 14.0", dot)
	}

	// Test Euclidean
	euclidean := CalculateSimilarity(a, b, Euclidean)
	if math.Abs(euclidean-1.0) > 0.0001 {
		t.Errorf("CalculateSimilarity(Euclidean) = %v, expected 1.0", euclidean)
	}
}
