# Maestro - Multi-Container Claude

Maestro (Multi-Container Claude) is a tool for managing isolated Docker containers for Claude Code development. It enables you to run multiple Claude instances in parallel, each in their own isolated environment with proper branch management and network firewalls.

## Requirements

- **Docker**: Must be installed and running on your system
  - Install from [docker.com](https://www.docker.com/get-started)
  - Ensure Docker daemon is running before using Maestro
- **Go**: Required for building Maestro (1.20 or later)
- **Claude Code**: Authentication required (handled via `maestro auth`)

### Optional Dependencies

- **terminal-notifier** (macOS only): For custom notification icons
  ```bash
  brew install terminal-notifier
  ```
  Without this, notifications still work using macOS's built-in `osascript`

## Features

### Core Features
- 🔥 **Firewall Protection**: Each container has network restrictions to prevent accidental data access
- 🌳 **Automatic Branch Management**: Creates and manages git branches for each task
- 📦 **Full Isolation**: Each container gets a complete copy of your project
- 🔔 **Attention Monitoring**: See which containers need your attention via tmux bell detection
- 🎯 **Easy Navigation**: Connect/disconnect from containers without losing state
- 🧰 **Full Dev Environment**: Includes Node.js, Go, Python (with UV), and all standard tools
- ♻️ **Persistent Caches**: npm, UV, and command history persist across container restarts

### Advanced Features
- 🤖 **Background Daemon**: Auto-monitors containers for token expiry and attention needs
- 📬 **Smart Notifications**: Desktop alerts when containers need attention (configurable delay)
- 📊 **Enhanced Status View**: Git status indicators (Δ changes, ↑ ahead, ↓ behind, ✓ clean)
- 🔄 **Quick Restart**: Restart crashed Claude processes without stopping containers
- ⏰ **Quiet Hours**: Configure notification-free hours for uninterrupted sleep

## Installation

### Quick Install (Recommended)

Install Maestro with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/uprockcom/maestro/main/install.sh | bash
```

This will:
- Download the latest release for your platform
- Install the `maestro` binary to `/usr/local/bin`
- Pull the Docker image
- Create the config directory at `~/.maestro`
- Set up the example configuration

**Prerequisites:**
- Docker must be installed and running
- `curl` command available

### Build from Source (Optional)

For development or if you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/uprockcom/maestro.git
cd maestro

# Build everything (binary + Docker image)
make all

# Install to /usr/local/bin (requires sudo)
sudo make install
```

Available make targets:
```bash
make build          # Build maestro binary only
make docker         # Build Docker image only
make all            # Build both binary and image
make install        # Install maestro to /usr/local/bin
make test           # Run Go tests
make clean          # Remove built binaries
make help           # Show all available targets
```

### Configuration

The installer creates a default configuration at `~/.maestro/config.yml`. Edit it to customize:

```bash
# Edit the config to add your additional folders and domains
nano ~/.maestro/config.yml
```

If you built from source, manually copy the example config:
```bash
cp config.yml.example ~/.maestro/config.yml
```

## Usage

### Creating a New Container

```bash
# Quick task description
maestro new "implement OAuth authentication"

# From a specification file
maestro new -f specs/feature-design.md

# Interactive mode
maestro new
```

This will:
1. Use Claude to generate an appropriate branch name
2. Create a new container with incremented numbering (e.g., `maestro-feat-oauth-1`)
3. Copy your entire project into the container
4. Create and checkout the new git branch
5. Start tmux with Claude in planning mode
6. Present you with the container info

### Managing Containers

```bash
# List all containers with status indicators
maestro list        # or: maestro ls, maestro ps

# Connect to a container
maestro connect feat-oauth-1

# Restart a crashed Claude process (preserves container state)
maestro restart feat-oauth-1

# Full container restart (if needed)
maestro restart feat-oauth-1 --full

# Stop a specific container
maestro stop feat-oauth-1

# Stop all dormant containers (where Claude has exited)
maestro stop

# Clean up stopped containers and their volumes
maestro cleanup

# Remove all containers (including running) and their volumes
maestro cleanup --all

# Clean up orphaned volumes (volumes without containers)
maestro cleanup-volumes
```

The `maestro list` command shows comprehensive status indicators:
- **GIT**: Git status (Δ79 = 79 changes, ↑2 = 2 commits ahead, ↓1 = 1 behind, ✓ = clean)
- **ACTIVITY**: Time since last activity
- **AUTH**: Token expiration (✓ valid, ⚠ expiring soon, ✗ expired)
- **🔔**: Container needs attention
- **💤**: Container is dormant (Claude has exited)

### Inside the Container

When connected to a container via `maestro connect`:

- **Window 0**: Claude Code running in yolo mode (auto-approved)
- **Window 1**: Shell for manual commands
- **Switch windows**: `Ctrl+b 0` (Claude) or `Ctrl+b 1` (shell)
- **Detach**: `Ctrl+b d` (returns you to host, container keeps running)

The tmux status line shows:
- Container name
- Current git branch
- Bell indicator when Claude needs attention

### Network Management

```bash
# Add a domain temporarily to a running container
maestro add-domain feat-oauth-1 api.example.com

# The tool will offer to add it to ~/config.yml for permanent access
```

## Background Daemon

Maestro includes a background daemon that monitors your containers for token expiry and attention needs.

### Starting the Daemon

```bash
# Start the daemon (runs in background)
maestro daemon start

# Check daemon status
maestro daemon status

# View daemon logs
maestro daemon logs

# Stop the daemon
maestro daemon stop
```

The daemon will:
- ✅ Auto-check containers every 30 minutes (configurable)
- ✅ Send desktop notifications when containers need attention (after 5 minute delay)
- ✅ Warn when tokens are expiring (< 1 hour remaining)
- ✅ Support quiet hours to avoid night interruptions
- ✅ Rate limit notifications (max once per 30 minutes per container)

### Daemon Features

**Token Monitoring**: Automatically checks token expiration and will support auto-refresh in the future.

**Smart Notifications**: Only notifies after containers have needed attention for a configurable threshold (default 5 minutes), preventing notification spam.

**Quiet Hours**: Configure time ranges when notifications should be suppressed (e.g., 23:00-08:00).

**Activity Tracking**: Monitors container activity and Claude process health.

**Custom Notification Icons**: On macOS, install `terminal-notifier` for custom icon support:
```bash
brew install terminal-notifier
```

Without `terminal-notifier`, notifications still work via macOS's built-in `osascript`, but will use the Terminal/iTerm icon. The daemon will automatically detect and use `terminal-notifier` if available.

## Configuration

Edit `~/config.yml` to customize:

```yaml
claude:
  config_path: ~/.claude       # Your Claude auth directory
  mcl_claude_path: ~/.maestro/.claude  # Maestro's centralized auth storage
  default_mode: yolo           # Auto-approve mode

containers:
  prefix: mcl-                 # Container name prefix
  image: mcl:latest            # Docker image name
  resources:
    memory: 4g                 # Memory limit
    cpus: "2"                  # CPU limit

firewall:
  allowed_domains:             # Whitelisted domains
    - github.com
    - pypi.org
    - api.anthropic.com
    # Add your domains here

sync:
  additional_folders:          # Folders to copy as siblings
    - ~/Documents/Code/mcp-servers
    - ~/Documents/Code/helpers

github:
  enabled: false               # Enable GitHub CLI (gh) integration
  config_path: ~/.maestro/gh       # Path to gh config directory (managed by mcl auth)

daemon:
  check_interval: 30m          # How often to check containers
  show_nag: true               # Show reminder to start daemon (disable: false)
  token_refresh:
    enabled: true              # Auto-refresh expiring tokens
    threshold: 6h              # Refresh when < 6h remaining
  notifications:
    enabled: true              # Send desktop notifications
    attention_threshold: 5m    # Wait 5m before notifying
    notify_on:
      - attention_needed       # Notify when container needs attention
      - token_expiring         # Notify when token < 1h
    quiet_hours:
      start: "23:00"           # Optional: quiet hours start (24h format)
      end: "08:00"             # Optional: quiet hours end
```

### Configuration Notes

- **show_nag**: Set to `false` to disable the "start daemon" reminder in `maestro list`
- **check_interval**: How often the daemon checks containers (e.g., "30m", "1h", "15m")
- **attention_threshold**: How long to wait before sending notification (prevents spam)
- **quiet_hours**: Optional. Leave empty (`""`) to disable quiet hours
- Time formats: Use Go duration format ("30m", "6h") or 24-hour time ("23:00")

## Architecture

### Container Structure

```
/workspace/              # Your main project (copied from host)
/workspace/../mcp-servers/  # Additional folders (from config)
/home/node/.claude/      # Claude config (mounted read-only)
```

### Persistent Volumes

Each container has named volumes for:
- npm cache (`<container>-npm`)
- UV cache (`<container>-uv`)
- Command history (`<container>-history`)

### Network Firewall

Containers can only access:
- Whitelisted domains from config
- GitHub API endpoints (auto-detected)
- Local Docker network
- DNS resolution

## Authentication

Maestro manages authentication for both Claude Code and GitHub CLI in a centralized location (`~/.maestro/`).

### Initial Setup

Run the authentication command once:

```bash
maestro auth
```

This will:
1. **Claude Code Authentication**: Start a container where you'll authenticate with Claude via OAuth
2. **Credential Sync**: Automatically sync new credentials to all running containers
3. **GitHub CLI Authentication** (optional): Optionally authenticate with GitHub for PR reviews and other features

All credentials are stored in `~/.maestro/` and shared (read-only) with all containers.

### Re-authenticating

If your authentication expires, simply run `maestro auth` again. By default, it will automatically update credentials in all running containers. Use `--no-sync` to skip this:

```bash
maestro auth --no-sync
```

### GitHub CLI Integration

When you run `maestro auth`, you'll be prompted to set up GitHub CLI. This enables features like:
- `gh pr review 123`
- `gh issue list`
- `gh repo view`

To enable GitHub CLI in containers, set in `~/config.yml`:
```yaml
github:
  enabled: true
  config_path: ~/.maestro/gh  # Managed by mcl auth
```

### Security Note

Authentication tokens are stored in `~/.maestro/` and mounted read-only in containers. GitHub CLI integration is opt-in during setup.

## Token Management

Claude authentication tokens automatically expire after approximately 1 week. Maestro provides tools to manage token expiration gracefully.

### Checking Token Status

The `maestro list` command shows authentication status for each running container:

```bash
maestro list
```

Output example:
```
NAME                STATUS   BRANCH           AUTH STATUS  ATTENTION
----                ------   ------           -----------  ---------
feat-oauth-1        running  feat/oauth       ✓ 147.2h
fix-api-bug-1       running  fix/api-bug      ⚠ 2.3h
refactor-db-1       running  refactor/db      ✗ EXPIRED    💤 DORMANT
```

Auth status indicators:
- **✓ Xh**: Token is valid for X hours (green indicator)
- **⚠ Xh**: Token expires in less than 24 hours (warning)
- **✗ EXPIRED**: Token has expired and needs refresh

### Refreshing Tokens

Claude CLI automatically refreshes tokens when actively used in a container. Use `maestro refresh-tokens` to find and propagate the freshest token:

```bash
maestro refresh-tokens
```

This command:
1. Scans all running containers and the host for credentials
2. Finds the container with the freshest token (Claude auto-refreshes during normal use)
3. Copies the fresh token to all other containers and the host
4. Ensures new containers will use the fresh token

**Example output:**
```
Scanning for credentials...
  ✓ Host: EXPIRED 2.8h ago
  ✓ maestro-feat-oauth-1: Valid for 147.2h
  ✓ maestro-fix-api-bug-1: EXPIRED 2.8h ago

✓ Found fresh token in maestro-feat-oauth-1
  Expires: Sat, 18 Oct 2025 06:33:27 PDT
  Status: Valid for 147.2h

Syncing credentials...
  ✓ Synced to host
  ✓ Synced to maestro-fix-api-bug-1

✅ Refresh complete! Synced to 2 location(s).
```

### Re-authenticating

If all tokens are expired, `maestro refresh-tokens` will prompt you to run `maestro auth`:

```bash
maestro auth
```

This will:
1. Start a temporary authentication container
2. Complete OAuth flow in your browser
3. Automatically sync new credentials to all running containers

### Token Expiration Warnings

When creating a new container, Maestro will warn you if tokens are expired or expiring soon:

```bash
maestro new "implement feature"

⚠️  WARNING: Authentication token is EXPIRED!
   Status: EXPIRED 2.8h ago
   Run 'maestro auth' or 'maestro refresh-tokens' to get a fresh token.

Continue creating container with expired token? (y/N):
```

### Best Practices

- **Check token status regularly**: Run `maestro list` to see auth status for all containers
- **Use `refresh-tokens` first**: If you see expired tokens, try `maestro refresh-tokens` before running `maestro auth` (it's faster and reuses existing fresh tokens)
- **Run `auth` when needed**: Only run `maestro auth` if all tokens are expired or `refresh-tokens` fails
- **Monitor expiration warnings**: If you see "⚠" warnings in `maestro list`, consider refreshing tokens soon

## Tips

### Working with Multiple Versions

The numbering system (`-1`, `-2`, etc.) lets you create multiple implementations:

```bash
maestro new "implement auth with OAuth"     # Creates maestro-feat-auth-oauth-1
maestro new "implement auth with JWT"       # Creates maestro-feat-auth-jwt-1
maestro new "implement auth with OAuth"     # Creates maestro-feat-auth-oauth-2
```

### Restarting Stopped Containers

```bash
# Restart a stopped container
docker start maestro-feat-oauth-1

# Then connect normally
maestro connect feat-oauth-1
```

### Monitoring Activity

The `maestro list` command shows status indicators for each container:
- **🔔 NEEDS ATTENTION**: Detected via tmux's bell feature when Claude has activity
- **💤 DORMANT**: Container is running but Claude process has exited

You can quickly clean up dormant containers with `maestro stop` (no arguments).

## Troubleshooting

### Container won't start

Check Docker logs:
```bash
docker logs maestro-feat-name-1
```

### Firewall blocking needed domain

Add it temporarily:
```bash
maestro add-domain container-name api.example.com
```

Then add to `~/config.yml` for permanent access.

### Can't connect to container

Ensure it's running:
```bash
docker ps
docker start <container-name>
```

### Claude not authenticated

Ensure your `~/.claude/` directory contains valid auth tokens. The directory is mounted read-only into each container.

## Development

To modify maestro itself:

```bash
# Make changes to Go files

# Run tests
make test

# Rebuild binary
make build

# Test changes
./bin/maestro --help

# Install updated binary
make install
```

To modify the container image:

```bash
# Edit docker/Dockerfile

# Rebuild image
make docker

# Or rebuild everything
make all
```

## License

Apache 2.0 - see [LICENSE](LICENSE) file for details