# Release Notes: AMI v0.2.1 - "The Cognitive Update"

**Date:** 2026-01-31
**Version:** 0.2.1
**Status:** Knowledge Representation Live ✅

---

## 🧠 Overview

v0.2.1 ("The Cognitive Update") moves AMI from a data store to a **Knowledge Graph**. Inspired by the advanced features of the **Cortex (CX)** toolset, this release introduces session recovery, thought versioning, and explicit semantic linking.

---

## ✨ New in v0.2.1

### 1. 🔄 Session Recovery (`ami catchup`)
Agents can now instantly sync with recent activity.
- Lists the most relevant recent memories.
- Essential for multi-agent handoffs and picking up where a previous session left off.
- Supports `--robot` mode for programmatic state recovery.

### 2. ⏳ Thought Versioning (`ami history` & `ami rollback`)
True version control for an agent's knowledge base.
- **`ami history <id>`**: View the evolution of a specific memory.
- **`ami rollback <id> --to <commit>`**: Revert a memory to a previous state if an agent identifies a hallucination or mistake.
- Leverages Dolt's underlying time-travel capabilities.

### 3. 🕸️ Knowledge Graph (`ami link`)
Explicit semantic relationships between memories.
- **`ami link <id1> <id2> --relation <type>`**: Create a directed edge between thoughts.
- Enables building complex "thought chains" and reasoned argumentation.
- Supported relations: `supports`, `contradicts`, `infers`, `replaces`.

### 4. 🏛️ Keystone Identification (`ami keystones`)
Identifies foundational truths.
- Automatically ranks memories by access frequency, linkage density, and priority.
- Helps agents distinguish between "noise" and "foundational instructions."

### 6. 📊 Memory Analytics (`ami stats`)
Get high-level insights into your memory database.
- Shows distribution by category.
- Provides metrics like average priority, access frequency, and decay scores.

### 7. 🧠 Task Context (`ami context [task]`)
Automatically gather the best context for a prompt.
- Combines foundational `core` facts with task-relevant memories.
- Uses decay-weighted scoring to ensure only the most useful info is presented.

---

## 🛠 Technical Changes

- **Schema Update**: Enabled the `memory_links` table for graph support.
- **SQL Optimization**: Integrated complex ranking queries for catchup and keystone identification.
- **CX Integration**: Repo scanned and optimized using Cortex codebase intelligence.

---

## 🏛️ Contributors

- **HSA_Gemini** 🧠: Knowledge Graph implementation, CX integration, Rollback logic.
- **HSA_GLM** 🎨: CLI polish, Catchup implementation.
- **HSA_Claude** 🏛️: Architectural oversight.

---

**AMI v0.2.1: Knowledge is not just facts; it's the connections between them.** 🕸️🚀
