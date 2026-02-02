package store

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hargabyte/ami/internal/db"
	"github.com/hargabyte/ami/internal/models"
)

// Decision represents a tracked decision and its outcome
type Decision struct {
	ID           string    `json:"id"`
	TaskID       string    `json:"task_id"`
	MemoryIDs    []string  `json:"memory_ids"`
	DecisionText string    `json:"decision_text"`
	Outcome      float64   `json:"outcome"`
	Feedback     string    `json:"feedback"`
	CommitHash   string    `json:"commit_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

// TrackDecision tracks a new decision with the memories that informed it
func TrackDecision(taskID string, memoryIDs []string, decisionText string, source string) (*Decision, error) {
	// Generate UUID
	id := uuid.New().String()
	now := time.Now().Format("2006-01-02 15:04:05")

	// Convert memory IDs to JSON
	memoryIDsJSON, err := json.Marshal(memoryIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal memory IDs: %w", err)
	}

	// Escape single quotes
	escapedDecisionText := strings.ReplaceAll(decisionText, "'", "''")

	// Get current commit hash for temporal linking
	commitHash, _ := db.GetHeadCommit()

	// Insert decision
	query := fmt.Sprintf(`
		INSERT INTO decisions (id, task_id, memory_ids, decision_text, created_at, commit_hash)
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s')
	`, id, taskID, string(memoryIDsJSON), escapedDecisionText, now, commitHash)

	_, err = db.ExecDoltSQL(query)
	if err != nil {
		return nil, fmt.Errorf("failed to insert decision: %w", err)
	}

	// Create Dolt commit
	excerpt := decisionText
	if len(excerpt) > 50 {
		excerpt = excerpt[:50] + "..."
	}
	commitMsg := fmt.Sprintf("Track decision: %s", excerpt)

	if err := DoltCommit(commitMsg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create Dolt commit: %v\n", err)
	}

	createdTime, _ := time.Parse("2006-01-02 15:04:05", now)
	return &Decision{
		ID:           id,
		TaskID:       taskID,
		MemoryIDs:    memoryIDs,
		DecisionText: decisionText,
		Outcome:      0.0,
		CommitHash:   commitHash,
		CreatedAt:    createdTime,
	}, nil
}

// RecordOutcome records the outcome of a decision and reinforces linked memories if successful
func RecordOutcome(decisionID string, outcome float64, feedback string) error {
	// Get the decision first to retrieve linked memory IDs
	decision, err := GetDecision(decisionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve decision: %w", err)
	}

	// Escape single quotes in feedback
	escapedFeedback := strings.ReplaceAll(feedback, "'", "''")

	// Update the decision
	query := fmt.Sprintf(`
		UPDATE decisions
		SET outcome = %f, feedback = '%s'
		WHERE id = '%s'
	`, outcome, escapedFeedback, decisionID)

	_, err = db.ExecDoltSQL(query)
	if err != nil {
		return fmt.Errorf("failed to update decision: %w", err)
	}

	// Reinforcement logic: if outcome > 0.8, boost priority of linked memories
	if outcome > 0.8 && len(decision.MemoryIDs) > 0 {
		for _, memID := range decision.MemoryIDs {
			// Increase priority by 0.1
			boostQuery := fmt.Sprintf(`
				UPDATE memories
				SET priority = priority + 0.1, access_count = access_count + 1
				WHERE id = '%s'
			`, memID)
			_, err := db.ExecDoltSQL(boostQuery)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to boost memory %s: %v\n", memID, err)
			}
		}

		// Commit the reinforcement
		commitMsg := fmt.Sprintf("Reinforce memories for decision %s (outcome: %.2f)", decisionID, outcome)
		if err := DoltCommit(commitMsg); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create Dolt commit: %v\n", err)
		}
	}

	// Commit the outcome recording
	commitMsg := fmt.Sprintf("Record outcome for decision %s: %.2f", decisionID, outcome)
	if err := DoltCommit(commitMsg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create Dolt commit: %v\n", err)
	}

	return nil
}

// GetDecision retrieves a decision by ID
func GetDecision(decisionID string) (*Decision, error) {
	query := fmt.Sprintf(`
		SELECT id, task_id, memory_ids, decision_text, outcome, feedback, created_at, commit_hash
		FROM decisions
		WHERE id = '%s'
	`, decisionID)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve decision: %w", err)
	}

	return parseDecisionJSON(output)
}

// ListDecisions retrieves decisions by task ID
func ListDecisions(taskID string) ([]Decision, error) {
	var query string
	if taskID != "" {
		query = fmt.Sprintf(`
			SELECT id, task_id, memory_ids, decision_text, outcome, feedback, created_at, commit_hash
			FROM decisions
			WHERE task_id = '%s'
			ORDER BY created_at DESC
		`, taskID)
	} else {
		query = `
			SELECT id, task_id, memory_ids, decision_text, outcome, feedback, created_at, commit_hash
			FROM decisions
			ORDER BY created_at DESC
		`
	}

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve decisions: %w", err)
	}

	return parseDecisionsJSON(output)
}

// parseDecisionJSON parses a single decision from Dolt SQL JSON output
func parseDecisionJSON(output string) (*Decision, error) {
	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("decision not found")
	}

	row := result.Rows[0]

	var memoryIDs []string
	if memIDsStr := models.AsString(row["memory_ids"]); memIDsStr != "" {
		json.Unmarshal([]byte(memIDsStr), &memoryIDs)
	}

	createdAt := models.AsTime(row["created_at"])

	decision := &Decision{
		ID:           models.AsString(row["id"]),
		TaskID:       models.AsString(row["task_id"]),
		MemoryIDs:    memoryIDs,
		DecisionText: models.AsString(row["decision_text"]),
		Feedback:     models.AsString(row["feedback"]),
		CommitHash:   models.AsString(row["commit_hash"]),
		CreatedAt:    createdAt,
		Outcome:      models.AsFloat64(row["outcome"]),
	}

	return decision, nil
}

// parseDecisionsJSON parses multiple decisions from Dolt SQL JSON output
func parseDecisionsJSON(output string) ([]Decision, error) {
	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	decisions := make([]Decision, 0, len(result.Rows))
	for _, row := range result.Rows {
		var memoryIDs []string
		if memIDsStr := models.AsString(row["memory_ids"]); memIDsStr != "" {
			json.Unmarshal([]byte(memIDsStr), &memoryIDs)
		}

		decision := Decision{
			ID:           models.AsString(row["id"]),
			TaskID:       models.AsString(row["task_id"]),
			MemoryIDs:    memoryIDs,
			DecisionText: models.AsString(row["decision_text"]),
			Feedback:     models.AsString(row["feedback"]),
			CommitHash:   models.AsString(row["commit_hash"]),
			CreatedAt:    models.AsTime(row["created_at"]),
			Outcome:      models.AsFloat64(row["outcome"]),
		}

		decisions = append(decisions, decision)
	}

	return decisions, nil
}
