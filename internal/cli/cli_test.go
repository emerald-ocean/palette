package cli

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testTemplate = `bg = "{{.abyss1}}" -- palette {{version}}`

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

	photo := filepath.Join(dir, "photo.png")
	writeTestPNG(t, photo)
	recoloured := filepath.Join(dir, "recoloured.png")

	tests := []struct {
		name       string
		args       []string
		wantStdout string
		wantErr    string
		wantRules  bool
	}{
		{name: "no command", args: nil, wantErr: "missing command"},
		{name: "unknown command", args: []string{"paint"}, wantErr: `unknown command "paint"`},
		{name: "check embedded", args: []string{"check"}, wantStdout: "all 44 colours meet the rules"},
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
		{name: "recolour without flags", args: []string{"recolour"}, wantErr: "needs -in and -out"},
		{
			name:    "recolour unknown colour",
			args:    []string{"recolour", "-in", photo, "-out", recoloured, "-ramp", "abyss1,trench"},
			wantErr: `unknown colour "trench"`,
		},
		{
			name:    "recolour reversed input",
			args:    []string{"recolour", "-in", photo, "-out", recoloured, "-dark", "0.9"},
			wantErr: "must be below light",
		},
		{name: "recolour", args: []string{"recolour", "-in", photo, "-out", recoloured}},
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

func TestRecolourOutput(t *testing.T) {
	dir := t.TempDir()
	photo := filepath.Join(dir, "photo.png")
	writeTestPNG(t, photo)
	out := filepath.Join(dir, "recoloured.png")
	if err := Run([]string{"recolour", "-in", photo, "-out", out}, nil); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("png.Decode() error = %v", err)
	}

	tests := []struct {
		name     string
		x        int
		expected color.NRGBA
	}{
		{name: "black becomes abyss1", x: 0, expected: color.NRGBA{R: 0x01, G: 0x13, B: 0x0b, A: 255}},
		{name: "white becomes foam5", x: 1, expected: color.NRGBA{R: 0xdb, G: 0xea, B: 0xeb, A: 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := color.NRGBAModel.Convert(img.At(tt.x, 0)); got != tt.expected {
				t.Errorf("pixel (%d, 0) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func writeTestPNG(
	t *testing.T,
	path string,
) {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.SetNRGBA(0, 0, color.NRGBA{A: 255})
	img.SetNRGBA(1, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatalf("png.Encode() error = %v", err)
	}
	if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
