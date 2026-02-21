package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

// testWorkspace is a minimal JSONC workspace for testing.
const testWorkspace = `{
	// Test workspace
	"folders": [
		{
			"path": "./project-a",
			"name": "🟢 Project A",
			"x-pa": {
				"slack_channel": "C001",
				"description": "First project"
			}
		},
		{
			"path": "./project-b",
			"name": "🔵 Project B",
			"x-pa": {
				"slack_channel": "C002"
			}
		},
		{
			"path": "./project-c",
			"name": "Project C"
		}
	],
	"settings": {}
}
`

func writeTestWorkspace(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "test.code-workspace")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadWorkspace(t *testing.T) {
	dir := t.TempDir()

	t.Run("valid JSONC", func(t *testing.T) {
		path := writeTestWorkspace(t, dir, testWorkspace)
		ws, err := ReadWorkspace(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ws == nil {
			t.Fatal("expected workspace, got nil")
		}
		folders := GetFolders(ws)
		if len(folders) != 3 {
			t.Errorf("expected 3 folders, got %d", len(folders))
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		path := filepath.Join(dir, "bad.code-workspace")
		os.WriteFile(path, []byte(`{not valid`), 0644)
		_, err := ReadWorkspace(path)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := ReadWorkspace(filepath.Join(dir, "nope.code-workspace"))
		if err == nil {
			t.Fatal("expected error for missing file")
		}
	})

	t.Run("empty folders", func(t *testing.T) {
		emptyDir := t.TempDir()
		path := filepath.Join(emptyDir, "test.code-workspace")
		os.WriteFile(path, []byte(`{"folders": []}`), 0644)
		ws, err := ReadWorkspace(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		folders := GetFolders(ws)
		if len(folders) != 0 {
			t.Errorf("expected 0 folders, got %d", len(folders))
		}
	})
}

func TestGetFolders(t *testing.T) {
	t.Run("nil workspace", func(t *testing.T) {
		folders := GetFolders(nil)
		if folders != nil {
			t.Errorf("expected nil, got %v", folders)
		}
	})

	t.Run("no folders key", func(t *testing.T) {
		ws := map[string]any{"settings": map[string]any{}}
		folders := GetFolders(ws)
		if folders != nil {
			t.Errorf("expected nil, got %v", folders)
		}
	})

	t.Run("wrong type folders", func(t *testing.T) {
		ws := map[string]any{"folders": "not an array"}
		folders := GetFolders(ws)
		if folders != nil {
			t.Errorf("expected nil, got %v", folders)
		}
	})
}

func loadTestWS(t *testing.T) map[string]any {
	t.Helper()
	dir := t.TempDir()
	path := writeTestWorkspace(t, dir, testWorkspace)
	ws, err := ReadWorkspace(path)
	if err != nil {
		t.Fatal(err)
	}
	return ws
}

func TestFindFolder(t *testing.T) {
	ws := loadTestWS(t)

	t.Run("by path", func(t *testing.T) {
		f := FindFolder(ws, "project-a")
		if f == nil {
			t.Fatal("expected folder, got nil")
		}
		if GetFolderField(f, "path") != "./project-a" {
			t.Errorf("expected ./project-a, got %s", GetFolderField(f, "path"))
		}
	})

	t.Run("by name partial match", func(t *testing.T) {
		f := FindFolder(ws, "Project B")
		if f == nil {
			t.Fatal("expected folder, got nil")
		}
		if GetFolderField(f, "path") != "./project-b" {
			t.Errorf("expected ./project-b, got %s", GetFolderField(f, "path"))
		}
	})

	t.Run("case insensitive name", func(t *testing.T) {
		f := FindFolder(ws, "project c")
		if f == nil {
			t.Fatal("expected folder, got nil")
		}
		if GetFolderField(f, "path") != "./project-c" {
			t.Errorf("expected ./project-c, got %s", GetFolderField(f, "path"))
		}
	})

	t.Run("not found", func(t *testing.T) {
		f := FindFolder(ws, "nonexistent")
		if f != nil {
			t.Errorf("expected nil, got %v", f)
		}
	})

	t.Run("empty workspace", func(t *testing.T) {
		empty := map[string]any{"folders": []any{}}
		f := FindFolder(empty, "anything")
		if f != nil {
			t.Errorf("expected nil, got %v", f)
		}
	})
}

func TestGetChannelMappings(t *testing.T) {
	ws := loadTestWS(t)

	t.Run("extracts channels", func(t *testing.T) {
		mappings := GetChannelMappings(ws)
		if len(mappings) != 2 {
			t.Errorf("expected 2 mappings, got %d", len(mappings))
		}
		if mappings["C001"] != "./project-a" {
			t.Errorf("expected ./project-a for C001, got %s", mappings["C001"])
		}
		if mappings["C002"] != "./project-b" {
			t.Errorf("expected ./project-b for C002, got %s", mappings["C002"])
		}
	})

	t.Run("skips folders without channel", func(t *testing.T) {
		mappings := GetChannelMappings(ws)
		for ch := range mappings {
			if ch == "" {
				t.Error("found empty channel key")
			}
		}
	})

	t.Run("empty workspace", func(t *testing.T) {
		empty := map[string]any{"folders": []any{}}
		mappings := GetChannelMappings(empty)
		if len(mappings) != 0 {
			t.Errorf("expected 0 mappings, got %d", len(mappings))
		}
	})
}

func TestAddFolder(t *testing.T) {
	t.Run("with all fields", func(t *testing.T) {
		ws := map[string]any{"folders": []any{}}
		folder := AddFolder(ws, "./new", "New Project", "C999", "A new one")

		if GetFolderField(folder, "path") != "./new" {
			t.Errorf("expected ./new, got %s", GetFolderField(folder, "path"))
		}
		if GetFolderField(folder, "name") != "New Project" {
			t.Errorf("expected New Project, got %s", GetFolderField(folder, "name"))
		}
		if GetXPAField(folder, "slack_channel") != "C999" {
			t.Errorf("expected C999, got %s", GetXPAField(folder, "slack_channel"))
		}
		if GetXPAField(folder, "description") != "A new one" {
			t.Errorf("expected A new one, got %s", GetXPAField(folder, "description"))
		}

		folders := GetFolders(ws)
		if len(folders) != 1 {
			t.Errorf("expected 1 folder, got %d", len(folders))
		}
	})

	t.Run("minimal fields", func(t *testing.T) {
		ws := map[string]any{"folders": []any{}}
		folder := AddFolder(ws, "./min", "Minimal", "", "")

		if folder["x-pa"] != nil {
			t.Error("expected no x-pa for minimal folder")
		}
	})

	t.Run("channel only", func(t *testing.T) {
		ws := map[string]any{"folders": []any{}}
		folder := AddFolder(ws, "./ch", "WithChannel", "C123", "")

		if GetXPAField(folder, "slack_channel") != "C123" {
			t.Errorf("expected C123, got %s", GetXPAField(folder, "slack_channel"))
		}
	})

	t.Run("appends to existing", func(t *testing.T) {
		ws := map[string]any{"folders": []any{
			map[string]any{"path": "./existing", "name": "Existing"},
		}}
		AddFolder(ws, "./new", "New", "", "")

		folders := GetFolders(ws)
		if len(folders) != 2 {
			t.Errorf("expected 2 folders, got %d", len(folders))
		}
	})
}

func TestUpdateFolderStatus(t *testing.T) {
	t.Run("replace existing status", func(t *testing.T) {
		ws := loadTestWS(t)
		ok, err := UpdateFolderStatus(ws, "project-a", "blocked")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected update to succeed")
		}
		folder := FindFolder(ws, "project-a")
		name := GetFolderField(folder, "name")
		if name != "🔴 Project A" {
			t.Errorf("expected '🔴 Project A', got '%s'", name)
		}
	})

	t.Run("add status to folder without one", func(t *testing.T) {
		ws := loadTestWS(t)
		ok, err := UpdateFolderStatus(ws, "project-c", "active")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected update to succeed")
		}
		folder := FindFolder(ws, "project-c")
		name := GetFolderField(folder, "name")
		if name != "🟢 Project C" {
			t.Errorf("expected '🟢 Project C', got '%s'", name)
		}
	})

	t.Run("unknown status", func(t *testing.T) {
		ws := loadTestWS(t)
		_, err := UpdateFolderStatus(ws, "project-a", "invalid")
		if err == nil {
			t.Fatal("expected error for unknown status")
		}
	})

	t.Run("project not found", func(t *testing.T) {
		ws := loadTestWS(t)
		ok, err := UpdateFolderStatus(ws, "nonexistent", "active")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Error("expected false for missing project")
		}
	})

	t.Run("all status types", func(t *testing.T) {
		for status, emoji := range StatusEmojis {
			ws := loadTestWS(t)
			ok, err := UpdateFolderStatus(ws, "project-c", status)
			if err != nil {
				t.Fatalf("unexpected error for status %s: %v", status, err)
			}
			if !ok {
				t.Fatalf("expected update to succeed for status %s", status)
			}
			folder := FindFolder(ws, "project-c")
			name := GetFolderField(folder, "name")
			if name != emoji+" Project C" {
				t.Errorf("status %s: expected '%s Project C', got '%s'", status, emoji, name)
			}
		}
	})
}

