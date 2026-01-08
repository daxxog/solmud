// internal_behavior.go - Internal Behavioral Pattern Analysis (Phase 3.2.2.1)

package main

import (
	"regexp"
	"strings"
)

// InternalBehaviorParser analyzes internal behavioral patterns within Java source code
type InternalBehaviorParser struct {
	// Regex patterns for internal analysis
	regex_method_call      *regexp.Regexp // this.method() calls
	regex_field_read       *regexp.Regexp // field access (read)
	regex_field_write      *regexp.Regexp // field assignment (write)
	regex_for_loop         *regexp.Regexp // for loops
	regex_while_loop       *regexp.Regexp // while loops
	regex_iterator_usage   *regexp.Regexp // iterator patterns
	regex_stream_operation *regexp.Regexp // stream operations
	regex_method_decl      *regexp.Regexp // method declarations
	regex_getter_pattern   *regexp.Regexp // get* methods
	regex_setter_pattern   *regexp.Regexp // set* methods
	regex_builder_pattern  *regexp.Regexp // with*, add* methods
	regex_factory_pattern  *regexp.Regexp // create*, new* methods
	regex_event_pattern    *regexp.Regexp // on*, handle* methods
	regex_utility_pattern  *regexp.Regexp // format*, parse*, validate* methods

	// Enhanced array pattern analysis (Phase 3.2.2.2.1)
	regex_single_dim_array  *regexp.Regexp // array[index]
	regex_multi_dim_array   *regexp.Regexp // array[][][][] or array[x][y][z]
	regex_array_in_for      *regexp.Regexp // for (int i = 0; i < array.length; i++)
	regex_bulk_operations   *regexp.Regexp // System.arraycopy, Arrays.fill, Arrays.copyOf
	regex_array_calculation *regexp.Regexp // array[index] + something, array[index] * something
	regex_array_assignment  *regexp.Regexp // array[index] = value
}

// NewInternalBehaviorParser creates a new internal behavior parser
func NewInternalBehaviorParser() *InternalBehaviorParser {
	return &InternalBehaviorParser{
		// Enhanced regex for Jad output patterns (obfuscated names: aClass123, method456)
		regex_method_call: regexp.MustCompile(`(?:this\.)?([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(`),
		// Field reads: this.fieldName (not followed by =)
		regex_field_read: regexp.MustCompile(`this\.([a-zA-Z_$][a-zA-Z0-9_$]*)(?:\s*[;=]|\s*\.)`),
		// Field writes: this.fieldName = or this.fieldName +=
		regex_field_write: regexp.MustCompile(`this\.([a-zA-Z_$][a-zA-Z0-9_$]*)\s*[+\-*/]?=`),
		// For loops: for (...) {
		regex_for_loop: regexp.MustCompile(`\bfor\s*\(`),
		// While loops: while (...) {
		regex_while_loop: regexp.MustCompile(`\bwhile\s*\(`),
		// Iterator usage: .iterator() or Iterator<
		regex_iterator_usage: regexp.MustCompile(`(?:\.iterator\(\)|Iterator\s*<)`),
		// Stream operations: .stream(), .map(, .filter(, .collect(
		regex_stream_operation: regexp.MustCompile(`\.(?:stream|map|filter|collect|forEach)\s*\(`),
		// Method declarations: public/private/protected void methodName(
		regex_method_decl: regexp.MustCompile(`^\s*(?:public|private|protected)?\s*(?:static\s+)?(?:final\s+)?[a-zA-Z_$][a-zA-Z0-9_$<>\[\]]*\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(`),
		// Getter patterns: getXxx(
		regex_getter_pattern: regexp.MustCompile(`\bget[A-Z][a-zA-Z0-9_]*\b`),
		// Setter patterns: setXxx(
		regex_setter_pattern: regexp.MustCompile(`\bset[A-Z][a-zA-Z0-9_]*\b`),
		// Builder patterns: withXxx(, addXxx(
		regex_builder_pattern: regexp.MustCompile(`\b(?:with|add)[A-Z][a-zA-Z0-9_]*\b`),
		// Factory patterns: createXxx(, newInstanceXxx(
		regex_factory_pattern: regexp.MustCompile(`\b(?:create|newInstance)[A-Z][a-zA-Z0-9_]*\b`),
		// Event patterns: onXxx(, handleXxx(
		regex_event_pattern: regexp.MustCompile(`\b(?:on|handle)[A-Z][a-zA-Z0-9_]*\b`),
		// Utility patterns: formatXxx(, parseXxx(, validateXxx(
		regex_utility_pattern: regexp.MustCompile(`\b(?:format|parse|validate)[A-zA-Z0-9_]*\b`),

		// Enhanced array pattern analysis (Phase 3.2.2.2.1)
		// Single dimensional: array[index]
		regex_single_dim_array: regexp.MustCompile(`([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\[\s*[^]]+\s*\]`),
		// Multi-dimensional: array[][][] or array[x][y][z]
		regex_multi_dim_array: regexp.MustCompile(`([a-zA-Z_$][a-zA-Z0-9_$]*)(?:\s*\[\s*[^]]+\s*\])+(?:\s*\[\s*[^]]+\s*\])`),
		// Array access in for loops: for (int i = 0; i < array.length; i++)
		regex_array_in_for: regexp.MustCompile(`for\s*\(\s*[^;]*;\s*[^;]*<\s*([a-zA-Z_$][a-zA-Z0-9_$]*)\.length[^;]*;\s*[^)]*\)`),
		// Bulk operations: System.arraycopy, Arrays.fill, Arrays.copyOf
		regex_bulk_operations: regexp.MustCompile(`(?:System\.arraycopy|Arrays\.(?:fill|copyOf|sort))\s*\(`),
		// Array calculations: array[index] + something, array[index] * something
		regex_array_calculation: regexp.MustCompile(`([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\[\s*[^]]+\s*\]\s*[+\-*/%]\s*[^;]+`),
		// Array assignments: array[index] = value
		regex_array_assignment: regexp.MustCompile(`([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\[\s*[^]]+\s*\]\s*=\s*[^;]+`),
	}
}

