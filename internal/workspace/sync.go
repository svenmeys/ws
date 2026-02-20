package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// SyncToCapabilities syncs channel mappings from workspace to capabilities.yaml.
func SyncToCapabilities(wsPath, capsPath string, dryRun bool) (map[string]string, error) {
	if capsPath == "" {
		capsPath = GetCapabilitiesPath()
	}
	if wsPath == "" {
		var err error
		wsPath, err = GetDefaultWorkspace()
		if err != nil {
			return nil, err
		}
	}

	ws, err := ReadWorkspace(wsPath)
	if err != nil {
		return nil, err
	}

	mappings := GetChannelMappings(ws)
	if len(mappings) == 0 {
		return map[string]string{}, nil
	}

	// Resolve paths relative to workspace location
	wsDir := filepath.Dir(wsPath)
	home, _ := os.UserHomeDir()
	resolved := make(map[string]string)
	for channel, relPath := range mappings {
		absPath, _ := filepath.Abs(filepath.Join(wsDir, relPath))
		if rel, err := filepath.Rel(home, absPath); err == nil {
			resolved[channel] = "~/" + rel
		} else {
			resolved[channel] = absPath
		}
	}

	if dryRun {
		return resolved, nil
	}

	// Load existing capabilities
	caps := make(map[string]any)
	if data, err := os.ReadFile(capsPath); err == nil {
		_ = yaml.Unmarshal(data, &caps)
	}

	// Update channel_projects section
	slack, ok := caps["slack"].(map[string]any)
	if !ok {
		slack = make(map[string]any)
		caps["slack"] = slack
	}
	slack["channel_projects"] = resolved

	output, err := yaml.Marshal(caps)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	_ = os.MkdirAll(filepath.Dir(capsPath), 0755)

	tmp := capsPath + ".tmp"
	if err := os.WriteFile(tmp, output, 0644); err != nil {
		return nil, err
	}

	return resolved, os.Rename(tmp, capsPath)
}
