package ports

import (
	"context"

	"github.com/PundarikakshNTripathi/Kortex/internal/core/domain"
)

// AIProvider defines the contract for any AI service (like Gemini, GPT-4, etc.).
// By using an interface, we can easily switch AI models without changing the core logic.
type AIProvider interface {
	// GenerateStream sends the chat history to the AI and gets a stream of text back.
	// It also provides a list of tools the AI can use.
	GenerateStream(ctx context.Context, history []domain.Message, tools []interface{}) (<-chan string, error)
}

// VectorStore defines how Kortex remembers things.
// It's an interface for a database that can search by "meaning" (vectors) rather than just keywords.
type VectorStore interface {
	// Save stores a new memory fragment.
	Save(ctx context.Context, fragment *domain.MemoryFragment) error

	// Search finds the most relevant memories for a given query vector.
	// 'limit' controls how many results to return.
	Search(ctx context.Context, queryVector []float32, limit int) ([]domain.MemoryFragment, error)
}

// Browser defines the "Hands" and "Eyes" of Kortex.
// It abstracts the web browser automation so the core logic doesn't need to know
// if we're using Playwright, Selenium, or something else.
type Browser interface {
	// Navigate goes to a specific URL.
	Navigate(url string) error

	// GetSnapshot returns a simplified text representation of the current page.
	// This is what the AI "sees" - a tree of elements, roles, and names.
	GetSnapshot() (string, error)

	// Highlight draws a visual box around an element to show the user what Kortex is looking at.
	Highlight(selector, message string) error

	// Click simulates a mouse click on an element.
	Click(selector string) error

	// Type simulates typing text into an input field.
	Type(selector, text string) error
}
