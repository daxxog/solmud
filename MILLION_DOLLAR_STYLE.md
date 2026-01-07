# MILLION_DOLLAR_STYLE

## Core Philosophy

Our codebase embodies three fundamental principles, in order of priority:

1. Safety
2. Performance
3. Developer Experience

These principles are not in conflict. Through careful design and iteration, we can achieve all three simultaneously. This requires discipline and hard work, but the rewards are substantial: code that is reliable, efficient, and a joy to work with.

## Foundational Principles

### Zero Technical Debt

We solve problems when we find them. Technical debt is far more expensive to fix in production than during development. When we identify potential issues - whether they're performance bottlenecks, safety concerns, or maintainability challenges - we address them immediately.

A problem solved in production is many times more expensive than a problem solved in implementation, or better yet, a problem solved in design.

### Simplicity Through Iteration

Simplicity is not the first attempt but the hardest revision. We invest time upfront in design and iterate until we find elegant solutions that satisfy all our goals. An hour spent in design can save weeks in production.

### Think First, Code Later

Build a precise mental model before writing code. Document your understanding through comments and assertions. Use testing as verification, not as a substitute for understanding.

## Code Organization

### Package Structure

Every package must serve a single, well-defined purpose. Package documentation should explain:

1. The package's purpose and scope
2. Key components and their relationships
3. Usage patterns and examples
4. Thread safety guarantees
5. Performance characteristics

Example:
```go
// Package block provides core functionality for handling blockchain blocks.
//
// Key Components:
//   - Block: The primary implementation of IBlock
//   - BlockDump: Implementation of IBlockDump for high-level operations
//   - BlockFactory: Handles block creation and instantiation
//
// Thread Safety:
//   All public methods are safe for concurrent use.
//   Internal state is protected by mutexes where necessary.
//
// Performance Characteristics:
//   - Memory: O(1) for standard operations
//   - CPU: O(log n) for block validation
//   - Disk: Batch writes for efficiency
package block
```

### File Organization

1. Place all constants in a dedicated `const.go` file
2. Group related types and implementations together
3. Keep test files alongside the code they test
4. Place interfaces in `interfaces.go` when they define package boundaries

## Modern Go Practices

### Standard Library Usage

The Go standard library evolves to provide better solutions over time. Always use current, non-deprecated packages and functions:

1. Use modern IO operations:
   ```go
   // Prefer:
   data, err := io.ReadAll(reader)
   
   // Over deprecated:
   data, err := ioutil.ReadAll(reader)
   ```

2. Use `io` and `os` packages directly:
   ```go
   // Prefer:
   temp_file, err := os.CreateTemp("", "prefix")
   temp_dir, err := os.MkdirTemp("", "prefix")
   
   // Over deprecated:
   temp_file, err := ioutil.TempFile("", "prefix")
   temp_dir, err := ioutil.TempDir("", "prefix")
   ```

### Iterator Usage

Go's `iter.Seq` type provides a powerful iteration pattern. While it supports both callback-style and range-style iteration, we prefer the more readable range-style syntax:
```go
type (
	iter.Seq[V any]     func(yield func(V) bool)
	iter.Seq2[K, V any] func(yield func(K, V) bool)
)
```

1. Prefer range-style iteration:
   ```go
   // Good - clear and idiomatic:
   for entry := range collection.Entries() {
       process(entry)
   }

   // Avoid - more verbose callback style:
   collection.Entries()(func(entry Entry) bool {
       process(entry)
       return true
   })
   ```

2. Use nested range loops for hierarchical iteration:
   ```go
   // Good - clear hierarchy and scope:
   for entry := range acl.Entries() {
       fmt.Printf("Entry: %s\n", entry.Name())
       
       for item := range entry.Items() {
           fmt.Printf("  Item: %s\n", item.Name())
       }
   }

   // Avoid - nested callbacks are harder to follow:
   acl.Entries()(func(entry Entry) bool {
       fmt.Printf("Entry: %s\n", entry.Name())
       
       entry.Items()(func(item Item) bool {
           fmt.Printf("  Item: %s\n", item.Name())
           return true
       })
       return true
   })
   ```

