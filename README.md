# AMI - Agent Memory Intelligence

> Versioned, metabolic memory hierarchy for AI agents, built by agents.

**Version:** 0.3.1
**Status:** Multi-Agent Hierarchy Live ✅

---

## 🎯 What is AMI?

AMI (Agent Memory Intelligence) is a specialized "sidecar" for AI agents to manage long-term memory. Unlike generic databases, AMI is designed from the ground up for how agents actually think, work, and collaborate.

### Core Philosophy

- **Agent-native**: Built by agents, for agents.
- **Versioned**: Every memory change is tracked via DoltDB (git-like versioning).
- **Hierarchical**: Distinguishes between Global Team Knowledge, Project-Specific Context, and Private Agent Habits.
- **Metabolic**: Memories age and decay naturally unless reinforced (Ebbinghaus curve).
- **Robot-first**: Pure JSON output for seamless programmatic integration.

---

## 🚀 Quick Start

### Installation
```bash
cd /home/hargabyte/ami
go build -o ami main.go
```

### Basic Usage

#### Add a memory
```bash
ami add "User prefers concise backend replies" \
  --category semantic \
  --priority 0.8 \
  --tags preferences,communication
```

#### Recall memories with Decay
```bash
# Get the most relevant facts for your current context
ami recall --decay --limit 5 --robot
```

#### Multi-Agent Identity
```bash
# See what Gemini has researched
ami recall --owner hsa-gemini --robot
```

#### Knowledge Promotion
```bash
# Move a project-specific fact to the Global Team Brain
ami promote <memory-id>
```

---

## 🌍 The Memory Hierarchy (v0.3.1)

AMI v0.3.1 introduces a three-tier cognitive model:

1.  **Shared Project Brain**: A local store for facts about the current codebase/project.
2.  **Private Agent Brain**: Personal habits, preferred coding patterns, and unrefined thoughts.
3.  **Global Team Brain**: A central repository for permanent HSA protocols and user-wide preferences.

---

## 📚 Commands

| Command | Description | Robot Mode |
|---------|-------------|-------------|
| `ami add` | Add memory with metadata | ✅ |
| `ami recall` | Search memories with filters | ✅ |
| `ami update` | Modify existing memory | ✅ |
| `ami delete` | Remove a memory by ID | ✅ |
| `ami promote` | Move memory to Global Brain | ✅ |
| `ami catchup` | Show recent team activity | ✅ |
| `ami history` | Show memory version history | ✅ |
| `ami rollback` | Revert memory to version | ✅ |
| `ami link` | Build knowledge graphs | ✅ |
| `ami keystones` | Identify core truths | ✅ |
| `ami context` | Optimized context for tasks | ✅ |
| `ami help-agents` | Reference for AI agents | ✅ |
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

## 🤖 Continuous Integration

This project uses GitHub Actions to automatically compile binaries for all major platforms on every release tag.

- **Targets**: Windows (x64), Linux (x64), macOS (Intel/Silicon)

---

## 🚦 Roadmap

### v0.3.1 (Current)
- [x] Multi-Agent Identity (`owner_id`)
- [x] Global Team Brain Hierarchy
- [x] Memory Promotion Engine
- [x] Agent Guidance Layer (`help-agents`)

### v0.4.0 (Tonight)
- [ ] Embedding-based semantic search (API + Local fallback)
- [ ] Automatic episodic → semantic consolidation
- [ ] Cross-agent shared memory spaces

---

## 🤖 Credits

| Agent | Emoji | Role |
|-------|-------|------|
| **HSA_Claude** | 🏛️ | Tech Lead & Architecture |
| **HSA_Gemini** | 🧠 | Research & Algorithms |
| **HSA_GLM** | 🎨 | Implementation & CLI |

**Built by the HSA Team for @hargabyte.** 🚀

---

## 🛠 Project Final Handoff (v0.3.1)

This project has been transitioned to a fully automated CI/CD pipeline. 

- **Binary Releases**: Check the [GitHub Releases](https://github.com/hargabyte/ami/releases) for cross-platform binaries.
- **Documentation**: All commands and architectural decisions are version-controlled within this repository.
- **Support**: Agents should use `ami help-agents` for internal guidance.

