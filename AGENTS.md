# Agent Integration Guide

How to use `ws` in your agent, daemon, or automation.

## Channel Resolution (Primary Use Case)

```bash
ws resolve-channel C0ABC123
# stdout: /absolute/path/to/project
# exit 0 on success, exit 1 if not found

ws resolve-channel C0ABC123 --json
# {"channel":"C0ABC123","path":"/absolute/path","workspace":"/path/to.code-workspace"}
```

This scans **all** workspace files in `$WS_WORKSPACES_DIR` (default `~/workspaces/`). A channel mapped in any workspace will be found.

## Querying Projects

```bash
# All projects as JSON array
ws list --json
# [{"name":"🟢 My Project","path":"./my-project","slack_channel":"C0ABC123","description":"..."}]

# Everything in one call (projects + channel mappings)
ws dump-config
# {"workspace_path":"...","projects":[...],"channel_mappings":{"C0ABC123":"./my-project"}}

# Just channel mappings
ws dump-config --section channels
# {"C0ABC123":"./my-project","C0DEF456":"./other"}

# Just projects
ws dump-config --section projects
```

## Modifying State

```bash
# Set project status
ws status my-project --active     # 🟢
ws status my-project --blocked    # 🔴

# Link a Slack channel
ws channel my-project --set C0ABC123

# Add a new project
ws add ./path --name "🟢 Name" --slack C0ABC123 --desc "Description"
```

Writes are atomic (temp file + rename) with automatic backup (`.backup` file).

## Activity Indicators

For real-time "who's working on what" in VS Code sidebar:

```bash
ws activity working --project my-project   # Adds ⚡
ws activity waiting --project my-project   # Adds ❓
ws activity idle --project my-project      # Clears indicator
```

## Claude Code Hook

The `hook` command reads Claude Code hook events from stdin and auto-updates activity indicators.

**Quick setup:** Merge [`hooks/claude-code.json`](hooks/claude-code.json) into `~/.claude/settings.json`. See [`hooks/README.md`](hooks/README.md) for details.

Input format:
```json
{"hook_event_name": "PreToolUse", "cwd": "/path/to/project"}
```

Events → indicators:
- `PreToolUse`, `SubagentStart`, `UserPromptSubmit` → `working` (⚡)
- `Notification` (idle/permission), `PermissionRequest` → `waiting` (❓)
- `Stop`, `PostToolUse`, `SubagentStop`, `SessionEnd` → `idle` (clear)

## Fuzzy Matching

Project arguments match against both path and name (case-insensitive, emoji-stripped):

```bash
ws channel slop          # Matches "./Projects/slop"
ws channel "Project B"   # Matches "🔵 Project B"
ws channel asyncbot      # Matches "./asyncmode/asyncbot"
```

## Workspace Auto-Detection

Resolution order:
1. `--workspace` flag or `$WS_WORKSPACE` env var
2. Scan `$WS_WORKSPACES_DIR` for workspace containing current directory
3. First `*.code-workspace` file in `$WS_WORKSPACES_DIR` (alphabetical)

## Custom Fields

Metadata lives in `x-pa` inside each folder entry. VS Code ignores it.

```jsonc
{
  "path": "./my-project",
  "name": "🟢 My Project",
  "x-pa": {
    "slack_channel": "C0ABC123",
    "description": "Project description"
  }
}
```

Read fields: `ws list --json` or `ws dump-config`.
Write fields: `ws channel --set`, `ws add --slack/--desc`, or edit the workspace file directly.