3. The range syntax maintains proper scoping while being more readable:
   ```go
   // Good - clear variable scoping:
   for parent := range hierarchy.Parents() {
       parent_name := parent.Name()
       
       for child := range parent.Children() {
           // parent_name is clearly scoped
           process(parent_name, child.Name())
       }
   }
   ```

4. Use range when iteration order matters:
   ```go
   // Good - order is clear:
   for step := range workflow.OrderedSteps() {
       execute(step)
   }

   // Avoid - order dependency is hidden in callback:
   workflow.OrderedSteps()(func(step Step) bool {
       execute(step)
       return true
   })
   ```

This approach aligns with Go's design philosophy of clarity over cleverness and makes code easier to read, maintain, and reason about. The range-style iteration clearly shows the flow of data through your program while maintaining all the benefits of the `iter.Seq` type system.

Remember: Code is read far more often than it is written. The range-style iteration pattern produces more readable and maintainable code while preserving type safety and iterator semantics.

# Type-Safe Enum Pattern

## Core Philosophy

Go lacks built-in enum support, but we can create a superior pattern using generics and phantom types. Our enum pattern provides compile-time type safety, runtime validation, and excellent developer experience—all aligning with our core principles of safety, performance, and developer ergonomics.

## Package Organization

We organize enum-related code according to a specific directory structure that promotes clarity and separation of concerns:

```
pkg/
├── enum/
│   └── enum.go        # Core implementation of the generic enum pattern
└── enums/             # Directory containing specific enum types
    ├── currency_type/
    │   └── currency_type.go
    └── payment_method/
        └── payment_method.go
```

Each specific enum type lives in its own package under the `enums/` directory. This organization provides several benefits:

1. Clean import paths that clearly indicate the enum type
2. Proper encapsulation of enum-specific logic
3. Ability to add enum-specific methods or constants
4. Self-documenting code organization

## Core Implementation

The `enum` package provides the generic foundation for all enum types:

```go
package enum

// Value is a generic struct that holds the string value and implements type safety.
type Value[P any] struct {
    value string
}

// IEnumSet defines the interface for enum set operations with type safety.
type IEnumSet[P any] interface {
    Valid(value Value[P]) bool
    FromString(str string) (Value[P], error)
    All() iter.Seq[Value[P]]
}

// Set creates a new IEnumSet from a variable number of values.
func Set[P any](values ...Value[P]) IEnumSet[P] {
    // Implementation...
}

// Enum returns a function that constructs enum values of a specific type.
func Enum[P any]() func(string) Value[P] {
    // Implementation...
}
```

This core implementation is shared by all specific enum types.

## Specific Enum Type Implementation

Each enum type is defined in its own package with a consistent structure:

```go
// File: pkg/enums/payment_method/payment_method.go
package payment_method

import (
    "stack-track/pkg/enum"
)

// Define the phantom type
type EPaymentMethod struct{}

// Define the enum values
var (
    // Define constructor
    xPaymentMethod = enum.Enum[EPaymentMethod]()
    
    // Define enum values with clean names (no prefix needed)
    CREDIT = xPaymentMethod("credit")
    DEBIT  = xPaymentMethod("debit")
    PAYPAL = xPaymentMethod("paypal")
    
    // Create the enum set
    PaymentMethod = enum.Set(
        CREDIT,
        DEBIT,
        PAYPAL,
    )
)

// Export the enum set with a short name for convenient access
var (
    E = PaymentMethod
)
```

## Naming Conventions

The package-based organization allows us to use cleaner naming conventions:

1. **Phantom types**: Use `E` prefix with the enum name in PascalCase, e.g., `EPaymentMethod`

