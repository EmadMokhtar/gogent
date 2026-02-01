package memory

import (
	"context"
	"testing"
	"time"
)

func TestConversation(t *testing.T) {
	ctx := context.Background()
	
	t.Run("add and get messages", func(t *testing.T) {
		mem := NewConversation()
		
		msg1 := NewUserMessage("Hello")
		msg2 := NewAssistantMessage("Hi there!")
		
		if err := mem.Add(ctx, msg1); err != nil {
			t.Fatalf("Add() error = %v", err)
		}
		
		if err := mem.Add(ctx, msg2); err != nil {
			t.Fatalf("Add() error = %v", err)
		}
		
		messages, err := mem.Get(ctx, 10)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		
		if len(messages) != 2 {
			t.Errorf("Get() returned %d messages, want 2", len(messages))
		}
	})
	
	t.Run("search messages", func(t *testing.T) {
		mem := NewConversation()
		
		msg1 := NewUserMessage("What is Go?")
		msg2 := NewAssistantMessage("Go is a programming language.")
		msg3 := NewUserMessage("Tell me about Python")
		
		_ = mem.Add(ctx, msg1)
		_ = mem.Add(ctx, msg2)
		_ = mem.Add(ctx, msg3)
		
		results, err := mem.Search(ctx, "Go")
		if err != nil {
			t.Fatalf("Search() error = %v", err)
		}
		
		if len(results) < 1 {
			t.Errorf("Search() returned %d results, want at least 1", len(results))
		}
	})
	
	t.Run("clear messages", func(t *testing.T) {
		mem := NewConversation()
		
		_ = mem.Add(ctx, NewUserMessage("Test"))
		
		if err := mem.Clear(ctx); err != nil {
			t.Fatalf("Clear() error = %v", err)
		}
		
		count, _ := mem.Count(ctx)
		if count != 0 {
			t.Errorf("Count() = %d after Clear(), want 0", count)
		}
	})
	
	t.Run("max size limit", func(t *testing.T) {
		mem := NewConversation(WithMaxSize(2))
		
		_ = mem.Add(ctx, NewUserMessage("Message 1"))
		_ = mem.Add(ctx, NewUserMessage("Message 2"))
		_ = mem.Add(ctx, NewUserMessage("Message 3"))
		
		count, _ := mem.Count(ctx)
		if count > 2 {
			t.Errorf("Count() = %d, want <= 2", count)
		}
	})
}

func TestShortTerm(t *testing.T) {
	ctx := context.Background()
	
	t.Run("add and get messages", func(t *testing.T) {
		mem := NewShortTerm(WithTTL(1 * time.Hour))
		
		msg := NewUserMessage("Test message")
		if err := mem.Add(ctx, msg); err != nil {
			t.Fatalf("Add() error = %v", err)
		}
		
		messages, err := mem.Get(ctx, 10)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		
		if len(messages) != 1 {
			t.Errorf("Get() returned %d messages, want 1", len(messages))
		}
	})
	
	t.Run("message expiry", func(t *testing.T) {
		mem := NewShortTerm(WithTTL(1 * time.Millisecond))
		
		msg := NewUserMessage("Expiring message")
		_ = mem.Add(ctx, msg)
		
		// Wait for message to expire
		time.Sleep(10 * time.Millisecond)
		
		count, _ := mem.Count(ctx)
		if count != 0 {
			t.Errorf("Count() = %d after expiry, want 0", count)
		}
	})
}

func TestLongTerm(t *testing.T) {
	ctx := context.Background()
	
	t.Run("add with importance", func(t *testing.T) {
		mem := NewLongTerm()
		
		msg := NewUserMessage("Important message")
		if err := mem.AddWithImportance(ctx, msg, 0.9); err != nil {
			t.Fatalf("AddWithImportance() error = %v", err)
		}
		
		messages, err := mem.GetByImportance(ctx, 0.8, 10)
		if err != nil {
			t.Fatalf("GetByImportance() error = %v", err)
		}
		
		if len(messages) != 1 {
			t.Errorf("GetByImportance() returned %d messages, want 1", len(messages))
		}
	})
	
	t.Run("importance threshold", func(t *testing.T) {
		mem := NewLongTerm()
		
		_ = mem.AddWithImportance(ctx, NewUserMessage("High importance"), 0.9)
		_ = mem.AddWithImportance(ctx, NewUserMessage("Low importance"), 0.3)
		
		messages, _ := mem.GetByImportance(ctx, 0.7, 10)
		
		if len(messages) != 1 {
			t.Errorf("GetByImportance() returned %d messages, want 1", len(messages))
		}
	})
}
