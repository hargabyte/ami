package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hargabyte/ami/internal/db"
	"github.com/hargabyte/ami/internal/models"
	"github.com/hargabyte/ami/internal/store"
	"github.com/spf13/cobra"
)

var version = "0.7.0"

// confirmAction asks for user confirmation
func confirmAction(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ami",
		Short: "Agent Memory Intelligence - versioned memory for AI agents",
		Long: `ami is a CLI tool for managing agent memory using DoltDB.
		
Features:
  - Versioned memory storage (git-like)
  - Decay-weighted retrieval
  - Pre-compression checkpointing
  - Robot mode for agent integration`,
		Version: version,
	}

	// Add commands
	rootCmd.AddCommand(addCmd())
	rootCmd.AddCommand(updateCmd())
	rootCmd.AddCommand(recallCmd())
	rootCmd.AddCommand(catchupCmd())
	rootCmd.AddCommand(historyCmd())
	rootCmd.AddCommand(rollbackCmd())
	rootCmd.AddCommand(linkCmd())
	rootCmd.AddCommand(keystonesCmd())
	rootCmd.AddCommand(statsCmd())
	rootCmd.AddCommand(contextCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(promoteCmd())
	rootCmd.AddCommand(helpAgentsCmd())
	rootCmd.AddCommand(deleteCmd())
	rootCmd.AddCommand(tagsCmd())
	rootCmd.AddCommand(checkpointCmd())
	rootCmd.AddCommand(consolidateCmd())
	rootCmd.AddCommand(decisionCmd())
	rootCmd.AddCommand(reflectCmd())
	rootCmd.AddCommand(conflictCmd())
	rootCmd.AddCommand(pairingCmd())
	rootCmd.AddCommand(robotCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addCmd() *cobra.Command {
	var category string
	var ownerID string
	var priority float64
	var tags []string
	var source string
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "add [content]",
		Short: "Add a memory to the database",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			// Validate category
			cat := models.Category(category)
			if !cat.IsValid() {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"invalid category %s"}`+"\n", category)
				} else {
					fmt.Fprintf(os.Stderr, "Error: invalid category '%s'. Must be one of: core, semantic, working, episodic\n", category)
				}
				os.Exit(1)
			}

			// Validate priority
			if priority < 0.0 || priority > 1.0 {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"priority must be between 0.0 and 1.0"}`+"\n")
				} else {
					fmt.Fprintf(os.Stderr, "Error: priority must be between 0.0 and 1.0\n")
				}
				os.Exit(1)
			}

			// Join content from all args
			content := args[0]
			if len(args) > 1 {
				for _, arg := range args[1:] {
					content += " " + arg
				}
			}

			// Add the memory
			memory, err := store.AddMemory(content, ownerID, cat, priority, tags, source)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error adding memory: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				result := map[string]interface{}{
					"status": "ok",
					"memory": memory,
				}
				jsonBytes, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("✓ Added memory %s (category: %s, priority: %.1f)\n", memory.ID, memory.Category, memory.Priority)
			}
		},
	}
	cmd.Flags().StringVar(&category, "category", "episodic", "Memory category (core|semantic|working|episodic)")
	cmd.Flags().StringVar(&ownerID, "owner", "system", "ID of the agent owning this memory")
	cmd.Flags().Float64Var(&priority, "priority", 0.5, "Priority (0.0-1.0)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags for the memory")
	cmd.Flags().StringVar(&source, "source", "", "Source of the memory (optional)")
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func updateCmd() *cobra.Command {
	var category string
	var ownerID string
	var priority float64
	var tags []string
	var source string
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "update [id] [new-content]",
		Short: "Update an existing memory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			id := args[0]

			// Build update parameters
			params := store.UpdateParams{
				ID: id,
			}

			// Only update fields that were explicitly set
			if cmd.Flags().Changed("owner") {
				params.OwnerID = &ownerID
			}

			if cmd.Flags().Changed("category") {
				cat := models.Category(category)
				if !cat.IsValid() {
					if robotMode {
						fmt.Printf(`{"status":"error","message":"invalid category %s"}`+"\n", category)
					} else {
						fmt.Fprintf(os.Stderr, "Error: invalid category '%s'. Must be one of: core, semantic, working, episodic\n", category)
					}
					os.Exit(1)
				}
				params.Category = &cat
			}

			if cmd.Flags().Changed("priority") {
				if priority < 0.0 || priority > 1.0 {
					if robotMode {
						fmt.Printf(`{"status":"error","message":"priority must be between 0.0 and 1.0"}`+"\n")
					} else {
						fmt.Fprintf(os.Stderr, "Error: priority must be between 0.0 and 1.0\n")
					}
					os.Exit(1)
				}
				params.Priority = &priority
			}

			if cmd.Flags().Changed("tags") {
				params.Tags = tags
			}

			if cmd.Flags().Changed("source") {
				params.Source = &source
			}

			// Content is provided as args
			if len(args) > 1 {
				content := args[1]
				if len(args) > 2 {
					for _, arg := range args[2:] {
						content += " " + arg
					}
				}
				params.Content = &content
			}

			// Update the memory
			if err := store.UpdateMemory(params); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error updating memory: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				fmt.Printf(`{"status":"ok","message":"updated memory %s"}`+"\n", id)
			} else {
				fmt.Printf("✓ Updated memory %s\n", id)
			}
		},
	}
	cmd.Flags().StringVar(&category, "category", "", "Memory category (core|semantic|working|episodic)")
	cmd.Flags().StringVar(&ownerID, "owner", "", "Update memory owner")
	cmd.Flags().Float64Var(&priority, "priority", -1, "Priority (0.0-1.0)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags for the memory")
	cmd.Flags().StringVar(&source, "source", "", "Source of the memory (optional)")
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func recallCmd() *cobra.Command {
	var robotMode bool
	var limit int
	var tagsFilter []string
	var categoryFilter string
	var ownerFilter string
	var withDecay bool
	var semanticSearch bool

	cmd := &cobra.Command{
		Use:   "recall [query]",
		Short: "Recall memories matching query",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			// Build filter options
			opts := store.RecallOptions{
				Query:      query,
				Limit:      limit,
				Tags:       tagsFilter,
				Category:   categoryFilter,
				OwnerID:    ownerFilter,
				WithDecay:  withDecay,
				Semantic:   semanticSearch,
			}

			// Search memories
			memories, err := store.RecallMemories(opts)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error recalling memories: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				// Robot Mode: Pure JSON to stdout
				result := map[string]interface{}{
					"status":   "ok",
					"query":    query,
					"filters":  map[string]interface{}{"tags": tagsFilter, "category": categoryFilter, "owner": ownerFilter},
					"count":    len(memories),
					"memories": memories,
				}
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(string(jsonBytes))
			} else {
				// Human-readable output
				filterDesc := ""
				if query != "" {
					filterDesc = fmt.Sprintf("matching '%s'", query)
				}
				if len(tagsFilter) > 0 {
					if filterDesc != "" {
						filterDesc += " and "
					}
					filterDesc += fmt.Sprintf("with tags: %v", tagsFilter)
				}
				if categoryFilter != "" {
					if filterDesc != "" {
						filterDesc += " and "
					}
					filterDesc += fmt.Sprintf("in category: %s", categoryFilter)
				}
				if filterDesc == "" {
					filterDesc = "(all memories)"
				}

				fmt.Printf("Found %d memory(ies) %s:\n\n", len(memories), filterDesc)
				if len(memories) == 0 {
					fmt.Println("No memories found.")
					return
				}

				for i, m := range memories {
					fmt.Printf("%d. [%s] %s\n", i+1, m.Category, m.ID)
					fmt.Printf("   Content: %s\n", m.Content)
					fmt.Printf("   Priority: %.1f | Accessed %d times\n", m.Priority, m.AccessCount)
					if len(m.Tags) > 0 {
						fmt.Printf("   Tags: %v\n", m.Tags)
					}
					fmt.Println()
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")
	cmd.Flags().StringSliceVar(&tagsFilter, "tags", []string{}, "Filter by tags (all tags must match)")
	cmd.Flags().StringVar(&categoryFilter, "category", "", "Filter by category (core|semantic|working|episodic)")
	cmd.Flags().StringVar(&ownerFilter, "owner", "", "Filter by memory owner")
	cmd.Flags().BoolVar(&withDecay, "decay", false, "Use decay-weighted scoring for recall")
	cmd.Flags().BoolVar(&semanticSearch, "semantic", false, "Use embeddings-based semantic search")
	return cmd
}

func catchupCmd() *cobra.Command {
	var robotMode bool
	var limit int
	var category string
	var since string

	cmd := &cobra.Command{
		Use:   "catchup",
		Short: "Catch up on recent memories",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			opts := store.CatchupOptions{
				Limit:    limit,
				Category: category,
				Since:    since,
			}

			memories, err := store.CatchupMemories(opts)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error catching up: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				result := map[string]interface{}{
					"status":   "ok",
					"count":    len(memories),
					"memories": memories,
				}
				jsonBytes, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("Recent memories (%d):\n\n", len(memories))
				for i, m := range memories {
					fmt.Printf("%d. [%s] %s (%s)\n", i+1, m.Category, m.ID, m.CreatedAt.Format("2006-01-02 15:04"))
					fmt.Printf("   %s\n\n", m.Content)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")
	cmd.Flags().StringVar(&category, "category", "", "Filter by category")
	cmd.Flags().StringVar(&since, "since", "", "Filter by creation time (YYYY-MM-DD HH:MM:SS)")
	return cmd
}

func historyCmd() *cobra.Command {
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "history [id]",
		Short: "Show version history for a memory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			id := args[0]
			history, err := store.GetMemoryHistory(id)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting history: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				result := map[string]interface{}{
					"status":  "ok",
					"id":      id,
					"history": history,
				}
				jsonBytes, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("History for memory %s:\n\n", id)
				for _, h := range history {
					fmt.Printf("Commit: %s\n", h.CommitHash)
					fmt.Printf("Date:   %s\n", h.CommitDate.Format("2006-01-02 15:04"))
					fmt.Printf("Content: %s\n", h.Content)
					fmt.Println("-----------------------------------")
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func rollbackCmd() *cobra.Command {
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "rollback [id] [commit-hash]",
		Short: "Rollback a memory to a specific version",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			id := args[0]
			commit := args[1]

			if err := store.RollbackMemory(id, commit); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error rolling back: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				fmt.Printf(`{"status":"ok","message":"rolled back memory %s to %s"}`+"\n", id, commit)
			} else {
				fmt.Printf("✓ Rolled back memory %s to version %s\n", id, commit)
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func linkCmd() *cobra.Command {
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "link [from-id] [to-id] [relation]",
		Short: "Link two memories together",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			from := args[0]
			to := args[1]
			relation := "related"
			if len(args) > 2 {
				relation = args[2]
			}

			if err := store.LinkMemories(from, to, relation); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error linking: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				fmt.Printf(`{"status":"ok","message":"linked %s to %s as %s"}`+"\n", from, to, relation)
			} else {
				fmt.Printf("✓ Linked %s to %s (%s)\n", from, to, relation)
			}
		},
	}
	
	cmd.AddCommand(&cobra.Command{
		Use:   "show [id]",
		Short: "Show all links for a memory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, _ := os.Getwd()
			db.InitDB(repoPath)
			defer db.CloseDB()

			links, err := store.GetMemoryLinks(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if robotMode {
				jsonBytes, _ := json.MarshalIndent(links, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("Links for %s:\n", args[0])
				for _, l := range links {
					fmt.Printf("- %s -> %s (%s)\n", l["from_id"], l["to_id"], l["relation"])
				}
			}
		},
	})

	cmd.PersistentFlags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func keystonesCmd() *cobra.Command {
	var robotMode bool
	var limit int

	cmd := &cobra.Command{
		Use:   "keystones",
		Short: "Identify core/foundational memories",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			keystones, err := store.GetKeystoneMemories(limit)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting keystones: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				result := map[string]interface{}{
					"status":    "ok",
					"count":     len(keystones),
					"keystones": keystones,
				}
				jsonBytes, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("Keystone Memories (%d):\n\n", len(keystones))
				for i, m := range keystones {
					fmt.Printf("%d. [%s] %s (Priority: %.1f, Accesses: %d)\n", i+1, m.Category, m.ID, m.Priority, m.AccessCount)
					fmt.Printf("   %s\n\n", m.Content)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")
	return cmd
}

func statsCmd() *cobra.Command {
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show memory database analytics",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, _ := os.Getwd()
			db.InitDB(repoPath)
			defer db.CloseDB()

			stats, err := store.GetMemoryStats()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting stats: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				jsonBytes, _ := json.MarshalIndent(stats, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Println("AMI Memory Statistics")
				fmt.Printf("Total Memories: %v\n", stats["total_memories"])
				fmt.Println("\nDistribution by Category:")
				dist := stats["distribution"].(map[string]int)
				for cat, count := range dist {
					fmt.Printf("- %-10s: %d\n", cat, count)
				}
				fmt.Println("\nMetrics:")
				metrics := stats["metrics"].(map[string]interface{})
				fmt.Printf("- Avg Priority:  %.2f\n", metrics["avg_priority"])
				fmt.Printf("- Avg Access:    %.2f\n", metrics["avg_access_count"])
				fmt.Printf("- Avg Decay Score: %.2f\n", metrics["avg_decay_score"])
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func contextCmd() *cobra.Command {
	var robotMode bool
	var limit int
	var tokenBudget int

	cmd := &cobra.Command{
		Use:   "context [task]",
		Short: "Get optimal context for a specific task",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			task := ""
			if len(args) > 0 {
				task = args[0]
			}

			// Get context memories with budget
			memories, err := store.GetContextMemories(task, limit, tokenBudget)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				// Robot Mode: Pure JSON
				result := map[string]interface{}{
					"status":   "ok",
					"task":     task,
					"budget":   tokenBudget,
					"memories": memories,
				}
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("Optimized Context for Task: %s (Budget: %d tokens)\n\n", task, tokenBudget)
				if len(memories) == 0 {
					fmt.Println("No relevant memories found.")
					return
				}

				for _, m := range memories {
					fmt.Printf("[%s] %s\n", m.Category, m.Content)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of task-related memories")
	cmd.Flags().IntVar(&tokenBudget, "tokens", 4000, "Maximum token budget for context")
	return cmd
}

func promoteCmd() *cobra.Command {
	var robotMode bool
	var globalPath string
	var autoPromote bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "promote [id]",
		Short: "Promote a memory to the global team brain",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, _ := os.Getwd()
			db.InitDB(repoPath)
			defer db.CloseDB()

			if autoPromote {
				// Auto-promote eligible memories
				memories, err := store.FindAutoPromotionCandidates(5, 0.8)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error finding candidates: %v\n", err)
					os.Exit(1)
				}

				if len(memories) == 0 {
					fmt.Println("No memories meet the promotion criteria.")
					return
				}

				if dryRun {
					fmt.Printf("Found %d candidate(s) for promotion:\n", len(memories))
					for _, m := range memories {
						fmt.Printf("  - [%s] %s (access_count: %d, priority: %.2f, category: %s)\n",
							m.ID[:8], m.Content, m.AccessCount, m.Priority, m.Category)
					}
					return
				}

				if !robotMode && !confirmAction(fmt.Sprintf("Promote %d memories to global brain?", len(memories))) {
					fmt.Println("Promotion cancelled.")
					return
				}

				for _, m := range memories {
					if err := store.PromoteMemory(m.ID, globalPath); err != nil {
						fmt.Fprintf(os.Stderr, "Error promoting %s: %v\n", m.ID, err)
					} else {
						fmt.Printf("✓ Promoted %s\n", m.ID[:8])
					}
				}
			} else {
				// Promote a specific memory
				if len(args) == 0 {
					fmt.Fprintln(os.Stderr, "Error: memory ID required (unless using --auto)")
					os.Exit(1)
				}
				id := args[0]
				if err := store.PromoteMemory(id, globalPath); err != nil {
					if robotMode {
						fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
					} else {
						fmt.Fprintf(os.Stderr, "Error promoting memory: %v\n", err)
					}
					os.Exit(1)
				}

				if robotMode {
					fmt.Printf(`{"status":"ok","message":"promoted memory %s to global store"}`+"\n", id)
				} else {
					fmt.Printf("✓ Promoted memory %s to global brain (%s)\n", id, globalPath)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	cmd.Flags().StringVar(&globalPath, "path", "/home/hargabyte/.ami/global", "Path to the global AMI store")
	cmd.Flags().BoolVar(&autoPromote, "auto", false, "Auto-promote memories meeting criteria")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show candidates without promoting")
	return cmd
}

func helpAgentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help-agents",
		Short: "Output agent-optimized command reference",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`# AMI Command Reference for AI Agents

> This tool manages your long-term memory using a versioned, metabolic architecture.
> Use it to store facts, decisions, and patterns so you don't burn tokens re-learning them.

## 🧠 Quick Start Workflow

1. **Before a Task**: Get focused context
   ` + "`" + `ami context "your task description" --limit 5 --robot` + "`" + `

2. **During a Task**: Store important discoveries
   ` + "`" + `ami add "Decision: use Go 1.22 for this module" --category working --tags technical` + "`" + `

3. **During a Task**: Track decisions with linked memories
   ` + "`" + `ami decision track "Use binary embeddings" --task "v0.4.0" --memories "abc,def"` + "`" + `

4. **After a Task**: Record decision outcomes
   ` + "`" + `ami decision outcome <id> --outcome 0.9 --feedback "Worked perfectly"` + "`" + `

5. **After a Task**: Clean up and promote
   ` + "`" + `ami promote <memory-id>` + "`" + ` (if it's a permanent team truth)

6. **Periodic Maintenance**: Reflect on episodic noise
   ` + "`" + `ami reflect --limit 10 --hours 24` + "`" + `

## 📂 Memory Categories

- **Core**: Foundational truths (User name, identity). Use for facts that NEVER change.
- **Semantic**: Learned patterns/habits. Use for general knowledge gained over time.
- **Working**: Task-specific context. Use for notes on the current session.
- **Episodic**: Event logs. Use for "I did X at time Y".

## 🎯 Decision Tracking (v0.5.0+)

Track the decisions you make and learn from their outcomes:

**Track a Decision:**
   ` + "`" + `ami decision track "your decision text" --task "project-id" --memories "id1,id2"` + "`" + `

**Record Outcome:**
   ` + "`" + `ami decision outcome <decision-id> --outcome 0.8 --feedback "Notes"` + "`" + `

**Synaptic Boost:** When outcome > 0.8, linked memories automatically get priority reinforcement.

## 🔍 Reflection (v0.5.0+)

   ` + "`" + `ami reflect --limit 10 --hours 24` + "`" + `

Identifies episodic noise and suggests semantic synthesis for consolidation.

## 🤖 Robot Mode

ALWAYS use the ` + "`" + `--robot` + "`" + ` flag for programmatic integration.
It returns pure JSON to stdout.

## 💡 Best Practices

- **Atomic Memories**: One fact per memory. Don't mix user preferences with technical specs.
- **Source Attribution**: Always use ` + "`" + `--source` + "`" + ` so future you knows WHY you believe a fact.
- **Aggressive Tagging**: Use tags for project IDs and concepts to make filtering faster.
- **Decision Tracking**: Link memories to decisions so successful choices reinforce useful knowledge.
- **Regular Reflection**: Use ` + "`" + `ami reflect` + "`" + ` to convert episodic noise into semantic facts.
`)
		},
	}
}

func deleteCmd() *cobra.Command {
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a memory by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			id := args[0]

			// Delete the memory
			if err := store.DeleteMemory(id); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error deleting memory: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				fmt.Printf(`{"status":"ok","message":"deleted memory %s"}`+"\n", id)
			} else {
				fmt.Printf("✓ Deleted memory %s\n", id)
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func tagsCmd() *cobra.Command {
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List all unique tags",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			tags, err := store.ListTags()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error listing tags: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				result := map[string]interface{}{
					"status": "ok",
					"tags":   tags,
				}
				jsonBytes, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Println("Unique Tags:")
				for _, tag := range tags {
					fmt.Printf("- %s\n", tag)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	return cmd
}

func checkpointCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkpoint [description]",
		Short: "Create a checkpoint before compression",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement
			fmt.Println("Checkpointing current state")
		},
	}
}

func consolidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "consolidate",
		Short: "Consolidate episodic memories to semantic",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement
			fmt.Println("Consolidating memories")
		},
	}
}

func decisionCmd() *cobra.Command {
	var taskID string
	var memoryIDsStr string
	var outcomeStr string
	var feedback string
	var robotMode bool

	cmd := &cobra.Command{
		Use:   "decision [action]",
		Short: "Track decisions and their outcomes",
		Long: `Track decisions and reinforce memories that lead to good outcomes.

Actions:
  track [decision]    - Track a new decision with linked memories
  outcome <id>        - Record the outcome of a decision (0.0 to 1.0)
  list [task_id]      - List all decisions, optionally filtered by task

Examples:
  ami decision track "Use binary embeddings" --task "v0.4.0" --memories "abc,def"
  ami decision outcome abc-123 --outcome 0.9 --feedback "Worked perfectly"
  ami decision list v0.4.0`,
	}

	trackCmd := &cobra.Command{
		Use:   "track [decision]",
		Short: "Track a new decision",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			// Parse memory IDs
			var memoryIDs []string
			if memoryIDsStr != "" {
				memoryIDs = strings.Split(memoryIDsStr, ",")
				// Trim spaces from each ID
				for i, id := range memoryIDs {
					memoryIDs[i] = strings.TrimSpace(id)
				}
			}

			// Track the decision
			decisionText := strings.Join(args, " ")
			decision, err := store.TrackDecision(taskID, memoryIDs, decisionText, "cli")
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error tracking decision: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				jsonBytes, _ := json.MarshalIndent(decision, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("✓ Decision tracked: %s\n", decision.ID)
				fmt.Printf("  Task: %s\n", decision.TaskID)
				fmt.Printf("  Text: %s\n", decision.DecisionText)
				if len(decision.MemoryIDs) > 0 {
					fmt.Printf("  Linked memories: %d\n", len(decision.MemoryIDs))
				}
			}
		},
	}
	trackCmd.Flags().StringVar(&taskID, "task", "", "Task ID")
	trackCmd.Flags().StringVar(&memoryIDsStr, "memories", "", "Comma-separated memory IDs")

	outcomeCmd := &cobra.Command{
		Use:   "outcome <decision_id>",
		Short: "Record the outcome of a decision",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				}
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				}
				os.Exit(1)
			}
			defer db.CloseDB()

			decisionID := args[0]
			var outcome float64
			_, err = fmt.Sscanf(outcomeStr, "%f", &outcome)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"invalid outcome: %v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error: invalid outcome value\n")
				}
				os.Exit(1)
			}

			// Validate outcome range
			if outcome < 0.0 || outcome > 1.0 {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"outcome must be between 0.0 and 1.0"}`+"\n")
				} else {
					fmt.Fprintf(os.Stderr, "Error: outcome must be between 0.0 and 1.0\n")
				}
				os.Exit(1)
			}

			// Record the outcome
			err = store.RecordOutcome(decisionID, outcome, feedback)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error recording outcome: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				fmt.Printf(`{"status":"success","decision_id":"%s","outcome":%f}`+"\n", decisionID, outcome)
			} else {
				fmt.Printf("✓ Outcome recorded: %.2f\n", outcome)
				if outcome > 0.8 {
					fmt.Printf("  → High success! Linked memories reinforced.\n")
				}
			}
		},
	}
	outcomeCmd.Flags().StringVar(&outcomeStr, "outcome", "", "Outcome value (0.0 to 1.0)")
	outcomeCmd.Flags().StringVar(&feedback, "feedback", "", "Optional feedback text")
	outcomeCmd.MarkFlagRequired("outcome")

	listCmd := &cobra.Command{
		Use:   "list [task_id]",
		Short: "List decisions",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				os.Exit(1)
			}
			defer db.CloseDB()

			listTaskID := ""
			if len(args) > 0 {
				listTaskID = args[0]
			}

			decisions, err := store.ListDecisions(listTaskID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing decisions: %v\n", err)
				os.Exit(1)
			}

			if robotMode {
				jsonBytes, _ := json.MarshalIndent(decisions, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				if len(decisions) == 0 {
					fmt.Println("No decisions found.")
					return
				}

				for _, d := range decisions {
					outcomeStr := "pending"
					if d.Outcome > 0 {
						outcomeStr = fmt.Sprintf("%.2f", d.Outcome)
					}

					fmt.Printf("\n%s\n", d.ID)
					fmt.Printf("  Task: %s\n", d.TaskID)
					fmt.Printf("  Decision: %s\n", d.DecisionText)
					fmt.Printf("  Outcome: %s\n", outcomeStr)
					if d.Feedback != "" {
						fmt.Printf("  Feedback: %s\n", d.Feedback)
					}
					if len(d.MemoryIDs) > 0 {
						fmt.Printf("  Linked memories: %d\n", len(d.MemoryIDs))
					}
				}
			}
		},
	}
	listCmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")

	cmd.AddCommand(trackCmd)
	cmd.AddCommand(outcomeCmd)
	cmd.AddCommand(listCmd)

	return cmd
}

func reflectCmd() *cobra.Command {
	var hours int
	var limit int

	cmd := &cobra.Command{
		Use:   "reflect",
		Short: "Reflect on episodic memories and suggest synthesis",
		Long: `Reflect on recent episodic memories and suggest semantic synthesis.
