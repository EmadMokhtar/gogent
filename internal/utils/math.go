package utils

import (
	"math"
)

// DotProduct calculates the dot product of two vectors
func DotProduct(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// Magnitude calculates the magnitude (L2 norm) of a vector
func Magnitude(v []float32) float32 {
	var sum float32
	for _, val := range v {
		sum += val * val
	}
	return float32(math.Sqrt(float64(sum)))
}

// CosineSimilarity calculates the cosine similarity between two vectors
// Returns a value between -1 and 1, where 1 means identical direction
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	
	dot := DotProduct(a, b)
	magA := Magnitude(a)
	magB := Magnitude(b)
	
	if magA == 0 || magB == 0 {
		return 0
	}
	
	return dot / (magA * magB)
}

// EuclideanDistance calculates the Euclidean distance between two vectors
func EuclideanDistance(a, b []float32) float32 {
	if len(a) != len(b) {
		return math.MaxFloat32
	}
	
	var sum float32
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	
	return float32(math.Sqrt(float64(sum)))
}

// Normalize normalizes a vector to unit length
func Normalize(v []float32) []float32 {
	mag := Magnitude(v)
	if mag == 0 {
		return v
	}
	
	normalized := make([]float32, len(v))
	for i, val := range v {
		normalized[i] = val / mag
	}
	
	return normalized
}
