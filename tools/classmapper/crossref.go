package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ReferenceType represents different types of cross-references
type ReferenceType int

const (
	RefMethodCall    ReferenceType = iota + 1 // Highest weight
	RefInstantiation                          // High weight
	RefFieldAccess                            // Medium weight
	RefStaticCall                             // Lower weight
	RefArrayAccess                            // Array access patterns
)

// String representation for debugging
func (rt ReferenceType) String() string {
	switch rt {
	case RefMethodCall:
		return "MethodCall"
	case RefInstantiation:
		return "Instantiation"
	case RefFieldAccess:
		return "FieldAccess"
	case RefStaticCall:
		return "StaticCall"
	case RefArrayAccess:
		return "ArrayAccess"
	default:
		return "Unknown"
	}
}

// CrossReference represents a single cross-reference between classes
type CrossReference struct {
	Type          ReferenceType `json:"type"`
	TargetClass   string        `json:"target_class"`
	TargetMember  string        `json:"target_member"`
	Context       string        `json:"context,omitempty"`
	LineNumber    int           `json:"line_number"`
	Frequency     int           `json:"frequency"`
	IsChained     bool          `json:"is_chained,omitempty"`
	IsArrayAccess bool          `json:"is_array_access,omitempty"`
}

// ClassCrossReferences contains all cross-references for a class
type ClassCrossReferences struct {
	// References organized by source method for context
	MethodReferences map[string][]CrossReference `json:"method_references"`

	// Aggregate statistics
	UniqueTargets     map[string]int        `json:"unique_targets"`     // class -> frequency
	ReferenceCounts   map[ReferenceType]int `json:"reference_counts"`   // type -> count
	TotalReferences   int                   `json:"total_references"`   // sum of all references
	ChainedReferences int                   `json:"chained_references"` // complex patterns
	ArrayReferences   int                   `json:"array_references"`   // array access patterns

	// Normalized data for comparison
	NormalizedTargets map[string]float64 `json:"normalized_targets,omitempty"`
	ReferenceVector   map[string]float64 `json:"reference_vector,omitempty"`
}

// CrossReferenceParser handles parsing of cross-references from Java source files
type CrossReferenceParser struct {
	// Regex patterns optimized for Jad decompiler output
	regex_method_call   *regexp.Regexp // obj.method(args)
	regex_field_access  *regexp.Regexp // obj.field
	regex_instantiation *regexp.Regexp // new Class(args)
	regex_static_call   *regexp.Regexp // Class.method(args)
	regex_array_access  *regexp.Regexp // array[index]
	regex_chain_pattern *regexp.Regexp // obj1.obj2.method()
	regex_method_decl   *regexp.Regexp // method declarations

	// Class registry for validation
	projectClasses map[string]bool
	anchorMappings map[string]string

	// Caching and optimization
	patternCache   map[string][]CrossReference
	exclusionCache map[string]bool
}

// Java built-in classes and packages to exclude from analysis
var BUILTIN_EXCLUSIONS = map[string]bool{
	// java.lang.*
	"String": true, "Object": true, "Class": true, "Math": true,
	"Integer": true, "Double": true, "Boolean": true, "System": true,
	"Thread": true, "Exception": true, "Runtime": true,
	"Throwable": true, "Error": true, "NullPointerException": true,

	// java.io.*
	"IOException": true, "EOFException": true, "DataInputStream": true,
	"BufferedReader": true, "FileReader": true, "InputStream": true,
	"OutputStream": true, "PrintStream": true,

	// java.net.*
	"URL": true, "URLConnection": true, "Socket": true,
	"InetAddress": true, "HttpURLConnection": true,

	// java.awt.*
	"Component": true, "Container": true, "Graphics": true,
	"Color": true, "Font": true, "Image": true,

	// java.applet.*
	"AppletContext": true, "Applet": true,

	// java.math.*
	"BigInteger": true,

	// java.util.zip.*
	"CRC32": true, "GZIPInputStream": true,

	// Primitives and wrappers (though unlikely to appear)
	"int": true, "boolean": true, "double": true, "long": true,
	"void": true, "byte": true, "char": true, "short": true,
	"float": true,

	// Arrays and generic types
	"Arrays": true, "Collections": true, "List": true, "ArrayList": true,
	"HashMap": true, "HashSet": true, "Map": true, "Set": true,
}

