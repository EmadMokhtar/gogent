package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EmadMokhtar/gogent/llm/providers/openai"
	"github.com/EmadMokhtar/gogent/rag"
)

func main() {
	ctx := context.Background()

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create RAG system with custom chunker
	embedder := openai.NewEmbedder(openai.WithAPIKey(apiKey))
	store := rag.NewInMemoryVectorStore(rag.WithMetric(rag.Cosine))
	
	// Use sentence-based chunker for better semantic units
	chunker := rag.NewSentenceChunker(rag.WithMaxChunkSize(300))

	ragSys := rag.New(
		rag.WithVectorStore(store),
		rag.WithEmbedder(embedder),
		rag.WithChunker(chunker),
		rag.WithTopK(5),
	)

	// Add documents with rich metadata
	docs := []rag.Document{
		{
			ID: "python_basics",
			Content: `Python is a high-level, interpreted programming language. 
It emphasizes code readability with significant whitespace. 
Python supports multiple programming paradigms including procedural, object-oriented, and functional programming.
It was created by Guido van Rossum and first released in 1991.`,
			Metadata: map[string]interface{}{
				"language":   "python",
				"topic":      "basics",
				"difficulty": "beginner",
				"year":       1991,
			},
		},
		{
			ID: "python_data",
			Content: `Python is widely used for data science and machine learning. 
Libraries like NumPy, Pandas, and Scikit-learn make it powerful for data analysis.
TensorFlow and PyTorch are popular deep learning frameworks in Python.`,
			Metadata: map[string]interface{}{
				"language":   "python",
				"topic":      "data_science",
				"difficulty": "intermediate",
			},
		},
		{
			ID: "go_basics",
			Content: `Go (Golang) is a compiled, statically typed language developed at Google.
It features built-in concurrency support with goroutines and channels.
Go has a simple syntax and fast compilation times.
It's designed for building scalable and efficient systems.`,
			Metadata: map[string]interface{}{
				"language":   "go",
				"topic":      "basics",
				"difficulty": "beginner",
				"year":       2009,
			},
		},
		{
			ID: "go_web",
			Content: `Go is excellent for building web services and APIs.
The standard library includes a powerful HTTP server.
Popular frameworks include Gin, Echo, and Fiber for web development.`,
			Metadata: map[string]interface{}{
				"language":   "go",
				"topic":      "web",
				"difficulty": "intermediate",
			},
		},
	}

	fmt.Println("Adding documents to RAG system...")
	if err := ragSys.AddDocuments(ctx, docs); err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}

	stats, err := ragSys.Stats(ctx)
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}
	fmt.Printf("Vector store stats: %d documents (chunks), dimension: %d\n\n", stats.TotalDocuments, stats.Dimension)

	// Example 1: Search without filters
	fmt.Println("=== Search without filters ===")
	query1 := "Tell me about web development"
	results, err := ragSys.Retrieve(ctx, query1, rag.WithRetrievalK(3))
	if err != nil {
		log.Fatalf("Failed to retrieve: %v", err)
	}

	fmt.Printf("Query: %s\n", query1)
	for _, result := range results {
		fmt.Printf("  [Score: %.3f, Lang: %v] %s\n",
			result.Score,
			result.Document.Metadata["language"],
			truncate(result.Document.Content, 60))
	}
	fmt.Println()

	// Example 2: Search with language filter
	fmt.Println("=== Search with filter (Go only) ===")
	query2 := "programming language features"
	filter := &rag.Filter{
		Metadata: map[string]interface{}{
			"language": "go",
		},
	}
	results, err = ragSys.Retrieve(ctx, query2, rag.WithFilter(filter), rag.WithRetrievalK(3))
	if err != nil {
		log.Fatalf("Failed to retrieve: %v", err)
	}

	fmt.Printf("Query: %s (filtered by language=go)\n", query2)
	for _, result := range results {
		fmt.Printf("  [Score: %.3f, Topic: %v] %s\n",
			result.Score,
			result.Document.Metadata["topic"],
			truncate(result.Document.Content, 60))
	}
	fmt.Println()

	// Example 3: Search with multiple filters
	fmt.Println("=== Search with multiple filters ===")
	query3 := "beginner friendly languages"
	filter2 := &rag.Filter{
		Metadata: map[string]interface{}{
			"topic":      "basics",
			"difficulty": "beginner",
		},
	}
	results, err = ragSys.Retrieve(ctx, query3, rag.WithFilter(filter2), rag.WithRetrievalK(3))
	if err != nil {
		log.Fatalf("Failed to retrieve: %v", err)
	}

	fmt.Printf("Query: %s (filtered by topic=basics, difficulty=beginner)\n", query3)
	for _, result := range results {
		fmt.Printf("  [Score: %.3f, Lang: %v, Year: %v] %s\n",
			result.Score,
			result.Document.Metadata["language"],
			result.Document.Metadata["year"],
			truncate(result.Document.Content, 60))
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
