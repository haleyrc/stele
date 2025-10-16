# Release Plan: GitHub Actions Automation for Stele CLI

## Overview

This document outlines a comprehensive plan for automating releases of the `stele` CLI using GitHub Actions and GoReleaser. The plan builds on the existing GoReleaser configuration and testing infrastructure already in place.

## Current State

### Existing Infrastructure
- **GoReleaser**: Configured in `.goreleaser.yaml` with:
  - Darwin/ARM64 (macOS Apple Silicon) builds
  - CGO enabled
  - UPX compression
  - Version injection via ldflags
  - Pre-build validation hook (`./bin/check`)
  - Changelog generation (filters docs, test, chore commits)

- **GitHub Actions**: Single test workflow (`.github/workflows/test.yml`)
  - Uses reusable workflow from `haleyrc/actions/.github/workflows/go-test.yml@main`
  - Runs on all pushes

- **Versioning**: Following semver with beta releases (v1.0.0-beta.3 is latest)

### Gaps
- No automated release workflow
- Manual tag creation and release publishing

## Proposed Solution

### 1. Release Workflow Architecture

#### Trigger Strategy
Release on **semantic version tags** (e.g., `v1.0.0`, `v1.0.0-beta.4`):
- Tags matching `v*` pattern trigger release workflow
- Beta/pre-release tags (`v*-beta*`, `v*-alpha*`, etc.) create pre-releases
- Stable tags (`v1.0.0`) create production releases

#### Workflow File: `.github/workflows/release.yml`

**Key components:**
1. **Trigger**: On tag push matching `v*`
2. **Permissions**: Write access for contents (to create releases and upload artifacts)
3. **Jobs**:
   - `release`: Run GoReleaser to build and publish

### 2. GoReleaser Configuration

The existing `.goreleaser.yaml` configuration is already set up correctly for darwin/arm64 builds with:
- CGO enabled
- UPX compression for smaller binaries
- Version injection via ldflags
- Pre-build validation hook

No changes are required to the GoReleaser configuration for basic release automation.

### 3. Release Workflow Implementation

#### Basic Release Workflow
```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for changelog

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'

      - name: Run pre-release checks
        run: ./bin/check

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Pre-Release Detection
Automatically mark beta/alpha releases as pre-releases:
```yaml
      - name: Determine release type
        id: release_type
        run: |
          if [[ $GITHUB_REF_NAME =~ (alpha|beta|rc) ]]; then
            echo "prerelease=true" >> $GITHUB_OUTPUT
          else
            echo "prerelease=false" >> $GITHUB_OUTPUT
          fi
```

Update GoReleaser config to use this:
```yaml
release:
  prerelease: auto
  name_template: "{{.Version}}"
  github:
    owner: haleyrc
    name: stele
```

### 4. Release Process Workflow

#### Development to Release Flow
```
1. Development
   └─> Feature branches → main branch

2. Quality Assurance
   └─> All tests pass on main
   └─> Run bin/check locally
   └─> Review CHANGELOG

3. Version Tagging
   └─> Create annotated tag: git tag -a v1.0.0 -m "Release v1.0.0"
   └─> Push tag: git push origin v1.0.0

4. Automated Release (GitHub Actions)
   └─> Trigger release workflow
   └─> Run tests
   └─> Run bin/check
   └─> Build darwin/arm64 binary
   └─> Generate changelog
   └─> Create GitHub release
   └─> Upload artifacts

5. Post-Release
   └─> Verify release artifacts
   └─> Test installation: gh release download
   └─> Update documentation if needed
   └─> Announce release
```

#### Manual Steps
1. **Before tagging**: Update version-specific docs, test thoroughly
2. **Create tag**: `git tag -a v1.0.0 -m "Release v1.0.0"`
3. **Push tag**: `git push origin v1.0.0`
4. **Monitor**: Watch GitHub Actions workflow completion
5. **Verify**: Test downloaded binaries

### 5. Security Considerations

#### Secrets Management
Required secrets in GitHub repository settings:
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions (no setup required)

#### Dependency Management
- Use `go mod verify` in workflow
- Run security scanners (gosec already in tools)
- Keep GitHub Actions versions pinned with SHA

### 6. Changelog Management

#### Current State
`.goreleaser.yaml` has basic changelog config that filters out:
- `docs:` commits
- `test:` commits
- `chore:` commits

#### Recommended Enhancement

Add grouping to organize changelog entries by type:

```yaml
changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
  groups:
    - title: "Features"
      regexp: "^feat"
      order: 0
    - title: "Bug Fixes"
      regexp: "^fix"
      order: 1
    - title: "Other Changes"
      regexp: "^refactor|^perf|^style"
      order: 2
```

This produces organized release notes like:
```markdown
## Features
- feat: Add dark mode support
- feat: Include file paths in error messages

## Bug Fixes
- fix: Resolve navigation race condition
- fix: Update urlf format strings

