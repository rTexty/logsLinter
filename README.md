# logsLinter

Production-ready Go analyzer for validating log messages in `log/slog` and `go.uber.org/zap` codebases.

## Goals

- Enforce consistent log message style
- Catch non-English or non-ASCII log text
- Flag decorative punctuation and emoji in log messages
- Prevent accidental logging of potentially sensitive data
- Integrate with standard Go analysis tooling and `golangci-lint`

## Planned Rules

- Log messages must start with a lowercase letter
- Log messages must contain English ASCII text only
- Log messages must not contain decorative special characters or emoji
- Log messages must not contain sensitive keywords such as `password`, `token`, or `secret`

## Development

### Requirements

- Go 1.23+

### Commands

```bash
# Build
go build ./...

# Test
go test ./... -race -count=1

# Format
gofmt -w .

# Tidy dependencies
go mod tidy
```

## Status

The repository is currently being scaffolded. The standalone analyzer entrypoint, analyzer package, test fixtures, and `golangci-lint` integration examples are implemented incrementally from the feature plan in `Context/Features/001-LogsLinter/Steps.md`.
