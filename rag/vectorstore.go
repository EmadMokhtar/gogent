package rag

import "context"

// VectorStore defines the interface for storing and retrieving document embeddings
type VectorStore interface {
	// Add stores documents with their embeddings in the vector store
	Add(ctx context.Context, docs []Document) error

	// Search finds the k most similar documents to the query embedding
	Search(ctx context.Context, query []float32, k int, filter *Filter) ([]SearchResult, error)

	// Delete removes documents by their IDs
	Delete(ctx context.Context, ids []string) error

	// Get retrieves a document by its ID
	Get(ctx context.Context, id string) (*Document, error)

	// Stats returns statistics about the vector store
	Stats(ctx context.Context) (*VectorStats, error)
}