// NewCrossReferenceParser creates a new cross-reference parser
func NewCrossReferenceParser(projectClasses []string, anchors map[string]string) *CrossReferenceParser {
	parser := &CrossReferenceParser{
		// Enhanced regex for Jad output patterns (obfuscated names: aClass123, method456)
		regex_method_call:   regexp.MustCompile(`([a-zA-Z_$][a-zA-Z0-9_$]*)\.([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(`),
		regex_field_access:  regexp.MustCompile(`([a-zA-Z_$][a-zA-Z0-9_$]*)\.([a-zA-Z_$][a-zA-Z0-9_$]*)(?:\s*[;=]|\s*\.)`),
		regex_instantiation: regexp.MustCompile(`new\s+([A-Z][a-zA-Z0-9_$]*)\s*\(`),
		regex_static_call:   regexp.MustCompile(`([A-Z][a-zA-Z0-9_$]*)\.([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(`),
		regex_array_access:  regexp.MustCompile(`([a-zA-Z_$][a-zA-Z0-9_$]*)\[(.*?)\]`),
		regex_chain_pattern: regexp.MustCompile(`(?:[a-zA-Z_$][a-zA-Z0-9_$]*\.){2,}[a-zA-Z_$][a-zA-Z0-9_$]*\s*\(`),
		regex_method_decl:   regexp.MustCompile(`^\s*(?:private|public|protected)?\s*(?:static\s+)?(?:final\s+)?[a-zA-Z_$][a-zA-Z0-9_$<>\[\]]*\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(`),

		projectClasses: makeProjectClassMap(projectClasses),
		anchorMappings: anchors,
		patternCache:   make(map[string][]CrossReference),
		exclusionCache: make(map[string]bool),
	}

	return parser
}

// makeProjectClassMap converts string slice to map for fast lookup
func makeProjectClassMap(classes []string) map[string]bool {
	classMap := make(map[string]bool)
	for _, class := range classes {
		// Handle both simple names and full qualified names
		classMap[class] = true
		// Also add without package prefix if present
		if lastDot := strings.LastIndex(class, "."); lastDot >= 0 {
			classMap[class[lastDot+1:]] = true
		}
	}
	return classMap
}

// isBuiltinClass checks if a class should be excluded from analysis
func (p *CrossReferenceParser) isBuiltinClass(className string) bool {
	// Check cache first
	if excluded, exists := p.exclusionCache[className]; exists {
		return excluded
	}

	// Check direct match
	if BUILTIN_EXCLUSIONS[className] {
		p.exclusionCache[className] = true
		return true
	}

	// Check package prefixes
	excluded := strings.HasPrefix(className, "java.") ||
		strings.HasPrefix(className, "javax.") ||
		strings.HasPrefix(className, "sun.") ||
		strings.HasPrefix(className, "com.sun.") ||
		strings.HasPrefix(className, "org.")

	p.exclusionCache[className] = excluded
	return excluded
}

// isProjectClass checks if a class is part of the project
func (p *CrossReferenceParser) isProjectClass(className string) bool {
	return p.projectClasses[className]
}

// shouldIncludeReference determines if a reference should be included in analysis
func (p *CrossReferenceParser) shouldIncludeReference(ref CrossReference) bool {
	// Exclude built-in classes
	if p.isBuiltinClass(ref.TargetClass) {
		return false
	}

	// Include if it's a known project class
	if p.isProjectClass(ref.TargetClass) {
		return true
	}

	// Include if it's a mapped anchor class
	if _, isMapped := p.anchorMappings[ref.TargetClass]; isMapped {
		return true
	}

	// For now, be conservative and only include known classes
	// This can be relaxed later for broader analysis
	return false
}

// ParseFile extracts cross-references from a Java source file
func (p *CrossReferenceParser) ParseFile(filePath string) (*ClassCrossReferences, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(file), "\n")
	refs := &ClassCrossReferences{
		MethodReferences: make(map[string][]CrossReference),
		UniqueTargets:    make(map[string]int),
		ReferenceCounts:  make(map[ReferenceType]int),
	}

	var currentMethod string

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Track current method for context
		if methodMatch := p.extractMethodSignature(line); methodMatch != "" {
			currentMethod = methodMatch
			continue
		}

		// Extract all reference types from this line
		lineRefs := p.extractReferencesFromLine(line, lineNum)

		// Filter and categorize references
		for _, ref := range lineRefs {
			if p.shouldIncludeReference(ref) {
				refs.MethodReferences[currentMethod] = append(refs.MethodReferences[currentMethod], ref)
				refs.aggregateReference(ref)
			}
		}
	}

	// Normalize frequencies for comparison
	refs.normalizeFrequencies()

	return refs, nil
}

