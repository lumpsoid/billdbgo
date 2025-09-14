#!/usr/bin/env bash
set -euo pipefail

# Config — edit if needed
MAIN_PKG="./cmd/server"
BINARY_BASE="billdbgo"
BUILD_DIR="dist"
GOOS_LIST=("linux")
GOARCH_LIST=("amd64" "arm64")
VERSION_FILE="VERSION"
GIT_REMOTE="${GIT_REMOTE:-origin}"
GIT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"

# GitHub repo info (owner/repo). Try to infer from origin if not set.
GITHUB_REPO="${GITHUB_REPO:-}"
if [[ -z "$GITHUB_REPO" ]]; then
  ORIGIN_URL="$(git config --get "remote.${GIT_REMOTE}.url" || true)"
  # handle URLs like:
  # git@github.com:owner/repo.git
  # https://github.com/owner/repo.git
  # https://github.com/owner/repo
  if [[ "$ORIGIN_URL" =~ github.com[:/](.+/.+)$ ]]; then
    GITHUB_REPO="${BASH_REMATCH[1]}"
    # strip optional .git suffix if present
    GITHUB_REPO="${GITHUB_REPO%.git}"
    echo "$GITHUB_REPO"
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
# read current version file (if exists)
if [[ -f "$VERSION_FILE" ]]; then
  CURRENT_VERSION="$(tr -d '[:space:]' < "$VERSION_FILE")"
else
  CURRENT_VERSION=""
fi

# read latest git tag (if any)
if git rev-parse --verify --quiet refs/tags >/dev/null; then
  LATEST_TAG="$(git describe --tags --abbrev=0 2>/dev/null || true)"
else
  LATEST_TAG="$(git describe --tags --abbrev=0 2>/dev/null || true)"
fi
LATEST_TAG="${LATEST_TAG:-<none>}"

# decide default VERSION if no VERSION file
if [[ -n "$CURRENT_VERSION" ]]; then
  VERSION="$CURRENT_VERSION"
else
  VERSION="v$(date -u '+%Y%m%d%H%M%S')"
fi

# prompt developer
echo "Latest git tag: ${LATEST_TAG}"
echo "Current VERSION file: ${CURRENT_VERSION:-<none>}"
echo "Proposed release version: ${VERSION}"
printf "Is this appropriate? (y/N): "
read -r reply
case "$reply" in
  [yY]) ;;
  *)
    echo "Aborting per developer response."
    exit 1
    ;;
esac

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
    case "$GOARCH" in
      amd64) CC_COMPILER="x86_64-linux-gnu-gcc" ;;
      arm64) CC_COMPILER="aarch64-linux-gnu-gcc" ;;
      *) CC_COMPILER="" ;;
    esac

    if [[ -n "$CC_COMPILER" ]]; then
      env GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=1 CC="$CC_COMPILER" go build -ldflags "-X main.version=${VERSION}" -o "$OUT_PATH" "$MAIN_PKG"
    else
      env GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o "$OUT_PATH" "$MAIN_PKG"
    fi

    # Create a temporary staging directory for the release layout
    RELEASE_DIR="${BUILD_DIR}/${BINARY_BASE}-${VERSION}-${GOOS}-${GOARCH}"
    rm -rf "$RELEASE_DIR"
    mkdir -p "$RELEASE_DIR"

    # Place the binary under a directory named after BINARY_BASE-VERSION
    cp "$OUT_PATH" "$RELEASE_DIR/${BINARY_BASE}"

    # Copy web assets (templates and static) into the release dir
    # Adjust source paths if your project layout differs
    if [[ -d "./web/templates" ]]; then
      cp -r ./web/templates "$RELEASE_DIR/templates"
    else
      echo "Warning: ./web/templates not found"
    fi

    if [[ -d "./web/static" ]]; then
      cp -r ./web/static "$RELEASE_DIR/static"
    else
      echo "Warning: ./web/static not found"
    fi

    # Create tarball that contains the single top-level dir: billdbgo-VERSION/...
    TAR_PATH="${BUILD_DIR}/${BINARY_BASE}-${VERSION}-${GOOS}-${GOARCH}.tar.gz"
    (
      cd "$BUILD_DIR" || exit 1
      tar -czf "$(basename "$TAR_PATH")" "$(basename "$RELEASE_DIR")"
    )
    artifacts+=("$TAR_PATH")

    # cleanup release staging but keep built binary if you want (optional)
    rm -rf "$RELEASE_DIR"
  done
done

printf 'Built artifacts:\n'
printf '%s\n' "${artifacts[@]}"

# Commit source changes (optional) and push tag
git add "$VERSION_FILE"
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

