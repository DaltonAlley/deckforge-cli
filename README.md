# DeckForge CLI

> Generate printable MTG deck PDFs from Archidekt CSV exports

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)

A command-line tool for creating professional-quality Magic: The Gathering deck PDFs with configurable bleed margins and clean progress tracking.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Smart Output Naming**: Automatically names PDFs after input CSV files
- **Professional Printing**: Configurable bleed margins for clean cutting
- **Clean Progress Display**: Indexed progress tracking `[current/total]` format
- **Error Resilience**: Graceful handling of invalid cards with detailed reporting
- **Flexible Output**: Custom filenames and quiet mode for automation
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Installation

### Prerequisites

- Go 1.21 or later
- Internet connection for card image downloads

### Install from Source

```bash
git clone https://github.com/daltonalley/deckforge-cli.git
cd deckforge-cli
go build -o deckforge-cli ./cmd/deckforge
```

### Verify Installation

```bash
./deckforge-cli --help
```

## Quick Start

```bash
# Basic usage - output defaults to deck.pdf
deckforge deck.csv

# Custom bleed for professional printing
deckforge --bleed 3.0 deck.csv
```

## Usage

```bash
deckforge [options] <decklist.csv>

Options:
  -o, --output string    Output PDF filename (defaults to CSV name)
  --bleed float          Bleed margin in mm around each card (default: 3.0)
  --quiet                Suppress progress output
  -h, --help             Show help
  -v, --version          Show version

Arguments:
  decklist.csv    Path to Archidekt CSV export file
```

## Configuration

### Bleed Margins

Control extra margin around cards for professional printing:

- **Default**: 3.0mm (recommended for most printers)
- **Range**: 0.0mm (no bleed) to 10.0mm (maximum bleed)
- **Usage**: Cards are centered within the bleed area
- **Purpose**: Provides safe area for cutting and prevents white edges

### Progress Display

- **Format**: `[current/total] Operation description`
- **Operations Tracked**: Card fetching, page generation, PDF assembly
- **Quiet Mode**: Use `--quiet` to suppress all progress output
- **Error Display**: Errors appear in status area with card ID context

## Examples

### Basic Usage

```bash
# Automatic output naming (deck.csv → deck.pdf)
deckforge deck.csv

# Process multiple decks
deckforge commander.csv
deckforge standard.csv
```

### Advanced Options

```bash
# Custom output filename
deckforge --output my_commander_deck.pdf commander.csv

# Professional printing with bleed margins
deckforge --bleed 5.0 deck.csv

# Quiet mode for scripts/automation
deckforge --quiet deck.csv

# Combined options for production use
deckforge --bleed 3.0 --output production_deck.pdf deck.csv
```

### Error Handling

The tool gracefully handles various error conditions:

```bash
# Invalid cards are skipped with error reporting
deckforge problematic_deck.csv
# Output: [5/10] Assembling PDF
#         ❌ PDF generation completed with 2 error(s):
#            • invalid-card-1: Error: 404 Not Found
#            • invalid-card-2: Error: Network timeout
#            • Successfully processed 8 cards
```

## Contributing

We welcome contributions! Please follow these guidelines:

### Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/deckforge-cli.git
   cd deckforge-cli
   ```
3. Install Go 1.21+
4. Build the project:
   ```bash
   go build -o deckforge-cli ./cmd/deckforge
   ```
5. Run tests:
   ```bash
   go test ./...
   ```

### Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature
   ```
2. Write tests for new functionality
3. Ensure all tests pass:
   ```bash
   go test ./...
   go vet ./...
   go fmt ./...
   ```
4. Follow Go best practices and existing code style
5. Update documentation as needed

### Pull Requests

- Provide clear description of changes
- Reference any related issues
- Ensure CI checks pass
- Squash commits for clean history
- Update CHANGELOG.md for user-facing changes

### Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Include unit tests for new features
- Handle errors appropriately
- Follow the guidelines in `AGENTS.md`
---

Built with ❤️ for the MTG community
