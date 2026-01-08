package main

import (
	"math"
	"strings"
	// ReferenceType is defined in crossref.go
)

const (
	WEIGHT_INTERFACE_MATCH      = 20.0
	WEIGHT_SUPERCLASS_MATCH     = 25.0
	WEIGHT_FIELD_COUNT_MATCH    = 5.0
	WEIGHT_FIELD_TYPE_SIM       = 10.0
	WEIGHT_FIELD_PATTERN_SIM    = 8.0 // New: field type patterns
	WEIGHT_METHOD_COUNT_MATCH   = 5.0
	WEIGHT_METHOD_SIG_SIM       = 20.0
	WEIGHT_METHOD_NAME_SIM      = 7.0 // New: method name similarity
	WEIGHT_CONSTRUCTOR_MATCH    = 5.0
	WEIGHT_ACCESS_MOD_MATCH     = 5.0
	WEIGHT_FUNCTIONAL_PATTERN   = 12.0 // New: functional grouping patterns
	WEIGHT_CROSSREF_SIMILARITY  = 15.0 // New: cross-reference similarity
	WEIGHT_UNIQUE_PATTERNS      = 8.0  // New: unique behavioral patterns
	WEIGHT_BEHAVIORAL_SIGNATURE = 12.0 // New: behavioral signature matching

	// Phase 3.2.2.1: Internal behavioral patterns
	WEIGHT_METHOD_CALL_GRAPH  = 10.0 // Intra-class method call patterns
	WEIGHT_STATE_MANIPULATION = 12.0 // Field read/write patterns
	WEIGHT_ITERATION_PATTERNS = 8.0  // Loop and iteration usage
	WEIGHT_SEMANTIC_METHODS   = 6.0  // Method naming patterns (getters/setters/etc)

	// Phase 3.2.2.2.2: Data structure type classification
	WEIGHT_ARRAY_PATTERN_ANALYSIS = 15.0 // Base array analysis weight
	WEIGHT_MULTIDIMENSIONAL       = 8.0  // Multi-dimensional bonus
	WEIGHT_GRAPHICS_ARRAYS        = 10.0 // Graphics-specific patterns
	WEIGHT_BULK_OPERATIONS        = 6.0  // Bulk array operations
	WEIGHT_3D_VERTEX_DATA         = 12.0 // Model class specific
	WEIGHT_TEXTURE_DATA           = 10.0 // Texture class specific
	WEIGHT_WORLD_DATA             = 9.0  // WorldController specific

	// Phase 1.1A: Magic Constants & Static Fields Pattern
	WEIGHT_MAGIC_CONSTANTS = 20.0 // Distinctive magic numbers and constants

	SIZE_DIFFERENCE_PENALTY   = 10.0
	MAX_SIZE_DIFFERENCE_RATIO = 1.5
)

type ScoreBreakdown struct {
	InterfaceMatch      float64
	SuperclassMatch     float64
	FieldCountMatch     float64
	FieldSimilarity     float64
	FieldPatternSim     float64
	MethodCountMatch    float64
	MethodSimilarity    float64
	MethodNameSim       float64
	ConstructorMatch    float64
	AccessMatch         float64
	FunctionalPattern   float64
	CrossrefSimilarity  float64
	UniquePatterns      float64
	BehavioralSignature float64

	// Phase 3.2.2.1: Internal behavioral patterns
	MethodCallGraph   float64
	StateManipulation float64
	IterationPatterns float64
	SemanticMethods   float64

	// Phase 3.2.2.2.2: Data structure classification
	GraphicsArrayBonus float64

	// Phase 1.1A: Magic Constants & Static Fields Pattern
	MagicConstants float64

	SizePenalty float64
}

type Scorer struct {
	anchor_mappings map[string]string
}

func NewScorer(anchors map[string]string) *Scorer {
	return &Scorer{
		anchor_mappings: anchors,
	}
}

