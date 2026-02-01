package rag

import (
	"context"
	"fmt"
)

// Retriever handles document retrieval
type Retriever struct {
	vectorStore      VectorStore
	embedder         Embedder
	topK             int
	similarityMetric SimilarityMetric
}

// RetrieverConfig holds configuration for the retriever
type RetrieverConfig struct {
	VectorStore      VectorStore
	Embedder         Embedder
	TopK             int
	SimilarityMetric SimilarityMetric
}

// NewRetriever creates a new retriever
func NewRetriever(config RetrieverConfig) *Retriever {
	if config.TopK <= 0 {
		config.TopK = 5
	}
	if config.SimilarityMetric == "" {
		config.SimilarityMetric = CosineSimilarity
	}
	
	return &Retriever{
		vectorStore:      config.VectorStore,
		embedder:         config.Embedder,
		topK:             config.TopK,
		similarityMetric: config.SimilarityMetric,
	}
}

// Retrieve retrieves relevant chunks for a query
func (r *Retriever) Retrieve(ctx context.Context, query string, filters map[string]interface{}) ([]RetrievalResult, error) {
	return r.RetrieveWithTopK(ctx, query, r.topK, filters)
}

// RetrieveWithTopK retrieves relevant chunks for a query with a specific topK
func (r *Retriever) RetrieveWithTopK(ctx context.Context, query string, topK int, filters map[string]interface{}) ([]RetrievalResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}
	
	if topK <= 0 {
		topK = r.topK
	}
	
	// Generate embedding for the query
	queryEmbedding, err := r.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}
	
	// Search vector store
	results, err := r.vectorStore.Search(ctx, queryEmbedding, topK, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}
	
	return results, nil
}

// RetrieveWithScore retrieves relevant chunks with a minimum score threshold
func (r *Retriever) RetrieveWithScore(ctx context.Context, query string, minScore float32, filters map[string]interface{}) ([]RetrievalResult, error) {
	results, err := r.Retrieve(ctx, query, filters)
	if err != nil {
		return nil, err
	}
	
	// Filter by minimum score
	filtered := make([]RetrievalResult, 0)
	for _, result := range results {
		if result.Score >= minScore {
			filtered = append(filtered, result)
		}
	}
	
	return filtered, nil
}
