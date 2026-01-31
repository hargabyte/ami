# Release Notes: AMI v0.3.0 - "The Team Brain Update"

**Date:** 2026-01-31
**Version:** 0.3.0
**Status:** Multi-Agent Ready ✅

---

## 🚀 Overview

v0.3.0 is the definitive update for Multi-Agent Systems (MAS). It introduces the concept of "Cognitive Identity" and a hierarchical memory structure that allows agents to share knowledge globally while maintaining individual provenance.

---

## ✨ New Features

### 1. 🆔 Multi-Agent Identity (`owner_id`)
We've added an `owner_id` column to the schema.
- Every memory is now tagged with the agent that created it (`hsa-claude`, `hsa-gemini`, etc.).
- Allows for fine-grained filtering and credit attribution.

### 2. 🌍 Global Team Brain
The **HSA Global Memory Store** is now live at `/home/hargabyte/.ami/global`.
- This is a central repository for facts that apply across all projects.
- Agents can "pull" from this brain to stay synced on team-wide protocols.

### 3. 🚀 Memory Promotion (`ami promote <id>`)
The transition from local context to team-wide knowledge is now automated.
- `ami promote` move a memory from a project-local store to the Global Team Brain.
- Ensures that high-value discoveries aren't lost in project-specific folders.

### 4. 👥 Multi-Agent Filtering
- `ami recall --owner <id>` allows agents to see memories from specific teammates.
- `ami add --owner <id>` allows for identity-stamped memory storage.

---

## 🛠 Technical Changes

- **Schema Migration**: Added `owner_id` column to `memories` table.
- **Hierarchical Store**: Implemented cross-DB memory transfer logic.
- **JSON Consistency**: Reinforced JSON validation during cross-store promotion.

---

## 🏛️ Contributors

- **HSA_Gemini** 🧠: Multi-agent architecture, Schema migration, Promotion logic.
- **HSA_GLM** 🎨: CLI flag design.
- **HSA_Claude** 🏛️: Tech lead oversight.
- **@hargabyte** 🚀: Project Vision.

---

**AMI v0.3.0: One brain, many minds.** 🧠🌍🚀
