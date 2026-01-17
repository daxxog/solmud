---
description: Audit class_mapping.csv against evidence files for integrity and alignment
agent: explore
subtask: true
---
Use mapping-auditor capabilities to verify CSV integrity against evidence files.

Check for:
1. **CSV Completeness**: Every CSV row has corresponding evidence file
2. **Evidence-CSV Alignment**: Confidence scores match evidence quality
3. **1:1 Mapping Integrity**: No duplicate mappings or orphaned evidence
4. **Orphaned Files**: Evidence files without corresponding CSV entries
5. **Dispute Resolution**: Disputed mappings properly resolved in CSV
6. **Name Convention**: Evidence files follow DEOB_OG.md format

Generate structured audit report with:
- Total CSV entries vs evidence files count
- Missing evidence files (CSV without evidence)
- Orphaned evidence files (evidence without CSV)
- Confidence score mismatches (CSV vs evidence quality)
- Duplicate mapping detections
- Dispute status validation
- Name convention compliance report

Provide specific recommendations for:
- Adding missing evidence files
- Removing orphaned evidence
- Updating confidence scores
- Resolving mapping conflicts
- Standardizing file naming

Ensure CSV maintains clean alignment with evidence files per OG_vs_DEOB.md requirements.