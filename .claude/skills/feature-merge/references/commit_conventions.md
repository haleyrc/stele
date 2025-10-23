# Commit Conventions

This project follows [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages.

## Commit Message Format

```
<type>: <description>

[optional body]

[optional footer(s)]
```

## Commit Types

- **feat**: A new feature for the user
- **fix**: A bug fix
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **docs**: Documentation only changes
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (updating dependencies, build scripts, etc.)

## Examples

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

## Changelog Generation

- Commits with `feat:` and `fix:` prefixes will appear in release changelogs
- Commits with `docs:`, `test:`, and `chore:` prefixes are excluded from changelogs
- Commits with `refactor:`, `perf:`, and `style:` are grouped under "Other Changes"
