package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Category represents the type of memory
type Category string

const (
	CategoryCore     Category = "core"
	CategorySemantic Category = "semantic"
	CategoryWorking  Category = "working"
	CategoryEpisodic Category = "episodic"
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

// Status represents the verification status of a memory
type Status string

const (
	StatusVerified     Status = "verified"
	StatusUnderReview Status = "under_review"
	StatusDeprecated  Status = "deprecated"
)

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusVerified, StatusUnderReview, StatusDeprecated:
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
	ID              string    `json:"id"`
	Content         string    `json:"content"`
	OwnerID         string    `json:"owner_id"`
	Category        Category  `json:"category"`
	Priority        float64   `json:"priority"`
	CreatedAt       time.Time `json:"created_at"`
	AccessedAt      time.Time `json:"accessed_at"`
	AccessCount     int       `json:"access_count"`
	Source          string    `json:"source,omitempty"`
	Tags            Tags      `json:"tags,omitempty"`
	Embedding       []float32 `json:"embedding,omitempty"`
	EmbeddingCached bool      `json:"embedding_cached"`
	Status          Status    `json:"status,omitempty"`
	TeamID          string    `json:"team_id,omitempty"`
}

// Helper functions for type conversion

func AsString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func AsFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0.0
	}
}

func AsInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

func AsTime(v interface{}) time.Time {
	if v == nil {
		return time.Time{}
	}
	s := AsString(v)
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Time{}
}
