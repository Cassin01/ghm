# ghm - GitHub Manager

`ghm` is an extended version of `ghq` that allows managing multiple instances of the same repository in different directories.

## Features

- Manage multiple instances of the same repository
- Instances are distinguished by numbers (`.1`, `.2`, `.3`...)
- `ghq`-compatible basic functionality

## Installation

### From Source

```bash
# Using Go
go build -o ghm ./cmd/ghm

# Using Make
make build
```

### From Release

Download the latest release from the [releases page](https://github.com/Cassin01/ghm/releases).

## Usage

### Clone Repository

```bash
# Basic clone
ghm get https://github.com/user/repo

# Clone with specific instance number
ghm get https://github.com/user/repo -n 1

# Auto-assign next available instance number
ghm get https://github.com/user/repo --auto
```

### List Repositories

```bash
# List all repositories
ghm list

# List all instances of a specific repository
ghm list github.com/user/repo
```

### Show Root Directory

```bash
ghm root
```

### Remove Repository

```bash
# Remove main instance
ghm remove github.com/user/repo

# Remove specific instance
ghm remove github.com/user/repo.1
```

## Directory Structure

```
$GHM_ROOT/
├── github.com/
│   └── user/
│       ├── repo/          # Main instance
│       ├── repo.1/        # First instance
│       ├── repo.2/        # Second instance
│       └── ...
└── gitlab.com/
    └── ...
```

## Configuration

### Environment Variables

- `GHM_ROOT`: Repository management directory (default: `~/ghm`)

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Development workflow (format, lint, test, build)
make dev
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make coverage

# Run linter
make lint
```

### Other Commands

```bash
# Clean build artifacts
make clean

# Install to GOPATH/bin
make install

# Show all available commands
make help
```

## License

MIT License
