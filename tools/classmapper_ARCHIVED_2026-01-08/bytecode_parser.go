package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type BytecodeParser struct {
	regex_class_decl *regexp.Regexp
	regex_extends    *regexp.Regexp
	regex_implements *regexp.Regexp
	regex_field      *regexp.Regexp
	regex_method     *regexp.Regexp

	// Cross-reference analysis
	bytecodeCrossrefParser *BytecodeCrossReferenceParser
}

func NewBytecodeParser(projectClasses []string, reverseAnchors map[string]string) *BytecodeParser {
	return &BytecodeParser{
		regex_class_decl: regexp.MustCompile(`^\s*(?:(?:public|private|protected|final|abstract|static)\s+)*(class|interface)\s+(\w+)`),
		regex_extends:    regexp.MustCompile(`extends (\w+)`),
		regex_implements: regexp.MustCompile(`implements ([\w,\s\.]+)`),
		regex_field:      regexp.MustCompile(`^private|public|protected?\s+(?:static\s+)?(?:final\s+)?([\[\]\w]+)\s+(\w+)`),
		regex_method:     regexp.MustCompile(`^public|private|protected?\s+(?:static\s+)?([\[\]\w]+)\s+(\w+)\(([^)]*)\)`),

		bytecodeCrossrefParser: NewBytecodeCrossReferenceParser(projectClasses, reverseAnchors),
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

	// Extract cross-references from bytecode if parser is available
	if self.bytecodeCrossrefParser != nil {
		crossRefs, err := self.bytecodeCrossrefParser.ParseBytecodeFile(file_path)
		if err != nil {
			// Log warning but don't fail parsing
			fmt.Fprintf(os.Stderr, "Warning: failed to extract cross-references from %s: %v\n", file_path, err)
		} else {
			class_info.CrossReferences = crossRefs
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
