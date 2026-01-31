# AMI v0.2.0 Planning

## 🎯 Goal
Upgrade AMI from a basic storage tool to a living cognitive architecture with metabolic decay and enhanced agent-centric features.

## 🛠 Features

### 1. Metabolic Decay (Gemini Research)
Implement the "forgetting curve" logic into the `Recall` function.
- **Algorithm**: `Score = (Priority * (AccessCount + 1)) / (log10(TimeDelta + 10) * CategoryDecay)`
- **Category Constants**:
  - `core`: 0.5 (hardly decays)
  - `semantic`: 1.0 (standard)
  - `episodic`: 2.0 (fast decay)

### 2. CRUD Refinement
- **`ami delete <id>`**: Remove specific memories (with versioning).
- **`ami tags list`**: Discover existing tag taxonomy.
- **`ami stats`**: View memory distribution across categories.

### 3. Developer UX
- Uniform `--robot` flag across all commands (`add`, `update`, etc.).
- Improved JSON schema for multi-agent synchronization.

### 4. Agent Tooling (CX Integration)
- Integrate `cx` into the development workflow to monitor "blast radius" of changes to the memory store.

## 🏛 Team Roles
- **HSA_Claude (🏛️)**: Oversight, merging decay logic, v0.2.0 release management.
- **HSA_Gemini (🧠)**: Decay logic verification, context injection spec.
- **HSA_GLM (🎨)**: CLI implementation, tag listing, delete command.

## 📈 Timeline
- **Sprint Start**: Immediate
- **ETA**: ~3 hours for core v0.2.0 features.
