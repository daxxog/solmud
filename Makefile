# Show available targets
help:
	@echo "Available targets:"
	@echo "  mapping  - Generate class mappings (CSV and JSON) in bytecode/mapping/"
	@echo "  help     - Show this help message"

# Generate class mappings to bytecode/mapping/
mapping:
	@echo "Generating class mappings..." >&2
	@mkdir -p bytecode/mapping
	@cd tools/classmapper && go run . -mode csv -obf ../../bytecode/client > ../../bytecode/mapping/class_mapping.csv
	@cd tools/classmapper && go run . -mode json -obf ../../bytecode/client > ../../bytecode/mapping/class_mapping.json
	@echo "Done! Files created in bytecode/mapping/" >&2

.PHONY: help mapping