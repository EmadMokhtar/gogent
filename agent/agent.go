package agent

import (
	"context"
	"fmt"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/rag"
)

// Agent represents an AI agent
type Agent struct {
	llm          llm.LLM
	rag          *rag.RAG
	systemPrompt string
	ragTopK      int
}

// Response represents an agent response with sources
type Response struct {
	Content string
	Sources []rag.SearchResult
	Model   string
}

// Option defines configuration options for Agent
type Option func(*Agent)

// WithLLM sets the language model
func WithLLM(llmProvider llm.LLM) Option {
	return func(a *Agent) {
		a.llm = llmProvider
	}
}

// WithRAG sets the RAG system
func WithRAG(ragSystem *rag.RAG) Option {
	return func(a *Agent) {
		a.rag = ragSystem
	}
}

// WithSystemPrompt sets the system prompt
func WithSystemPrompt(prompt string) Option {
	return func(a *Agent) {
		a.systemPrompt = prompt
	}
}

// WithRAGTopK sets the number of documents to retrieve from RAG
func WithRAGTopK(k int) Option {
	return func(a *Agent) {
		a.ragTopK = k
	}
}

// New creates a new agent
func New(opts ...Option) *Agent {
	agent := &Agent{
		ragTopK: 3,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

// Run executes the agent with the given query
func (a *Agent) Run(ctx context.Context, query string) (*Response, error) {
	if a.llm == nil {
		return nil, fmt.Errorf("LLM is required")
	}

	var messages []llm.Message
	systemPrompt := a.systemPrompt

	// If RAG is enabled, retrieve relevant context
	var sources []rag.SearchResult
	if a.rag != nil {
		results, err := a.rag.Retrieve(ctx, query, rag.WithRetrievalK(a.ragTopK))
		if err != nil {
			// Log error but continue without RAG context
			fmt.Printf("Warning: RAG retrieval failed: %v\n", err)
		} else if len(results) > 0 {
			sources = results

			// Build context from retrieved documents
			contextStr := "\n\nRelevant information:\n"
			for _, result := range results {
				contextStr += fmt.Sprintf("- %s\n", result.Document.Content)
			}

			// Inject context into system prompt
			systemPrompt += contextStr
		}
	}

	// Add system message if present
	if systemPrompt != "" {
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	// Add user query
	messages = append(messages, llm.Message{
		Role:    "user",
		Content: query,
	})

	// Generate response
	llmResponse, err := a.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	response := &Response{
		Content: llmResponse.Content,
		Sources: sources,
		Model:   llmResponse.Model,
	}

	return response, nil
}