// extractMethodSignature extracts method name from method declarations
func (p *CrossReferenceParser) extractMethodSignature(line string) string {
	matches := p.regex_method_decl.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractReferencesFromLine extracts all types of references from a single line
func (p *CrossReferenceParser) extractReferencesFromLine(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	// Check for method chains first (highest priority)
	if p.regex_chain_pattern.MatchString(line) {
		chainRefs := p.parseMethodChains(line, lineNum)
		refs = append(refs, chainRefs...)
	}

	// Extract individual reference types
	methodRefs := p.parseMethodCalls(line, lineNum)
	refs = append(refs, methodRefs...)

	fieldRefs := p.parseFieldAccesses(line, lineNum)
	refs = append(refs, fieldRefs...)

	instRefs := p.parseInstantiations(line, lineNum)
	refs = append(refs, instRefs...)

	staticRefs := p.parseStaticCalls(line, lineNum)
	refs = append(refs, staticRefs...)

	arrayRefs := p.parseArrayAccesses(line, lineNum)
	refs = append(refs, arrayRefs...)

	return refs
}

// parseMethodChains handles complex method chains like: aSpotAnim_1568.aAnimation_407.method258(anInt1569)
func (p *CrossReferenceParser) parseMethodChains(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	// Split by dots to analyze chain segments
	segments := strings.Split(line, ".")

	for i := 0; i < len(segments)-1; i++ {
		currentSegment := strings.TrimSpace(segments[i])
		nextSegment := strings.TrimSpace(segments[i+1])

		// Extract method name from next segment
		methodName := p.extractMethodName(nextSegment)
		if methodName == "" {
			continue
		}

		// Check if current segment is a valid class reference
		if p.isValidClassReference(currentSegment) {
			refs = append(refs, CrossReference{
				Type:         RefMethodCall,
				TargetClass:  currentSegment,
				TargetMember: methodName,
				Context:      line,
				LineNumber:   lineNum,
				IsChained:    true,
			})
		}
	}

	return refs
}

// parseMethodCalls extracts regular method calls like obj.method(args)
func (p *CrossReferenceParser) parseMethodCalls(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	matches := p.regex_method_call.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		targetClass := match[1]
		methodName := match[2]

		// Skip if this matches a method chain (already handled)
		if strings.Contains(line, targetClass+"."+methodName) &&
			p.regex_chain_pattern.MatchString(line) {
			continue
		}

		refs = append(refs, CrossReference{
			Type:         RefMethodCall,
			TargetClass:  targetClass,
			TargetMember: methodName,
			Context:      line,
			LineNumber:   lineNum,
			IsChained:    false,
		})
	}

	return refs
}

// parseFieldAccesses extracts field accesses like obj.field
func (p *CrossReferenceParser) parseFieldAccesses(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	matches := p.regex_field_access.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		targetClass := match[1]
		fieldName := match[2]

		// Skip array access patterns (handled separately)
		if strings.Contains(fieldName, "[") {
			continue
		}

		refs = append(refs, CrossReference{
			Type:         RefFieldAccess,
			TargetClass:  targetClass,
			TargetMember: fieldName,
			Context:      line,
			LineNumber:   lineNum,
		})
	}

	return refs
}

// parseInstantiations extracts constructor calls like new Class(args)
func (p *CrossReferenceParser) parseInstantiations(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	matches := p.regex_instantiation.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		targetClass := match[1]

		refs = append(refs, CrossReference{
			Type:         RefInstantiation,
			TargetClass:  targetClass,
			TargetMember: "<init>", // Constructor
			Context:      line,
			LineNumber:   lineNum,
		})
	}

	return refs
}

