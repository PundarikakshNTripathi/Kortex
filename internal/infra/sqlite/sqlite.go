package sqlite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PundarikakshNTripathi/Kortex/internal/core/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// SQLiteVectorStore implements the VectorStore interface using a local SQLite database.
// It uses the 'sqlite-vec' extension (conceptually) to perform vector similarity searches.
type SQLiteVectorStore struct {
	db *gorm.DB // The GORM database connection
}

// NewSQLiteVectorStore opens or creates the database file.
func NewSQLiteVectorStore(dbPath string) (*SQLiteVectorStore, error) {
	// Open the SQLite database file
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable foreign keys for data integrity
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// AutoMigrate creates the tables if they don't exist.
	// It looks at our struct definitions (Session, Message, etc.) and makes matching tables.
	if err := db.AutoMigrate(&domain.Session{}, &domain.Message{}, &domain.MemoryFragment{}); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &SQLiteVectorStore{db: db}, nil
}

// Save stores a memory fragment in the database.
func (s *SQLiteVectorStore) Save(ctx context.Context, fragment *domain.MemoryFragment) error {
	return s.db.WithContext(ctx).Create(fragment).Error
}

// Search finds memories that are semantically similar to the query vector.
// It uses a cosine distance function: smaller distance = more similar.
func (s *SQLiteVectorStore) Search(ctx context.Context, queryVector []float32, limit int) ([]domain.MemoryFragment, error) {
	var fragments []domain.MemoryFragment

	// We convert the query vector to a JSON string to pass it to the SQL query.
	queryVectorJSON, err := json.Marshal(queryVector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query vector: %w", err)
	}

	// This raw SQL query uses the vector extension function `vec_distance_cosine`.
	// It calculates the distance between the stored embedding and our query.
	// We order by distance (ascending) to get the closest matches.
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
