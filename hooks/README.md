# Claude Code Activity Hooks

These hooks make `ws` update your VS Code sidebar in real-time to show which project Claude is working on.

## What You Get

- **⚡** appears next to the project name when Claude is actively working (tool calls, writing code)
- **❓** appears when Claude is waiting for your input (permission prompts, idle)
- Indicator clears automatically when Claude stops or the session ends

## Install

### Option 1: Global (all projects)

Merge `claude-code.json` into `~/.claude/settings.json`:

```bash
# If you don't have hooks yet, copy the whole hooks section:
# Open ~/.claude/settings.json and add the "hooks" key from claude-code.json
```

### Option 2: Per-project

Copy the hooks section into your project's `.claude/settings.json`:

```bash
# In your project root:
mkdir -p .claude
# Add the "hooks" key from claude-code.json to .claude/settings.json
```

## Requirements

- `ws` binary in your `$PATH` (see main README for install)
- A VS Code `.code-workspace` file in `~/workspaces/` with your projects
- Projects must have a `path` that matches your working directory

## How It Works

1. Claude Code fires hook events (PreToolUse, Stop, etc.) as JSON to stdin
2. `ws hook` reads the event and determines the indicator type (working/waiting/idle)
3. `ws` finds which project matches the `cwd` from the event
4. Updates the project name in your `.code-workspace` file with the indicator emoji
5. VS Code picks up the file change and refreshes the sidebar

The hook runs with a 5-second timeout so it never blocks Claude.
