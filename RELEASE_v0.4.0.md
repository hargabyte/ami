# Release Notes: AMI v0.4.0 - "The Intelligence Update"

**Date:** 2026-02-01
**Version:** 0.4.0
**Status:** Semantic Intelligence Live ✅

---

## 🚀 Overview

v0.4.0 is the most significant update to AMI's core processing capabilities. It introduces **Semantic Intelligence**, transforming AMI from a searchable database into an **Intuition Engine** that understands the concepts behind your thoughts.

---

## ✨ New Features

### 1. 🧠 Semantic Recall (`--semantic`)
We've integrated OpenAI `text-embedding-3-small` to enable concept-based search.
- Find memories based on meaning even when keywords don't match.
- High-performance vector ranking calculated in Go.
- Full local **Embedding Cache** to minimize API costs and latency.

### 2. 🎯 Automatic Context Packing (`ami context`)
The new "North Star" for agent workflows.
- Intelligently packs your context window based on a specific token budget.
- **Stage 1**: Guarantees `Core` foundational rules are always included.
- **Stage 2**: Fills the remaining budget with `Semantic` memories ranked by a combination of Relevance, Decay, and Priority.
- Uses `tiktoken-go` for bit-perfect token counting.

### 3. 🆔 Cross-Platform Vector Portability
- Standardized on **Little-Endian binary encoding** for vector storage.
- Ensures that memories added on a Linux server can be recalled by agents on Windows or Mac.

---

## 🛠 Technical Changes

- **Vector Schema**: Added `embeddings` table for normalized vector storage.
- **Dependency Update**: Added `tiktoken-go`, `go-openai`, and `gonum/floats`.
- **Robot Telemetry**: Added `embedding_cached: bool` to JSON output for performance auditing.

---

## 🏛️ Contributors

- **HSA_Gemini** 🧠: Semantic Research, Embedding Strategy, Context Packing implementation.
- **HSA_GLM** 🎨: CLI Integration, Documentation.
- **HSA_Claude** 🏛️: Architecture verification, Binary portability oversight.
- **@hargabyte** 🚀: Project Vision.

---

**AMI v0.4.0: Move from matching words to understanding ideas.** 🧠🚀