func TestWriteWorkspace(t *testing.T) {
	dir := t.TempDir()

	t.Run("roundtrip", func(t *testing.T) {
		path := writeTestWorkspace(t, dir, testWorkspace)
		ws, err := ReadWorkspace(path)
		if err != nil {
			t.Fatal(err)
		}

		outPath := filepath.Join(dir, "out.code-workspace")
		// Write a file first so backup logic triggers
		os.WriteFile(outPath, []byte("old"), 0644)
		if err := WriteWorkspace(outPath, ws); err != nil {
			t.Fatal(err)
		}

		// Read back
		ws2, err := ReadWorkspace(outPath)
		if err != nil {
			t.Fatalf("failed to read written workspace: %v", err)
		}
		folders := GetFolders(ws2)
		if len(folders) != 3 {
			t.Errorf("expected 3 folders after roundtrip, got %d", len(folders))
		}
	})

	t.Run("creates backup", func(t *testing.T) {
		path := filepath.Join(dir, "backup-test.code-workspace")
		os.WriteFile(path, []byte("original"), 0644)

		ws := map[string]any{"folders": []any{}}
		if err := WriteWorkspace(path, ws); err != nil {
			t.Fatal(err)
		}

		backup, err := os.ReadFile(path + ".backup")
		if err != nil {
			t.Fatalf("expected backup file: %v", err)
		}
		if string(backup) != "original" {
			t.Errorf("backup content mismatch: %s", string(backup))
		}
	})
}

