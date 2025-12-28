# Justfile for minilisp

# Default recipe - show available commands
default:
    @just --list

# Run tests with coverage percentage
test-coverage:
    go test -cover ./...

# Generate coverage profile
coverage:
    go test -coverprofile=coverage.out ./...

# View coverage summary in terminal
coverage-func: coverage
    go tool cover -func=coverage.out

# Open HTML coverage report in browser
coverage-html: coverage
    go tool cover -html=coverage.out

# Generate coverage with count mode (shows execution frequency)
coverage-count:
    go test -coverprofile=coverage.out -covermode=count ./...
    go tool cover -html=coverage.out

# Run all tests
test:
    go test ./...

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Clean coverage files
clean:
    rm -f coverage.out

# Run tests and show coverage report in one command
check: coverage coverage-func
