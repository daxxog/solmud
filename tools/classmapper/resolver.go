package main

import (
	"fmt"
	"os"
	"strings"
)

type Resolver struct {
	anchors          map[string]string
	reverse_anchors  map[string]string
	scorer           *Scorer
	deob_inheritance map[string][]string // superclass -> [subclasses]
	obf_inheritance  map[string][]string // superclass -> [subclasses]
	deob_references  map[string][]string // class -> [classes it references]
	obf_references   map[string][]string // class -> [classes it references]
}

func NewResolver(anchors map[string]string) *Resolver {
	reverse_anchors := make(map[string]string)
	for deob, obf := range anchors {
		reverse_anchors[obf] = deob
	}

	return &Resolver{
		anchors:          anchors,
		reverse_anchors:  reverse_anchors,
		scorer:           NewScorer(anchors),
		deob_inheritance: make(map[string][]string),
		obf_inheritance:  make(map[string][]string),
		deob_references:  make(map[string][]string),
		obf_references:   make(map[string][]string),
	}
}

func (self *Resolver) buildInheritanceMaps(deob_classes []ClassInfo, obf_classes []ClassInfo) {
	// Build deobfuscated inheritance map
	for _, class := range deob_classes {
		if class.Superclass != "" {
			self.deob_inheritance[class.Superclass] = append(self.deob_inheritance[class.Superclass], class.Name)
		}
		// Build reference maps based on field types and cross-references
		if class.CrossReferences != nil {
			for className := range class.CrossReferences.UniqueTargets {
				self.deob_references[class.Name] = append(self.deob_references[class.Name], className)
			}
		}
		// Also include field type references as fallback
		for _, field := range class.Fields {
			if self.is_known_class_reference(field.TypeName) {
				found := false
				for _, existing := range self.deob_references[class.Name] {
					if existing == field.TypeName {
						found = true
						break
					}
				}
				if !found {
					self.deob_references[class.Name] = append(self.deob_references[class.Name], field.TypeName)
				}
			}
		}
	}

	// Build obfuscated inheritance map
	for _, class := range obf_classes {
		if class.Superclass != "" {
			self.obf_inheritance[class.Superclass] = append(self.obf_inheritance[class.Superclass], class.Name)
		}
		// Build reference maps based on field types
		for _, field := range class.Fields {
			if self.is_known_class_reference(field.TypeName) {
				self.obf_references[class.Name] = append(self.obf_references[class.Name], field.TypeName)
			}
		}
	}
}

func (self *Resolver) is_known_class_reference(typeName string) bool {
	// Check if this type references a known mapped class
	// Remove array brackets and check if it's a known class
	cleanType := strings.TrimPrefix(typeName, "[")
	cleanType = strings.TrimPrefix(cleanType, "L")
	cleanType = strings.TrimSuffix(cleanType, ";")

	_, isMapped := self.anchors[cleanType]
	_, isReverseMapped := self.reverse_anchors[cleanType]
	return isMapped || isReverseMapped
}

