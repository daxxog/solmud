package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	HIGH_CONFIDENCE_THRESHOLD = 85.0
	MIN_CONFIDENCE_THRESHOLD  = 65.0
	SIZE_DIFFERENCE_PENALTY   = 10.0
	MAX_SIZE_DIFFERENCE_RATIO = 1.5

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
)

var ANCHOR_MAPPINGS = map[string]string{
	"sign/signlink": "sign/signlink",
	"Node":          "PKVMXVTO",
	"NodeSub":       "PPOHBEGB",
	"RSApplet":      "KHACHIFW",
	"client":        "client",
	// Phase 1: Enhanced inheritance chain anchors
	"Animable": "XHHRODPC",
	"Entity":   "GQOSZKJC",
	"Model":    "ZKARKDQW",
	"Stream":   "CRRWDRTI", // Reassigned from MBMGIXGO
	// Phase 2: Immediate high-confidence matches
	"RSFrame": "FPVKJCAH", // Extends Frame, identical constructor signatures
	// Phase 2: File loading pattern matches
	"IDK":       "TAVAECED", // Loads "idk.dat"
	"VarBit":    "SXYSOXTR", // Loads "varbit.dat"
	"ItemDef":   "DJRMEMXO", // Loads "obj.dat"
	"ObjectDef": "YZDBYLRM", // Loads "loc.dat"
	"EntityDef": "CKDEJADD", // Loads "npc.dat"
	"Animation": "LKGEGIEW", // Loads "seq.dat"
	"Flo":       "MNHKFPQO", // Loads "flo.dat"
	"Varp":      "VGXVBFVC", // Loads "varp.dat"
	"SpotAnim":  "MUDLUUBC", // Loads "spotanim.dat"
	// Phase 2: Inheritance and complexity matches
	"TextDrawingArea": "YXVQXWYR", // Extends DrawingArea, text rendering
	"WorldController": "MBMGIXGO", // 91 methods, world management
	"RSInterface":     "RKAYAFDQ", // 51 methods, interface management
	"OnDemandFetcher": "CRRWDRTI", // 46 methods, network fetching (reassigned from Stream)
}

const (
	ERR_MISSING_JAVAP      = "javap not found in PATH"
	ERR_INVALID_CLASS_FILE = "invalid or corrupt class file: %s"
	ERR_INVALID_BYTECODE   = "invalid bytecode file: %s"
	ERR_NO_CLASSES_FOUND   = "no class files found in: %s"
	ERR_PARSE_FAILED       = "failed to parse class: %s"
	ERR_WRITE_FAILED       = "failed to write output: %w"
	ERR_DIR_NOT_FOUND      = "directory not found: %s"
	ERR_JAVAP_EXECUTION    = "failed to execute javap: %w"
)

type ClassInfo struct {
	Name            string
	PackageName     string
	AccessModifiers []string
	Superclass      string
	Interfaces      []string
	Fields          []FieldInfo
	Methods         []MethodInfo
	Constructors    []MethodInfo
	Source          ClassSource
}

type ClassSource int

const (
	SourceDeobfuscated ClassSource = iota
	SourceObfuscated
)

type FieldInfo struct {
	AccessModifiers []string
	TypeName        string
	Name            string
}

type MethodInfo struct {
	AccessModifiers []string
	ReturnType      string
	Name            string
	Parameters      []string
}

type MatchResult struct {
	DeobfuscatedClass string
	ObfuscatedClass   string
	ConfidenceScore   float64
	ScoreBreakdown    ScoreBreakdown
	Details           string
}

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

type BytecodeParser struct {
	regex_class_decl *regexp.Regexp
	regex_extends    *regexp.Regexp
	regex_implements *regexp.Regexp
	regex_field      *regexp.Regexp
	regex_method     *regexp.Regexp
}

