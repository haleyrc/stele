# Release Checklist

This checklist ensures a smooth and reliable release process for Stele.

## Pre-Release Checklist

Before creating a release tag, verify the following:

### Code Quality
- [ ] Pre-release checks pass: `./bin/check` (includes tests)
- [ ] No pending critical bug fixes
- [ ] Code has been reviewed (if applicable)

### Documentation
- [ ] CHANGELOG reflects recent changes (automatic, but review for accuracy)
- [ ] README is up to date
- [ ] Version-specific documentation updated (if needed)
- [ ] Breaking changes are documented (for major versions)

### Version Planning
- [ ] Determined appropriate version number following [semver](https://semver.org/):
  - **PATCH** (v1.0.X): Bug fixes only
  - **MINOR** (v1.X.0): New features, backwards compatible
  - **MAJOR** (vX.0.0): Breaking changes
  - **PRE-RELEASE** (v1.0.0-beta.X): Beta/alpha/rc versions

### Repository State
- [ ] On the correct branch (usually `main`)
- [ ] Branch is up to date with remote: `git pull`
- [ ] Working directory is clean: `git status`
- [ ] All intended commits are merged

## Release Steps

### 1. Create the Release Tag

```bash
# For a stable release
git tag -a v1.0.0 -m "Release v1.0.0"

# For a pre-release (beta/alpha/rc)
git tag -a v1.0.0-beta.1 -m "Release v1.0.0-beta.1"
```

### 2. Push the Tag

```bash
git push origin v1.0.0
```

### 3. Monitor GitHub Actions

- Navigate to the [Actions tab](https://github.com/haleyrc/stele/actions)
- Find the "Release" workflow for your tag
- Watch the workflow progress:
  - ✅ Test job completes successfully
  - ✅ Release job runs GoReleaser
  - ✅ Artifacts are uploaded

### 4. Verify Release on GitHub

Once the workflow completes:

- [ ] Visit the [Releases page](https://github.com/haleyrc/stele/releases)
- [ ] Confirm the new release appears
- [ ] Verify release notes are generated correctly
- [ ] Check that changelog groups are properly formatted:
  - Features section
  - Bug Fixes section
  - Other Changes section
- [ ] Pre-release is marked correctly (if applicable)
- [ ] Assets are attached (darwin_arm64 tarball)

## Post-Release Verification

### 5. Download and Test the Binary

```bash
# Download the release
gh release download v1.0.0 -R haleyrc/stele

# Extract
tar -xzf stele_Darwin_arm64.tar.gz

# Test the binary
./stele version

# Verify version matches release
./stele --help
```

### 6. Functional Testing (Optional but Recommended)

Run through key workflows with the released binary:

- [ ] Basic command execution works
- [ ] Core features function correctly
- [ ] No obvious regressions

## Rollback Procedures

### If the Release Workflow Fails

The release was never created, so it's safe to delete and retry:

```bash
# Delete remote tag
git push --delete origin v1.0.0

# Delete local tag
git tag -d v1.0.0

# Fix the issue in the codebase
# ... make necessary changes ...

# Re-create and push the tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### If the Release Succeeds but Has Critical Issues

1. **Document the issue**
   - Edit the release notes on GitHub
   - Add a warning at the top about the issue
   - Point users to the upcoming fixed version

2. **Create a hotfix**
   ```bash
   # Fix the issue in your codebase
   git commit -m "fix: Critical issue in v1.0.0"

   # Create a patch release
   git tag -a v1.0.1 -m "Release v1.0.1 - Hotfix for v1.0.0"
   git push origin v1.0.1
   ```

3. **Update the old release notes**
   - Add a note at the top: "⚠️ This release has a known issue with X. Please use v1.0.1 instead."

### If the Release Has Minor Issues

For non-critical issues:

- [ ] Add known issues to release notes
- [ ] Create an issue to track the bug
- [ ] Fix in the next scheduled release

## Common Issues and Solutions

### Workflow Fails on Tests

**Cause**: Code doesn't pass tests or pre-release checks

**Solution**:
1. Delete the tag (see rollback procedure)
2. Fix failing tests
3. Re-tag and push

### Workflow Fails on GoReleaser

**Cause**: GoReleaser configuration issue or missing dependencies

**Solution**:
1. Check GoReleaser logs in GitHub Actions
2. Test locally: `goreleaser release --snapshot --clean`
3. Fix configuration issues
4. Delete tag, fix, and retry