func TestGetFolderField(t *testing.T) {
	folder := map[string]any{
		"path": "./test",
		"name": "Test",
		"num":  42,
	}

	t.Run("existing string", func(t *testing.T) {
		if v := GetFolderField(folder, "path"); v != "./test" {
			t.Errorf("expected ./test, got %s", v)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		if v := GetFolderField(folder, "missing"); v != "" {
			t.Errorf("expected empty string, got %s", v)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		if v := GetFolderField(folder, "num"); v != "" {
			t.Errorf("expected empty string for non-string, got %s", v)
		}
	})
}

func TestGetXPAField(t *testing.T) {
	t.Run("with x-pa", func(t *testing.T) {
		folder := map[string]any{
			"x-pa": map[string]any{
				"slack_channel": "C123",
			},
		}
		if v := GetXPAField(folder, "slack_channel"); v != "C123" {
			t.Errorf("expected C123, got %s", v)
		}
	})

	t.Run("without x-pa", func(t *testing.T) {
		folder := map[string]any{"path": "./test"}
		if v := GetXPAField(folder, "slack_channel"); v != "" {
			t.Errorf("expected empty, got %s", v)
		}
	})

	t.Run("missing field in x-pa", func(t *testing.T) {
		folder := map[string]any{
			"x-pa": map[string]any{},
		}
		if v := GetXPAField(folder, "slack_channel"); v != "" {
			t.Errorf("expected empty, got %s", v)
		}
	})
}

func TestStatusEmojis(t *testing.T) {
	expected := map[string]string{
		"active":   "🟢",
		"paused":   "🟡",
		"blocked":  "🔴",
		"progress": "🔵",
		"dormant":  "⚪",
	}
	for status, emoji := range expected {
		if StatusEmojis[status] != emoji {
			t.Errorf("StatusEmojis[%s] = %s, want %s", status, StatusEmojis[status], emoji)
		}
	}
	if len(StatusEmojis) != len(expected) {
		t.Errorf("StatusEmojis has %d entries, expected %d", len(StatusEmojis), len(expected))
	}
}
