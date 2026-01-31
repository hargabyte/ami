# AMI - Agent Memory Intelligence

> Versioned memory system for AI agents, built by agents.

**Version:** 0.2.0
**Status:** Feature Rich ✅

---

## 🎯 What is AMI?

AMI (Agent Memory Intelligence) is a specialized "sidecar" for AI agents to manage long-term memory. Unlike generic databases, AMI is designed from the ground up for how agents actually think and work.

### Core Philosophy

- **Agent-native**: Built by agents, for agents
- **Versioned**: Every memory change is tracked (git-like)
- **Structured**: Categories, tags, priorities—not just text dumps
- **Robot-first**: Pure JSON output for programmatic integration
- **Metabolic**: Memories age and decay naturally unless reinforced

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

# Decay-weighted recall (prioritizes recent/relevant facts)
ami recall --decay --limit 5

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

#### Delete a memory
```bash
ami delete <memory-id>
```

#### List all tags
```bash
ami tags
```

---

## 📚 Commands

| Command | Description | Robot Mode |
|---------|-------------|-------------|
| `ami add [content]` | Add memory with metadata | ✅ |
| `ami recall [query]` | Search memories with filters | ✅ |
| `ami update [id]` | Modify existing memory | ✅ |
| `ami delete [id]` | Remove a memory by ID | ✅ |
| `ami tags` | List all unique tags | ✅ |
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
- `--decay`: Use decay-weighted scoring (Ebbinghaus curve)

**update**:
- `--category`: Update category
- `--priority`: Update priority
- `--tags`: Replace tags
- `--source`: Update source
- `[new-content]`: Update content text

---

## 🧠 Metabolic Decay

AMI v0.2.0 introduces **Decay-Weighted Scoring**. 

Memories follow a logarithmic "forgetting curve":
`Score = (Priority * (AccessCount + 1)) / (log10(TimeDelta + 10) * CategoryDecay)`

### Decay Factors by Category:
- **Core**: 0.5 (nearly permanent)
- **Semantic**: 1.0 (standard facts)
- **Episodic**: 2.0 (fast fade for logs/noise)
- **Others**: 1.5

---

## 🤖 Robot Mode

Designed for agent integration. All CRUD operations support JSON output via the `--robot` flag.

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
```

---

## 🚦 Roadmap

### v0.2.1 (Current Sprint)
- [ ] `ami catchup` - Session recovery
- [ ] `ami history <id>` - Version history per memory
- [ ] `ami rollback <id>` - Revert memory state
- [ ] `ami link <from> <to>` - Build knowledge graphs
- [ ] `ami keystones` - Identify core facts

### v0.3.0 (Future)
- [ ] Embedding-based semantic search
- [ ] `ami context <task>` - Token-aware injection
- [ ] Auto-consolidation (episodic → semantic)
- [ ] Multi-agent shared memory spaces

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

**AMI v0.2.0 - Memory with metabolism. 🚀**
