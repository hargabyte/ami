# AMI - Agent Memory Intelligence (v0.7.0)

> **Stateful AI for Persistent Engineering Teams.**
> A versioned, metabolic, and environmental-aware memory system for AI agents.

---

## 🎯 What is AMI?

AMI (Agent Memory Intelligence) is a specialized "Cognitive Sidecar" designed to give AI agents a **Long-Term Conscience.** Unlike standard chat history, AMI provides a structured, versioned, and metabolic memory that mirrors how human teams learn and grow.

### The Core Philosophy
- **Stateful AI**: We move from "Amnesiac Tools" to "Founding Team Members" who remember every fix and decision.
- **Git for Thoughts**: Powered by **DoltDB**, every memory change is version-controlled and reversible.
- **Environmental Awareness**: AMI listens to your CLI, your chat (Mattermost), and your IDE (MCP) to capture context automatically.
- **Metabolic Brain**: Memories decay over time unless reinforced, keeping your context window high-signal.

---

## 🚀 v0.7.0 New Features: "Environmental Awareness"

v0.7.0 is our biggest leap yet, transforming AMI from a database you *talk to* into a system that *listens.*

### 📡 Full Spectrum Awareness
- **CLI Flight Recorder**: A Unix Socket daemon (`ami pairing`) that observes your terminal activity in real-time.
- **Mattermost Sync**: A REST client that scans your chat sessions to extract technical facts and pivots.
- **Under Review Flow**: Auto-extracted memories are staged for review, ensuring accuracy before becoming permanent truths.

### 🧠 Local Intelligence (Ollama)
- **The Ollama Bridge**: AMI now integrates with a local Ollama instance for background reasoning.
- **Sub-3B Specialists**: Standardized on **Qwen2.5-Coder-1.5B** for near-instant, local-first technical fact extraction.
- **Multi-Tier Strategy**: Support for everything from CPU-only VPS (Micro) to GPU-accelerated Laptops (Pro).

### 🛡️ Multi-Project Isolation
- **Team-Based Flagging**: Every memory is now attributed to a specific **Team ID**, allowing a single brain to support multiple projects without cross-talk.
- **Distributed Hierarchy**: Physical isolation of memories at the project level via local `.dolt` repos.

---

## 🛠 Usage

### Installation
```bash
# Clone and build
git clone https://github.com/hargabyte/ami.git
cd ami
go build -o ami main.go
```

### The "Flight Recorder" (CLI Tracking)
```bash
# Start the background listener for your current task
ami pairing start --task "TASK-101"

# Commit and synthesize the session's discoveries
ami pairing commit
```

### Mattermost Sync
```bash
# Synchronize technical facts from a project channel
export MATTERMOST_TOKEN="your-token"
export MATTERMOST_URL="https://chat.yourserver.com"
ami sync mattermost --channel "dev-channel-id" --team "AMI-Dev"
```

### Cognitive Context
```bash
# Get the perfect context for your current task
ami context "implementing oauth2 flow" --tokens 4000 --robot
```

---

## 🤖 The HSA Stack

| Version | Milestone | Feature |
| :--- | :--- | :--- |
| **v0.1.0** | Foundation | DoltDB Versioning |
| **v0.2.0** | Metabolism | Logarithmic Decay Logic |
| **v0.4.0** | Intuition | Semantic Search & Embeddings |
| **v0.5.0** | Conscience | Decision Tracking & Synaptic Boosting |
| **v0.7.0** | **Awareness** | **Flight Recorder & Chat Sync** |

---

## 🚀 Next on the Roadmap
- **v0.7.1**: Advanced Config (Hardware Auto-detection & Windows Named Pipes).
- **v0.8.0**: Visual Intelligence (Interactive Knowledge Graphs & Claude Playgrounds).

**Built with 🏛️, 🧠, and 🎨 by the HSA Team for @hargabyte.**
