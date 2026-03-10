# GitHub Ruleset Recommendations

This repository contains the automation and policy files that can live in Git.
GitHub repository rulesets themselves are configured in the GitHub UI and are not reliably portable as repository files.

## Recommended Branch Ruleset

Target branch pattern:

- `main`

Recommended protections:

- require pull requests before merge
- require at least 1 approval
- dismiss stale approvals on new commits
- require conversation resolution before merge
- require status checks before merge
- block force pushes
- block branch deletion
- require linear history

## Recommended Required Status Checks

- `build-test`
- `dependency-review`
- `analyze`
- `actionlint`

## Recommended Repository Settings

- enable vulnerability alerts
- enable Dependabot security updates
- enable automatic deletion of merged branches
- enable private vulnerability reporting
