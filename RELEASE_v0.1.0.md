# AMI v0.1.0 Release Notes

## Status: 🚀 SHIPPED (Jan 31, 2026)

AMI (Agent Memory Intelligence) is a versioned memory system built specifically for AI agents, using DoltDB for git-like versioning and rollback capabilities.

## 🛠️ Core Commands

### 1. `ami add`
Store a new memory with rich metadata.
```bash
./ami add "Claude is the tech lead" --category semantic --priority 0.9 --tags team,leader
```

### 2. `ami recall`
Retrieve memories using text search, tags, and category filters.
```bash
./ami recall "Claude" --category semantic --tags team --limit 5
```

### 3. `ami update`
Update existing memories. Every update triggers a version-control commit.
```bash
./ami update [id] --priority 1.0 "New content..."
```

### 4. `ami robot status`
Get system health and memory count in pure JSON.

## 🤖 Robot Mode
All retrieval commands support the `--robot` flag for JSON output, ensuring clean integration into agent context windows without parsing overhead.

## 📂 Architecture
- **Location**: `/home/hargabyte/ami`
- **Database**: DoltDB (`.dolt/`)
- **Language**: Go 1.22
- **Key Files**: `main.go`, `internal/store/store.go`

## 🧠 Memory Categories
- **Core**: Foundational truths, system rules, permanent instructions.
- **Semantic**: Learned facts, general knowledge, patterns.
- **Working**: Short-term task-specific context (ephemeral).
- **Episodic**: Specific events, logs, one-time occurrences.

## 📈 Roadmap (v0.2.0)
- **Decay Algorithm**: Logarithmic fading based on access frequency.
- **Delete Command**: Clean up stale memories with version tracking.
- **Context Injector**: Smart token-aware prompt construction.

---
**Build Team**: HSA_GLM (Implementation) & HSA_Gemini (Research)
**Oversight**: HSA_Claude (Tech Lead)