func (self *Scorer) CalculateScore(deob_class *ClassInfo, obf_class *ClassInfo) ScoreBreakdown {
	breakdown := ScoreBreakdown{}

	interface_points := self.calculate_interface_score(deob_class, obf_class)
	breakdown.InterfaceMatch = interface_points

	superclass_points := self.calculate_superclass_score(deob_class, obf_class)
	breakdown.SuperclassMatch = superclass_points

	if len(deob_class.Fields) == len(obf_class.Fields) {
		breakdown.FieldCountMatch = WEIGHT_FIELD_COUNT_MATCH
	}

	field_sim := self.calculate_field_similarity(deob_class, obf_class)
	breakdown.FieldSimilarity = field_sim * WEIGHT_FIELD_TYPE_SIM

	field_pattern_sim := self.calculate_field_pattern_similarity(deob_class, obf_class)
	breakdown.FieldPatternSim = field_pattern_sim * WEIGHT_FIELD_PATTERN_SIM

	if len(deob_class.Methods) == len(obf_class.Methods) {
		breakdown.MethodCountMatch = WEIGHT_METHOD_COUNT_MATCH
	}

	method_sim := self.calculate_method_similarity(deob_class, obf_class)
	breakdown.MethodSimilarity = method_sim * WEIGHT_METHOD_SIG_SIM

	method_name_sim := self.calculate_method_name_similarity(deob_class, obf_class)
	breakdown.MethodNameSim = method_name_sim * WEIGHT_METHOD_NAME_SIM

	functional_pattern := self.calculate_functional_pattern_score(deob_class, obf_class)
	breakdown.FunctionalPattern = functional_pattern * WEIGHT_FUNCTIONAL_PATTERN

	// Cross-reference similarity scoring
	crossref_sim := self.calculate_crossref_similarity(deob_class, obf_class)
	breakdown.CrossrefSimilarity = crossref_sim * WEIGHT_CROSSREF_SIMILARITY

	unique_patterns := self.calculate_unique_patterns(deob_class, obf_class)
	breakdown.UniquePatterns = unique_patterns * WEIGHT_UNIQUE_PATTERNS

	behavioral_sig := self.calculate_behavioral_signature(deob_class, obf_class)
	breakdown.BehavioralSignature = behavioral_sig * WEIGHT_BEHAVIORAL_SIGNATURE

	// Phase 3.2.2.1: Internal behavioral pattern scoring
	method_call_graph := self.calculate_method_call_graph_similarity(deob_class, obf_class)
	breakdown.MethodCallGraph = method_call_graph * WEIGHT_METHOD_CALL_GRAPH

	state_manipulation := self.calculate_state_manipulation_similarity(deob_class, obf_class)
	breakdown.StateManipulation = state_manipulation * WEIGHT_STATE_MANIPULATION

	iteration_patterns := self.calculate_iteration_pattern_similarity(deob_class, obf_class)
	breakdown.IterationPatterns = iteration_patterns * WEIGHT_ITERATION_PATTERNS

	semantic_methods := self.calculate_semantic_method_similarity(deob_class, obf_class)
	breakdown.SemanticMethods = semantic_methods * WEIGHT_SEMANTIC_METHODS

	// Phase 3.2.2.2.2: Data structure classification scoring
	graphics_bonus := self.calculateGraphicsArrayBonus(deob_class, obf_class)
	breakdown.GraphicsArrayBonus = graphics_bonus

	// Phase 1.1A: Magic Constants & Static Fields Pattern scoring
	magic_constants := self.calculateMagicConstantsScore(deob_class, obf_class)
	breakdown.MagicConstants = magic_constants

	if len(deob_class.Constructors) == len(obf_class.Constructors) {
		breakdown.ConstructorMatch = WEIGHT_CONSTRUCTOR_MATCH
	}

	if self.access_modifiers_match(deob_class, obf_class) {
		breakdown.AccessMatch = WEIGHT_ACCESS_MOD_MATCH
	}

	if self.significant_size_difference(deob_class, obf_class) {
		breakdown.SizePenalty = -SIZE_DIFFERENCE_PENALTY
	}

	return breakdown
}

// calculate_crossref_similarity computes Jaccard similarity between cross-reference sets
func (self *Scorer) calculate_crossref_similarity(deob, obf *ClassInfo) float64 {
	// Adaptive scoring: work with partial cross-reference data
	return self.calculateAdaptiveCrossrefSimilarity(deob, obf)
}

// calculateAdaptiveCrossrefSimilarity handles cases where one or both sides may lack cross-reference data
func (self *Scorer) calculateAdaptiveCrossrefSimilarity(deob, obf *ClassInfo) float64 {
	switch {
	case deob.CrossReferences != nil && obf.CrossReferences != nil:
		// Both sides have data - use full Jaccard similarity
		return self.calculateFullCrossrefSimilarity(deob, obf)
	case deob.CrossReferences != nil && deob.CrossReferences.TotalReferences > 0:
		// Only deobfuscated side has data - use structural fallback
		return self.calculateSourceToBytecodeSimilarity(deob, obf)
	case obf.CrossReferences != nil && obf.CrossReferences.TotalReferences > 0:
		// Only obfuscated side has data - use structural fallback
		return self.calculateBytecodeToSourceSimilarity(deob, obf)
	default:
		// Neither side has cross-reference data
		return 0.0
	}
}

// calculateFullCrossrefSimilarity computes full Jaccard similarity when both sides have data
func (self *Scorer) calculateFullCrossrefSimilarity(deob, obf *ClassInfo) float64 {
	if deob.CrossReferences.TotalReferences == 0 || obf.CrossReferences.TotalReferences == 0 {
		return 0.0
	}

	// Use Jaccard similarity for target class overlap
	jaccard := deob.CrossReferences.CalculateJaccardSimilarity(obf.CrossReferences)

	// Weight by reference density (more references = more reliable similarity)
	deob_density := float64(deob.CrossReferences.TotalReferences) / float64(len(deob.Fields)+len(deob.Methods)+1)
	obf_density := float64(obf.CrossReferences.TotalReferences) / float64(len(obf.Fields)+len(obf.Methods)+1)

	// Average density factor (normalize to 0-1 range, cap at 1.0)
	density_factor := math.Min(deob_density, obf_density) / math.Max(deob_density, obf_density)
	density_factor = math.Min(density_factor, 1.0)

	return jaccard * density_factor
}

