package main

import (
	"math"
	"strings"
)

const (
	WEIGHT_INTERFACE_MATCH    = 20.0
	WEIGHT_SUPERCLASS_MATCH   = 25.0
	WEIGHT_FIELD_COUNT_MATCH  = 5.0
	WEIGHT_FIELD_TYPE_SIM     = 10.0
	WEIGHT_FIELD_PATTERN_SIM  = 8.0 // New: field type patterns
	WEIGHT_METHOD_COUNT_MATCH = 5.0
	WEIGHT_METHOD_SIG_SIM     = 20.0
	WEIGHT_METHOD_NAME_SIM    = 7.0 // New: method name similarity
	WEIGHT_CONSTRUCTOR_MATCH  = 5.0
	WEIGHT_ACCESS_MOD_MATCH   = 5.0
	WEIGHT_FUNCTIONAL_PATTERN = 12.0 // New: functional grouping patterns
	SIZE_DIFFERENCE_PENALTY   = 10.0
	MAX_SIZE_DIFFERENCE_RATIO = 1.5
)

type ScoreBreakdown struct {
	InterfaceMatch    float64
	SuperclassMatch   float64
	FieldCountMatch   float64
	FieldSimilarity   float64
	FieldPatternSim   float64
	MethodCountMatch  float64
	MethodSimilarity  float64
	MethodNameSim     float64
	ConstructorMatch  float64
	AccessMatch       float64
	FunctionalPattern float64
	SizePenalty       float64
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
