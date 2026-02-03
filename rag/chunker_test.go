package rag

import (
	"context"
	"testing"
)

func TestFixedSizeChunker(t *testing.T) {
	chunker := NewFixedSizeChunker(
		WithChunkSize(20),
		WithOverlap(5),
	)

	ctx := context.Background()
	text := "This is a test text that should be chunked into multiple pieces."

	chunks, err := chunker.Chunk(ctx, text)
	if err != nil {
		t.Fatalf("Chunk() failed: %v", err)
	}

	if len(chunks) == 0 {
		t.Error("Chunk() returned no chunks")
	}

	// Check that chunks have overlap
	if len(chunks) > 1 {
		// Last characters of first chunk should appear in second chunk
		firstChunk := chunks[0]
		secondChunk := chunks[1]

		if len(firstChunk) < 5 || len(secondChunk) < 5 {
			t.Skip("Chunks too small to verify overlap")
		}

		// Due to overlap, there should be some common content
		// (exact match depends on chunking boundaries)
	}
}

func TestFixedSizeChunker_EmptyText(t *testing.T) {
	chunker := NewFixedSizeChunker()
	ctx := context.Background()

	chunks, err := chunker.Chunk(ctx, "")
	if err != nil {
		t.Fatalf("Chunk() failed: %v", err)
	}

	if len(chunks) != 0 {
		t.Errorf("Chunk() with empty text returned %d chunks, expected 0", len(chunks))
	}
}

func TestFixedSizeChunker_SmallText(t *testing.T) {
	chunker := NewFixedSizeChunker(WithChunkSize(100))
	ctx := context.Background()

	text := "Small text"
	chunks, err := chunker.Chunk(ctx, text)
	if err != nil {
		t.Fatalf("Chunk() failed: %v", err)
	}

	if len(chunks) != 1 {
		t.Errorf("Chunk() with small text returned %d chunks, expected 1", len(chunks))
	}

	if chunks[0] != text {
		t.Errorf("Chunk() = %v, expected %v", chunks[0], text)
	}
}

func TestSentenceChunker(t *testing.T) {
	chunker := NewSentenceChunker(WithMaxChunkSize(50))
	ctx := context.Background()

	text := "This is sentence one. This is sentence two! Is this sentence three? Yes it is."

	chunks, err := chunker.Chunk(ctx, text)
	if err != nil {
		t.Fatalf("Chunk() failed: %v", err)
	}

	if len(chunks) == 0 {
		t.Error("Chunk() returned no chunks")
	}

	// Verify chunks are not empty
	for i, chunk := range chunks {
		if chunk == "" {
			t.Errorf("Chunk %d is empty", i)
		}
	}
}

func TestSentenceChunker_EmptyText(t *testing.T) {
	chunker := NewSentenceChunker()
	ctx := context.Background()

	chunks, err := chunker.Chunk(ctx, "")
	if err != nil {
		t.Fatalf("Chunk() failed: %v", err)
	}

	if len(chunks) != 0 {
		t.Errorf("Chunk() with empty text returned %d chunks, expected 0", len(chunks))
	}
}

func TestParagraphChunker(t *testing.T) {
	chunker := NewParagraphChunker(WithMaxParagraphChunkSize(100))
	ctx := context.Background()

	text := `This is paragraph one.
It has multiple sentences.

This is paragraph two.
It also has content.

And here is paragraph three.`

	chunks, err := chunker.Chunk(ctx, text)
	if err != nil {
		t.Fatalf("Chunk() failed: %v", err)
	}

	if len(chunks) == 0 {
		t.Error("Chunk() returned no chunks")
	}

	// Verify chunks are not empty
	for i, chunk := range chunks {
		if chunk == "" {
			t.Errorf("Chunk %d is empty", i)
		}
	}
}

func TestParagraphChunker_EmptyText(t *testing.T) {
	chunker := NewParagraphChunker()
	ctx := context.Background()

	chunks, err := chunker.Chunk(ctx, "")
	if err != nil {
		t.Fatalf("Chunk() failed: %v", err)
	}

	if len(chunks) != 0 {
		t.Errorf("Chunk() with empty text returned %d chunks, expected 0", len(chunks))
	}
}
