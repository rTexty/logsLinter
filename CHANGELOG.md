# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog and the project uses Semantic Versioning.

## [Unreleased]

### Added

-

### Changed

-

### Fixed

-

### Security

-

## [0.1.0] - 2026-03-10

### Added

- Initial public release of `logsLinter`, a Go static analyzer for validating literal log messages in `log/slog` and `go.uber.org/zap` codebases.
- Support for `log/slog` top-level APIs, context variants, structured variants, and `*slog.Logger` chains including `With(...)` and `WithGroup(...)`.
- Support for `*zap.Logger` message methods and `*zap.SugaredLogger` key-value methods with explicit message arguments.
- Four core validation rules: lowercase-start, English ASCII only, no decorative special characters or emoji, and no sensitive-data keywords.
- Standalone analyzer binary powered by `golang.org/x/tools/go/analysis`.
- `golangci-lint` module plugin integration with example configuration and runtime settings.
- Runtime configuration for enabling or disabling individual rules and extending sensitive keyword detection.
- Safe `SuggestedFix` support for lowercase-start violations on compatible string literals.
- Unit tests, `analysistest` integration coverage, release automation, repository CI, CodeQL, and dependency management workflows.

### Changed

- Standardized repository automation, release packaging, and checksum generation for multi-platform builds.
- Organized project documentation around standalone usage, plugin integration, verification, and release policy.

### Fixed

- Aligned custom `golangci-lint` workflow versions with a compatible build path for CI.
- Resolved workflow lint issues in the release pipeline shell scripts.

### Security

- Added sensitive-data detection defaults to reduce the risk of shipping obvious credential-like log messages.
