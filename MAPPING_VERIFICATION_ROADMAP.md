# 1:1 Mapping Verification Implementation Roadmap
**Project**: OG vs DEOB Bytecode Forensic Mapping Verification  
**Date**: January 9, 2026  
**Status**: Active Implementation Phase  

## Executive Summary

This roadmap provides a phase-by-phase execution plan to complete forensic verification of all 73 classes in the 317 deobfuscation mapping. The project requires creating 15 missing evidence files, resolving naming inconsistencies, and implementing quality control protocols.

## Current State Analysis

### Completed Work
- ✅ **58 evidence files** created and verified
- ✅ **74 class mappings** identified in class_mapping.csv
- ✅ **Evidence template** established
- ✅ **Quality standards** defined

### Critical Gaps
- ❌ **15 missing evidence files** need creation
- ❌ **Naming inconsistencies** in existing files (20+ files)
- ❌ **Quality assurance** validation pending
- ❌ **Disputed resolutions** incomplete (7 files in disputed folder)

---

## Phase-by-Phase Execution Plan

### Phase 1: Emergency Evidence Creation (Days 1-2)

#### Daily Targets - Day 1
**Time**: 8:00 AM - 6:00 PM with breaks  
**Priority**: Critical Path Classes (High-confidence mappings)

**Morning Batch (8:00 AM - 12:00 PM)**
- `Class13_HZTFWEML.md` (100% confidence)
- `OnDemandFetcher_GHOWLKWN.md` (100% confidence)
- `Class6_CLRWXPOI.md` (100% confidence)
- `MRUNodes_GCPOSBWX.md` (100% confidence)

**Afternoon Batch (1:00 PM - 5:00 PM)**
- `Class11_LLORVYLP.md` (100% confidence)
- `Class29_SQHJOGRT.md` (100% confidence)
- `TextDrawingArea_YXVQXWYR.md` (100% confidence)
- `Class4_CDEJWOSB.md` (100% confidence)

#### Daily Targets - Day 2
**Morning Batch (8:00 AM - 12:00 PM)**
- `Class32_QPNUVGRI.md` (100% confidence)
- `Stream_MBMGIXGO.md` (100% confidence)
- `OnDemandFetcherParent_VJKFYAWG.md` (100% confidence)
- `ISAACRandomGen_JOCFVBOI.md` (100% confidence)

**Afternoon Batch (1:00 PM - 5:00 PM)**
- `NodeCache_ARZPHHDH.md` (100% confidence)
- `Class33_RJXWGZGD.md` (96% confidence)
- `Class21_ZARDZRHZ.md` (90% confidence)
- `Class18_XPBACSMK.md` (90% confidence)

#### Deliverables Phase 1
- ✅ All 15 missing evidence files created
- ✅ Each file meets forensic evidence standards
- ✅ Bash commands tested and functional
- ✅ Cross-reference validation complete

---

### Phase 2: Quality Control & Standardization (Day 3)

#### Morning: Quality Audit (8:00 AM - 12:00 PM)
**File Renaming Operations**
```bash
# Standardize naming convention issues
find bytecode/mapping/evidence/verified -name "*.md" | grep -v "^[A-Z]"
# Fix inconsistent cases like TEXTCLASS_ZTQFNQRH.md → TextClass_ZTQFNQRH.md
```

**Quality Validation Checklist**
- [ ] All files use relative paths (no /Users/daxxog/Desktop)
- [ ] Bash commands are functional and tested
- [ ] Evidence meets all 12 checklist items from OG_vs_DEOB.md
- [ ] Multiple lines of evidence shown (not single-line grep)
- [ ] Cross-references are non-contradictory

#### Afternoon: Dispute Resolution (1:00 PM - 5:00 PM)
**Resolve 7 Disputed Files**
- `NodeCache_ARZPHHDH.md` → Move to verified with proper evidence
- `Class33_RJXWGZGD.md` → Resolve duplicate with verified version
- `Class32_QPNUVGRI.md` → Complete evidence file
- `Class11_LLORVYLP.md` → Complete evidence file
- `CLASS29_SQHJOGRT.md` → Fix naming and complete
- `Class4_CDEJWOSB.md` → Complete evidence file
- `Class6_CLRWXPOI.md` → Complete evidence file

---

## Subagent Task Specifications

### Task Briefing Template for Subagents

**You are a subagent working on the OG vs DEOB forensic mapping verification project.**

**Primary Mission**: Create forensic evidence files for specific class mappings using exact standards defined in OG_vs_DEOB.md.

