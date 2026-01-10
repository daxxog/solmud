# Forensic Evidence Template: <DEOBFILENAME> → <OBFUSCATED_NAME>

## Evidence Summary
- **Deobfuscated Class**: `<DEOBFILENAME>`
- **Obfuscated Class**: `<OBFUSCATED_NAME>`
- **Confidence Score**: `<score>/100`
- **Evidence Type**: `<Anchor|Inheritance|Signature|Behavioral>`
- **Verification Status**: `<Verified|Pending|Disputed>`

## Primary Forensic Evidence

### 1. Structural Fingerprints
- **Class Modifiers**: `<public|private|final>`
- **Superclass**: `<superclass_match>`
- **Interfaces**: `<interface_list>`
- **Field Count**: `<deob>/<obf>`
- **Method Count**: `<deob>/<obf>`

### 2. Behavioral Analysis
- **Constructor Pattern**: `<signature_analysis>`
- **Magic Constants**: `<constants_found>`
- **Cross-References**: `<external_dependencies>`
- **Unique Operations**: `<distinctive_bytecode_patterns>`

### 3. Forensic Technique Evidence

#### Anchor Evidence (if applicable)
- **File Loading Pattern**: `<e.g., loads "seq.dat">`
- **Implementation Match**: `<exact_behavioral_match>`
- **Method Signature**: `<critical_identifying_method>`

#### Inheritance Chain Evidence (if applicable)
- **Parent Class**: `<confirmed_parent_mapping>`
- **Sibling Classes**: `<related_mappings>`
- **Chain Validation**: `<hierarchy_consistency>`

#### Signature Evidence (if applicable)
- **Field Type Patterns**: `<type_distribution_analysis>`
- **Method Signature Similarity**: `<signature_correlation_score>`
- **Access Modifier Match**: `<modifier_consistency>`

#### Behavioral Evidence (if applicable)
- **Bytecode Instruction Patterns**: `<unique_instruction_sequences>`
- **Method Call Graph**: `<caller_callee_relationships>`
- **State Manipulation**: `<field_access_patterns>`

## Detailed Analysis

### Critical Evidence Points
1. **[Evidence Point 1]**: `<detailed_description_with_bytecode_references>`
2. **[Evidence Point 2]**: `<detailed_description_with_bytecode_references>`
3. **[Evidence Point 3]**: `<detailed_description_with_bytecode_references>`

### Cross-Reference Validation
- **Dependencies**: `<referenced_classes_and_their_mappings>`
- **References From**: `<classes_that_reference_this_class>`
- **Validation Loop**: `<how_related_mappings_confirm_this_match>`

### Counter-Evidence Analysis
- **Potential Conflicts**: `<any_discrepancies_or_alternatives>`
- **Risk Assessment**: `<confidence_limiting_factors>`
- **Alternative Hypotheses**: `<other_possible_matches>`

## Forensic Methodology

### Detection Method Used
- **Primary Method**: `<anchor|inheritance|signature|behavioral|hybrid>`
- **Supporting Methods**: `<secondary_techniques>`
- **Validation Techniques**: `<cross_reference_patterns>`

### Evidence Strength Assessment
- **Structural Match**: `<percentage>/25`
- **Behavioral Match**: `<percentage>/25`
- **Cross-Reference Match**: `<percentage>/25`
- **Unique Identifiers**: `<percentage>/25`

## Verification History
- **Initial Match**: `<date_and_method>`
- **Cross-Reference Validation**: `<date_and_results>`
- **Manual Review**: `<date_and_findings>`
- **Final Confirmation**: `<date_and_confidence_level>`

## Sources and References
- **Deobfuscated Source**: `srcAllDummysRemoved/src/<DEOBFILENAME>.java`
- **Obfuscated Bytecode**: `bytecode/client/<OBFUSCATED_NAME>.bytecode.txt`
- **Related Evidence**: `<links_to_other_evidence_files>`
- **Analysis Tools**: `<classmapper_version_and_settings>`