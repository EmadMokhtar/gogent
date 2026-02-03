package openai

import (
	"context"
	"testing"
)

func TestEmbedder_Dimension(t *testing.T) {
	embedder := NewEmbedder(
		WithAPIKey("test-key"),
		WithModel("text-embedding-ada-002"),
	)

	if embedder.Dimension() != DefaultDimension {
		t.Errorf("Dimension() = %d, expected %d", embedder.Dimension(), DefaultDimension)
	}
}

func TestEmbedder_EmptyText(t *testing.T) {
	embedder := NewEmbedder(WithAPIKey("test-key"))
	ctx := context.Background()

	_, err := embedder.Embed(ctx, "")
	if err != ErrEmptyText {
		t.Errorf("Embed() with empty text error = %v, expected %v", err, ErrEmptyText)
	}
}

func TestEmbedder_NoAPIKey(t *testing.T) {
	embedder := NewEmbedder() // No API key
	ctx := context.Background()

	_, err := embedder.Embed(ctx, "test")
	if err != ErrEmptyAPIKey {
		t.Errorf("Embed() without API key error = %v, expected %v", err, ErrEmptyAPIKey)
	}
}

func TestEmbedder_CacheEnabled(t *testing.T) {
	embedder := NewEmbedder(
		WithAPIKey("test-key"),
		WithCache(true),
	)

	if !embedder.useCache {
		t.Error("Expected cache to be enabled")
	}

	if embedder.cache == nil {
		t.Error("Expected cache map to be initialized")
	}
}

func TestEmbedder_CacheDisabled(t *testing.T) {
	embedder := NewEmbedder(
		WithAPIKey("test-key"),
		WithCache(false),
	)

	if embedder.useCache {
		t.Error("Expected cache to be disabled")
	}
}

func TestEmbedder_CustomBaseURL(t *testing.T) {
	customURL := "https://custom.api.com/v1"
	embedder := NewEmbedder(
		WithAPIKey("test-key"),
		WithBaseURL(customURL),
	)

	if embedder.baseURL != customURL {
		t.Errorf("baseURL = %s, expected %s", embedder.baseURL, customURL)
	}
}

func TestEmbedder_CustomModel(t *testing.T) {
	customModel := "text-embedding-3-small"
	embedder := NewEmbedder(
		WithAPIKey("test-key"),
		WithModel(customModel),
	)

	if embedder.model != customModel {
		t.Errorf("model = %s, expected %s", embedder.model, customModel)
	}
}

func TestLLM_NoAPIKey(t *testing.T) {
	llm := New() // No API key
	ctx := context.Background()

	_, err := llm.Generate(ctx, nil)
	if err == nil {
		t.Error("Generate() without API key should return error")
	}
}

func TestLLM_Model(t *testing.T) {
	customModel := "gpt-4"
	llm := New(
		WithLLMAPIKey("test-key"),
		WithLLMModel(customModel),
	)

	if llm.Model() != customModel {
		t.Errorf("Model() = %s, expected %s", llm.Model(), customModel)
	}
}