**Critical Instructions**:
1. Read OG_vs_DEOB.md COMPLETELY before starting any work
2. Use only relative paths (never absolute paths like /Users/daxxog/Desktop)
3. Test all bash commands to ensure they work
4. Provide multiple lines of evidence, not single-line grep results
5. Follow the exact naming convention: DEOBfilename_OBFUSCATED.md

**Quality Standards Checklist** (12 items from OG_vs_DEOB.md):
- [ ] Class overview and purpose
- [ ] Architecture relationships with mermaid diagrams
- [ ] Bash commands showing specific bytecode areas
- [ ] DEOB source code correlation
- [ ] Javap cache correlation  
- [ ] Multiple lines of context (not single lines)
- [ ] Commands verified to work
- [ ] 1:1 mapping validation
- [ ] No absolute paths
- [ ] Non-contradictory evidence
- [ ] Cross-reference validation
- [ ] Proper file naming

### Batch Assignment Strategy

**Batch 1: Anchor Classes (100% confidence)**
```bash
# 8 files for subagent batch 1
Class13_HZTFWEML
OnDemandFetcher_GHOWLKWN  
Class6_CLRWXPOI
MRUNodes_GCPOSBWX
Class11_LLORVYLP
Class29_SQHJOGRT
TextDrawingArea_YXVQXWYR
Class4_CDEJWOSB
```

**Batch 2: High-Confidence Classes**
```bash
# 7 files for subagent batch 2
Class32_QPNUVGRI
Stream_MBMGIXGO
OnDemandFetcherParent_VJKFYAWG
ISAACRandomGen_JOCFVBOI
NodeCache_ARZPHHDH
Class33_RJXWGZGD
Class21_ZARDZRHZ
```

**Batch 3: Quality Assurance Team**
```bash
# Tasks for QA subagent team
1. Validate all evidence files meet 12-point checklist
2. Test all bash commands for functionality
3. Check for absolute paths and replace with relative
4. Resolve naming inconsistencies (TEXTCLASS → TextClass)
5. Move verified disputed files to proper location
6. Create final verification report
```

---

## Risk Management Protocol

### Pre-Operation Backup Procedures
```bash
# Create timestamped backup before any file operations
BACKUP_DIR="backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"
cp -r bytecode/mapping/evidence "$BACKUP_DIR/"
git commit -am "Pre-operation backup $(date)"
```

### Validation Checkpoints
**After Each File Creation**:
- [ ] Read OG_vs_DEOB.md to confirm standards compliance
- [ ] Test all bash commands in the file
- [ ] Verify file follows naming convention
- [ ] Check for absolute paths
- [ ] Validate 1:1 mapping consistency

**After Each Phase**:
- [ ] Complete file count verification
- [ ] Bash command functionality test
- [ ] Cross-reference validation
- [ ] Quality checkpoint review

### Rollback Strategies
**If Issues Detected**:
1. **Immediate Stop**: Halt all operations
2. **Assess Impact**: Determine scope of problem
3. **Rollback**: `git checkout HEAD~1 -- bytecode/mapping/evidence/`
4. **Analyze**: Review what went wrong
5. **Correct**: Fix the issue before proceeding
6. **Verify**: Re-run affected operations

**Emergency Rollback Command**:
```bash
# Complete rollback to last known good state
git log --oneline -n 10  # Identify last good commit
git reset --hard <commit_hash>  # Rollback to that commit
```

---

## Evidence Creation Templates

### Template 1: Standard Evidence File Structure
```markdown
# Forensic Evidence: <OBFUSCATED_NAME> → <DEOBFILENAME>

## Class Identification
- **Obfuscated Name**: <OBFUSCATED_NAME>
- **Deobfuscated Name**: <DEOBFILENAME>
- **Confidence**: <score>% 
- **Evidence Type**: <Anchor|Inheritance|Signature|Behavioral>
- **Date**: <YYYY-MM-DD>

## Primary Evidence
### 1. Structural Analysis
- **Class Modifiers**: <public|private|final>
- **Superclass**: <superclass_match>
- **Interfaces**: <interface_list>
- **Field/Method Count**: <deob>/<obf>

### 2. Forensic Evidence
#### Bytecode Evidence
```bash
# Show specific bytecode patterns
grep -A 5 -B 5 "<pattern>" bytecode/client/<OBFUSCATED_NAME>.bytecode.txt
```

#### DEOB Source Evidence
```bash
# Show corresponding DEOB source
grep -A 3 -B 3 "<pattern>" srcAllDummysRemoved/src/<DEOBFILENAME>.java
```

#### Javap Cache Evidence
```bash
# Show javap cache correlation
grep -A 3 -B 3 "<pattern>" srcAllDummysRemoved/.javap_cache/<DEOBFILENAME>.javap.cache
```

### 3. Cross-Reference Validation
- **Dependencies**: <list of referenced classes>
- **References From**: <classes that reference this>
- **Validation Loop**: <how related mappings confirm>

## Unique Identifiers
- **Magic Constants**: <specific constants>
- **Method Signatures**: <distinctive patterns>
- **Field Types**: <type distribution>
- **Behavioral Patterns**: <unique operations>

## Conclusion
**Confidence**: <score>% - <summary of evidence strength>
```

