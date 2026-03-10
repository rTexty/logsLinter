<div align="center">
  <h1>logsLinter</h1>
  <p><strong>Clear, safe, consistent Go log messages.</strong></p>
  <p>Static analysis for <code>log/slog</code> and <code>go.uber.org/zap</code></p>
  <p>
    <a href="https://github.com/rTexty/logsLinter/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/rTexty/logsLinter/actions/workflows/ci.yml/badge.svg?branch=main"></a>
    <a href="https://github.com/rTexty/logsLinter/actions/workflows/codeql.yml"><img alt="CodeQL" src="https://github.com/rTexty/logsLinter/actions/workflows/codeql.yml/badge.svg?branch=main"></a>
    <a href="https://github.com/rTexty/logsLinter/actions/workflows/release.yml"><img alt="Release" src="https://github.com/rTexty/logsLinter/actions/workflows/release.yml/badge.svg"></a>
    <a href="https://github.com/rTexty/logsLinter/releases"><img alt="Latest Release" src="https://img.shields.io/github/v/release/rTexty/logsLinter"></a>
    <a href="https://pkg.go.dev/github.com/rTexty/logsLinter"><img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/rTexty/logsLinter.svg"></a>
    <a href="https://goreportcard.com/report/github.com/rTexty/logsLinter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/rTexty/logsLinter"></a>
    <a href="https://github.com/rTexty/logsLinter/blob/main/LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/license-MIT-green.svg"></a>
  </p>
</div>

Literal log messages are one of the easiest places for drift, inconsistency, and accidental secrets to slip into a Go codebase. `logsLinter` turns that into an enforceable rule set for `log/slog` and `go.uber.org/zap`, with a standalone analyzer, `golangci-lint` module plugin support, and safe autofix coverage for the simplest lowercase-start violations.

Built on `golang.org/x/tools/go/analysis`, it is meant to fit into normal Go workflows instead of becoming its own ecosystem.

## Why Use It

- Keep log messages consistent across large teams and services.
- Catch obvious sensitive-data wording before it reaches production logs.
- Enforce log quality in CI instead of code review after the fact.
- Run the same analyzer as a standalone binary or through `golangci-lint`.

## What It Checks

- Log messages must start with a lowercase letter.
- Log messages must contain English ASCII text only.
- Log messages must not contain decorative special characters or emoji.
- Log messages must not contain sensitive keywords such as `password`, `token`, or `secret`.

For the lowercase-start rule, `logsLinter` also emits a safe `SuggestedFix` when the message is an interpreted string literal whose first rune is an ASCII uppercase letter.

## Quick Start

Build the standalone analyzer:

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

## Tech Spec Coverage

| Task from `tech_spec.pdf` | Status |
| --- | --- |
| Support `log/slog` log calls | Done |
| Support `go.uber.org/zap` log calls | Done |
| Enforce lowercase-start log messages | Done |
| Enforce English-only ASCII log messages | Done |
| Reject decorative special characters and emoji | Done |
| Detect potentially sensitive data keywords in log messages | Done |
| Integrate as a `go/analysis` analyzer | Done |
| Work as a standalone binary | Done |
| Work through `golangci-lint` module plugin flow | Done |
| Skip dynamic and non-literal messages to avoid false positives | Done |
| Emit precise diagnostics with file position and rule meaning | Done |
| Cover behavior with `analysistest` and fixture-based tests | Done |
| Provide README build, usage, and integration documentation | Done |
| Bonus: rule configuration via `.golangci.yml` | Done |
| Bonus: safe `SuggestedFixes` for auto-correctable cases | Done |
| Bonus: CI/CD pipeline for build, test, and release automation | Done |

## Example

Input:

```go
slog.Info("Starting auth token rotation")
```

Reported diagnostics:

```text
internal/service/auth.go:42:18: log message must start with a lowercase letter
internal/service/auth.go:42:18: log message may contain sensitive data
```

Safe autofix example:

```go
// before
slog.Info("Starting server")

// after SuggestedFix
slog.Info("starting server")
```

## Supported APIs

- `log/slog` top-level calls: `Debug`, `Info`, `Warn`, `Error`
- `log/slog` context variants: `DebugContext`, `InfoContext`, `WarnContext`, `ErrorContext`
- `log/slog` structured variants: `Log`, `LogAttrs`
- `*slog.Logger` methods, including chained `With(...)` and `WithGroup(...)`
- `*zap.Logger` methods: `Debug`, `Info`, `Warn`, `Error`
- `*zap.SugaredLogger` methods with explicit message argument: `Debugw`, `Infow`, `Warnw`, `Errorw`

