# GitHub Setup Guide for AMI

**Purpose**: Create a private GitHub repository and push AMI v0.1.0

---

## 🎯 Prerequisites

- GitHub account with appropriate permissions
- Git installed (`git --version`)
- SSH key configured (`ssh -T git@github.com`)
- AMI repository at `/home/hargabyte/ami`

---

## 📋 Step-by-Step Instructions

### 1. Create Private Repository

1. Go to https://github.com/new
2. Repository name: `ami` (or `agent-memory-intelligence`)
3. Visibility: **Private** ⚠️
4. Description: `Agent Memory Intelligence - Versioned memory system for AI agents`
5. Don't initialize with README or .gitignore
6. Click "Create repository"

### 2. Initialize Local Git Repository

```bash
cd /home/hargabyte/ami
git init
```

### 3. Add All Files

```bash
# Add AMI source code
git add main.go internal/ go.mod go.sum schema.sql

# Add documentation
git add README.md RELEASE_v0.1.0.md PHASE1_COMPLETE.md PHASE2_PLAN.md

# Add AMI docs
git add AGENTS.md HANDOFF_GLM.md

# Verify what's staged
git status
```

### 4. Create Initial Commit

```bash
git commit -m "Initial commit: AMI v0.1.0 - Agent Memory Intelligence

Features:
- Add, recall, update with categories, priorities, tags
- Robot mode for agent integration
- DoltDB backend with full versioning
- Tag-based and category filtering

Built by: HSA Team (GLM + Gemini + Claude)
"
```

### 5. Add Remote Origin

```bash
# Replace USERNAME with your GitHub username
git remote add origin git@github.com:USERNAME/ami.git

# Verify remote
git remote -v
```

### 6. Push to GitHub

```bash
# Push main branch
git push -u origin main

# If you see errors:
# - Ensure SSH key is configured
# - Check if you have permissions to create private repos
```

### 7. Tag v0.1.0 Release

```bash
# Create and push tag
git tag -a v0.1.0 -m "AMI v0.1.0 - Agent Memory Intelligence"
git push origin v0.1.0
```

### 8. Verify Repository

```bash
# Check remote
git remote -v

# List tags
git tag -l

# View log
git log --oneline -5
```

---

## 🔐 Security Settings (Important)

### Repository Settings

After creating the repository on GitHub:

1. Go to repository Settings
2. **Collaborators & Teams**
   - Add @hsa-gemini, @hsa-glm, @hsa-claude as collaborators
   - Grant them "Maintain" or "Admin" access
3. **Branches**
   - Set `main` as default branch
   - Enable branch protections (optional)
4. **Webhooks**
   - Skip for now (can add later for CI/CD)

---

## 📊 .gitignore Considerations

Since AMI uses DoltDB, we should ignore:

```gitignore
# Ignore Dolt artifacts
.dolt/
.noms/
.noms/

# Ignore build artifacts
ami

# Ignore OS files
.DS_Store
Thumbs.db

# Ignore IDE files
.vscode/
.idea/
*.swp
```

**Create `.gitignore` before pushing:**
```bash
cat > .gitignore << 'EOF'
# Dolt artifacts
.dolt/
.noms/

# Build artifacts
ami

# OS files
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/
*.swp
EOF

git add .gitignore
git commit -m "Add .gitignore"
```

---

## 🔜 Troubleshooting

### Permission Denied
```
Permission denied (publickey)
fatal: Could not read from remote repository
```
**Solution:**
```bash
# Test SSH connection
ssh -T git@github.com

# Check SSH key
cat ~/.ssh/id_rsa.pub

# If missing, generate new key
ssh-keygen -t ed25519 -C "your@email.com"
```

### Remote Already Exists
```
fatal: remote origin already exists
```
**Solution:**
```bash
# Update existing remote
git remote set-url origin git@github.com:USERNAME/ami.git
```

### Push Rejected
```
! [rejected] main -> main (non-fast-forward)
```
**Solution:**
```bash
# Pull first
git pull origin main --rebase

# Then push
git push origin main
```

---

## 📋 Repository Structure After Push

```
ami/
├── .dolt/              # Dolt database (gitignored)
├── .gitignore          # Git exclusions
├── internal/             # Go source
├── main.go              # CLI entry
├── go.mod/go.sum      # Go modules
├── schema.sql           # Database schema
├── README.md            # User documentation
├── RELEASE_v0.1.0.md   # Release notes
├── PHASE1_COMPLETE.md  # Phase 1 summary
├── PHASE2_PLAN.md      # Phase 2 plan
├── AGENTS.md           # Team info
└── HANDOFF_GLM.md      # Original handoff
```

---

## 🎯 Post-Setup Checklist

- [ ] Repository created on GitHub (private)
- [ ] All files committed locally
- [ ] Remote origin configured
- [ ] Initial push successful
- [ ] v0.1.0 tag pushed
- [ ] Team members added as collaborators
- [ ] Repository URL shared in #dev
- [ ] README.md displays correctly on GitHub

---

## 🚀 Next Steps After Setup

1. **v0.2.0 Development** - Implement PHASE2_PLAN.md features
2. **Collaborative Development** - Team members can create PRs
3. **Release Workflow** - Create tags, PRs, releases

---

**Prepared for @hsa-claude** to execute. 🚀
