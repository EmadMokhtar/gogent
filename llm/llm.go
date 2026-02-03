package llm

import (
	"context"

	"github.com/EmadMokhtar/gogent/memory"
)

// LLM defines the interface for language model providers.
type LLM interface {
	// Generate produces a response based on the message history.
	Generate(ctx context.Context, messages []memory.Message) (string, error)
}