### Template 2: Missing TextDrawingArea Evidence (YXVQXWYR)
```markdown
# Forensic Evidence: YXVQXWYR → TextDrawingArea

## Class Identification
- **Obfuscated Name**: YXVQXWYR
- **Deobfuscated Name**: TextDrawingArea
- **Confidence**: 100%
- **Evidence Type**: Anchor
- **Date**: January 9, 2026

## Primary Evidence

### 1. Text Rendering Methods (IRREFUTABLE)
```bash
# Look for text drawing method signatures
grep -A 10 -B 2 "drawText\|renderText\|method[0-9]*.([IILjava/lang/String;)" bytecode/client/YXVQXWYR.bytecode.txt
```

```bash
# Corresponding DEOB text drawing methods
grep -A 5 -B 2 "drawString\|drawText\|method433" srcAllDummysRemoved/src/TextDrawingArea.java
```

### 2. Font/Character Array Evidence
```bash
# Look for character arrays and font data
grep -A 15 -B 3 "char.*\[\|font.*\[\|anIntArray" bytecode/client/YXVQXWYR.bytecode.txt
```

```bash
# DEOB font data structures
grep -A 10 -B 2 "characterPixels\|fontWidth\|fontHeight" srcAllDummysRemoved/src/TextDrawingArea.java
```

### 3. Text Measurement Methods
```bash
# String width calculation methods
grep -A 8 -B 2 "getWidth\|measureText\|stringWidth" bytecode/client/YXVQXWYR.bytecode.txt
```

### 4. Color/Style Operations
```bash
# Text color and styling methods
grep -A 5 -B 2 "setColor\|textColor\|anInt.*[0-9]" bytecode/client/YXVQXWYR.bytecode.txt
```

## Unique Identifiers
- **Character Arrays**: Multiple char[] for font data
- **Text Width Calculation**: String measurement methods
- **Color Management**: Text color operations
- **Font Loading**: Character pixel data processing
- **Rendering Pipeline**: Text drawing with positioning
```

---

## Quality Improvement Checklist

### File Quality Standards (Must Pass All 15 Checks)

#### Content Standards (7 checks)
- [ ] **Class Overview**: Complete description of purpose and functionality
- [ ] **Architecture Context**: Relationship with other classes documented
- [ ] **Mermaid Diagram**: Class relationship visualization included
- [ ] **Multiple Evidence Lines**: Not single-line grep results
- [ ] **Context Provided**: Before/after context for all evidence
- [ ] **Non-Contradictory**: All evidence consistent and verifiable
- [ ] **1:1 Validation**: Each obfuscated class maps to single DEOB class

#### Technical Standards (8 checks)
- [ ] **Bash Commands**: All commands functional and tested
- [ ] **Relative Paths**: No absolute paths like /Users/daxxog/Desktop
- [ ] **Bytecode Evidence**: Specific bytecode areas shown with grep
- [ ] **DEOB Source**: Corresponding source code sections shown
- [ ] **Javap Cache**: Cache correlation evidence provided
- [ ] **File Naming**: Correct DEOB_OBFUSCATED.md format
- [ ] **Unique Identifiers**: Magic constants, signatures, patterns identified
- [ ] **Cross-Reference**: Dependencies and references validated

### Naming Convention Conversion Guide

#### Correct Naming Pattern
```
Format: <DEOBClassName>_<OBFUSCATED_NAME>.md
Example: TextDrawingArea_YXVQXWYR.md
```

#### Current Inconsistencies to Fix
```bash
# Files with incorrect naming
TEXTCLASS_ZTQFNQRH.md → TextClass_ZTQFNQRH.md
STREAMLOADER_XTGLDHGX.md → StreamLoader_XTGLDHGX.md
SOUNDS_JHDAGNBV.md → Sounds_JHDAGNBV.md
SKILLS_YUXCUCXD.md → Skills_YUXCUCXD.md
NPC_CWNCPMLX.md → NPC_CWNCPMLX.md (already correct)
OBJECT5_OPNPFUJE.md → Object5_OPNPFUJE.md (already correct)
NODESUBLIST_LHGXPZPG.md → NodeSubList_LHGXPZPG.md
```