// ParseInternalBehavior analyzes internal behavioral patterns from Java source
func (p *InternalBehaviorParser) ParseInternalBehavior(javaSource string, methods []MethodInfo) *InternalBehavioralPatterns {
	patterns := &InternalBehavioralPatterns{
		MethodCallGraph:   p.parseMethodCallGraph(javaSource, methods),
		StatePatterns:     p.parseStateManipulation(javaSource, methods),
		IterationPatterns: p.parseIterationPatterns(javaSource),
		SemanticPatterns:  p.parseSemanticPatterns(methods),
		ArrayPatterns:     p.detectArrayPatterns(javaSource),
	}

	return patterns
}

// parseMethodCallGraph builds the intra-class method call graph
func (p *InternalBehaviorParser) parseMethodCallGraph(javaSource string, methods []MethodInfo) *InternalCallGraph {
	graph := &InternalCallGraph{
		MethodCalls: make(map[string][]string),
		CallFreq:    make(map[string]int),
		Recursion:   make(map[string]bool),
	}

	lines := strings.Split(javaSource, "\n")
	var currentMethod string

	// Method name to signature mapping for validation
	methodNames := make(map[string]bool)
	for _, method := range methods {
		methodNames[method.Name] = true
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Track current method context
		if methodMatch := p.extractCurrentMethod(line); methodMatch != "" {
			currentMethod = methodMatch
			continue
		}

		// Only analyze within method bodies
		if currentMethod == "" {
			continue
		}

		// Find method calls in this line
		calls := p.extractMethodCalls(line, methodNames)

		// Add calls to graph (exclude self-calls for recursion detection)
		for _, callee := range calls {
			if callee != currentMethod {
				// Add to call graph
				found := false
				for _, existing := range graph.MethodCalls[currentMethod] {
					if existing == callee {
						found = true
						break
					}
				}
				if !found {
					graph.MethodCalls[currentMethod] = append(graph.MethodCalls[currentMethod], callee)
				}
				graph.CallFreq[callee]++
			} else {
				// Self-call = recursion
				graph.Recursion[currentMethod] = true
			}
		}
	}

	return graph
}

// parseStateManipulation analyzes field read/write patterns
func (p *InternalBehaviorParser) parseStateManipulation(javaSource string, methods []MethodInfo) *StateManipulationPatterns {
	patterns := &StateManipulationPatterns{
		FieldReads:   make(map[string]int),
		FieldWrites:  make(map[string]int),
		SetStateFreq: 0,
	}

	lines := strings.Split(javaSource, "\n")

	// For now, we'll analyze field access patterns in method bodies
	// In a full implementation, we'd need field declarations from the parser
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Analyze field reads: this.fieldName (not assignment)
		if matches := p.regex_field_read.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				if len(match) >= 2 {
					fieldName := match[1]
					patterns.FieldReads[fieldName]++
				}
			}
		}

		// Analyze field writes: this.fieldName =
		if matches := p.regex_field_write.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				if len(match) >= 2 {
					fieldName := match[1]
					patterns.FieldWrites[fieldName]++

					// Check for setState patterns
					if strings.Contains(line, "setState") {
						patterns.SetStateFreq++
					}
				}
			}
		}
	}

	return patterns
}

