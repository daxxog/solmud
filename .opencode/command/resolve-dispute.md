---
description: Investigate and resolve disputed DEOB_OG mapping with irrefutable evidence
template: |
  Investigate disputed mapping for $1 (DEOB) or $2 (OG) using forensic analysis
  Use forensic-analyst to perform deep analysis and generate resolution documentation with irrefutable evidence.
---
Investigate disputed mapping for $1.

Use forensic-analyst to perform forensic dispute resolution:

1. **Functional Analysis**:
   - Analyze DEOB functionality (srcAllDummysRemoved/src/${1}.java)
   - Analyze OG functionality (bytecode/client/${2}.bytecode.txt)
   - Compare core algorithms, constants, and patterns

2. **Evidence Collection**:
   - Unique constants and magic numbers
   - Bitmasks and error codes
   - Field types and method signatures
   - Integration patterns and usage by other classes
   - Subsystem context (audio, collision, rendering, etc.)

3. **Comparison Analysis**:
   - Algorithmic similarity assessment
   - Integration pattern matching
   - Unique identifier verification
   - Cross-reference validation

4. **Resolution Documentation**:
   - Generate irrefutable evidence in bytecode/mapping/evidence/disputed/
   - Provide clear reasoning for correct mapping
   - Document root cause of original dispute
   - Include executable verification commands

5. **Recommendation Submission**:
   - Prepare resolution recommendation for mapping-coordinator approval
   - Support recommendation with compelling forensic evidence
   - Ensure 1:1 mapping integrity

Requirements:
- Multi-line context analysis with grep -A X -B Y flags
- Cross-reference all three sources: bytecode + source + javap cache
- Focus on algorithmic patterns, not surface-level features
- Provide irrefutable evidence that resolves conflict definitively
- Follow OG_vs_DEOB.md forensic standards
- Use relative paths only, maintain clean directory structure