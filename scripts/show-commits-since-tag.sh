#!/usr/bin/env bash
set -euo pipefail

# Allow specifying a tag as first argument, otherwise use latest
if [[ $# -gt 0 ]]; then
  LATEST_TAG="$1"
  # Verify the tag exists
  if ! git rev-parse "$LATEST_TAG" >/dev/null 2>&1; then
    echo "Error: Tag '$LATEST_TAG' does not exist"
    exit 1
  fi
else
  # Get the latest tag
  LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
fi

if [[ -z "$LATEST_TAG" ]]; then
  echo "No tags found in repository. Showing all commits:"
  git log --pretty=format:"%h - %s" HEAD
else
  echo "Commits since tag $LATEST_TAG:"
  echo "----------------------------------------"
  git log "${LATEST_TAG}..HEAD" --pretty=format:"%s"
fi

echo ""  # Add newline at end
