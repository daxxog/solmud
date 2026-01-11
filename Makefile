.PHONY: help verify-evidence test-evidence-file validate-mapping-file check-duplicates validate-all clean-evidence status report

# Default target
help:
	@echo "Available targets:"
	@echo "  verify-evidence OG_CLASS=DEOB_CLASS  - Verify single mapping with all evidence checks"
	@echo "  test-evidence-file FILE=...          - Test all bash commands in an evidence file"
	@echo "  validate-mapping-file FILE=...       - Validate evidence file against OG_vs_DEOB.md checklist"
	@echo "  check-duplicates                     - Check for duplicate OG class mappings"
	@echo "  validate-all                         - Validate all evidence files"
	@echo "  clean-evidence                       - Remove non-compliant files"
	@echo "  status                              - Show verification progress"
	@echo "  report                              - Generate mapping report"
	@echo "  help                                - Show this help message"

# Verify single mapping with all evidence checks
verify-evidence:
	@if [ -z "$(OG_CLASS)" ] || [ -z "$(DEOB_CLASS)" ]; then \
		echo "Usage: make verify-evidence OG_CLASS=DEOB_CLASS"; \
		echo "Example: make verify-evidence client=client"; \
		exit 1; \
	fi
	@echo "Verifying mapping: $(DEOB_CLASS)_$(OG_CLASS)"
	@FILE="bytecode/mapping/evidence/verified/$(DEOB_CLASS)_$(OG_CLASS).md"; \
	if [ -f "$$FILE" ]; then \
		echo "Testing evidence file: $$FILE"; \
		$(MAKE) test-evidence-file FILE="$$FILE"; \
		$(MAKE) validate-mapping-file FILE="$$FILE"; \
	else \
		echo "Evidence file not found: $$FILE"; \
		exit 1; \
	fi

# Test all bash commands in an evidence file
test-evidence-file:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make test-evidence-file FILE=path/to/evidence.md"; \
		exit 1; \
	fi
	@echo "Testing bash commands in: $(FILE)"
	@awk '/^```$$/{flag=!flag; next} flag && /^$$/ {next} flag && /^[^$$]/ {print $$0}' "$(FILE)" | while read -r cmd; do \
		if [ -n "$$cmd" ]; then \
			echo "Testing: $$cmd"; \
			eval "$$cmd" || { echo "FAILED: $$cmd"; exit 1; }; \
		fi; \
	done
	@echo "All bash commands executed successfully"

# Validate evidence file against OG_vs_DEOB.md checklist
validate-mapping-file:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make validate-mapping-file FILE=path/to/evidence.md"; \
		exit 1; \
	fi
	@echo "Validating: $(FILE)"
	@# Check for absolute paths
	@if grep -q "/Users/daxxog/Desktop" "$(FILE)"; then \
		echo "❌ FAIL: Contains absolute paths"; \
		exit 1; \
	else \
		echo "✓ PASS: No absolute paths found"; \
	fi
	@# Check for bash code blocks
	@if grep -q '```bash' "$(FILE)"; then \
		echo "✓ PASS: Contains bash code blocks"; \
	else \
		echo "❌ FAIL: Missing bash code blocks"; \
		exit 1; \
	fi
	@# Check for overview section
	@if grep -q "## Overview" "$(FILE)"; then \
		echo "✓ PASS: Contains overview section"; \
	else \
		echo "❌ FAIL: Missing overview section"; \
		exit 1; \
	fi
	@# Check for mermaid diagrams (optional but recommended)
	@if grep -q "```mermaid" "$(FILE)"; then \
		echo "✓ PASS: Contains mermaid diagrams"; \
	else \
		echo "⚠️  WARN: No mermaid diagrams found"; \
	fi
	@echo "Validation completed"

# Check for duplicate OG class mappings
check-duplicates:
	@echo "Checking for duplicate OG class mappings..."
	@find bytecode/mapping/evidence -name "*.md" -type f | sed 's/.*\///' | sed 's/\.md$$//' | awk -F'_' '{print $$2}' | sort | uniq -c | sort -nr | while read count og_class; do \
		if [ "$$count" -gt 1 ]; then \
			echo "❌ DUPLICATE: $$og_class appears $$count times"; \
		fi; \
	done
	@echo "Duplicate check completed"

