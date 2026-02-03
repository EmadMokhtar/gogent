# Gogent

Gogent is a high-performance, multi-agent framework for Go with built-in Retrieval-Augmented Generation (RAG) capabilities.

## Features

- 🤖 **Intelligent Agents**: Build AI agents with LLM integration
- 🔍 **RAG System**: Semantic search and context injection with vector embeddings
- 💾 **Vector Store**: In-memory vector store with multiple similarity metrics
- 📄 **Document Chunking**: Multiple chunking strategies (fixed-size, sentence, paragraph)
- 🎯 **Metadata Filtering**: Filter search results by metadata
- 🔌 **OpenAI Integration**: Built-in support for OpenAI embeddings and chat models
- 🚀 **High Performance**: Optimized for concurrent access and fast retrieval

## Installation

```bash
go get github.com/EmadMokhtar/gogent
```

## Quick Start

### Basic RAG Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/EmadMokhtar/gogent/llm/providers/openai"
    "github.com/EmadMokhtar/gogent/rag"
)

func main() {
    ctx := context.Background()

    // Create RAG components
    embedder := openai.NewEmbedder(openai.WithAPIKey("your-api-key"))
    store := rag.NewInMemoryVectorStore()
    chunker := rag.NewFixedSizeChunker(
        rag.WithChunkSize(500),
        rag.WithOverlap(50),
    )

    // Create RAG system
    ragSys := rag.New(
        rag.WithVectorStore(store),
        rag.WithEmbedder(embedder),
        rag.WithChunker(chunker),
    )

    // Add documents
    docs := []rag.Document{
        {
            ID:      "doc1",
            Content: "Go is a statically typed, compiled programming language...",
            Metadata: map[string]interface{}{"category": "programming"},
        },
    }
    
    if err := ragSys.AddDocuments(ctx, docs); err != nil {
        log.Fatal(err)
    }

    // Retrieve relevant documents
    results, err := ragSys.Retrieve(ctx, "What is Go?", rag.WithRetrievalK(3))
    if err != nil {
        log.Fatal(err)
    }

    for _, result := range results {
        fmt.Printf("Score: %.3f - %s\n", result.Score, result.Document.Content)
    }
}
```

### Agent with RAG

```go
package main

import (
    "context"
    "fmt"

    "github.com/EmadMokhtar/gogent/agent"
    "github.com/EmadMokhtar/gogent/llm/providers/openai"
    "github.com/EmadMokhtar/gogent/rag"
)

func main() {
    ctx := context.Background()

    // Setup RAG
    embedder := openai.NewEmbedder(openai.WithAPIKey("your-api-key"))
    store := rag.NewInMemoryVectorStore()
    ragSys := rag.New(
        rag.WithVectorStore(store),
        rag.WithEmbedder(embedder),
    )

    // Add knowledge base
    ragSys.AddDocuments(ctx, []rag.Document{
        {ID: "1", Content: "The company was founded in 2020."},
        {ID: "2", Content: "Our main product is an AI assistant."},
    })

    // Create agent with RAG
    llm := openai.New(openai.WithLLMAPIKey("your-api-key"))
    ag := agent.New(
        agent.WithLLM(llm),
        agent.WithRAG(ragSys),
        agent.WithSystemPrompt("You are a helpful assistant."),
    )

    // Query with RAG context
    response, _ := ag.Run(ctx, "When was the company founded?")
    fmt.Println(response.Content)
    
    // Check sources
    for _, src := range response.Sources {
        fmt.Printf("Source: %s (score: %.3f)\n", src.Content, src.Score)
    }
}
```

## Core Components

### RAG System

The RAG system orchestrates document chunking, embedding, storage, and retrieval:

```go
ragSys := rag.New(
    rag.WithVectorStore(store),
    rag.WithEmbedder(embedder),
    rag.WithChunker(chunker),
    rag.WithSimilarityMetric(rag.Cosine),
    rag.WithTopK(5),
)
```

### Vector Store

In-memory vector store with support for multiple similarity metrics:

```go
store := rag.NewInMemoryVectorStore(
    rag.WithMetric(rag.Cosine), // or rag.DotProduct, rag.Euclidean
)
```

**Similarity Metrics:**
- `Cosine`: Cosine similarity (most common for text, range: -1 to 1)
- `DotProduct`: Dot product similarity (faster, good for normalized vectors)
- `Euclidean`: Euclidean distance (L2 distance, converted to similarity)

### Embeddings

OpenAI embedder with caching and batch processing:

```go
embedder := openai.NewEmbedder(
    openai.WithAPIKey("your-api-key"),
    openai.WithModel("text-embedding-ada-002"),
    openai.WithCache(true),
)
```

### Document Chunking

Three chunking strategies available:

**Fixed-Size Chunker:**
```go
chunker := rag.NewFixedSizeChunker(
    rag.WithChunkSize(500),
    rag.WithOverlap(50),
)
```

**Sentence-Based Chunker:**
```go
chunker := rag.NewSentenceChunker(
    rag.WithMaxChunkSize(1000),
)
```

**Paragraph-Based Chunker:**
```go
chunker := rag.NewParagraphChunker(
    rag.WithMaxParagraphChunkSize(2000),
)
```

### Metadata Filtering

Filter search results by metadata:

```go
filter := &rag.Filter{
    Metadata: map[string]interface{}{
        "category": "programming",
        "year":     2023,
    },
}

results, _ := ragSys.Retrieve(ctx, "query", rag.WithFilter(filter))
```

## Examples

See the `examples/` directory for complete working examples:

- `examples/basic_rag/`: Basic RAG operations
- `examples/rag_agent/`: Agent with RAG integration
- `examples/advanced_rag/`: Metadata filtering and custom chunking

To run examples:

```bash
export OPENAI_API_KEY="your-api-key"
go run examples/basic_rag/main.go
```

## API Reference

### RAG Interface

```go
type RAG struct { ... }

// AddDocuments adds documents to the RAG system
func (r *RAG) AddDocuments(ctx context.Context, docs []Document) error

// Retrieve retrieves relevant documents for a query
func (r *RAG) Retrieve(ctx context.Context, query string, opts ...RetrieveOption) ([]SearchResult, error)

// Stats returns statistics about the RAG system
func (r *RAG) Stats(ctx context.Context) (*VectorStats, error)
```

### VectorStore Interface

```go
type VectorStore interface {
    Add(ctx context.Context, docs []Document) error
    Search(ctx context.Context, query []float32, k int, filter *Filter) ([]SearchResult, error)
    Delete(ctx context.Context, ids []string) error
    Get(ctx context.Context, id string) (*Document, error)
    Stats(ctx context.Context) (*VectorStats, error)
}
```

### Embedder Interface

```go
type Embedder interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
    Dimension() int
}
```

### Chunker Interface

```go
type Chunker interface {
    Chunk(ctx context.Context, text string) ([]string, error)
}
```

## Performance Considerations

- **Batch Processing**: Embeddings are batched (up to 100 at a time) for efficiency
- **Caching**: Optional embedding cache to avoid re-computation
- **Concurrent Access**: Thread-safe operations with read/write locks
- **Brute-Force Search**: Efficient for up to ~10,000 documents
- **Memory Usage**: All vectors stored in memory for fast access

## Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
