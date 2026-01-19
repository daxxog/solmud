package assert

import "fmt"

// IAssert defines interface for runtime assertions.
//
// Provides methods for validating invariants and catching
// programming errors that should never occur in correct code.
type IAssert interface {
	// NotNil validates that value is not nil.
	NotNil(value interface{}, message string)

	// True validates that condition is true.
	True(condition bool, message string)

	// InRange validates that value is within [min, max] range.
	InRange(value, min, max int, name string)
}

// Assert provides production assertion implementation.
//
// Panics on assertion failures with descriptive messages.
// Used for invariants that MUST always be true.
type Assert struct{}

// NewAssert creates a production assertion implementation.
func NewAssert() IAssert {
	return &Assert{}
}

func (a *Assert) NotNil(value interface{}, message string) {
	if value == nil {
		panic(fmt.Sprintf("ASSERTION FAILED: %s - value is nil", message))
	}
}

func (a *Assert) True(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("ASSERTION FAILED: %s", message))
	}
}

func (a *Assert) InRange(value, min, max int, name string) {
	if value < min || value > max {
		panic(fmt.Sprintf("ASSERTION FAILED: %s=%d outside range [%d, %d]",
			name, value, min, max))
	}
}
