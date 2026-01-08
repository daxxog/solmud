// Quick test to verify data structure classification
package main

import (
	"fmt"
	"os"
)

func testDataStructureClassification() {
	// Read the Model.java source file
	sourcePath := "../../srcAllDummysRemoved/src/Model.java"
	sourceContent, err := os.ReadFile(sourcePath)
	if err != nil {
		fmt.Printf("Error reading source file: %v\n", err)
		return
	}

	// Create parser and test array detection
	parser := NewInternalBehaviorParser()
	actualMetrics := parser.detectArrayPatterns(string(sourceContent))

	fmt.Printf("Array Pattern Detection Results for Model.java:\n")
	fmt.Printf("  SingleDimAccess: %d\n", actualMetrics.SingleDimAccess)
	fmt.Printf("  MultiDimAccess: %d\n", actualMetrics.MultiDimAccess)
	fmt.Printf("  BulkOperations: %d\n", actualMetrics.BulkOperations)
	fmt.Printf("  LoopArrayAccess: %d\n", actualMetrics.LoopArrayAccess)
	fmt.Printf("  NestedArrayDepth: %d\n", actualMetrics.NestedArrayDepth)
	fmt.Printf("  ArrayAlgorithms: %v\n", actualMetrics.ArrayAlgorithms)

	// Create classifier and test with actual metrics
	classifier := NewDataStructureClassifier()
	signature := classifier.classifyDataStructure(actualMetrics, string(sourceContent))

	fmt.Printf("Data Structure Classification Results for Model.java:\n")
	fmt.Printf("  StructureType: %d\n", signature.StructureType)
	fmt.Printf("  DimCount: %d\n", signature.DimCount)
	fmt.Printf("  ElementSize: %d\n", signature.ElementSize)
	fmt.Printf("  AccessPattern: %s\n", signature.AccessPattern)
	fmt.Printf("  LoopNesting: %d\n", signature.LoopNesting)
	fmt.Printf("  DomainSpecific: %+v\n", signature.DomainSpecific)

	// Test Texture classification
	texturePath := "../../srcAllDummysRemoved/src/Texture.java"
	textureContent, err := os.ReadFile(texturePath)
	if err != nil {
		fmt.Printf("Error reading texture file: %v\n", err)
		return
	}

	textureMetrics := parser.detectArrayPatterns(string(textureContent))
	textureSignature := classifier.classifyDataStructure(textureMetrics, string(textureContent))
	fmt.Printf("\nData Structure Classification Results for Texture.java:\n")
	fmt.Printf("  StructureType: %d\n", textureSignature.StructureType)
	fmt.Printf("  SingleDimAccess: %d\n", textureMetrics.SingleDimAccess)
	fmt.Printf("  DomainSpecific: %+v\n", textureSignature.DomainSpecific)
}