// parseIterationPatterns analyzes loop and iteration usage
func (p *InternalBehaviorParser) parseIterationPatterns(javaSource string) *IterationPatterns {
	patterns := &IterationPatterns{}

	lines := strings.Split(javaSource, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Count for loops
		if p.regex_for_loop.MatchString(line) {
			patterns.ForLoops++
		}

		// Count while loops
		if p.regex_while_loop.MatchString(line) {
			patterns.WhileLoops++
		}

		// Count iterator usage
		if p.regex_iterator_usage.MatchString(line) {
			patterns.IteratorUsage++
		}

		// Count stream operations
		if p.regex_stream_operation.MatchString(line) {
			patterns.StreamProcessing++
		}
	}

	return patterns
}

// parseSemanticPatterns analyzes method naming patterns
func (p *InternalBehaviorParser) parseSemanticPatterns(methods []MethodInfo) *SemanticMethodPatterns {
	patterns := &SemanticMethodPatterns{}

	for _, method := range methods {
		methodName := method.Name

		// Classify by naming patterns
		if p.regex_getter_pattern.MatchString(methodName) {
			patterns.GetterMethods = append(patterns.GetterMethods, methodName)
		}
		if p.regex_setter_pattern.MatchString(methodName) {
			patterns.SetterMethods = append(patterns.SetterMethods, methodName)
		}
		if p.regex_builder_pattern.MatchString(methodName) {
			patterns.BuilderMethods = append(patterns.BuilderMethods, methodName)
		}
		if p.regex_factory_pattern.MatchString(methodName) {
			patterns.FactoryMethods = append(patterns.FactoryMethods, methodName)
		}
		if p.regex_event_pattern.MatchString(methodName) {
			patterns.EventHandlers = append(patterns.EventHandlers, methodName)
		}
		if p.regex_utility_pattern.MatchString(methodName) {
			patterns.UtilityMethods = append(patterns.UtilityMethods, methodName)
		}
	}

	return patterns
}

// extractCurrentMethod extracts the current method name from a method declaration
func (p *InternalBehaviorParser) extractCurrentMethod(line string) string {
	matches := p.regex_method_decl.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractMethodCalls extracts method calls from a line, filtering to known methods
func (p *InternalBehaviorParser) extractMethodCalls(line string, knownMethods map[string]bool) []string {
	var calls []string

	// Skip lines that are method declarations
	if p.regex_method_decl.MatchString(line) {
		return calls
	}

	matches := p.regex_method_call.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			methodName := match[1]
			// Only include calls to known methods in this class
			if knownMethods[methodName] {
				calls = append(calls, methodName)
			}
		}
	}

	return calls
}

// detectArrayPatterns performs comprehensive array pattern analysis
func (p *InternalBehaviorParser) detectArrayPatterns(javaSource string) *ArrayPatternMetrics {
	metrics := &ArrayPatternMetrics{
		ArrayAlgorithms: make([]string, 0),
	}

	lines := strings.Split(javaSource, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Count single dimensional array access
		if p.regex_single_dim_array.MatchString(line) {
			metrics.SingleDimAccess++

			// Check if this single-dim access is actually part of multi-dim
			if p.regex_multi_dim_array.MatchString(line) {
				metrics.MultiDimAccess++
			}
		}

		// Count multi-dimensional array access (additional detection)
		if matches := p.regex_multi_dim_array.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				if len(match) > 1 {
					arrayName := match[1]
					// Count brackets to determine depth
					depth := strings.Count(line, "[") - strings.Count(line, "]")
					if depth > metrics.NestedArrayDepth {
						metrics.NestedArrayDepth = depth
					}
					metrics.MultiDimAccess++
					// Use arrayName to avoid unused variable warning
					_ = arrayName
				}
			}
		}

		// Count array access in loops
		if p.regex_array_in_for.MatchString(line) {
			metrics.LoopArrayAccess++
		}

		// Count bulk operations
		if p.regex_bulk_operations.MatchString(line) {
			metrics.BulkOperations++

			// Detect specific algorithms
			if strings.Contains(line, "Arrays.sort") {
				if !contains(metrics.ArrayAlgorithms, "sort") {
					metrics.ArrayAlgorithms = append(metrics.ArrayAlgorithms, "sort")
				}
			}
			if strings.Contains(line, "System.arraycopy") {
				if !contains(metrics.ArrayAlgorithms, "copy") {
					metrics.ArrayAlgorithms = append(metrics.ArrayAlgorithms, "copy")
				}
			}
		}

		// Additional algorithm detection from patterns
		if p.detectSortingPattern(line) && !contains(metrics.ArrayAlgorithms, "sort") {
			metrics.ArrayAlgorithms = append(metrics.ArrayAlgorithms, "sort")
		}
		if p.detectSearchPattern(line) && !contains(metrics.ArrayAlgorithms, "search") {
			metrics.ArrayAlgorithms = append(metrics.ArrayAlgorithms, "search")
		}
		if p.detectFilterPattern(line) && !contains(metrics.ArrayAlgorithms, "filter") {
			metrics.ArrayAlgorithms = append(metrics.ArrayAlgorithms, "filter")
		}
	}

	return metrics
}

