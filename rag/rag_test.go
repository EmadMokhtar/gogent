package rag

import (
	"context"
	"testing"
)

func TestChunker(t *testing.T) {
	t.Run("chunk document", func(t *testing.T) {
		chunker := NewChunker(ChunkerConfig{
			ChunkSize:    100,
			ChunkOverlap: 20,
		})
		
		doc := NewDocument("test-doc", "This is a test document with enough content to be split into multiple chunks. Each chunk should have some overlap with the next chunk to maintain context.")
		
		chunks, err := chunker.ChunkDocument(doc)
		if err != nil {
			t.Fatalf("ChunkDocument() error = %v", err)
		}
		
		if len(chunks) == 0 {
			t.Error("ChunkDocument() returned no chunks")
		}
		
		for i, chunk := range chunks {
			if chunk.DocumentID != doc.ID {
				t.Errorf("Chunk %d DocumentID = %s, want %s", i, chunk.DocumentID, doc.ID)
			}
			if chunk.Index != i {
				t.Errorf("Chunk %d Index = %d, want %d", i, chunk.Index, i)
			}
		}
	})
	
	t.Run("empty document", func(t *testing.T) {
		chunker := NewChunker(ChunkerConfig{ChunkSize: 100})
		doc := NewDocument("empty", "")
		
		_, err := chunker.ChunkDocument(doc)
		if err == nil {
			t.Error("ChunkDocument() expected error for empty document")
		}
	})
}

func TestVectorStore(t *testing.T) {
	ctx := context.Background()
	
	t.Run("add and get chunk", func(t *testing.T) {
		vs := NewInMemoryVectorStore()
		
		chunk := Chunk{
			ID:        "chunk-1",
			Content:   "Test content",
			Embedding: []float32{0.1, 0.2, 0.3},
		}
		
		if err := vs.Add(ctx, chunk); err != nil {
			t.Fatalf("Add() error = %v", err)
		}
		
		got, err := vs.Get(ctx, "chunk-1")
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		
		if got.ID != chunk.ID {
			t.Errorf("Get() ID = %s, want %s", got.ID, chunk.ID)
		}
	})
	
	t.Run("search chunks", func(t *testing.T) {
		vs := NewInMemoryVectorStore()
		
		chunks := []Chunk{
			{
				ID:        "chunk-1",
				Content:   "Go programming",
				Embedding: []float32{1.0, 0.0, 0.0},
			},
			{
				ID:        "chunk-2",
				Content:   "Python programming",
				Embedding: []float32{0.0, 1.0, 0.0},
			},
		}
		
		if err := vs.AddBatch(ctx, chunks); err != nil {
			t.Fatalf("AddBatch() error = %v", err)
		}
		
		query := []float32{0.9, 0.1, 0.0}
		results, err := vs.Search(ctx, query, 2, nil)
		if err != nil {
			t.Fatalf("Search() error = %v", err)
		}
		
		if len(results) != 2 {
			t.Errorf("Search() returned %d results, want 2", len(results))
		}
		
		// First result should be closer to query
		if results[0].Score < results[1].Score {
			t.Error("Search() results not sorted by score")
		}
	})
	
	t.Run("delete chunk", func(t *testing.T) {
		vs := NewInMemoryVectorStore()
		
		chunk := Chunk{
			ID:        "chunk-1",
			Content:   "Test",
			Embedding: []float32{0.1, 0.2},
		}
		
		_ = vs.Add(ctx, chunk)
		
		if err := vs.Delete(ctx, "chunk-1"); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		
		_, err := vs.Get(ctx, "chunk-1")
		if err == nil {
			t.Error("Get() expected error after delete")
		}
	})
	
	t.Run("clear store", func(t *testing.T) {
		vs := NewInMemoryVectorStore()
		
		chunk := Chunk{
			ID:        "chunk-1",
			Content:   "Test",
			Embedding: []float32{0.1, 0.2},
		}
		
		_ = vs.Add(ctx, chunk)
		
		if err := vs.Clear(ctx); err != nil {
			t.Fatalf("Clear() error = %v", err)
		}
		
		count, _ := vs.Count(ctx)
		if count != 0 {
			t.Errorf("Count() = %d after Clear(), want 0", count)
		}
	})
}

func TestRAG(t *testing.T) {
	ctx := context.Background()
	
	t.Run("add documents and query", func(t *testing.T) {
		ragSystem := New()
		
		docs := []Document{
			NewDocument("doc1", "Go is a programming language designed at Google."),
			NewDocument("doc2", "Python is a high-level programming language."),
		}
		
		if err := ragSystem.AddDocuments(ctx, docs); err != nil {
			t.Fatalf("AddDocuments() error = %v", err)
		}
		
		count, _ := ragSystem.Count(ctx)
		if count == 0 {
			t.Error("Count() = 0 after adding documents")
		}
		
		results, err := ragSystem.Query(ctx, "programming language", 2)
		if err != nil {
			t.Fatalf("Query() error = %v", err)
		}
		
		if len(results) == 0 {
			t.Error("Query() returned no results")
		}
	})
}
