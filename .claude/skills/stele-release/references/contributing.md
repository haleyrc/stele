# Contributing to Stele

Thank you for your interest in contributing to Stele! This document provides guidelines for contributing to the project.

## Commit Message Conventions

We follow [Conventional Commits](https://www.conventionalcommits.org/) to maintain a clean and meaningful commit history. This also enables automatic changelog generation.

### Commit Message Format

```
<type>: <description>

[optional body]

[optional footer(s)]
```

### Commit Types

- **feat**: A new feature for the user
- **fix**: A bug fix
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **docs**: Documentation only changes
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (updating dependencies, build scripts, etc.)

### Examples

```bash
# New feature
git commit -m "feat: Add support for custom templates"

# Bug fix
git commit -m "fix: Resolve navigation race condition"

# Documentation
git commit -m "docs: Update installation instructions"

# Performance improvement
git commit -m "perf: Optimize template compilation"
```

### Notes on Changelog Generation

- Commits with `feat:` and `fix:` prefixes will appear in release changelogs
- Commits with `docs:`, `test:`, and `chore:` prefixes are excluded from changelogs
- Commits with `refactor:`, `perf:`, and `style:` are grouped under "Other Changes"

## Release Process

Stele uses automated releases via GitHub Actions and GoReleaser. Releases are triggered by pushing semantic version tags.

### Creating a Release

1. **Ensure all tests pass**
   ```bash
   go test ./...
   ./bin/check
   ```

2. **Review recent commits and plan the version**
   ```bash
   git log --oneline
   ```

   Follow [Semantic Versioning](https://semver.org/):
   - **MAJOR** (v2.0.0): Breaking changes
   - **MINOR** (v1.1.0): New features, backwards compatible
   - **PATCH** (v1.0.1): Bug fixes, backwards compatible
   - **PRE-RELEASE** (v1.0.0-beta.1): Pre-release versions

3. **Create an annotated tag**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   ```

4. **Push the tag to trigger the release workflow**
   ```bash
   git push origin v1.0.0
   ```

5. **Monitor the GitHub Actions workflow**
   - Visit the Actions tab in the GitHub repository
   - Watch the "Release" workflow complete
   - Verify tests pass and artifacts are generated

6. **Verify the release**
   - Check the Releases page for the new release
   - Download and test the binary
   - Review the generated changelog

### Pre-Release Versions

For beta, alpha, or release candidate versions:

```bash
# Beta release
git tag -a v1.0.0-beta.1 -m "Release v1.0.0-beta.1"
git push origin v1.0.0-beta.1

# Alpha release
git tag -a v2.0.0-alpha.1 -m "Release v2.0.0-alpha.1"
git push origin v2.0.0-alpha.1

# Release candidate
git tag -a v1.0.0-rc.1 -m "Release v1.0.0-rc.1"
git push origin v1.0.0-rc.1
```

These will automatically be marked as pre-releases on GitHub.

### Rollback Procedures

#### If the Release Workflow Fails

1. Delete the remote tag:
   ```bash
   git push --delete origin v1.0.0
   ```

2. Delete the local tag:
   ```bash
   git tag -d v1.0.0
   ```

3. Fix the issues in the codebase

4. Re-create and push the tag once fixed

#### If the Release Succeeds but Has Issues

**Do not delete the release** - users may have already downloaded it.

1. Fix the issue in the codebase
2. Create a new patch release (e.g., v1.0.1)
3. Optionally edit the old release notes to point to the fixed version

## Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Run tests and checks (`go test ./...` and `./bin/check`)
5. Commit with conventional commit messages
6. Push to your fork
7. Create a pull request

## Questions?

If you have questions about contributing, please open an issue on GitHub.
