package rag

import (
	"context"
	"sync"
	"testing"
)

// TestInMemoryVectorStore_ConcurrentAccess tests thread safety
func TestInMemoryVectorStore_ConcurrentAccess(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	// Number of concurrent goroutines
	numGoroutines := 10
	numOperationsPerGoroutine := 20

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperationsPerGoroutine; j++ {
				doc := Document{
					ID:        string(rune('a'+id)) + string(rune('0'+j)),
					Content:   "content",
					Embedding: []float32{float32(id), float32(j), 1.0},
				}
				_ = store.Add(ctx, []Document{doc})
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	stats, err := store.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats() failed: %v", err)
	}

	expectedDocs := numGoroutines * numOperationsPerGoroutine
	if stats.TotalDocuments != expectedDocs {
		t.Errorf("Expected %d documents after concurrent writes, got %d",
			expectedDocs, stats.TotalDocuments)
	}

	// Concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			query := []float32{1, 1, 1}
			_, err := store.Search(ctx, query, 5, nil)
			if err != nil {
				t.Errorf("Concurrent Search() failed: %v", err)
			}
		}()
	}

	wg.Wait()
}

// TestInMemoryVectorStore_ConcurrentReadWrite tests concurrent reads and writes
func TestInMemoryVectorStore_ConcurrentReadWrite(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	// Add initial documents
	for i := 0; i < 10; i++ {
		doc := Document{
			ID:        string(rune('a' + i)),
			Content:   "content",
			Embedding: []float32{float32(i), 1.0, 1.0},
		}
		_ = store.Add(ctx, []Document{doc})
	}

	var wg sync.WaitGroup
	numReaders := 5
	numWriters := 5

	// Start readers
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			query := []float32{1, 1, 1}
			for j := 0; j < 50; j++ {
				_, err := store.Search(ctx, query, 3, nil)
				if err != nil {
					t.Errorf("Concurrent read failed: %v", err)
				}
			}
		}()
	}

	// Start writers
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				doc := Document{
					ID:        string(rune('A'+id)) + string(rune('0'+(j%10))),
					Content:   "new content",
					Embedding: []float32{float32(id), float32(j), 2.0},
				}
				_ = store.Add(ctx, []Document{doc})
			}
		}(i)
	}

	wg.Wait()

	// Verify store is still functional
	stats, err := store.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats() after concurrent access failed: %v", err)
	}

	if stats.TotalDocuments == 0 {
		t.Error("Expected documents in store after concurrent operations")
	}
}

// TestInMemoryVectorStore_ConcurrentDelete tests concurrent deletions
func TestInMemoryVectorStore_ConcurrentDelete(t *testing.T) {
	store := NewInMemoryVectorStore()
	ctx := context.Background()

	// Add documents
	for i := 0; i < 100; i++ {
		doc := Document{
			ID:        string(rune('a' + (i % 26))) + string(rune('0' + (i / 26))),
			Content:   "content",
			Embedding: []float32{float32(i), 1.0, 1.0},
		}
		_ = store.Add(ctx, []Document{doc})
	}

	var wg sync.WaitGroup
	numDeleters := 5

	// Concurrent deletions
	wg.Add(numDeleters)
	for i := 0; i < numDeleters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				docID := string(rune('a'+(id*4+j)%26)) + string(rune('0'+((id*4+j)/26)))
				_ = store.Delete(ctx, []string{docID})
			}
		}(i)
	}

	wg.Wait()

	// Store should still be functional
	_, err := store.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats() after concurrent deletes failed: %v", err)
	}
}
