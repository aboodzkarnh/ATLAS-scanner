package storage

import (
    "fmt"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// ScanRepo defines the interface for scan storage operations
type ScanRepo interface {
    Save(scan interface{}) error
}

// FindingRepo defines the interface for finding storage operations
type FindingRepo interface {
    Save(finding interface{}) error
}

// EvidenceStore defines the interface for evidence file storage
type EvidenceStore interface {
    Store(data []byte, filename string) (string, error)
}

// SQLiteAdapter implements ScanRepo and FindingRepo using SQLite
type SQLiteAdapter struct {
    DB *gorm.DB
}

// NewSQLiteAdapter creates a new SQLite adapter
func NewSQLiteAdapter(dbPath string) (*SQLiteAdapter, error) {
    if dbPath == "" {
        return nil, fmt.Errorf("database path is empty")
    }
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("cannot open database: %w", err)
    }
    return &SQLiteAdapter{DB: db}, nil
}

// Save stores a scan (stub)
func (a *SQLiteAdapter) Save(scan interface{}) error {
    return a.DB.Create(scan).Error
}