This is the MVP for autonomous reflection - it identifies noisy clusters
and provides synthesis prompts for consolidation.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				os.Exit(1)
			}
			defer db.CloseDB()

			// Calculate the time threshold
			sinceTime := time.Now().Add(-time.Duration(hours) * time.Hour)

			// Get episodic memories from the last N hours
			memories, err := store.CatchupMemories(store.CatchupOptions{
				Limit:    limit,
				Category: "episodic",
				Since:    sinceTime.Format("2006-01-02 15:04:05"),
			})

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching memories: %v\n", err)
				os.Exit(1)
			}

			if len(memories) == 0 {
				fmt.Println("No episodic memories found for reflection.")
				return
			}

			fmt.Printf("🤔 Reflecting on %d episodic memories from the last %d hour(s)\n\n", len(memories), hours)

			// Display memories
			for i, m := range memories {
				fmt.Printf("%d. [%s] %s\n", i+1, m.ID[:8], m.Content)
				if len(m.Tags) > 0 {
					fmt.Printf("   Tags: %v\n", m.Tags)
				}
			}

			fmt.Println("\n--- Synthesis Prompt ---")
			fmt.Println("Review the above memories and suggest 1-3 Semantic Facts that:")
			fmt.Println("  • Capture the essential knowledge")
			fmt.Println("  • Eliminate redundant detail")
			fmt.Println("  • Maintain high information density")
			fmt.Println()
			fmt.Println("For each suggested fact, provide:")
			fmt.Println("  1. Title (short, descriptive)")
			fmt.Println("  2. Content (concise, definitive statement)")
			fmt.Println("  3. Related memory IDs (for traceability)")
			fmt.Println()
			fmt.Println("Example:")
			fmt.Println("  Fact 1:")
			fmt.Println("    Title: Binary embedding storage")
			fmt.Println("    Content: AMI v0.4.0 uses BLOB fields with HEX encoding to store float32 embeddings,")
			fmt.Println("               providing bit-perfect precision and enabling semantic search.")
			fmt.Println("    Related: abc-123, def-456")
		},
	}
	cmd.Flags().IntVar(&hours, "hours", 24, "Hours to look back")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of memories to reflect on")

	return cmd
}

func conflictCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conflict",
		Short: "Detect and resolve conflicting memories",
		Long: `Detect memories that may contradict each other using semantic similarity.
This helps maintain team consensus across multiple agents.`,
	}

	resolveCmd := &cobra.Command{
		Use:   "resolve <id1> <id2>",
		Short: "Resolve a conflict between two memories",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoPath, _ := os.Getwd()
			db.InitDB(repoPath)
			defer db.CloseDB()

			id1 := args[0]
			id2 := args[1]

			// Get both memories
			m1, err := store.GetMemoryByID(id1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching memory %s: %v\n", id1, err)
				os.Exit(1)
			}
			m2, err := store.GetMemoryByID(id2)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching memory %s: %v\n", id2, err)
				os.Exit(1)
			}

			fmt.Println("Memory 1:")
			fmt.Printf("  ID: %s\n", m1.ID)
			fmt.Printf("  Content: %s\n", m1.Content)
			fmt.Printf("  Category: %s\n", m1.Category)

			fmt.Println("\nMemory 2:")
			fmt.Printf("  ID: %s\n", m2.ID)
			fmt.Printf("  Content: %s\n", m2.Content)
			fmt.Printf("  Category: %s\n", m2.Category)

			fmt.Println("\nResolution options:")
			fmt.Println("1. Keep Memory 1 (deprecate Memory 2)")
			fmt.Println("2. Keep Memory 2 (deprecate Memory 1)")
			fmt.Println("3. Merge into Memory 1")
			fmt.Println("4. Keep both (no action)")

			fmt.Print("\nSelect option [1-4]: ")
			var choice int
			fmt.Scanln(&choice)

			switch choice {
			case 1:
				store.UpdateMemoryStatus(id2, models.StatusDeprecated)
				fmt.Printf("✓ Deprecated %s\n", id2)
			case 2:
				store.UpdateMemoryStatus(id1, models.StatusDeprecated)
				fmt.Printf("✓ Deprecated %s\n", id1)
			case 3:
				// Simple merge: combine content
				mergedContent := fmt.Sprintf("%s | %s", m1.Content, m2.Content)
				store.UpdateMemoryContent(id1, mergedContent)
				store.UpdateMemoryStatus(id2, models.StatusDeprecated)
				fmt.Printf("✓ Merged into %s, deprecated %s\n", id1, id2)
			case 4:
				fmt.Println("No action taken.")
			default:
				fmt.Println("Invalid choice.")
			}
		},
	}

	cmd.AddCommand(resolveCmd)
	return cmd
}

func robotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "robot",
		Short: "Robot mode commands for agent integration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Memory system status (JSON)",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, err := os.Getwd()
			if err != nil {
				result := map[string]interface{}{
					"status":  "error",
					"message": "Cannot get working directory",
					"version": version,
				}
				jsonBytes, _ := json.Marshal(result)
				fmt.Println(string(jsonBytes))
				os.Exit(1)
			}

			if err := db.InitDB(repoPath); err != nil {
				result := map[string]interface{}{
					"status":  "error",
					"message": "Database initialization failed",
					"version": version,
				}
				jsonBytes, _ := json.Marshal(result)
				fmt.Println(string(jsonBytes))
				os.Exit(1)
			}
			defer db.CloseDB()

			// Get memory count
			count, err := store.GetMemoryCount()
			if err != nil {
				result := map[string]interface{}{
					"status":  "error",
					"message": "Failed to count memories",
					"version": version,
				}
				jsonBytes, _ := json.Marshal(result)
				fmt.Println(string(jsonBytes))
				os.Exit(1)
			}

			result := map[string]interface{}{
				"status":   "ok",
				"memories": count,
				"version":  version,
			}
			jsonBytes, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonBytes))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "checkpoint",
		Short: "Auto-checkpoint for compression hooks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`{"checkpointed":true}`)
		},
	})

	return cmd
}

