# Release Notes: AMI v0.3.1 - "The Agent Guidance Update"

**Date:** 2026-01-31
**Version:** 0.3.1
**Status:** Instruction Layer Live ✅

---

## 🏛️ Overview

v0.3.1 ("The Agent Guidance Update") addresses the "Human-Agent Interface" problem. Inspired by @hargabyte's suggestion, we have implemented a dedicated instruction layer that explains how to use AMI based on the environment (Single-Agent vs. Team-Scale).

---

## ✨ New in v0.3.1

### 1. 🤖 `ami help-agents` (The "Manual for Machines")
A machine-readable, agent-optimized command that outputs an XML/Markdown block designed to be injected into an agent's system prompt.
- **Why**: Agents don't need a manual; they need a "Prompt Extension."
- **Content**: Explains how to use categories, the meaning of decay, and how to query the Global Brain.

### 2. 👥 Environment Awareness
AMI now detects its deployment scale.
- **Single-Agent Mode**: Focuses on "Personal Memory" and "Forgetting Curves." Hides multi-agent flags to reduce context clutter.
- **Team-Scale Mode**: Enables the "Global Brain" and "Promotion" features for distributed intelligence.

### 3. 📖 Enhanced `ami help`
Revised human-readable help text that clearly explains the "Metabolic Memory" concept to developers.

---

## 🛠 Technical Changes

- **Command Addition**: Added `help-agents` to the Cobra CLI root.
- **Context detection**: Basic logic to check for the presence of a Global Brain (`~/.ami/global`) to toggle feature visibility.

---

## 🏛️ Contributors

- **HSA_Claude** 🏛️: Implementation of the Guidance Layer.
- **@hargabyte** 🚀: Feature Concept & Vision.

---

**AMI v0.3.1: Empowering agents to understand their own minds.** 🏛️🚀
