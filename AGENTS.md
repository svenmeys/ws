# Workspace CLI

CLI tool for managing VS Code workspace files with agent integration.

## Overview

`ws` is a Go CLI that manipulates VS Code `.code-workspace` files, adding custom `x-pa` fields for agent workflows (Slack channel mappings, project metadata). Single binary, zero runtime dependencies.

## Tech Stack

- **Language**: Go 1.22+
- **CLI Framework**: Cobra
- **JSON Handling**: hujson (JSONC/HuJSON support)
- **Config**: gopkg.in/yaml.v3
- **Build**: `go build`

## Action Execution (CRITICAL)

*Execute first. Report completed actions. Never promise future actions.*

| DON'T | DO |
|-------|-----|
| "I will read the file" | [Read tool] → "Read complete: found X" |
| "Let me delegate this" | [Task tool] → "Delegated to Explore agent" |

Tool calls come BEFORE response text. Observable results only.

---

## Key Files

```
main.go                        # Entry point
cmd/                           # Command definitions (one per file)
  root.go                      # Root command, shared helpers
  list.go, add.go, status.go   # Project management
  channel.go, dumpconfig.go    # Config operations
  resolvechannel.go            # Channel → project lookup
  activity.go, hook.go         # Claude Code integration
internal/workspace/            # Business logic
  config.go                    # Environment-based configuration
  workspace.go                 # Core workspace file operations
  sync.go                      # Export to capabilities.yaml
  hooks.go                     # Claude Code hooks
```

## Commands Reference

```bash
ws list              # List projects
ws list-all          # List all workspaces and projects
ws add <path>        # Add project to workspace
ws status <project>  # Update status emoji
ws channel <project> # Get/set Slack channel
ws dump-config       # Export config as JSON
ws validate          # Validate workspace file
ws resolve-channel   # Get project path for channel ID
ws activity          # Set activity indicator
ws hook              # Claude Code hook handler
ws completion        # Shell completions (bash/zsh/fish)
```

## Environment Variables

- `WS_WORKSPACE` - Path to workspace file
- `WS_WORKSPACES_DIR` - Directory containing workspace files (default: `~/workspaces/`)
- `WS_CAPABILITIES` - Path to capabilities.yaml for sync

## Integration Points

- **Daemon integration**: Use `ws resolve-channel` for channel→project lookup
- **Config export**: `ws dump-config` for machine-readable workspace state
- **Claude Code hooks**: `ws hook` reads stdin for event-driven title updates

## Building

```bash
go build -o ws .
go install .  # Installs as 'workspace-cli' to GOPATH/bin
```

---

*Last updated: 2026-02-20*