// parseStaticCalls extracts static method calls like Class.method(args)
func (p *CrossReferenceParser) parseStaticCalls(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	matches := p.regex_static_call.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		targetClass := match[1]
		methodName := match[2]

		refs = append(refs, CrossReference{
			Type:         RefStaticCall,
			TargetClass:  targetClass,
			TargetMember: methodName,
			Context:      line,
			LineNumber:   lineNum,
		})
	}

	return refs
}

// parseArrayAccesses extracts array access patterns like array[index]
func (p *CrossReferenceParser) parseArrayAccesses(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	matches := p.regex_array_access.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		arrayName := match[1]

		// Try to extract class name from array access
		if p.isValidClassReference(arrayName) {
			refs = append(refs, CrossReference{
				Type:          RefArrayAccess,
				TargetClass:   arrayName,
				TargetMember:  "array_access",
				Context:       line,
				LineNumber:    lineNum,
				IsArrayAccess: true,
			})
		}
	}

	return refs
}

// isValidClassReference checks if a string could be a valid class reference
func (p *CrossReferenceParser) isValidClassReference(name string) bool {
	// Skip very short names (likely variables)
	if len(name) < 2 {
		return false
	}

	// Skip obvious variable names (start with lowercase)
	if strings.ToLower(name) == name && !strings.HasPrefix(name, "a") {
		return false
	}

	// Check if it's a known class or builtin (we'll filter later)
	return p.isProjectClass(name) || p.isBuiltinClass(name) ||
		strings.Contains(name, "_") // Obfuscated names often have underscores
}

// extractMethodName extracts method name from a segment like "method258(anInt1569)"
func (p *CrossReferenceParser) extractMethodName(segment string) string {
	// Find opening parenthesis
	parenIndex := strings.Index(segment, "(")
	if parenIndex == -1 {
		return ""
	}

	methodName := segment[:parenIndex]
	return strings.TrimSpace(methodName)
}

// aggregateReference adds a reference to the aggregate statistics
func (refs *ClassCrossReferences) aggregateReference(ref CrossReference) {
	refs.TotalReferences++

	// Count by type
	refs.ReferenceCounts[ref.Type]++

	// Count unique targets
	refs.UniqueTargets[ref.TargetClass]++

	// Track special patterns
	if ref.IsChained {
		refs.ChainedReferences++
	}
	if ref.IsArrayAccess {
		refs.ArrayReferences++
	}
}

// normalizeFrequencies normalizes reference frequencies for comparison
func (refs *ClassCrossReferences) normalizeFrequencies() {
	if refs.TotalReferences == 0 {
		return
	}

	refs.NormalizedTargets = make(map[string]float64)
	for class, count := range refs.UniqueTargets {
		refs.NormalizedTargets[class] = float64(count) / float64(refs.TotalReferences)
	}

	// Create reference vector for similarity calculations
	refs.ReferenceVector = make(map[string]float64)
	for refType, count := range refs.ReferenceCounts {
		refs.ReferenceVector[refType.String()] = float64(count) / float64(refs.TotalReferences)
	}
}

