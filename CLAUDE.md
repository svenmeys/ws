# Workspace CLI (`ws`)

Single-binary Go CLI for managing VS Code `.code-workspace` files. Extends them with `x-pa` custom fields for project metadata (Slack channels, status, descriptions) that VS Code ignores.

See [AGENTS.md](AGENTS.md) for agent integration details.

## Project Structure

```
main.go
cmd/                           # Cobra command definitions (12 commands)
  root.go                      # loadWorkspace() helper, flag setup
  list.go, add.go, status.go, channel.go, dumpconfig.go,
  syncchannels.go, listall.go, resolvechannel.go,
  validate.go, activity.go, hook.go
internal/workspace/            # Core business logic
  workspace.go                 # Read/write workspace, folder ops, status updates
  config.go                    # Workspace discovery, path resolution
  hooks.go                     # Activity indicators, Claude Code hook handling
  sync.go                      # Sync channel mappings to capabilities.yaml
hooks/                         # Claude Code hook config (copy to ~/.claude/settings.json)
  claude-code.json             # Hook definitions for activity indicators
```

## Key Patterns

- **`loadWorkspace()`** in `cmd/root.go` — shared helper used by 7 commands to resolve + read workspace
- **`map[string]any`** throughout — workspace data is dynamic JSON, not typed structs (VS Code workspace format varies)
- **Atomic writes** — `WriteWorkspace` uses temp file + `os.Rename` with `.backup` of previous version
- **Fuzzy matching** — `FindFolder` matches by path substring or emoji-stripped case-insensitive name
- **JSONC support** — `hujson.Standardize` strips comments before JSON parsing

## Environment Variables

- `WS_WORKSPACE` — explicit workspace file path (skips auto-detect)
- `WS_WORKSPACES_DIR` — where to find `*.code-workspace` files (default `~/workspaces/`)
- `WS_CAPABILITIES` — sync target for channel mappings

## Testing

```bash
go test ./...    # 50 tests in internal/workspace/
go build .       # Single binary, no CGO
```
