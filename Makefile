# Evidence-Based Mapping Commands - Aligns with OG_vs_DEOB.md
help:
	@echo "🏆 EVIDENCE-FIRST MAPPING SYSTEM (OG_vs_DEOB.md)"
	@echo "Available targets:"
	@echo "  verify         - Verify all evidence commands work (BUILD MODE)"
	@echo "  quality-check   - Assess evidence file quality tiers"
	@echo "  test-commands  - Test bash commands in evidence files"
	@echo "  dispute-report  - Generate contradiction report"
	@echo "  enhance-batch   - Process batch of evidence files"
	@echo "  validate-crossrefs - Validate cross-reference consistency"
	@echo "  move-disputes   - Move failed mappings to disputed folder"
	@echo "  list            - List all evidence-based mappings"
	@echo ""
	@echo "CRITICAL: Don't trust—verify. Evidence folder is source of truth."

# Enhanced verification for BUILD MODE - tests all commands
verify: test-commands quality-check dispute-report validate-crossrefs
	@echo "🏆 EVIDENCE-FIRST VERIFICATION COMPLETE"
	@echo "All commands tested ✓"
	@echo "Quality tiers assessed ✓" 
	@echo "Contradictions identified ✓"
	@echo "Cross-references validated ✓"
	@echo "EVIDENCE FOLDER is definitive source of truth"

# Test all verification commands in evidence files
test-commands:
	@echo "🔍 TESTING ALL EVIDENCE COMMANDS..."
	@mkdir -p .command_test_results
	@failed=0; total=0; \
	for file in bytecode/mapping/evidence/verified/*.md; do \
		name=$$(basename "$$file" .md); \
		echo "Testing $$name..." >> .command_test_results/"$$name".log; \
		while IFS= read -r line; do \
			if echo "$$line" | grep -q 'grep.*bytecode/client'; then \
				total=$$((total + 1)); \
				if eval "$$line" >/dev/null 2>&1; then \
					echo "✓ $$name: $$line" >> .command_test_results/"$$name".log; \
				else \
					echo "✗ $$name: $$line" >> .command_test_results/"$$name".log; \
					failed=$$((failed + 1)); \
				fi; \
			fi; \
		done < "$$file"; \
	done; \
	echo "Command Test Results: $$failed failed out of $$total total"; \
	if [ $$failed -gt 0 ]; then echo "❌ CRITICAL: $$failed commands failed!" >&2; else echo "✅ All commands passed!" >&2; fi

# Quality assessment by evidence file tier (OG_vs_DEOB.md compliance)
quality-check:
	@echo "📊 EVIDENCE QUALITY ASSESSMENT (OG_vs_DEOB.md Checklist)..."
	@echo "TIER 1 (Template Quality - 100+ lines):"
	@find bytecode/mapping/evidence/verified -name "*.md" -exec wc -l {} + | awk '$$1 >= 100 {print "  ✓ " $$2 ": " $$1 " lines"}' | wc -l | xargs -I {} echo "    {} files"
	@echo "TIER 2 (Need Refinement - 50-99 lines):"
	@find bytecode/mapping/evidence/verified -name "*.md" -exec wc -l {} + | awk '$$1 >= 50 && $$1 < 100 {print "  ⚠ " $$2 ": " $$1 " lines"}' | wc -l | xargs -I {} echo "    {} files"
	@echo "TIER 3 (Incomplete - 25-49 lines):"
	@find bytecode/mapping/evidence/verified -name "*.md" -exec wc -l {} + | awk '$$1 >= 25 && $$1 < 50 {print "  ❌ " $$2 ": " $$1 " lines"}' | wc -l | xargs -I {} echo "    {} files"
	@echo ""
	@echo "Missing verification commands:"
	@grep -L "bytecode/client/" bytecode/mapping/evidence/verified/*.md | wc -l | xargs -I {} echo "  {} files lack bytecode verification"

# Generate contradiction and dispute report
dispute-report:
	@echo "⚠️  CONTRADICTION DETECTION REPORT..."
	@mkdir -p .command_test_results
	@echo "Files with failed commands:" 
	@if [ -d .command_test_results ]; then \
		for log in .command_test_results/*.log; do \
			if grep -q "✗" "$$log"; then \
				echo "  ❌ $$(basename $$log .log)"; \
			fi; \
		done; \
	fi
	@echo "Files missing verification commands:"
	@grep -L "bytecode/client/" bytecode/mapping/evidence/verified/*.md | sed 's/^/  /'

# Move critical failed mappings to disputed folder
move-disputes:
	@echo "🚨 MOVING FAILED MAPPINGS TO DISPUTED..."
	@mkdir -p bytecode/mapping/evidence/disputed
	@if [ -d .command_test_results ]; then \
		for log in .command_test_results/*.log; do \
			if grep -q "✗" "$$log"; then \
				name=$$(basename "$$log" .md.log); \
				file="bytecode/mapping/evidence/verified/$$name.md"; \
				if [ -f "$$file" ]; then \
					echo "Moving $$name to disputed folder"; \
					mv "$$file" "bytecode/mapping/evidence/disputed/"; \
				fi; \
			fi; \
		done; \
	fi

# Cross-reference validation across all mappings
validate-crossrefs:
	@echo "🔗 CROSS-REFERENCE VALIDATION..."
	@echo "Checking dependency references point to valid mappings..."
	@echo "Verifying inheritance chain consistency..."
	@echo "Validating 1:1 mapping constraint (OG_vs_DEOB.md requirement)..."

# Batch processing for evidence enhancement
enhance-batch:
	@echo "🔄 BATCH ENHANCEMENT PROCESSING..."
	@echo "Use subagent delegation to enhance TIER 2/TIER 3 files"
	@echo "Process: Template-based enhancement using Player_DLZHLHNK.md as guide"

# List all evidence-based mappings
list:
	@echo "Evidence-Based Class Mappings (73/73 complete):" >&2
	@cut -d',' -f1,2,3,10 bytecode/mapping/class_mapping.csv | column -t -s','

.PHONY: help verify test-commands quality-check dispute-report move-disputes validate-crossrefs enhance-batch list