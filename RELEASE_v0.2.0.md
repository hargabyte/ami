# Release Notes: AMI v0.2.0 - "The Metabolic Update"

**Date:** 2026-01-31
**Version:** 0.2.0
**Status:** Feature Rich ✅

---

## 🚀 Overview

v0.2.0 is a major milestone that transforms AMI from a structured database into a **metabolic cognitive architecture**. This update introduces natural memory decay, full CRUD uniformity, and essential cleanup tools.

---

## ✨ New Features

### 1. 🧠 Metabolic Decay (`--decay`)
We've implemented a research-backed **Ebbinghaus Forgetting Curve** for recall.
- Memories now "age" naturally based on their category.
- Frequent access "strengthens" the memory, keeping it at the top of the recall stack.
- `core` facts are practically permanent, while `episodic` logs fade quickly to reduce noise.

### 2. 🤖 Uniform Robot Mode
Total consistency for agent integration.
- Every command (`add`, `recall`, `update`, `delete`, `tags`) now supports the `--robot` flag.
- Agents can now handle the entire memory lifecycle using pure JSON.

### 3. 🧹 Memory Deletion (`ami delete`)
The much-requested cleanup tool is live.
- Allows permanent (but version-controlled) removal of incorrect or outdated memories.
- Every deletion is tracked in the Dolt history.

### 4. 🏷️ Tag Discovery (`ami tags`)
Agents can now discover what they've learned.
- `ami tags` lists all unique tags currently in the database.
- Supports `--robot` mode for programmatic discovery.

---

## 🛠 Technical Changes

- **Schema Stability**: No schema changes required, leveraging existing `accessed_at` and `access_count` fields.
- **Improved Parsing**: Switched to JSON-based SQL output parsing for 100% reliability in Robot Mode.
- **Versioning**: Each new feature is fully integrated with Dolt versioning.

---

## 📋 v0.2.0 Usage Stats

| Command | Status | New in v0.2.0 |
|---------|--------|---------------|
| `add` | ✅ | `--robot` support |
| `recall` | ✅ | `--decay` scoring |
| `update` | ✅ | `--robot` support |
| `delete` | ✅ | **NEW** |
| `tags` | ✅ | **NEW** |
| `status` | ✅ | v0.2.0 tracking |

---

## 🏛️ Contributors

- **HSA_Gemini** 🧠: Research, Decay Algorithm, Robot Mode refactor.
- **HSA_GLM** 🎨: CLI Implementation, Documentation, Testing.
- **HSA_Claude** 🏛️: Architectural validation.
- **@hargabyte** 🚀: Project Vision.

---

**AMI v0.2.0: Because a mind that remembers everything, remembers nothing.** 🧠🚀
