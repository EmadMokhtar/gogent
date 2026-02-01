package rag

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// VectorStore is the interface for storing and retrieving vectors
type VectorStore interface {
	// Add adds a chunk with its embedding
	Add(ctx context.Context, chunk Chunk) error
	
	// AddBatch adds multiple chunks with their embeddings
	AddBatch(ctx context.Context, chunks []Chunk) error
	
	// Search performs similarity search and returns top-k results
	Search(ctx context.Context, queryEmbedding []float32, topK int, filters map[string]interface{}) ([]RetrievalResult, error)
	
	// Get retrieves a chunk by ID
	Get(ctx context.Context, id string) (*Chunk, error)
	
	// Delete removes a chunk by ID
	Delete(ctx context.Context, id string) error
	
	// Clear removes all chunks
	Clear(ctx context.Context) error
	
	// Count returns the number of stored chunks
	Count(ctx context.Context) (int, error)
}

// InMemoryVectorStore implements an in-memory vector store
type InMemoryVectorStore struct {
	mu               sync.RWMutex
	chunks           map[string]Chunk
	similarityMetric SimilarityMetric
}

// VectorStoreOption is a functional option for configuring the vector store
type VectorStoreOption func(*InMemoryVectorStore)

// NewInMemoryVectorStore creates a new in-memory vector store
func NewInMemoryVectorStore(opts ...VectorStoreOption) *InMemoryVectorStore {
	vs := &InMemoryVectorStore{
		chunks:           make(map[string]Chunk),
		similarityMetric: CosineSimilarity,
	}
	
	for _, opt := range opts {
		opt(vs)
	}
	
	return vs
}

// WithSimilarityMetric sets the similarity metric
func WithSimilarityMetric(metric SimilarityMetric) VectorStoreOption {
	return func(vs *InMemoryVectorStore) {
		vs.similarityMetric = metric
	}
}

// Add adds a chunk with its embedding
func (vs *InMemoryVectorStore) Add(ctx context.Context, chunk Chunk) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	
	if chunk.ID == "" {
		return fmt.Errorf("chunk ID cannot be empty")
	}
	
	if len(chunk.Embedding) == 0 {
		return fmt.Errorf("chunk embedding cannot be empty")
	}
	
	vs.chunks[chunk.ID] = chunk
	return nil
}

// AddBatch adds multiple chunks with their embeddings
func (vs *InMemoryVectorStore) AddBatch(ctx context.Context, chunks []Chunk) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	
	for _, chunk := range chunks {
		if chunk.ID == "" {
			return fmt.Errorf("chunk ID cannot be empty")
		}
		
		if len(chunk.Embedding) == 0 {
			return fmt.Errorf("chunk embedding cannot be empty for chunk %s", chunk.ID)
		}
		
		vs.chunks[chunk.ID] = chunk
	}
	
	return nil
}

// Search performs similarity search and returns top-k results
func (vs *InMemoryVectorStore) Search(ctx context.Context, queryEmbedding []float32, topK int, filters map[string]interface{}) ([]RetrievalResult, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	
	if len(queryEmbedding) == 0 {
		return nil, fmt.Errorf("query embedding cannot be empty")
	}
	
	if topK <= 0 {
		topK = 5
	}
	
	results := make([]RetrievalResult, 0)
	
	// Calculate similarity for all chunks
	for _, chunk := range vs.chunks {
		// Apply filters if provided
		if !matchesFilters(chunk, filters) {
			continue
		}
		
		similarity := CalculateSimilarity(queryEmbedding, chunk.Embedding, vs.similarityMetric)
		
		result := RetrievalResult{
			Chunk:    chunk,
			Score:    similarity,
			Metadata: chunk.Metadata,
		}
		
		results = append(results, result)
	}
	
	// Sort by score (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	// Return top-k results
	if topK < len(results) {
		results = results[:topK]
	}
	
	return results, nil
}

// Get retrieves a chunk by ID
func (vs *InMemoryVectorStore) Get(ctx context.Context, id string) (*Chunk, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	
	chunk, exists := vs.chunks[id]
	if !exists {
		return nil, fmt.Errorf("chunk %s not found", id)
	}
	
	return &chunk, nil
}

// Delete removes a chunk by ID
func (vs *InMemoryVectorStore) Delete(ctx context.Context, id string) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	
	if _, exists := vs.chunks[id]; !exists {
		return fmt.Errorf("chunk %s not found", id)
	}
	
	delete(vs.chunks, id)
	return nil
}

// Clear removes all chunks
func (vs *InMemoryVectorStore) Clear(ctx context.Context) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	
	vs.chunks = make(map[string]Chunk)
	return nil
}

// Count returns the number of stored chunks
func (vs *InMemoryVectorStore) Count(ctx context.Context) (int, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	
	return len(vs.chunks), nil
}

// matchesFilters checks if a chunk matches the provided filters
func matchesFilters(chunk Chunk, filters map[string]interface{}) bool {
	if len(filters) == 0 {
		return true
	}
	
	for key, value := range filters {
		chunkValue, exists := chunk.Metadata[key]
		if !exists || chunkValue != value {
			return false
		}
	}
	
	return true
}
