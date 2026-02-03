package rag

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrNoVectorStore is returned when RAG is used without a vector store
	ErrNoVectorStore = errors.New("vector store is required")
	// ErrNoEmbedder is returned when RAG is used without an embedder
	ErrNoEmbedder = errors.New("embedder is required")
)

// RAG orchestrates the retrieval-augmented generation workflow
type RAG struct {
	vectorStore VectorStore
	embedder    Embedder
	chunker     Chunker
	metric      SimilarityMetric
	topK        int
}

// New creates a new RAG system
func New(opts ...RAGOption) *RAG {
	rag := &RAG{
		metric: Cosine,
		topK:   3,
	}

	for _, opt := range opts {
		opt(rag)
	}

	return rag
}

// AddDocuments adds documents to the RAG system (chunks, embeds, and stores)
func (r *RAG) AddDocuments(ctx context.Context, docs []Document) error {
	if r.vectorStore == nil {
		return ErrNoVectorStore
	}
	if r.embedder == nil {
		return ErrNoEmbedder
	}

	var processedDocs []Document

	for _, doc := range docs {
		// If document has no content, skip it
		if doc.Content == "" {
			continue
		}

		// Chunk the document if chunker is provided
		if r.chunker != nil {
			chunks, err := r.chunker.Chunk(ctx, doc.Content)
			if err != nil {
				return fmt.Errorf("failed to chunk document %s: %w", doc.ID, err)
			}

			// Generate embeddings for chunks
			embeddings, err := r.embedder.EmbedBatch(ctx, chunks)
			if err != nil {
				return fmt.Errorf("failed to embed chunks for document %s: %w", doc.ID, err)
			}

			// Create chunk documents
			for i, chunk := range chunks {
				chunkDoc := Document{
					ID:         fmt.Sprintf("%s_chunk_%d", doc.ID, i),
					Content:    chunk,
					Embedding:  embeddings[i],
					Metadata:   doc.Metadata,
					ChunkIndex: i,
					Source:     doc.ID,
				}
				processedDocs = append(processedDocs, chunkDoc)
			}
		} else {
			// No chunking, embed the whole document
			embedding, err := r.embedder.Embed(ctx, doc.Content)
			if err != nil {
				return fmt.Errorf("failed to embed document %s: %w", doc.ID, err)
			}

			doc.Embedding = embedding
			processedDocs = append(processedDocs, doc)
		}
	}

	// Store all processed documents
	if err := r.vectorStore.Add(ctx, processedDocs); err != nil {
		return fmt.Errorf("failed to store documents: %w", err)
	}

	return nil
}

// Retrieve retrieves relevant documents for a query
func (r *RAG) Retrieve(ctx context.Context, query string, opts ...RetrieveOption) ([]SearchResult, error) {
	if r.vectorStore == nil {
		return nil, ErrNoVectorStore
	}
	if r.embedder == nil {
		return nil, ErrNoEmbedder
	}

	// Parse options
	options := retrieveOptions{
		k: r.topK,
	}
	for _, opt := range opts {
		opt(&options)
	}

	// Embed the query
	queryEmbedding, err := r.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// Search the vector store
	results, err := r.vectorStore.Search(ctx, queryEmbedding, options.k, options.filter)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return results, nil
}

// Stats returns statistics about the RAG system
func (r *RAG) Stats(ctx context.Context) (*VectorStats, error) {
	if r.vectorStore == nil {
		return nil, ErrNoVectorStore
	}

	return r.vectorStore.Stats(ctx)
}
