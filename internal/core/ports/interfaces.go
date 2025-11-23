package ports

import (
	"context"

	"github.com/PundarikakshNTripathi/Kortex/internal/core/domain"
)

// AIProvider defines the interface for interacting with AI models.
type AIProvider interface {
	GenerateStream(ctx context.Context, history []domain.Message, tools []interface{}) (<-chan string, error)
}

// VectorStore defines the interface for storing and searching memory fragments.
type VectorStore interface {
	Save(ctx context.Context, fragment *domain.MemoryFragment) error
	Search(ctx context.Context, queryVector []float32, limit int) ([]domain.MemoryFragment, error)
}

// Browser defines the interface for browser interactions.
type Browser interface {
	Navigate(url string) error
	GetSnapshot() (string, error)
	Highlight(selector, message string) error
}
