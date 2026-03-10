# logsLinter

[![CI](https://github.com/rTexty/logsLinter/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/rTexty/logsLinter/actions/workflows/ci.yml)
[![CodeQL](https://github.com/rTexty/logsLinter/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/rTexty/logsLinter/actions/workflows/codeql.yml)
[![Release](https://github.com/rTexty/logsLinter/actions/workflows/release.yml/badge.svg)](https://github.com/rTexty/logsLinter/actions/workflows/release.yml)
[![Latest Release](https://img.shields.io/github/v/release/rTexty/logsLinter)](https://github.com/rTexty/logsLinter/releases)

Production-ready Go analyzer for validating log messages in `log/slog` and `go.uber.org/zap` codebases.

The analyzer is implemented as `golang.org/x/tools/go/analysis` and can run both as a standalone binary and through the `golangci-lint` module plugin workflow.

## Status Panel

| Signal | Status |
| --- | --- |
| CI | formatting, vet, tests, build, repository lint |
| Security | CodeQL, dependency review, Dependabot, SECURITY policy |
| Releases | tag-based GitHub Release workflow with packaged binaries and checksums |
| Tooling | standalone CLI, golangci-lint module plugin examples |

## Supported APIs

- `log/slog` top-level calls: `Debug`, `Info`, `Warn`, `Error`
- `log/slog` context variants: `DebugContext`, `InfoContext`, `WarnContext`, `ErrorContext`
- `log/slog` structured variants: `Log`, `LogAttrs`
- `*slog.Logger` methods, including chained `With(...)` and `WithGroup(...)`
- `*zap.Logger` methods: `Debug`, `Info`, `Warn`, `Error`
- `*zap.SugaredLogger` methods with explicit message argument: `Debugw`, `Infow`, `Warnw`, `Errorw`

## Goals

- Enforce consistent log message style
- Catch non-English or non-ASCII log text
- Flag decorative punctuation and emoji in log messages
- Prevent accidental logging of potentially sensitive data
- Integrate with standard Go analysis tooling and `golangci-lint`

## Rules

- Log messages must start with a lowercase letter
- Log messages must contain English ASCII text only
- Log messages must not contain decorative special characters or emoji
- Log messages must not contain sensitive keywords such as `password`, `token`, or `secret`

## Skipped Cases

- Non-literal messages such as variables, function calls, and `fmt.Sprintf(...)`
- Literal plus variable concatenation such as `"password: " + secret`
- `zap.SugaredLogger` print-style methods such as `Info(...)`, `Warn(...)`, `Error(...)`
- `zap.SugaredLogger` format-style methods such as `Infof(...)`, `Warnf(...)`, `Errorf(...)`

## Development

### Requirements

- Go 1.23+

### Commands

```bash
# Build
go build ./...

# Build standalone analyzer binary
go build -o ./bin/logslinter ./cmd/logslinter

# Test
go test ./... -race -count=1

# Format
gofmt -w .

# Tidy dependencies
go mod tidy
```

## Standalone Usage

Build the analyzer:

```bash
go build -o ./bin/logslinter ./cmd/logslinter
```

Run it directly on packages:

```bash
./bin/logslinter ./...
```

Or use it through `go vet` as a `vettool`:

```bash
go vet -vettool=$(pwd)/bin/logslinter ./...
```

The command exits with a non-zero status when diagnostics are reported or package loading fails.

Example diagnostics:

```text
internal/service/auth.go:42:18: log message must start with a lowercase letter
internal/service/auth.go:42:18: log message may contain sensitive data
```

Example:

```go
slog.Info("Starting auth token rotation")
```

This call reports:

- `log message must start with a lowercase letter`
- `log message may contain sensitive data`

## golangci-lint Module Plugin

The repository contains example plugin configuration for the current module plugin workflow:

- [.custom-gcl.yml](/Users/rtexty/Documents/MyProjects/Programming/logsLinter/.custom-gcl.yml) builds a custom `golangci-lint` binary with `logsLinter` linked in
- [.golangci.yml](/Users/rtexty/Documents/MyProjects/Programming/logsLinter/.golangci.yml) enables the custom linter as a module plugin

Build a custom `golangci-lint` binary:

```bash
golangci-lint custom
```

With the example config in this repository, the custom binary is written to `./bin/custom-golangci-lint`.

Run the custom binary:

```bash
./bin/custom-golangci-lint run ./...
```

Minimal local-path plugin config:

```yaml
version: v2.11.2
name: custom-golangci-lint
destination: ./bin

plugins:
	- module: github.com/rTexty/logsLinter
		path: .
```

```yaml
version: "2"

linters:
	default: none
	enable:
		- logslinter
	settings:
		custom:
			logslinter:
				type: module
				description: Validate literal slog and zap log messages with logsLinter.
```

## Verification

- Unit coverage exists for rule evaluation, extraction, diagnostics, and logger call inspection
- Integration coverage runs through `analysistest` fixtures for `slog`, `zap`, and mixed edge cases
- Current verification baseline is `go test ./...`

## Repository Automation

- GitHub Actions CI runs formatting checks, `go vet`, tests, build, and repository linting.
- CodeQL runs on pull requests and on a weekly schedule.
- Dependency review runs on pull requests.
- Dependabot tracks both Go modules and GitHub Actions.
- Releases are built from Git tags matching `v*` and publish packaged binaries plus `SHA256SUMS.txt`.
- Changelog and release note conventions are documented in `CHANGELOG.md` and `docs/release-policy.md`.
- Recommended GitHub branch ruleset settings are documented in `docs/github-ruleset.md`.

## Release Process

Create and push a semantic version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow will:

- run tests before packaging
- build `logslinter` for Linux, macOS, and Windows
- upload `.tar.gz` and `.zip` artifacts
- publish a GitHub Release with generated notes and `SHA256SUMS.txt`

Release note categorization is driven by pull request labels configured in `.github/release.yml`.

## Status

The analyzer MVP is functional: supported `slog` and `zap` calls are inspected, literal messages are validated against the four core rules, and diagnostics are covered by unit tests plus `analysistest` fixtures.
