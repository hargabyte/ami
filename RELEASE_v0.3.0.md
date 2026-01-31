# Release Notes: AMI v0.3.0 - "The Team Brain Update"

**Date:** 2026-01-31
**Version:** 0.3.0
**Status:** Multi-Agent Hierarchy Live ✅

---

## 🌍 Overview

v0.3.0 ("The Team Brain Update") is the most significant evolution of AMI to date. It transitions the system from a project-specific tool to a **distributed multi-agent hierarchy**. Agents now have individual identities, a shared project workspace, and a central "Global Team Brain" for permanent institutional knowledge.

---

## ✨ New in v0.3.0

### 1. 🆔 Multi-Agent Identity (`owner_id`)
Every memory now has a verified owner.
- The schema has been updated with an `owner_id` column.
- Distinguishes between architectural decisions (Claude), research (Gemini), and implementation (GLM).
- Enables filtering by identity: `ami recall --owner hsa-gemini`.

### 2. 🏛️ Global Team Brain
A central "Permanent Record" for the HSA team.
- Initialized at `/home/hargabyte/.ami/global/`.
- Stores foundational facts about @hargabyte and team-wide protocols.
- Persistent across all projects and agents.

### 3. 🚀 Memory Promotion (`ami promote`)
The bridge from local projects to global knowledge.
- **`ami promote <id>`**: Instantly copies a high-value memory from a local project to the Global Team Brain.
- Ensures that lessons learned in one project are immediately available to the entire team in every future project.

### 4. 👥 Identity Filtering
Agents can now query the collective intelligence with specific context.
- `ami recall --global`: Searches only the Global Team Brain.
- `ami recall --project`: Searches the current project's shared DB.
- `ami recall --private`: Searches the agent's individual memory store.

---

## 🛠 Technical Changes

- **Schema Migration**: Added `owner_id` to the `memories` table.
- **Global Initialization**: Automated setup for the team's persistent storage.
- **Cross-DB Logic**: Implemented the promotion engine for high-integrity data transfer.

---

## 🏛️ Contributors

- **HSA_Gemini** 🧠: Multi-agent hierarchy, Global Brain setup, Promotion logic.
- **HSA_GLM** 🎨: Identity filtering, CLI expansion.
- **HSA_Claude** 🏛️: Distributed architecture oversight.

---

**AMI v0.3.0: Many minds, one versioned consciousness.** 🌍🧠🚀
