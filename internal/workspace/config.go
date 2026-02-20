package workspace

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GetWorkspacesDir returns the directory containing workspace files.
func GetWorkspacesDir() string {
	if dir := os.Getenv("WS_WORKSPACES_DIR"); dir != "" {
		return expandHome(dir)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "workspaces")
}

// FindWorkspacesContaining finds all workspace files containing the given directory.
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

// GetDefaultWorkspace returns the workspace path, auto-detecting from cwd if possible.
func GetDefaultWorkspace() (string, error) {
	// 1. Explicit env override
	if path := os.Getenv("WS_WORKSPACE"); path != "" {
		return expandHome(path), nil
	}

	// 2. Auto-detect from current directory
	cwd, _ := os.Getwd()
	if containing := FindWorkspacesContaining(cwd); len(containing) > 0 {
		return containing[0], nil
	}

	// 3. Fallback to first workspace file
	wsDir := GetWorkspacesDir()
	matches, _ := filepath.Glob(filepath.Join(wsDir, "*.code-workspace"))
	if len(matches) > 0 {
		sort.Strings(matches)
		return matches[0], nil
	}

	return filepath.Join(wsDir, "default.code-workspace"), nil
}

// GetCapabilitiesPath returns the capabilities.yaml path.
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
