# Release Policy

## Versioning

The repository uses Semantic Versioning tags:

- `MAJOR` for breaking changes
- `MINOR` for backward-compatible features
- `PATCH` for backward-compatible fixes

Release tags must use the format `vX.Y.Z`.

## Changelog Policy

- Keep pending changes under `Unreleased` in `CHANGELOG.md`.
- Move relevant entries from `Unreleased` into the released version section when creating a tag.
- Prefer concise user-facing entries over internal implementation detail.
- Security-sensitive details should only be disclosed when safe to publish.

## Release Notes Policy

- GitHub Release notes are generated automatically from merged pull requests and labels.
- Use labels such as `feature`, `enhancement`, `bug`, `security`, `documentation`, `chore`, and `breaking-change` to place changes in the correct section.
- Use `skip-release-notes` for internal changes that should not appear in public notes.

## Recommended Release Checklist

1. Ensure CI is green on the release commit.
2. Update `CHANGELOG.md`.
3. Confirm release labels are accurate on merged pull requests.
4. Create and push a `vX.Y.Z` tag.
5. Verify artifacts and `SHA256SUMS.txt` in the published GitHub Release.
