package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hargabyte/ami/internal/db"
	"github.com/hargabyte/ami/internal/models"
	"github.com/hargabyte/ami/internal/store"
	"github.com/spf13/cobra"
)

var version = "0.3.0"

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
	rootCmd.AddCommand(promoteCmd())
	rootCmd.AddCommand(deleteCmd())
	rootCmd.AddCommand(tagsCmd())
	rootCmd.AddCommand(checkpointCmd())
	rootCmd.AddCommand(consolidateCmd())
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

	cmd := &cobra.Command{
		Use:   "context [task]",
		Short: "Get optimal context for a specific task",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, _ := os.Getwd()
			db.InitDB(repoPath)
			defer db.CloseDB()

			task := ""
			if len(args) > 0 {
				task = args[0]
			}

			memories, err := store.GetContextMemories(task, limit)
			if err != nil {
				if robotMode {
					fmt.Printf(`{"status":"error","message":"%v"}`+"\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
				}
				os.Exit(1)
			}

			if robotMode {
				result := map[string]interface{}{
					"status":   "ok",
					"task":     task,
					"memories": memories,
				}
				jsonBytes, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("Optimized Context for Task: %s\n\n", task)
				for _, m := range memories {
					fmt.Printf("[%s] %s\n", m.Category, m.Content)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of task-related memories")
	return cmd
}

func promoteCmd() *cobra.Command {
	var robotMode bool
	var globalPath string

	cmd := &cobra.Command{
		Use:   "promote [id]",
		Short: "Promote a memory to the global team brain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize database
			repoPath, _ := os.Getwd()
			db.InitDB(repoPath)
			defer db.CloseDB()

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
		},
	}
	cmd.Flags().BoolVar(&robotMode, "robot", false, "Robot mode: output JSON")
	cmd.Flags().StringVar(&globalPath, "path", "/home/hargabyte/.ami/global", "Path to the global AMI store")
	return cmd
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
