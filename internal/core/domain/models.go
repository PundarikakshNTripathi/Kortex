package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a continuous conversation between the user and Kortex.
// Think of this like a chat history in a messaging app.
// It groups related messages together so the agent can remember the context.
type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"` // Unique identifier for this conversation
	Context   string    `json:"context"`                          // Optional summary or topic of the session
	CreatedAt time.Time `json:"created_at"`                       // When this conversation started
}

// BeforeCreate is a GORM hook that runs before a new Session is saved to the database.
// It ensures that every session has a unique ID.
func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

// Message represents a single piece of communication in a session.
// It could be from the User ("Find flights"), the Model ("I'm looking..."),
// or a Tool ("Navigated to google.com").
type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`  // Unique ID for this specific message
	SessionID  uuid.UUID `gorm:"type:uuid;index" json:"session_id"` // Links this message to a specific conversation
	Role       string    `json:"role"`                              // Who sent this? "user", "model", or "tool"
	Content    string    `json:"content"`                           // The actual text content
	ToolCallID *string   `json:"tool_call_id,omitempty"`            // If this was a tool use, what was the ID?
	ToolResult *string   `json:"tool_result,omitempty"`             // If this was a tool output, what was the result?
	CreatedAt  time.Time `json:"created_at"`                        // When this message was sent
}

// BeforeCreate ensures every message has a unique ID before saving.
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

// MemoryFragment represents a piece of long-term memory for Kortex.
// Unlike short-term session history, these are stored permanently and retrieved
// based on similarity to the current topic.
//
// Example: If you ask about "Tokyo" today, Kortex can search these fragments
// to remember you asked about "Japan" last week.
type MemoryFragment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"` // Unique ID for this memory
	Content   string    `json:"content"`                          // The text content to remember
	Embedding []float32 `gorm:"serializer:json" json:"embedding"` // A vector (list of numbers) representing the meaning of the text
	Tags      string    `gorm:"type:json" json:"tags"`            // Extra labels for filtering (JSON string)
	CreatedAt time.Time `json:"created_at"`                       // When this memory was formed
}

// BeforeCreate ensures every memory fragment has a unique ID.
func (mf *MemoryFragment) BeforeCreate(tx *gorm.DB) (err error) {
	if mf.ID == uuid.Nil {
		mf.ID = uuid.New()
	}
	return
}
