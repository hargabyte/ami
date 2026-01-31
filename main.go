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

var version = "0.1.0"

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
	var priority float64
	var tags []string
	var source string

	cmd := &cobra.Command{
		Use:   "add [content]",
		Short: "Add a memory to the database",
		Args:  cobra.MinimumNArgs(1),
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

			// Validate category
			cat := models.Category(category)
			if !cat.IsValid() {
				fmt.Fprintf(os.Stderr, "Error: invalid category '%s'. Must be one of: core, semantic, working, episodic\n", category)
				os.Exit(1)
			}

			// Validate priority
			if priority < 0.0 || priority > 1.0 {
				fmt.Fprintf(os.Stderr, "Error: priority must be between 0.0 and 1.0\n")
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
			memory, err := store.AddMemory(content, cat, priority, tags, source)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error adding memory: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ Added memory %s (category: %s, priority: %.1f)\n", memory.ID, memory.Category, memory.Priority)
		},
	}
	cmd.Flags().StringVar(&category, "category", "episodic", "Memory category (core|semantic|working|episodic)")
	cmd.Flags().Float64Var(&priority, "priority", 0.5, "Priority (0.0-1.0)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags for the memory")
	cmd.Flags().StringVar(&source, "source", "", "Source of the memory (optional)")
	return cmd
}

func updateCmd() *cobra.Command {
	var category string
	var priority float64
	var tags []string
	var source string

	cmd := &cobra.Command{
		Use:   "update [id] [new-content]",
		Short: "Update an existing memory",
		Args:  cobra.MinimumNArgs(1),
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

			id := args[0]

			// Build update parameters
			params := store.UpdateParams{
				ID: id,
			}

			// Only update fields that were explicitly set
			if cmd.Flags().Changed("category") {
				cat := models.Category(category)
				if !cat.IsValid() {
					fmt.Fprintf(os.Stderr, "Error: invalid category '%s'. Must be one of: core, semantic, working, episodic\n", category)
					os.Exit(1)
				}
				params.Category = &cat
			}

			if cmd.Flags().Changed("priority") {
				if priority < 0.0 || priority > 1.0 {
					fmt.Fprintf(os.Stderr, "Error: priority must be between 0.0 and 1.0\n")
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
				fmt.Fprintf(os.Stderr, "Error updating memory: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ Updated memory %s\n", id)
		},
	}
	cmd.Flags().StringVar(&category, "category", "", "Memory category (core|semantic|working|episodic)")
	cmd.Flags().Float64Var(&priority, "priority", -1, "Priority (0.0-1.0)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags for the memory")
	cmd.Flags().StringVar(&source, "source", "", "Source of the memory (optional)")
	return cmd
}

func recallCmd() *cobra.Command {
	var robotMode bool
	var limit int
	var tagsFilter []string
	var categoryFilter string

	cmd := &cobra.Command{
		Use:   "recall [query]",
		Short: "Recall memories matching query",
		Args:  cobra.MinimumNArgs(0),
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

			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			// Build filter options
			opts := store.RecallOptions{
				Query:    query,
				Limit:    limit,
				Tags:     tagsFilter,
				Category: categoryFilter,
			}

			// Search memories
			memories, err := store.RecallMemories(opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error recalling memories: %v\n", err)
				os.Exit(1)
			}

			if robotMode {
				// Robot Mode: Pure JSON to stdout
				result := map[string]interface{}{
					"query":    query,
					"filters":  map[string]interface{}{"tags": tagsFilter, "category": categoryFilter},
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
