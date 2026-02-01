package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Category represents the type of memory
type Category string

const (
	CategoryCore      Category = "core"
	CategorySemantic   Category = "semantic"
	CategoryWorking    Category = "working"
	CategoryEpisodic   Category = "episodic"
)

// IsValid checks if the category is valid
func (c Category) IsValid() bool {
	switch c {
	case CategoryCore, CategorySemantic, CategoryWorking, CategoryEpisodic:
		return true
	default:
		return false
	}
}

// Tags is a JSON-able slice of strings
type Tags []string

// Scan implements sql.Scanner for Tags
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, t)
}

// Value implements driver.Valuer for Tags
func (t Tags) Value() (driver.Value, error) {
	if len(t) == 0 {
		return "[]", nil
	}
	return json.Marshal(t)
}

// Memory represents a stored memory
type Memory struct {
	ID          string     `json:"id"`
	Content     string     `json:"content"`
	OwnerID     string     `json:"owner_id"`
	Category    Category   `json:"category"`
	Priority    float64    `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	AccessedAt  time.Time  `json:"accessed_at"`
	AccessCount int        `json:"access_count"`
	Source      string     `json:"source,omitempty"`
	Tags        Tags       `json:"tags,omitempty"`
	Embedding   []float32  `json:"embedding,omitempty"`
}
