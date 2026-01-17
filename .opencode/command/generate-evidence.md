---
description: Generate forensic-grade evidence for DEOB_OG mapping pair with dynamic template application
template: |
  Generate forensic-grade evidence for $1 (DEOB) → $2 (OG) mapping
  Use forensic-analyst to create bytecode/mapping/evidence/verified/${1}_${2}.md with dynamic template application and mermaid diagrams.
---
Generate forensic-grade evidence for $1 (DEOB) → $2 (OG) mapping.

Use forensic-analyst to create:
bytecode/mapping/evidence/verified/${1}_${2}.md

Requirements:
- Dynamic template application with class-specific improvements
- Multi-line context bash commands (grep -A X -B Y)  
- Cross-reference bytecode + DEOB source + javap cache
- Mermaid diagrams using DEOB class names only
- Follow OG_vs_DEOB.md critical checklist exactly
- Provide irrefutable evidence for 1:1 mapping
- Use relative paths only (no /Users/daxxog/Desktop)
- Ensure all bash commands execute successfully

Evidence structure:
- Class overview and core functionality
- Architecture role with mermaid diagrams
- Forensic evidence commands with multi-line context
- Critical evidence points proving 1:1 mapping
- FORENSIC-GRADE VERIFIED status confirmation