2. **Constructor variables**: Use `x` prefix with the enum name in camelCase, e.g., `xPaymentMethod`

3. **Enum values**: Use UPPERCASE without type prefix, e.g., `CREDIT` instead of `PAYMENT_METHOD_CREDIT`

4. **Set variables**: Use PascalCase enum name, e.g., `PaymentMethod`

5. **Exported set reference**: Use short name `E` for convenient access

## Importing and Using Enums

Enums are imported and used with clean, self-documenting syntax:

```go
import (
    "fmt"
    "stack-track/pkg/enum"
    PaymentMethod "stack-track/pkg/enums/payment_method"
)

func main() {
    // Access enum values directly from the package
    if PaymentMethod.E.Valid(PaymentMethod.CREDIT) {
        fmt.Println("Valid payment method:", PaymentMethod.CREDIT)
    }
    
    // Iterate over all values
    fmt.Println("All payment methods:")
    for v := range PaymentMethod.E.All() {
        fmt.Println("-", v)
    }
    
    // Convert from string with validation
    method, err := PaymentMethod.E.FromString("paypal")
    if err == nil {
        process_payment(method)
    }
}

// Type-safe function that only accepts payment methods
func process_payment(method enum.Value[PaymentMethod.EPaymentMethod]) {
    // Implementation...
}
```

This approach provides a clean, intuitive API while maintaining full type safety.

## Enum Package Implementation Guidelines

When implementing a new enum package, follow these guidelines:

1. **Create a dedicated package** for each enum type under the `enums/` directory

2. **Name the package file** the same as the directory name for clarity

3. **Define a phantom type** specific to this enum

4. **Create enum values** with clean, context-free names (the package provides the context)

5. **Export the enum set** as both the full name and as `E` for convenience

6. **Add documentation** that explains the purpose and usage of the enum type

## Benefits of Package-Based Organization

This organization pattern provides several advantages:

1. **Namespace Clarity**: The package name provides natural namespacing for enum values

2. **Import Control**: Consumers can import only the enum types they need

3. **Extension Point**: Each enum package can add specific behavior or constants related to that enum type

4. **Documentation**: Package documentation can provide enum-specific usage guidelines

5. **Evolution**: Enum types can evolve independently without affecting other types

## Advanced: Adding Behavior to Enum Packages

Enum packages can include behavior specific to that enum type:

```go
// In pkg/enums/currency_type/currency_type.go
package currency_type

import (
    "fmt"
    "stack-track/pkg/enum"
)

// Define type and values...

// Add currency-specific formatting function
func Format(amount float64, currency enum.Value[ECurrencyType]) string {
    switch currency {
    case FIAT:
        return fmt.Sprintf("$%.2f", amount)
    case CRYPTO_GEN1, CRYPTO_GEN2:
        return fmt.Sprintf("%.8f BTC", amount)
    default:
        return fmt.Sprintf("%.2f", amount)
    }
}
```

This allows for a rich API that extends beyond simple enumeration while maintaining type safety.

## Type Safety Across Package Boundaries

The phantom type system ensures type safety even across package boundaries:

```go
// This compiles because the function accepts the correct type
process_payment(PaymentMethod.CREDIT)

// This won't compile - type mismatch
// process_payment(CurrencyType.FIAT) 
```

Even though both values are ultimately `enum.Value[T]` instances, the phantom type parameter prevents mixing different enum types.

## Migration Strategy

When migrating existing code to use this pattern:

1. Identify all string constants or iota-based enums that represent the same conceptual type

2. Create a dedicated package for each enum type

3. Convert each enum value to use the type-safe pattern

4. Update function signatures to specify the enum type they accept

5. Add validation for values coming from external sources

This incremental approach allows for gradual adoption while maintaining backward compatibility.

## Why This Pattern Matters

The package-based type-safe enum pattern offers a compelling combination of benefits:

