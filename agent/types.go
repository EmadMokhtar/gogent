// Package agent provides the core agent implementation for the Gogent framework.
package agent

import (
	"context"

	"github.com/EmadMokhtar/gogent/types"
)

// RunContext holds the context for a single agent run.
type RunContext struct {
	ctx              context.Context
	messages         []types.Message
	systemPrompt     string
	maxIterations    int
	currentIteration int
}

// NewRunContext creates a new run context.
func NewRunContext(ctx context.Context, systemPrompt string, maxIterations int) *RunContext {
	return &RunContext{
		ctx:              ctx,
		messages:         make([]types.Message, 0),
		systemPrompt:     systemPrompt,
		maxIterations:    maxIterations,
		currentIteration: 0,
	}
}

// AddMessage adds a message to the context.
func (rc *RunContext) AddMessage(msg types.Message) {
	rc.messages = append(rc.messages, msg)
}

// Messages returns all messages in the context.
func (rc *RunContext) Messages() []types.Message {
	if rc.systemPrompt != "" {
		// Prepend system prompt
		return append([]types.Message{{Role: "system", Content: rc.systemPrompt}}, rc.messages...)
	}
	return rc.messages
}

// Context returns the underlying context.Context.
func (rc *RunContext) Context() context.Context {
	return rc.ctx
}

// IncrementIteration increments the iteration counter.
func (rc *RunContext) IncrementIteration() {
	rc.currentIteration++
}

// CanIterate checks if we can perform another iteration.
func (rc *RunContext) CanIterate() bool {
	return rc.currentIteration < rc.maxIterations
}

// CurrentIteration returns the current iteration number.
func (rc *RunContext) CurrentIteration() int {
	return rc.currentIteration
}
