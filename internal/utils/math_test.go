package utils

import (
	"math"
	"testing"
)

func TestDotProduct(t *testing.T) {
	tests := []struct {
		name string
		a    []float32
		b    []float32
		want float32
	}{
		{
			name: "simple dot product",
			a:    []float32{1, 2, 3},
			b:    []float32{4, 5, 6},
			want: 32, // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
		},
		{
			name: "zero vectors",
			a:    []float32{0, 0, 0},
			b:    []float32{1, 2, 3},
			want: 0,
		},
		{
			name: "different lengths",
			a:    []float32{1, 2},
			b:    []float32{1, 2, 3},
			want: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DotProduct(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("DotProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	tests := []struct {
		name string
		v    []float32
		want float32
	}{
		{
			name: "simple magnitude",
			v:    []float32{3, 4},
			want: 5, // sqrt(9 + 16) = sqrt(25) = 5
		},
		{
			name: "zero vector",
			v:    []float32{0, 0, 0},
			want: 0,
		},
		{
			name: "unit vector",
			v:    []float32{1, 0, 0},
			want: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Magnitude(tt.v)
			if math.Abs(float64(got-tt.want)) > 0.0001 {
				t.Errorf("Magnitude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    []float32
		b    []float32
		want float32
	}{
		{
			name: "identical vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{1, 2, 3},
			want: 1.0,
		},
		{
			name: "orthogonal vectors",
			a:    []float32{1, 0},
			b:    []float32{0, 1},
			want: 0.0,
		},
		{
			name: "opposite vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{-1, -2, -3},
			want: -1.0,
		},
		{
			name: "different lengths",
			a:    []float32{1, 2},
			b:    []float32{1, 2, 3},
			want: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CosineSimilarity(tt.a, tt.b)
			if math.Abs(float64(got-tt.want)) > 0.0001 {
				t.Errorf("CosineSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEuclideanDistance(t *testing.T) {
	tests := []struct {
		name string
		a    []float32
		b    []float32
		want float32
	}{
		{
			name: "simple distance",
			a:    []float32{0, 0},
			b:    []float32{3, 4},
			want: 5, // sqrt((3-0)^2 + (4-0)^2) = sqrt(25) = 5
		},
		{
			name: "identical vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{1, 2, 3},
			want: 0,
		},
		{
			name: "different lengths",
			a:    []float32{1, 2},
			b:    []float32{1, 2, 3},
			want: math.MaxFloat32,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EuclideanDistance(tt.a, tt.b)
			if math.Abs(float64(got-tt.want)) > 0.0001 {
				t.Errorf("EuclideanDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		v    []float32
	}{
		{
			name: "simple vector",
			v:    []float32{3, 4},
		},
		{
			name: "already normalized",
			v:    []float32{1, 0},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.v)
			mag := Magnitude(got)
			if math.Abs(float64(mag-1.0)) > 0.0001 && mag != 0 {
				t.Errorf("Normalize() magnitude = %v, want 1.0", mag)
			}
		})
	}
}
