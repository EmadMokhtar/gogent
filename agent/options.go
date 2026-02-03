package agent

import (
	"github.com/EmadMokhtar/gogent/llm"
	"github.com/EmadMokhtar/gogent/memory"
)

// Option is a function that configures an Agent.
type Option func(*Agent)

// WithLLM sets the LLM for the agent.
func WithLLM(l llm.LLM) Option {
	return func(a *Agent) {
		a.llm = l
	}
}

// WithMemory sets the memory for the agent.
func WithMemory(m memory.Memory) Option {
	return func(a *Agent) {
		a.memory = m
	}
}

// WithSystemPrompt sets the system prompt for the agent.
func WithSystemPrompt(prompt string) Option {
	return func(a *Agent) {
		a.systemPrompt = prompt
	}
}
