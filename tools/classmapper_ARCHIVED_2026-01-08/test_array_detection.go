// Quick test to verify array pattern detection
package main

import (
	"fmt"
	"os"
)

// testArrayDetection verifies that our enhanced array pattern detection is working
func testArrayDetection() {
	// Read the Model.java source file
	sourcePath := "../../srcAllDummysRemoved/src/Model.java"
	sourceContent, err := os.ReadFile(sourcePath)
	if err != nil {
		fmt.Printf("Error reading source file: %v\n", err)
		return
	}

	// Create parser and test array detection
	parser := NewInternalBehaviorParser()
	patterns := parser.detectArrayPatterns(string(sourceContent))

	fmt.Printf("Array Pattern Detection Results for Model.java:\n")
	fmt.Printf("  SingleDimAccess: %d\n", patterns.SingleDimAccess)
	fmt.Printf("  MultiDimAccess: %d\n", patterns.MultiDimAccess)
	fmt.Printf("  BulkOperations: %d\n", patterns.BulkOperations)
	fmt.Printf("  LoopArrayAccess: %d\n", patterns.LoopArrayAccess)
	fmt.Printf("  NestedArrayDepth: %d\n", patterns.NestedArrayDepth)
	fmt.Printf("  ArrayAlgorithms: %v\n", patterns.ArrayAlgorithms)

	// Test some specific array context analysis
	context := parser.analyzeAccessContext(string(sourceContent), "anIntArray1627")
	fmt.Printf("\nArray Context Analysis for 'anIntArray1627':\n")
	fmt.Printf("  InLoop: %v\n", context.InLoop)
	fmt.Printf("  InMethodCall: %v\n", context.InMethodCall)
	fmt.Printf("  InCalculation: %v\n", context.InCalculation)
	fmt.Printf("  InAssignment: %v\n", context.InAssignment)
	fmt.Printf("  AccessFrequency: %.4f\n", context.AccessFrequency)
}
