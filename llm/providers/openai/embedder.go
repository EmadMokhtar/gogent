package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	// ErrEmptyAPIKey is returned when no API key is provided
	ErrEmptyAPIKey = errors.New("OpenAI API key is required")
	// ErrEmptyText is returned when trying to embed empty text
	ErrEmptyText = errors.New("text cannot be empty")
)

const (
	// DefaultModel is the default embedding model
	DefaultModel = "text-embedding-ada-002"
	// DefaultDimension is the dimension for ada-002
	DefaultDimension = 1536
	// MaxBatchSize is the maximum number of texts per API call
	MaxBatchSize = 100
)

// Embedder generates embeddings using OpenAI's API
type Embedder struct {
	apiKey     string
	model      string
	dimension  int
	baseURL    string
	httpClient *http.Client
	cache      map[string][]float32
	cacheMu    sync.RWMutex
	useCache   bool
}

// EmbedderOption defines configuration options for Embedder
type EmbedderOption func(*Embedder)

// WithAPIKey sets the OpenAI API key
func WithAPIKey(apiKey string) EmbedderOption {
	return func(e *Embedder) {
		e.apiKey = apiKey
	}
}

// WithModel sets the embedding model
func WithModel(model string) EmbedderOption {
	return func(e *Embedder) {
		e.model = model
	}
}

// WithCache enables or disables caching
func WithCache(useCache bool) EmbedderOption {
	return func(e *Embedder) {
		e.useCache = useCache
		if useCache && e.cache == nil {
			e.cache = make(map[string][]float32)
		}
	}
}

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) EmbedderOption {
	return func(e *Embedder) {
		e.baseURL = baseURL
	}
}

// NewEmbedder creates a new OpenAI embedder
func NewEmbedder(opts ...EmbedderOption) *Embedder {
	embedder := &Embedder{
		model:      DefaultModel,
		dimension:  DefaultDimension,
		baseURL:    "https://api.openai.com/v1",
		httpClient: &http.Client{Timeout: 30 * time.Second},
		useCache:   false,
	}

	for _, opt := range opts {
		opt(embedder)
	}

	return embedder
}

// Embed generates an embedding for a single text
func (e *Embedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	// Check cache
	if e.useCache {
		e.cacheMu.RLock()
		if cached, ok := e.cache[text]; ok {
			e.cacheMu.RUnlock()
			return cached, nil
		}
		e.cacheMu.RUnlock()
	}

	embeddings, err := e.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	if len(embeddings) == 0 {
		return nil, errors.New("no embeddings returned")
	}

	return embeddings[0], nil
}

// EmbedBatch generates embeddings for multiple texts
func (e *Embedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if e.apiKey == "" {
		return nil, ErrEmptyAPIKey
	}

	if len(texts) == 0 {
		return nil, nil
	}

	// Process in batches if necessary
	var allEmbeddings [][]float32
	for i := 0; i < len(texts); i += MaxBatchSize {
		end := i + MaxBatchSize
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]
		embeddings, err := e.embedBatch(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("batch %d-%d failed: %w", i, end, err)
		}

		allEmbeddings = append(allEmbeddings, embeddings...)
	}

	return allEmbeddings, nil
}

// embedBatch processes a single batch (internal)
func (e *Embedder) embedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	// Check cache for all texts
	if e.useCache {
		var uncached []string
		var uncachedIndices []int
		results := make([][]float32, len(texts))

		e.cacheMu.RLock()
		for i, text := range texts {
			if cached, ok := e.cache[text]; ok {
				results[i] = cached
			} else {
				uncached = append(uncached, text)
				uncachedIndices = append(uncachedIndices, i)
			}
		}
		e.cacheMu.RUnlock()

		if len(uncached) == 0 {
			return results, nil
		}

		// Fetch embeddings for uncached texts
		freshEmbeddings, err := e.callAPI(ctx, uncached)
		if err != nil {
			return nil, err
		}

		// Cache and merge results
		e.cacheMu.Lock()
		for i, embedding := range freshEmbeddings {
			idx := uncachedIndices[i]
			results[idx] = embedding
			e.cache[texts[idx]] = embedding
		}
		e.cacheMu.Unlock()

		return results, nil
	}

	return e.callAPI(ctx, texts)
}

// callAPI makes the actual API call to OpenAI
func (e *Embedder) callAPI(ctx context.Context, texts []string) ([][]float32, error) {
	reqBody := map[string]interface{}{
		"input": texts,
		"model": e.model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := e.baseURL + "/embeddings"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
			Index     int       `json:"index"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Data) != len(texts) {
		return nil, fmt.Errorf("expected %d embeddings, got %d", len(texts), len(response.Data))
	}

	embeddings := make([][]float32, len(texts))
	for _, item := range response.Data {
		embeddings[item.Index] = item.Embedding
	}

	return embeddings, nil
}

// Dimension returns the dimensionality of the embeddings
func (e *Embedder) Dimension() int {
	return e.dimension
}
