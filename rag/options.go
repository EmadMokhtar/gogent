package rag

// RAGOption defines configuration options for RAG
type RAGOption func(*RAG)

// WithVectorStore sets the vector store
func WithVectorStore(store VectorStore) RAGOption {
	return func(r *RAG) {
		r.vectorStore = store
	}
}

// WithEmbedder sets the embedder
func WithEmbedder(embedder Embedder) RAGOption {
	return func(r *RAG) {
		r.embedder = embedder
	}
}

// WithChunker sets the chunker
func WithChunker(chunker Chunker) RAGOption {
	return func(r *RAG) {
		r.chunker = chunker
	}
}

// WithSimilarityMetric sets the similarity metric (for vector store)
func WithSimilarityMetric(metric SimilarityMetric) RAGOption {
	return func(r *RAG) {
		r.metric = metric
	}
}

// WithTopK sets the default number of results to retrieve
func WithTopK(k int) RAGOption {
	return func(r *RAG) {
		r.topK = k
	}
}

// RetrieveOption defines options for retrieval
type RetrieveOption func(*retrieveOptions)

// WithFilter sets a metadata filter for retrieval
func WithFilter(filter *Filter) RetrieveOption {
	return func(o *retrieveOptions) {
		o.filter = filter
	}
}

// WithRetrievalK overrides the default top-k for a specific retrieval
func WithRetrievalK(k int) RetrieveOption {
	return func(o *retrieveOptions) {
		o.k = k
	}
}

type retrieveOptions struct {
	filter *Filter
	k      int
}
