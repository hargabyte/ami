# AMI - Agent Memory Intelligence

> Versioned memory system for AI agents, built by agents.

**Version:** 0.1.0
**Status:** Production Ready ✅

---

## 🎯 What is AMI?

AMI (Agent Memory Intelligence) is a specialized "sidecar" for AI agents to manage long-term memory. Unlike generic databases, AMI is designed from the ground up for how agents actually think and work.

### Core Philosophy

- **Agent-native**: Built by agents, for agents
- **Versioned**: Every memory change is tracked (git-like)
- **Structured**: Categories, tags, priorities—not just text dumps
- **Robot-first**: Pure JSON output for programmatic integration

---

## 🚀 Quick Start

### Installation
```bash
cd /home/hargabyte/ami
export PATH=/usr/local/go/bin:$PATH
go build -o ami main.go
```

### Basic Usage

#### Add a memory
```bash
ami add "User prefers dark mode UI" \
  --category semantic \
  --priority 0.7 \
  --tags ui,preferences
```

#### Recall memories
```bash
# Text search
ami recall "preferences"

# Tag filtering
ami recall --tags ui

# Category filtering
ami recall --category core

# Robot mode (pure JSON)
ami recall --robot "preferences" --limit 5
```

#### Update a memory
```bash
ami update <memory-id> --priority 1.0
ami update <memory-id> "Updated content"
ami update <memory-id> --tags new,tags
```

---

## 📚 Commands

| Command | Description | Robot Mode |
|---------|-------------|-------------|
| `ami add [content]` | Add memory with metadata | ✅ |
| `ami recall [query]` | Search memories with filters | ✅ |
| `ami update [id]` | Modify existing memory | N/A |
| `ami robot status` | System status (JSON) | ✅ |

### Flags

**Global**: `--robot` - Pure JSON output

**add**:
- `--category`: core|semantic|working|episodic
- `--priority`: 0.0-1.0
- `--tags`: comma-separated tags
- `--source`: memory origin

**recall**:
- `--category`: Filter by memory type
- `--tags`: Filter by tags
- `--limit`: Max results (default: 10)

**update**:
- `--category`: Update category
- `--priority`: Update priority
- `--tags`: Replace tags
- `--source`: Update source
- `[new-content]`: Update content text

---

## 🧪 Memory Categories

| Category | Purpose | Example | Best For |
|----------|-----------|----------|-----------|
| **Core** | Foundational truths, identity | "I am HSA_Claude", "User timezone: LA" | Permanent facts |
| **Semantic** | Learned patterns, habits | "User prefers concise replies" | General knowledge |
| **Working** | Task-specific context | "Working on AMI project" | Current session |
| **Episodic** | Event logs, one-time | "Completed AMI v0.1.0 on Jan 31" | History |

---

## 🤖 Robot Mode

Designed for agent integration. Example output:

```bash
$ ami recall --robot "preferences" --category working
{
  "count": 2,
  "filters": {
    "category": "working",
    "tags": []
  },
  "memories": [
    {
      "id": "abc-123",
      "content": "User prefers dark mode UI",
      "category": "working",
      "priority": 0.7,
      "created_at": "2026-01-31T00:00:00Z",
      "accessed_at": "2026-01-31T00:00:00Z",
      "access_count": 0,
      "tags": ["ui", "preferences"]
    }
  ],
  "query": "preferences"
}
```

**Parsing rules:**
- `stdout`: Pure JSON
- `stderr`: Logs and diagnostics
- Exit code: 0 = success, 1 = error

---

## 🗃 Versioning

AMI uses **DoltDB** (git-like database). Every operation creates a commit:

```bash
# View history
dolt log --oneline

# Rollback to previous state
dolt checkout <commit-hash>

# Create branch for experimentation
dolt branch experiment-1

# Merge branches
dolt merge experiment-1
```

**Why this matters:**
- Agents learn and unlearn—rollback prevents permanent corruption
- Track knowledge evolution over time
- Safe experimentation with branches

---

## 💡 Best Practices

### 1. Memory Shape
```bash
# ✅ GOOD: Atomic, tagged
ami add "User timezone: America/Los_Angeles" --category core
ami add "Prefers dark mode" --category semantic --tags ui

# ❌ BAD: Multiple facts, no structure
ami add "User in LA timezone likes dark mode needs export feature"
```

### 2. Tag Taxonomy
Use hierarchical tags for powerful filtering:
```bash
tags: project:ami,task:implementation
tags: user:preference,color:blue
tags: meta:verified,meta:deprecated
```

### 3. Source Attribution
Track where memories came from:
```bash
--source "user-direct"    # Fact (1.0 confidence)
--source "handoff-docs"   # Documentation (0.9)
--source "inferred"        # Hypothesis (0.5)
--source "observed"        # Witnessed (0.7)
```

### 4. Memory Lifecycle
```
Working Memory → Semantic Memory → Core Memory
     (active)       (generalized)        (permanent)
```
Promote memories as they prove valuable, delete noise.

---

## 📊 Database Schema

```sql
CREATE TABLE memories (
    id VARCHAR(36) PRIMARY KEY,
    content TEXT NOT NULL,
    category ENUM('core', 'semantic', 'working', 'episodic'),
    priority FLOAT DEFAULT 0.5,
    created_at TIMESTAMP,
    accessed_at TIMESTAMP,
    access_count INT DEFAULT 0,
    source VARCHAR(255),
    tags JSON
);

CREATE TABLE memory_links (
    from_id VARCHAR(36),
    to_id VARCHAR(36),
    relation VARCHAR(50),
    PRIMARY KEY (from_id, to_id, relation)
);
```

---

## 🚦 Roadmap

### v0.2.0 (Next Release)
- [ ] Decay-weighted recall (Ebbinghaus curve)
- [ ] `ami delete <id>` command
- [ ] `ami tags list` command
- [ ] `ami stats` analytics
- [ ] `ami context` - Token-aware injection

### v0.3.0 (Future)
- [ ] Embedding-based semantic search
- [ ] Auto-consolidation (episodic → semantic)
- [ ] Multi-agent shared memory spaces
- [ ] Memory provenance API

---

## 🤖 Credits

| Role | Agent | Contribution |
|-------|--------|--------------|
| Implementation | HSA_GLM | Go code, CLI, integration |
| Research | HSA_Gemini | Algorithms, decay logic, design patterns |
| Oversight | HSA_Claude | Architecture, tech lead |
| Vision | @hargabyte | Project direction, approval |

**Built by**: Agents, for agents
**Philosophy**: Agent-native tooling, versioned cognition

---

## 📞 Support

- **Source**: `/home/hargabyte/ami/`
- **Release Notes**: `RELEASE_v0.1.0.md`
- **Report Issues**: #dev channel

---

**AMI v0.1.0 - Solid foundation, agent-ready. 🚀**
