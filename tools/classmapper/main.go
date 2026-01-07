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
	WEIGHT_METHOD_COUNT_MATCH = 5.0
	WEIGHT_METHOD_SIG_SIM     = 20.0
	WEIGHT_CONSTRUCTOR_MATCH  = 5.0
	WEIGHT_ACCESS_MOD_MATCH   = 5.0
)

var ANCHOR_MAPPINGS = map[string]string{
	"sign/signlink": "sign/signlink",
	"Node":          "PKVMXVTO",
	"NodeSub":       "PPOHBEGB",
	"RSApplet":      "KHACHIFW",
	"client":        "client",
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
	InterfaceMatch   float64
	SuperclassMatch  float64
	FieldCountMatch  float64
	FieldSimilarity  float64
	MethodCountMatch float64
	MethodSimilarity float64
	ConstructorMatch float64
	AccessMatch      float64
	SizePenalty      float64
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
		regex_class_decl: regexp.MustCompile(`^public (class|interface) (\w+)`),
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
			class_info.Name = matches[2]
		}

		if matches := self.regex_extends.FindStringSubmatch(line); matches != nil {
			class_info.Superclass = matches[1]
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
}

func NewJavapParser() *JavapParser {
	return &JavapParser{
		regex_class_decl: regexp.MustCompile(`^(public|final|abstract)?\s*(class|interface)\s+(\w+)`),
		regex_extends:    regexp.MustCompile(`extends\s+(\w+)`),
		regex_implements: regexp.MustCompile(`implements\s+([\w,\.]+)`),
		regex_field:      regexp.MustCompile(`^\s*(private|public|protected)?\s*(static\s+)?(final\s+)?([\[\]\w]+)\s+(\w+)`),
		regex_method:     regexp.MustCompile(`^\s*(public|private|protected)?\s*(static\s+)?(final\s+)?([\[\]\w]+)\s+(\w+)\(([^)]*)\)`),
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

	return classes, nil
}

func (self *JavapParser) parse_class(file_path string, class_name string) (*ClassInfo, error) {
	cmd := exec.Command("javap", "-c", "-p", file_path)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(ERR_JAVAP_EXECUTION, err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	class_info := &ClassInfo{
		Name:   class_name,
		Source: SourceDeobfuscated,
	}

	for scanner.Scan() {
		line := scanner.Text()

		if matches := self.regex_class_decl.FindStringSubmatch(line); matches != nil {
			class_info.Name = matches[3]
		}

		if matches := self.regex_extends.FindStringSubmatch(line); matches != nil {
			class_info.Superclass = matches[1]
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

	if len(deob_class.Methods) == len(obf_class.Methods) {
		breakdown.MethodCountMatch = WEIGHT_METHOD_COUNT_MATCH
	}

	method_sim := self.calculate_method_similarity(deob_class, obf_class)
	breakdown.MethodSimilarity = method_sim * WEIGHT_METHOD_SIG_SIM

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
				matching_methods++
				break
			}
		}
	}

	return float64(matching_methods) / total_methods
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
	anchors         map[string]string
	reverse_anchors map[string]string
	scorer          *Scorer
}

func NewResolver(anchors map[string]string) *Resolver {
	reverse_anchors := make(map[string]string)
	for deob, obf := range anchors {
		reverse_anchors[obf] = deob
	}

	return &Resolver{
		anchors:         anchors,
		reverse_anchors: reverse_anchors,
		scorer:          NewScorer(anchors),
	}
}

func (self *Resolver) ResolveAll(deob_classes []ClassInfo, obf_classes []ClassInfo) []MatchResult {
	var results []MatchResult

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

	for _, deob_class := range pending_deob {
		best_match, best_score := self.find_best_match(&deob_class, pending_obf)
		if best_match != nil && best_score >= HIGH_CONFIDENCE_THRESHOLD {
			result := self.create_match(&deob_class, best_match)
			result.ConfidenceScore = best_score
			result.Details = fmt.Sprintf("Matched by inheritance - score %.2f", best_score)
			results = append(results, result)
		}
	}

	fmt.Fprintln(os.Stderr, "Pass 3: Matching by signatures...")
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
		breakdown.FieldCountMatch + breakdown.FieldSimilarity +
		breakdown.MethodCountMatch + breakdown.MethodSimilarity +
		breakdown.ConstructorMatch + breakdown.AccessMatch +
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

	flag.Parse()

	if *mode != "csv" && *mode != "json" {
		panic("mode must be 'csv' or 'json'")
	}

	if *threshold < 0 || *threshold > 100 {
		panic("threshold must be between 0 and 100")
	}

	fmt.Fprintf(os.Stderr, "Parsing deobfuscated classes from: %s\n", *deob_dir)
	javap_parser := NewJavapParser()
	deob_classes, err := javap_parser.ParseAll(*deob_dir)
	if err != nil {
		panic(fmt.Sprintf("failed to parse deobfuscated classes: %v", err))
	}
	fmt.Fprintf(os.Stderr, "  Found %d deobfuscated classes\n", len(deob_classes))

	fmt.Fprintf(os.Stderr, "Parsing obfuscated classes from: %s\n", *obf_dir)
	bytecode_parser := NewBytecodeParser()
	obf_classes, err := bytecode_parser.ParseAll(*obf_dir)
	if err != nil {
		panic(fmt.Sprintf("failed to parse obfuscated classes: %v", err))
	}
	fmt.Fprintf(os.Stderr, "  Found %d obfuscated classes\n", len(obf_classes))

	resolver_instance := NewResolver(ANCHOR_MAPPINGS)

	fmt.Fprintln(os.Stderr, "Resolving class matches...")
	matches := resolver_instance.ResolveAll(deob_classes, obf_classes)
	fmt.Fprintf(os.Stderr, "  Found %d total matches\n", len(matches))

	valid_matches := filter_by_threshold(matches, *threshold)
	fmt.Fprintf(os.Stderr, "  %d matches above threshold %.2f\n", len(valid_matches), *threshold)

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

	fmt.Fprintln(os.Stderr, "\nDone!")
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