# Validate all evidence files
validate-all:
	@echo "Validating all evidence files..."
	@find bytecode/mapping/evidence -name "*.md" -type f | while read file; do \
		echo "Validating $$file"; \
		$(MAKE) validate-mapping-file FILE="$$file" || echo "VALIDATION FAILED: $$file"; \
	done
	@$(MAKE) check-duplicates
	@echo "All validation completed"

# Clean evidence directory of non-compliant files
clean-evidence:
	@echo "Removing non-compliant evidence files..."
	@find bytecode/mapping/evidence -name "*.md" -type f | while read file; do \
		if grep -q "/Users/daxxog/Desktop" "$$file" || ! grep -q "## Overview" "$$file"; then \
			echo "Removing non-compliant file: $$file"; \
			rm "$$file"; \
		fi; \
	done
	@echo "Cleanup completed"

# Show verification progress
status:
	@echo "=== Verification Status ==="
	@echo "Total evidence files:"
	@find bytecode/mapping/evidence -name "*.md" -type f | wc -l
	@echo "Files with bash blocks:"
	@find bytecode/mapping/evidence -name "*.md" -type f -exec grep -l '```bash' {} \; | wc -l
	@echo "Files with overview sections:"
	@find bytecode/mapping/evidence -name "*.md" -type f -exec grep -l '## Overview' {} \; | wc -l
	@echo "Files with mermaid diagrams:"
	@find bytecode/mapping/evidence -name "*.md" -type f -exec grep -l '```mermaid' {} \; | wc -l
	@echo "Files with absolute paths (should be 0):"
	@find bytecode/mapping/evidence -name "*.md" -type f -exec grep -l '/Users/daxxog/Desktop' {} \; | wc -l
	@echo "=== Classes to Map ==="
	@echo "DEOB classes available:"
	@ls -1 srcAllDummysRemoved/.javap_cache/*.javap.cache | xargs -n 1 basename | sed 's/\.javap\.cache$$//' | wc -l
	@echo "OG classes available:"
	@ls -1 bytecode/client/*.bytecode.txt | xargs -n 1 basename | sed 's/\.bytecode\.txt$$//' | wc -l
	@echo "Current mappings in CSV:"
	@if [ -f bytecode/mapping/class_mapping.csv ]; then wc -l bytecode/mapping/class_mapping.csv; else echo "No CSV found"; fi

# Generate comprehensive mapping report
report:
	@echo "=== Mapping Report ===" > bytecode/mapping/report.md
	@echo "Generated on: $$(date)" >> bytecode/mapping/report.md
	@echo "" >> bytecode/mapping/report.md
	@echo "## Summary" >> bytecode/mapping/report.md
	@echo "- Total DEOB classes: $$(ls -1 srcAllDummysRemoved/.javap_cache/*.javap.cache | xargs -n 1 basename | sed 's/\.javap\.cache$$//' | wc -l | tr -d ' ')" >> bytecode/mapping/report.md
	@echo "- Total OG classes: $$(ls -1 bytecode/client/*.bytecode.txt | xargs -n 1 basename | sed 's/\.bytecode\.txt$$//' | wc -l | tr -d ' ')" >> bytecode/mapping/report.md
	@echo "- Evidence files created: $$(find bytecode/mapping/evidence -name "*.md" -type f | wc -l | tr -d ' ')" >> bytecode/mapping/report.md
	@echo "" >> bytecode/mapping/report.md
	@echo "## Mapped Classes" >> bytecode/mapping/report.md
	@find bytecode/mapping/evidence -name "*.md" -type f | sed 's/.*\///' | sed 's/\.md$$//' | sort | while read mapping; do \
		echo "- $$mapping" >> bytecode/mapping/report.md; \
	done
	@echo "" >> bytecode/mapping/report.md
	@echo "## Validation Issues" >> bytecode/mapping/report.md
	@find bytecode/mapping/evidence -name "*.md" -type f -exec grep -l '/Users/daxxog/Desktop' {} \; | while read file; do \
		echo "- Absolute paths in: $$file" >> bytecode/mapping/report.md; \
	done
	@echo "Report generated: bytecode/mapping/report.md"