// calculateSourceToBytecodeSimilarity provides fallback when only source has cross-reference data
func (self *Scorer) calculateSourceToBytecodeSimilarity(deob, obf *ClassInfo) float64 {
	// Use structural similarity as proxy for behavioral similarity
	structural_score := self.calculate_method_similarity(deob, obf)*WEIGHT_METHOD_SIG_SIM +
		self.calculate_field_similarity(deob, obf)*WEIGHT_FIELD_TYPE_SIM

	// Reduce weight since we don't have actual behavioral data
	return structural_score * 0.3 // 30% of normal behavioral weight
}

// calculateBytecodeToSourceSimilarity provides fallback when only bytecode has cross-reference data
func (self *Scorer) calculateBytecodeToSourceSimilarity(deob, obf *ClassInfo) float64 {
	// Similar to source-to-bytecode but with even lower confidence
	structural_score := self.calculate_method_similarity(deob, obf)*WEIGHT_METHOD_SIG_SIM +
		self.calculate_field_similarity(deob, obf)*WEIGHT_FIELD_TYPE_SIM

	// Even lower weight for bytecode-only data
	return structural_score * 0.2 // 20% of normal behavioral weight
}

// calculate_unique_patterns scores classes with distinctive behavioral patterns
func (self *Scorer) calculate_unique_patterns(deob, obf *ClassInfo) float64 {
	if deob.CrossReferences == nil || obf.CrossReferences == nil {
		return 0.0
	}

	score := 0.0

	// Reward chained references (complex interactions)
	deob_chained := float64(deob.CrossReferences.ChainedReferences)
	obf_chained := float64(obf.CrossReferences.ChainedReferences)

	if deob_chained > 0 && obf_chained > 0 {
		chain_ratio := math.Min(deob_chained, obf_chained) / math.Max(deob_chained, obf_chained)
		score += chain_ratio * 0.4
	}

	// Reward array access patterns (data structures)
	deob_arrays := float64(deob.CrossReferences.ArrayReferences)
	obf_arrays := float64(obf.CrossReferences.ArrayReferences)

	if deob_arrays > 0 && obf_arrays > 0 {
		array_ratio := math.Min(deob_arrays, obf_arrays) / math.Max(deob_arrays, obf_arrays)
		score += array_ratio * 0.3
	}

	// Reward unique reference combinations
	deob_unique := len(deob.CrossReferences.UniqueTargets)
	obf_unique := len(obf.CrossReferences.UniqueTargets)

	if deob_unique > 0 && obf_unique > 0 {
		unique_ratio := math.Min(float64(deob_unique), float64(obf_unique)) / math.Max(float64(deob_unique), float64(obf_unique))
		score += unique_ratio * 0.3
	}

	return score
}

// calculate_behavioral_signature analyzes reference type distributions
func (self *Scorer) calculate_behavioral_signature(deob, obf *ClassInfo) float64 {
	if deob.CrossReferences == nil || obf.CrossReferences == nil {
		return 0.0
	}

	// Compare reference type distributions (method calls, field accesses, etc.)
	similarity := 0.0
	total_weight := 0.0

	reference_types := []ReferenceType{RefMethodCall, RefInstantiation, RefFieldAccess, RefStaticCall, RefArrayAccess}
	weights := []float64{0.4, 0.3, 0.2, 0.05, 0.05} // Method calls most important

	for i, refType := range reference_types {
		deob_count := deob.CrossReferences.ReferenceCounts[refType]
		obf_count := obf.CrossReferences.ReferenceCounts[refType]

		if deob_count == 0 && obf_count == 0 {
			continue // No references of this type
		}

		// Calculate similarity for this reference type
		if deob_count > 0 && obf_count > 0 {
			ratio := math.Min(float64(deob_count), float64(obf_count)) / math.Max(float64(deob_count), float64(obf_count))
			similarity += ratio * weights[i]
		}

		total_weight += weights[i]
	}

	if total_weight == 0 {
		return 0.0
	}

	return similarity / total_weight
}

func (self *Scorer) calculate_interface_score(deob, obf *ClassInfo) float64 {
	resolved_interfaces := self.resolve_interfaces(deob.Interfaces)
	matches := 0

	for _, deob_iface := range resolved_interfaces {
		for _, obf_iface := range obf.Interfaces {
			if deob_iface == obf_iface {
				matches++
				break
			}
		}
	}

	if matches == len(deob.Interfaces) && len(deob.Interfaces) > 0 {
		return WEIGHT_INTERFACE_MATCH
	}
	return 0.0
}

func (self *Scorer) calculate_superclass_score(deob, obf *ClassInfo) float64 {
	resolved_super := self.resolve_superclass(deob.Superclass)
	if resolved_super == obf.Superclass {
		return WEIGHT_SUPERCLASS_MATCH
	}
	return 0.0
}

func (self *Scorer) calculate_field_similarity(deob, obf *ClassInfo) float64 {
	if len(deob.Fields) == 0 || len(obf.Fields) == 0 {
		return 0.0
	}

	total_fields := math.Max(float64(len(deob.Fields)), float64(len(obf.Fields)))
	matching_fields := 0

	for _, deob_field := range deob.Fields {
		for _, obf_field := range obf.Fields {
			if deob_field.TypeName == obf_field.TypeName {
				matching_fields++
				break
			}
		}
	}

	return float64(matching_fields) / total_fields
}

