# Installation

## Requirements

- Go 1.21 or later
- Aha.io account with API access

## Install as Library

Add `aha-go` to your Go project:

```bash
go get github.com/grokify/aha-go
```

Import in your code:

```go
import aha "github.com/grokify/aha-go"
```

## Install CLI Tool

Install the `aha` command-line tool:

```bash
go install github.com/grokify/aha-go/cmd/aha@latest
```

Verify installation:

```bash
aha --version
```

## Build from Source

Clone and build locally:

```bash
git clone https://github.com/grokify/aha-go.git
cd aha-go
go build ./cmd/aha
```

## Next Steps

- [Set up authentication](authentication.md)
- [Quick start guide](quickstart.md)
