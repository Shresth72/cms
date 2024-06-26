package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Shresth72/server/internals/utils"
)

type MetadataModel struct {
	db *sql.DB
}

type Metadata struct {
	ID          int    `json:"id"`
  Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	Url         string `json:"url"`
	Author      string `json:"author"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func VerifyMetadata(v *utils.Validator, metadata *Metadata) {
	v.Check(utils.MinChars(metadata.Title, 4), "title", "Title must have atleast 4 characters")
	v.Check(utils.MaxChars(metadata.Title, 200), "title", "Title must have atmost 200 characters")
	// More later
}

// repository functions
func (m *MetadataModel) CreateFile(metadata *Metadata) error {
  query := `INSERT INTO metadata (key, title, desc, author, url) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, created_at`

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  args := []any{metadata.Key, metadata.Title, metadata.Description, metadata.Author, metadata.Url}

  // Returning
  err := m.db.QueryRowContext(ctx, query, args...).Scan(
    &metadata.ID, 
    &metadata.Title, 
    &metadata.CreatedAt,
  )
  if err != nil {
    return err
  }
  return nil
}

func (m *MetadataModel) FindFiles() ([]Metadata, error) {
  query := `SELECT id, title, url, author, created_at FROM metadata`

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  rows, err := m.db.QueryContext(ctx, query)
  if err != nil {
    return nil, err
  }

  var metadata []Metadata = make([]Metadata, 0)
  for rows.Next() {
    var data Metadata
    err = rows.Scan(
      &data.ID,
      &data.Title,
      &data.Url,
      &data.Author,
      &data.CreatedAt,
    )

    if err != nil {
      return nil, err
    }
    metadata = append(metadata, data)
  }

  if rows.Err() != nil {
    return nil, err
  }
  return metadata, nil
}

func (m *MetadataModel) FindFileById(id int) (*Metadata, error) {
  var metadata Metadata

  query := `SELECT id, title, desc, author, url, created_at, updated_at FROM metadata WHERE id = $1`

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  err := m.db.QueryRowContext(ctx, query, id).Scan(
    &metadata.ID,
    &metadata.Title,
    &metadata.Description,
    &metadata.Author,
    &metadata.Url,
    &metadata.CreatedAt,
    &metadata.UpdatedAt,
  )

  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, ErrNotFound
    }
    return nil, err
  }
  return &metadata, nil
}

func (m *MetadataModel) DeleteFile(id int) (*Metadata, error) {
  var metadata Metadata
  query := `DELETE FROM metadata WHERE if = $1 RETURNING id, title, author`

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  err := m.db.QueryRowContext(ctx, query, id).Scan(
    &metadata.ID,
    &metadata.Title,
    &metadata.Author,
  )

  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, ErrNotFound
    }
    return nil, err
  }
  return &metadata, nil
}
