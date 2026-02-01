package rag

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

// RAG orchestrates the entire RAG pipeline
type RAG struct {
	vectorStore VectorStore
	embedder    Embedder
	chunker     *Chunker
	retriever   *Retriever
}

// RAGOption is a functional option for configuring RAG
type RAGOption func(*RAG)

// New creates a new RAG instance
func New(opts ...RAGOption) *RAG {
	r := &RAG{
		vectorStore: NewInMemoryVectorStore(),
		embedder:    NewMockEmbedder(1536),
		chunker:     NewChunker(ChunkerConfig{ChunkSize: 500, ChunkOverlap: 50}),
	}
	
	for _, opt := range opts {
		opt(r)
	}
	
	// Create retriever
	r.retriever = NewRetriever(RetrieverConfig{
		VectorStore:      r.vectorStore,
		Embedder:         r.embedder,
		TopK:             5,
		SimilarityMetric: CosineSimilarity,
	})
	
	return r
}

// WithVectorStore sets the vector store
func WithVectorStore(vs VectorStore) RAGOption {
	return func(r *RAG) {
		r.vectorStore = vs
	}
}

// WithEmbedder sets the embedder
func WithEmbedder(e Embedder) RAGOption {
	return func(r *RAG) {
		r.embedder = e
	}
}

// WithChunkSize sets the chunk size
func WithChunkSize(size int) RAGOption {
	return func(r *RAG) {
		if r.chunker == nil {
			r.chunker = NewChunker(ChunkerConfig{ChunkSize: size})
		} else {
			r.chunker.config.ChunkSize = size
		}
	}
}

// WithChunkOverlap sets the chunk overlap
func WithChunkOverlap(overlap int) RAGOption {
	return func(r *RAG) {
		if r.chunker == nil {
			r.chunker = NewChunker(ChunkerConfig{ChunkOverlap: overlap})
		} else {
			r.chunker.config.ChunkOverlap = overlap
		}
	}
}

// AddDocuments adds documents to the RAG system
func (r *RAG) AddDocuments(ctx context.Context, docs []Document) error {
	if len(docs) == 0 {
		return fmt.Errorf("no documents provided")
	}
	
	// Chunk documents
	chunks, err := r.chunker.ChunkDocuments(docs)
	if err != nil {
		return fmt.Errorf("failed to chunk documents: %w", err)
	}
	
	if len(chunks) == 0 {
		return fmt.Errorf("no chunks generated from documents")
	}
	
	// Generate embeddings for chunks in parallel
	embeddings, err := r.generateEmbeddingsParallel(ctx, chunks)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}
	
	// Attach embeddings to chunks
	for i := range chunks {
		chunks[i].Embedding = embeddings[i]
	}
	
	// Add chunks to vector store
	if err := r.vectorStore.AddBatch(ctx, chunks); err != nil {
		return fmt.Errorf("failed to add chunks to vector store: %w", err)
	}
	
	return nil
}

// generateEmbeddingsParallel generates embeddings for chunks in parallel
func (r *RAG) generateEmbeddingsParallel(ctx context.Context, chunks []Chunk) ([][]float32, error) {
	embeddings := make([][]float32, len(chunks))
	var mu sync.Mutex
	
	g, ctx := errgroup.WithContext(ctx)
	
	// Process chunks in batches
	batchSize := 10
	for i := 0; i < len(chunks); i += batchSize {
		end := i + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}
		
		batch := chunks[i:end]
		batchStart := i
		
		g.Go(func() error {
			texts := make([]string, len(batch))
			for j, chunk := range batch {
				texts[j] = chunk.Content
			}
			
			batchEmbeddings, err := r.embedder.EmbedBatch(ctx, texts)
			if err != nil {
				return err
			}
			
			mu.Lock()
			for j, emb := range batchEmbeddings {
				embeddings[batchStart+j] = emb
			}
			mu.Unlock()
			
			return nil
		})
	}
	
	if err := g.Wait(); err != nil {
		return nil, err
	}
	
	return embeddings, nil
}

// Query performs a RAG query and returns relevant chunks
func (r *RAG) Query(ctx context.Context, query string, topK int) ([]RetrievalResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}
	
	if topK <= 0 {
		topK = 5
	}
	
	// Update retriever topK
	r.retriever.topK = topK
	
	// Retrieve relevant chunks
	results, err := r.retriever.Retrieve(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve: %w", err)
	}
	
	return results, nil
}

// QueryWithFilters performs a RAG query with metadata filters
func (r *RAG) QueryWithFilters(ctx context.Context, query string, topK int, filters map[string]interface{}) ([]RetrievalResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}
	
	if topK <= 0 {
		topK = 5
	}
	
	// Update retriever topK
	r.retriever.topK = topK
	
	// Retrieve relevant chunks with filters
	results, err := r.retriever.Retrieve(ctx, query, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve: %w", err)
	}
	
	return results, nil
}

// Clear clears all documents from the RAG system
func (r *RAG) Clear(ctx context.Context) error {
	return r.vectorStore.Clear(ctx)
}

// Count returns the number of chunks in the RAG system
func (r *RAG) Count(ctx context.Context) (int, error) {
	return r.vectorStore.Count(ctx)
}
