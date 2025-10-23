#!/usr/bin/env python3
"""
Version analysis script for stele releases.

Analyzes commits since the last tag to suggest the next semantic version.
"""

import subprocess
import sys
import re
from typing import Tuple, List


def run_git_command(args: List[str]) -> str:
    """Run a git command and return its output."""
    try:
        result = subprocess.run(
            ["git"] + args,
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"Error running git command: {e}", file=sys.stderr)
        sys.exit(1)


def get_latest_tag() -> str:
    """Get the latest semantic version tag."""
    tags = run_git_command(["tag", "-l", "v*", "--sort=-version:refname"])
    if not tags:
        return "v0.0.0"
    return tags.split("\n")[0]


def parse_version(tag: str) -> Tuple[int, int, int, str]:
    """Parse a version tag into (major, minor, patch, prerelease)."""
    # Remove 'v' prefix
    version = tag[1:] if tag.startswith("v") else tag

    # Split prerelease suffix if present
    if "-" in version:
        base_version, prerelease = version.split("-", 1)
    else:
        base_version, prerelease = version, ""

    # Parse major.minor.patch
    parts = base_version.split(".")
    major = int(parts[0]) if len(parts) > 0 else 0
    minor = int(parts[1]) if len(parts) > 1 else 0
    patch = int(parts[2]) if len(parts) > 2 else 0

    return major, minor, patch, prerelease


def get_commits_since_tag(tag: str) -> List[str]:
    """Get commit messages since the specified tag."""
    commit_range = f"{tag}..HEAD" if tag != "v0.0.0" else "HEAD"
    commits = run_git_command(["log", "--pretty=format:%s", commit_range])
    return commits.split("\n") if commits else []


def analyze_commits(commits: List[str]) -> Tuple[bool, bool, bool]:
    """
    Analyze commits to determine version bump type.

    Returns: (has_breaking, has_features, has_fixes)
    """
    has_breaking = False
    has_features = False
    has_fixes = False

    for commit in commits:
        # Check for breaking changes
        if "BREAKING CHANGE" in commit or commit.startswith("!"):
            has_breaking = True

        # Check for features
        if commit.startswith("feat:") or commit.startswith("feat("):
            has_features = True

        # Check for fixes
        if commit.startswith("fix:") or commit.startswith("fix("):
            has_fixes = True

    return has_breaking, has_features, has_fixes


def suggest_next_version(current: str, commits: List[str]) -> dict:
    """Suggest the next version based on commits."""
    major, minor, patch, prerelease = parse_version(current)
    has_breaking, has_features, has_fixes = analyze_commits(commits)

    # Determine version bump
    if has_breaking:
        bump_type = "major"
        next_major, next_minor, next_patch = major + 1, 0, 0
    elif has_features:
        bump_type = "minor"
        next_major, next_minor, next_patch = major, minor + 1, 0
    elif has_fixes:
        bump_type = "patch"
        next_major, next_minor, next_patch = major, minor, patch + 1
    else:
        bump_type = "none"
        next_major, next_minor, next_patch = major, minor, patch

    return {
        "current_version": current,
        "current_parsed": f"{major}.{minor}.{patch}",
        "prerelease_suffix": prerelease,
        "bump_type": bump_type,
        "next_version": f"v{next_major}.{next_minor}.{next_patch}",
        "has_breaking": has_breaking,
        "has_features": has_features,
        "has_fixes": has_fixes,
        "commit_count": len([c for c in commits if c.strip()])
    }


def format_output(analysis: dict) -> str:
    """Format the analysis results."""
    output = []
    output.append(f"Current version: {analysis['current_version']}")
    output.append(f"Commits analyzed: {analysis['commit_count']}")
    output.append("")
    output.append("Commit analysis:")
    output.append(f"  - Breaking changes: {'Yes' if analysis['has_breaking'] else 'No'}")
    output.append(f"  - New features: {'Yes' if analysis['has_features'] else 'No'}")
    output.append(f"  - Bug fixes: {'Yes' if analysis['has_fixes'] else 'No'}")
    output.append("")

    if analysis['bump_type'] == 'none':
        output.append("Recommendation: No version bump needed")
    else:
        output.append(f"Recommendation: {analysis['bump_type'].upper()} version bump")
        output.append(f"Suggested version: {analysis['next_version']}")

        if analysis['prerelease_suffix']:
            beta_version = f"{analysis['next_version']}-{analysis['prerelease_suffix']}"
            output.append(f"Or for pre-release: {beta_version}")

    return "\n".join(output)


def main():
    """Main entry point."""
    print("Analyzing stele repository for version bump...\n")

    # Get latest tag
    latest_tag = get_latest_tag()

    # Get commits since tag
    commits = get_commits_since_tag(latest_tag)

    # Analyze and suggest
    analysis = suggest_next_version(latest_tag, commits)

    # Print results
    print(format_output(analysis))


if __name__ == "__main__":
    main()
