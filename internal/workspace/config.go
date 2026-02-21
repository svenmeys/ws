package workspace

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetWorkspacesDir() string {
	if dir := os.Getenv("WS_WORKSPACES_DIR"); dir != "" {
		return expandHome(dir)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "workspaces")
}

func FindWorkspacesContaining(dir string) []string {
	dir, _ = filepath.Abs(dir)
	wsDir := GetWorkspacesDir()
	matches, err := filepath.Glob(filepath.Join(wsDir, "*.code-workspace"))
	if err != nil || len(matches) == 0 {
		return nil
	}
	var result []string
	for _, wsPath := range matches {
		ws, err := ReadWorkspace(wsPath)
		if err != nil {
			continue
		}
		for _, folder := range GetFolders(ws) {
			if s, ok := folder["path"].(string); ok {
				abs, _ := filepath.Abs(filepath.Join(filepath.Dir(wsPath), s))
				if strings.HasPrefix(dir, abs) {
					result = append(result, wsPath)
					break
				}
			}
		}
	}
	sort.Strings(result)
	return result
}

func GetDefaultWorkspace() (string, error) {
	if path := os.Getenv("WS_WORKSPACE"); path != "" {
		return expandHome(path), nil
	}
	cwd, _ := os.Getwd()
	if containing := FindWorkspacesContaining(cwd); len(containing) > 0 {
		return containing[0], nil
	}
	wsDir := GetWorkspacesDir()
	matches, _ := filepath.Glob(filepath.Join(wsDir, "*.code-workspace"))
	if len(matches) > 0 {
		sort.Strings(matches)
		return matches[0], nil
	}
	return filepath.Join(wsDir, "default.code-workspace"), nil
}

func GetCapabilitiesPath() string {
	if path := os.Getenv("WS_CAPABILITIES"); path != "" {
		return expandHome(path)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "workspace-cli", "capabilities.yaml")
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
