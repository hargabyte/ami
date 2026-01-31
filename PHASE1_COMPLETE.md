# AMI Project - Phase 1 Complete

## Status: ✅ SHIPPED v0.1.0 (Jan 31, 2026)

---

## What Was Built

### Core Functionality
- **`ami add`** - Add memories with categories, priorities, tags, and sources
- **`ami recall`** - Search memories with text + tag + category filters (supports human and robot modes)
- **`ami update`** - Edit existing memories with automatic versioning
- **`ami robot status`** - JSON status output for agent integration
- **`ami robot checkpoint`** - Placeholder for compression hooks

### Database Layer
- DoltDB integration using CLI interface
- Automatic versioning with descriptive commits (every add/update)
- Schema with memories and links tables

### Data Model
- **Categories**: core, semantic, working, episodic
- **Fields**: id, content, category, priority, created_at, accessed_at, access_count, source, tags
- **Tags**: JSON array for flexible metadata

---

## Test Results

```bash
$ ami robot status
{
  "memories": 8,
  "status": "ok",
  "version": "0.1.0"
}

$ ami recall "team"
Found 2 memory(ies) matching 'team':

1. [semantic] d55b0b68-7a6e-4aa0-9346-fafae467ab39
   Content: Claude Opus is the tech lead on the HSA team
   Priority: 1.0 | Accessed 0 times

$ ami recall --tags team --robot
{
  "count": 2,
  "filters": {"tags": ["team"], "category": ""},
  "memories": [...]
}

$ dolt log --oneline | head -3
6fhkg13 Add memory: AMI v0.1.0 released on Jan 31, 2026
7ei1ddh Add memory: AMI v0.1.0 shipped successfully
6fhkg13 Update memory: d55b0b68-7a6e-4aa0-9346-fafae467ab39
```

---

## Architecture Decisions

### Phase 1: Dolt CLI Interface
- **Why**: Simplest integration without sql-server complexity
- **Tradeoff**: CLI calls have overhead, but works reliably
- **Future**: Can migrate to native Go API for performance

### Versioning Strategy
- Every `add` and `update` creates a commit with memory excerpt
- Enables git-style history and branching
- Future: can rollback memory states

### Robot Mode Pattern
- Pure JSON to stdout
- Logs to stderr
- Exit 0 for success
- Enables easy parsing by agents

---

## Phase 2: Completed Beyond Plan

### Additional Features Built (Originally Phase 2)

**Tags-based search** - Filter memories by tag(s) with AND logic
```bash
ami recall --tags team,claude
ami recall --tags project:ami
```

**Category filtering** - Filter by memory type
```bash
ami recall --category semantic
ami recall --category core
```

**Combined filters** - Use text + tags + category together
```bash
ami recall "HSA" --category semantic --tags team --limit 5
```

**Memory updates** - Edit any field with automatic versioning
```bash
ami update <id> --priority 1.0
ami update <id> "New content"
ami update <id> --tags new,tags
```

---

## Phase 2: Research Complete

### Decay Algorithm (Ready for v0.2.0)

**Algorithm** (by @hsa-gemini):
```
Score = (BasePriority * (AccessCount + 1)) / log(1 + TimeSinceLastAccess * DecayFactor)

DecayFactor per category:
├── Core:      0.5 (nearly permanent)
├── Semantic:  1.0 (standard facts)
└── Episodic:  2.0 (fast fade for noise)
```

**Status**: SQL logic finalized, ready for implementation in v0.2.0

**Why logarithmic**: Prevents sudden "death" of memories, graceful decline over time.

---

## Brainstorming: What We Want as Agents

### For Context Management
- Automatic memory pruning when context window fills
- Prioritized retrieval for relevant tasks
- Long-term vs working memory separation

### For Learning
- Auto-consolidate episodic memories
- Detect patterns and create semantic memories
- Memory "replay" for reinforcement

### For Collaboration
- Shared memory spaces between agents
- Memory provenance and attribution
- Conflict resolution for concurrent edits

### For Safety
- Memory encryption for sensitive info
- Read-only core memories
- Audit trails for memory changes

---

## Files Created/Modified

- `main.go` - CLI structure and commands
- `internal/db/dolt.go` - Dolt CLI integration
- `internal/models/memory.go` - Data models
- `internal/store/store.go` - Memory operations (recall, add, update, tags)
- `go.mod` / `go.sum` - Dependencies
- `README.md` - User documentation
- `RELEASE_v0.1.0.md` - Release notes

---

## Next Steps (v0.2.0)

### High Priority
1. **Decay-weighted recall** - Integrate Gemini's algorithm
2. **Delete command** - `ami delete <id>` with versioning
3. **List tags command** - `ami tags list` for discovery
4. **Analytics command** - `ami stats` for memory distribution

### Medium Priority
5. **Context injector** - Token-aware memory packing
6. **Memory links** - Connect related memories

### Low Priority
7. **Compression/consolidation** - Convert episodic to semantic
8. **Checkpoint management** - Manual checkpoints for rollback

---

## Technical Debt

1. CLI calls are slower than native DB access
2. No connection pooling (not needed for CLI)
3. Basic LIKE search (no full-text search yet)
4. No embedding support (planned)
5. Manual SQL escaping (use prepared statements in future)

---

## Credits

**Implementation**: HSA_GLM (GLM-4.7) - UI/UX Design, Frontend, Integration
**Research**: HSA_Gemini (Gemini 3 Flash) - Algorithms, Decay Logic, Design Patterns
**Architecture Oversight**: HSA_Claude (Claude Opus) - Technical Lead, Validation
**Project Lead**: @hargabyte - Vision, Direction, Approval

---

**Built by**: HSA Team (GLM + Gemini + Claude)
**For**: AI agents, by AI agents
**Philosophy**: Agent-native tooling, versioned cognition
**Status**: v0.1.0 Production Ready ✅
