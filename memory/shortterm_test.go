package memory

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestShortTermMemory_Add(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(WithTTL(5 * time.Second))
	defer mem.Stop()

	msg := Message{
		Role:    "user",
		Content: "Hello",
	}

	err := mem.Add(ctx, msg)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	messages, _ := mem.Get(ctx, 0)
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
}

func TestShortTermMemory_TTLExpiration(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(WithTTL(100 * time.Millisecond))
	defer mem.Stop()

	msg := Message{
		Role:    "user",
		Content: "Should expire",
	}

	mem.Add(ctx, msg)

	// Message should exist initially
	messages, _ := mem.Get(ctx, 0)
	if len(messages) != 1 {
		t.Errorf("Expected 1 message initially, got %d", len(messages))
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Message should be expired now
	messages, _ = mem.Get(ctx, 0)
	if len(messages) != 0 {
		t.Errorf("Expected 0 messages after expiration, got %d", len(messages))
	}
}

func TestShortTermMemory_CustomTTL(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(WithTTL(5 * time.Second)) // default TTL
	defer mem.Stop()

	// Message with custom TTL
	msg := Message{
		Role:    "user",
		Content: "Custom TTL",
		TTL:     100 * time.Millisecond,
	}

	mem.Add(ctx, msg)

	// Should exist initially
	messages, _ := mem.Get(ctx, 0)
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	// Wait for custom TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	messages, _ = mem.Get(ctx, 0)
	if len(messages) != 0 {
		t.Errorf("Expected 0 messages after custom TTL, got %d", len(messages))
	}
}

func TestShortTermMemory_BackgroundCleanup(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(
		WithTTL(100*time.Millisecond),
		WithCleanupInterval(200*time.Millisecond),
	)
	defer mem.Stop()

	// Add multiple messages
	for i := 0; i < 5; i++ {
		mem.Add(ctx, Message{
			Role:    "user",
			Content: "test",
		})
	}

	// Check initial count
	stats1, _ := mem.Stats(ctx)
	if stats1.TotalMessages != 5 {
		t.Errorf("Expected 5 messages initially, got %d", stats1.TotalMessages)
	}

	// Wait for TTL expiration
	time.Sleep(150 * time.Millisecond)

	// Messages are expired but may not be cleaned up yet
	stats2, _ := mem.Stats(ctx)
	if stats2.TotalMessages != 0 {
		t.Errorf("Expected 0 non-expired messages, got %d", stats2.TotalMessages)
	}

	// Wait for cleanup to run
	time.Sleep(100 * time.Millisecond)

	// Now the expired messages should be removed from storage
	mem.mu.RLock()
	actualCount := len(mem.messages)
	mem.mu.RUnlock()

	if actualCount != 0 {
		t.Errorf("Expected 0 messages in storage after cleanup, got %d", actualCount)
	}
}

func TestShortTermMemory_Search(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(WithTTL(1 * time.Minute))
	defer mem.Stop()

	messages := []Message{
		{Role: "user", Content: "I like Go"},
		{Role: "assistant", Content: "Go is great"},
		{Role: "user", Content: "Python is cool"},
	}

	for _, msg := range messages {
		mem.Add(ctx, msg)
	}

	// Search for "Go"
	results, _ := mem.Search(ctx, "Go", 0)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'Go', got %d", len(results))
	}

	// Search with limit
	results, _ = mem.Search(ctx, "Go", 1)
	if len(results) != 1 {
		t.Errorf("Expected 1 result with limit, got %d", len(results))
	}
}

func TestShortTermMemory_Clear(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm()
	defer mem.Stop()

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

func TestShortTermMemory_Stats(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(WithTTL(1 * time.Minute))
	defer mem.Stop()

	mem.Add(ctx, Message{Role: "user", Content: "Hello"})
	mem.Add(ctx, Message{Role: "assistant", Content: "Hi"})
	mem.Add(ctx, Message{Role: "system", Content: "System"})

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
}

func TestShortTermMemory_Concurrent(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(WithTTL(1 * time.Minute))
	defer mem.Stop()

	var wg sync.WaitGroup
	iterations := 50

	// Multiple goroutines adding messages
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mem.Add(ctx, Message{
					Role:    "user",
					Content: "test",
				})
			}
		}()
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

func TestShortTermMemory_GetOrdering(t *testing.T) {
	ctx := context.Background()
	mem := NewShortTerm(WithTTL(1 * time.Minute))
	defer mem.Stop()

	// Add messages with delays to ensure different timestamps
	for i := 0; i < 5; i++ {
		mem.Add(ctx, Message{
			Role:    "user",
			Content: string(rune('A' + i)),
		})
		time.Sleep(10 * time.Millisecond)
	}

	messages, _ := mem.Get(ctx, 0)
	
	// Should be in chronological order
	for i := 0; i < len(messages)-1; i++ {
		if messages[i].Timestamp.After(messages[i+1].Timestamp) {
			t.Error("Messages should be in chronological order")
		}
	}
}
