package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

// InitDB initializes the connection to the local Dolt database
// We use the dolt CLI directly via shell commands for simplicity
func InitDB(repoPath string) error {
	// For Phase 1, we'll use dolt CLI directly via shell commands
	// This is simpler and avoids the sql-server requirement
	fmt.Fprintf(os.Stderr, "[DB] Using Dolt CLI directly (repo: %s)\n", repoPath)
	return nil
}

// ExecDoltSQL executes a SQL command via the dolt CLI
func ExecDoltSQL(sql string) (string, error) {
	repoPath, err := GetRepoPath()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("dolt", "sql", "-q", sql, "-r", "json")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("dolt sql failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// GetDB returns the database connection (nil for CLI mode)
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetRepoPath returns the path to the Dolt repository
func GetRepoPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree to find the .dolt directory
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, ".dolt")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a Dolt repository")
		}
		dir = parent
	}
}

// DoltCommit creates a Dolt commit for the current database state
func DoltCommit(message string) error {
	repoPath, err := GetRepoPath()
	if err != nil {
		return fmt.Errorf("failed to get repo path: %w", err)
	}

	// Add all changes
	cmd := exec.Command("dolt", "add", "-A")
	cmd.Dir = repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("dolt add failed: %w\nOutput: %s", err, string(output))
	}

	// Commit
	cmd = exec.Command("dolt", "commit", "-m", message)
	cmd.Dir = repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		// Ignore "nothing to commit" errors
		if !strings.Contains(string(output), "nothing to commit") {
			return fmt.Errorf("dolt commit failed: %w\nOutput: %s", err, string(output))
		}
	}

	return nil
}
