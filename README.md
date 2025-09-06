![Scruffy](docs/logo/logo.mini.128x128.png)

# Scruffy

A Go CLI tool for cleaning Cloudflare IP Access rules.

## Features

- **Clean all records**: Remove all IP Access rules from a zone
- **Clean by prefix**: Remove records that start with a specific prefix
- **Clean specific record**: Remove a specific IP address, CIDR block, or ASN
- **Clean by description**: Remove records containing specific text in their description/notes
- **Dry run mode**: Preview what would be deleted without making changes
- **Zone resolution**: Use zone name instead of zone ID for convenience

## Installation

### From source

```bash
go install github.com/scottbrown/scruffy/cmd/scruffy@latest
```

### Using Go Task

```bash
git clone https://github.com/scottbrown/scruffy.git
cd scruffy
task install
```

## Authentication

Scruffy requires a Cloudflare API token with Zone:Zone:Read and Zone:Zone Settings:Edit permissions.

### Using environment variable (recommended)

```bash
export CLOUDFLARE_API_TOKEN="your-api-token-here"
```

### Using command line flag (not recommended for security)

```bash
scruffy --token "your-api-token-here" [command]
```

## Usage

### Basic Commands

```bash
# Clean all IP Access rules in a zone
scruffy --zone-id abc123 clean all

# Clean all rules with dry run (preview only)
scruffy --zone-id abc123 clean all --dry-run

# Clean rules starting with specific prefix
scruffy --zone-name example.com clean prefix "192.168."

# Clean a specific IP/CIDR/ASN
scruffy --zone-id abc123 clean target "203.0.113.0/24"

# Clean rules by description
scruffy --zone-id abc123 clean description "temporary block"
```

### Zone Specification

You can specify a zone using either:

- `--zone-id`: Cloudflare Zone ID (faster)
- `--zone-name`: Domain name (will be resolved to Zone ID)

```bash
# Using Zone ID
scruffy --zone-id abc123def456 clean all

# Using Zone name (domain)
scruffy --zone-name example.com clean all
```

### Dry Run Mode

Always test your commands with `--dry-run` first:

```bash
scruffy --zone-id abc123 clean all --dry-run
```

## Examples

### Clean all temporary blocks

```bash
# Preview what would be deleted
scruffy --zone-name mydomain.com clean description "temp" --dry-run

# Actually delete them
scruffy --zone-name mydomain.com clean description "temp"
```

### Clean specific IP ranges

```bash
# Clean all 10.x.x.x private IPs
scruffy --zone-id abc123 clean prefix "10."

# Clean specific CIDR block
scruffy --zone-id abc123 clean target "192.168.1.0/24"
```

### Clean specific ASN

```bash
scruffy --zone-id abc123 clean target "AS64496"
```

## Record Types Supported

- **IP addresses**: `203.0.113.1`
- **CIDR blocks**: `203.0.113.0/24`
- **ASN (Autonomous System Numbers)**: `AS64496`

## Development

### Prerequisites

- Go 1.19+
- [Task](https://taskfile.dev/) (optional, for task runner)

### Building

```bash
# Using Go
go build -o .build/scruffy ./cmd/scruffy

# Using Task (builds to .build/scruffy)
task build

# Build for all platforms
task build-all

# Create release artifacts
task release
```

### Testing

```bash
# Run tests
task test

# Run tests with coverage (outputs to .test/ directory)
task coverage

# View coverage report
open .test/coverage.html
```

### Code Quality

```bash
# Format code
task format

# Run linting
task lint

# Security analysis
task sast

# Vulnerability scanning
task vuln
```

### Development

```bash
# Download and tidy dependencies
task deps

# Clean all build artifacts
task clean
```

## License

MIT License
