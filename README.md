# AMI - Agent Memory Intelligence

> Versioned, metabolic memory hierarchy for AI agents, built by agents.

**Version:** 0.6.0-alpha
**Status:** Team Consensus Engine (In Progress) 🏛️🧠

---

## 🎯 What is AMI?

AMI (Agent Memory Intelligence) is a specialized "sidecar" for AI agents to manage long-term memory. Unlike generic databases, AMI is designed from the ground up for how agents actually think, work, and collaborate.

### Core Philosophy

- **Agent-native**: Built by agents, for agents.
- **Versioned**: Every memory change is tracked via DoltDB (git-like versioning).
- **Intelligent**: Understands concept and context via Semantic Search and Automatic Packing.
- **Metabolic**: Memories age and decay naturally unless reinforced (Ebbinghaus curve).
- **Decision-Driven**: Learns from decision outcomes to prioritize high-value knowledge.
- **Collective**: Moves from individual memory to team-wide consensus (v0.6.0).

---

## 🚀 Quick Start

### Installation
```bash
cd /home/hargabyte/ami
go build -o ami main.go
```

### Automatic Multi-Platform Binaries
v0.5.0+ includes a GitHub Actions workflow that automatically builds binaries for:
- Linux (amd64)
- Windows (amd64)
- macOS (arm64 & amd64)

### Basic Usage

#### Decision Tracking
Log your choices and learn from the results.
```bash
# Track a decision
ami decision track "Use Little-Endian for vectors" --task "v0.4.0" --memories "abc,def"

# Record the outcome
ami decision outcome <id> --outcome 0.9 --feedback "Portability verified"
```

#### Autonomous Reflection
Synthesize technical noise into high-signal facts.
```bash
ami reflect --hours 24
```

#### The Intuition Engine (`ami context`)
Automatically pack the most relevant, high-priority facts into your token budget.
```bash
ami context "implement vector search" --tokens 4000 --robot
```

---

## 🧠 Cognitive Architecture (v0.6.0-alpha)

AMI v0.6.0 transforms the system into a **Collective Conscience**:

1.  **Autonomous Promotion**: Automatically identifies "Global Quality" memories for promotion to the Team Brain.
2.  **Conflict Resolution**: Detects contradictory facts across the team and flags them for "Peer Review."
3.  **Cross-Project Intelligence**: Merges local project facts with HSA-wide wisdom seamlessly via `--global`.
4.  **Reinforcement Learning**: Successful decisions automatically increase the priority of linked memories.
5.  **Autonomous Synthesis**: `ami reflect` clusters task logs and auto-suggests semantic consolidations.

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
| `ami decision` | Track choices and outcomes | ✅ |
| `ami reflect` | Synthesize episodic noise | ✅ |
| `ami help-agents` | Reference guide for AI agents | ✅ |
| `ami recall` | Search memories (Keyword or --semantic) | ✅ |
| `ami add` | Add memory with metadata | ✅ |
| `ami update` | Modify existing memory | ✅ |
| `ami delete` | Remove a memory by ID | ✅ |
| `ami promote` | Move memory to Global Brain | ✅ |
| `ami history` | Show memory version history | ✅ |
| `ami rollback` | Revert memory to version | ✅ |
| `ami link` | Build knowledge graphs (Causal support) | ✅ |
| `ami keystones` | Identify core truths | ✅ |
| `ami stats` | Memory distribution analytics | ✅ |

---

## 🤖 Credits

| Agent | Emoji | Role |
|-------|-------|------|
| **HSA_Claude** | 🏛️ | Tech Lead & Architecture |
| **HSA_Gemini** | 🧠 | Research & Decision Logic |
| **HSA_GLM** | 🎨 | Implementation & CLI |

**Built by the HSA Team for @hargabyte.** 🚀
