package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const syncTestWorkspace = `{
	"folders": [
		{
			"path": "./alpha",
			"name": "Alpha",
			"x-pa": {"slack_channel": "C100"}
		},
		{
			"path": "./beta",
			"name": "Beta",
			"x-pa": {"slack_channel": "C200"}
		},
		{
			"path": "./gamma",
			"name": "Gamma"
		}
	]
}
`

func TestSyncToCapabilities(t *testing.T) {
	t.Run("dry run returns mappings without writing", func(t *testing.T) {
		dir := t.TempDir()
		wsPath := filepath.Join(dir, "test.code-workspace")
		os.WriteFile(wsPath, []byte(syncTestWorkspace), 0644)
		capsPath := filepath.Join(dir, "caps.yaml")

		mappings, err := SyncToCapabilities(wsPath, capsPath, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(mappings) != 2 {
			t.Errorf("expected 2 mappings, got %d", len(mappings))
		}

		// Caps file should not exist after dry run
		if _, err := os.Stat(capsPath); !os.IsNotExist(err) {
			t.Error("capabilities file should not exist after dry run")
		}
	})

	t.Run("actual sync writes capabilities", func(t *testing.T) {
		dir := t.TempDir()
		wsPath := filepath.Join(dir, "test.code-workspace")
		os.WriteFile(wsPath, []byte(syncTestWorkspace), 0644)
		capsPath := filepath.Join(dir, "config", "caps.yaml")

		mappings, err := SyncToCapabilities(wsPath, capsPath, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(mappings) != 2 {
			t.Errorf("expected 2 mappings, got %d", len(mappings))
		}

		// Verify file was written
		data, err := os.ReadFile(capsPath)
		if err != nil {
			t.Fatalf("failed to read capabilities: %v", err)
		}

		var caps map[string]any
		if err := yaml.Unmarshal(data, &caps); err != nil {
			t.Fatalf("failed to parse capabilities: %v", err)
		}

		slack, ok := caps["slack"].(map[string]any)
		if !ok {
			t.Fatal("expected slack section in capabilities")
		}
		channelProjects, ok := slack["channel_projects"].(map[string]any)
		if !ok {
			t.Fatal("expected channel_projects in slack section")
		}
		if len(channelProjects) != 2 {
			t.Errorf("expected 2 channel_projects, got %d", len(channelProjects))
		}
	})

	t.Run("preserves existing capabilities", func(t *testing.T) {
		dir := t.TempDir()
		wsPath := filepath.Join(dir, "test.code-workspace")
		os.WriteFile(wsPath, []byte(syncTestWorkspace), 0644)
		capsPath := filepath.Join(dir, "caps.yaml")

		// Write existing capabilities
		existing := map[string]any{
			"other_key": "preserved",
			"slack":     map[string]any{"webhook": "https://example.com"},
		}
		data, _ := yaml.Marshal(existing)
		os.WriteFile(capsPath, data, 0644)

		_, err := SyncToCapabilities(wsPath, capsPath, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		result, _ := os.ReadFile(capsPath)
		var caps map[string]any
		yaml.Unmarshal(result, &caps)

		if caps["other_key"] != "preserved" {
			t.Error("expected other_key to be preserved")
		}
	})

	t.Run("empty channel mappings", func(t *testing.T) {
		dir := t.TempDir()
		wsPath := filepath.Join(dir, "test.code-workspace")
		os.WriteFile(wsPath, []byte(`{"folders": [{"path": "./no-channel", "name": "None"}]}`), 0644)
		capsPath := filepath.Join(dir, "caps.yaml")

		mappings, err := SyncToCapabilities(wsPath, capsPath, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(mappings) != 0 {
			t.Errorf("expected 0 mappings, got %d", len(mappings))
		}
	})

	t.Run("paths resolve to home-relative", func(t *testing.T) {
		dir := t.TempDir()
		wsPath := filepath.Join(dir, "test.code-workspace")
		os.WriteFile(wsPath, []byte(syncTestWorkspace), 0644)
		capsPath := filepath.Join(dir, "caps.yaml")

		mappings, err := SyncToCapabilities(wsPath, capsPath, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		home, _ := os.UserHomeDir()
		for _, path := range mappings {
			// Path should be either home-relative (~/) or absolute
			if !strings.HasPrefix(path, "~/") && !filepath.IsAbs(path) {
				t.Errorf("expected home-relative or absolute path, got %s", path)
			}
			// Should not be a relative path
			if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
				t.Errorf("expected resolved path, got relative: %s", path)
			}
			_ = home // used above in prefix check
		}
	})

	t.Run("missing workspace file", func(t *testing.T) {
		_, err := SyncToCapabilities("/nonexistent/ws.code-workspace", "/tmp/caps.yaml", false)
		if err == nil {
			t.Fatal("expected error for missing workspace")
		}
	})
}
