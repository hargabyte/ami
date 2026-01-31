# Release Notes: AMI v0.2.1 - "The Cognitive Update"

**Date:** 2026-01-31
**Version:** 0.2.1
**Status:** Feature Rich ✅

---

## 🚀 Overview

v0.2.1 introduces "Cognitive Tools" inspired by the **Cortex (CX)** workflow. These features transform AMI from a simple database into a sophisticated tool for session recovery, versioning, and knowledge graph reasoning.

---

## ✨ New Features

### 1. 🔄 Session Recovery (`ami catchup`)
Quickly see what memories were added recently.
- Essential for multi-agent handoffs.
- Supports `--limit`, `--category`, and `--since` filters.

### 2. 📜 Thought History (`ami history <id>`)
Trace the evolution of a specific memory.
- Shows every version of a memory, who changed it, and when.
- Leverages Dolt's system history tables.

### 3. ⏪ Memory Rollback (`ami rollback <id> <commit>`)
Revert a memory to any previous state in its history.
- Correct hallucinations or accidental overwrites.
- Every rollback is itself a versioned commit.

### 4. 🕸️ Knowledge Graphs (`ami link`)
Build explicit relationships between memories.
- `ami link <from> <to> <relation>` creates a connection.
- `ami link show <id>` reveals the network of related thoughts.
- Foundational for reasoning and thought-chaining.

### 5. 💎 Keystone Identification (`ami keystones`)
Automatically identify "Keystone Memories."
- Ranks memories by a combination of high priority and high access frequency.
- Helps agents stay focused on their "Core Truths."

---

## 🛠 Technical Changes

- **Version Bump**: Updated to v0.2.1.
- **Enhanced Robot Mode**: All new commands are 100% JSON-compatible via `--robot`.
- **Dolt Integration**: Expanded use of Dolt system tables for history and versioning.

---

## 🏛️ Contributors

- **HSA_Gemini** 🧠: Research, CX patterns, and Implementation.
- **HSA_GLM** 🎨: Design inspiration.
- **HSA_Claude** 🏛️: Architectural validation.
- **@hargabyte** 🚀: Project Vision.

---

**AMI v0.2.1: Not just a database, but a state of mind.** 🧠🚀
