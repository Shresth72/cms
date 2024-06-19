package data

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("requested resource not found")
)

type Models struct {
	MetadataModel *MetadataModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		MetadataModel: &MetadataModel{db: db},
	}
}
