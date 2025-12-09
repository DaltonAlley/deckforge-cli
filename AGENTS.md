# AGENTS.md

## Build/Lint/Test Commands
- **Build**: `go build -o deckforge-cli .`
- **Test all**: `go test ./...`
- **Test single**: `go test -run TestName ./package`
- **Format**: `gofmt -w .`
- **Vet**: `go vet ./...`

## Code Style Guidelines

### Imports
- Standard library imports first
- Third-party imports second
- Local imports last
- Blank line between groups

### Naming Conventions
- Exported types/functions: PascalCase (Card, FindCardByID)
- Unexported: camelCase
- Constants: PascalCase

### Error Handling
- Functions return (result, error)
- Use fmt.Errorf for wrapping
- Caller handles errors appropriately

### Types & Structs
- Use JSON struct tags for API structs
- Meaningful type names
- Group related fields

### Testing
- Use testify/require for assertions
- Test functions: TestFunctionName
- Integration tests in separate files

### Logging
- Use zerolog for structured logging
- Log errors with context

### Formatting
- Use gofmt standard formatting
- No trailing whitespace
- Consistent indentation