## Other Changes
- refactor: Simplify compiler logic
```

**Implementation**: Update `.goreleaser.yaml` and document commit conventions in `CONTRIBUTING.md`.

### 7. Testing Strategy

#### Pre-Release Testing
Add to release workflow:
```yaml
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
      - name: Run tests
        run: go test ./...
      - name: Run checks
        run: ./bin/check

  release:
    needs: test  # Don't release if tests fail
    runs-on: ubuntu-latest
    # ... rest of release job
```

#### Post-Release Verification
Create a verification job:
```yaml
  verify:
    needs: release
    runs-on: macos-latest
    steps:
      - name: Download release
        run: |
          gh release download ${{ github.ref_name }} -R haleyrc/stele
      - name: Extract and test
        run: |
          tar -xzf stele_*.tar.gz
          ./stele version
```

### 8. Rollback Strategy

#### If Release Fails
1. **Delete tag**: `git push --delete origin v1.0.0`
2. **Delete local tag**: `git tag -d v1.0.0`
3. **Delete GitHub release**: Via UI or `gh release delete v1.0.0`
4. **Fix issues**: Address problems in codebase
5. **Re-tag**: Create new tag with same version

#### If Release Succeeds but Has Issues
1. **Don't delete release** (breaks users who already downloaded)
2. **Create hotfix**: Fix issue in code
3. **Release patch**: Tag and release v1.0.1
4. **Mark old release**: Edit release notes to indicate issue and point to fixed version

### 9. Monitoring and Notifications

#### Workflow Status
Monitor via:
- GitHub Actions UI
- Email notifications (configured per-user)
- Slack/Discord webhooks (optional)

#### Release Notifications
```yaml
  notify:
    needs: release
    runs-on: ubuntu-latest
    if: success()
    steps:
      - name: Send notification
        run: |
          # Post to Slack/Discord/etc
          curl -X POST ${{ secrets.WEBHOOK_URL }} \
            -H 'Content-Type: application/json' \
            -d '{"text":"Released stele ${{ github.ref_name }}"}'
```

### 10. Documentation Updates

Files to create/update:
- `.github/workflows/release.yml` - Main release workflow
- `.goreleaser.yaml` - Add changelog grouping configuration
- `CONTRIBUTING.md` - Commit conventions and release process (created)
- `docs/RELEASE_CHECKLIST.md` - Create manual checklist

### 11. Implementation Path

#### Phase 1: Basic Automation
1. Create `.github/workflows/release.yml`
2. Test with next beta release (v1.0.0-beta.4)
3. Verify artifacts are generated correctly

#### Phase 2: Documentation
1. Create `docs/RELEASE_CHECKLIST.md`
2. Update `CONTRIBUTING.md` with release process
3. Document rollback procedures

### 12. Success Metrics

Track the following to measure success:
- **Time to release**: Should drop from manual (~30 min) to automated (~5 min)
- **Release frequency**: Should increase with reduced friction
- **Download stats**: Monitor via GitHub release metrics
- **Error rate**: Workflow failures should be < 5%

### 13. Potential Issues and Mitigations

| Issue | Risk | Mitigation |
|-------|------|------------|
| UPX compression failures | Low | Monitor logs, UPX can be disabled if needed |
| Large binary sizes | Low | Already using UPX and strip flags |
| Tag conflicts | Low | Use protected tags, documented versioning scheme |
| Workflow quota limits | Low | GitHub provides generous free tier, monitor usage |

## Implementation Checklist

### Prerequisites
- [ ] Review current `.goreleaser.yaml` configuration (already set up correctly)
- [ ] Ensure `./bin/check` script works correctly

### Phase 1: Basic Release Automation
- [ ] Create `.github/workflows/release.yml`
- [ ] Configure workflow permissions
- [ ] Add GoReleaser action with pre-release checks
- [ ] Update `.goreleaser.yaml` with changelog grouping
- [ ] Test with beta release tag
- [ ] Verify artifact generation and download

### Phase 2: Documentation
- [x] Create `CONTRIBUTING.md` with commit conventions and release process
- [ ] Create `docs/RELEASE_CHECKLIST.md`
- [ ] Document rollback procedures
- [ ] Add troubleshooting guide

## Maintenance

### Regular Tasks
- **Weekly**: Monitor workflow execution, review failures
- **Monthly**: Update GitHub Actions versions, review dependencies
- **Per Release**: Follow release checklist, verify artifacts
- **Quarterly**: Review and update documentation

### Action Version Updates
Pin action versions and update regularly:
```yaml
# Current versions as of Oct 2024
- uses: actions/checkout@v4
- uses: actions/setup-go@v5
- uses: goreleaser/goreleaser-action@v6
```

Check for updates monthly via Dependabot:
```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
```

## Conclusion

This release plan provides a streamlined approach to automating releases for the `stele` CLI. By implementing this plan, the project will benefit from:

- **Reduced release time**: From 30 minutes to 5 minutes
- **Increased reliability**: Automated testing and checks
- **Consistent process**: Documented workflows and checklists
- **Better tracking**: GitHub releases with automatic changelogs

The implementation is straightforward: create the GitHub Actions workflow and test with the next beta release. Documentation can be added as time allows.

## References

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
