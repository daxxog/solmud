package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type JavapParser struct {
	regex_class_decl *regexp.Regexp
	regex_extends    *regexp.Regexp
	regex_implements *regexp.Regexp
	regex_field      *regexp.Regexp
	regex_method     *regexp.Regexp
	cacheManager     *CacheManager
	cacheHits        int
	cacheMisses      int

	// Cross-reference parsing
	crossRefParser  *CrossReferenceParser
	projectClasses  []string
	anchorMappings  map[string]string
	enableCrossRefs bool

	// Internal behavior parsing (Phase 3.2.2.1)
	internalBehaviorParser *InternalBehaviorParser
	enableInternalBehavior bool
}

func NewJavapParser(cacheManager *CacheManager, projectClasses []string, anchors map[string]string, enableCrossRefs bool) *JavapParser {
	parser := &JavapParser{
		regex_class_decl:       regexp.MustCompile(`^\s*(?:(?:public|private|protected|final|abstract|static)\s+)*(class|interface)\s+(\w+)`),
		regex_extends:          regexp.MustCompile(`extends\s+(\w+)`),
		regex_implements:       regexp.MustCompile(`implements\s+([\w,\.]+)`),
		regex_field:            regexp.MustCompile(`^\s*(private|public|protected)?\s*(static\s+)?(final\s+)?([\[\]\w]+)\s+(\w+)`),
		regex_method:           regexp.MustCompile(`^\s*(public|private|protected)?\s*(static\s+)?(final\s+)?([\[\]\w]+)\s+(\w+)\(([^)]*)\)`),
		cacheManager:           cacheManager,
		projectClasses:         projectClasses,
		anchorMappings:         anchors,
		enableCrossRefs:        enableCrossRefs,
		enableInternalBehavior: true, // Enable by default for Phase 3.2.2.1
	}

	// Initialize cross-reference parser if enabled
	if enableCrossRefs {
		parser.crossRefParser = NewCrossReferenceParser(projectClasses, anchors)
	}

	// Initialize internal behavior parser
	parser.internalBehaviorParser = NewInternalBehaviorParser()

	return parser
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
			class_info.Name = strings.TrimSpace(matches[2])
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

	// Parse cross-references from source file if enabled
	if self.enableCrossRefs && self.crossRefParser != nil {
		if self.cacheManager != nil && self.cacheManager.verbose {
			fmt.Fprintf(os.Stderr, "  [CROSSREF] Processing %s...\n", class_name)
		}
		// Convert class file path to source file path
		sourcePath := self.convertClassToSourcePath(file_path)
		if sourcePath != "" {
			if self.cacheManager != nil && self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  [CROSSREF] Found source: %s\n", sourcePath)
			}
			crossRefs, err := self.crossRefParser.ParseFile(sourcePath)
			if err != nil {
				// Log warning but don't fail the entire parse
				if self.cacheManager != nil && self.cacheManager.verbose {
					fmt.Fprintf(os.Stderr, "  Warning: Failed to parse cross-references for %s: %v\n", class_name, err)
				}
				crossRefs = &ClassCrossReferences{} // Empty fallback
			}
			if self.cacheManager != nil && self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  [CROSSREF] Parsed %d references for %s\n", crossRefs.TotalReferences, class_name)
			}
			class_info.CrossReferences = crossRefs
		} else {
			if self.cacheManager != nil && self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  [CROSSREF] No source file found for %s\n", class_name)
			}
		}
	}

	// Parse internal behavioral patterns from source file (Phase 3.2.2.1)
	if self.enableInternalBehavior && self.internalBehaviorParser != nil {
		if self.cacheManager != nil && self.cacheManager.verbose {
			fmt.Fprintf(os.Stderr, "  [INTERNAL] Processing %s...\n", class_name)
		}
		// Convert class file path to source file path
		sourcePath := self.convertClassToSourcePath(file_path)
		if sourcePath != "" {
			if self.cacheManager != nil && self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  [INTERNAL] Found source: %s\n", sourcePath)
			}
			sourceContent, err := os.ReadFile(sourcePath)
			if err != nil {
				if self.cacheManager != nil && self.cacheManager.verbose {
					fmt.Fprintf(os.Stderr, "  Warning: Failed to read source file for internal behavior analysis %s: %v\n", class_name, err)
				}
			} else {
				internalBehavior := self.internalBehaviorParser.ParseInternalBehavior(string(sourceContent), class_info.Methods)
				if self.cacheManager != nil && self.cacheManager.verbose {
					fmt.Fprintf(os.Stderr, "  [INTERNAL] Parsed internal behavior for %s\n", class_name)
				}
				class_info.InternalBehavior = internalBehavior
			}
		} else {
			if self.cacheManager != nil && self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  [INTERNAL] No source file found for %s\n", class_name)
			}
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

// convertClassToSourcePath converts a .class file path to the corresponding .java source file path
func (self *JavapParser) convertClassToSourcePath(classPath string) string {
	// Assume class files are in srcAllDummysRemoved/bin/
	// Source files are in srcAllDummysRemoved/src/

	// Replace /bin/ with /src/ and change extension from .class to .java
	sourcePath := strings.Replace(classPath, "/bin/", "/src/", 1)
	sourcePath = strings.TrimSuffix(sourcePath, ".class") + ".java"

	// Check if the source file exists
	if _, err := os.Stat(sourcePath); err == nil {
		if self.cacheManager != nil && self.cacheManager.verbose {
			fmt.Fprintf(os.Stderr, "  Found source file: %s\n", sourcePath)
		}
		return sourcePath
	}

	// Try alternative: look in src/ directory from project root
	// Extract just the class name and look in srcAllDummysRemoved/src/
	if strings.Contains(classPath, "srcAllDummysRemoved/bin/") {
		className := strings.TrimSuffix(filepath.Base(classPath), ".class")
		altPath := filepath.Join("srcAllDummysRemoved", "src", className+".java")
		if _, err := os.Stat(altPath); err == nil {
			if self.cacheManager != nil && self.cacheManager.verbose {
				fmt.Fprintf(os.Stderr, "  Found alt source file: %s\n", altPath)
			}
			return altPath
		}

		// Source file not found
		if self.cacheManager != nil && self.cacheManager.verbose {
			fmt.Fprintf(os.Stderr, "  Source file not found for: %s (tried: %s, %s)\n", classPath, sourcePath, altPath)
		}
	}
	return ""
}