1. **Type Safety**: Prevents mixing different enum types at compile time

2. **Self-Documentation**: Package imports and function signatures clearly indicate which enum types are in use

3. **IDE Support**: Autocompletion shows available enum values for each type

4. **Runtime Validation**: Easily validate external values before using them

5. **Organization**: Clean package structure keeps related code together

6. **Evolution**: Enums can evolve independently with their own version history

By organizing enums into their own packages and using the type-safe implementation, we create code that is safer, more maintainable, and more pleasant to work with.

## Type Design and Interfaces

### Type Visibility

Go's type system provides powerful tools for API design and encapsulation. Our approach emphasizes clear public interfaces while maintaining strong encapsulation of implementation details.

1. Make types public by default:
   ```go
   // Good:
   type BlockProcessor struct {
       current_offset int64
       parent_hash   string
   }

   // Bad - unnecessary privacy:
   type blockProcessor struct {
       current_offset int64
       parent_hash   string
   }
   ```

2. Keep struct fields private by default:
   ```go
   // Good:
   type BlockProcessor struct {
       current_offset int64  // Private field
       parent_hash   string  // Private field
   }

   // Bad:
   type BlockProcessor struct {
       CurrentOffset int64  // Unnecessarily public
       ParentHash   string  // Unnecessarily public
   }
   ```

3. Use private types only when absolutely necessary, through either:
   ```go
   // Anonymous struct (preferred for truly internal types):
   var config = struct {
       retry_count    int
       timeout_ms     int64
   }{
       retry_count:    3,
       timeout_ms:     5000,
   }

   // Or camelCase naming (when type must be named but internal):
   type blockMetadata struct {
       hash      string
       timestamp int64
   }
   ```

### Interface Design

Interfaces form the backbone of our public APIs. They should be designed thoughtfully to provide clear, minimal contracts while hiding implementation details.

1. Prefix interfaces with "I" to clearly distinguish them from implementations:
   ```go
   // Good:
   type IBlockProcessor interface {
       ProcessChunk(data []byte) error  // Public method in PascalCase
       Status() string                  // Property accessor in PascalCase
   }

   // Bad:
   type BlockProcessor interface {       // Missing "I" prefix
       process_chunk(data []byte) error  // Wrong - interface methods cannot be private
       get_status() string              // Wrong - interface methods cannot be private
   }
   ```

2. Handle private functionality in the implementation, not the interface:
   ```go
   // Good:
   type IBlockProcessor interface {
       ProcessChunk(data []byte) error  // Public interface method
   }

   type blockProcessor struct {
       chunk_size int64
   }

   func (self *blockProcessor) ProcessChunk(data []byte) error {
       // Public implementation of interface method
       return self.process_chunk_internal(data)  // Private helper function
   }

   func (self *blockProcessor) process_chunk_internal(data []byte) error {
       // Private implementation details
       return nil
   }

   // Bad:
   type IBlockProcessor interface {
       ProcessChunk(data []byte) error
       process_internal(data []byte) error  // Wrong - private methods don't belong in interfaces
   }
   ```

2. Return interfaces from constructors to hide implementations:
   ```go
   // Good:
   func NewBlockProcessor(config *Config) IBlockProcessor {
       return &blockProcessor{config: config}
   }

   // Bad:
   func NewBlockProcessor(config *Config) *BlockProcessor {
       return &BlockProcessor{config: config}
   }
   ```

3. Design properties as methods without "Get" prefixes:
   ```go
   // Good:
   type IBlockInfo interface {
       Hash() string           // Interface methods use PascalCase
       ParentHash() string     // No Get prefix needed
       Height() int64         // Properties are nouns
   }

   // Bad:
   type IBlockInfo interface {
       GetHash() string          // Don't use Get prefix
       get_parent_hash() string  // Don't use snake_case in interfaces
       Height() int64           // This one is correct
   }

   // The implementation can then use snake_case internally:
   type blockInfo struct {
       hash_value    string
       parent_hash   string
       block_height  int64
   }

   func (self *blockInfo) Hash() string {
       return self.hash_value
   }

   func (self *blockInfo) ParentHash() string {
       return self.parent_hash
   }

   func (self *blockInfo) Height() int64 {
       return self.block_height
   }
   ```

