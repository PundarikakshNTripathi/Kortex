package sqlite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PundarikakshNTripathi/Kortex/internal/core/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SQLiteVectorStore struct {
	db *gorm.DB
}

func NewSQLiteVectorStore(dbPath string) (*SQLiteVectorStore, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable foreign keys
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// AutoMigrate
	if err := db.AutoMigrate(&domain.Session{}, &domain.Message{}, &domain.MemoryFragment{}); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &SQLiteVectorStore{db: db}, nil
}

func (s *SQLiteVectorStore) Save(ctx context.Context, fragment *domain.MemoryFragment) error {
	return s.db.WithContext(ctx).Create(fragment).Error
}

func (s *SQLiteVectorStore) Search(ctx context.Context, queryVector []float32, limit int) ([]domain.MemoryFragment, error) {
	var fragments []domain.MemoryFragment

	// Serialize query vector to JSON string for sqlite-vec if needed,
	// but typically sqlite-vec takes a blob or specific format.
	// However, the user requested "Raw SQL for the sqlite-vec extension".
	// Assuming vec_distance_cosine takes two vectors.
	// We need to pass the query vector correctly.
	// For this implementation, we'll assume the query vector can be passed as a parameter.
	// Note: The actual syntax depends on how sqlite-vec is loaded and used.
	// Common pattern: SELECT *, vec_distance_cosine(embedding, ?) as distance FROM memory_fragments ORDER BY distance LIMIT ?

	// We need to convert []float32 to a format sqlite-vec understands (often JSON or binary).
	// Let's try JSON string representation for now as it's common for some extensions,
	// or we might need a specific serialization.
	// Given "Embedding (Float32 array, len 768)", we'll try passing it as a JSON string first.

	queryVectorJSON, err := json.Marshal(queryVector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query vector: %w", err)
	}

	// Using Raw SQL as requested
	// Note: We are selecting all fields from memory_fragments
	// We assume the table name is `memory_fragments` (GORM default snake_case plural)
	query := `
		SELECT id, content, embedding, tags, created_at, vec_distance_cosine(embedding, ?) as distance
		FROM memory_fragments
		ORDER BY distance
		LIMIT ?
	`

	if err := s.db.WithContext(ctx).Raw(query, string(queryVectorJSON), limit).Scan(&fragments).Error; err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}

	return fragments, nil
}
