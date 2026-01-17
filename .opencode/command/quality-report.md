---
description: Generate scored quality report for all evidence files with improvement recommendations
agent: general
subtask: true
---
Use evidence-reviewer to analyze all evidence files in bytecode/mapping/evidence/ and generate comprehensive quality report.

For each evidence file, evaluate:
1. **Bash Command Execution** (20 pts): % commands that execute successfully
2. **Multi-line Context Usage** (20 pts): % commands with -A X -B Y flags  
3. **Non-contradictory Evidence** (15 pts): Cross-reference validation success
4. **OG_vs_DEOB.md Checklist Compliance** (25 pts): Item-by-item verification
5. **Documentation Quality** (20 pts): Forensic-grade standards assessment

Generate structured report with:
- Quality distribution by scoring tier:
  - FORENSIC-GRADE (90-100): Template candidates
  - GOOD (75-89): Acceptable with minor improvements  
  - NEEDS WORK (50-74): Major improvements needed
  - LOW QUALITY (0-49): Regeneration required
- Top 5 best evidence files (quality templates)
- Files needing improvement (specific issues identified)
- Files to regenerate (auto-queued for forensic-analyst)
- Conflict detection report (1:1 mapping violations)
- Improvement recommendations for each evidence file

Focus on quality metrics that support continuous improvement toward 100% FORENSIC-GRADE VERIFIED evidence.