func (self *Scorer) calculate_field_pattern_similarity(deob, obf *ClassInfo) float64 {
	if len(deob.Fields) == 0 || len(obf.Fields) == 0 {
		return 0.0
	}

	// Count field types
	deob_types := make(map[string]int)
	obf_types := make(map[string]int)

	for _, field := range deob.Fields {
		deob_types[field.TypeName]++
	}
	for _, field := range obf.Fields {
		obf_types[field.TypeName]++
	}

	// Calculate similarity based on type distribution
	total_types := 0
	matching_types := 0

	for typ, count := range deob_types {
		total_types++
		if obf_count, exists := obf_types[typ]; exists {
			// Consider it a match if counts are similar (within 2)
			if math.Abs(float64(count-obf_count)) <= 2 {
				matching_types++
			}
		}
	}

	if total_types == 0 {
		return 0.0
	}
	return float64(matching_types) / float64(total_types)
}

func (self *Scorer) calculate_method_similarity(deob, obf *ClassInfo) float64 {
	if len(deob.Methods) == 0 || len(obf.Methods) == 0 {
		return 0.0
	}

	total_methods := math.Max(float64(len(deob.Methods)), float64(len(obf.Methods)))
	matching_methods := 0

	for _, deob_method := range deob.Methods {
		for _, obf_method := range obf.Methods {
			if deob_method.ReturnType == obf_method.ReturnType &&
				len(deob_method.Parameters) == len(obf_method.Parameters) {
				// Enhanced: also check parameter types
				param_match := true
				for i, param := range deob_method.Parameters {
					if i >= len(obf_method.Parameters) || param != obf_method.Parameters[i] {
						param_match = false
						break
					}
				}
				if param_match {
					matching_methods++
					break
				}
			}
		}
	}

	return float64(matching_methods) / total_methods
}

func (self *Scorer) calculate_method_name_similarity(deob, obf *ClassInfo) float64 {
	if len(deob.Methods) == 0 || len(obf.Methods) == 0 {
		return 0.0
	}

	// For deobfuscated code, method names are meaningful
	// For obfuscated code, method names are single characters
	// Look for patterns in method names
	deob_named_methods := 0
	obf_single_char_methods := 0

	for _, method := range deob.Methods {
		if len(method.Name) > 1 { // Meaningful name
			deob_named_methods++
		}
	}

	for _, method := range obf.Methods {
		if len(method.Name) == 1 { // Single character (obfuscated)
			obf_single_char_methods++
		}
	}

	// If deob has many named methods and obf has many single-char methods, high similarity
	total_deob_methods := float64(len(deob.Methods))
	total_obf_methods := float64(len(obf.Methods))

	if total_deob_methods == 0 || total_obf_methods == 0 {
		return 0.0
	}

	deob_ratio := float64(deob_named_methods) / total_deob_methods
	obf_ratio := float64(obf_single_char_methods) / total_obf_methods

	// Perfect match if patterns are consistent
	if (deob_ratio > 0.8 && obf_ratio > 0.8) || (deob_ratio < 0.2 && obf_ratio < 0.2) {
		return 1.0
	}

	return 0.5 // Partial match
}

func (self *Scorer) calculate_functional_pattern_score(deob, obf *ClassInfo) float64 {
	score := 0.0

	// Check for utility class patterns
	if self.is_utility_class_pattern(deob, obf) {
		score += 0.3
	}

	// Check for data structure patterns
	if self.is_data_structure_pattern(deob, obf) {
		score += 0.3
	}

	// Check for interface implementation patterns
	if self.is_interface_implementation_pattern(deob, obf) {
		score += 0.4
	}

	return score
}

func (self *Scorer) is_utility_class_pattern(deob, obf *ClassInfo) bool {
	// Utility classes typically have:
	// - Few fields, many static methods
	// - No inheritance beyond Object
	// - Short class names in obfuscated code

	return len(deob.Fields) <= 3 && len(obf.Fields) <= 5 &&
		len(deob.Methods) >= 5 && len(obf.Methods) >= 3
}

func (self *Scorer) is_data_structure_pattern(deob, obf *ClassInfo) bool {
	// Data structures typically have:
	// - Many fields, few methods
	// - Array fields
	// - Node-like inheritance

	has_arrays := false
	for _, field := range deob.Fields {
		if strings.Contains(field.TypeName, "[") {
			has_arrays = true
			break
		}
	}

	return has_arrays && len(deob.Fields) > len(deob.Methods) &&
		len(obf.Fields) > len(obf.Methods)
}

func (self *Scorer) is_interface_implementation_pattern(deob, obf *ClassInfo) bool {
	// Interface implementations typically have:
	// - Methods that match interface signatures
	// - Similar inheritance patterns

	return len(deob.Interfaces) > 0 && len(deob.Methods) > len(deob.Fields) &&
		len(obf.Methods) > len(obf.Fields)
}

