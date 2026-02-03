package rag

import (
	"context"
	"errors"
	"testing"
)

// MockEmbedder is a mock implementation of Embedder for testing
type MockEmbedder struct {
	dimension int
}

func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, errors.New("empty text")
	}
	// Simple mock: convert text length to a vector
	// This is just for testing purposes
	vec := make([]float32, m.dimension)
	for i := range vec {
		vec[i] = float32(len(text)) / float32(m.dimension)
	}
	return vec, nil
}

func (m *MockEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var embeddings [][]float32
	for _, text := range texts {
		emb, err := m.Embed(ctx, text)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, emb)
	}
	return embeddings, nil
}

func (m *MockEmbedder) Dimension() int {
	return m.dimension
}

func TestRAG_AddDocuments(t *testing.T) {
	ctx := context.Background()

	embedder := &MockEmbedder{dimension: 3}
	store := NewInMemoryVectorStore()
	ragSys := New(
		WithVectorStore(store),
		WithEmbedder(embedder),
	)

	docs := []Document{
		{
			ID:      "doc1",
			Content: "test content",
			Metadata: map[string]interface{}{
				"category": "test",
			},
		},
	}

	err := ragSys.AddDocuments(ctx, docs)
	if err != nil {
		t.Fatalf("AddDocuments() failed: %v", err)
	}

	// Verify document was added
	stats, err := ragSys.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats() failed: %v", err)
	}

	if stats.TotalDocuments != 1 {
		t.Errorf("Stats() TotalDocuments = %v, expected 1", stats.TotalDocuments)
	}
}

func TestRAG_AddDocuments_WithChunker(t *testing.T) {
	ctx := context.Background()

	embedder := &MockEmbedder{dimension: 3}
	store := NewInMemoryVectorStore()
	chunker := NewFixedSizeChunker(WithChunkSize(10), WithOverlap(2))

	ragSys := New(
		WithVectorStore(store),
		WithEmbedder(embedder),
		WithChunker(chunker),
	)

	docs := []Document{
		{
			ID:      "doc1",
			Content: "This is a longer text that will be chunked into multiple pieces.",
		},
	}

	err := ragSys.AddDocuments(ctx, docs)
	if err != nil {
		t.Fatalf("AddDocuments() failed: %v", err)
	}

	// Verify multiple chunks were created
	stats, err := ragSys.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats() failed: %v", err)
	}

	// Should have more than 1 document (chunks)
	if stats.TotalDocuments <= 1 {
		t.Errorf("Stats() TotalDocuments = %v, expected > 1 (chunks)", stats.TotalDocuments)
	}
}

func TestRAG_Retrieve(t *testing.T) {
	ctx := context.Background()

	embedder := &MockEmbedder{dimension: 3}
	store := NewInMemoryVectorStore()

	ragSys := New(
		WithVectorStore(store),
		WithEmbedder(embedder),
		WithTopK(2),
	)

	docs := []Document{
		{
			ID:      "doc1",
			Content: "short",
			Metadata: map[string]interface{}{
				"category": "test",
			},
		},
		{
			ID:      "doc2",
			Content: "medium length text",
			Metadata: map[string]interface{}{
				"category": "test",
			},
		},
		{
			ID:      "doc3",
			Content: "this is a much longer piece of text",
			Metadata: map[string]interface{}{
				"category": "other",
			},
		},
	}

	err := ragSys.AddDocuments(ctx, docs)
	if err != nil {
		t.Fatalf("AddDocuments() failed: %v", err)
	}

	// Test retrieval
	results, err := ragSys.Retrieve(ctx, "short query")
	if err != nil {
		t.Fatalf("Retrieve() failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Retrieve() returned %d results, expected 2", len(results))
	}

	// Check that results are ranked
	for i, result := range results {
		if result.Rank != i+1 {
			t.Errorf("Result %d has rank %d, expected %d", i, result.Rank, i+1)
		}
	}
}

func TestRAG_Retrieve_WithFilter(t *testing.T) {
	ctx := context.Background()

	embedder := &MockEmbedder{dimension: 3}
	store := NewInMemoryVectorStore()

	ragSys := New(
		WithVectorStore(store),
		WithEmbedder(embedder),
	)

	docs := []Document{
		{
			ID:      "doc1",
			Content: "test1",
			Metadata: map[string]interface{}{
				"category": "A",
			},
		},
		{
			ID:      "doc2",
			Content: "test2",
			Metadata: map[string]interface{}{
				"category": "B",
			},
		},
		{
			ID:      "doc3",
			Content: "test3",
			Metadata: map[string]interface{}{
				"category": "A",
			},
		},
	}

	err := ragSys.AddDocuments(ctx, docs)
	if err != nil {
		t.Fatalf("AddDocuments() failed: %v", err)
	}

	// Test retrieval with filter
	filter := &Filter{
		Metadata: map[string]interface{}{
			"category": "A",
		},
	}

	results, err := ragSys.Retrieve(ctx, "query", WithFilter(filter), WithRetrievalK(10))
	if err != nil {
		t.Fatalf("Retrieve() failed: %v", err)
	}

	// Should only return documents with category A
	if len(results) != 2 {
		t.Errorf("Retrieve() with filter returned %d results, expected 2", len(results))
	}

	for _, result := range results {
		if result.Document.Metadata["category"] != "A" {
			t.Errorf("Result has category %v, expected A", result.Document.Metadata["category"])
		}
	}
}

func TestRAG_NoVectorStore(t *testing.T) {
	embedder := &MockEmbedder{dimension: 3}
	ragSys := New(WithEmbedder(embedder))

	ctx := context.Background()
	err := ragSys.AddDocuments(ctx, []Document{{ID: "1", Content: "test"}})
	if err != ErrNoVectorStore {
		t.Errorf("AddDocuments() error = %v, expected %v", err, ErrNoVectorStore)
	}
}

func TestRAG_NoEmbedder(t *testing.T) {
	store := NewInMemoryVectorStore()
	ragSys := New(WithVectorStore(store))

	ctx := context.Background()
	err := ragSys.AddDocuments(ctx, []Document{{ID: "1", Content: "test"}})
	if err != ErrNoEmbedder {
		t.Errorf("AddDocuments() error = %v, expected %v", err, ErrNoEmbedder)
	}
}

func TestRAG_EmptyDocument(t *testing.T) {
	ctx := context.Background()

	embedder := &MockEmbedder{dimension: 3}
	store := NewInMemoryVectorStore()
	ragSys := New(
		WithVectorStore(store),
		WithEmbedder(embedder),
	)

	docs := []Document{
		{
			ID:      "doc1",
			Content: "", // Empty content
		},
	}

	// Should not error, just skip empty documents
	err := ragSys.AddDocuments(ctx, docs)
	if err != nil {
		t.Fatalf("AddDocuments() with empty content failed: %v", err)
	}

	stats, _ := ragSys.Stats(ctx)
	if stats.TotalDocuments != 0 {
		t.Errorf("Expected 0 documents after adding empty document, got %d", stats.TotalDocuments)
	}
}
