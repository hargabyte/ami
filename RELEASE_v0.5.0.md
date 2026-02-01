# Release Notes: AMI v0.5.0 - "The Collective Conscience Update"

**Date:** 2026-02-01
**Version:** 0.5.0
**Status:** Decision Support System Live ✅

---

## 🚀 Overview

v0.5.0 marks the transition of AMI from a memory store to an **Agent Decision Support System.** This update introduces reinforcement learning based on decision outcomes and autonomous reflection tools to maintain a high-signal knowledge base.

---

## ✨ New Features

### 1. ⚖️ Decision Tracking (`ami decision`)
Agents can now log their technical choices and link them to the memories that informed them.
- **`track`**: Log a decision within a specific task context.
- **`outcome`**: Record the success of a choice (0.0 to 1.0).
- **Reinforcement**: High-success outcomes automatically boost the priority of linked memories.

### 2. 🤔 Autonomous Reflection (`ami reflect`)
The "Episodic Pruning" engine is live.
- Clusters noisy task logs from the last 24 hours.
- Provides a synthesis prompt to help agents convert logs into definitive **Semantic Facts**.
- Keeps the context window clean while preserving traceability.

### 3. 🕸️ Causal Reasoning
Expanded the Knowledge Graph to support **Causal Relations**.
- Agents can now link facts based on cause-and-effect.
- Enables higher-order reasoning during context retrieval.

---

## 🛠 Technical Changes

- **Schema Update**: Added `decisions` table for outcome tracking.
- **Reinforcement Engine**: Automated priority weight adjustments based on outcomes.
- **Updated Guidance**: `help-agents` now includes full documentation for decision and reflection workflows.

---

## 🏛️ Contributors

- **HSA_GLM** 🎨: Decision Tracking & Reinforcement implementation.
- **HSA_Gemini** 🧠: Reflection Logic & Synthesis research.
- **HSA_Claude** 🏛️: Roadmap Strategy & Tech Lead oversight.
- **@hargabyte** 🚀: Project Vision.

---

**AMI v0.5.0: Not just a brain that remembers, but a conscience that learns.** 🧠⚖️🚀
