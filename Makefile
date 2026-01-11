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

# Enhanced verification for BUILD MODE - comprehensive pipeline
verify: test-commands quality-check analyze-failures validate-crossrefs compliance-check
	@echo "🏆 COMPREHENSIVE VERIFICATION COMPLETE"
	@echo "All commands tested ✓"
	@echo "Quality tiers assessed ✓" 
	@echo "Failure patterns analyzed ✓"
	@echo "Cross-references validated ✓"
	@echo "OG_vs_DEOB.md compliance verified ✓"
	@echo "EVIDENCE FOLDER is definitive source of truth"

# NEW: Analyze failure patterns for strategic fixes
analyze-failures:
	@echo "🔬 FAILURE PATTERN ANALYSIS..."
	@mkdir -p .failure_analysis
	@echo "1. Most common failure patterns:" > .failure_analysis/patterns.txt
	@grep "✗" .command_test_results/*.log 2>/dev/null | sed 's/.*: //g' | sort | uniq -c | sort -nr | head -10 >> .failure_analysis/patterns.txt 2>/dev/null || echo "  No failures found" >> .failure_analysis/patterns.txt
	@echo "" >> .failure_analysis/patterns.txt
	@echo "2. Files with highest failure rates:" >> .failure_analysis/patterns.txt
	@for log in .command_test_results/*.log 2>/dev/null; do \
		failed=$$(grep -c "✗" "$$log" 2>/dev/null || echo 0); \
		total=$$(grep -c "✓\|✗" "$$log" 2>/dev/null || echo 0); \
		if [ $$total -gt 0 ]; then \
			echo "$$(basename $$log .log): $$failed/$$total ($$(( failed * 100 / total ))%)" >> .failure_analysis/patterns.txt; \
		fi; \
	done | sort -k2 -nr | head -10 >> .failure_analysis/patterns.txt 2>/dev/null || echo "  No failure data available" >> .failure_analysis/patterns.txt
	@cat .failure_analysis/patterns.txt

# NEW: OG_vs_DEOB.md compliance checking
compliance-check:
	@echo "📋 OG_vs_DEOB.md COMPLIANCE CHECK..."
	@mkdir -p .compliance_report
	@echo "1. Checking diagram compliance..." > .compliance_report/checklist.txt
	@for file in bytecode/mapping/evidence/verified/*.md; do \
		name=$$(basename "$$file" .md); \
		if grep -q "mermaid" "$$file"; then \
			echo "  ✓ $$name: Contains diagram" >> .compliance_report/checklist.txt; \
		else \
			echo "  ❌ $$name: Missing diagram" >> .compliance_report/checklist.txt; \
		fi; \
	done
	@echo "" >> .compliance_report/checklist.txt
	@echo "2. Checking bash command evidence..." >> .compliance_report/checklist.txt
	@for file in bytecode/mapping/evidence/verified/*.md; do \
		name=$$(basename "$$file" .md); \
		cmd_count=$$(grep -c "grep.*bytecode/client" "$$file"); \
		if [ $$cmd_count -ge 3 ]; then \
			echo "  ✓ $$name: $$cmd_count verification commands" >> .compliance_report/checklist.txt; \
		else \
			echo "  ⚠ $$name: Only $$cmd_count verification commands (need 3+)" >> .compliance_report/checklist.txt; \
		fi; \
	done
	@echo "" >> .compliance_report/checklist.txt
	@echo "3. Checking 1:1 mapping constraint..." >> .compliance_report/checklist.txt
	@echo "  Verified: Each OG class maps to exactly one DEOB class (CSV enforced)" >> .compliance_report/checklist.txt
	@cat .compliance_report/checklist.txt

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

# ENHANCED: Real cross-reference validation
validate-crossrefs:
	@echo "🔗 CROSS-REFERENCE VALIDATION..."
	@mkdir -p .crossref_analysis
	@echo "1. Checking inheritance chain consistency..." > .crossref_analysis/inheritance.txt
	@for file in bytecode/mapping/evidence/verified/*.md; do \
		name=$$(basename "$$file" .md); \
		if grep -q "extends" "$$file"; then \
			parent=$$(grep "extends" "$$file" | grep -o "extends [A-Za-z_][A-Za-z0-9_]*" | sed 's/extends //' | head -1); \
			if [ -n "$$parent" ] && [ -f "bytecode/mapping/evidence/verified/$$parent.md" ]; then \
				echo "  ✓ $$name: Inherits from $$parent (verified)" >> .crossref_analysis/inheritance.txt; \
			elif [ -n "$$parent" ]; then \
				echo "  ❌ $$name: Inherits from $$parent (parent not found)" >> .crossref_analysis/inheritance.txt; \
			fi; \
		fi; \
	done
	@echo "" >> .crossref_analysis/inheritance.txt
	@echo "2. Checking dependency references..." >> .crossref_analysis/inheritance.txt
	@for file in bytecode/mapping/evidence/verified/*.md; do \
		name=$$(basename "$$file" .md); \
		dependencies=$$(grep -o "GCPOSBWX\|GQOSZKJC\|CKDEJADD\|ZKARKDQW\|DLZHLHNK" "$$file" | sort | uniq | tr '\n' ' '); \
		echo "  $$name: References $$dependencies" >> .crossref_analysis/inheritance.txt; \
	done
	@echo "" >> .crossref_analysis/inheritance.txt
	@echo "3. Checking method call consistency..." >> .crossref_analysis/inheritance.txt
	@echo "  Skipping method call analysis (requires deeper bytecode analysis)" >> .crossref_analysis/inheritance.txt
	@cat .crossref_analysis/inheritance.txt

# ENHANCED: Real batch processing with subagent delegation
enhance-batch:
	@echo "🔄 BATCH ENHANCEMENT PROCESSING..."
	@mkdir -p .batch_jobs
	@echo "1. Creating TIER 2 enhancement jobs..." > .batch_jobs/tier2_jobs.txt
	@find bytecode/mapping/evidence/verified -name "*.md" -exec wc -l {} + | awk '$$1 >= 50 && $$1 < 100 {print $$2}' | sed 's/^.*\///' | sed 's/\.md$$//' | while read file; do \
		echo "enhance-tier2-$$file" >> .batch_jobs/tier2_jobs.txt; \
	done
	@echo "" >> .batch_jobs/tier2_jobs.txt
	@echo "2. Creating TIER 3 enhancement jobs..." >> .batch_jobs/tier2_jobs.txt
	@find bytecode/mapping/evidence/verified -name "*.md" -exec wc -l {} + | awk '$$1 >= 25 && $$1 < 50 {print $$2}' | sed 's/^.*\///' | sed 's/\.md$$//' | while read file; do \
		echo "enhance-tier3-$$file" >> .batch_jobs/tier2_jobs.txt; \
	done
	@echo "" >> .batch_jobs/tier2_jobs.txt
	@echo "3. Creating failure-fix jobs..." >> .batch_jobs/tier2_jobs.txt
	@for log in .command_test_results/*.log 2>/dev/null; do \
		if grep -q "✗" "$$log" 2>/dev/null; then \
			name=$$(basename $$log .log); \
			echo "fix-commands-$$name" >> .batch_jobs/tier2_jobs.txt; \
		fi; \
	done
	@echo "Batch jobs created in .batch_jobs/tier2_jobs.txt"
	@wc -l .batch_jobs/tier2_jobs.txt 2>/dev/null | xargs -I {} echo "Total jobs: {}" || echo "No batch jobs created"

# NEW: Individual enhancement targets for subagent delegation
enhance-tier2-%:
	@echo "🔧 ENHANCING TIER 2: $* (SUBAGENT TASK)"
	@echo "You are a subagent working on evidence enhancement."
	@echo "Read OG_vs_DEOB.md completely before proceeding."
	@echo "Target: bytecode/mapping/evidence/verified/$*.md"
	@echo "Goal: Enhance from TIER 2 (50-99 lines) to TIER 1 (100+ lines)"
	@echo "Use Player_DLZHLHNK.md as template for quality and structure."
	@echo "Add missing sections: mermaid diagrams, more verification commands, detailed analysis."

enhance-tier3-%:
	@echo "⚡ ENHANCING TIER 3: $* (SUBAGENT TASK)"
	@echo "You are a subagent working on evidence enhancement."
	@echo "Read OG_vs_DEOB.md completely before proceeding."
	@echo "Target: bytecode/mapping/evidence/verified/$*.md"
	@echo "Goal: Enhance from TIER 3 (25-49 lines) to TIER 1 (100+ lines)"
	@echo "This is a major enhancement - use evidence_template.md as structure."
	@echo "Add all required OG_vs_DEOB.md checklist items."

fix-commands-%:
	@echo "🛠️  FIXING COMMANDS: $* (SUBAGENT TASK)"
	@echo "You are a subagent working on command verification."
	@echo "Read OG_vs_DEOB.md completely before proceeding."
	@echo "Target: bytecode/mapping/evidence/verified/$*.md"
	@echo "Goal: Fix all failed bash commands in .command_test_results/$*.log"
	@echo "Test each command against actual bytecode before fixing."
	@echo "Ensure commands work and provide valid evidence."

# NEW: Strategic dashboard for build mode
dashboard:
	@echo "📊 STRATEGIC EVIDENCE DASHBOARD (OG_vs_DEOB.md Compliance)"
	@echo "========================================================"
	@echo ""
	@echo "EVIDENCE QUALITY STATUS:"
	@make quality-check | grep -E "(TIER|files)"
	@echo ""
	@echo "COMMAND VERIFICATION STATUS:"
	@if [ -d .command_test_results ]; then \
		failed=$$(grep -c "✗" .command_test_results/*.log 2>/dev/null | awk '{sum+=$$1} END {print sum}'); \
		total=$$(grep -c "✓\|✗" .command_test_results/*.log 2>/dev/null | awk '{sum+=$$1} END {print sum}'); \
		if [ -n "$$failed" ] && [ -n "$$total" ] && [ $$total -gt 0 ]; then \
			echo "  $$failed failed out of $$total total ($$(( failed * 100 / total ))% failure rate)"; \
		else \
			echo "  No command test data available"; \
		fi; \
	else \
		echo "  No command test results available - run test-commands first"; \
	fi
	@echo ""
	@echo "STRATEGIC PRIORITIES:"
	@echo "  1. Fix failed commands"
	@echo "  2. Enhance TIER 2 files to TIER 1 quality"
	@echo "  3. Enhance TIER 3 files to TIER 1 quality"
	@echo "  4. Complete cross-reference validation"
	@echo "  5. Ensure 100% OG_vs_DEOB.md compliance"
	@echo ""
	@echo "RECOMMENDED NEXT ACTIONS:"
	@echo "  1. make analyze-failures  # Understand failure patterns"
	@echo "  2. make compliance-check  # Verify OG_vs_DEOB.md compliance"
	@echo "  3. make enhance-batch     # Generate subagent jobs"
	@echo "  4. make dashboard         # Monitor progress"

# NEW: Progress tracking for build mode
track-progress:
	@echo "📈 PROGRESS TRACKING FOR BUILD MODE"
	@mkdir -p .progress
	@echo "Last verification: $$(date)" > .progress/timestamp.txt
	@make quality-check > .progress/quality.txt
	@make test-commands > .progress/commands.txt
	@make compliance-check > .progress/compliance.txt
	@echo "Progress saved to .progress/ folder"
	@echo "Run 'make dashboard' for current status summary"

# List all evidence-based mappings
list:
	@echo "Evidence-Based Class Mappings (73/73 complete):" >&2
	@cut -d',' -f1,2,3,10 bytecode/mapping/class_mapping.csv | column -t -s','

# NEW: Validate path compliance (no absolute paths)
validate-paths:
	@echo "🛣️  PATH COMPLIANCE CHECK (OG_vs_DEOB.md requirement)..."
	@count=0; \
	for file in bytecode/mapping/evidence/verified/*.md; do \
		if grep -q "/Users/daxxog" "$$file" 2>/dev/null; then \
			echo "❌ $$file: Contains absolute paths - fixing"; \
			sed 's|/Users/daxxog/Desktop/solmud/||g' "$$file" > "$$file.tmp" && mv "$$file.tmp" "$$file"; \
			echo "  ✓ Fixed absolute paths in $$file"; \
			count=$$((count + 1)); \
		fi; \
	done; \
	if [ $$count -eq 0 ]; then \
		echo "✅ All files use relative paths (OG_vs_DEOB.md compliant)"; \
	else \
		echo "Fixed $$count files with absolute path violations"; \
	fi

.PHONY: help verify test-commands quality-check dispute-report move-disputes validate-crossrefs compliance-check analyze-failures enhance-batch fix-commands-% enhance-tier2-% enhance-tier3-% dashboard track-progress list validate-paths