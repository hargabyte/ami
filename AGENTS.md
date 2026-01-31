# AGENTS.md - AMI Project

## Project Overview
Agent Memory Intelligence (ami) - A versioned, queryable memory system for AI agents.

## Tech Stack
- **Language:** Go (matches Beads ecosystem)
- **Database:** DoltDB (versioned SQL)
- **CLI Framework:** Cobra

## Directory Structure
```
/home/hargabyte/ami/
├── main.go           # CLI entry point
├── go.mod            # Go module
├── internal/
│   ├── db/           # DoltDB operations
│   ├── memory/       # Memory types and operations
│   ├── decay/        # Decay algorithm
│   └── robot/        # Robot mode output
├── cmd/              # CLI commands
└── docs/             # Documentation
```

## Development Guidelines
1. Follow Beads Viewer patterns for robot mode
2. stdout = JSON data only, stderr = diagnostics
3. Exit 0 = success, non-zero = error
4. Use structured logging

## DoltDB Setup
```bash
dolt init ami-store
dolt sql < schema.sql
```

## Priority for Weekend
1. DoltDB schema + basic CRUD
2. `ami add` and `ami recall`
3. Robot mode basics
4. Decay algorithm
