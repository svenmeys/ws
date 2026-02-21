package workspace

import (
	"strings"
	"testing"
)

func TestHandleHook(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantCwd       string
		wantIndicator string
		wantOk        bool
	}{
		{
			name:          "PreToolUse",
			input:         `{"hook_event_name": "PreToolUse", "cwd": "/projects/foo"}`,
			wantCwd:       "/projects/foo",
			wantIndicator: "working",
			wantOk:        true,
		},
		{
			name:          "SubagentStart",
			input:         `{"hook_event_name": "SubagentStart", "cwd": "/projects/bar"}`,
			wantCwd:       "/projects/bar",
			wantIndicator: "working",
			wantOk:        true,
		},
		{
			name:          "UserPromptSubmit",
			input:         `{"hook_event_name": "UserPromptSubmit", "cwd": "/home/user"}`,
			wantCwd:       "/home/user",
			wantIndicator: "working",
			wantOk:        true,
		},
		{
			name:          "Notification idle_prompt",
			input:         `{"hook_event_name": "Notification", "cwd": "/home/user", "matcher": "idle_prompt"}`,
			wantCwd:       "/home/user",
			wantIndicator: "waiting",
			wantOk:        true,
		},
		{
			name:          "Notification permission_prompt",
			input:         `{"hook_event_name": "Notification", "cwd": "/home/user", "matcher": "permission_prompt"}`,
			wantCwd:       "/home/user",
			wantIndicator: "waiting",
			wantOk:        true,
		},
		{
			name:          "Notification other",
			input:         `{"hook_event_name": "Notification", "cwd": "/home/user", "matcher": "something_else"}`,
			wantCwd:       "/home/user",
			wantIndicator: "idle",
			wantOk:        true,
		},
		{
			name:          "PermissionRequest",
			input:         `{"hook_event_name": "PermissionRequest", "cwd": "/projects/baz"}`,
			wantCwd:       "/projects/baz",
			wantIndicator: "waiting",
			wantOk:        true,
		},
		{
			name:          "Stop",
			input:         `{"hook_event_name": "Stop", "cwd": "/projects/foo"}`,
			wantCwd:       "/projects/foo",
			wantIndicator: "idle",
			wantOk:        true,
		},
		{
			name:          "PostToolUse",
			input:         `{"hook_event_name": "PostToolUse", "cwd": "/projects/foo"}`,
			wantCwd:       "/projects/foo",
			wantIndicator: "idle",
			wantOk:        true,
		},
		{
			name:          "SubagentStop",
			input:         `{"hook_event_name": "SubagentStop", "cwd": "/projects/foo"}`,
			wantCwd:       "/projects/foo",
			wantIndicator: "idle",
			wantOk:        true,
		},
		{
			name:          "SessionEnd",
			input:         `{"hook_event_name": "SessionEnd", "cwd": "/projects/foo"}`,
			wantCwd:       "/projects/foo",
			wantIndicator: "idle",
			wantOk:        true,
		},
		{
			name:   "unknown event",
			input:  `{"hook_event_name": "SomeNewEvent", "cwd": "/foo"}`,
			wantOk: false,
		},
		{
			name:   "invalid JSON",
			input:  `{not json`,
			wantOk: false,
		},
		{
			name:   "empty input",
			input:  ``,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwd, indicator, ok := HandleHook(strings.NewReader(tt.input))
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if !tt.wantOk {
				return
			}
			if cwd != tt.wantCwd {
				t.Errorf("cwd = %s, want %s", cwd, tt.wantCwd)
			}
			if indicator != tt.wantIndicator {
				t.Errorf("indicator = %s, want %s", indicator, tt.wantIndicator)
			}
		})
	}
}

func TestUpdateActivityIndicator(t *testing.T) {
	t.Run("add indicator", func(t *testing.T) {
		folder := map[string]any{"name": "🟢 My Project"}
		changed := UpdateActivityIndicator(folder, "working")
		if !changed {
			t.Error("expected change")
		}
		name := GetFolderField(folder, "name")
		if name != "🟢 My Project ⚡" {
			t.Errorf("expected '🟢 My Project ⚡', got '%s'", name)
		}
	})

	t.Run("replace indicator", func(t *testing.T) {
		folder := map[string]any{"name": "🟢 My Project ⚡"}
		changed := UpdateActivityIndicator(folder, "waiting")
		if !changed {
			t.Error("expected change")
		}
		name := GetFolderField(folder, "name")
		if name != "🟢 My Project ❓" {
			t.Errorf("expected '🟢 My Project ❓', got '%s'", name)
		}
	})

	t.Run("remove indicator (idle)", func(t *testing.T) {
		folder := map[string]any{"name": "🟢 My Project ⚡"}
		changed := UpdateActivityIndicator(folder, "idle")
		if !changed {
			t.Error("expected change")
		}
		name := GetFolderField(folder, "name")
		if name != "🟢 My Project" {
			t.Errorf("expected '🟢 My Project', got '%s'", name)
		}
	})

	t.Run("no change when already idle", func(t *testing.T) {
		folder := map[string]any{"name": "🟢 My Project"}
		changed := UpdateActivityIndicator(folder, "idle")
		if changed {
			t.Error("expected no change")
		}
	})

	t.Run("no change when same indicator", func(t *testing.T) {
		folder := map[string]any{"name": "🟢 My Project ⚡"}
		changed := UpdateActivityIndicator(folder, "working")
		if changed {
			t.Error("expected no change when same indicator")
		}
	})
}

func TestIndicators(t *testing.T) {
	expected := map[string]string{
		"working": "⚡",
		"waiting": "❓",
		"idle":    "",
	}
	for key, emoji := range expected {
		if Indicators[key] != emoji {
			t.Errorf("Indicators[%s] = %s, want %s", key, Indicators[key], emoji)
		}
	}
	if len(Indicators) != len(expected) {
		t.Errorf("Indicators has %d entries, expected %d", len(Indicators), len(expected))
	}
}
