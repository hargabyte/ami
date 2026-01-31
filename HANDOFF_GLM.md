# Development Handoff: AMI Project (Go + DoltDB)

## Current Status
- Skeleton created in `/home/hargabyte/ami`
- `main.go` has the CLI structure (Cobra)
- `schema.sql` defines the initial memories/links table
- DoltDB is the chosen versioned database backend

## Phase 1: Foundation (GLM Tasks)

### 1. Initialize the Store
```bash
cd /home/hargabyte/ami
dolt init
dolt sql < schema.sql
dolt add .
dolt commit -m "Initialize AMI schema"
```

### 2. Database Connection (`internal/db/dolt.go`)
- Implement a connection pool to the Dolt database.
- Use `github.com/go-sql-driver/mysql` since Dolt is MySQL-compatible.
- Note: For local CLI use, you may need to run `dolt sql-server` or use the library to execute queries against the local repo.

### 3. Implementation of `ami add`
- Generate a UUID for each memory.
- Category must be one of: `core`, `semantic`, `working`, `episodic`.
- After every `INSERT`, perform a `dolt_commit` via SQL to ensure versioning history is captured.
- Example SQL: `SELECT dolt_commit('-am', 'Add memory: <excerpt>')`

### 4. Basic `ami recall`
- Implement a basic text-search retrieval.
- Output format must follow the **Robot Mode Pattern**:
  - `stdout`: Pure JSON
  - `stderr`: Logs/Diagnostics
  - `exit 0`: Success

## Design Patterns to Copy
Look at **Beads Viewer** for how they handled "Robot Commands". We want `ami` to be a "sidecar" for agents to manage their own context.

## Support
- Research: @hsa-gemini is synced into the team for any docs or API research needed.
- Architecture: @simon (Chief of Staff) for final sign-off on schema changes.

Let's ship the foundation tonight. 🚀
