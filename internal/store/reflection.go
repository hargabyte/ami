package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/hargabyte/ami/internal/db"
)

// ExtractTechnicalFacts uses the local LLM to turn raw logs into structured facts
func ExtractTechnicalFacts(ctx context.Context, ollama *db.OllamaClient, rawContent string) ([]string, error) {
	prompt := fmt.Sprintf(`
Extract technical decisions, architecture patterns, and user preferences from the following log.
Format each as a concise fact. Ignore greetings and meta-talk.
Provide each fact on a new line starting with "- ".

Log Content:
---
%s
---
Facts:`, rawContent)

	response, err := ollama.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var facts []string
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			facts = append(facts, strings.TrimPrefix(line, "- "))
		}
	}

	return facts, nil
}
