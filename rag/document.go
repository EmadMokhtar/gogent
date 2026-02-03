package rag

// Document represents a document in the vector store
type Document struct {
	// ID is the unique identifier for the document
	ID string
	// Content is the text content of the document
	Content string
	// Embedding is the vector representation of the content
	Embedding []float32
	// Metadata contains additional information about the document
	Metadata map[string]interface{}
	// ChunkIndex indicates the position if this document is part of a larger document
	ChunkIndex int
	// Source is the original document ID if this is a chunk
	Source string
}

// SearchResult represents a document returned from a search query
type SearchResult struct {
	// Document is the retrieved document
	Document Document
	// Score is the similarity score (higher is better)
	Score float64
	// Rank is the position in the result set (1-indexed)
	Rank int
}

// VectorStats contains statistics about the vector store
type VectorStats struct {
	// TotalDocuments is the number of documents in the store
	TotalDocuments int
	// Dimension is the embedding dimension
	Dimension int
	// MemoryUsage is an estimate of memory usage in bytes
	MemoryUsage int64
}
