package workspace

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
)

// Indicators maps activity states to their emoji indicators.
var Indicators = map[string]string{
	"working": "⚡",
	"waiting": "❓",
	"idle":    "",
}

// GetProjectFromCwd finds the workspace and project matching the current working directory.
func GetProjectFromCwd(cwd string) (string, map[string]any) {
	wsPath, err := GetDefaultWorkspace()
	if err != nil {
		return "", nil
	}

	ws, err := ReadWorkspace(wsPath)
	if err != nil {
		return "", nil
	}

	absPath, _ := filepath.Abs(cwd)
	for _, folder := range GetFolders(ws) {
		if path, ok := folder["path"].(string); ok {
			folderAbs, _ := filepath.Abs(filepath.Join(filepath.Dir(wsPath), path))
			if strings.HasPrefix(absPath, folderAbs) {
				return wsPath, folder
			}
		}
	}

	return wsPath, nil
}

// UpdateActivityIndicator updates the activity indicator in a folder's name.
func UpdateActivityIndicator(folder map[string]any, indicatorType string) bool {
	name, _ := folder["name"].(string)
	original := name

	for _, ind := range Indicators {
		if ind != "" && strings.Contains(name, ind) {
			name = strings.TrimSpace(strings.ReplaceAll(name, ind, ""))
		}
	}

	indicator := Indicators[indicatorType]
	if indicator != "" {
		folder["name"] = name + " " + indicator
	} else {
		folder["name"] = name
	}

	return folder["name"] != original
}

// HandleHook reads hook event data from stdin and returns the cwd and indicator type.
func HandleHook(stdin io.Reader) (cwd string, indicatorType string, ok bool) {
	var hookData map[string]any
	if err := json.NewDecoder(stdin).Decode(&hookData); err != nil {
		return "", "", false
	}

	event, _ := hookData["hook_event_name"].(string)
	cwd, _ = hookData["cwd"].(string)

	switch event {
	case "PreToolUse", "SubagentStart", "UserPromptSubmit":
		return cwd, "working", true
	case "Notification":
		matcher, _ := hookData["matcher"].(string)
		if matcher == "idle_prompt" || matcher == "permission_prompt" {
			return cwd, "waiting", true
		}
		return cwd, "idle", true
	case "PermissionRequest":
		return cwd, "waiting", true
	case "Stop", "PostToolUse", "SubagentStop", "SessionEnd":
		return cwd, "idle", true
	default:
		return "", "", false
	}
}
