# Enhanced Makefile for 317 Deob Mapping Project
# Optimized workflow with agent enhancements

.PHONY: help verify-all verify-evidence verify-cleanliness

# Helper targets
help:
	@echo "317 Deob Mapping Helper Commands"
	@echo "========================================="
	@echo "verify-evidence  - Use evidence-reviewer subagent for verification"
	@echo "verify-csv        - Agent-driven CSV verification"
	@echo "verify-count      - Count evidence files"
	@echo "evidence-stats    - Evidence directory statistics"
	@echo "verify-naming     - Check file naming conventions"
	@echo "check-paths       - Check for absolute paths"
	@echo "verify-structure   - Verify directory structure"
	@echo "verify-all        - Complete verification"
	@echo "verify-cleanliness- Run OG_vs_DEOB.md compliance verification"

# OG_vs_DEOB.md compliance verification using skill
verify-cleanliness:
	@bash tools/verify_cleanliness.sh
	@echo "OG_vs_DEOB.md compliance complete."

# Evidence verification
verify-evidence:
	@echo "=== EVIDENCE VERIFICATION ==="
	@echo "Use evidence-reviewer subagent for binary verification"
	@echo ""
	@echo "Verification Gates:"
	@echo "  - Command Execution: All bash commands work?"
	@echo "  - Three-Way Evidence: Bytecode + Source + Javap Cache?"
	@echo "  - Multi-Line Context: grep -A X -B Y flags?"
	@echo "  - No Template Patterns: No 'unique to obfuscated' sections?"
	@echo "  - DEOB Diagrams Only: No obfuscated names in diagrams?"
	@echo "  - Cross-Reference Verification: Uniqueness proven?"
	@echo ""
	@echo "NO SCORING - Binary PASS/FAIL gates only"

# Enhanced verification commands
verify-count:
	@echo "Evidence Files Count:"
	@ls bytecode/mapping/evidence/verified/ | wc -l
	@echo "Expected: 73"
	@echo "Disputed Files:"
	@ls bytecode/mapping/evidence/disputed/ | wc -l

evidence-stats:
	@echo "=== EVIDENCE DIRECTORY STATISTICS ==="
	@echo "Verified files:"
	@ls bytecode/mapping/evidence/verified/ | wc -l
	@echo "Disputed files:"
	@ls bytecode/mapping/evidence/disputed/ | wc -l
	@echo "Total evidence files:"
	@find bytecode/mapping/evidence -name "*.md" | wc -l

verify-csv:
	@echo "Verifying CSV alignment with evidence files..."
	@echo "Total CSV entries:"
	@tail -n +2 bytecode/mapping/class_mapping.csv | wc -l
	@echo "Total evidence files:"
	@ls bytecode/mapping/evidence/verified/ | wc -l

verify-naming:
	@echo "Checking file naming convention (DEOB_OG.md)..."
	@ls bytecode/mapping/evidence/verified/ | head -5

check-paths:
	@echo "Checking for absolute paths in evidence files..."
	@grep -r "/Users/daxxog/Desktop" bytecode/mapping/evidence/ | wc -l
	@echo "Absolute paths found (should be 0)"

verify-structure:
	@echo "Verifying clean directory structure..."
	@echo "Files in bytecode/mapping root:"
	@ls -la bytecode/mapping/ | grep -v "^d"
	@echo "Evidence subdirectories:"
	@ls bytecode/mapping/evidence/



# Complete enhanced verification
verify-all: verify-count verify-csv evidence-stats verify-naming check-paths verify-structure
	@echo "=== VERIFICATION SUMMARY ==="
	@echo "Use verify-evidence for agent-driven quality verification"
	@echo "Use verify-cleanliness for OG_vs_DEOB.md compliance"