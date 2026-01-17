#!/bin/bash
# OG_vs_DEOB.md Compliance Verification Script

echo "=== OG_vs_DEOB.md COMPLIANCE VERIFICATION ==="
echo "Running automated checks for directory cleanliness..."

# Directory structure check
tree bytecode/mapping | tail -3
echo ""

# Check for forbidden markdown files at root level
echo "Checking for forbidden markdown files at root level..."
if ls bytecode/mapping/*.md >/dev/null 2>&1; then
    echo "❌ VIOLATION: Markdown files found at root level"
    ls bytecode/mapping/*.md
    echo "Remove these files to comply with OG_vs_DEOB.md"
    exit 1
else
    echo "✅ No markdown spam at root level"
fi
echo ""

# Check for SUMMARY/REPORT files
echo "Checking for SUMMARY/REPORT files..."
forbidden_files=$(ls bytecode/mapping | grep -E "(SUMMARY|REPORT)$" || true)
if [ -n "$forbidden_files" ]; then
    echo "❌ VIOLATION: Forbidden documentation files detected"
    echo "$forbidden_files"
    echo "Remove these to comply with OG_vs_DEOB.md"
    exit 1
else
    echo "✅ No forbidden documentation files"
fi
echo ""

# Count evidence files
echo "Counting evidence files..."
evidence_count=$(find bytecode/mapping/evidence/verified -name "*.md" | wc -l)
if [ "$evidence_count" -eq 73 ]; then
    echo "✅ Correct evidence file count: $evidence_count"
else
    echo "❌ VIOLATION: Expected 73 evidence files, found $evidence_count"
    exit 1
fi
echo ""

# Check file naming conventions
echo "Checking file naming conventions..."
bad_names=$(ls bytecode/mapping/evidence/verified/ | grep -v "_.*\.md$" | wc -l)
if [ "$bad_names" -gt 0 ]; then
    echo "❌ VIOLATION: Files with incorrect naming convention"
    ls bytecode/mapping/evidence/verified/ | grep -v "_.*\.md$"
    exit 1
else
    echo "✅ All files follow DEOB_OG.md naming convention"
fi
echo ""

echo "🎉 OG_vs_DEOB.md COMPLIANCE VERIFICATION PASSED"