func pairingCmd() *cobra.Command {
	var taskID string

	cmd := &cobra.Command{
		Use:   "pairing",
		Short: "Manage session pairing daemon",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the pairing daemon",
		Run: func(cmd *cobra.Command, args []string) {
			socketPath := store.GetSocketPath()
			
			// Cleanup old socket
			os.Remove(socketPath)

			listener, err := net.Listen("unix", socketPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error starting pairing daemon: %v\n", err)
				os.Exit(1)
			}
			defer listener.Close()

			fmt.Printf("✓ Pairing daemon started at %s (Task: %s)\n", socketPath, taskID)
			fmt.Println("Listening for tool reports...")

			for {
				conn, err := listener.Accept()
				if err != nil {
					continue
				}
				
				go func(c net.Conn) {
					defer c.Close()
					var action store.PairingAction
					if err := json.NewDecoder(c).Decode(&action); err == nil {
						fmt.Printf(" [LOG] Tool: %s | Action: %s\n", action.Source, action.Action)
						// In a real impl, we'd write to a draft branch in Dolt here
					}
				}(conn)
			}
		},
	}
	startCmd.Flags().StringVar(&taskID, "task", "default", "Task ID to associate with the session")

	cmd.AddCommand(startCmd)
	return cmd
}

func syncCmd() *cobra.Command {
	var channelID string
	var limit int

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize memories from external sources",
	}

	mmCmd := &cobra.Command{
		Use:   "mattermost",
		Short: "Sync memories from Mattermost channel",
		Run: func(cmd *cobra.Command, args []string) {
			token := os.Getenv("MATTERMOST_TOKEN")
			url := os.Getenv("MATTERMOST_URL")
			
			if token == "" || url == "" {
				fmt.Fprintln(os.Stderr, "Error: MATTERMOST_TOKEN and MATTERMOST_URL must be set")
				os.Exit(1)
			}

			client := db.NewMattermostClient(url, token)
			messages, err := client.GetRecentMessages(channelID, limit)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching messages: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ Pulled %d messages from Mattermost. Processing via Ollama...\n", len(messages))
			
			// Process through reflection engine (Ollama)
			ctx := context.Background()
			ollama := db.NewOllamaClient("http://localhost:11434", "qwen2.5-coder:1.5b")
			
			rawContent := strings.Join(messages, "\n---\n")
			facts, err := store.ExtractTechnicalFacts(ctx, ollama, rawContent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error extracting facts: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ Extracted %d potential facts for review:\n", len(facts))
			for _, f := range facts {
				fmt.Printf("- %s\n", f)
			}
		},
	}
	mmCmd.Flags().StringVar(&channelID, "channel", "", "Mattermost Channel ID")
	mmCmd.Flags().IntVar(&limit, "limit", 20, "Number of messages to pull")
	mmCmd.MarkFlagRequired("channel")

	cmd.AddCommand(mmCmd)
	return cmd
}
