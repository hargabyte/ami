package store

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hargabyte/ami/internal/db"
	"github.com/hargabyte/ami/internal/models"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

// CatchupOptions specifies filters for memory catchup
type CatchupOptions struct {
	Limit    int
	Category string
	Since    string
}

// RecallOptions specifies filters for memory recall
type RecallOptions struct {
	Query      string
	Limit      int
	Tags       []string
	Category   string
	OwnerID    string
	TeamID     string
	WithDecay  bool
	Semantic   bool
}

// UpdateParams specifies fields to update on a memory
type UpdateParams struct {
	ID       string
	Content  *string
	OwnerID  *string
	Category *models.Category
	Priority *float64
	Source   *string
	Tags     []string
}

// DoltCommit is a wrapper for db.DoltCommit
func DoltCommit(message string) error {
	return db.DoltCommit(message)
}

// ExecDoltSQLJSON executes a SQL query via dolt CLI and returns JSON output
func ExecDoltSQLJSON(query string) (string, error) {
	repoPath, err := db.GetRepoPath()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("dolt", "sql", "-q", query, "-r", "json")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("dolt sql failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// AddMemory adds a new memory to the database and creates a Dolt commit
func AddMemory(content string, ownerID string, category models.Category, priority float64, tags []string, source string, teamID string) (*models.Memory, error) {
	// Generate UUID
	id := uuid.New().String()
	now := time.Now().Format("2006-01-02 15:04:05")

	// Set default owner if empty
	if ownerID == "" {
		ownerID = "system"
	}

	// Set default team if empty
	if teamID == "" {
		teamID = "system"
	}

	// Convert tags to JSON for SQL
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	// Escape single quotes in content
	escapedContent := strings.ReplaceAll(content, "'", "''")
	escapedSource := strings.ReplaceAll(source, "'", "''")

	// 1. Calculate embedding if enabled (v0.4.0)
	embeddingHex := "NULL"
	if os.Getenv("OPENAI_API_KEY") != "" {
		vector, err := GetEmbedding(content)
		if err == nil {
			binaryData := Float32ToBinary(vector)
			embeddingHex = fmt.Sprintf("X'%x'", binaryData)
		}
	}

	// 2. Insert memory using dolt CLI
	query := fmt.Sprintf(`
		INSERT INTO memories (id, content, owner_id, category, priority, created_at, accessed_at, access_count, source, tags, embedding, team_id)
		VALUES ('%s', '%s', '%s', '%s', %f, '%s', '%s', 0, '%s', '%s', %s, '%s')
	`, id, escapedContent, ownerID, string(category), priority, now, now, escapedSource, string(tagsJSON), embeddingHex, teamID)

	_, err = db.ExecDoltSQL(query)
	if err != nil {
		return nil, fmt.Errorf("failed to insert memory: %w", err)
	}

	// Create Dolt commit for versioning
	excerpt := content
	if len(excerpt) > 50 {
		excerpt = excerpt[:50] + "..."
	}
	commitMsg := fmt.Sprintf("Add memory: %s", excerpt)

	if err := DoltCommit(commitMsg); err != nil {
		// Log warning but don't fail the memory add
		fmt.Fprintf(os.Stderr, "Warning: failed to create Dolt commit: %v\n", err)
	}

	// Return the created memory
	createdTime, _ := time.Parse("2006-01-02 15:04:05", now)
	return &models.Memory{
		ID:          id,
		Content:     content,
		OwnerID:     ownerID,
		Category:    category,
		Priority:    priority,
		CreatedAt:   createdTime,
		AccessedAt:  createdTime,
		AccessCount: 0,
		Source:      source,
		Tags:        models.Tags(tags),
		TeamID:      teamID,
	}, nil
}

// CatchupMemories returns the most recent memories
func CatchupMemories(opts CatchupOptions) ([]models.Memory, error) {
	whereClauses := []string{}

	if opts.Category != "" {
		cat := models.Category(opts.Category)
		if cat.IsValid() {
			whereClauses = append(whereClauses, fmt.Sprintf("category = '%s'", string(cat)))
		}
	}

	if opts.Since != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= '%s'", opts.Since))
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, content, category, priority, created_at, accessed_at, access_count, source, tags, status
		FROM memories
		%s
		ORDER BY created_at DESC
		LIMIT %d
	`, whereClause, opts.Limit)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, err
	}

	return parseMemoriesJSON(output)
}

// RecallMemories performs a basic text search on memories with optional filters
func RecallMemories(opts RecallOptions) ([]models.Memory, error) {
	// 1. Build WHERE clause
	whereClauses := []string{}
	// ... (rest of where clause building)

	// Text search
	if opts.Query != "" {
		escapedQuery := strings.ReplaceAll(opts.Query, "'", "''")
		whereClauses = append(whereClauses, fmt.Sprintf("content LIKE '%%%s%%'", escapedQuery))
	}

	// Category filter
	if opts.Category != "" {
		cat := models.Category(opts.Category)
		if cat.IsValid() {
			whereClauses = append(whereClauses, fmt.Sprintf("category = '%s'", string(cat)))
		}
	}

	// Owner filter
	if opts.OwnerID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("owner_id = '%s'", opts.OwnerID))
	}

	// Team filter
	if opts.TeamID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("team_id = '%s'", opts.TeamID))
	}

	// Tags filter - check JSON_CONTAINS
	if len(opts.Tags) > 0 {
		for _, tag := range opts.Tags {
			escapedTag := strings.ReplaceAll(tag, "'", "''")
			whereClauses = append(whereClauses, fmt.Sprintf("JSON_CONTAINS(tags, '\"%s\"')", escapedTag))
		}
	}

	// Combine WHERE clauses
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Build query
	var searchQuery string
	if opts.Semantic {
		// Fetch all memories with embeddings for in-memory ranking
		searchQuery = fmt.Sprintf(`
			SELECT id, content, owner_id, category, priority, created_at, accessed_at, access_count, source, tags, embedding, embedding_cached, status
			FROM memories
			%s
		`, whereClause)
	} else if opts.WithDecay {
		// Use logarithmic decay scoring:
		// Score = (Priority * (AccessCount + 1)) / (log10(TimeDelta + 10) * CategoryDecay)
		searchQuery = fmt.Sprintf(`
			SELECT id, content, owner_id, category, priority, created_at, accessed_at, access_count, source, tags, status,
			(priority * (access_count + 1)) / (LOG10(TIMESTAMPDIFF(SECOND, accessed_at, NOW()) + 10) *
			CASE
				WHEN category = 'core' THEN 0.5
				WHEN category = 'semantic' THEN 1.0
				WHEN category = 'episodic' THEN 2.0
				ELSE 1.5
			END) as recall_score
			FROM memories
			%s
			ORDER BY recall_score DESC
			LIMIT %d
		`, whereClause, opts.Limit)
	} else {
		searchQuery = fmt.Sprintf(`
			SELECT id, content, owner_id, category, priority, created_at, accessed_at, access_count, source, tags, status
			FROM memories
			%s
			ORDER BY priority DESC, accessed_at DESC
			LIMIT %d
		`, whereClause, opts.Limit)
	}

	output, err := ExecDoltSQLJSON(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %w", err)
	}

	// Parse JSON output
	memories, err := parseMemoriesJSON(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse memories: %w", err)
	}

	if opts.Semantic && opts.Query != "" {
		// 1. Get embedding for the query
		queryVector, err := GetEmbedding(opts.Query)
		if err != nil {
			return nil, fmt.Errorf("failed to get query embedding: %w", err)
		}

		// 2. Rank memories by similarity
		type memoryWithScore struct {
			models.Memory
			score float32
		}
		var ranked []memoryWithScore
		for _, m := range memories {
			if len(m.Embedding) > 0 {
				score := CosineSimilarity(queryVector, m.Embedding)
				ranked = append(ranked, memoryWithScore{m, score})
			}
		}

		// Sort by score
		sort.Slice(ranked, func(i, j int) bool {
			return ranked[i].score > ranked[j].score
		})

		// Convert back to []models.Memory and apply limit
		finalMemories := make([]models.Memory, 0, len(ranked))
		for i := 0; i < len(ranked) && i < opts.Limit; i++ {
			finalMemories = append(finalMemories, ranked[i].Memory)
		}
		return finalMemories, nil
	}

	return memories, nil
}

// MemoryHistory represents a version of a memory in history
type MemoryHistory struct {
	models.Memory
	CommitHash string    `json:"commit_hash"`
	Committer  string    `json:"committer"`
	CommitDate time.Time `json:"commit_date"`
}

// GetMemoryHistory returns the version history of a memory
func GetMemoryHistory(id string) ([]MemoryHistory, error) {
	query := fmt.Sprintf(`
		SELECT id, content, category, priority, created_at, accessed_at, access_count, source, tags, commit_hash, committer, commit_date
		FROM dolt_history_memories
		WHERE id = '%s'
		ORDER BY commit_date DESC
	`, id)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, err
	}

	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, err
	}

	var history []MemoryHistory
	for _, row := range result.Rows {
		var h MemoryHistory
		h.ID = asString(row["id"])
		h.Content = asString(row["content"])
		h.Category = models.Category(asString(row["category"]))
		h.Priority = asFloat64(row["priority"])
		h.CreatedAt = asTime(row["created_at"])
		h.AccessedAt = asTime(row["accessed_at"])
		h.AccessCount = asInt(row["access_count"])
		h.Source = asString(row["source"])
		h.CommitHash = asString(row["commit_hash"])
		h.Committer = asString(row["committer"])
		h.CommitDate = asTime(row["commit_date"])

		tagsJSON := asString(row["tags"])
		var tags models.Tags
		if tagsJSON != "" && tagsJSON != "[]" {
			json.Unmarshal([]byte(tagsJSON), &tags)
		}
		h.Tags = tags

		history = append(history, h)
	}

	return history, nil
}

// LinkMemories creates a link between two memories
func LinkMemories(fromID, toID, relation string) error {
	query := fmt.Sprintf(`
		INSERT INTO memory_links (from_id, to_id, relation)
		VALUES ('%s', '%s', '%s')
		ON DUPLICATE KEY UPDATE relation = VALUES(relation)
	`, fromID, toID, relation)

	if _, err := db.ExecDoltSQL(query); err != nil {
		return err
	}

	commitMsg := fmt.Sprintf("Link memory %s to %s (%s)", fromID, toID, relation)
	return DoltCommit(commitMsg)
}

// GetMemoryLinks returns all links for a specific memory
func GetMemoryLinks(id string) ([]map[string]string, error) {
	query := fmt.Sprintf(`
		SELECT from_id, to_id, relation
		FROM memory_links
		WHERE from_id = '%s' OR to_id = '%s'
	`, id, id)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, err
	}

	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, err
	}

	var links []map[string]string
	for _, row := range result.Rows {
		links = append(links, map[string]string{
			"from_id":  asString(row["from_id"]),
			"to_id":    asString(row["to_id"]),
			"relation": asString(row["relation"]),
		})
	}

	return links, nil
}

// GetKeystoneMemories returns high-priority and high-access memories
func GetKeystoneMemories(limit int) ([]models.Memory, error) {
	// Formula: (Priority * 2) + (AccessCount / 10)
	query := fmt.Sprintf(`
		SELECT id, content, category, priority, created_at, accessed_at, access_count, source, tags
		FROM memories
		ORDER BY (priority * 2) + (access_count / 10.0) DESC
		LIMIT %d
	`, limit)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, err
	}

	return parseMemoriesJSON(output)
}

// CountTokens counts tokens in a string
func CountTokens(text string) int {
	// Initialize tiktoken for Claude/GPT-4 encoding
	encoding := "cl100k_base"
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		// Fallback to rough estimate: 1 token ~= 4 chars or 0.75 words
		return len(text) / 4
	}
	token := tke.Encode(text, nil, nil)
	return len(token)
}

// GetContextMemories returns memories optimized for prompt context
func GetContextMemories(task string, limit int, tokenBudget int) ([]models.Memory, error) {
	// 1. Get high-priority core facts first
	coreOpts := RecallOptions{
		Category: "core",
		Limit:    10,
	}
	coreMemories, _ := RecallMemories(coreOpts)

	// 2. Get task-relevant memories with semantic search if task is provided
	var taskMemories []models.Memory
	if task != "" {
		taskOpts := RecallOptions{
			Query:     task,
			Limit:     limit,
			WithDecay: true,
			Semantic:  true,
		}
		taskMemories, _ = RecallMemories(taskOpts)
	}

	// 3. Pack memories into the budget
	seen := make(map[string]bool)
	var final []models.Memory
	currentTokens := 0

	// Pack Core first
	for _, m := range coreMemories {
		if !seen[m.ID] {
			tokens := CountTokens(m.Content)
			if currentTokens+tokens <= tokenBudget {
				final = append(final, m)
				seen[m.ID] = true
				currentTokens += tokens
			}
		}
	}

	// Pack Semantic/Task context
	for _, m := range taskMemories {
		if !seen[m.ID] {
			tokens := CountTokens(m.Content)
			if currentTokens+tokens <= tokenBudget {
				final = append(final, m)
				seen[m.ID] = true
				currentTokens += tokens
			}
		}
	}

	return final, nil
}

// PromoteMemory moves a memory from local store to global store
func PromoteMemory(id string, globalStorePath string) error {
	// 1. Get memory from local
	query := fmt.Sprintf("SELECT content, owner_id, category, priority, source, tags FROM memories WHERE id = '%s'", id)
	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return err
	}

	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	json.Unmarshal([]byte(output), &result)

	if len(result.Rows) == 0 {
		return fmt.Errorf("memory %s not found in local store", id)
	}

	row := result.Rows[0]
	content := asString(row["content"])
	owner := asString(row["owner_id"])
	category := asString(row["category"])
	priority := asFloat64(row["priority"])
	source := asString(row["source"])
	tagsJSON := asString(row["tags"])

	// Handle tags correctly
	var tags []string
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		// If it's not a JSON array string, it might be a raw string from tabular output
		// But ExecDoltSQLJSON should return JSON. 
		// Actually, dolt_history sometimes returns comma-separated strings or JSON.
		// Let's just try to ensure it's valid JSON for the insert.
	}
	tagsBytes, _ := json.Marshal(tags)
	finalTags := string(tagsBytes)

	// 2. Add to global
	insertQuery := fmt.Sprintf(`
		INSERT INTO memories (id, content, owner_id, category, priority, created_at, accessed_at, access_count, source, tags)
		VALUES ('%s', '%s', '%s', '%s', %f, NOW(), NOW(), 0, '%s', '%s')
		ON DUPLICATE KEY UPDATE content = VALUES(content), priority = VALUES(priority)
	`, id, strings.ReplaceAll(content, "'", "''"), owner, category, priority, strings.ReplaceAll(source, "'", "''"), finalTags)

	cmd := exec.Command("dolt", "sql", "-q", insertQuery)
	cmd.Dir = globalStorePath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to insert into global store: %w\nOutput: %s", err, string(output))
	}

	// 3. Commit global
	commitCmd := exec.Command("dolt", "commit", "-am", fmt.Sprintf("Promoted memory %s from project store", id))
	commitCmd.Dir = globalStorePath
	commitCmd.Run() // Ignore "nothing to commit" errors

	return nil
}

// BinaryToFloat32 converts a binary BLOB to []float32
func BinaryToFloat32(data []byte) []float32 {
	floats := make([]float32, len(data)/4)
	for i := 0; i < len(floats); i++ {
		bits := binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
		floats[i] = math.Float32frombits(bits)
	}
	return floats
}

// Float32ToBinary converts []float32 to a binary BLOB
func Float32ToBinary(floats []float32) []byte {
	data := make([]byte, len(floats)*4)
	for i, f := range floats {
		binary.LittleEndian.PutUint32(data[i*4:(i+1)*4], math.Float32bits(f))
	}
	return data
}

// GetEmbedding fetches the embedding vector for a string using OpenAI
func GetEmbedding(text string) ([]float32, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)
	resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

// CosineSimilarity calculates the similarity between two vectors
func CosineSimilarity(a, b []float32) float32 {
	var dotProduct float32
	var normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// UpdateMemory updates an existing memory
func UpdateMemory(params UpdateParams) error {
	// Build SET clause
	var setClauses []string

	if params.Content != nil {
		escapedContent := strings.ReplaceAll(*params.Content, "'", "''")
		setClauses = append(setClauses, fmt.Sprintf("content = '%s'", escapedContent))
	}

	if params.OwnerID != nil {
		setClauses = append(setClauses, fmt.Sprintf("owner_id = '%s'", *params.OwnerID))
	}

	if params.Category != nil {
		setClauses = append(setClauses, fmt.Sprintf("category = '%s'", string(*params.Category)))
	}

	if params.Priority != nil {
		setClauses = append(setClauses, fmt.Sprintf("priority = %f", *params.Priority))
	}

	if params.Source != nil {
		escapedSource := strings.ReplaceAll(*params.Source, "'", "''")
		setClauses = append(setClauses, fmt.Sprintf("source = '%s'", escapedSource))
	}

	if params.Tags != nil {
		tagsJSON, err := json.Marshal(params.Tags)
		if err != nil {
			return fmt.Errorf("failed to marshal tags: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("tags = '%s'", string(tagsJSON)))
	}

	// Update accessed_at to refresh timestamp
	now := time.Now().Format("2006-01-02 15:04:05")
	setClauses = append(setClauses, fmt.Sprintf("accessed_at = '%s'", now))

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields specified for update")
	}

	// Build UPDATE query
	updateQuery := fmt.Sprintf(`
		UPDATE memories
		SET %s
		WHERE id = '%s'
	`, strings.Join(setClauses, ", "), params.ID)

	if _, err := db.ExecDoltSQL(updateQuery); err != nil {
		return fmt.Errorf("failed to update memory: %w", err)
	}

	// Create Dolt commit
	commitMsg := fmt.Sprintf("Update memory: %s", params.ID)
	if err := DoltCommit(commitMsg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create Dolt commit: %v\n", err)
	}

	return nil
}

// RollbackMemory rolls back a memory to a specific commit
func RollbackMemory(id string, commitHash string) error {
	// 1. Get the content/metadata from history for that commit
	query := fmt.Sprintf(`
		SELECT content, category, priority, source, tags
		FROM dolt_history_memories
		WHERE id = '%s' AND commit_hash = '%s'
		LIMIT 1
	`, id, commitHash)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return err
	}

	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return err
	}

	if len(result.Rows) == 0 {
		return fmt.Errorf("no history found for id %s and commit %s", id, commitHash)
	}

	row := result.Rows[0]
	content := asString(row["content"])
	category := asString(row["category"])
	priority := asFloat64(row["priority"])
	source := asString(row["source"])
	tagsJSON := asString(row["tags"])

	// 2. Update the current memories table
	updateQuery := fmt.Sprintf(`
		UPDATE memories
		SET content = '%s', category = '%s', priority = %f, source = '%s', tags = '%s', accessed_at = NOW()
		WHERE id = '%s'
	`, strings.ReplaceAll(content, "'", "''"), category, priority, strings.ReplaceAll(source, "'", "''"), tagsJSON, id)

	if _, err := db.ExecDoltSQL(updateQuery); err != nil {
		return err
	}

	// 3. Commit the rollback
	commitMsg := fmt.Sprintf("Rollback memory %s to commit %s", id, commitHash)
	return DoltCommit(commitMsg)
}

// DeleteMemory deletes a memory by ID
func DeleteMemory(id string) error {
	query := fmt.Sprintf("DELETE FROM memories WHERE id = '%s'", id)

	if _, err := db.ExecDoltSQL(query); err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	// Create Dolt commit
	commitMsg := fmt.Sprintf("Delete memory: %s", id)
	if err := DoltCommit(commitMsg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create Dolt commit: %v\n", err)
	}

	return nil
}

// ListTags returns all unique tags in the database
func ListTags() ([]string, error) {
	// Use JSON_OVERLAPS or similar if supported, but simple way is to select all and unique in Go
	// Better: SELECT DISTINCT JSON_UNQUOTE(JSON_EXTRACT(tags, '$[*]')) FROM memories
	// Dolt supports JSON functions. Let's try to flatten.
	// Actually, easier to fetch all tags and unique them in Go for now.
	query := "SELECT tags FROM memories WHERE tags IS NOT NULL AND tags != '[]'"
	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, err
	}

	tagMap := make(map[string]bool)
	for _, row := range result.Rows {
		v, ok := row["tags"]
		if !ok || v == nil {
			continue
		}

		var tags []string
		switch val := v.(type) {
		case string:
			if val != "" && val != "[]" {
				json.Unmarshal([]byte(val), &tags)
			}
		case []interface{}:
			for _, t := range val {
				tags = append(tags, fmt.Sprintf("%v", t))
			}
		}

		for _, tag := range tags {
			tagMap[tag] = true
		}
	}

	uniqueTags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		uniqueTags = append(uniqueTags, tag)
	}

	return uniqueTags, nil
}

// GetMemoryStats returns analytics about the memory database
func GetMemoryStats() (map[string]interface{}, error) {
	// Category distribution
	distQuery := `
		SELECT category, COUNT(*) as count 
		FROM memories 
		GROUP BY category
	`
	output, err := ExecDoltSQLJSON(distQuery)
	if err != nil {
		return nil, err
	}

	var distResult struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	json.Unmarshal([]byte(output), &distResult)

	distribution := make(map[string]int)
	total := 0
	for _, row := range distResult.Rows {
		cat := asString(row["category"])
		count := asInt(row["count"])
		distribution[cat] = count
		total += count
	}

	// Average priority and decay
	metricsQuery := `
		SELECT 
			AVG(priority) as avg_priority,
			AVG(access_count) as avg_access,
			AVG(recall_score) as avg_decay_score
		FROM (
			SELECT priority, access_count,
			(priority * (access_count + 1)) / (LOG10(TIMESTAMPDIFF(SECOND, accessed_at, NOW()) + 10) * 
			CASE 
				WHEN category = 'core' THEN 0.5 
				WHEN category = 'semantic' THEN 1.0 
				WHEN category = 'episodic' THEN 2.0 
				ELSE 1.5 
			END) as recall_score
			FROM memories
		) as metrics
	`
	output, err = ExecDoltSQLJSON(metricsQuery)
	if err != nil {
		return nil, err
	}

	var metricsResult struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	json.Unmarshal([]byte(output), &metricsResult)

	avgPriority := 0.0
	avgAccess := 0.0
	avgDecay := 0.0
	if len(metricsResult.Rows) > 0 {
		row := metricsResult.Rows[0]
		avgPriority = asFloat64(row["avg_priority"])
		avgAccess = asFloat64(row["avg_access"])
		avgDecay = asFloat64(row["avg_decay_score"])
	}

	return map[string]interface{}{
		"total_memories": total,
		"distribution":   distribution,
		"metrics": map[string]interface{}{
			"avg_priority":     avgPriority,
			"avg_access_count": avgAccess,
			"avg_decay_score":  avgDecay,
		},
	}, nil
}

// GetMemoryCount returns total number of memories
func GetMemoryCount() (int, error) {
	output, err := ExecDoltSQLJSON("SELECT COUNT(*) as count FROM memories")
	if err != nil {
		return 0, fmt.Errorf("failed to count memories: %w", err)
	}

	// Parse JSON output
	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return 0, fmt.Errorf("failed to parse count result: %w\nOutput: %s", err, output)
	}

	if len(result.Rows) == 0 {
		return 0, nil
	}

	// Get the count from the first column
	for _, v := range result.Rows[0] {
		return asInt(v), nil
	}

	return 0, nil
}

// parseMemoriesJSON parses JSON output from Dolt
func parseMemoriesJSON(output string) ([]models.Memory, error) {
	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w\nOutput: %s", err, output)
	}

	var memories []models.Memory
	for _, row := range result.Rows {
		var m models.Memory

		// Parse row values by column name
		m.ID = asString(row["id"])
		m.Content = asString(row["content"])
		m.OwnerID = asString(row["owner_id"])
		m.Category = models.Category(asString(row["category"]))
		m.Priority = asFloat64(row["priority"])
		m.CreatedAt = asTime(row["created_at"])
		m.AccessedAt = asTime(row["accessed_at"])
		m.AccessCount = asInt(row["access_count"])
		m.Source = asString(row["source"])
		m.EmbeddingCached = asInt(row["embedding_cached"]) == 1
		m.Status = models.Status(asString(row["status"]))
		m.TeamID = asString(row["team_id"])

		// Parse embedding if present
		if row["embedding"] != nil {
			embeddingStr := asString(row["embedding"])
			if data, err := base64.StdEncoding.DecodeString(embeddingStr); err == nil {
				m.Embedding = BinaryToFloat32(data)
			}
		}

		// Parse tags JSON
		tagsJSON := asString(row["tags"])
		var tags models.Tags
		if tagsJSON != "" && tagsJSON != "[]" {
			json.Unmarshal([]byte(tagsJSON), &tags)
		}
		m.Tags = tags

		memories = append(memories, m)
	}

	return memories, nil
}

// Helper functions for type conversion

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func asFloat64(v interface{}) float64 {
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

func asInt(v interface{}) int {
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

func asTime(v interface{}) time.Time {
	if v == nil {
		return time.Time{}
	}
	s := asString(v)
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Time{}
}

// FindAutoPromotionCandidates finds memories eligible for promotion to global brain
// Criteria: high access count, linked to successful decisions, semantic or core category
func FindAutoPromotionCandidates(minAccessCount int, minOutcome float64) ([]models.Memory, error) {
	// Find memories that were linked to successful decisions
	query := fmt.Sprintf(`
		SELECT DISTINCT m.id, m.content, m.owner_id, m.category, m.priority,
		       m.created_at, m.accessed_at, m.access_count, m.source, m.tags, m.status
		FROM memories m
		JOIN decisions d ON JSON_CONTAINS(d.memory_ids, CONCAT('"', m.id, '"'))
		WHERE m.access_count >= %d
		  AND d.outcome >= %f
		  AND m.category IN ('semantic', 'core')
		  AND m.status = 'verified'
		ORDER BY m.access_count DESC, m.priority DESC
	`, minAccessCount, minOutcome)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query promotion candidates: %w", err)
	}

	return parseMemoriesJSON(output)
}

// GetMemoryByID retrieves a specific memory by ID
func GetMemoryByID(id string) (*models.Memory, error) {
	query := fmt.Sprintf(`
		SELECT id, content, owner_id, category, priority, created_at, accessed_at, access_count, source, tags, status
		FROM memories
		WHERE id = '%s'
	`, id)

	output, err := ExecDoltSQLJSON(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve memory: %w", err)
	}

	memories, err := parseMemoriesJSON(output)
	if err != nil {
		return nil, err
	}

	if len(memories) == 0 {
		return nil, fmt.Errorf("memory not found: %s", id)
	}

	return &memories[0], nil
}

// UpdateMemoryStatus updates the status of a memory
func UpdateMemoryStatus(id string, status models.Status) error {
	query := fmt.Sprintf(`
		UPDATE memories
		SET status = '%s'
		WHERE id = '%s'
	`, string(status), id)

	_, err := ExecDoltSQLJSON(query)
	if err != nil {
		return fmt.Errorf("failed to update memory status: %w", err)
	}

	commitMsg := fmt.Sprintf("Update status of memory %s to %s", id, status)
	DoltCommit(commitMsg)

	return nil
}

// UpdateMemoryContent updates the content of a memory
func UpdateMemoryContent(id string, content string) error {
	escapedContent := strings.ReplaceAll(content, "'", "''")
	query := fmt.Sprintf(`
		UPDATE memories
		SET content = '%s'
		WHERE id = '%s'
	`, escapedContent, id)

	_, err := ExecDoltSQLJSON(query)
	if err != nil {
		return fmt.Errorf("failed to update memory content: %w", err)
	}

	commitMsg := fmt.Sprintf("Update content of memory %s", id)
	DoltCommit(commitMsg)

	return nil
}
