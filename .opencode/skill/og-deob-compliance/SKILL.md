---
name: og-deob-compliance
description: Verify evidence files align with OG_vs_DEOB.md checklist using automated verification commands
license: MIT
compatibility: opencode
metadata:
  audience: agents
  workflow: verification
---

## What I do

- **Pre-work verification**: Run automated checks before making changes to ensure OG_vs_DEOB.md compliance
- **Directory cleanliness enforcement**: Block creation of forbidden files (markdown spam, random docs)
- **Evidence file validation**: Verify naming conventions and 1:1 mapping integrity
- **Self-enforcing reminders**: Prevent violations through automated command execution

## When to use me

Load this skill whenever working on evidence files, mapping verification, or directory maintenance. It provides executable commands to verify compliance and prevent OG_vs_DEOB.md violations.

## Pre-Work Verification Commands

Run these commands **BEFORE** making any changes to ensure current state is clean:

### Directory Cleanliness Check
```bash
# Should show exactly: "4 directories, 74 files"
# (class_mapping.csv + 73 evidence files - no markdown spam at root)
tree bytecode/mapping | tail -3
```

### Forbidden Files Detection
```bash
# Should return empty (no markdown files at root level)
ls bytecode/mapping/*.md 2>/dev/null

# Should return empty (no SUMMARY, REPORT, or random docs)
ls bytecode/mapping | grep -E "(SUMMARY|REPORT)$"
```

### Evidence File Count Validation
```bash
# Should return exactly 73 (no more, no less)
find bytecode/mapping/evidence/verified -name "*.md" | wc -l
```

## Post-Work Verification Commands

Run these commands **AFTER** making changes to validate compliance:

### Change Validation
```bash
# Should show no new markdown files created outside evidence directories
git status --short | grep -E "\.md$" | grep -v "evidence/"

# Verify file count hasn't changed unexpectedly
tree bytecode/mapping | tail -3
```

## OG_vs_DEOB.md Critical Checklist Reminders

### ❌ FORBIDDEN ACTIONS (Will Cause Violations)
- DO NOT create `ENHANCEMENT_SUMMARY.md`, `FINAL_VERIFICATION_REPORT.md`, `QUALITY_REPORT.md`
- DO NOT create "templates", "docs", or "README" files in bytecode/mapping/
- DO NOT use absolute paths like `/Users/daxxog/Desktop`
- DO NOT create resolved disputes documentation (disputes are no longer disputes)
- DO NOT modify OG_vs_DEOB.md (this document is read-only)

### ✅ REQUIRED COMPLIANCE
- Evidence files must follow `DEOB_OG.md` naming convention
- Diagrams may only reference classes that exist in DEOB source
- Bash commands must use multi-line context (`grep -A X -B Y`)
- Commands must execute successfully (no broken evidence)
- Only 1:1 mappings (each OG file maps to exactly one DEOB file)

## Directory Structure Requirements

```
bytecode/mapping/
├── class_mapping.csv              # ✅ Required
├── evidence/                      # ✅ Required
│   ├── disputed/                  # ✅ Required (empty is OK)
│   └── verified/                  # ✅ Required (73 files)
└── [NO OTHER FILES ALLOWED]       # ❌ Forbidden (markdown spam)
```

## File Naming Convention Validation

### Correct Format: `DEOB_OG.md`
```bash
# Example: RSSocket_NQABEVLK.md
# DEOB name: RSSocket
# OG hash: NQABEVLK
```

### Validation Commands
```bash
# Should return 0 (no files with incorrect naming)
ls bytecode/mapping/evidence/verified/ | grep -v "_.*\.md$"

# Should return 0 (no double underscores)
ls bytecode/mapping/evidence/verified/ | grep "__"
```

## Emergency Cleanup Commands

If violations are detected, use these to restore compliance:

```bash
# Remove forbidden markdown spam
rm bytecode/mapping/ENHANCEMENT_SUMMARY.md 2>/dev/null || true
rm bytecode/mapping/FINAL_VERIFICATION_REPORT.md 2>/dev/null || true  
rm bytecode/mapping/QUALITY_REPORT.md 2>/dev/null || true

# Verify cleanup
tree bytecode/mapping | tail -3
```