func (self *Resolver) ResolveAll(deob_classes []ClassInfo, obf_classes []ClassInfo) []MatchResult {
	var results []MatchResult

	// Build inheritance maps for enhanced matching
	self.buildInheritanceMaps(deob_classes, obf_classes)

	fmt.Fprintln(os.Stderr, "Pass 1: Matching anchor classes...")
	for _, deob_class := range deob_classes {
		if obf_name, exists := self.anchors[deob_class.Name]; exists {
			obf_class := self.find_obfuscated_class(obf_classes, obf_name)
			if obf_class != nil {
				result := self.create_match(&deob_class, obf_class)
				result.ConfidenceScore = 100.0
				result.Details = "Exact match - anchor class"
				results = append(results, result)
			}
		}
	}

	fmt.Fprintln(os.Stderr, "Pass 2: Matching by inheritance hierarchy...")
	pending_deob := self.remove_matched(deob_classes, results)
	pending_obf := self.remove_matched(obf_classes, results)

	// Enhanced inheritance chain matching
	fmt.Fprintln(os.Stderr, "  Propagating inheritance chains...")
	inherited_matches := self.propagateInheritanceChains(pending_deob, pending_obf)
	results = append(results, inherited_matches...)

	// Update pending lists after inheritance propagation
	pending_deob = self.remove_matched(deob_classes, results)
	pending_obf = self.remove_matched(obf_classes, results)

	for _, deob_class := range pending_deob {
		best_match, best_score := self.find_best_match(&deob_class, pending_obf)
		if best_match != nil && best_score >= HIGH_CONFIDENCE_THRESHOLD {
			result := self.create_match(&deob_class, best_match)
			result.ConfidenceScore = best_score
			result.Details = fmt.Sprintf("Matched by inheritance - score %.2f", best_score)
			results = append(results, result)
		}
	}

	fmt.Fprintln(os.Stderr, "Pass 3: Matching by signatures and cross-references...")
	pending_deob = self.remove_matched(deob_classes, results)
	pending_obf = self.remove_matched(obf_classes, results)

	// Enhanced cross-reference matching
	cross_ref_matches := self.find_cross_reference_matches(pending_deob, pending_obf)
	results = append(results, cross_ref_matches...)

	// Update pending lists after cross-reference matching
	pending_deob = self.remove_matched(deob_classes, results)
	pending_obf = self.remove_matched(obf_classes, results)

	for _, deob_class := range pending_deob {
		best_match, best_score := self.find_best_match(&deob_class, pending_obf)
		if best_match != nil && best_score >= MIN_CONFIDENCE_THRESHOLD {
			result := self.create_match(&deob_class, best_match)
			result.ConfidenceScore = best_score
			result.Details = fmt.Sprintf("Matched by signature - score %.2f", best_score)
			results = append(results, result)
		}
	}

	return results
}

func (self *Resolver) find_obfuscated_class(obf_classes []ClassInfo, name string) *ClassInfo {
	for _, obf_class := range obf_classes {
		if obf_class.Name == name {
			return &obf_class
		}
	}
	return nil
}

func (self *Resolver) propagateInheritanceChains(deob_classes []ClassInfo, obf_classes []ClassInfo) []MatchResult {
	var results []MatchResult

	// For each anchor mapping, try to propagate to subclasses
	for deob_super, obf_super := range self.anchors {
		// Find deobfuscated subclasses of this superclass that are still pending
		deob_subclasses := self.deob_inheritance[deob_super]
		obf_subclasses := self.obf_inheritance[obf_super]

		// Filter to only pending classes
		var pending_deob_subs []string
		for _, subclass := range deob_subclasses {
			if self.find_deobfuscated_class(deob_classes, subclass) != nil {
				pending_deob_subs = append(pending_deob_subs, subclass)
			}
		}

		var pending_obf_subs []string
		for _, subclass := range obf_subclasses {
			if self.find_obfuscated_class(obf_classes, subclass) != nil {
				pending_obf_subs = append(pending_obf_subs, subclass)
			}
		}

		if len(pending_deob_subs) > 0 && len(pending_obf_subs) > 0 {
			// Try to match subclasses by relative position in hierarchy
			matched_subs := self.matchSubclassesByHierarchy(pending_deob_subs, pending_obf_subs, deob_classes, obf_classes)
			results = append(results, matched_subs...)
		}
	}

	return results
}

func (self *Resolver) matchSubclassesByHierarchy(deob_subclasses []string, obf_subclasses []string,
	deob_classes []ClassInfo, obf_classes []ClassInfo) []MatchResult {

	var results []MatchResult

	// Simple heuristic: if counts match, try to match by complexity (method/field count)
	if len(deob_subclasses) == len(obf_subclasses) {
		// Sort both lists by complexity (method count + field count)
		deob_sorted := self.sortClassesByComplexity(deob_subclasses, deob_classes)
		obf_sorted := self.sortClassesByComplexity(obf_subclasses, obf_classes)

		// Match by relative complexity
		for i, deob_name := range deob_sorted {
			if i >= len(obf_sorted) {
				break
			}
			obf_name := obf_sorted[i]

			// Find the actual class objects
			deob_class := self.find_deobfuscated_class(deob_classes, deob_name)
			obf_class := self.find_obfuscated_class(obf_classes, obf_name)

			if deob_class != nil && obf_class != nil {
				// Since these are subclasses of anchor-matched classes, we can be confident
				result := self.create_match(deob_class, obf_class)
				result.ConfidenceScore = 85.0 // High confidence for inheritance chain match
				result.Details = fmt.Sprintf("Inherited from anchor class hierarchy")
				results = append(results, result)
			}
		}
	}

	return results
}

