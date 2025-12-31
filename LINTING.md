# Code Quality & Linting Guide

This project uses [golangci-lint](https://golangci-lint.run/) for comprehensive
code quality checks.

## Quick Start

### Install golangci-lint

**macOS (Homebrew):**

```bash
brew install golangci-lint
```

**Linux/Windows:**

```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.1
```

**Or use the Makefile (auto-installs):**

```bash
make lint
```

## Available Commands

### Linting

```bash
# Run all linters
make lint

# Run linters with auto-fix
make lint-fix
```

### Code Formatting

```bash
# Check if code is formatted
make fmt

# Auto-fix formatting issues
make fmt-fix
```

### Other Quality Checks

```bash
# Run go vet
make vet

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run all checks (fmt, vet, lint, test)
make check

# CI-friendly check (fails on any issue)
make ci-check
```

## Enabled Linters

The following linters are enabled for production-grade code quality:

### Code Quality

- **errcheck** - Check for unchecked errors
- **errorlint** - Error handling best practices
- **exportloopref** - Detect exporting loop variables
- **gci** - Import ordering
- **gofmt** - Code formatting
- **goimports** - Import management
- **golint** - Go linting
- **ineffassign** - Detect ineffectual assignments
- **misspell** - Spelling mistakes
- **nakedret** - Naked returns
- **nilerr** - Detect nil errors
- **noctx** - Missing context in function signatures
- **nolintlint** - Check nolint directives
- **unconvert** - Unnecessary conversions
- **unparam** - Unused parameters
- **unused** - Unused code
- **varcheck** - Unused variables
- **whitespace** - Trailing whitespace

### Security

- **gosec** - Security issues (medium severity/confidence)

### Static Analysis

- **govet** - Go vet checks
- **staticcheck** - Static analysis
- **structcheck** - Unused struct fields
- **typecheck** - Type checking

### Style

- **revive** - Fast, configurable linter with Go best practices

## Configuration

The linting configuration is in `.golangci.yml`. Key settings:

- **Timeout**: 5 minutes
- **Max issues per linter**: 50
- **Max same issues**: 3
- **Test files**: Included but with relaxed rules
- **Excluded directories**: `vendor/`, `tmp/`, `bin/`, `logs/`, `docs/`

## Pre-commit Hooks (Optional)

To run linting before each commit, install a pre-commit hook:

```bash
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
make check
EOF
chmod +x .git/hooks/pre-commit
```

Or use [pre-commit](https://pre-commit.com/):

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.60.1
    hooks:
      - id: golangci-lint
```

## CI/CD Integration

A GitHub Actions workflow is included at `.github/workflows/lint.yml` that runs:

- golangci-lint
- Tests
- Formatting checks
- go vet

The workflow runs on:

- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

## Common Issues & Fixes

### Unchecked Errors

```go
// ❌ Bad
result, _ := someFunction()

// ✅ Good
result, err := someFunction()
if err != nil {
    return err
}
```

### Missing Context

```go
// ❌ Bad
func GetUser(id string) (*User, error)

// ✅ Good
func GetUser(ctx context.Context, id string) (*User, error)
```

### Error Wrapping

```go
// ❌ Bad
return fmt.Errorf("failed: %v", err)

// ✅ Good
return fmt.Errorf("failed: %w", err)
```

### Import Ordering

```go
// ✅ Good (auto-fixed by goimports)
import (
    "context"
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/mrhpn/go-rest-api/internal/errors"
)
```

## Best Practices

1. **Run `make check` before committing** - Catches most issues early
2. **Fix auto-fixable issues** - Run `make lint-fix` and `make fmt-fix`
3. **Address warnings** - Even warnings can indicate potential issues
4. **Review CI results** - All PRs are automatically linted
5. **Keep dependencies updated** - Run `make tidy` regularly

## Disabling Linters (When Necessary)

If you need to disable a linter for a specific line:

```go
//nolint:errcheck
result, _ := someFunction()
```

For a specific reason:

```go
//nolint:errcheck // This error is intentionally ignored because...
result, _ := someFunction()
```

**Note**: Use `//nolint` sparingly and always document why.

## Resources

- [golangci-lint Documentation](https://golangci-lint.run/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
