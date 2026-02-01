package rag

import (
	"github.com/EmadMokhtar/gogent/internal/utils"
)

// SimilarityMetric represents different similarity metrics
type SimilarityMetric string

const (
	// CosineSimilarity measures the cosine of the angle between vectors
	CosineSimilarity SimilarityMetric = "cosine"
	
	// DotProduct measures the dot product between vectors
	DotProduct SimilarityMetric = "dot"
	
	// EuclideanDistance measures the Euclidean distance between vectors
	EuclideanDistance SimilarityMetric = "euclidean"
)

// CalculateSimilarity calculates similarity between two vectors using the specified metric
func CalculateSimilarity(a, b []float32, metric SimilarityMetric) float32 {
	switch metric {
	case CosineSimilarity:
		return utils.CosineSimilarity(a, b)
	case DotProduct:
		return utils.DotProduct(a, b)
	case EuclideanDistance:
		// Convert distance to similarity (smaller distance = higher similarity)
		dist := utils.EuclideanDistance(a, b)
		if dist == 0 {
			return 1.0
		}
		return 1.0 / (1.0 + dist)
	default:
		return utils.CosineSimilarity(a, b)
	}
}
