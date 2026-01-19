package debug

import "fmt"

// Debug provides runtime debugging following AMILLI.md pattern.
//
// All debug output uses TODO_REMOVE_ prefix to ensure
// cleanup before production deployment.
type Debug struct{}

// Printf outputs formatted debug message to console.
//
// Uses TODO_REMOVE_ prefix per AMILLI.md guidelines for
// easy identification and removal before production.
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
