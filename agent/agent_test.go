// Package agent provides the core agent implementation for the Gogent framework.
package agent

import (
	"context"
	"testing"

	"github.com/EmadMokhtar/gogent/types"
)

func TestNewRunContext(t *testing.T) {
	ctx := context.Background()
	runCtx := NewRunContext(ctx, "Test prompt", 5)
	
	if runCtx.systemPrompt != "Test prompt" {
		t.Errorf("expected system prompt 'Test prompt', got %s", runCtx.systemPrompt)
	}
	
	if runCtx.maxIterations != 5 {
		t.Errorf("expected max iterations 5, got %d", runCtx.maxIterations)
	}
}

func TestRunContext_AddMessage(t *testing.T) {
	ctx := context.Background()
	runCtx := NewRunContext(ctx, "", 5)
	
	msg := types.Message{Role: "user", Content: "Hello"}
	runCtx.AddMessage(msg)
	
	messages := runCtx.Messages()
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
	
	if messages[0].Content != "Hello" {
		t.Errorf("expected content 'Hello', got %s", messages[0].Content)
	}
}

func TestRunContext_Messages_WithSystemPrompt(t *testing.T) {
	ctx := context.Background()
	runCtx := NewRunContext(ctx, "System prompt", 5)
	
	msg := types.Message{Role: "user", Content: "Hello"}
	runCtx.AddMessage(msg)
	
	messages := runCtx.Messages()
	if len(messages) != 2 {
		t.Errorf("expected 2 messages (system + user), got %d", len(messages))
	}
	
	if messages[0].Role != "system" {
		t.Errorf("expected first message role 'system', got %s", messages[0].Role)
	}
	
	if messages[0].Content != "System prompt" {
		t.Errorf("expected system prompt content, got %s", messages[0].Content)
	}
}

func TestRunContext_CanIterate(t *testing.T) {
	ctx := context.Background()
	runCtx := NewRunContext(ctx, "", 2)
	
	if !runCtx.CanIterate() {
		t.Error("should be able to iterate initially")
	}
	
	runCtx.IncrementIteration()
	if !runCtx.CanIterate() {
		t.Error("should still be able to iterate after one increment")
	}
	
	runCtx.IncrementIteration()
	if runCtx.CanIterate() {
		t.Error("should not be able to iterate after reaching max")
	}
}

func TestRunContext_CurrentIteration(t *testing.T) {
	ctx := context.Background()
	runCtx := NewRunContext(ctx, "", 5)
	
	if runCtx.CurrentIteration() != 0 {
		t.Errorf("expected initial iteration 0, got %d", runCtx.CurrentIteration())
	}
	
	runCtx.IncrementIteration()
	if runCtx.CurrentIteration() != 1 {
		t.Errorf("expected iteration 1, got %d", runCtx.CurrentIteration())
	}
}

func TestAgent_New(t *testing.T) {
	ag := New()
	
	if ag == nil {
		t.Error("agent should not be nil")
	}
	
	if ag.toolRegistry == nil {
		t.Error("tool registry should not be nil")
	}
	
	if ag.toolExecutor == nil {
		t.Error("tool executor should not be nil")
	}
}

func TestAgent_WithSystemPrompt(t *testing.T) {
	ag := New(WithSystemPrompt("Custom prompt"))
	
	if ag.systemPrompt != "Custom prompt" {
		t.Errorf("expected system prompt 'Custom prompt', got %s", ag.systemPrompt)
	}
}

func TestAgent_WithMaxIterations(t *testing.T) {
	ag := New(WithMaxIterations(10))
	
	if ag.maxIterations != 10 {
		t.Errorf("expected max iterations 10, got %d", ag.maxIterations)
	}
}
