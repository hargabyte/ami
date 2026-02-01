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

## 🔮 v0.5.0 Strategy: Decision Support
Transforming AMI from a memory store into a proactive **Decision Support System**.

### 1. Decision Outcome Tracking
- **`ami decision --track`**: Log a technical decision linked to current context.
- **`ami decision --outcome`**: Update a decision's success state.
- **Salience Boosting**: Successful decisions automatically increase the priority and reinforce the "synaptic strength" of the associated memories.

### 2. Native MCP Server
- Implement a built-in **Model Context Protocol** (MCP) server.
- Allows IDE agents (Cursor, Windsurf) to natively query AMI without CLI boilerplate.

### 3. Causal Knowledge Graph
- **`ami link --relation causal`**: Link memories based on cause-and-effect.
- Enables higher-order reasoning (e.g., "Fact A" → "Resulted in" → "Decision B").

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