func NewBytecodeParser() *BytecodeParser {
	return &BytecodeParser{
		regex_class_decl: regexp.MustCompile(`^(public|final|abstract)?\s*(class|interface)\s+(\w+)`),
		regex_extends:    regexp.MustCompile(`extends (\w+)`),
		regex_implements: regexp.MustCompile(`implements ([\w,\s\.]+)`),
		regex_field:      regexp.MustCompile(`^private|public|protected?\s+(?:static\s+)?(?:final\s+)?([\[\]\w]+)\s+(\w+)`),
		regex_method:     regexp.MustCompile(`^public|private|protected?\s+(?:static\s+)?([\[\]\w]+)\s+(\w+)\(([^)]*)\)`),
	}
}

func (self *BytecodeParser) ParseAll(source_path string) ([]ClassInfo, error) {
	files, err := os.ReadDir(source_path)
	if err != nil {
		return nil, fmt.Errorf(ERR_DIR_NOT_FOUND, source_path)
	}

	var classes []ClassInfo

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".bytecode.txt") {
			continue
		}

		file_path := filepath.Join(source_path, file.Name())
		class_info, err := self.parse_file(file_path)
		if err != nil {
			return nil, fmt.Errorf(ERR_PARSE_FAILED, file_path)
		}

		classes = append(classes, *class_info)
	}

	if len(classes) == 0 {
		panic(fmt.Sprintf(ERR_NO_CLASSES_FOUND, source_path))
	}

	return classes, nil
}

func (self *BytecodeParser) parse_file(file_path string) (*ClassInfo, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	class_info := &ClassInfo{
		Source: SourceObfuscated,
	}

	for scanner.Scan() {
		line := scanner.Text()

		if matches := self.regex_class_decl.FindStringSubmatch(line); matches != nil {
			class_info.Name = matches[3]
		}

		if matches := self.regex_extends.FindStringSubmatch(line); matches != nil {
			class_info.Superclass = strings.TrimSpace(matches[1])
		}

		if matches := self.regex_implements.FindStringSubmatch(line); matches != nil {
			interfaces := strings.Split(matches[1], ",")
			for _, iface := range interfaces {
				class_info.Interfaces = append(class_info.Interfaces, strings.TrimSpace(iface))
			}
		}

		if matches := self.regex_field.FindStringSubmatch(line); matches != nil {
			class_info.Fields = append(class_info.Fields, FieldInfo{
				TypeName: matches[1],
				Name:     matches[2],
			})
		}

		if matches := self.regex_method.FindStringSubmatch(line); matches != nil {
			class_info.Methods = append(class_info.Methods, MethodInfo{
				ReturnType: matches[1],
				Name:       matches[2],
				Parameters: self.parse_parameters(matches[3]),
			})
		}
	}

	return class_info, nil
}

func (self *BytecodeParser) parse_parameters(param_str string) []string {
	if strings.TrimSpace(param_str) == "" {
		return []string{}
	}
	return strings.Split(param_str, ",")
}

type JavapParser struct {
	regex_class_decl *regexp.Regexp
	regex_extends    *regexp.Regexp
	regex_implements *regexp.Regexp
	regex_field      *regexp.Regexp
	regex_method     *regexp.Regexp
	cacheManager     *CacheManager
	cacheHits        int
	cacheMisses      int
}

func NewJavapParser(cacheManager *CacheManager) *JavapParser {
	return &JavapParser{
		regex_class_decl: regexp.MustCompile(`^(public|final|abstract)?\s*(class|interface)\s+(\w+)`),
		regex_extends:    regexp.MustCompile(`extends\s+(\w+)`),
		regex_implements: regexp.MustCompile(`implements\s+([\w,\.]+)`),
		regex_field:      regexp.MustCompile(`^\s*(private|public|protected)?\s*(static\s+)?(final\s+)?([\[\]\w]+)\s+(\w+)`),
		regex_method:     regexp.MustCompile(`^\s*(public|private|protected)?\s*(static\s+)?(final\s+)?([\[\]\w]+)\s+(\w+)\(([^)]*)\)`),
		cacheManager:     cacheManager,
	}
}

