# Homebrew Distribution Setup Guide

This guide walks you through setting up Homebrew distribution for `stele`, which provides the best installation experience for macOS users and eliminates the Gatekeeper "unverified developer" warning without requiring an Apple Developer account.

## Table of Contents

1. [Why Homebrew?](#why-homebrew)
2. [Prerequisites](#prerequisites)
3. [Overview of the Process](#overview-of-the-process)
4. [Step 1: Create a Homebrew Tap Repository](#step-1-create-a-homebrew-tap-repository)
5. [Step 2: Create the Formula](#step-2-create-the-formula)
6. [Step 3: Configure GoReleaser](#step-3-configure-goreleaser)
7. [Step 4: Test Locally](#step-4-test-locally)
8. [Step 5: Release and Verify](#step-5-release-and-verify)
9. [User Installation Instructions](#user-installation-instructions)
10. [Maintenance](#maintenance)
11. [Troubleshooting](#troubleshooting)

## Why Homebrew?

Homebrew distribution solves several problems:

- **No Gatekeeper Warning**: Binaries installed via Homebrew don't trigger macOS Gatekeeper warnings
- **No Apple Developer Account Required**: Unlike code signing/notarization, Homebrew is completely free
- **Better UX**: Users get automatic updates via `brew upgrade`
- **Standard Practice**: Most popular Go CLI tools (hugo, gh, jq, etc.) use Homebrew
- **Cross-platform**: While this guide focuses on macOS, Homebrew also works on Linux

## Prerequisites

Before you begin, ensure you have:

- [x] A GitHub account
- [x] Existing releases of `stele` on GitHub with binary assets (you already have this!)
- [x] GoReleaser installed locally for testing (optional but recommended)
- [x] Homebrew installed on your local machine for testing

## Overview of the Process

Here's what we'll do:

1. Create a new GitHub repository called `homebrew-tap` (or similar name)
2. Create a formula file that tells Homebrew how to install `stele`
3. Configure GoReleaser to automatically update this formula on each release
4. Test the formula locally
5. Users can then install with `brew install haleyrc/tap/stele`

## Step 1: Create a Homebrew Tap Repository

A "tap" is a GitHub repository that contains Homebrew formulas.

### 1.1 Create the Repository

Create a new GitHub repository with one of these naming patterns:
- `homebrew-tap` (recommended - simple and clear)
- `homebrew-stele` (if you want a dedicated tap)

**Important naming rules:**
- Must start with `homebrew-`
- The part after `homebrew-` becomes the tap name
- Example: `homebrew-tap` → users reference it as `haleyrc/tap`

### 1.2 Initialize the Repository

```bash
# Create and navigate to your tap directory
mkdir homebrew-tap
cd homebrew-tap

# Initialize git
git init
git branch -M main

# Create the Formula directory (required)
mkdir Formula

# Create a basic README
cat > README.md << 'EOF'
# Homebrew Tap for haleyrc

This tap provides Homebrew formulas for haleyrc projects.

## Installation

\`\`\`bash
brew install haleyrc/tap/stele
\`\`\`

## Available Formulas

- **stele**: An opinionated static site generator
EOF

# Create .gitignore
cat > .gitignore << 'EOF'
.DS_Store
EOF

# Commit and push
git add .
git commit -m "Initial commit"
git remote add origin git@github.com:haleyrc/homebrew-tap.git
git push -u origin main
```

## Step 2: Create the Formula

A formula is a Ruby file that tells Homebrew how to install your software.

### 2.1 Option A: Generate Formula Automatically (Recommended)

You can use Homebrew to generate an initial formula:

```bash
# This will create a template based on your latest release
brew create --tap haleyrc/tap https://github.com/haleyrc/stele/releases/download/v1.0.0-beta.3/stele_Darwin_arm64.tar.gz
```

This creates `Formula/stele.rb` in your tap repository.

### 2.2 Option B: Create Formula Manually

Create `Formula/stele.rb` with this content:

```ruby
class Stele < Formula
  desc "An opinionated static site generator for people with analysis paralysis"
  homepage "https://github.com/haleyrc/stele"
  version "1.0.0-beta.3"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/haleyrc/stele/releases/download/v1.0.0-beta.3/stele_Darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/haleyrc/stele/releases/download/v1.0.0-beta.3/stele_Linux_arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
    if Hardware::CPU.intel?
      url "https://github.com/haleyrc/stele/releases/download/v1.0.0-beta.3/stele_Linux_x86_64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
  end

  def install
    bin.install "stele"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/stele --version")
  end
end
```

### 2.3 Calculate SHA256 Checksums

For each binary, calculate its SHA256:

```bash
# Download the release asset
curl -L -o stele_Darwin_arm64.tar.gz \
  https://github.com/haleyrc/stele/releases/download/v1.0.0-beta.3/stele_Darwin_arm64.tar.gz

# Calculate SHA256
shasum -a 256 stele_Darwin_arm64.tar.gz
```

Replace `REPLACE_WITH_ACTUAL_SHA256` in the formula with the actual checksums.

**Note:** GoReleaser already generates a checksums file! Check your release assets for `stele_1.0.0-beta.3_checksums.txt`.

### 2.4 Commit the Formula

```bash
cd homebrew-tap
git add Formula/stele.rb
git commit -m "Add stele formula v1.0.0-beta.3"
git push
```

## Step 3: Configure GoReleaser

Now let's automate formula updates on each release by configuring GoReleaser.

### 3.1 Create a GitHub Personal Access Token

GoReleaser needs permission to push to your tap repository:

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token with these scopes:
   - `repo` (full control)
3. Copy the token value

### 3.2 Add Token to stele Repository

Add the token as a repository secret:

1. Go to your `stele` repository settings
2. Navigate to Secrets and variables → Actions
3. Create a new secret:
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Value: (paste your token)

### 3.3 Update .goreleaser.yaml

Add this `brews` section to your `.goreleaser.yaml`:

```yaml
brews:
  - name: stele
    repository:
      owner: haleyrc
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"

    directory: Formula

    homepage: "https://github.com/haleyrc/stele"
    description: "An opinionated static site generator for people with analysis paralysis"
    license: "MIT"

    # Only build for macOS in the formula (matching your current builds config)
    install: |
      bin.install "stele"

    test: |
      assert_match version.to_s, shell_output("#{bin}/stele --version")

    # This tells Homebrew it depends on Go's runtime (optional, since binary is statically linked)
    # dependencies:
    #   - name: go
    #     type: optional
```

### 3.4 Update Your Release Process

When you create releases, make sure the environment variable is available:

**If releasing via GitHub Actions**, update your workflow to include:

```yaml
- name: Run GoReleaser
  uses: goreleaser/goreleaser-action@v5
  with:
    version: latest
    args: release --clean
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

**If releasing locally**, export the variable:

```bash
export HOMEBREW_TAP_GITHUB_TOKEN="your_token_here"
goreleaser release --clean
```

## Step 4: Test Locally

Before releasing, test your formula locally:

### 4.1 Install from Your Tap

```bash
# Add your tap
brew tap haleyrc/tap

# Install stele
brew install haleyrc/tap/stele

# Test it works
stele --version
```

### 4.2 Test the Formula Directly

You can audit the formula for common issues:

```bash
brew audit --strict haleyrc/tap/stele
brew test haleyrc/tap/stele
```

### 4.3 Uninstall for Clean Testing

```bash
brew uninstall stele
brew untap haleyrc/tap
```

## Step 5: Release and Verify

### 5.1 Create a New Release

Create a new tag and release as you normally would:

```bash
git tag v1.0.0-beta.4
git push origin v1.0.0-beta.4
goreleaser release --clean
```

### 5.2 Verify Formula Update

After the release completes:

1. Check your `homebrew-tap` repository
2. You should see a new commit from GoReleaser updating `Formula/stele.rb`
3. The commit message will be something like "Brew formula update for stele version v1.0.0-beta.4"

### 5.3 Test Installation

```bash
# Update your local tap
brew update

# Install the new version
brew upgrade stele
# or if not installed yet:
brew install haleyrc/tap/stele

# Verify version
stele --version
```

## User Installation Instructions

Once everything is set up, users can install `stele` with:

```bash
# Add the tap and install in one command
brew install haleyrc/tap/stele

# Or add tap first, then install
brew tap haleyrc/tap
brew install stele
```

### Upgrading

```bash
brew upgrade stele
```

### Uninstalling

```bash
brew uninstall stele
```

## Maintenance

### Updating the Formula

With GoReleaser configured, formulas are automatically updated on each release. However, you may occasionally need to manually update:

```bash
cd homebrew-tap
# Edit Formula/stele.rb as needed
git add Formula/stele.rb
git commit -m "Update stele formula: [reason]"
git push
```

### Common Manual Updates

You might manually update for:
- Changing the description
- Adding dependencies
- Modifying the test block
- Adding caveats for users

### Example: Adding User Caveats

```ruby
class Stele < Formula
  # ... other config ...

  def caveats
    <<~EOS
      To get started with stele:
        stele dev

      For more information:
        https://github.com/haleyrc/stele
    EOS
  end
end
```

## Troubleshooting

### Formula Fails to Install

**Problem**: `brew install haleyrc/tap/stele` fails with checksum mismatch

**Solution**:
- The SHA256 in the formula doesn't match the actual file
- Download the release asset and recalculate: `shasum -a 256 file.tar.gz`
- Update the formula with correct SHA256

### GoReleaser Doesn't Update Formula

**Problem**: New releases don't update the homebrew-tap repository

**Solution**:
- Verify `HOMEBREW_TAP_GITHUB_TOKEN` is set correctly
- Check token has `repo` scope permissions
- Look at GoReleaser output for errors
- Ensure the `brews` section in `.goreleaser.yaml` is properly formatted

### Formula Audit Warnings

**Problem**: `brew audit` reports warnings

**Solution**:
- Review the warning message
- Common issues:
  - Missing license field
  - URL not using HTTPS
  - Incorrect homepage URL
  - Test block doesn't properly verify installation

### Users Can't Find the Tap

**Problem**: `brew tap haleyrc/tap` fails

**Solution**:
- Verify repository name follows `homebrew-*` pattern
- Ensure repository is public
- Check repository exists at `github.com/haleyrc/homebrew-tap`

### Architecture Support

**Problem**: Formula only works on ARM64 Macs

**Solution**:
- Update `.goreleaser.yaml` to build for more architectures:

```yaml
builds:
  - main: .
    env:
      - CGO_ENABLED=1
    goos:
      - darwin
      - linux
    goarch:
      - amd64
      - arm64
```

Then update the formula to include all architectures (or let GoReleaser generate it).

## Additional Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GoReleaser Homebrew Documentation](https://goreleaser.com/customization/homebrew/)
- [Homebrew Tap Documentation](https://docs.brew.sh/Taps)

## Next Steps

After Homebrew distribution is set up, consider:

1. **Update README.md**: Make Homebrew the primary installation method
2. **Remove manual download warnings**: Users won't see Gatekeeper warnings via Homebrew
3. **Add installation badge**: Show Homebrew install command prominently
4. **Consider Homebrew Core**: Once stable, you could submit to the main Homebrew repository (though your tap works great for most cases)

---

**Note**: This guide assumes you're maintaining `stele` as an open-source project. All tools and services mentioned (GitHub, Homebrew, GoReleaser) are free for open-source use.
