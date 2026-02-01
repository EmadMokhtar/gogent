package rag

import (
	"fmt"
	"strings"
)

// ChunkerConfig holds configuration for document chunking
type ChunkerConfig struct {
	ChunkSize    int
	ChunkOverlap int
	Separator    string
}

// Chunker splits documents into chunks
type Chunker struct {
	config ChunkerConfig
}

// NewChunker creates a new chunker with the given configuration
func NewChunker(config ChunkerConfig) *Chunker {
	if config.ChunkSize <= 0 {
		config.ChunkSize = 500
	}
	if config.ChunkOverlap < 0 {
		config.ChunkOverlap = 0
	}
	if config.ChunkOverlap >= config.ChunkSize {
		config.ChunkOverlap = config.ChunkSize / 4
	}
	if config.Separator == "" {
		config.Separator = " "
	}
	
	return &Chunker{
		config: config,
	}
}

// ChunkDocument splits a document into chunks
func (c *Chunker) ChunkDocument(doc Document) ([]Chunk, error) {
	if doc.Content == "" {
		return nil, fmt.Errorf("document content is empty")
	}
	
	chunks := make([]Chunk, 0)
	text := doc.Content
	
	// Simple character-based chunking with overlap
	start := 0
	index := 0
	
	for start < len(text) {
		end := start + c.config.ChunkSize
		if end > len(text) {
			end = len(text)
		}
		
		// Try to break at word boundary
		if end < len(text) {
			// Look for the last space within chunk
			for i := end - 1; i > start && i > end-50; i-- {
				if text[i] == ' ' || text[i] == '\n' {
					end = i + 1
					break
				}
			}
		}
		
		content := strings.TrimSpace(text[start:end])
		if content != "" {
			chunkID := fmt.Sprintf("%s_chunk_%d", doc.ID, index)
			chunk := NewChunk(chunkID, doc.ID, content, index)
			
			// Copy metadata from document
			for k, v := range doc.Metadata {
				chunk.Metadata[k] = v
			}
			chunk.Metadata["chunk_index"] = index
			
			chunks = append(chunks, chunk)
			index++
		}
		
		// Move start position with overlap
		newStart := end - c.config.ChunkOverlap
		// Ensure we always move forward
		if newStart <= start {
			start = end
		} else {
			start = newStart
		}
	}
	
	return chunks, nil
}

// ChunkDocuments splits multiple documents into chunks
func (c *Chunker) ChunkDocuments(docs []Document) ([]Chunk, error) {
	allChunks := make([]Chunk, 0)
	
	for _, doc := range docs {
		chunks, err := c.ChunkDocument(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to chunk document %s: %w", doc.ID, err)
		}
		allChunks = append(allChunks, chunks...)
	}
	
	return allChunks, nil
}