func (self *JavapParser) ParseAll(source_path string) ([]ClassInfo, error) {
	if _, err := exec.LookPath("javap"); err != nil {
		panic(ERR_MISSING_JAVAP)
	}

	files, err := os.ReadDir(source_path)
	if err != nil {
		return nil, fmt.Errorf(ERR_DIR_NOT_FOUND, source_path)
	}

	var classes []ClassInfo

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".class") {
			if file.Name() == "GUI.class" {
				continue
			}
			continue
		}

		class_name := strings.TrimSuffix(file.Name(), ".class")
		file_path := filepath.Join(source_path, file.Name())

		class_info, err := self.parse_class(file_path, class_name)
		if err != nil {
			return nil, fmt.Errorf(ERR_PARSE_FAILED, file_path)
		}

		classes = append(classes, *class_info)
	}

	if len(classes) == 0 {
		panic(fmt.Sprintf(ERR_NO_CLASSES_FOUND, source_path))
	}

	// Report cache statistics if verbose
	if self.cacheManager != nil && self.cacheManager.verbose {
		total := self.cacheHits + self.cacheMisses
		if total > 0 {
			hitRate := float64(self.cacheHits) / float64(total) * 100.0
			fmt.Fprintf(os.Stderr, "  Cache statistics: %d hits, %d misses, %.0f%% hit rate\n",
				self.cacheHits, self.cacheMisses, hitRate)
		}
	}

	return classes, nil
}

func (self *JavapParser) parse_class(file_path string, class_name string) (*ClassInfo, error) {
	var output []byte
	var err error

	// Try to use cache if available
	if self.cacheManager != nil {
		cached, cacheErr := self.cacheManager.GetCachedOutput(file_path)
		if cacheErr == nil {
			// Cache hit
			output = cached
			self.cacheHits++
			if self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  Cache HIT: %s\n", class_name)
			}
		} else {
			// Cache miss - run javap
			self.cacheMisses++
			if self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  Cache MISS: %s - running javap\n", class_name)
			}
			cmd := exec.Command("javap", "-c", "-p", file_path)
			output, err = cmd.Output()
			if err != nil {
				return nil, fmt.Errorf(ERR_JAVAP_EXECUTION, err)
			}

			// Store in cache
			if cacheStoreErr := self.cacheManager.StoreCachedOutput(file_path, output); cacheStoreErr != nil {
				if self.cacheManager.verbose {
					fmt.Fprintf(os.Stderr, "  Warning: failed to store cache for %s: %v\n", class_name, cacheStoreErr)
				}
				// Don't fail the parsing, just warn
			}
		}
	} else {
		// No cache - run javap directly
		cmd := exec.Command("javap", "-c", "-p", file_path)
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf(ERR_JAVAP_EXECUTION, err)
		}
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	class_info := &ClassInfo{
		Name:   class_name,
		Source: SourceDeobfuscated,
	}

	for scanner.Scan() {
		line := scanner.Text()

		if matches := self.regex_class_decl.FindStringSubmatch(line); matches != nil {
			class_info.Name = strings.TrimSpace(matches[3])
		}

		if matches := self.regex_extends.FindStringSubmatch(line); matches != nil {
			class_info.Superclass = strings.TrimSpace(matches[1])
		}

		if matches := self.regex_implements.FindStringSubmatch(line); matches != nil {
			interfaces := strings.Split(matches[1], ",")
			for _, iface := range interfaces {
				class_info.Interfaces = append(class_info.Interfaces, strings.TrimSpace(iface))
			}
		}

		if matches := self.regex_field.FindStringSubmatch(line); matches != nil {
			class_info.Fields = append(class_info.Fields, FieldInfo{
				TypeName: matches[4],
				Name:     matches[5],
			})
		}

		if matches := self.regex_method.FindStringSubmatch(line); matches != nil {
			class_info.Methods = append(class_info.Methods, MethodInfo{
				ReturnType: matches[4],
				Name:       matches[5],
				Parameters: self.parse_parameters(matches[6]),
			})
		}
	}

	return class_info, nil
}