func (self *Scorer) access_modifiers_match(deob, obf *ClassInfo) bool {
	return len(deob.AccessModifiers) == len(obf.AccessModifiers)
}

func (self *Scorer) significant_size_difference(deob, obf *ClassInfo) bool {
	deob_size := len(deob.Fields) + len(deob.Methods)
	obf_size := len(obf.Fields) + len(obf.Methods)

	if deob_size == 0 || obf_size == 0 {
		return false
	}

	ratio := math.Max(float64(deob_size), float64(obf_size)) / math.Min(float64(deob_size), float64(obf_size))
	return ratio > MAX_SIZE_DIFFERENCE_RATIO
}

func (self *Scorer) resolve_interfaces(interfaces []string) []string {
	resolved := make([]string, len(interfaces))
	for i, iface := range interfaces {
		if resolved_name, exists := self.anchor_mappings[iface]; exists {
			resolved[i] = resolved_name
		} else {
			resolved[i] = iface
		}
	}
	return resolved
}

func (self *Scorer) resolve_superclass(superclass string) string {
	if resolved, exists := self.anchor_mappings[superclass]; exists {
		return resolved
	}
	return superclass
}

// Phase 3.2.2.1: Internal behavioral pattern scoring methods

// calculate_method_call_graph_similarity compares intra-class method call patterns
func (self *Scorer) calculate_method_call_graph_similarity(deob, obf *ClassInfo) float64 {
	if deob.InternalBehavior == nil || deob.InternalBehavior.MethodCallGraph == nil ||
		obf.InternalBehavior == nil || obf.InternalBehavior.MethodCallGraph == nil {
		return 0.0
	}

	deob_graph := deob.InternalBehavior.MethodCallGraph
	obf_graph := obf.InternalBehavior.MethodCallGraph

	// Compare method call frequencies
	total_similarity := 0.0
	comparisons := 0

	// Check if both have similar call patterns
	deob_total_calls := 0
	for _, freq := range deob_graph.CallFreq {
		deob_total_calls += freq
	}

	obf_total_calls := 0
	for _, freq := range obf_graph.CallFreq {
		obf_total_calls += freq
	}

	if deob_total_calls > 0 && obf_total_calls > 0 {
		// Similarity based on call density (calls per method)
		deob_density := float64(deob_total_calls) / float64(len(deob.Methods))
		obf_density := float64(obf_total_calls) / float64(len(obf.Methods))

		if deob_density > 0 && obf_density > 0 {
			density_ratio := math.Min(deob_density, obf_density) / math.Max(deob_density, obf_density)
			total_similarity += density_ratio * 0.6
			comparisons++
		}

		// Similarity based on recursion patterns
		deob_recursive_count := 0
		for _, isRecursive := range deob_graph.Recursion {
			if isRecursive {
				deob_recursive_count++
			}
		}

		obf_recursive_count := 0
		for _, isRecursive := range obf_graph.Recursion {
			if isRecursive {
				obf_recursive_count++
			}
		}

		if deob_recursive_count > 0 || obf_recursive_count > 0 {
			recursion_ratio := math.Min(float64(deob_recursive_count), float64(obf_recursive_count)) /
				math.Max(float64(deob_recursive_count), float64(obf_recursive_count))
			if math.Max(float64(deob_recursive_count), float64(obf_recursive_count)) > 0 {
				total_similarity += recursion_ratio * 0.4
				comparisons++
			}
		}
	}

	if comparisons == 0 {
		return 0.0
	}

	return total_similarity / float64(comparisons)
}

// calculate_state_manipulation_similarity compares field access patterns
func (self *Scorer) calculate_state_manipulation_similarity(deob, obf *ClassInfo) float64 {
	if deob.InternalBehavior == nil || deob.InternalBehavior.StatePatterns == nil ||
		obf.InternalBehavior == nil || obf.InternalBehavior.StatePatterns == nil {
		return 0.0
	}

	deob_state := deob.InternalBehavior.StatePatterns
	obf_state := obf.InternalBehavior.StatePatterns

	total_similarity := 0.0
	comparisons := 0

	// Compare field read patterns
	deob_reads := 0
	for _, count := range deob_state.FieldReads {
		deob_reads += count
	}

	obf_reads := 0
	for _, count := range obf_state.FieldReads {
		obf_reads += count
	}

	if deob_reads > 0 || obf_reads > 0 {
		read_ratio := math.Min(float64(deob_reads), float64(obf_reads)) /
			math.Max(float64(deob_reads), float64(obf_reads))
		if math.Max(float64(deob_reads), float64(obf_reads)) > 0 {
			total_similarity += read_ratio * 0.5
			comparisons++
		}
	}

	// Compare field write patterns
	deob_writes := 0
	for _, count := range deob_state.FieldWrites {
		deob_writes += count
	}

	obf_writes := 0
	for _, count := range obf_state.FieldWrites {
		obf_writes += count
	}

	if deob_writes > 0 || obf_writes > 0 {
		write_ratio := math.Min(float64(deob_writes), float64(obf_writes)) /
			math.Max(float64(deob_writes), float64(obf_writes))
		if math.Max(float64(deob_writes), float64(obf_writes)) > 0 {
			total_similarity += write_ratio * 0.5
			comparisons++
		}
	}

	if comparisons == 0 {
		return 0.0
	}

	return total_similarity / float64(comparisons)
}

