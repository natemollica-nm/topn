# topn

Fast file size scanner to find and remove large files on Unix-like systems with beautiful terminal UI.

## Features

- **🎨 Beautiful TUI** with interactive file selection and removal
- **⚡ Fast concurrent scanning** with configurable worker pools
- **💾 Memory efficient** using min-heap to track only top N files
- **🗑️ Interactive file removal** with confirmation and selection
- **📏 Flexible size filtering** (supports K, M, G, T suffixes)
- **🚫 Pattern exclusion** using glob patterns
- **🖥️ Cross-platform** support for Linux, macOS, and Unix systems
- **🎯 Dual modes** - Classic CLI and modern TUI

## Installation

### Homebrew

```bash
brew install natemollica/tap/topn
```

### Go Install

```bash
go install github.com/natemollica/topn/cmd/topn@latest
```

### Download Binary

Download from [releases](https://github.com/natemollica/topn/releases).

## Usage

### Interactive TUI Mode (Recommended)

```bash
# Launch interactive terminal UI
topn -tui

# TUI with custom settings
topn -tui -dir /var/log -min 100M -top 20

# Remove mode automatically uses TUI
topn -remove
```

**TUI Controls:**
- `Space` - Select/deselect files
- `Enter` - Remove selected files
- `r` - Rescan directory
- `q` - Quit

### Classic CLI Mode

```bash
# Find top 50 files >= 1GB in home directory
topn

# Scan specific directory with custom threshold
topn -dir /var/log -min 100M -top 20

# Exclude patterns
topn -exclude "*.log" -exclude "node_modules" -exclude "/tmp/*"

# Custom worker count
topn -workers 8
```

### Options

- `-dir`: Root directory to scan (default: $HOME)
- `-min`: Minimum file size threshold (default: 1G)
- `-top`: Number of largest files to keep (default: 50)
- `-workers`: Number of concurrent workers (default: 4*GOMAXPROCS)
- `-exclude`: Glob patterns to exclude (repeatable)
- `-remove`: Enable interactive file removal (uses TUI)
- `-tui`: Force interactive terminal UI mode

### Size Format

Supports standard size suffixes:
- `K`, `KB`: Kilobytes (1024 bytes)
- `M`, `MB`: Megabytes (1024² bytes)
- `G`, `GB`: Gigabytes (1024³ bytes)
- `T`, `TB`: Terabytes (1024⁴ bytes)

## Examples

```bash
# Interactive cleanup with beautiful UI
topn -tui -dir ~/Downloads -min 100M

# Quick CLI scan
topn -dir /var/log -min 50M

# Exclude common directories
topn -exclude ".git" -exclude "node_modules" -exclude "*.tmp"

# Safe interactive removal
topn -remove -min 500M
```

## Development

### Building

```bash
# Build binary
make build

# Install locally
make install

# Run tests
make test
```

### Releasing

```bash
# Dry run release (no tags created)
make release-dry-run

# Interactive release with automated versioning
make release

# Manual tag creation
make tag-release
```

The release process automatically:
- Validates git state and runs tests
- Creates version tags with semantic versioning
- Builds cross-platform binaries
- Publishes to GitHub releases
- Updates Homebrew tap

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## License

MIT