// analyzeAccessContext analyzes the context of array access patterns
func (p *InternalBehaviorParser) analyzeAccessContext(javaSource string, arrayName string) *ArrayAccessContext {
	context := &ArrayAccessContext{}

	lines := strings.Split(javaSource, "\n")
	totalLines := len(lines)
	arrayAccesses := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Check if this line contains access to the specific array
		arrayPattern := regexp.MustCompile(arrayName + `\s*\[`)
		if arrayPattern.MatchString(line) {
			arrayAccesses++

			// Check context
			if p.isInLoop(line) {
				context.InLoop = true
			}
			if p.isInMethodCall(line) {
				context.InMethodCall = true
			}
			if p.regex_array_calculation.MatchString(line) {
				context.InCalculation = true
			}
			if p.regex_array_assignment.MatchString(line) {
				context.InAssignment = true
			}
		}
	}

	// Calculate access frequency
	if totalLines > 0 {
		context.AccessFrequency = float64(arrayAccesses) / float64(totalLines)
	}

	return context
}

// Helper functions for array pattern detection

// isInLoop checks if a line is inside a loop construct
func (p *InternalBehaviorParser) isInLoop(line string) bool {
	return p.regex_for_loop.MatchString(line) || p.regex_while_loop.MatchString(line)
}

// isInMethodCall checks if array access is in a method call context
func (p *InternalBehaviorParser) isInMethodCall(line string) bool {
	// Simple heuristic: if there are parentheses after array access
	return strings.Contains(line, "(") && strings.Contains(line, ")")
}

// detectSortingPattern detects basic sorting algorithm patterns
func (p *InternalBehaviorParser) detectSortingPattern(line string) bool {
	// Look for nested loops with array swaps (bubble sort pattern)
	swapPatterns := []string{
		`temp\s*=.*\[.*\];.*\[.*\]\s*=.*\[.*\];.*\[.*\]\s*=\s*temp`, // swap pattern
		`if\s*\(.*\[.*\]\s*>\s*.*\[.*\]\)`,                          // comparison pattern
	}
	for _, pattern := range swapPatterns {
		if regexp.MustCompile(pattern).MatchString(line) {
			return true
		}
	}
	return false
}

// detectSearchPattern detects search algorithm patterns
func (p *InternalBehaviorParser) detectSearchPattern(line string) bool {
	searchPatterns := []string{
		`for\s*\([^;]*;\s*[^;]*<\s*.*\.length[^;]*;`, // linear search loop
		`while\s*\([^;]*&&\s*.*\[.*\]\s*!=\s*`,       // search condition
	}
	for _, pattern := range searchPatterns {
		if regexp.MustCompile(pattern).MatchString(line) {
			return true
		}
	}
	return false
}

// detectFilterPattern detects filtering/copying patterns
func (p *InternalBehaviorParser) detectFilterPattern(line string) bool {
	filterPatterns := []string{
		`if\s*\(.*\[.*\]\s*[<>!=]+\s*`,               // conditional array access
		`[a-zA-Z_$][a-zA-Z0-9_$]*\[.*\]\s*=.*\[.*\]`, // array to array copy
	}
	for _, pattern := range filterPatterns {
		if regexp.MustCompile(pattern).MatchString(line) {
			return true
		}
	}
	return false
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
