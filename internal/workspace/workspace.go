package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/tailscale/hujson"
)

// StatusEmojis maps status names to their emoji indicators.
var StatusEmojis = map[string]string{
	"active":   "🟢",
	"paused":   "🟡",
	"blocked":  "🔴",
	"progress": "🔵",
	"dormant":  "⚪",
}

// ReadWorkspace reads a JSONC workspace file.
func ReadWorkspace(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("workspace file not found: %s", path)
	}

	standardized, err := hujson.Standardize(data)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace file: %w", err)
	}

	var ws map[string]any
	if err := json.Unmarshal(standardized, &ws); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return ws, nil
}

// WriteWorkspace writes workspace data with backup and atomic replace.
func WriteWorkspace(path string, data map[string]any) error {
	// Create backup
	if existing, err := os.ReadFile(path); err == nil {
		_ = os.WriteFile(path+".backup", existing, 0644)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode workspace: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0644); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

// GetFolders returns all folder entries from a workspace.
func GetFolders(ws map[string]any) []map[string]any {
	folders, ok := ws["folders"].([]any)
	if !ok {
		return nil
	}
	result := make([]map[string]any, 0, len(folders))
	for _, f := range folders {
		if m, ok := f.(map[string]any); ok {
			result = append(result, m)
		}
	}
	return result
}

// FindFolder finds a folder by name (partial match) or path.
func FindFolder(ws map[string]any, nameOrPath string) map[string]any {
	for _, folder := range GetFolders(ws) {
		if path, ok := folder["path"].(string); ok {
			if strings.Contains(path, nameOrPath) {
				return folder
			}
		}
		if name, ok := folder["name"].(string); ok {
			clean := strings.Map(func(r rune) rune {
				if r < 128 || unicode.IsLetter(r) || unicode.IsDigit(r) {
					return r
				}
				return -1
			}, name)
			if strings.Contains(strings.ToLower(clean), strings.ToLower(nameOrPath)) {
				return folder
			}
		}
	}
	return nil
}

// GetChannelMappings extracts channel -> path mappings from workspace.
func GetChannelMappings(ws map[string]any) map[string]string {
	mappings := make(map[string]string)
	for _, folder := range GetFolders(ws) {
		xpa := getXPA(folder)
		if channel, ok := xpa["slack_channel"].(string); ok && channel != "" {
			if path, ok := folder["path"].(string); ok {
				mappings[channel] = path
			}
		}
	}
	return mappings
}

// AddFolder adds a new folder entry to workspace.
func AddFolder(ws map[string]any, path, name, slackChannel, description string) map[string]any {
	folder := map[string]any{
		"path": path,
		"name": name,
	}

	if slackChannel != "" || description != "" {
		xpa := map[string]any{}
		if slackChannel != "" {
			xpa["slack_channel"] = slackChannel
		}
		if description != "" {
			xpa["description"] = description
		}
		folder["x-pa"] = xpa
	}

	folders, _ := ws["folders"].([]any)
	ws["folders"] = append(folders, folder)
	return folder
}

// UpdateFolderStatus updates the status emoji in a folder's name.
func UpdateFolderStatus(ws map[string]any, nameOrPath, status string) (bool, error) {
	emoji, ok := StatusEmojis[strings.ToLower(status)]
	if !ok {
		return false, fmt.Errorf("unknown status: %s (use: active, paused, blocked, progress, dormant)", status)
	}

	folder := FindFolder(ws, nameOrPath)
	if folder == nil {
		return false, nil
	}

	name, _ := folder["name"].(string)
	for _, oldEmoji := range StatusEmojis {
		if strings.Contains(name, oldEmoji) {
			folder["name"] = strings.Replace(name, oldEmoji, emoji, 1)
			return true, nil
		}
	}

	folder["name"] = emoji + " " + name
	return true, nil
}

func getXPA(folder map[string]any) map[string]any {
	if xpa, ok := folder["x-pa"].(map[string]any); ok {
		return xpa
	}
	return map[string]any{}
}

// GetFolderField safely gets a string field from a folder.
func GetFolderField(folder map[string]any, key string) string {
	if v, ok := folder[key].(string); ok {
		return v
	}
	return ""
}

// GetXPAField safely gets a string field from a folder's x-pa section.
func GetXPAField(folder map[string]any, key string) string {
	return GetFolderField(getXPA(folder), key)
}
