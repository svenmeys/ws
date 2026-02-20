# Nexus

> *"I've seen things you people wouldn't believe... twenty-one projects across four workspaces, all connected by a single command."*

**Your workspace's central nervous system.** One binary. Zero dependencies. Instant channel-to-project routing.

Nexus manages VS Code `.code-workspace` files — mapping Slack channels to projects, tracking status, and giving your agents a single source of truth. Built for humans who run too many projects and the agents that help them.

## Install

```bash
# From source
go install github.com/svenmeys/workspace-cli@latest
cp "$(go env GOPATH)/bin/workspace-cli" /usr/local/bin/ws

# Or just build it
git clone https://github.com/svenmeys/workspace-cli.git
cd workspace-cli && go build -o /usr/local/bin/ws .
```

## Usage

```bash
# What's in the workspace?
ws list
ws list --json              # For agents

# Add a project
ws add ./my-project --name "🟢 My Project" --slack C0ABC123

# Route a Slack channel to its project
ws resolve-channel C0ABC123
# → /Users/you/Projects/my-project

# Update project status
ws status my-project --active     # 🟢
ws status my-project --blocked    # 🔴
ws status my-project --paused     # 🟡
ws status my-project --progress   # 🔵

# Channel management
ws channel my-project             # Get channel ID
ws channel my-project --set C0X   # Set channel ID

# Dump everything for scripts/daemons
ws dump-config --section channels

# Validate workspace integrity
ws validate

# Shell completions
ws completion zsh >> ~/.zshrc
```

## How It Works

Nexus reads and writes VS Code `.code-workspace` files, extending them with `x-pa` custom fields that VS Code silently ignores:

```json
{
  "folders": [
    {
      "path": "./my-project",
      "name": "🟢 My Project",
      "x-pa": {
        "slack_channel": "C0ABC123",
        "description": "The one that actually makes money"
      }
    }
  ]
}
```

Your workspace file becomes the single source of truth. No separate config. No database. Just JSON that VS Code already knows how to render.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WS_WORKSPACE` | auto-detect from cwd | Path to workspace file |
| `WS_WORKSPACES_DIR` | `~/workspaces/` | Directory containing workspace files |
| `WS_CAPABILITIES` | `~/.config/workspace-cli/capabilities.yaml` | Sync target for channel mappings |

## For Agents

- **Always use `--json`** — structured output, no surprises
- **`dump-config`** gives you everything in one call
- **`resolve-channel`** is your hot path — channel ID in, project path out
- **Project matching is fuzzy** — partial names and paths work
- **`ws hook`** reads Claude Code events from stdin for real-time status

## License

MIT
