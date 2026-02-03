package memory

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestConversationMemory_Add(t *testing.T) {
	ctx := context.Background()
	mem := NewConversation()

	msg := Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	err := mem.Add(ctx, msg)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	messages, err := mem.Get(ctx, 0)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if messages[0].Content != "Hello, world!" {
		t.Errorf("Expected 'Hello, world!', got '%s'", messages[0].Content)
	}

	if messages[0].Timestamp.IsZero() {
		t.Error("Timestamp should be set automatically")
	}
}

func TestConversationMemory_MaxSize(t *testing.T) {
	ctx := context.Background()
	mem := NewConversation(WithMaxSize(3))

	// Add 5 messages
	for i := 0; i < 5; i++ {
		msg := Message{
			Role:    "user",
			Content: string(rune('A' + i)),
		}
		mem.Add(ctx, msg)
	}

	messages, _ := mem.Get(ctx, 0)
	if len(messages) != 3 {
		t.Errorf("Expected 3 messages (max size), got %d", len(messages))
	}

	// Should have the last 3 messages (C, D, E)
	if messages[0].Content != "C" {
		t.Errorf("Expected first message to be 'C', got '%s'", messages[0].Content)
	}
}

func TestConversationMemory_Get(t *testing.T) {
	ctx := context.Background()
	mem := NewConversation()

	// Add 10 messages
	for i := 0; i < 10; i++ {
		msg := Message{
			Role:    "user",
			Content: string(rune('A' + i)),
		}
		mem.Add(ctx, msg)
	}

	// Get all messages
	all, _ := mem.Get(ctx, 0)
	if len(all) != 10 {
		t.Errorf("Expected 10 messages, got %d", len(all))
	}

	// Get last 5 messages
	last5, _ := mem.Get(ctx, 5)
	if len(last5) != 5 {
		t.Errorf("Expected 5 messages, got %d", len(last5))
	}

	// Should be F, G, H, I, J
	if last5[0].Content != "F" {
		t.Errorf("Expected first of last 5 to be 'F', got '%s'", last5[0].Content)
	}
}

func TestConversationMemory_Search(t *testing.T) {
	ctx := context.Background()
	mem := NewConversation()

	messages := []Message{
		{Role: "user", Content: "I love programming"},
		{Role: "assistant", Content: "That's great!"},
		{Role: "user", Content: "Especially Go programming"},
		{Role: "assistant", Content: "Go is awesome"},
	}

	for _, msg := range messages {
		mem.Add(ctx, msg)
	}

	// Search for "programming"
	results, _ := mem.Search(ctx, "programming", 0)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'programming', got %d", len(results))
	}

	// Search with limit
	results, _ = mem.Search(ctx, "programming", 1)
	if len(results) != 1 {
		t.Errorf("Expected 1 result with limit, got %d", len(results))
	}

	// Case insensitive search
	results, _ = mem.Search(ctx, "PROGRAMMING", 0)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for case-insensitive search, got %d", len(results))
	}
}

func TestConversationMemory_Clear(t *testing.T) {
	ctx := context.Background()
	mem := NewConversation()

	// Add messages
	for i := 0; i < 5; i++ {
		mem.Add(ctx, Message{Role: "user", Content: "test"})
	}

	err := mem.Clear(ctx)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	messages, _ := mem.Get(ctx, 0)
	if len(messages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(messages))
	}
}

func TestConversationMemory_Stats(t *testing.T) {
	ctx := context.Background()
	mem := NewConversation()

	// Add different types of messages
	mem.Add(ctx, Message{Role: "user", Content: "Hello"})
	time.Sleep(10 * time.Millisecond)
	mem.Add(ctx, Message{Role: "assistant", Content: "Hi there"})
	time.Sleep(10 * time.Millisecond)
	mem.Add(ctx, Message{Role: "system", Content: "System message"})

	stats, err := mem.Stats(ctx)
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	if stats.TotalMessages != 3 {
		t.Errorf("Expected 3 total messages, got %d", stats.TotalMessages)
	}

	if stats.UserMessages != 1 {
		t.Errorf("Expected 1 user message, got %d", stats.UserMessages)
	}

	if stats.AssistantMessages != 1 {
		t.Errorf("Expected 1 assistant message, got %d", stats.AssistantMessages)
	}

	if stats.SystemMessages != 1 {
		t.Errorf("Expected 1 system message, got %d", stats.SystemMessages)
	}

	if stats.OldestMessage == nil {
		t.Error("OldestMessage should not be nil")
	}

	if stats.NewestMessage == nil {
		t.Error("NewestMessage should not be nil")
	}

	if stats.TotalSize <= 0 {
		t.Error("TotalSize should be greater than 0")
	}
}

func TestConversationMemory_Concurrent(t *testing.T) {
	ctx := context.Background()
	mem := NewConversation()

	var wg sync.WaitGroup
	iterations := 100

	// Multiple goroutines adding messages
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mem.Add(ctx, Message{
					Role:    "user",
					Content: "test",
				})
			}
		}(i)
	}

	// Multiple goroutines reading messages
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mem.Get(ctx, 10)
			}
		}()
	}

	wg.Wait()

	stats, _ := mem.Stats(ctx)
	if stats.TotalMessages != 10*iterations {
		t.Errorf("Expected %d messages, got %d", 10*iterations, stats.TotalMessages)
	}
}