func (self *Resolver) sortClassesByComplexity(class_names []string, classes []ClassInfo) []string {
	type classComplexity struct {
		name       string
		complexity int
	}

	var complexities []classComplexity
	for _, name := range class_names {
		class := self.find_by_name(classes, name)
		if class != nil {
			complexity := len(class.Methods) + len(class.Fields)
			complexities = append(complexities, classComplexity{name, complexity})
		}
	}

	// Sort by complexity descending
	for i := 0; i < len(complexities)-1; i++ {
		for j := i + 1; j < len(complexities); j++ {
			if complexities[i].complexity < complexities[j].complexity {
				complexities[i], complexities[j] = complexities[j], complexities[i]
			}
		}
	}

	var result []string
	for _, c := range complexities {
		result = append(result, c.name)
	}
	return result
}

func (self *Resolver) find_deobfuscated_class(classes []ClassInfo, name string) *ClassInfo {
	for _, class := range classes {
		if class.Name == name {
			return &class
		}
	}
	return nil
}

func (self *Resolver) find_by_name(classes []ClassInfo, name string) *ClassInfo {
	for _, class := range classes {
		if class.Name == name {
			return &class
		}
	}
	return nil
}

func (self *Resolver) classes_have_similar_inheritance(deob, obf *ClassInfo) bool {
	// For inheritance chain propagation, we want classes that both extend anchor classes
	// Check if both classes extend anchor-mapped superclasses
	_, deob_super_is_anchor := self.anchors[deob.Superclass]
	_, obf_super_is_anchor := self.reverse_anchors[obf.Superclass]

	return deob_super_is_anchor && obf_super_is_anchor
}

func (self *Resolver) find_cross_reference_matches(deob_classes []ClassInfo, obf_classes []ClassInfo) []MatchResult {
	var results []MatchResult

	// For each pending deobfuscated class, look for obfuscated classes with similar reference patterns
	for _, deob_class := range deob_classes {
		// Use cross-reference data if available
		if deob_class.CrossReferences != nil && len(deob_class.CrossReferences.UniqueTargets) > 0 {
			var refNames []string
			for className := range deob_class.CrossReferences.UniqueTargets {
				refNames = append(refNames, className)
			}
			best_match := self.find_best_cross_reference_match(&deob_class, refNames, obf_classes)
			if best_match != nil {
				result := self.create_match(&deob_class, best_match)
				result.ConfidenceScore = 75.0 // Moderate confidence for cross-reference match
				result.Details = "Matched by cross-reference analysis"
				results = append(results, result)
			}
		}
	}

	return results
}

func (self *Resolver) find_best_cross_reference_match(deob_class *ClassInfo, deob_refs []string, obf_classes []ClassInfo) *ClassInfo {
	// Look for obfuscated classes that reference the same mapped classes
	for _, obf_class := range obf_classes {
		obf_refs := self.obf_references[obf_class.Name]
		if len(obf_refs) == 0 {
			continue
		}

		// Check if they reference similar known classes
		if self.references_are_similar(deob_refs, obf_refs) {
			// Additional validation: similar size/complexity
			if self.classes_have_similar_complexity(deob_class, &obf_class) {
				return &obf_class
			}
		}
	}
	return nil
}

