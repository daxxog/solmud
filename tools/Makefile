# Generate class mappings to bytecode/mapping/
mapping:
	@echo "Generating class mappings..." >&2
	@mkdir -p ../bytecode/mapping
	@cd .. && ./tools/classmapper/classmapper -mode csv > bytecode/mapping/class_mapping.csv
	@cd .. && ./tools/classmapper/classmapper -mode json > bytecode/mapping/class_mapping.json
	@echo "Done! Files created in bytecode/mapping/" >&2