# Contributing to Stele

Thank you for your interest in contributing to Stele! This document outlines the development workflow and conventions used in this project.

## Development Setup

### Prerequisites
- Go 1.25 or later
- [templ](https://templ.guide/) for template generation

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/haleyrc/stele.git
   cd stele
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run checks and tests:
   ```bash
   ./bin/check
   go test ./...
   ```

4. Start the development server:
   ```bash
   go run . dev
   ```

## Commit Message Conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for clear and structured commit history.

### Format

```
<type>: <description>

[optional body]
```

### Types

- **feat**: New feature or functionality
  ```
  feat: Add support for code syntax highlighting
  ```

- **fix**: Bug fix
  ```
  fix: Resolve race condition in live reload server
  ```

- **docs**: Documentation changes only (excluded from changelog)
  ```
  docs: Update installation instructions
  ```

- **test**: Test additions or modifications (excluded from changelog)
  ```
  test: Add integration tests for compiler
  ```

- **chore**: Maintenance tasks, dependency updates (excluded from changelog)
  ```
  chore: Update templ to v0.3.943
  ```

- **refactor**: Code restructuring without changing behavior
  ```
  refactor: Simplify template rendering logic
  ```

- **perf**: Performance improvements
  ```
  perf: Optimize markdown parsing
  ```

- **style**: Code formatting, whitespace changes
  ```
  style: Run go fmt on all files
  ```

### Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Keep the first line under 72 characters
- Reference issues and PRs when applicable

### Examples

Good commits:
```
feat: Add frontmatter validation for posts
fix: Handle empty tag arrays in post metadata
docs: Add deployment guide for Cloudflare Pages
chore: Update dependencies to latest versions
```

## Testing

Before submitting changes:

1. Run the validation script:
   ```bash
   ./bin/check
   ```
   This runs:
   - `go fmt` - Code formatting
   - `go vet` - Static analysis
   - `staticcheck` - Additional linting
   - `gosec` - Security checks

2. Run tests:
   ```bash
   go test ./...
   ```

3. Test manually with the dev server:
   ```bash
   go run . dev
   ```

## Release Process

Releases are automated using GitHub Actions and GoReleaser. See [docs/RELEASE_PLAN.md](docs/RELEASE_PLAN.md) for full details.

### For Maintainers

1. Ensure all tests pass and code is ready
2. Create an annotated tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   ```
3. Push the tag:
   ```bash
   git push origin v1.0.0
   ```
4. GitHub Actions will automatically:
   - Run tests and validation
   - Build the binary
   - Generate changelog from commits
   - Create GitHub release with artifacts

### Beta Releases

Use beta tags for pre-release versions:
```bash
git tag -a v1.0.0-beta.4 -m "Beta release v1.0.0-beta.4"
git push origin v1.0.0-beta.4
```

These are automatically marked as pre-releases on GitHub.

## Code Style

- Follow standard Go conventions
- Use `go fmt` for formatting
- Write clear, self-documenting code
- Add comments for non-obvious logic
- Keep functions focused and small

## Questions?

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Be respectful and constructive in discussions

## License

By contributing to Stele, you agree that your contributions will be licensed under the project's license.