### Naming Consistency

When implementing types and interfaces, maintain consistent naming conventions that clearly distinguish between different elements of our code:

1. Type names use specific conventions based on visibility:
   ```go
   // Public types use PascalCase
   type BlockProcessor struct {
       internal_count int64
   }

   // Private types use camelCase
   type blockMetadata struct {
       block_hash string
   }
   ```

2. All fields, methods, and variables use snake_case regardless of where they appear:
   ```go
   // Good - consistent snake_case for fields
   type blockProcessor struct {
       current_state   processState
       retry_count     int64
       error_channel   chan error
   }

   // Bad - mixing naming conventions
   type blockProcessor struct {
       currentState    processState  // Should be current_state
       retryCount     int64         // Should be retry_count
       errorChannel   chan error    // Should be error_channel
   }
   ```

### Visual Alignment and Readability

Code should be visually aligned to make relationships clear and improve readability:

1. Align struct fields with consistent spacing:
   ```go
   type BlockProcessor struct {
       current_state    processState  // Fields aligned
       retry_count     int64         // for visual
       error_channel   chan error    // clarity
   }
   ```

2. Align related variable declarations:
   ```go
   var (
       source_buffer  []byte       // Related buffers
       target_buffer  []byte       // aligned for
       error_buffer   []byte       // easy comparison
   )
   ```

3. Align method parameters for complex signatures:
   ```go
   func process_block(
       ctx         context.Context,
       block_data  []byte,
       options    *Options,
   ) (*ProcessedBlock, error) {
       // Implementation
   }
   ```

## Safety

### Assertions

Assertions enforce programming invariants that must always be true. Unlike error handling, which deals with expected failure cases, assertions catch programming errors that should never occur in correct code.

Using assertions:
```go
// Good: Enforcing a programming contract
func NewProcessor(config *Config) IProcessor {
    assert.NotNil(config, fmt.Errorf("config cannot be nil"))
    return &processor{config: config}
}

// Bad: Handling an expected error condition
func process_data(data []byte) error {
    assert.NotNil(data, fmt.Errorf("data cannot be nil"))
    // Should instead return an error for nil data
}
```

Key principles for assertions:
1. Use them only for true invariants - conditions that must always be true
2. Provide clear error messages that help identify the violated invariant
3. Remember that assertions will panic - they are not for handling recoverable errors

### Control Flow

1. Use simple, explicit control flow
2. Put limits on everything - all loops must have fixed upper bounds
3. Hard limit of 70 lines per function
4. Prefer early exits over nested conditionals:
   ```go
   // Prefer:
   func process_user_request(user *User, request *Request) error {
       if !user.is_authenticated() {
           return fmt.Errorf(ERR_NOT_AUTHENTICATED)
       }
       // We know user is authenticated here
       
       if !user.has_permission(request.permission) {
           return fmt.Errorf(ERR_INSUFFICIENT_PERMISSIONS)
       }
       // We know user is authenticated and has permission
       
       if request.is_expired() {
           return fmt.Errorf(ERR_REQUEST_EXPIRED)
       }
       // We know: user is authenticated, has permission, request is valid
       
       return process_valid_request(user, request)
   }

   // Over:
   func process_user_request(user *User, request *Request) error {
       if user.is_authenticated() {
           if user.has_permission(request.permission) {
               if !request.is_expired() {
                   return process_valid_request(user, request)
               } else {
                   return fmt.Errorf(ERR_REQUEST_EXPIRED)
               }
           } else {
               return fmt.Errorf(ERR_INSUFFICIENT_PERMISSIONS)
           }
       } else {
           return fmt.Errorf(ERR_NOT_AUTHENTICATED)
       }
   }
   ```

