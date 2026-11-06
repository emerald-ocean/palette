package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testTemplate = `bg = "{{.abyss}}" -- palette {{version}}`

func TestRun(t *testing.T) {
	dir := t.TempDir()
	tmpl := filepath.Join(dir, "palette.lua.tmpl")
	out := filepath.Join(dir, "palette.lua")
	if err := os.WriteFile(tmpl, []byte(testTemplate), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	broken := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(broken, []byte(`{
		"space": "oklch",
		"groups": {"abyss": {"abyss": {"oklch": [0.165, 0.033, 166], "hex": "#000000"}}}
	}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name       string
		args       []string
		wantStdout string
		wantErr    string
		wantRules  bool
	}{
		{name: "no command", args: nil, wantErr: "missing command"},
		{name: "unknown command", args: []string{"paint"}, wantErr: `unknown command "paint"`},
		{name: "check embedded", args: []string{"check"}, wantStdout: "all 15 colours meet the rules"},
		{
			name:       "check broken",
			args:       []string{"check", "-palette", broken},
			wantStdout: "abyss: hex #000000 does not match",
			wantErr:    "1 problems",
			wantRules:  true,
		},
		{name: "report", args: []string{"report"}, wantStdout: "deuteranopia"},
		{name: "generate without flags", args: []string{"generate"}, wantErr: "needs -template and -out"},
		{name: "generate", args: []string{"generate", "-template", tmpl, "-out", out}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout strings.Builder
			err := Run(tt.args, &stdout)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Errorf("Run() error = %v, want it to contain %q", err, tt.wantErr)
			}
			if got := errors.Is(err, errRules); got != tt.wantRules {
				t.Errorf("errors.Is(Run(), errRules) = %v, want %v", got, tt.wantRules)
			}
			if !strings.Contains(stdout.String(), tt.wantStdout) {
				t.Errorf("Run() stdout = %q, want it to contain %q", stdout.String(), tt.wantStdout)
			}
		})
	}

	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if expected := `bg = "#01130b" -- palette ` + version(); string(got) != expected {
		t.Errorf("generated file = %q, want %q", got, expected)
	}
}
