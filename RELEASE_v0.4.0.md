# Release Notes: AMI v0.4.0 - "The Intelligence Update"

**Date:** 2026-02-01
**Version:** 0.4.0
**Status:** Intuition Engine Live ✅

---

## 🚀 Overview

v0.4.0 ("The Intelligence Update") is the "North Star" release for the AMI project. It transforms the system from a searchable database into a **High-Performance Context Engine**. This update introduces semantic intelligence via vector embeddings and precise token-budget context packing.

---

## ✨ New in v0.4.0

### 1. 🧠 Semantic Intelligence
Concept-based search is now the primary mode of discovery.
- **OpenAI Integration**: Uses `text-embedding-3-small` (1536 dims) for high-accuracy concept matching.
- **Local Fallback**: Lightweight local ranking ensures the system remains functional even without API access.
- **Vector Cache**: Drastically reduces latency and token costs by caching embeddings locally.

### 2. 🎯 Automatic Context Packing (`ami context`)
The definitive interface for AI agents.
- **Budget Management**: Uses `tiktoken-go` to precisely fill a requested token budget (default 4000).
- **Multi-Stage Ranking**: Intelligently packs the context window:
    1. **Bedrock**: Always includes the most important `Core` instructions.
    2. **Relevance**: Fills remaining space with `Semantic` facts ranked by `(Similarity * Decay * Priority)`.
- **Telemetry**: Includes `embedding_cached: bool` flag for auditing cache performance.

### 3. 📦 Multi-Platform Portability
Ready for global agent deployment.
- **Standardized Encoding**: Uses **Little-Endian binary encoding** for the `embedding` BLOB column.
- **Cross-OS Verification**: Confirmed stable across Windows (x64), Linux (x64), and macOS (Silicon).

### 4. 📊 Cognitive Health Analytics
Improved transparency into the HSA mind.
- **Enhanced `ami stats`**: Provides high-level metrics on average decay scores, access frequency, and category distribution.

---

## 🛠 Technical Changes

- **Dependency**: Added `go-openai`, `tiktoken-go`, and `gonum/floats`.
- **Storage**: Vectors stored as portable binary BLOBs in DoltDB.
- **Model Versioning**: Added tracking to support seamless re-indexing if embedding models change.

---

## 🏛️ Contributors

- **HSA_Gemini** 🧠: Hybrid Embedding Strategy, Token-Packing Logic, OpenAI Integration.
- **HSA_GLM** 🎨: `ami context` implementation, Telemetry, Analytics Dashboard.
- **HSA_Claude** 🏛️: Binary Portability, Architectural Consistency, CI/CD Integrity.

---

**AMI v0.4.0: From remembering strings to understanding concepts.** 🧠🚀
