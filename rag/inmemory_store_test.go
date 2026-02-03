package rag

import (
	"context"
	"testing"
)

func TestInMemoryVectorStore_AddAndGet(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	doc := Document{
		ID:        "test1",
		Content:   "test content",
		Embedding: []float32{1, 2, 3},
		Metadata:  map[string]interface{}{"key": "value"},
	}

	// Test Add
	err := store.Add(ctx, []Document{doc})
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Test Get
	retrieved, err := store.Get(ctx, "test1")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if retrieved.ID != doc.ID {
		t.Errorf("Get() ID = %v, expected %v", retrieved.ID, doc.ID)
	}
	if retrieved.Content != doc.Content {
		t.Errorf("Get() Content = %v, expected %v", retrieved.Content, doc.Content)
	}

	// Test Get non-existent
	_, err = store.Get(ctx, "nonexistent")
	if err != ErrDocumentNotFound {
		t.Errorf("Get() error = %v, expected %v", err, ErrDocumentNotFound)
	}
}

func TestInMemoryVectorStore_Search(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	docs := []Document{
		{
			ID:        "doc1",
			Content:   "about go",
			Embedding: []float32{1, 0, 0},
			Metadata:  map[string]interface{}{"language": "go"},
		},
		{
			ID:        "doc2",
			Content:   "about python",
			Embedding: []float32{0, 1, 0},
			Metadata:  map[string]interface{}{"language": "python"},
		},
		{
			ID:        "doc3",
			Content:   "about go advanced",
			Embedding: []float32{0.9, 0.1, 0},
			Metadata:  map[string]interface{}{"language": "go"},
		},
	}

	err := store.Add(ctx, docs)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Search for go-related content
	query := []float32{1, 0, 0}
	results, err := store.Search(ctx, query, 2, nil)
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Search() returned %d results, expected 2", len(results))
	}

	// First result should be doc1 (exact match)
	if results[0].Document.ID != "doc1" {
		t.Errorf("Search() first result ID = %v, expected doc1", results[0].Document.ID)
	}

	// Check that results are ranked
	if results[0].Rank != 1 {
		t.Errorf("Search() first result rank = %v, expected 1", results[0].Rank)
	}
	if results[1].Rank != 2 {
		t.Errorf("Search() second result rank = %v, expected 2", results[1].Rank)
	}
}

func TestInMemoryVectorStore_SearchWithFilter(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	docs := []Document{
		{
			ID:        "doc1",
			Content:   "go basics",
			Embedding: []float32{1, 0, 0},
			Metadata:  map[string]interface{}{"language": "go", "level": "beginner"},
		},
		{
			ID:        "doc2",
			Content:   "python basics",
			Embedding: []float32{1, 0, 0},
			Metadata:  map[string]interface{}{"language": "python", "level": "beginner"},
		},
		{
			ID:        "doc3",
			Content:   "go advanced",
			Embedding: []float32{1, 0, 0},
			Metadata:  map[string]interface{}{"language": "go", "level": "advanced"},
		},
	}

	err := store.Add(ctx, docs)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Search with filter
	query := []float32{1, 0, 0}
	filter := &Filter{
		Metadata: map[string]interface{}{
			"language": "go",
		},
	}

	results, err := store.Search(ctx, query, 10, filter)
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Search() with filter returned %d results, expected 2", len(results))
	}

	// Check that all results match the filter
	for _, result := range results {
		if result.Document.Metadata["language"] != "go" {
			t.Errorf("Result has language = %v, expected go", result.Document.Metadata["language"])
		}
	}
}

func TestInMemoryVectorStore_Delete(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	doc := Document{
		ID:        "test1",
		Content:   "test",
		Embedding: []float32{1, 2, 3},
	}

	// Add document
	err := store.Add(ctx, []Document{doc})
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Delete document
	err = store.Delete(ctx, []string{"test1"})
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Verify it's gone
	_, err = store.Get(ctx, "test1")
	if err != ErrDocumentNotFound {
		t.Errorf("Get() after Delete() error = %v, expected %v", err, ErrDocumentNotFound)
	}
}

func TestInMemoryVectorStore_Stats(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	docs := []Document{
		{ID: "1", Content: "test", Embedding: []float32{1, 2, 3}},
		{ID: "2", Content: "test2", Embedding: []float32{4, 5, 6}},
	}

	err := store.Add(ctx, docs)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	stats, err := store.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats() failed: %v", err)
	}

	if stats.TotalDocuments != 2 {
		t.Errorf("Stats() TotalDocuments = %v, expected 2", stats.TotalDocuments)
	}

	if stats.Dimension != 3 {
		t.Errorf("Stats() Dimension = %v, expected 3", stats.Dimension)
	}

	if stats.MemoryUsage == 0 {
		t.Error("Stats() MemoryUsage = 0, expected > 0")
	}
}

func TestInMemoryVectorStore_EmptyEmbedding(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	doc := Document{
		ID:        "test1",
		Content:   "test",
		Embedding: []float32{},
	}

	err := store.Add(ctx, []Document{doc})
	if err != ErrEmptyEmbedding {
		t.Errorf("Add() with empty embedding error = %v, expected %v", err, ErrEmptyEmbedding)
	}
}
