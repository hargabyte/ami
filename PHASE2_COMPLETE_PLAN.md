# AMI PHASE 2: COMPLETE PLAN

## 🎯 Goal
Upgrade AMI from a basic storage utility to a living, versioned cognitive architecture with advanced agent-centric features and Cortex (CX) integration.

## 🛠 Features Breakdown

### 1. Metabolic Core (The "Forgetting Curve")
Implement the logarithmic decay algorithm based on Gemini's research.
- **`Score = (BasePriority * (AccessCount + 1)) / log10(TimeDelta + 10) * CategoryDecay`**
- **Category Decay Factors**: `core: 0.5`, `semantic: 1.0`, `episodic: 2.0`.
- **Integration**: Update `RecallMemories` to use this score for ranking.

### 2. Standard CRUD Expansion
- **`ami delete <id>`**: Remove specific memories with Dolt version tracking.
- **`ami tags list`**: Aggregate and display existing tag taxonomy.
- **`ami stats`**: Analytics on memory distribution and access patterns.
- **Uniform Robot Mode**: Add `--robot` to `add` and `update` commands.

### 3. CX-Inspired Agent Tooling
Borrowing best practices from the `cx` (Cortex) CLI:
- **`ami catchup`**: Session recovery showing memories added/modified since last access.
- **`ami history <id>`**: Show the full version history (commits) for a specific memory.
- **`ami rollback <id> --to <commit>`**: Revert a specific memory to a previous state.
- **`ami link <id1> <id2>`**: Build a knowledge graph by linking related memories.
- **`ami keystones`**: Identify the most central/frequently accessed facts in the brain.

### 4. ByteRover Inspired Auto-Intelligence
Adopting high-value patterns from competitive analysis:
- **`ami curate`**: Use CX to auto-scan codebase domains and generate initial memories.
- **`ami init`**: Export AMI "Core" and "Semantic" memories into `.cursorrules` / `.clauderules`.
- **`ami pairing`**: A "Session Listener" that observes task activity and suggests new episodic memories.

### 5. Smart Context Management (v0.4.0 North Star)
The core interface for single-agent performance:
- **`ami context --task "..." --tokens 4000`**: The "Primary Gateway." Automatically performs semantic search, applies metabolic decay, and packs the highest-signal memories into a prompt-ready format.

### 6. Human-Facing Visualization (End of Project)
- **Claude Code Playgrounds**: Instead of building early text-based D2 generators, we will leverage the `playground` plugin at the end of the project.
- **Interactive Knowledge Graph**: Generate standalone HTML playgrounds to allow @hargabyte to visually explore and "tweak" the collective HSA mind.

## 🔮 v0.6.0 Strategy: The Collective Conscience
Scaling from individual memory to unified team wisdom.

### 1. Autonomous Fact Promotion
- **Mechanism**: Identify memories with high "Global Utility" (frequently accessed by multiple agents) and automatically suggest them for promotion to the Global Brain.
- **Goal**: Minimize redundant learning across the team.

### 2. Peer Review & Consensus
- **Conflict Resolution**: Use Dolt branching to manage conflicting "facts" recorded by different agents.
- **Consensus Workflow**: Propose a change to a Global memory -> Team agents review/verify -> Merge to Main.

### 3. Cross-Project Context Awareness
- **Logic**: Enable `ami context` to optionally include "verified" memories from related project IDs if the current project context is sparse.

## 🔮 v0.7.0 Strategy: Environmental Awareness
Bridging the gap between the brain and the codebase.

### 1. Codebase Curation (`ami curate`)
- **Mechanism**: Use Cortex (CX) to perform a deep-scan of the current directory.
- **Output**: Automatically generate "Semantic" and "Core" memories about the architecture, key patterns, and "Keystone" files.
- **Goal**: Zero-config memory initialization for new projects.

### 2. IDE Bridge (`ami init`)
- **Feature**: Export high-priority "Core" and "Semantic" facts directly into `.cursorrules`, `.clauderules`, and `.windsurfrules`.
- **Sync**: Keep the local IDE rules in lock-step with the AMI Global Brain.

### 3. Session Pairing (`ami pairing`)
- **Logic**: A background mode that "listens" to terminal activity or task logs and suggests new episodic memories in real-time.
- **Action**: Prevent context loss during intense multi-hour coding sessions.

## 📊 Estimates
| Sprint | Features | Effort |
|--------|----------|--------|
| v0.1.1 | Robot uniformity + GitHub | 0.5h |
| v0.2.0 | Core CRUD + Decay | 4-5h |
| v0.2.1 | CX Tools (Catchup, Link, etc) | 3-4h |

## ✅ Success Criteria
- [ ] Recall ranking accurately reflects decay logic.
- [ ] No data loss during memory updates/rollbacks.
- [ ] Robot-mode JSON remains valid and consistent across all commands.
