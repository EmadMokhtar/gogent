package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/rag"
)

// MockLLM is a mock LLM for testing
type MockLLM struct {
	model string
}

func (m *MockLLM) Generate(ctx context.Context, messages []llm.Message) (*llm.Response, error) {
	if len(messages) == 0 {
		return nil, errors.New("no messages")
	}

	// Simple mock: echo back the last user message
	lastMessage := messages[len(messages)-1]
	return &llm.Response{
		Content: "Response to: " + lastMessage.Content,
		Model:   m.model,
	}, nil
}

func (m *MockLLM) Model() string {
	return m.model
}

func TestAgent_New(t *testing.T) {
	llm := &MockLLM{model: "test-model"}
	ag := New(
		WithLLM(llm),
		WithSystemPrompt("test prompt"),
		WithRAGTopK(5),
	)

	if ag.llm == nil {
		t.Error("Expected LLM to be set")
	}

	if ag.systemPrompt != "test prompt" {
		t.Errorf("systemPrompt = %s, expected 'test prompt'", ag.systemPrompt)
	}

	if ag.ragTopK != 5 {
		t.Errorf("ragTopK = %d, expected 5", ag.ragTopK)
	}
}

func TestAgent_Run_NoLLM(t *testing.T) {
	ag := New() // No LLM
	ctx := context.Background()

	_, err := ag.Run(ctx, "test query")
	if err == nil {
		t.Error("Run() without LLM should return error")
	}
}

func TestAgent_Run_WithLLM(t *testing.T) {
	llm := &MockLLM{model: "test-model"}
	ag := New(WithLLM(llm))
	ctx := context.Background()

	response, err := ag.Run(ctx, "test query")
	if err != nil {
		t.Fatalf("Run() failed: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Content == "" {
		t.Error("Expected response content to be non-empty")
	}

	if response.Model != "test-model" {
		t.Errorf("Response model = %s, expected 'test-model'", response.Model)
	}
}

func TestAgent_Run_WithSystemPrompt(t *testing.T) {
	llm := &MockLLM{model: "test-model"}
	ag := New(
		WithLLM(llm),
		WithSystemPrompt("You are a helpful assistant."),
	)
	ctx := context.Background()

	response, err := ag.Run(ctx, "test query")
	if err != nil {
		t.Fatalf("Run() failed: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}
}

// MockEmbedder for testing RAG integration
type MockEmbedder struct {
	dimension int
}

func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	vec := make([]float32, m.dimension)
	for i := range vec {
		vec[i] = float32(len(text)) / float32(m.dimension)
	}
	return vec, nil
}

func (m *MockEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var embeddings [][]float32
	for _, text := range texts {
		emb, err := m.Embed(ctx, text)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, emb)
	}
	return embeddings, nil
}

func (m *MockEmbedder) Dimension() int {
	return m.dimension
}

func TestAgent_Run_WithRAG(t *testing.T) {
	ctx := context.Background()

	// Setup RAG
	embedder := &MockEmbedder{dimension: 3}
	store := rag.NewInMemoryVectorStore()
	ragSys := rag.New(
		rag.WithVectorStore(store),
		rag.WithEmbedder(embedder),
	)

	// Add documents
	docs := []rag.Document{
		{
			ID:      "doc1",
			Content: "The company was founded in 2020.",
		},
		{
			ID:      "doc2",
			Content: "We have 50 employees.",
		},
	}

	err := ragSys.AddDocuments(ctx, docs)
	if err != nil {
		t.Fatalf("AddDocuments() failed: %v", err)
	}

	// Create agent with RAG
	llm := &MockLLM{model: "test-model"}
	ag := New(
		WithLLM(llm),
		WithRAG(ragSys),
		WithSystemPrompt("You are a helpful assistant."),
		WithRAGTopK(2),
	)

	// Run query
	response, err := ag.Run(ctx, "When was the company founded?")
	if err != nil {
		t.Fatalf("Run() with RAG failed: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	// Check sources
	if len(response.Sources) == 0 {
		t.Error("Expected sources from RAG, got none")
	}
}

func TestAgent_Run_RAGRetrievalError(t *testing.T) {
	ctx := context.Background()

	// Setup RAG without embedder (will cause error)
	store := rag.NewInMemoryVectorStore()
	ragSys := rag.New(
		rag.WithVectorStore(store),
		// No embedder - will cause error
	)

	// Create agent with broken RAG
	llm := &MockLLM{model: "test-model"}
	ag := New(
		WithLLM(llm),
		WithRAG(ragSys),
	)

	// Should not fail the agent run, just skip RAG
	response, err := ag.Run(ctx, "test query")
	if err != nil {
		t.Fatalf("Run() failed: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response even with RAG error")
	}

	// Sources should be empty due to RAG error
	if len(response.Sources) != 0 {
		t.Errorf("Expected no sources due to RAG error, got %d", len(response.Sources))
	}
}
