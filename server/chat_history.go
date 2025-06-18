package main

import (
	"sync"
)
type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "assistant" or "system"
	Content string `json:"content"` // Message text
}

// Thread-safe storage of chat history
type ChatMemory struct {
	mu     sync.RWMutex
	history map[string][]ChatMessage // userId -> history
}

var memory = ChatMemory{
	history: make(map[string][]ChatMessage),
}

// Add a message to a user thread
func (c *ChatMemory) Add(userID string, msg ChatMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.history[userID] = append(c.history[userID], msg)
}

// Get a user’s conversation history
func (c *ChatMemory) Get(userID string) []ChatMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.history[userID]
}