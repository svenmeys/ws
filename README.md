# ws — World Serpent

> *"I've seen things you people wouldn't believe... twenty-one projects across four workspaces, all connected by a single command."*

**Your workspace's central nervous system.** One binary. Zero dependencies. Instant channel-to-project routing.

`ws` is short for **workspace** — and named after [Jörmungandr](https://en.wikipedia.org/wiki/J%C3%B6rmungandr), the World Serpent of Norse mythology. The serpent so vast it wraps around all of Midgard and bites its own tail, connecting everything. That's what `ws` does for your projects.

## Why This Exists

I run 20+ projects across multiple VS Code workspaces. I also run an AI personal assistant daemon that receives Slack messages and needs to figure out which project they belong to. The question was always: **where does this channel ID map to?**

The answer used to be a hand-maintained YAML file. That broke constantly.

The fix: extend the VS Code `.code-workspace` file with custom fields (`x-pa`) that VS Code silently ignores, then build a CLI to query it. The workspace file becomes the single source of truth for project metadata — names, status, Slack channels, descriptions. No separate config. No database. Just JSON that VS Code already renders.

My PA daemon now calls `ws resolve-channel C0ABC123` and gets back `/Users/me/Projects/my-project`. Done.

## Install

```bash
go install github.com/svenmeys/ws@latest
cp "$(go env GOPATH)/bin/ws" /usr/local/bin/ws
```

Or build from source:

```bash
git clone https://github.com/svenmeys/ws.git
cd ws && go build -o /usr/local/bin/ws .
```

## Usage

```bash
ws list                              # What's in the workspace?
ws list --json                       # Machine-readable

ws add ./my-project --name "🟢 My Project" --slack C0ABC123

ws resolve-channel C0ABC123          # → /absolute/path/to/project

ws status my-project --active        # 🟢
ws status my-project --blocked       # 🔴
ws status my-project --paused        # 🟡
ws status my-project --progress      # 🔵
ws status my-project --dormant       # ⚪

ws channel my-project                # Get channel ID
ws channel my-project --set C0X      # Set channel ID

ws dump-config                       # Full JSON for scripts/daemons
ws dump-config --section channels    # Just channel mappings
ws dump-config --section projects    # Just project list

ws validate                          # Check workspace integrity

ws list-all                          # All workspaces, all projects

ws activity working --project foo    # Set activity indicator (⚡/❓)
ws hook                              # Claude Code hook handler (stdin)

ws completion zsh >> ~/.zshrc        # Shell completions
```

## How It Works

`ws` reads and writes VS Code `.code-workspace` files, extending them with `x-pa` custom fields:

```jsonc
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

VS Code ignores `x-pa`. Your agents and scripts don't.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WS_WORKSPACE` | auto-detect from cwd | Path to workspace file |
| `WS_WORKSPACES_DIR` | `~/workspaces/` | Directory containing workspace files |
| `WS_CAPABILITIES` | `~/.config/workspace-cli/capabilities.yaml` | Sync target for channel mappings |

## Claude Code Integration

`ws` includes a `hook` command that reads Claude Code hook events from stdin and updates project activity indicators in real-time:

- **⚡** when Claude is working (tool calls, subagents)
- **❓** when Claude is waiting (permission prompts, idle prompts)
- Cleared when idle or session ends

This shows up directly in VS Code's sidebar — you can see which project Claude is actively working on.

Ready-to-use hook configuration is in [`hooks/`](hooks/). Merge `hooks/claude-code.json` into your `~/.claude/settings.json` to enable it globally.

## For Agents

See [AGENTS.md](AGENTS.md) for integration details. The short version:

- Use `--json` for structured output
- `dump-config` gives you everything in one call
- `resolve-channel` is your hot path — channel ID in, project path out
- Project matching is fuzzy — partial names and paths work

## License

MIT
