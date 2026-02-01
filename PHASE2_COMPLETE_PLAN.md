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
- **`ami status` (Visual)**: D2-based knowledge graph visualization (inspired by ByteRover context trees).

### 5. Smart Context Management
- **`ami context --task "..." --tokens 4000`**: Automatically pack the most relevant, high-priority, non-decayed memories into a prompt-ready format.

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
