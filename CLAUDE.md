# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Scruffy is a Go CLI tool for cleaning Cloudflare IP Access rules. It uses the official Cloudflare Go SDK to interact with the Cloudflare API and provides multiple filtering options for bulk cleanup operations.

## Development Commands

### Building and Testing
```bash
# Build binary (outputs to .build/scruffy)
task build
# or: go build -o .build/scruffy ./cmd/scruffy

# Run all tests
task test
# or: go test -v ./...

# Run tests with coverage (generates coverage.html)
task test-coverage

# Run specific test
go test -v -run TestFilterRulesByPrefix

# Run linter
task lint

# Clean build artifacts
task clean

# Install binary
task install
```

### Running the Tool
```bash
# Set API token (required)
export CLOUDFLARE_API_TOKEN="your-token"

# Test with dry run first
./.build/scruffy --zone-name example.com clean all --dry-run

# Clean operations
./.build/scruffy --zone-id abc123 clean all
./.build/scruffy --zone-name example.com clean prefix "192.168."
./.build/scruffy --zone-id abc123 clean target "203.0.113.0/24"
./.build/scruffy --zone-id abc123 clean description "temp"
```

## Architecture

### Package Structure
The project follows Go's standard CLI layout with business logic in the root package:

- **Root Package (`scruffy`)**: Contains all business logic and CLI commands
- **Main Entry Point**: `cmd/scruffy/main.go` - minimal main function that calls `scruffy.Execute()`
- **Command Structure**: Uses cobra for CLI with hierarchical subcommands

### Key Components

**Authentication & Client Management**:
- API token via `CLOUDFLARE_API_TOKEN` environment variable (preferred) or `--token` flag
- Zone specification via either `--zone-id` (direct) or `--zone-name` (resolved to ID)
- `Client` struct wraps the Cloudflare SDK with zone-specific operations

**Command Architecture**:
- `rootCmd`: Main command with global flags (token, zone-id, zone-name)
- `cleanCmd`: Parent command for all cleanup operations
- Subcommands: `all`, `prefix`, `target`, `description` - each with specific filtering logic
- All operations support `--dry-run` for safe preview

**Core Business Logic**:
- `AccessRule` struct: Normalized representation of Cloudflare access rules
- Filtering functions: `FilterRulesByPrefix()`, `FilterRulesByTarget()`, `FilterRulesByDescription()`
- `setupClient()`: Handles zone resolution and client initialization
- `deleteRules()`: Common deletion logic with dry-run support and progress reporting

### Data Flow
1. **Command Parsing**: Cobra processes CLI args and validates required flags
2. **Client Setup**: Zone name resolution (if needed) → API client creation
3. **Rule Retrieval**: Fetch all access rules from Cloudflare API
4. **Filtering**: Apply specific filters based on command (prefix, target, description)
5. **Execution**: Preview (dry-run) or delete filtered rules with progress feedback

### Testing Strategy
- Unit tests for filtering logic and command validation
- Mock-friendly design with interfaces for external dependencies
- Test coverage focused on business logic rather than API integration
- Current coverage: ~40% (primarily business logic and utilities)

### Dependencies
- **github.com/spf13/cobra**: CLI framework and command structure
- **github.com/cloudflare/cloudflare-go**: Official Cloudflare API client
- Standard library for HTTP, context, and string operations