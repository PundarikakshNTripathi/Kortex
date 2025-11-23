package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a user session.
type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Context   string    `json:"context"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate generates a new UUID for the session if not present.
func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

// Message represents a message in a session.
type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	SessionID  uuid.UUID `gorm:"type:uuid;index" json:"session_id"`
	Role       string    `json:"role"` // user, model, tool
	Content    string    `json:"content"`
	ToolCallID *string   `json:"tool_call_id,omitempty"`
	ToolResult *string   `json:"tool_result,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// BeforeCreate generates a new UUID for the message if not present.
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

// MemoryFragment represents a piece of memory with embedding.
type MemoryFragment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Content   string    `json:"content"`
	Embedding []float32 `gorm:"serializer:json" json:"embedding"` // 768 dimensions
	Tags      string    `gorm:"type:json" json:"tags"`            // JSON string
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate generates a new UUID for the memory fragment if not present.
func (mf *MemoryFragment) BeforeCreate(tx *gorm.DB) (err error) {
	if mf.ID == uuid.Nil {
		mf.ID = uuid.New()
	}
	return
}
