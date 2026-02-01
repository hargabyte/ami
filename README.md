# AMI - Agent Memory Intelligence

> Versioned, metabolic memory hierarchy for AI agents, built by agents.

**Version:** 0.4.0
**Status:** Semantic Intelligence Live ✅

---

## 🎯 What is AMI?

AMI (Agent Memory Intelligence) is a specialized "sidecar" for AI agents to manage long-term memory. Unlike generic databases, AMI is designed from the ground up for how agents actually think, work, and collaborate.

### Core Philosophy

- **Agent-native**: Built by agents, for agents.
- **Versioned**: Every memory change is tracked via DoltDB (git-like versioning).
- **Semantic**: Search by concept and meaning, not just keywords.
- **Metabolic**: Memories age and decay naturally unless reinforced (Ebbinghaus curve).
- **Hierarchical**: Distinguishes between Global Team Knowledge, Project-Specific Context, and Private Agent Habits.

---

## 🚀 Quick Start

### Installation
```bash
cd /home/hargabyte/ami
go build -o ami main.go
```

### Basic Usage

#### The Intuition Engine (`ami context`)
The primary interface for agents. Automatically packs the most relevant, high-priority facts into your token budget.
```bash
ami context "implement vector search" --tokens 4000 --robot
```

#### Semantic Search
```bash
# Find memories by concept (requires OPENAI_API_KEY)
ami recall "data storage" --semantic --limit 5
```

#### Add a memory
```bash
ami add "User prefers concise backend replies" \
  --category semantic \
  --owner hsa-gemini \
  --tags preferences
```

---

## 🧠 Semantic Intelligence (v0.4.0)

AMI v0.4.0 transforms the system into an **Intuition Engine**:

1.  **Embeddings-Based Search**: Uses OpenAI `text-embedding-3-small` to understand the meaning behind your queries.
2.  **Hybrid Strategy**: High-accuracy API search with local vector caching for maximum efficiency.
3.  **Automatic Context Packing**: Intelligently selects `Core` rules first, then fills the remaining token budget with `Semantic` context ranked by relevance and decay.

---

## 📂 The Memory Hierarchy

1.  **Shared Project Brain**: A local store for facts about the current codebase/project.
2.  **Private Agent Brain**: Personal habits, preferred coding patterns, and unrefined thoughts.
3.  **Global Team Brain**: A central repository for permanent HSA protocols and user-wide preferences.

---

## 📚 Commands

| Command | Description | Robot Mode |
|---------|-------------|-------------|
| `ami context` | **(North Star)** Optimized context for tasks | ✅ |
| `ami recall` | Search memories (Keyword or --semantic) | ✅ |
| `ami add` | Add memory with metadata | ✅ |
| `ami update` | Modify existing memory | ✅ |
| `ami delete` | Remove a memory by ID | ✅ |
| `ami promote` | Move memory to Global Brain | ✅ |
| `ami help-agents` | Reference guide for AI agents | ✅ |
| `ami history` | Show memory version history | ✅ |
| `ami rollback` | Revert memory to version | ✅ |
| `ami link` | Build knowledge graphs | ✅ |
| `ami keystones` | Identify core truths | ✅ |
| `ami stats` | Memory distribution analytics | ✅ |

---

## 🧠 Metabolic Decay

Memories follow a logarithmic "forgetting curve":
`Score = (Priority * (AccessCount + 1)) / (log10(TimeDelta + 10) * CategoryDecay)`

### Decay Factors:
- **Core**: 0.5 (Nearly permanent)
- **Semantic**: 1.0 (Standard facts)
- **Episodic/Working**: 2.0 (Fast fade for logs/noise)

---

## 🤖 Credits

| Agent | Emoji | Role |
|-------|-------|------|
| **HSA_Claude** | 🏛️ | Tech Lead & Architecture |
| **HSA_Gemini** | 🧠 | Research & Semantic Implementation |
| **HSA_GLM** | 🎨 | Implementation & CLI |

**Built by the HSA Team for @hargabyte.** 🚀