// CalculateJaccardSimilarity calculates similarity between two cross-reference sets
func (refs *ClassCrossReferences) CalculateJaccardSimilarity(other *ClassCrossReferences) float64 {
	if refs.TotalReferences == 0 && other.TotalReferences == 0 {
		return 1.0
	}
	if refs.TotalReferences == 0 || other.TotalReferences == 0 {
		return 0.0
	}

	// Jaccard similarity for target classes
	intersection := 0
	union := 0

	// Count intersection and union of target classes
	allClasses := make(map[string]bool)
	for class := range refs.UniqueTargets {
		allClasses[class] = true
	}
	for class := range other.UniqueTargets {
		if refs.UniqueTargets[class] > 0 {
			intersection++
		}
		allClasses[class] = true
	}

	union = len(allClasses)

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// String provides a human-readable summary of cross-references
func (refs *ClassCrossReferences) String() string {
	if refs.TotalReferences == 0 {
		return "No cross-references found"
	}

	return fmt.Sprintf("Cross-references: %d total, %d unique targets, %d chained, %d arrays",
		refs.TotalReferences, len(refs.UniqueTargets), refs.ChainedReferences, refs.ArrayReferences)
}

// BytecodeCrossReferenceParser handles parsing of cross-references from bytecode files
type BytecodeCrossReferenceParser struct {
	// Regex patterns optimized for bytecode instruction format
	regex_new_instantiation   *regexp.Regexp // new #constant_pool_index // class ClassName
	regex_array_instantiation *regexp.Regexp // anewarray #constant_pool_index // class ClassName
	regex_invoke_method       *regexp.Regexp // invokevirtual/invokestatic/invokespecial #constant_pool_index // Method ClassName.methodName:signature
	regex_field_access        *regexp.Regexp // getfield/putfield/getstatic/putstatic #constant_pool_index // Field ClassName:fieldName:type
	regex_array_access        *regexp.Regexp // aaload/aastore for array operations

	// Class registry for validation and mapping
	projectClasses map[string]bool
	reverseAnchors map[string]string // obfuscated -> deobfuscated
}

// NewBytecodeCrossReferenceParser creates a new bytecode cross-reference parser
func NewBytecodeCrossReferenceParser(projectClasses []string, reverseAnchors map[string]string) *BytecodeCrossReferenceParser {
	parser := &BytecodeCrossReferenceParser{
		// Bytecode instruction patterns (space-sensitive for accurate parsing)
		regex_new_instantiation:   regexp.MustCompile(`new\s+#(\d+)\s+.*?//\s*class\s+([A-Z]{8}|[a-zA-Z_$\.]+)`),
		regex_array_instantiation: regexp.MustCompile(`anewarray\s+#(\d+)\s+.*?//\s*class\s+([A-Z]{8}|[a-zA-Z_$\.]+)`),
		regex_invoke_method:       regexp.MustCompile(`invoke\w+\s+#(\d+)\s+.*?//\s*Method\s+([A-Z]{8}|[a-zA-Z_$\.]+)\.([^:]+):`),
		regex_field_access:        regexp.MustCompile(`(?:get|put)(?:field|static)\s+#(\d+)\s+.*?//\s*Field\s+([A-Z]{8}|[a-zA-Z_$\.]+):([^:]+)`),
		regex_array_access:        regexp.MustCompile(`(?:aa|ia|ba|ca|sa|la|fa|da)(?:load|store)\s+.*?//\s*(?:array|reference).*?`),

		projectClasses: makeProjectClassMap(projectClasses),
		reverseAnchors: reverseAnchors,
	}

	return parser
}

// ParseBytecodeFile extracts cross-references from a bytecode file
func (p *BytecodeCrossReferenceParser) ParseBytecodeFile(bytecodePath string) (*ClassCrossReferences, error) {
	bytecode, err := os.ReadFile(bytecodePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytecode file %s: %w", bytecodePath, err)
	}

	lines := strings.Split(string(bytecode), "\n")
	refs := &ClassCrossReferences{
		MethodReferences: make(map[string][]CrossReference),
		UniqueTargets:    make(map[string]int),
		ReferenceCounts:  make(map[ReferenceType]int),
	}

	// Process each line for cross-references
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract references from this bytecode instruction
		lineRefs := p.extractReferencesFromBytecodeLine(line, lineNum)
		for _, ref := range lineRefs {
			if p.shouldIncludeBytecodeReference(ref) {
				refs.MethodReferences["<bytecode>"] = append(refs.MethodReferences["<bytecode>"], ref)
				refs.aggregateReference(ref)
			}
		}
	}

	// Normalize frequencies for comparison
	refs.normalizeFrequencies()

	return refs, nil
}

