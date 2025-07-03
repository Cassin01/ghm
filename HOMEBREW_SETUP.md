# Homebrew Automation Setup

This document explains how to set up automated Homebrew formula updates for the ghm project.

## GitHub Secrets Required

⚠️ **IMPORTANT**: The default `GITHUB_TOKEN` cannot access other repositories. To enable automated Homebrew formula updates, you **MUST** set up a Personal Access Token.

### HOMEBREW_TAP_TOKEN (Required for Automation)

To enable cross-repository access to the `homebrew-ghm` repository:

1. Create a Personal Access Token (PAT) with the following permissions:
   - `repo` (Full control of private repositories)
   - `workflow` (Update GitHub Action workflows)

2. Add the token as a repository secret:
   - Go to your repository on GitHub
   - Navigate to Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `HOMEBREW_TAP_TOKEN`
   - Value: Your PAT

**Note**: If `HOMEBREW_TAP_TOKEN` is not set, the Homebrew update step will be skipped and you'll need to update the formula manually.

## How It Works

1. **Release Trigger**: When a new tag is pushed (via tagpr), the release workflow runs
2. **Build & Release**: Creates binaries and GitHub release
3. **Homebrew Update**: Automatically updates the Homebrew formula with:
   - New version number
   - Updated SHA256 checksum
   - New download URL

## Automation Features

- **Version Detection**: Automatically extracts version from Git tag
- **Checksum Calculation**: Downloads release archive and calculates SHA256
- **Formula Update**: Updates the Homebrew formula file
- **Commit & Push**: Commits changes to homebrew-ghm repository
- **Pre-release Skip**: Skips Homebrew updates for pre-release versions (containing `-`)

## Manual Homebrew Formula Update

If you choose not to set up automation, you can update the Homebrew formula manually:

1. **Calculate SHA256 for new release:**
   ```bash
   curl -sL https://github.com/Cassin01/ghm/archive/v0.1.1.tar.gz | sha256sum
   ```

2. **Update the formula in homebrew-ghm repository:**
   - Clone: `git clone https://github.com/Cassin01/homebrew-ghm.git`
   - Edit `ghm.rb` with new version and SHA256
   - Commit: `git commit -am "Update ghm to v0.1.1"`
   - Push: `git push origin main`

## Testing Automation

To test the Homebrew update script manually (requires HOMEBREW_TAP_TOKEN):

```bash
# Set your token
export GITHUB_TOKEN=your_homebrew_tap_token

# Make script executable
chmod +x scripts/update-homebrew.sh

# Test with a version
./scripts/update-homebrew.sh v0.1.1
```

## Troubleshooting

### Permission Errors
- Ensure `HOMEBREW_TAP_TOKEN` has correct permissions
- Check that the token hasn't expired

### Network Errors
- GitHub API rate limits may cause temporary failures
- The workflow will retry automatically

### Formula Errors
- Check that the formula syntax is correct
- Verify SHA256 calculation is accurate

## Repository Structure

```
ghm/
├── scripts/
│   └── update-homebrew.sh    # Homebrew update script
├── .github/
│   └── workflows/
│       ├── release.yml       # Main release workflow
│       └── tagpr.yml         # tagpr automation
└── HOMEBREW_SETUP.md         # This documentation
```