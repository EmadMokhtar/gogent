package agent

import (
	"context"
	"testing"

	"github.com/EmadMokhtar/gogent/memory"
)

func TestAgent_Run_WithoutMemory(t *testing.T) {
	ctx := context.Background()
	ag := New()

	resp, err := ag.Run(ctx, "Hello")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if resp.Content != "Echo: Hello" {
		t.Errorf("Expected 'Echo: Hello', got '%s'", resp.Content)
	}
}

func TestAgent_Run_WithMemory(t *testing.T) {
	ctx := context.Background()
	mem := memory.NewConversation()

	ag := New(
		WithMemory(mem),
		WithSystemPrompt("You are a helpful assistant."),
	)

	// First interaction
	resp1, err := ag.Run(ctx, "My name is Alice")
	if err != nil {
		t.Fatalf("First Run failed: %v", err)
	}

	if resp1.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", resp1.Role)
	}

	// Check memory - should have system, user, and assistant messages
	stats, _ := mem.Stats(ctx)
	if stats.TotalMessages < 2 {
		t.Errorf("Expected at least 2 messages in memory (user + assistant), got %d", stats.TotalMessages)
	}

	// Second interaction
	resp2, err := ag.Run(ctx, "What's my name?")
	if err != nil {
		t.Fatalf("Second Run failed: %v", err)
	}

	// Memory should have more messages now
	stats, _ = mem.Stats(ctx)
	if stats.TotalMessages < 4 {
		t.Errorf("Expected at least 4 messages after second interaction, got %d", stats.TotalMessages)
	}

	if stats.UserMessages < 2 {
		t.Errorf("Expected at least 2 user messages, got %d", stats.UserMessages)
	}

	if stats.AssistantMessages < 2 {
		t.Errorf("Expected at least 2 assistant messages, got %d", stats.AssistantMessages)
	}

	_ = resp2 // Silence unused warning
}

func TestAgent_WithSystemPrompt(t *testing.T) {
	ctx := context.Background()
	mem := memory.NewConversation()

	ag := New(
		WithMemory(mem),
		WithSystemPrompt("You are a test assistant."),
	)

	ag.Run(ctx, "Hello")

	// System prompt is used for LLM context but not stored in memory
	// Only user and assistant messages are stored
	messages, _ := mem.Get(ctx, 0)
	
	// Should have user and assistant messages
	if len(messages) < 2 {
		t.Errorf("Expected at least 2 messages (user + assistant), got %d", len(messages))
	}

	// Verify no system messages are stored (they're only used for context)
	for _, msg := range messages {
		if msg.Role == "system" {
			t.Error("System prompts should not be stored in memory")
		}
	}
}

func TestAgent_MemoryContext(t *testing.T) {
	ctx := context.Background()
	mem := memory.NewConversation()

	ag := New(
		WithMemory(mem),
	)

	// Add some messages
	ag.Run(ctx, "First message")
	ag.Run(ctx, "Second message")
	ag.Run(ctx, "Third message")

	// Memory should have all messages
	messages, _ := mem.Get(ctx, 0)
	if len(messages) < 6 { // 3 user + 3 assistant
		t.Errorf("Expected at least 6 messages in memory, got %d", len(messages))
	}

	// Get recent messages only
	recentMessages, _ := mem.Get(ctx, 2)
	if len(recentMessages) != 2 {
		t.Errorf("Expected 2 recent messages, got %d", len(recentMessages))
	}
}
