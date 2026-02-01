package agent

import (
	"context"
	"sync"
	"time"
)

// Context manages the execution context for an agent
type Context struct {
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	sessionID string
	metadata  map[string]interface{}
	startTime time.Time
}

// NewContext creates a new agent context
func NewContext(ctx context.Context, sessionID string) *Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	
	return &Context{
		ctx:       ctxWithCancel,
		cancel:    cancel,
		sessionID: sessionID,
		metadata:  make(map[string]interface{}),
		startTime: time.Now(),
	}
}

// Context returns the underlying context
func (c *Context) Context() context.Context {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ctx
}

// Cancel cancels the context
func (c *Context) Cancel() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancel != nil {
		c.cancel()
	}
}

// SessionID returns the session ID
func (c *Context) SessionID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessionID
}

// SetMetadata sets a metadata value
func (c *Context) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metadata[key] = value
}

// GetMetadata gets a metadata value
func (c *Context) GetMetadata(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.metadata[key]
	return val, ok
}

// GetAllMetadata returns all metadata
func (c *Context) GetAllMetadata() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[string]interface{}, len(c.metadata))
	for k, v := range c.metadata {
		result[k] = v
	}
	return result
}

// Elapsed returns the time elapsed since context creation
func (c *Context) Elapsed() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Since(c.startTime)
}