// calculate_iteration_pattern_similarity compares loop and iteration usage
func (self *Scorer) calculate_iteration_pattern_similarity(deob, obf *ClassInfo) float64 {
	if deob.InternalBehavior == nil || deob.InternalBehavior.IterationPatterns == nil ||
		obf.InternalBehavior == nil || obf.InternalBehavior.IterationPatterns == nil {
		return 0.0
	}

	deob_iter := deob.InternalBehavior.IterationPatterns
	obf_iter := obf.InternalBehavior.IterationPatterns

	total_similarity := 0.0
	comparisons := 0

	// Compare for loops
	if deob_iter.ForLoops > 0 || obf_iter.ForLoops > 0 {
		for_ratio := math.Min(float64(deob_iter.ForLoops), float64(obf_iter.ForLoops)) /
			math.Max(float64(deob_iter.ForLoops), float64(obf_iter.ForLoops))
		if math.Max(float64(deob_iter.ForLoops), float64(obf_iter.ForLoops)) > 0 {
			total_similarity += for_ratio * 0.4
			comparisons++
		}
	}

	// Compare while loops
	if deob_iter.WhileLoops > 0 || obf_iter.WhileLoops > 0 {
		while_ratio := math.Min(float64(deob_iter.WhileLoops), float64(obf_iter.WhileLoops)) /
			math.Max(float64(deob_iter.WhileLoops), float64(obf_iter.WhileLoops))
		if math.Max(float64(deob_iter.WhileLoops), float64(obf_iter.WhileLoops)) > 0 {
			total_similarity += while_ratio * 0.3
			comparisons++
		}
	}

	// Compare iterator usage
	if deob_iter.IteratorUsage > 0 || obf_iter.IteratorUsage > 0 {
		iterator_ratio := math.Min(float64(deob_iter.IteratorUsage), float64(obf_iter.IteratorUsage)) /
			math.Max(float64(deob_iter.IteratorUsage), float64(obf_iter.IteratorUsage))
		if math.Max(float64(deob_iter.IteratorUsage), float64(obf_iter.IteratorUsage)) > 0 {
			total_similarity += iterator_ratio * 0.2
			comparisons++
		}
	}

	// Compare stream operations
	if deob_iter.StreamProcessing > 0 || obf_iter.StreamProcessing > 0 {
		stream_ratio := math.Min(float64(deob_iter.StreamProcessing), float64(obf_iter.StreamProcessing)) /
			math.Max(float64(deob_iter.StreamProcessing), float64(obf_iter.StreamProcessing))
		if math.Max(float64(deob_iter.StreamProcessing), float64(obf_iter.StreamProcessing)) > 0 {
			total_similarity += stream_ratio * 0.1
			comparisons++
		}
	}

	if comparisons == 0 {
		return 0.0
	}

	return total_similarity / float64(comparisons)
}

// calculate_semantic_method_similarity compares method naming patterns
func (self *Scorer) calculate_semantic_method_similarity(deob, obf *ClassInfo) float64 {
	if deob.InternalBehavior == nil || deob.InternalBehavior.SemanticPatterns == nil ||
		obf.InternalBehavior == nil || obf.InternalBehavior.SemanticPatterns == nil {
		return 0.0
	}

	deob_semantic := deob.InternalBehavior.SemanticPatterns
	obf_semantic := obf.InternalBehavior.SemanticPatterns

	total_similarity := 0.0
	comparisons := 0

	// Compare getter methods (deobfuscated should have meaningful names, obfuscated won't)
	if len(deob_semantic.GetterMethods) > 0 || len(obf_semantic.GetterMethods) > 0 {
		// For obfuscated code, we won't detect getters by name, so similarity is low
		// This helps distinguish deobfuscated from obfuscated classes
		if len(deob_semantic.GetterMethods) > 0 && len(obf_semantic.GetterMethods) == 0 {
			total_similarity += 0.8 // High similarity for this pattern
		} else if len(deob_semantic.GetterMethods) == 0 && len(obf_semantic.GetterMethods) == 0 {
			total_similarity += 0.9 // Both have no getters (possibly obfuscated)
		} else {
			total_similarity += 0.2 // Mismatch
		}
		comparisons++
	}

	// Compare setter methods
	if len(deob_semantic.SetterMethods) > 0 || len(obf_semantic.SetterMethods) > 0 {
		if len(deob_semantic.SetterMethods) > 0 && len(obf_semantic.SetterMethods) == 0 {
			total_similarity += 0.8
		} else if len(deob_semantic.SetterMethods) == 0 && len(obf_semantic.SetterMethods) == 0 {
			total_similarity += 0.9
		} else {
			total_similarity += 0.2
		}
		comparisons++
	}

	// Compare builder/factory methods
	builder_factory_deob := len(deob_semantic.BuilderMethods) + len(deob_semantic.FactoryMethods)
	builder_factory_obf := len(obf_semantic.BuilderMethods) + len(obf_semantic.FactoryMethods)

	if builder_factory_deob > 0 || builder_factory_obf > 0 {
		ratio := math.Min(float64(builder_factory_deob), float64(builder_factory_obf)) /
			math.Max(float64(builder_factory_deob), float64(builder_factory_obf))
		if math.Max(float64(builder_factory_deob), float64(builder_factory_obf)) > 0 {
			total_similarity += ratio * 0.6
			comparisons++
		}
	}

	// Compare event handlers
	if len(deob_semantic.EventHandlers) > 0 || len(obf_semantic.EventHandlers) > 0 {
		ratio := math.Min(float64(len(deob_semantic.EventHandlers)), float64(len(obf_semantic.EventHandlers))) /
			math.Max(float64(len(deob_semantic.EventHandlers)), float64(len(obf_semantic.EventHandlers)))
		if math.Max(float64(len(deob_semantic.EventHandlers)), float64(len(obf_semantic.EventHandlers))) > 0 {
			total_similarity += ratio * 0.4
			comparisons++
		}
	}

	// Compare utility methods
	if len(deob_semantic.UtilityMethods) > 0 || len(obf_semantic.UtilityMethods) > 0 {
		ratio := math.Min(float64(len(deob_semantic.UtilityMethods)), float64(len(obf_semantic.UtilityMethods))) /
			math.Max(float64(len(deob_semantic.UtilityMethods)), float64(len(obf_semantic.UtilityMethods)))
		if math.Max(float64(len(deob_semantic.UtilityMethods)), float64(len(obf_semantic.UtilityMethods))) > 0 {
			total_similarity += ratio * 0.3
			comparisons++
		}
	}

	if comparisons == 0 {
		return 0.0
	}

	return total_similarity / float64(comparisons)
}

