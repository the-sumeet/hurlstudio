# Embedding Hurl Binary Setup

## Directory Structure

Create this directory structure in your project:

```
hurlstudio/
├── binaries/
│   ├── darwin-amd64/
│   │   └── hurl
│   ├── darwin-arm64/
│   │   └── hurl
│   ├── linux-amd64/
│   │   └── hurl
│   ├── linux-arm64/
│   │   └── hurl
│   ├── windows-amd64/
│   │   └── hurl.exe
│   └── windows-arm64/
│       └── hurl.exe
├── hurl.go
└── ...
```

## How to Get Hurl Binaries

### Option 1: Download from Official Releases

1. Go to https://github.com/Orange-OpenSource/hurl/releases
2. Download the binaries for each platform you want to support
3. Extract and place them in the appropriate directories

### Option 2: Build from Source

For each platform:

```bash
# Clone Hurl repository
git clone https://github.com/Orange-OpenSource/hurl.git
cd hurl

# Build for your target platform
cargo build --release

# The binary will be in target/release/hurl
```

### Option 3: Install via Package Manager (Current Platform Only)

**macOS (Homebrew):**
```bash
brew install hurl
# Binary is at /opt/homebrew/bin/hurl (Apple Silicon)
# or /usr/local/bin/hurl (Intel)
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install hurl

# Or use cargo
cargo install hurl
```

**Windows:**
```bash
# Using Scoop
scoop install hurl

# Using Chocolatey
choco install hurl
```

## Quick Setup for Development (macOS only)

If you just want to get started quickly on macOS:

```bash
# Install hurl
brew install hurl

# Find the binary
which hurl  # e.g., /opt/homebrew/bin/hurl

# Copy to your project structure
mkdir -p binaries/darwin-arm64
cp $(which hurl) binaries/darwin-arm64/

# For Intel Macs
mkdir -p binaries/darwin-amd64
cp $(which hurl) binaries/darwin-amd64/
```

## Alternative: System Hurl (No Embedding)

If you don't want to embed the binary, the code will automatically fall back to using `hurl` from the system PATH if available.

Just make sure hurl is installed:
```bash
# macOS
brew install hurl

# The app will find it automatically
```

## Usage in Your App

The `hurl.go` module provides these functions:

```go
// Run a hurl file
output, err := app.RunHurl("/path/to/file.hurl")

// Run with options
output, err := app.RunHurlWithOptions("/path/to/file.hurl", []string{"--verbose", "--json"})
```

## Minimum Requirement

For development on macOS, you only need:
- `binaries/darwin-arm64/hurl` (Apple Silicon)
- OR `binaries/darwin-amd64/hurl` (Intel)

For distribution, include binaries for all platforms you want to support.