func (self *JavapParser) parse_parameters(param_str string) []string {
	if strings.TrimSpace(param_str) == "" {
		return []string{}
	}
	return strings.Split(param_str, ",")
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
		// Build reference maps based on field types
		for _, field := range class.Fields {
			if self.is_known_class_reference(field.TypeName) {
				self.deob_references[class.Name] = append(self.deob_references[class.Name], field.TypeName)
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
		deob_refs := self.deob_references[deob_class.Name]
		if len(deob_refs) == 0 {
			continue
		}

		best_match := self.find_best_cross_reference_match(&deob_class, deob_refs, obf_classes)
		if best_match != nil {
			result := self.create_match(&deob_class, best_match)
			result.ConfidenceScore = 75.0 // Moderate confidence for cross-reference match
			result.Details = "Matched by cross-reference analysis"
			results = append(results, result)
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

type StdoutWriter struct {
	writer *csv.Writer
}

func NewStdoutWriter() *StdoutWriter {
	writer := csv.NewWriter(os.Stdout)
	return &StdoutWriter{
		writer: writer,
	}
}

func (self *StdoutWriter) WriteHighConfidence(results []MatchResult) error {
	if len(results) == 0 {
		return nil
	}

	header := []string{
		"deobfuscated_name",
		"obfuscated_name",
		"confidence_score",
		"superclass_match",
		"interface_count",
		"field_count_deob",
		"field_count_obf",
		"method_count_deob",
		"method_count_obf",
		"notes",
	}

	if err := self.writer.Write(header); err != nil {
		panic(fmt.Sprintf(ERR_WRITE_FAILED, err))
	}

	for _, result := range results {
		if result.ConfidenceScore < HIGH_CONFIDENCE_THRESHOLD {
			continue
		}

		row := []string{
			result.DeobfuscatedClass,
			result.ObfuscatedClass,
			fmt.Sprintf("%.2f", result.ConfidenceScore),
			fmt.Sprintf("%t", result.ScoreBreakdown.SuperclassMatch > 0),
			fmt.Sprintf("%.0f", result.ScoreBreakdown.InterfaceMatch/WEIGHT_INTERFACE_MATCH),
			"N/A",
			"N/A",
			"N/A",
			"N/A",
			result.Details,
		}

		if err := self.writer.Write(row); err != nil {
			panic(fmt.Sprintf(ERR_WRITE_FAILED, err))
		}
	}

	self.writer.Flush()
	return nil
}

func (self *StdoutWriter) WriteUncertain(results []MatchResult) error {
	if len(results) == 0 {
		return nil
	}

	header := []string{
		"deobfuscated_name",
		"obfuscated_name",
		"confidence_score",
		"notes",
	}

	if err := self.writer.Write(header); err != nil {
		panic(fmt.Sprintf(ERR_WRITE_FAILED, err))
	}

	for _, result := range results {
		if result.ConfidenceScore >= HIGH_CONFIDENCE_THRESHOLD {
			continue
		}

		row := []string{
			result.DeobfuscatedClass,
			result.ObfuscatedClass,
			fmt.Sprintf("%.2f", result.ConfidenceScore),
			result.Details,
		}

		if err := self.writer.Write(row); err != nil {
			panic(fmt.Sprintf(ERR_WRITE_FAILED, err))
		}
	}

	self.writer.Flush()
	return nil
}

type JSONWriter struct{}

func NewJSONWriter() *JSONWriter {
	return &JSONWriter{}
}

func (self *JSONWriter) WriteHighConfidence(results []MatchResult) error {
	return nil
}

func (self *JSONWriter) WriteUncertain(results []MatchResult) error {
	return nil
}

func (self *JSONWriter) WriteDetailed(results []MatchResult) error {
	if len(results) == 0 {
		return nil
	}

	detailed := struct {
		Summary struct {
			TotalMatches          int `json:"total_matches"`
			HighConfidenceCount   int `json:"high_confidence_count"`
			MediumConfidenceCount int `json:"medium_confidence_count"`
			LowConfidenceCount    int `json:"low_confidence_count"`
		} `json:"summary"`
		Matches []MatchResult `json:"matches"`
	}{
		Matches: results,
	}

	for _, result := range results {
		detailed.Summary.TotalMatches++
		if result.ConfidenceScore >= HIGH_CONFIDENCE_THRESHOLD {
			detailed.Summary.HighConfidenceCount++
		} else if result.ConfidenceScore >= MIN_CONFIDENCE_THRESHOLD {
			detailed.Summary.MediumConfidenceCount++
		} else {
			detailed.Summary.LowConfidenceCount++
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(detailed); err != nil {
		panic(fmt.Sprintf(ERR_WRITE_FAILED, err))
	}

	return nil
}

func main() {
	deob_dir := flag.String(
		"deob",
		"srcAllDummysRemoved/bin",
		"Path to deobfuscated class files",
	)
	obf_dir := flag.String(
		"obf",
		"bytecode/client",
		"Path to obfuscated bytecode files",
	)
	mode := flag.String(
		"mode",
		"csv",
		"Output mode: csv or json (default: csv)",
	)
	threshold := flag.Float64(
		"threshold",
		MIN_CONFIDENCE_THRESHOLD,
		"Minimum confidence threshold (0-100)",
	)
	cacheDir := flag.String(
		"cache-dir",
		"srcAllDummysRemoved/.javap_cache",
		"Javap cache directory",
	)
	cleanCache := flag.Bool(
		"clean-cache",
		false,
		"Clean cache and exit",
	)
	noCache := flag.Bool(
		"no-cache",
		false,
		"Disable caching",
	)
	verbose := flag.Bool(
		"verbose",
		false,
		"Verbose output with timing information",
	)
	workers := flag.Int(
		"workers",
		-1,
		"Number of worker goroutines (-1 = auto, runtime.NumCPU()-1)",
	)
	parallelProgress := flag.Bool(
		"parallel-progress",
		false,
		"Show real-time parallel progress",
	)

	flag.Parse()

	// Handle cache cleanup flag
	if *cleanCache {
		_, err := NewCacheManager(*cacheDir, *verbose)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize cache manager: %v", err))
		}

		if *verbose {
			fmt.Fprintf(os.Stderr, "Cleaning cache in %s...\n", *cacheDir)
		}

		if err := os.RemoveAll(*cacheDir); err != nil {
			panic(fmt.Sprintf("failed to clean cache: %v", err))
		}

		fmt.Fprintf(os.Stderr, "Cache cleaned successfully\n")
		return // Exit after cleaning
	}

	if *mode != "csv" && *mode != "json" {
		panic("mode must be 'csv' or 'json'")
	}

	if *threshold < 0 || *threshold > 100 {
		panic("threshold must be between 0 and 100")
	}

	// Initialize cache manager
	var cacheManager *CacheManager
	if *noCache {
		if *verbose {
			fmt.Fprintf(os.Stderr, "Caching disabled (-no-cache=true)\n")
		}
		cacheManager = nil
	} else {
		var err error
		cacheManager, err = NewCacheManager(*cacheDir, *verbose)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize cache manager: %v", err))
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Cache directory: %s\n", *cacheDir)
		}
	}

	fmt.Fprintf(os.Stderr, "Parsing deobfuscated classes from: %s\n", *deob_dir)

	var parseStart time.Time
	if *verbose {
		parseStart = time.Now()
	}

	javap_parser := NewJavapParser(cacheManager)
	deob_classes, err := javap_parser.ParseAll(*deob_dir)
	if err != nil {
		panic(fmt.Sprintf("failed to parse deobfuscated classes: %v", err))
	}
	fmt.Fprintf(os.Stderr, "  Found %d deobfuscated classes\n", len(deob_classes))

	if *verbose {
		parseDuration := time.Since(parseStart)
		fmt.Fprintf(os.Stderr, "  Parsing time: %.3fs\n", parseDuration.Seconds())

		// Estimate time saved by cache
		if cacheManager != nil {
			totalClasses := javap_parser.cacheHits + javap_parser.cacheMisses
			if totalClasses > 0 {
				avgJavapTime := 0.15 // estimated average javap execution time in seconds
				estimatedSavings := float64(javap_parser.cacheHits) * avgJavapTime
				if javap_parser.cacheHits > 0 {
					fmt.Fprintf(os.Stderr, "  Estimated time saved by cache: %.3fs\n", estimatedSavings)
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Parsing obfuscated classes from: %s\n", *obf_dir)

	var parse_start time.Time
	if *parallelProgress {
		parse_start = time.Now()
	}

	bytecode_parser := NewBytecodeParser()

	// Create progress tracker if enabled
	var progress *ProgressTracker
	if *parallelProgress {
		// Count files first for progress tracking
		if files, err := os.ReadDir(*obf_dir); err == nil {
			count := 0
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".bytecode.txt") {
					count++
				}
			}
			progress = NewProgressTracker(count, *parallelProgress)
			fmt.Fprintf(os.Stderr, "  Using %d worker goroutines\n", *workers)
		}
	}

	obf_classes, err := bytecode_parser.ParseAllParallel(*obf_dir, *workers, progress)
	if err != nil {
		panic(fmt.Sprintf("failed to parse obfuscated classes: %v", err))
	}

	if progress != nil {
		progress.Finalize()
	}

	fmt.Fprintf(os.Stderr, "  Found %d obfuscated classes\n", len(obf_classes))

	if *parallelProgress {
		parse_duration := time.Since(parse_start)
		fmt.Fprintf(os.Stderr, "  Parsing time: %.3fs\n", parse_duration.Seconds())
	}

	resolver_instance := NewResolver(ANCHOR_MAPPINGS)

	fmt.Fprintln(os.Stderr, "Resolving class matches...")
	matches := resolver_instance.ResolveAll(deob_classes, obf_classes)
	fmt.Fprintf(os.Stderr, "  Found %d total matches\n", len(matches))

	valid_matches := filter_by_threshold(matches, *threshold)
	fmt.Fprintf(os.Stderr, "  %d matches above threshold %.2f\n", len(valid_matches), *threshold)

	// Validate and refine matches
	valid_matches = resolver_instance.validate_and_refine_matches(valid_matches)

	high_conf := count_by_confidence(valid_matches, HIGH_CONFIDENCE_THRESHOLD)
	medium_conf := len(valid_matches) - high_conf

	if *mode == "json" {
		fmt.Fprintln(os.Stderr, "\nJSON Output:")
		json_writer := NewJSONWriter()
		if err := json_writer.WriteDetailed(valid_matches); err != nil {
			panic(fmt.Sprintf("failed to write JSON: %v", err))
		}
	} else {
		fmt.Fprintln(os.Stderr, "\nHigh Confidence Matches (CSV):")
		stdout_writer := NewStdoutWriter()
		if err := stdout_writer.WriteHighConfidence(valid_matches); err != nil {
			panic(err)
		}

		fmt.Fprintln(os.Stderr, "\nUncertain Matches (CSV):")
		if err := stdout_writer.WriteUncertain(valid_matches); err != nil {
			panic(err)
		}
	}

	fmt.Fprintln(os.Stderr, "\nSummary:")
	fmt.Fprintf(os.Stderr, "  High confidence (≥%.2f): %d\n", HIGH_CONFIDENCE_THRESHOLD, high_conf)
	fmt.Fprintf(os.Stderr, "  Medium confidence (%.2f-%.2f): %d\n", MIN_CONFIDENCE_THRESHOLD, HIGH_CONFIDENCE_THRESHOLD, medium_conf)
	fmt.Fprintf(os.Stderr, "  Low confidence (<%.2f): %d\n", HIGH_CONFIDENCE_THRESHOLD, len(valid_matches)-high_conf-medium_conf)

	// Save cache checksums for next run
	if cacheManager != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "Saving cache checksums...\n")
		}
		if err := cacheManager.SaveChecksums(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to save cache checksums: %v\n", err)
		} else if *verbose {
			fmt.Fprintf(os.Stderr, "Cache checksums saved successfully\n")
		}
	} else if *verbose {
		fmt.Fprintf(os.Stderr, "Cache manager is nil, skipping checksum save\n")
	}

	fmt.Fprintln(os.Stderr, "\nDone!")
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

func filter_by_threshold(matches []MatchResult, threshold float64) []MatchResult {
	var filtered []MatchResult
	for _, match := range matches {
		if match.ConfidenceScore >= threshold {
			filtered = append(filtered, match)
		}
	}
	return filtered
}

func count_by_confidence(matches []MatchResult, threshold float64) int {
	count := 0
	for _, match := range matches {
		if match.ConfidenceScore >= threshold {
			count++
		}
	}
	return count
}
