#!/usr/bin/env bash
set -euo pipefail

# Config — edit if needed
MAIN_PKG="./cmd/server"
BINARY_BASE="server"
BUILD_DIR="dist"
GOOS_LIST=("linux")
GOARCH_LIST=("amd64" "arm64")
VERSION_FILE="VERSION"
GIT_REMOTE="${GIT_REMOTE:-origin}"
GIT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"

# GitHub repo info (owner/repo). Try to infer from origin if not set.
GITHUB_REPO="${GITHUB_REPO:-}"
if [[ -z "$GITHUB_REPO" ]]; then
  ORIGIN_URL="$(git config --get remote.$GIT_REMOTE.url || true)"
  if [[ "$ORIGIN_URL" =~ github.com[:/](.+/.+)(\.git)?$ ]]; then
    GITHUB_REPO="${BASH_REMATCH[1]}"
  else
    echo "Cannot infer GitHub repo from origin. Set GITHUB_REPO env (owner/repo)."
    exit 1
  fi
fi

# Auth token
# Read GitHub token from file
GITHUB_TOKEN_FILE="${GITHUB_TOKEN_FILE:-$HOME/.secrets/billdb/token}"

if [[ ! -f "$GITHUB_TOKEN_FILE" ]]; then
  echo "GitHub token file not found: $GITHUB_TOKEN_FILE"
  echo "Create it with your token and set permissions: chmod 600 $GITHUB_TOKEN_FILE"
  exit 1
fi

# Optional: warn if file permissions are too permissive
perm=$(stat -c '%a' "$GITHUB_TOKEN_FILE")
if (( perm & 0077 )); then
  echo "Warning: $GITHUB_TOKEN_FILE permissions are $perm — consider 'chmod 600 $GITHUB_TOKEN_FILE'"
fi

# Read token (trim whitespace)
GITHUB_TOKEN="$(tr -d '[:space:]' < "$GITHUB_TOKEN_FILE")"
if [[ -z "$GITHUB_TOKEN" ]]; then
  echo "GitHub token file is empty: $GITHUB_TOKEN_FILE"
  exit 1
fi


# Version
if [[ -f "$VERSION_FILE" ]]; then
  VERSION="$(tr -d '[:space:]' < "$VERSION_FILE")"
  [[ -n "$VERSION" ]] || { echo "VERSION file empty"; exit 1; }
else
  VERSION="v$(date -u '+%Y%m%d%H%M%S')"
fi

echo "Repo: $GITHUB_REPO"
echo "Version: $VERSION"

# Build artifacts
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"
artifacts=()

for GOOS in "${GOOS_LIST[@]}"; do
  for GOARCH in "${GOARCH_LIST[@]}"; do
    OUT_NAME="${BINARY_BASE}-${GOOS}-${GOARCH}"
    [[ "$GOOS" != "windows" ]] || OUT_NAME="${OUT_NAME}.exe"
    OUT_PATH="${BUILD_DIR}/${OUT_NAME}"
    echo "Building $OUT_PATH (GOOS=$GOOS GOARCH=$GOARCH)"
    env GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -ldflags "-X main.version=${VERSION}" -o "$OUT_PATH" "$MAIN_PKG"
    TAR_PATH="${OUT_PATH}.tar.gz"
    tar -C "$BUILD_DIR" -czf "$TAR_PATH" "$(basename "$OUT_PATH")"
    artifacts+=("$TAR_PATH")
  done
done

printf 'Built artifacts:\n'
printf '%s\n' "${artifacts[@]}"

# Commit source changes (optional) and push tag
git add -A
git commit -m "chore(release): ${VERSION}" || echo "No changes to commit."
git tag -a "$VERSION" -m "Release $VERSION" || echo "Tag $VERSION exists."
git push "$GIT_REMOTE" "$GIT_BRANCH"
git push "$GIT_REMOTE" "$VERSION"

# Create release via GitHub API if not exists
API="https://api.github.com"
AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"
USER_AGENT_HEADER="User-Agent: build-and-upload-script"

# Check if release exists by tag
release_resp=$(curl -sS -H "$AUTH_HEADER" -H "$USER_AGENT_HEADER" "$API/repos/${GITHUB_REPO}/releases/tags/${VERSION}" || true)
release_id=$(jq -r '.id // empty' <<<"$release_resp" || true)

if [[ -n "$release_id" ]]; then
  echo "Using existing release id=$release_id for tag $VERSION"
else
  echo "Creating release for tag $VERSION"
  create_payload=$(jq -n --arg tag "$VERSION" --arg name "$VERSION" --arg body "Release $VERSION" \
    '{tag_name:$tag, name:$name, body:$body, draft:false, prerelease:false}')
  create_resp=$(curl -sS -X POST -H "$AUTH_HEADER" -H "Content-Type: application/json" -H "$USER_AGENT_HEADER" \
    -d "$create_payload" "$API/repos/${GITHUB_REPO}/releases")
  release_id=$(jq -r '.id' <<<"$create_resp")
  if [[ "$release_id" == "null" || -z "$release_id" ]]; then
    echo "Failed creating release. Response:"
    jq . <<<"$create_resp"
    exit 1
  fi
  echo "Created release id=$release_id"
fi

# Upload each artifact to release
for file in "${artifacts[@]}"; do
  filename="$(basename "$file")"
  echo "Uploading $filename"
  # GitHub upload URL: get upload_url from release
  upload_url=$(curl -sS -H "$AUTH_HEADER" -H "$USER_AGENT_HEADER" "$API/repos/${GITHUB_REPO}/releases/${release_id}" | jq -r '.upload_url')
  # upload_url contains {?name,label}; strip after {
  upload_url="${upload_url%%\{*}"
  # If asset with same name exists, delete it first
  existing_asset_id=$(curl -sS -H "$AUTH_HEADER" -H "$USER_AGENT_HEADER" "$API/repos/${GITHUB_REPO}/releases/${release_id}/assets" | jq -r --arg name "$filename" '.[] | select(.name==$name) | .id' | head -n1 || true)
  if [[ -n "$existing_asset_id" ]]; then
    echo "Deleting existing asset id=$existing_asset_id"
    curl -sS -X DELETE -H "$AUTH_HEADER" -H "$USER_AGENT_HEADER" "$API/repos/${GITHUB_REPO}/releases/assets/${existing_asset_id}"
  fi
  # Upload asset (must set Content-Type; octet-stream is fine)
  resp=$(curl -sS -X POST -H "$AUTH_HEADER" -H "$USER_AGENT_HEADER" -H "Content-Type: application/gzip" \
    --data-binary @"$file" "${upload_url}?name=$(printf '%s' "$filename" | jq -s -R -r @uri)")
  ok_id=$(jq -r '.id // empty' <<<"$resp")
  if [[ -z "$ok_id" ]]; then
    echo "Upload failed. Response:"
    jq . <<<"$resp"
    exit 1
  fi
  echo "Uploaded asset id=$ok_id"
done

echo "Release $VERSION published with assets:"
curl -sS -H "$AUTH_HEADER" -H "$USER_AGENT_HEADER" "$API/repos/${GITHUB_REPO}/releases/${release_id}" | jq -r '.assets[] | .name + " -> " + .browser_download_url'

