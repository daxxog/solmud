# Generate class mappings to bytecode/mapping/
mapping:
	@echo "Generating class mappings..."
	@mkdir -p ../bytecode/mapping
	@cd .. && ./tools/classmapper/classmapper -mode csv | grep -E "^(deobfuscated_name|Node|NodeSub|RSApplet|client|Uncertain)" > bytecode/mapping/class_mapping.csv
	@cd .. && ./tools/classmapper/classmapper -mode json | sed -n '/^{/,/}/p' > bytecode/mapping/class_mapping.json
	@echo "Done! Files created in bytecode/mapping/"