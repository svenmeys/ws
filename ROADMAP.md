# Roadmap

## v0.1.0 — Shipped

| Feature | Status |
|---------|--------|
| Core commands (`list`, `add`, `channel`, `status`) | ✅ |
| `resolve-channel` for daemon integration | ✅ |
| `dump-config` for agent workflows | ✅ |
| `sync-channels` to capabilities.yaml | ✅ |
| `validate` workspace integrity | ✅ |
| Backup on write | ✅ |
| Shell completions (bash/zsh/fish) | ✅ Built-in via cobra |
| `--workspace` flag per-command | ✅ |
| Rewrite in Go (single binary) | ✅ |

## Next Up

| Item | Notes |
|------|-------|
| Homebrew tap | `brew install ws` via goreleaser |
| CI/CD | GitHub Actions for release builds |
| `ws remove` | Remove project from workspace |

## Stretch Goals

| Feature | Unlock When |
|---------|-------------|
| Project templates | Adding 10+ projects |
| Watch mode | Frequent manual edits |

---

*Last updated: 2026-02-20*
