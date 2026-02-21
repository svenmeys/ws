package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandHome(t *testing.T) {
	home, _ := os.UserHomeDir()

	t.Run("with tilde prefix", func(t *testing.T) {
		result := expandHome("~/foo/bar")
		expected := filepath.Join(home, "foo/bar")
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("without tilde", func(t *testing.T) {
		result := expandHome("/absolute/path")
		if result != "/absolute/path" {
			t.Errorf("expected /absolute/path, got %s", result)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		result := expandHome("relative/path")
		if result != "relative/path" {
			t.Errorf("expected relative/path, got %s", result)
		}
	})
}

func TestGetWorkspacesDir(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		os.Unsetenv("WS_WORKSPACES_DIR")
		dir := GetWorkspacesDir()
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, "workspaces")
		if dir != expected {
			t.Errorf("expected %s, got %s", expected, dir)
		}
	})

	t.Run("env override", func(t *testing.T) {
		t.Setenv("WS_WORKSPACES_DIR", "/custom/dir")
		dir := GetWorkspacesDir()
		if dir != "/custom/dir" {
			t.Errorf("expected /custom/dir, got %s", dir)
		}
	})

	t.Run("env with tilde", func(t *testing.T) {
		t.Setenv("WS_WORKSPACES_DIR", "~/my-workspaces")
		dir := GetWorkspacesDir()
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, "my-workspaces")
		if dir != expected {
			t.Errorf("expected %s, got %s", expected, dir)
		}
	})
}

func TestGetCapabilitiesPath(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		os.Unsetenv("WS_CAPABILITIES")
		path := GetCapabilitiesPath()
		if !strings.HasSuffix(path, filepath.Join(".config", "workspace-cli", "capabilities.yaml")) {
			t.Errorf("unexpected default path: %s", path)
		}
	})

	t.Run("env override", func(t *testing.T) {
		t.Setenv("WS_CAPABILITIES", "/custom/caps.yaml")
		path := GetCapabilitiesPath()
		if path != "/custom/caps.yaml" {
			t.Errorf("expected /custom/caps.yaml, got %s", path)
		}
	})
}

func TestGetDefaultWorkspace(t *testing.T) {
	t.Run("explicit env", func(t *testing.T) {
		t.Setenv("WS_WORKSPACE", "/explicit/test.code-workspace")
		path, err := GetDefaultWorkspace()
		if err != nil {
			t.Fatal(err)
		}
		if path != "/explicit/test.code-workspace" {
			t.Errorf("expected /explicit/test.code-workspace, got %s", path)
		}
	})

	t.Run("fallback to first workspace in dir", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("WS_WORKSPACE", "")
		t.Setenv("WS_WORKSPACES_DIR", dir)

		// Create workspace files
		os.WriteFile(filepath.Join(dir, "aaa.code-workspace"), []byte(`{"folders":[]}`), 0644)
		os.WriteFile(filepath.Join(dir, "zzz.code-workspace"), []byte(`{"folders":[]}`), 0644)

		// Change to a dir not in any workspace
		origDir, _ := os.Getwd()
		tmpCwd := t.TempDir()
		os.Chdir(tmpCwd)
		defer os.Chdir(origDir)

		path, err := GetDefaultWorkspace()
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasSuffix(path, "aaa.code-workspace") {
			t.Errorf("expected first workspace (aaa), got %s", path)
		}
	})
}

func TestFindWorkspacesContaining(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("WS_WORKSPACES_DIR", dir)

	// Create a project directory inside the workspace dir (relative path ./myproject)
	projectDir := filepath.Join(dir, "myproject")
	os.MkdirAll(projectDir, 0755)

	// Create workspace that references the project with a relative path
	// FindWorkspacesContaining resolves paths relative to the workspace file's directory
	wsContent := `{"folders": [{"path": "./myproject"}]}`
	os.WriteFile(filepath.Join(dir, "test.code-workspace"), []byte(wsContent), 0644)

	t.Run("finds containing workspace", func(t *testing.T) {
		result := FindWorkspacesContaining(projectDir)
		if len(result) == 0 {
			t.Fatal("expected to find containing workspace")
		}
		if !strings.HasSuffix(result[0], "test.code-workspace") {
			t.Errorf("expected test.code-workspace, got %s", result[0])
		}
	})

	t.Run("no match", func(t *testing.T) {
		result := FindWorkspacesContaining("/nonexistent/path")
		if len(result) != 0 {
			t.Errorf("expected no results, got %v", result)
		}
	})

	t.Run("empty workspaces dir", func(t *testing.T) {
		emptyDir := t.TempDir()
		t.Setenv("WS_WORKSPACES_DIR", emptyDir)
		result := FindWorkspacesContaining(projectDir)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}
