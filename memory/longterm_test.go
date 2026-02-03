package memory

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestLongTermMemory_Add(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(10))

	msg := Message{
		Role:       "user",
		Content:    "Important information",
		Importance: 0.8,
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

func TestLongTermMemory_ImportanceThreshold(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(
		WithCapacity(10),
		WithImportanceThreshold(0.5),
	)

	// Below threshold - should be ignored
	mem.Add(ctx, Message{
		Role:       "user",
		Content:    "Not important",
		Importance: 0.3,
	})

	// Above threshold - should be added
	mem.Add(ctx, Message{
		Role:       "user",
		Content:    "Important",
		Importance: 0.7,
	})

	messages, _ := mem.Get(ctx, 0)
	if len(messages) != 1 {
		t.Errorf("Expected 1 message (above threshold), got %d", len(messages))
	}

	if messages[0].Content != "Important" {
		t.Errorf("Expected 'Important', got '%s'", messages[0].Content)
	}
}

func TestLongTermMemory_CapacityPruning(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(3))

	// Add 5 messages with different importance
	messages := []Message{
		{Role: "user", Content: "A", Importance: 0.5},
		{Role: "user", Content: "B", Importance: 0.9},
		{Role: "user", Content: "C", Importance: 0.3},
		{Role: "user", Content: "D", Importance: 0.7},
		{Role: "user", Content: "E", Importance: 0.6},
	}

	for _, msg := range messages {
		mem.Add(ctx, msg)
	}

	result, _ := mem.Get(ctx, 0)
	if len(result) != 3 {
		t.Errorf("Expected 3 messages (capacity), got %d", len(result))
	}

	// Check that we kept the most important ones
	// Should keep B (0.9), D (0.7), E (0.6)
	hasB, hasD, hasE := false, false, false
	for _, msg := range result {
		switch msg.Content {
		case "B":
			hasB = true
		case "D":
			hasD = true
		case "E":
			hasE = true
		}
	}

	if !hasB || !hasD || !hasE {
		t.Error("Should keep the 3 most important messages (B, D, E)")
	}
}

func TestLongTermMemory_Search(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(10))

	messages := []Message{
		{Role: "user", Content: "User prefers dark mode", Importance: 0.9},
		{Role: "user", Content: "User likes coffee", Importance: 0.5},
		{Role: "user", Content: "Dark chocolate is good", Importance: 0.3},
	}

	for _, msg := range messages {
		mem.Add(ctx, msg)
	}

	// Search for "dark"
	results, _ := mem.Search(ctx, "dark", 0)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'dark', got %d", len(results))
	}

	// Results should be sorted by importance (descending)
	if results[0].Importance < results[1].Importance {
		t.Error("Search results should be sorted by importance (descending)")
	}

	// Search with limit
	results, _ = mem.Search(ctx, "dark", 1)
	if len(results) != 1 {
		t.Errorf("Expected 1 result with limit, got %d", len(results))
	}

	// Should return the most important match
	if results[0].Importance != 0.9 {
		t.Errorf("Expected most important match (0.9), got %f", results[0].Importance)
	}
}

func TestLongTermMemory_AccessCount(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(5))

	// Add messages with same importance but we'll access one more
	mem.Add(ctx, Message{Role: "user", Content: "A", Importance: 0.5})
	mem.Add(ctx, Message{Role: "user", Content: "B", Importance: 0.5})

	// Access messages multiple times
	mem.Get(ctx, 0)
	mem.Get(ctx, 0)

	// Check that access count increased
	messages, _ := mem.Get(ctx, 0)
	for _, msg := range messages {
		if msg.AccessCount < 1 {
			t.Error("Access count should increase on Get()")
		}
	}
}

func TestLongTermMemory_Clear(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm()

	// Add messages
	for i := 0; i < 5; i++ {
		mem.Add(ctx, Message{
			Role:       "user",
			Content:    "test",
			Importance: 0.5,
		})
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

func TestLongTermMemory_Stats(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(10))

	mem.Add(ctx, Message{Role: "user", Content: "Hello", Importance: 0.7})
	mem.Add(ctx, Message{Role: "assistant", Content: "Hi", Importance: 0.6})
	mem.Add(ctx, Message{Role: "system", Content: "System", Importance: 0.5})

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

	if stats.OldestMessage == nil || stats.NewestMessage == nil {
		t.Error("Oldest and Newest message should not be nil")
	}
}

func TestLongTermMemory_Concurrent(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(200))

	var wg sync.WaitGroup
	iterations := 50

	// Multiple goroutines adding messages
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mem.Add(ctx, Message{
					Role:       "user",
					Content:    "test",
					Importance: 0.5,
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
	// Should have 200 messages (capacity limit)
	if stats.TotalMessages != 200 {
		t.Errorf("Expected 200 messages (capacity), got %d", stats.TotalMessages)
	}
}

func TestLongTermMemory_GetOrdering(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(10))

	// Add messages with delays to ensure different timestamps
	for i := 0; i < 5; i++ {
		mem.Add(ctx, Message{
			Role:       "user",
			Content:    string(rune('A' + i)),
			Importance: 0.5,
		})
		time.Sleep(10 * time.Millisecond)
	}

	messages, _ := mem.Get(ctx, 0)

	// Should be in chronological order
	for i := 0; i < len(messages)-1; i++ {
		if messages[i].Timestamp.After(messages[i+1].Timestamp) {
			t.Error("Messages from Get() should be in chronological order")
		}
	}
}

func TestLongTermMemory_PruningWithAccessCount(t *testing.T) {
	ctx := context.Background()
	mem := NewLongTerm(WithCapacity(2))

	// Add two messages with same importance
	mem.Add(ctx, Message{Role: "user", Content: "A", Importance: 0.5})
	mem.Add(ctx, Message{Role: "user", Content: "B", Importance: 0.5})

	// Access message A multiple times
	mem.Search(ctx, "A", 1)
	mem.Search(ctx, "A", 1)

	// Add a third message with same importance
	mem.Add(ctx, Message{Role: "user", Content: "C", Importance: 0.5})

	// Should keep A and B (A has higher access count, B was there first)
	// or keep most accessed ones
	messages, _ := mem.Get(ctx, 0)
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages after pruning, got %d", len(messages))
	}
}
