package rag

import (
	"math"
)

// SimilarityMetric defines the type of similarity calculation
type SimilarityMetric int

const (
	// Cosine similarity (most common for text embeddings)
	Cosine SimilarityMetric = iota
	// DotProduct similarity (faster, works well with normalized vectors)
	DotProduct
	// Euclidean distance (L2 distance, lower is better)
	Euclidean
)

// CosineSimilarity calculates the cosine similarity between two vectors
// Returns a value between -1 and 1, where 1 means identical direction
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// DotProductSimilarity calculates the dot product between two vectors
// Higher values indicate greater similarity
func DotProductSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
	}

	return dotProduct
}

// EuclideanDistance calculates the Euclidean (L2) distance between two vectors
// Lower values indicate greater similarity (this is a distance, not similarity)
func EuclideanDistance(a, b []float32) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}

	var sum float64
	for i := range a {
		diff := float64(a[i] - b[i])
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

// CalculateSimilarity calculates similarity using the specified metric
func CalculateSimilarity(a, b []float32, metric SimilarityMetric) float64 {
	switch metric {
	case Cosine:
		return CosineSimilarity(a, b)
	case DotProduct:
		return DotProductSimilarity(a, b)
	case Euclidean:
		// Convert distance to similarity (invert so higher is better)
		dist := EuclideanDistance(a, b)
		if math.IsInf(dist, 1) {
			return 0
		}
		return 1.0 / (1.0 + dist)
	default:
		return CosineSimilarity(a, b)
	}
}