---

## Success Metrics Dashboard

### Quantitative Targets

#### File Completion Metrics
- **Total Evidence Files Target**: 73
- **Currently Completed**: 58
- **Missing Files**: 15
- **Daily Target**: 7-8 files
- **Project Timeline**: 3 days

#### Quality Score Metrics
- **Evidence Quality Score**: Target 95%+ (15/15 checklist items)
- **Bash Command Success Rate**: Target 100%
- **Cross-Reference Validation**: Target 100%
- **Naming Convention Compliance**: Target 100%

### Qualitative Benchmarks

#### Evidence Strength Levels
- **IRREFUTABLE (100%)**: Exact byte-for-byte matches, magic constants
- **STRONG (95-99%)**: Multiple unique identifiers, cross-references
- **SOLID (90-94%)**: Strong correlation with minor variations
- **ADEQUATE (85-89%)**: Inheritance-based or signature-based matches

#### Verification Completeness
- **All 73 classes** have evidence files
- **All evidence** passes 12-point checklist validation
- **All cross-references** are validated and consistent
- **All bash commands** are functional and tested
- **All files** use relative paths

### Progress Tracking Methodology

#### Daily Progress Report Template
```markdown
## Day X Progress Report - YYYY-MM-DD

### Files Completed Today
- [ ] File1_ClassName_OBFUSCATED.md ✅/❌
- [ ] File2_ClassName_OBFUSCATED.md ✅/❌
- [ ] File3_ClassName_OBFUSCATED.md ✅/❌

### Quality Metrics
- **Evidence Files Created**: X/8 target
- **Quality Score Average**: XX%
- **Bash Command Success**: X/X tested
- **Checklist Compliance**: X/X items

### Blockers/Risks
- **Issue**: Description of any problems
- **Mitigation**: Steps taken to resolve
- **Impact**: Effect on timeline if any

### Tomorrow's Targets
- [ ] Primary files: List of priority files
- [ ] Secondary tasks: Quality checks, naming fixes
```

#### Cumulative Progress Tracking
```bash
# Progress tracking command
echo "=== PROGRESS TRACKER ==="
echo "Total Evidence Files: $(find bytecode/mapping/evidence/verified -name "*.md" | wc -l)"
echo "Files Meeting Standards: $(grep -l "## Primary Evidence" bytecode/mapping/evidence/verified/*.md | wc -l)"
echo "Disputed Files: $(find bytecode/mapping/evidence/disputed -name "*.md" | wc -l)"
echo "Missing Files: $(grep -v "deobfuscated_name" bytecode/mapping/class_mapping.csv | wc -l) - $(find bytecode/mapping/evidence/verified -name "*.md" | wc -l)"
```

---

## Implementation Commands Reference

### Essential Commands for Subagents

#### File Creation and Validation
```bash
# Create new evidence file from template
cp bytecode/mapping/evidence/evidence_template.md "bytecode/mapping/evidence/verified/ClassName_OBFUSCATED.md"

# Validate file meets standards
grep "## Primary Evidence" bytecode/mapping/evidence/verified/ClassName_OBFUSCATED.md

# Check for absolute paths (should return nothing)
grep "/Users/daxxog/Desktop" bytecode/mapping/evidence/verified/ClassName_OBFUSCATED.md
```

#### Evidence Collection Commands
```bash
# Search for specific patterns in bytecode
grep -A 10 -B 5 "method[0-9]*\|field[0-9]*" bytecode/client/OBFUSCATED.bytecode.txt

# Find corresponding DEOB patterns
grep -A 5 -B 2 "method.*\|field.*" srcAllDummysRemoved/src/ClassName.java

# Check javap cache for correlations
grep -A 3 -B 3 "Constant.*\|Method.*" srcAllDummysRemoved/.javap_cache/ClassName.javap.cache
```

#### Quality Assurance Commands
```bash
# Test all bash commands in evidence files
for file in bytecode/mapping/evidence/verified/*.md; do
  echo "Testing commands in $file"
  # Extract and test bash blocks (manual verification needed)
done

# Check naming convention compliance
find bytecode/mapping/evidence/verified -name "*.md" | grep -v "^[A-Z]"
```

This roadmap provides a complete, actionable plan for finishing the 1:1 mapping verification project with forensic-grade accuracy and quality control.