## Requirements

- Go 1.24+

## Configuration

### Standalone flags

The standalone analyzer exposes rule toggles and custom sensitive keywords:

```bash
./bin/logslinter \
  -lowercase-start=true \
  -english-ascii-only=true \
  -no-special-chars-or-emoji=true \
  -no-sensitive-data=true \
  -additional-sensitive-keywords=credential,session_id \
  ./...
```

### golangci-lint module plugin

The repository includes example module-plugin configuration:

- `.custom-gcl.yml` builds a custom `golangci-lint` binary with `logsLinter` linked in.
- `.golangci.yml` enables the custom linter and shows analyzer settings.

Build the custom binary:

```bash
golangci-lint custom
```

With the example config in this repository, the binary is written to `./bin/custom-golangci-lint`.

Run it:

```bash
./bin/custom-golangci-lint run ./...
```

Minimal local-path plugin config:

```yaml
version: v2.8.0
name: custom-golangci-lint
destination: ./bin

plugins:
  - module: github.com/rTexty/logsLinter
    import: github.com/rTexty/logsLinter/plugin
    path: .
```

Analyzer settings example:

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
        settings:
          rules:
            lowercase-start: true
            english-ascii-only: true
            no-special-chars-or-emoji: true
            no-sensitive-data: true
          sensitive-data:
            additional-keywords:
              - credential
              - session_id
```

Available plugin settings:

- `rules.lowercase-start`
- `rules.english-ascii-only`
- `rules.no-special-chars-or-emoji`
- `rules.no-sensitive-data`
- `sensitive-data.additional-keywords`

## Development

Common commands:

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

## Verification

- Unit coverage exists for rule evaluation, extraction, diagnostics, and logger call inspection.
- Integration coverage runs through `analysistest` fixtures for `slog`, `zap`, and mixed edge cases.
- Configurable rule toggles and additional sensitive keywords are covered in analyzer and plugin tests.
- Current local verification baseline is `go test ./...`.

## Scope Boundaries

Skipped cases:

- Non-literal messages such as variables, function calls, and `fmt.Sprintf(...)`
- Literal plus variable concatenation such as `"password: " + secret`
- `zap.SugaredLogger` print-style methods such as `Info(...)`, `Warn(...)`, `Error(...)`
- `zap.SugaredLogger` format-style methods such as `Infof(...)`, `Warnf(...)`, `Errorf(...)`

Known limitations:

- Only string literals and literal-only concatenations are analyzed.
- Dynamic messages such as variables, `fmt.Sprintf(...)`, and mixed literal-plus-variable expressions are intentionally skipped.
- `zap.SugaredLogger` print-style and format-style methods stay out of scope in the initial release.
- The lowercase-start `SuggestedFix` is intentionally limited to interpreted string literals with an ASCII uppercase first rune.

## Troubleshooting

If `golangci-lint custom` fails, make sure the version in `.custom-gcl.yml` matches the version used in CI and your local build flow. This repository currently pins the custom binary workflow to `v2.8.0` for compatibility with the current release setup.

If the standalone analyzer returns a non-zero exit code, that is expected when diagnostics are found. It does not automatically mean package loading failed.

## Release

Recommended release checklist:

- `gofmt -l .`
- `go test ./... -race -count=1`
- `go build ./...`
- `go build -o ./bin/logslinter ./cmd/logslinter`
- `golangci-lint custom -v`
- `./bin/custom-golangci-lint run ./...`
- Verify standalone output on a sample package or repository

Create and push the first release tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow will:

- run tests before packaging
- build `logslinter` for Linux, macOS, and Windows
- upload `.tar.gz` and `.zip` artifacts
- publish a GitHub Release with generated notes and `SHA256SUMS.txt`

Release notes and checklist policy are documented in `CHANGELOG.md`, `docs/release-policy.md`, and `.github/release.yml`.

## Repository Automation

- GitHub Actions CI runs formatting checks, `go vet`, tests, build, and repository linting.
- CodeQL runs on pull requests and on a weekly schedule.
- Dependency review runs on pull requests.
- Dependabot tracks both Go modules and GitHub Actions.
- Releases are built from Git tags matching `v*` and publish packaged binaries plus `SHA256SUMS.txt`.
- Recommended GitHub branch ruleset settings are documented in `docs/github-ruleset.md`.
