# Workspace CLI (`ws`)

CLI for managing VS Code workspace files. Single binary, zero runtime dependencies.

## Install

```bash
go install github.com/svenmeys/workspace-cli@latest
# Binary installs as 'workspace-cli', rename to 'ws':
cp "$(go env GOPATH)/bin/workspace-cli" /usr/local/bin/ws
```

## Quick Reference

```bash
# List projects
ws list
ws list --json  # Machine-readable output

# Add a project
ws add ./my-project --name "🟢 My Project" --slack C0ABC123 --desc "Description"

# Update status
ws status my-project --active    # 🟢
ws status my-project --progress  # 🔵
ws status my-project --paused    # 🟡
ws status my-project --blocked   # 🔴
ws status my-project --dormant   # ⚪

# Slack channel management
ws channel my-project                    # Get channel
ws channel my-project --set C0ABC123     # Set channel

# Config dump for daemons
ws dump-config                           # Full config as JSON
ws dump-config --section channels        # Just channel mappings
ws dump-config --section projects        # Just project list

# Resolve channel to project path
ws resolve-channel C0ABC123              # Returns absolute path

# Validation
ws validate

# Shell completions (built-in via cobra)
ws completion bash
ws completion zsh
ws completion fish
```

## Environment Variables

- `WS_WORKSPACE` - Path to workspace file (default: auto-detect from cwd, then first `*.code-workspace` in `~/workspaces/`)
- `WS_WORKSPACES_DIR` - Directory containing workspace files (default: `~/workspaces/`)
- `WS_CAPABILITIES` - Path to capabilities.yaml for sync

## Project Structure

```
main.go
cmd/                           # Command definitions (cobra)
  root.go, list.go, add.go, status.go, channel.go,
  dumpconfig.go, syncchannels.go, listall.go,
  resolvechannel.go, validate.go, activity.go, hook.go
internal/workspace/            # Core business logic
  config.go, workspace.go, sync.go, hooks.go
```

## For Agents

1. **Prefer `--json` output** for parsing: `ws list --json`
2. **Use `dump-config`** to get all workspace info in one call
3. **Project matching is fuzzy** - partial path or name works
4. **Status emojis are automatic** - just use `--active`, `--blocked`, etc.

## x-pa Custom Fields

The workspace file supports custom fields prefixed with `x-pa`:

```json
{
  "folders": [
    {
      "path": "./my-project",
      "name": "🟢 My Project",
      "x-pa": {
        "slack_channel": "C0ABC123",
        "description": "Project description"
      }
    }
  ]
}
```

VS Code ignores these fields. They're used for agent workflows.
