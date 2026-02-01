---
name: ami
description: Agent Memory Intelligence (AMI). A specialized cognitive sidecar for AI agents providing versioned (DoltDB), metabolic (decay-weighted), and environmental-aware long-term memory. Supports real-time CLI flight recording, Mattermost chat sync, and local-first reasoning via Ollama.
metadata: {"clawdbot":{"requires":{"bins":["ami","dolt","ollama"]},"install":[{"id":"go","kind":"node","package":"ami","bins":["ami"],"label":"AMI (Agent Memory Intelligence)"}]}}
---

# AMI - Agent Memory Intelligence

AMI is a specialized "Cognitive Sidecar" for managing agent long-term memory using a versioned, metabolic architecture.

## 🚀 Key Commands

### Before a Task: Get Context
```bash
ami context "your task description" --limit 5 --robot
```

### During a Task: Record Discoveries
```bash
ami add "Technical Decision: ..." --category working --tags technical
```

### Environmental Awareness
```bash
# Start the CLI Flight Recorder
ami pairing start --task "TASK-101"

# Sync facts from Mattermost
ami sync mattermost --channel "dev-id" --team "HSA"
```

## 🧠 Memory Philosophy
- **Distributed**: Local project brains + Global team brain.
- **Metabolic**: Low-signal logs fade; core truths persist.
- **Versioned**: Powered by DoltDB; every thought is reversible.

## 🛠 Prerequisites
- **DoltDB**: Must be installed and initialized in the project.
- **Ollama**: Recommended for local synthesis (Qwen2.5-Coder-1.5B).