5. Add comments after each early exit to document the guaranteed state:
   - These comments serve as progressive assertions about what we know to be true
   - They help readers understand the narrowing of possibilities
   - They make the code self-documenting about its invariants

6. Keep the happy path at the end and at the minimum indentation level:
   - This makes the main logic easy to find
   - It ensures all preconditions have been checked
   - It keeps the core logic clean and focused

## Error Handling

### Error Creation and Wrapping

1. Define error messages as constants:
   ```go
   const (
       ERR_INVALID_INPUT = "invalid input: %s"
       ERR_DECODE_FAILED = "failed to decode: %w"
   )
   ```

2. Always wrap errors with context:
   ```go
   if err != nil {
       return fmt.Errorf("failed to process block: %w", err)
   }
   ```

3. Use error channels for asynchronous operations:
   ```go
   type IErrorChanneler interface {
       ErrorChannel() <-chan error     // Interface methods use PascalCase
   }

   // Implementation still uses snake_case internally
   type errorChanneler struct {
       error_chan chan error
   }

   func (self *errorChanneler) ErrorChannel() <-chan error {
       return self.error_chan
   }
   ```

## Documentation

### Function Documentation

Document all exported functions thoroughly:
```go
// process_block_data transforms raw block data according to the specified options.
// It enforces size limits and performs validation before processing.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - block_data: Raw block data to process
//   - options: Processing configuration options
//
// Returns:
//   - *ProcessedBlock: The transformed block data
//   - error: Any error encountered during processing
//
// Thread Safety: Safe for concurrent use
// Performance: O(n) time complexity, O(1) space complexity
func (self *Processor) process_block_data(
    ctx context.Context,
    block_data []byte,
    options *Options,
) (*ProcessedBlock, error)
```

### Package Documentation

Every package should have comprehensive documentation:
```go
// Package processor provides robust data processing capabilities with
// strict safety guarantees and predictable performance characteristics.
//
// Key Features:
//   - Thread-safe processing
//   - Bounded memory usage
//   - Configurable timeout handling
//   - Comprehensive error reporting
//
// Usage Example:
//   processor := NewProcessor(options)
//   result, err := processor.process_data(ctx, input_data)
package processor
```

## Testing

### Test Organization

Use table-driven tests with clear naming:
```go
func TestProcessBlock(t *testing.T) {
    test_cases := []struct {
        name           string
        input_data    []byte
        want_result   string
        want_error    bool
    }{
        {
            name:        "valid_block",
            input_data: valid_test_data,
            want_result: "processed",
            want_error:  false,
        },
        // More test cases...
    }
    
    for _, tc := range test_cases {
        t.Run(tc.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Performance

### Design-Time Optimization

1. Consider performance during design - this is where the 1000x wins happen
2. Perform back-of-the-envelope calculations for:
   - Network bandwidth and latency
   - Disk I/O patterns
   - Memory access patterns
   - CPU utilization
3. Use batching to amortize costs
4. Profile early and often

## Evolution

This style guide is a living document that evolves with our project's needs while maintaining our commitment to safety, performance, and developer experience. Regular reviews and updates ensure it continues to serve the development team effectively.

Remember: The best time to solve performance is in the design phase, which is precisely when we can't measure or profile. Work with mechanical sympathy, like a carpenter working with the grain.

# Debugging Practices

## Core Philosophy

Our debugging approach follows the same principles as our main codebase: safety, performance, and developer experience. When adding debug code, we maintain a clear separation between production and debugging functionality while ensuring that debug information is comprehensive and actionable.

## Debug Package Usage

The debug package provides standardized logging and inspection capabilities. All temporary debugging code must be clearly marked for removal using the `TODO_REMOVE_debug` import alias. This serves three critical purposes:

1. Makes debug code instantly recognizable
2. Facilitates complete removal before production deployment
3. Enables systematic auditing of debug statements

Example of proper debug import and usage:
```go
import (
    "fmt"
    "io"
    
    TODO_REMOVE_debug "mdh/debug"
)

