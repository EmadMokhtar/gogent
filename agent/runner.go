package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/memory"
	"github.com/EmadMokhtar/gogent/rag"
)

// Runner executes agent workflows
type Runner struct {
	agent *Agent
}

// NewRunner creates a new agent runner
func NewRunner(agent *Agent) *Runner {
	return &Runner{
		agent: agent,
	}
}

// Run executes the agent with a query
func (r *Runner) Run(ctx context.Context, query string) (*Result, error) {
	return r.agent.Run(ctx, query)
}

// RunWithContext executes the agent with a custom context
func (r *Runner) RunWithContext(agentCtx *Context, query string) (*Result, error) {
	return r.agent.RunWithContext(agentCtx, query)
}

// Stream executes the agent with streaming
func (r *Runner) Stream(ctx context.Context, query string) (<-chan StreamChunk, error) {
	return r.agent.Stream(ctx, query)
}

// StreamChunk represents a chunk of streaming output
type StreamChunk struct {
	Content string
	Done    bool
	Error   error
	Result  *Result
}

// Result represents the result of an agent execution
type Result struct {
	Response  string                 `json:"response"`
	Sources   []rag.RetrievalResult  `json:"sources,omitempty"`
	ToolCalls []ToolCallResult       `json:"tool_calls,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Duration  time.Duration          `json:"duration"`
}

// ToolCallResult represents the result of a tool call
type ToolCallResult struct {
	ToolName string                 `json:"tool_name"`
	Input    map[string]interface{} `json:"input"`
	Output   interface{}            `json:"output"`
	Error    string                 `json:"error,omitempty"`
}

// prepareMessages prepares messages for LLM call
func (r *Runner) prepareMessages(ctx context.Context, query string) ([]llm.Message, error) {
	messages := make([]llm.Message, 0)
	
	// Add system prompt
	if r.agent.opts.SystemPrompt != "" {
		messages = append(messages, llm.NewSystemMessage(r.agent.opts.SystemPrompt))
	}
	
	// Add memory if available
	if r.agent.opts.Memory != nil {
		memoryMessages, err := r.agent.opts.Memory.Get(ctx, 10)
		if err == nil && len(memoryMessages) > 0 {
			for _, msg := range memoryMessages {
				messages = append(messages, llm.Message{
					Role:      llm.Role(msg.Role),
					Content:   msg.Content,
					Timestamp: msg.Timestamp,
				})
			}
		}
	}
	
	// Add RAG context if available
	var ragContext string
	if r.agent.opts.RAG != nil {
		results, err := r.agent.opts.RAG.Query(ctx, query, 3)
		if err == nil && len(results) > 0 {
			ragContext = "\n\nRelevant context:\n"
			for i, result := range results {
				ragContext += fmt.Sprintf("%d. %s\n", i+1, result.Chunk.Content)
			}
		}
	}
	
	// Add user query
	userContent := query
	if ragContext != "" {
		userContent = query + ragContext
	}
	messages = append(messages, llm.NewUserMessage(userContent))
	
	return messages, nil
}

// saveToMemory saves messages to memory
func (r *Runner) saveToMemory(ctx context.Context, userQuery, assistantResponse string) error {
	if r.agent.opts.Memory == nil {
		return nil
	}
	
	// Save user message
	if err := r.agent.opts.Memory.Add(ctx, memory.NewUserMessage(userQuery)); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}
	
	// Save assistant message
	if err := r.agent.opts.Memory.Add(ctx, memory.NewAssistantMessage(assistantResponse)); err != nil {
		return fmt.Errorf("failed to save assistant message: %w", err)
	}
	
	return nil
}
