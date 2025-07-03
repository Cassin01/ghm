#!/bin/bash
set -e

# Homebrew Formula Update Script
# This script updates the Homebrew formula with new version and checksum

VERSION="$1"
REPO_OWNER="Cassin01"
REPO_NAME="ghm"
HOMEBREW_TAP_REPO="homebrew-ghm"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v0.1.1"
    exit 1
fi

# Check if GITHUB_TOKEN is set for authentication
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ Error: GITHUB_TOKEN is not set"
    echo "This token is required to push updates to the homebrew-ghm repository"
    echo "Please set HOMEBREW_TAP_TOKEN secret in your repository settings"
    exit 1
fi

# Remove 'v' prefix if present
VERSION_NUMBER="${VERSION#v}"

echo "Updating Homebrew formula for version $VERSION_NUMBER"

# Download release archive and calculate SHA256
ARCHIVE_URL="https://github.com/$REPO_OWNER/$REPO_NAME/archive/$VERSION.tar.gz"
echo "Downloading archive from: $ARCHIVE_URL"

TEMP_DIR=$(mktemp -d)
ARCHIVE_PATH="$TEMP_DIR/archive.tar.gz"

curl -sL "$ARCHIVE_URL" -o "$ARCHIVE_PATH"
SHA256=$(sha256sum "$ARCHIVE_PATH" | cut -d' ' -f1)

echo "Calculated SHA256: $SHA256"

# Clone homebrew tap repository with authentication
TAP_DIR="$TEMP_DIR/$HOMEBREW_TAP_REPO"
git clone "https://$GITHUB_TOKEN@github.com/$REPO_OWNER/$HOMEBREW_TAP_REPO.git" "$TAP_DIR"

cd "$TAP_DIR"

# Update formula
FORMULA_FILE="ghm.rb"

# Create new formula content
cat > "$FORMULA_FILE" << EOF
class Ghm < Formula
  desc "GitHub Manager - manage multiple instances of the same repository"
  homepage "https://github.com/$REPO_OWNER/$REPO_NAME"
  url "https://github.com/$REPO_OWNER/$REPO_NAME/archive/$VERSION.tar.gz"
  sha256 "$SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/ghm"
  end

  test do
    system "#{bin}/ghm", "--help"
    system "#{bin}/ghm", "root"
  end
end
EOF

echo "Updated formula:"
cat "$FORMULA_FILE"

# Commit and push changes
git config user.name "github-actions[bot]"
git config user.email "github-actions[bot]@users.noreply.github.com"

git add "$FORMULA_FILE"
git commit -m "Update ghm to $VERSION_NUMBER

Auto-updated by release workflow

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push origin main

echo "Successfully updated Homebrew formula to version $VERSION_NUMBER"

# Cleanup
rm -rf "$TEMP_DIR"