func process_data(input []byte) error {
    TODO_REMOVE_debug.Printf("processing %d bytes", len(input))
    // ... implementation
}
```

## Debug Wrappers

The debug package provides wrapper types for standard interfaces like io.Reader and io.Writer. These wrappers enable transparent inspection of data flow without modifying the underlying code structure.

### Reader Wrapper Usage

When debugging data reading operations, wrap the reader with a debug reader:
```go
func read_data(r io.Reader) error {
    // Create named debug reader for clear logging
    debug_reader := TODO_REMOVE_debug.DebugReader("cache_read", r)
    
    // Use debug reader in place of original reader
    decoder := gob.NewDecoder(debug_reader)
    return decoder.Decode(&data)
}
```

### Writer Wrapper Usage

Similarly, wrap writers when debugging output operations:
```go
func write_data(w io.Writer, data []byte) error {
    // Create named debug writer for clear logging
    debug_writer := TODO_REMOVE_debug.DebugWriter("cache_write", w)
    
    // Use debug writer in place of original writer
    encoder := gob.NewEncoder(debug_writer)
    return encoder.Encode(data)
}
```

## Debug Message Formatting

Debug messages should provide clear context and relevant data values. They should follow these principles:

1. Use complete sentences in lowercase with proper punctuation
2. Include relevant numeric values (sizes, counts, etc.)
3. Reference specific identifiers (URLs, filenames, etc.)
4. Follow a consistent format within each component

Examples of good debug messages:
```go
TODO_REMOVE_debug.Printf("cache file size: %d bytes", size)
TODO_REMOVE_debug.Printf("reading cache for url: %s", url)
TODO_REMOVE_debug.Printf("added new cache entry for url: %s", url)
```

## Progressive Debugging

When investigating complex issues, add debug statements that progressively narrow down the problem area. This creates a breadcrumb trail that helps identify exactly where issues occur:

```go
func process_complex_operation() error {
    TODO_REMOVE_debug.Printf("starting complex operation")
    
    data, err := read_input()
    if err != nil {
        TODO_REMOVE_debug.Printf("failed to read input: %v", err)
        return fmt.Errorf("failed to read input: %w", err)
    }
    TODO_REMOVE_debug.Printf("successfully read %d bytes", len(data))
    
    result, err := transform_data(data)
    if err != nil {
        TODO_REMOVE_debug.Printf("failed to transform data: %v", err)
        return fmt.Errorf("failed to transform data: %w", err)
    }
    TODO_REMOVE_debug.Printf("successfully transformed data")
    
    return write_output(result)
}
```

## Debug Code Removal

Debug code should be removed before code is merged to production. This can be automated by searching for the `TODO_REMOVE_debug` import statement. The import name makes it impossible to accidentally leave debug code in production, as it won't compile if the import is removed.

Follow these steps when removing debug code:

1. Remove all TODO_REMOVE_debug import statements
2. Remove all debug wrapper usage
3. Remove any variables only used for debugging
4. Verify the code compiles and tests pass
5. Verify no debug output appears in logs

## When to Add Debug Code

Add debug statements when:

1. Investigating complex data flow issues
2. Debugging concurrent operations
3. Tracing resource usage patterns
4. Understanding system boundaries (file I/O, network calls)
5. Analyzing performance bottlenecks

The goal is to add enough debug information to understand the system's behavior without overwhelming the output with unnecessary details.

## A Note on Performance

While debug code can impact performance, our priority during debugging is clarity and completeness of information. The `TODO_REMOVE_debug` prefix ensures this code won't impact production performance.

Remember: A slow program that helps us understand the problem is better than a fast program that leaves us guessing.
