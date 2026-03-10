# logsLinter

[![CI](https://github.com/rTexty/logsLinter/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/rTexty/logsLinter/actions/workflows/ci.yml)
[![CodeQL](https://github.com/rTexty/logsLinter/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/rTexty/logsLinter/actions/workflows/codeql.yml)
[![Release](https://github.com/rTexty/logsLinter/actions/workflows/release.yml/badge.svg)](https://github.com/rTexty/logsLinter/actions/workflows/release.yml)
[![Latest Release](https://img.shields.io/github/v/release/rTexty/logsLinter)](https://github.com/rTexty/logsLinter/releases)

Production-ready Go analyzer for validating log messages in `log/slog` and `go.uber.org/zap` codebases.

## Status Panel

| Signal | Status |
| --- | --- |
| CI | formatting, vet, tests, build, repository lint |
| Security | CodeQL, dependency review, Dependabot, SECURITY policy |
| Releases | tag-based GitHub Release workflow with packaged binaries and checksums |
| Tooling | standalone CLI, golangci-lint module plugin examples |

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

The repository is currently being scaffolded. The standalone analyzer entrypoint, analyzer package, test fixtures, and `golangci-lint` integration examples are implemented incrementally from the feature plan in `Context/Features/001-LogsLinter/Steps.md`.
