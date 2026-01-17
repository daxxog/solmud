---
description: Verify all evidence files in bytecode/mapping/evidence/ for quality and compliance
template: |
  Verify all evidence files in bytecode/mapping/evidence/ using evidence-reviewer and forensic-analyst collaboration
  Use mapping-coordinator to orchestrate verification workflow across all evidence files.
---
Use mapping-coordinator to verify all evidence files in bytecode/mapping/evidence/:

1. Request evidence-reviewer to score all evidence files using 100-point methodology
2. Identify files needing regeneration (<50 points) and conflicts
3. Delegate to forensic-analyst for regeneration and conflict resolution
4. Ensure OG_vs_DEOB.md checklist compliance across all files
5. Generate comprehensive quality report with scoring distribution
6. Verify 1:1 mapping integrity and resolve any conflicts
7. Maintain clean, auditable directory structure

Focus on achieving FORENSIC-GRADE VERIFIED status (≥90 points) for all evidence files.