// Phase 3.2.2.2.2: Data structure classification scoring methods

// calculateGraphicsArrayBonus calculates domain-specific bonuses for graphics classes
func (self *Scorer) calculateGraphicsArrayBonus(deob, obf *ClassInfo) float64 {
	// Only the deobfuscated class has source code and internal behavior analysis
	if deob.InternalBehavior == nil || deob.InternalBehavior.ArrayStructure == nil {
		return 0.0
	}

	deobSig := deob.InternalBehavior.ArrayStructure
	bonus := 0.0

	// Apply bonus based on array structure classification type
	switch deobSig.StructureType {
	case VertexData:
		bonus += WEIGHT_3D_VERTEX_DATA
	case TextureData:
		bonus += WEIGHT_TEXTURE_DATA
	case WorldData:
		bonus += WEIGHT_WORLD_DATA
	}

	return bonus
}

// calculateVertexSimilarity compares vertex data structure signatures
func (self *Scorer) calculateVertexSimilarity(deobSig, obfSig *ArrayStructureSignature) float64 {
	similarity := 0.0
	comparisons := 0

	// Compare dimensionality (3D vertex data)
	if deobSig.DimCount >= 2 && obfSig.DimCount >= 2 {
		similarity += 0.3
		comparisons++
	}

	// Compare domain-specific patterns
	if deobSig.DomainSpecific["vertex_arrays"] && obfSig.DomainSpecific["vertex_arrays"] {
		similarity += 0.4
		comparisons++
	}
	if deobSig.DomainSpecific["3d_transforms"] && obfSig.DomainSpecific["3d_transforms"] {
		similarity += 0.2
		comparisons++
	}
	if deobSig.DomainSpecific["face_indices"] && obfSig.DomainSpecific["face_indices"] {
		similarity += 0.1
		comparisons++
	}

	if comparisons == 0 {
		return 0.0
	}

	return similarity / float64(comparisons)
}

// calculateTextureSimilarity compares texture data structure signatures
func (self *Scorer) calculateTextureSimilarity(deobSig, obfSig *ArrayStructureSignature) float64 {
	similarity := 0.0
	comparisons := 0

	// Compare dimensionality (2D texture data)
	if deobSig.DimCount >= 1 && obfSig.DimCount >= 1 {
		similarity += 0.3
		comparisons++
	}

	// Compare domain-specific patterns
	if deobSig.DomainSpecific["pixel_manipulation"] && obfSig.DomainSpecific["pixel_manipulation"] {
		similarity += 0.4
		comparisons++
	}
	if deobSig.DomainSpecific["image_dimensions"] && obfSig.DomainSpecific["image_dimensions"] {
		similarity += 0.2
		comparisons++
	}
	if deobSig.DomainSpecific["color_operations"] && obfSig.DomainSpecific["color_operations"] {
		similarity += 0.1
		comparisons++
	}

	if comparisons == 0 {
		return 0.0
	}

	return similarity / float64(comparisons)
}