// extractReferencesFromBytecodeLine extracts all types of references from a single bytecode line
func (p *BytecodeCrossReferenceParser) extractReferencesFromBytecodeLine(line string, lineNum int) []CrossReference {
	var refs []CrossReference

	// Extract class instantiations
	if matches := p.regex_new_instantiation.FindAllStringSubmatch(line, -1); matches != nil {
		for _, match := range matches {
			if len(match) >= 3 {
				targetClass := p.mapBytecodeClassName(match[2])
				if targetClass != "" {
					refs = append(refs, CrossReference{
						Type:         RefInstantiation,
						TargetClass:  targetClass,
						TargetMember: "<init>",
						Context:      line,
						LineNumber:   lineNum,
					})
				}
			}
		}
	}

	// Extract array instantiations
	if matches := p.regex_array_instantiation.FindAllStringSubmatch(line, -1); matches != nil {
		for _, match := range matches {
			if len(match) >= 3 {
				targetClass := p.mapBytecodeClassName(match[2])
				if targetClass != "" {
					refs = append(refs, CrossReference{
						Type:          RefArrayAccess,
						TargetClass:   targetClass,
						TargetMember:  "array_creation",
						Context:       line,
						LineNumber:    lineNum,
						IsArrayAccess: true,
					})
				}
			}
		}
	}

	// Extract method invocations
	if matches := p.regex_invoke_method.FindAllStringSubmatch(line, -1); matches != nil {
		for _, match := range matches {
			if len(match) >= 4 {
				targetClass := p.mapBytecodeClassName(match[2])
				methodName := match[3]
				if targetClass != "" && methodName != "" {
					refType := RefMethodCall
					if strings.Contains(line, "invokestatic") {
						refType = RefStaticCall
					}

					refs = append(refs, CrossReference{
						Type:         refType,
						TargetClass:  targetClass,
						TargetMember: methodName,
						Context:      line,
						LineNumber:   lineNum,
					})
				}
			}
		}
	}

	// Extract field accesses
	if matches := p.regex_field_access.FindAllStringSubmatch(line, -1); matches != nil {
		for _, match := range matches {
			if len(match) >= 4 {
				targetClass := p.mapBytecodeClassName(match[2])
				fieldName := match[3]
				if targetClass != "" && fieldName != "" {
					refs = append(refs, CrossReference{
						Type:         RefFieldAccess,
						TargetClass:  targetClass,
						TargetMember: fieldName,
						Context:      line,
						LineNumber:   lineNum,
					})
				}
			}
		}
	}

	// Extract array access patterns (simplified detection)
	if p.regex_array_access.MatchString(line) {
		// For array operations, we don't have specific class targets
		// but we can mark this as a general array access pattern
		refs = append(refs, CrossReference{
			Type:          RefArrayAccess,
			TargetClass:   "<array_operation>",
			TargetMember:  "array_access",
			Context:       line,
			LineNumber:    lineNum,
			IsArrayAccess: true,
		})
	}

	return refs
}

// mapBytecodeClassName converts bytecode class names to deobfuscated names when possible
func (p *BytecodeCrossReferenceParser) mapBytecodeClassName(bytecodeName string) string {
	// First check if this is an obfuscated class name (8 uppercase letters)
	if len(bytecodeName) == 8 && strings.ToUpper(bytecodeName) == bytecodeName {
		if deobfuscated, exists := p.reverseAnchors[bytecodeName]; exists {
			return deobfuscated
		}
	}

	// Check if this is a project class (either obfuscated or deobfuscated name)
	if p.projectClasses[bytecodeName] {
		return bytecodeName
	}

	// For standard Java classes, keep as-is but check if we should exclude them
	if p.isBuiltinClass(bytecodeName) {
		return "" // Exclude built-in classes
	}

	// If we can't map it but it's a valid class reference, keep the obfuscated name
	// This allows similarity comparison even for unmapped classes
	if len(bytecodeName) > 0 {
		return bytecodeName
	}

	return ""
}

// isBuiltinClass checks if a class should be excluded from analysis (same logic as source parser)
func (p *BytecodeCrossReferenceParser) isBuiltinClass(className string) bool {
	// Check direct match
	if BUILTIN_EXCLUSIONS[className] {
		return true
	}

	// Check package prefixes
	return strings.HasPrefix(className, "java.") ||
		strings.HasPrefix(className, "javax.") ||
		strings.HasPrefix(className, "sun.") ||
		strings.HasPrefix(className, "com.sun.") ||
		strings.HasPrefix(className, "org.")
}

// shouldIncludeBytecodeReference determines if a bytecode reference should be included in analysis
func (p *BytecodeCrossReferenceParser) shouldIncludeBytecodeReference(ref CrossReference) bool {
	// Exclude built-in classes
	if p.isBuiltinClass(ref.TargetClass) || ref.TargetClass == "" {
		return false
	}

	// Include if it's a known project class or mapped anchor class
	if p.projectClasses[ref.TargetClass] {
		return true
	}

	// Include if it's an obfuscated name that could match (we'll validate during similarity)
	if len(ref.TargetClass) == 8 && strings.ToUpper(ref.TargetClass) == ref.TargetClass {
		return true
	}

	// For now, be conservative and only include known classes
	return false
}
