package rag

import (
	"context"
	"errors"
	"sort"
	"sync"
)

var (
	// ErrDocumentNotFound is returned when a document ID is not found
	ErrDocumentNotFound = errors.New("document not found")
	// ErrEmptyEmbedding is returned when a document has no embedding
	ErrEmptyEmbedding = errors.New("document has empty embedding")
)

// InMemoryVectorStore is a thread-safe in-memory implementation of VectorStore
type InMemoryVectorStore struct {
	mu        sync.RWMutex
	documents map[string]Document // ID -> Document
	metric    SimilarityMetric
}

// InMemoryVectorStoreOption defines configuration options for InMemoryVectorStore
type InMemoryVectorStoreOption func(*InMemoryVectorStore)

// WithMetric sets the similarity metric for the vector store
func WithMetric(metric SimilarityMetric) InMemoryVectorStoreOption {
	return func(s *InMemoryVectorStore) {
		s.metric = metric
	}
}

// NewInMemoryVectorStore creates a new in-memory vector store
func NewInMemoryVectorStore(opts ...InMemoryVectorStoreOption) *InMemoryVectorStore {
	store := &InMemoryVectorStore{
		documents: make(map[string]Document),
		metric:    Cosine, // Default to cosine similarity
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

// Add stores documents with their embeddings in the vector store
func (s *InMemoryVectorStore) Add(ctx context.Context, docs []Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, doc := range docs {
		if len(doc.Embedding) == 0 {
			return ErrEmptyEmbedding
		}
		s.documents[doc.ID] = doc
	}

	return nil
}

// Search finds the k most similar documents to the query embedding
func (s *InMemoryVectorStore) Search(ctx context.Context, query []float32, k int, filter *Filter) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(query) == 0 {
		return nil, ErrEmptyEmbedding
	}

	// Calculate similarity for all documents
	var results []SearchResult
	for _, doc := range s.documents {
		// Skip documents that don't match the filter
		if filter != nil && !filter.Matches(doc) {
			continue
		}

		// Skip documents with no embedding
		if len(doc.Embedding) == 0 {
			continue
		}

		// Check dimension compatibility
		if len(doc.Embedding) != len(query) {
			continue
		}

		score := CalculateSimilarity(query, doc.Embedding, s.metric)
		results = append(results, SearchResult{
			Document: doc,
			Score:    score,
		})
	}

	// Sort by score (descending - higher scores are better)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit to top k results
	if k < len(results) {
		results = results[:k]
	}

	// Set ranks
	for i := range results {
		results[i].Rank = i + 1
	}

	return results, nil
}

// Delete removes documents by their IDs
func (s *InMemoryVectorStore) Delete(ctx context.Context, ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range ids {
		delete(s.documents, id)
	}

	return nil
}

// Get retrieves a document by its ID
func (s *InMemoryVectorStore) Get(ctx context.Context, id string) (*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, exists := s.documents[id]
	if !exists {
		return nil, ErrDocumentNotFound
	}

	return &doc, nil
}

// Stats returns statistics about the vector store
func (s *InMemoryVectorStore) Stats(ctx context.Context) (*VectorStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &VectorStats{
		TotalDocuments: len(s.documents),
	}

	// Calculate dimension and memory usage
	for _, doc := range s.documents {
		if len(doc.Embedding) > 0 {
			stats.Dimension = len(doc.Embedding)
			break
		}
	}

	// Rough estimate of memory usage
	for _, doc := range s.documents {
		stats.MemoryUsage += int64(len(doc.ID))
		stats.MemoryUsage += int64(len(doc.Content))
		stats.MemoryUsage += int64(len(doc.Source))
		stats.MemoryUsage += int64(len(doc.Embedding) * 4) // 4 bytes per float32
	}

	return stats, nil
}