// calculateWorldSimilarity compares world data structure signatures
func (self *Scorer) calculateWorldSimilarity(deobSig, obfSig *ArrayStructureSignature) float64 {
	similarity := 0.0
	comparisons := 0

	// Compare dimensionality (2D/3D world data)
	if deobSig.DimCount >= 2 && obfSig.DimCount >= 2 {
		similarity += 0.3
		comparisons++
	}

	// Compare domain-specific patterns
	if deobSig.DomainSpecific["tile_grid"] && obfSig.DomainSpecific["tile_grid"] {
		similarity += 0.3
		comparisons++
	}
	if deobSig.DomainSpecific["chunk_loading"] && obfSig.DomainSpecific["chunk_loading"] {
		similarity += 0.2
		comparisons++
	}
	if deobSig.DomainSpecific["heightmap_operations"] && obfSig.DomainSpecific["heightmap_operations"] {
		similarity += 0.2
		comparisons++
	}

	if comparisons == 0 {
		return 0.0
	}

	return similarity / float64(comparisons)
}

// Phase 1.1A: Magic Constants & Static Fields Pattern
// calculateMagicConstantsScore detects distinctive magic numbers and constants
func (self *Scorer) calculateMagicConstantsScore(deob, obf *ClassInfo) float64 {
	score := 0.0

	// TextClass: Base-37 string hashing methods
	if self.hasTextClassMethods(deob) && self.hasTextClassMethods(obf) {
		score += 8.0 // TextClass specific methods: longForName, nameForLong
	} else if self.hasTextClassMethods(deob) || self.hasTextClassMethods(obf) {
		score += 6.0 // Partial match - at least one class has TextClass-like methods
	}

	// ISAACRandomGen: Cryptographic constants and array sizes
	if self.hasISAACPatterns(deob) && self.hasISAACPatterns(obf) {
		score += 10.0 // ISAAC specific: 256-element arrays, getNextKey method
	}

	// Skills: 25 skill count and skill name arrays
	if self.hasSkillsPatterns(deob) && self.hasSkillsPatterns(obf) {
		score += 8.0 // Skills class with 25-element arrays and skill-related names
	}

	// SizeConstants: Large static arrays with specific patterns
	if self.hasSizeConstantsPatterns(deob) && self.hasSizeConstantsPatterns(obf) {
		score += 8.0 // SizeConstants: large static final arrays
	}

	return score
}

// Helper methods for magic constant detection

func (self *Scorer) hasTextClassMethods(class *ClassInfo) bool {
	hasLongForName := false
	hasNameForLong := false

	for _, method := range class.Methods {
		if strings.Contains(method.Name, "longForName") || strings.Contains(method.Name, "method585") {
			hasLongForName = true
		}
		if strings.Contains(method.Name, "nameForLong") || strings.Contains(method.Name, "method586") {
			hasNameForLong = true
		}
	}

	// TextClass has string conversion methods - be flexible
	return hasLongForName || hasNameForLong
}

func (self *Scorer) hasISAACPatterns(class *ClassInfo) bool {
	// ISAAC has getNextKey method and typically 2 int[] fields
	hasGetNextKey := false
	intArrayCount := 0

	for _, method := range class.Methods {
		if strings.Contains(method.Name, "getNextKey") {
			hasGetNextKey = true
		}
	}

	for _, field := range class.Fields {
		if strings.Contains(field.TypeName, "int[") || strings.Contains(field.TypeName, "[I") {
			intArrayCount++
		}
	}

	// ISAAC typically has getNextKey method and 2 int arrays (memory, results)
	return hasGetNextKey && intArrayCount >= 2
}

func (self *Scorer) hasSkillsPatterns(class *ClassInfo) bool {
	// Skills has 25-element arrays and skill-related field names
	hasSkillCount := false
	hasSkillNames := false

	// Check for skillsCount field or similar
	for _, field := range class.Fields {
		if strings.Contains(field.Name, "skillsCount") ||
			(strings.Contains(field.TypeName, "String[") && strings.Contains(field.TypeName, "25")) {
			hasSkillCount = true
		}
		if strings.Contains(field.Name, "skillNames") ||
			strings.Contains(field.Name, "skillEnabled") {
			hasSkillNames = true
		}
	}

	// Check method names for skill-related functionality
	methodCount := 0
	for _, method := range class.Methods {
		if strings.Contains(method.Name, "skill") ||
			strings.Contains(method.Name, "Skill") {
			methodCount++
		}
	}

	return hasSkillCount || (hasSkillNames && methodCount >= 2)
}

func (self *Scorer) hasSizeConstantsPatterns(class *ClassInfo) bool {
	// SizeConstants has large static final arrays
	staticFinalCount := 0
	largeArrayCount := 0

	for _, field := range class.Fields {
		// Look for static final modifiers (if available)
		isStatic := false
		isFinal := false
		for _, mod := range field.AccessModifiers {
			if mod == "static" {
				isStatic = true
			}
			if mod == "final" {
				isFinal = true
			}
		}

		if isStatic && isFinal {
			staticFinalCount++
		}

		// Look for large arrays (256+ elements or specific array names)
		if strings.Contains(field.Name, "anIntArray552") ||
			strings.Contains(field.Name, "packetSizes") ||
			strings.Contains(field.TypeName, "int[256]") ||
			strings.Contains(field.TypeName, "[I") {
			largeArrayCount++
		}
	}

	// SizeConstants typically has 2 large static final arrays
	return staticFinalCount >= 2 || largeArrayCount >= 2
}
