package palette

import (
	"strings"
	"testing"
)

const testVersion = "v9.9.9"

func TestRender(t *testing.T) {
	p, err := Default()
	if err != nil {
		t.Fatalf("Default() error = %v", err)
	}
	p.Version = testVersion
	tests := []struct {
		name     string
		text     string
		expected string
		wantErr  string
	}{
		{name: "colour", text: `fg = "{{.emerald}}"`, expected: `fg = "#83eab5"`},
		{name: "alpha", text: `{{alpha .glow 0.16}}`, expected: "#76dffb29"},
		{name: "opaque", text: `{{alpha .glow 1.0}}`, expected: "#76dffbff"},
		{name: "version", text: `palette {{version}}`, expected: "palette " + testVersion},
		{name: "unknown colour", text: `{{.seaweed}}`, wantErr: `map has no entry for key "seaweed"`},
		{name: "alpha out of range", text: `{{alpha .glow 1.5}}`, wantErr: "want a value from 0 to 1"},
		{name: "alpha bad hex", text: `{{alpha "green" 0.5}}`, wantErr: `parse hex "green"`},
		{name: "bad syntax", text: `{{.glow`, wantErr: "parse template test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			err := p.Render(&b, "test", tt.text)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Render() error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			if got := b.String(); got != tt.expected {
				t.Errorf("Render() = %q, want %q", got, tt.expected)
			}
		})
	}
}
