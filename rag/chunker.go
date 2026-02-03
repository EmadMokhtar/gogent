package rag

import (
	"context"
	"regexp"
	"strings"
	"unicode"
)

// Chunker defines the interface for splitting text into chunks
type Chunker interface {
	// Chunk splits text into smaller chunks
	Chunk(ctx context.Context, text string) ([]string, error)
}

// FixedSizeChunker splits text into fixed-size chunks with optional overlap
type FixedSizeChunker struct {
	chunkSize    int
	chunkOverlap int
}

// FixedSizeChunkerOption defines configuration options for FixedSizeChunker
type FixedSizeChunkerOption func(*FixedSizeChunker)

// WithChunkSize sets the chunk size in characters
func WithChunkSize(size int) FixedSizeChunkerOption {
	return func(c *FixedSizeChunker) {
		c.chunkSize = size
	}
}

// WithOverlap sets the overlap size in characters
func WithOverlap(overlap int) FixedSizeChunkerOption {
	return func(c *FixedSizeChunker) {
		c.chunkOverlap = overlap
	}
}

// NewFixedSizeChunker creates a new fixed-size chunker
func NewFixedSizeChunker(opts ...FixedSizeChunkerOption) *FixedSizeChunker {
	chunker := &FixedSizeChunker{
		chunkSize:    500,
		chunkOverlap: 50,
	}

	for _, opt := range opts {
		opt(chunker)
	}

	return chunker
}

// Chunk splits text into fixed-size chunks with overlap
func (c *FixedSizeChunker) Chunk(ctx context.Context, text string) ([]string, error) {
	if text == "" {
		return nil, nil
	}

	var chunks []string
	runes := []rune(text)
	length := len(runes)

	for start := 0; start < length; {
		end := start + c.chunkSize
		if end > length {
			end = length
		}

		chunk := string(runes[start:end])
		chunks = append(chunks, strings.TrimSpace(chunk))

		// Move forward by chunkSize - overlap
		start += c.chunkSize - c.chunkOverlap
		if start >= length {
			break
		}
	}

	return chunks, nil
}

// SentenceChunker splits text into chunks at sentence boundaries
type SentenceChunker struct {
	maxChunkSize int
}

// SentenceChunkerOption defines configuration options for SentenceChunker
type SentenceChunkerOption func(*SentenceChunker)

// WithMaxChunkSize sets the maximum chunk size in characters
func WithMaxChunkSize(size int) SentenceChunkerOption {
	return func(c *SentenceChunker) {
		c.maxChunkSize = size
	}
}

// NewSentenceChunker creates a new sentence-based chunker
func NewSentenceChunker(opts ...SentenceChunkerOption) *SentenceChunker {
	chunker := &SentenceChunker{
		maxChunkSize: 1000,
	}

	for _, opt := range opts {
		opt(chunker)
	}

	return chunker
}

// Chunk splits text into chunks at sentence boundaries
func (c *SentenceChunker) Chunk(ctx context.Context, text string) ([]string, error) {
	if text == "" {
		return nil, nil
	}

	sentences := c.splitIntoSentences(text)
	var chunks []string
	var currentChunk strings.Builder

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// If adding this sentence would exceed max size, start a new chunk
		if currentChunk.Len() > 0 && currentChunk.Len()+len(sentence)+1 > c.maxChunkSize {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
		}
		currentChunk.WriteString(sentence)
	}

	// Add the last chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks, nil
}

// splitIntoSentences splits text into sentences
func (c *SentenceChunker) splitIntoSentences(text string) []string {
	// Simple sentence splitting using regex
	// Matches: . ! ? followed by whitespace and capital letter, or end of string
	sentenceRegex := regexp.MustCompile(`[.!?]+(?:\s+|$)`)
	indices := sentenceRegex.FindAllStringIndex(text, -1)

	if len(indices) == 0 {
		return []string{text}
	}

	var sentences []string
	lastEnd := 0

	for _, match := range indices {
		end := match[1]
		sentence := text[lastEnd:end]
		sentences = append(sentences, sentence)
		lastEnd = end
	}

	// Add remaining text if any
	if lastEnd < len(text) {
		sentences = append(sentences, text[lastEnd:])
	}

	return sentences
}

// ParagraphChunker splits text into chunks at paragraph boundaries
type ParagraphChunker struct {
	maxChunkSize int
}

// ParagraphChunkerOption defines configuration options for ParagraphChunker
type ParagraphChunkerOption func(*ParagraphChunker)

// WithMaxParagraphChunkSize sets the maximum chunk size in characters
func WithMaxParagraphChunkSize(size int) ParagraphChunkerOption {
	return func(c *ParagraphChunker) {
		c.maxChunkSize = size
	}
}

// NewParagraphChunker creates a new paragraph-based chunker
func NewParagraphChunker(opts ...ParagraphChunkerOption) *ParagraphChunker {
	chunker := &ParagraphChunker{
		maxChunkSize: 2000,
	}

	for _, opt := range opts {
		opt(chunker)
	}

	return chunker
}

// Chunk splits text into chunks at paragraph boundaries
func (c *ParagraphChunker) Chunk(ctx context.Context, text string) ([]string, error) {
	if text == "" {
		return nil, nil
	}

	paragraphs := c.splitIntoParagraphs(text)
	var chunks []string
	var currentChunk strings.Builder

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		// If adding this paragraph would exceed max size, start a new chunk
		if currentChunk.Len() > 0 && currentChunk.Len()+len(para)+2 > c.maxChunkSize {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(para)
	}

	// Add the last chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks, nil
}

// splitIntoParagraphs splits text into paragraphs
func (c *ParagraphChunker) splitIntoParagraphs(text string) []string {
	// Split by double newline or multiple whitespace lines
	paragraphRegex := regexp.MustCompile(`\n\s*\n`)
	paragraphs := paragraphRegex.Split(text, -1)

	// Filter out empty paragraphs
	var result []string
	for _, para := range paragraphs {
		para = strings.TrimFunc(para, unicode.IsSpace)
		if para != "" {
			result = append(result, para)
		}
	}

	return result
}
