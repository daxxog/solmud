package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	HIGH_CONFIDENCE_THRESHOLD = 85.0
	MIN_CONFIDENCE_THRESHOLD  = 65.0
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

	// Cross-reference data
	CrossReferences *ClassCrossReferences `json:"cross_refs,omitempty"`
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
	DeobfuscatedClass string         `json:"deobfuscated_class"`
	ObfuscatedClass   string         `json:"obfuscated_class"`
	ConfidenceScore   float64        `json:"confidence_score"`
	ScoreBreakdown    ScoreBreakdown `json:"score_breakdown"`
	Details           string         `json:"details"`
}

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
	// Graphics classes
	"DrawingArea":     "AFCKELYG", // Base graphics class with pixel manipulation
	"Background":      "DSMJIEPN", // Loads texture data
	"Sprite":          "CXGZMTJK", // Image rendering class
	"RSImageProducer": "IVIFZQBK", // Image producer for rendering
	"Decompressor":    "IGSLDTHC", // Archive decompression with RandomAccessFile
	"RSSocket":        "NQABEVLK", // Network socket with Runnable interface
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
	enableCrossRefs := flag.Bool(
		"crossref",
		true,
		"Enable cross-reference analysis (default: true)",
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

	// First pass: get project class names for cross-reference analysis
	temp_parser := NewJavapParser(cacheManager, []string{}, ANCHOR_MAPPINGS, false)
	temp_classes, err := temp_parser.ParseAll(*deob_dir)
	if err != nil {
		panic(fmt.Sprintf("failed to parse deobfuscated classes: %v", err))
	}

	// Extract class names for cross-reference analysis
	projectClasses := make([]string, len(temp_classes))
	for i, class := range temp_classes {
		projectClasses[i] = class.Name
	}

	// Create parser with cross-reference support if enabled
	javap_parser := NewJavapParser(cacheManager, projectClasses, ANCHOR_MAPPINGS, *enableCrossRefs)
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
			totalClasses := len(deob_classes)
			if totalClasses > 0 {
				avgJavapTime := 0.15 // estimated average javap execution time in seconds
				estimatedSavings := float64(totalClasses) * avgJavapTime
				if totalClasses > 0 {
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
			actual_workers := *workers
			if actual_workers <= 0 {
				actual_workers = runtime.NumCPU() - 1
				if actual_workers < 1 {
					actual_workers = 1
				}
			}
			fmt.Fprintf(os.Stderr, "  Using %d worker goroutines\n", actual_workers)
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