func (self *Resolver) references_are_similar(deob_refs []string, obf_refs []string) bool {
	if len(deob_refs) == 0 || len(obf_refs) == 0 {
		return false
	}

	// Convert references to mapped class names for comparison
	mapped_deob_refs := make([]string, 0)
	for _, ref := range deob_refs {
		if mapped, exists := self.anchors[ref]; exists {
			mapped_deob_refs = append(mapped_deob_refs, mapped)
		} else {
			mapped_deob_refs = append(mapped_deob_refs, ref)
		}
	}

	mapped_obf_refs := make([]string, 0)
	for _, ref := range obf_refs {
		if mapped, exists := self.reverse_anchors[ref]; exists {
			mapped_obf_refs = append(mapped_obf_refs, mapped)
		} else {
			mapped_obf_refs = append(mapped_obf_refs, ref)
		}
	}

	// Check for overlap
	overlap := 0
	for _, dref := range mapped_deob_refs {
		for _, oref := range mapped_obf_refs {
			if dref == oref {
				overlap++
				break
			}
		}
	}

	return overlap > 0 && overlap >= len(deob_refs)/2 // At least half the references match
}

func (self *Resolver) classes_have_similar_complexity(deob, obf *ClassInfo) bool {
	deob_complexity := len(deob.Fields) + len(deob.Methods)
	obf_complexity := len(obf.Fields) + len(obf.Methods)

	// Allow 30% difference in complexity
	ratio := float64(deob_complexity) / float64(obf_complexity)
	return ratio >= 0.7 && ratio <= 1.3
}

func (self *Resolver) find_best_match(deob_class *ClassInfo, obf_classes []ClassInfo) (*ClassInfo, float64) {
	var best_match *ClassInfo
	best_score := 0.0

	for _, obf_class := range obf_classes {
		breakdown := self.scorer.CalculateScore(deob_class, &obf_class)
		score := breakdown.InterfaceMatch + breakdown.SuperclassMatch +
			breakdown.FieldCountMatch + breakdown.FieldSimilarity +
			breakdown.MethodCountMatch + breakdown.MethodSimilarity +
			breakdown.ConstructorMatch + breakdown.AccessMatch +
			breakdown.SizePenalty

		if score > best_score {
			best_score = score
			best_match = &obf_class
		}
	}

	return best_match, best_score
}

func (self *Resolver) create_match(deob, obf *ClassInfo) MatchResult {
	breakdown := self.scorer.CalculateScore(deob, obf)
	total_score := breakdown.InterfaceMatch + breakdown.SuperclassMatch +
		breakdown.FieldCountMatch + breakdown.FieldSimilarity + breakdown.FieldPatternSim +
		breakdown.MethodCountMatch + breakdown.MethodSimilarity + breakdown.MethodNameSim +
		breakdown.ConstructorMatch + breakdown.AccessMatch + breakdown.FunctionalPattern +
		breakdown.SizePenalty

	return MatchResult{
		DeobfuscatedClass: deob.Name,
		ObfuscatedClass:   obf.Name,
		ConfidenceScore:   total_score,
		ScoreBreakdown:    breakdown,
		Details:           "",
	}
}

func (self *Resolver) remove_matched(original []ClassInfo, matches []MatchResult) []ClassInfo {
	matched_names := make(map[string]bool)
	for _, match := range matches {
		matched_names[match.DeobfuscatedClass] = true
		matched_names[match.ObfuscatedClass] = true
	}

	var remaining []ClassInfo
	for _, class := range original {
		if !matched_names[class.Name] {
			remaining = append(remaining, class)
		}
	}

	return remaining
}

func (self *Resolver) validate_and_refine_matches(results []MatchResult) []MatchResult {
	fmt.Fprintln(os.Stderr, "Validating and refining matches...")

	// Build mapping lookup
	mappings := make(map[string]string)
	reverse_mappings := make(map[string]string)
	for _, result := range results {
		if result.ObfuscatedClass != "" {
			mappings[result.DeobfuscatedClass] = result.ObfuscatedClass
			if existing, exists := reverse_mappings[result.ObfuscatedClass]; exists {
				fmt.Fprintf(os.Stderr, "  Conflict detected: %s and %s both map to %s\n",
					existing, result.DeobfuscatedClass, result.ObfuscatedClass)
				// Keep the higher confidence match
				// (This is a simple conflict resolution - could be enhanced)
			} else {
				reverse_mappings[result.ObfuscatedClass] = result.DeobfuscatedClass
			}
		}
	}

	// For now, return all results - validation could be enhanced
	fmt.Fprintf(os.Stderr, "  Validation complete: %d mappings\n", len(results))
	return results
}
