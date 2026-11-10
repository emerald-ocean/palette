package palette

import (
	"strings"
	"testing"
)

const testRules = `"rules": {
	"minContrast": 4.5,
	"minDeltaE": 0.10,
	"backgrounds": ["abyss"],
	"neighbours": [["white", "glow"]]
}`

func testDocument(glow string) string {
	return `{
	"name": "Test",
	"space": "oklch",
	"groups": {
		"abyss": {"abyss": {"oklch": [0.165, 0.033, 166], "hex": "#01130b"}},
		"foam": {"white": {"oklch": [0.95, 0.004, 250], "hex": "#eceff1"}},
		"kelp": {"glow": ` + glow + `}
	},
	` + testRules + `
}`
}

func TestDefault(t *testing.T) {
	p, err := Default()
	if err != nil {
		t.Fatalf("Default() error = %v", err)
	}
	if got := len(p.Colours); got != 15 {
		t.Errorf("len(Default().Colours) = %d, want %d", got, 15)
	}
	if errs := p.Check(); len(errs) != 0 {
		t.Errorf("Default().Check() = %v, want no errors", errs)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "invalid json", input: `{`, expected: "decode palette"},
		{name: "wrong space", input: `{"space": "hsl"}`, expected: `space "hsl", want oklch`},
		{
			name: "duplicate colour",
			input: `{"space": "oklch", "groups": {
				"a": {"x": {"oklch": [0.5, 0, 0], "hex": "#636363"}},
				"b": {"x": {"oklch": [0.5, 0, 0], "hex": "#636363"}}
			}}`,
			expected: `colour "x" defined twice`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.input))
			if err == nil || !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("Parse() error = %v, want it to contain %q", err, tt.expected)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name     string
		glow     string
		expected []string
	}{
		{
			name:     "valid",
			glow:     `{"oklch": [0.85, 0.105, 218], "hex": "#76dffb"}`,
			expected: nil,
		},
		{
			name:     "out of gamut",
			glow:     `{"oklch": [0.88, 0.10, 230], "hex": "#8fe4ff"}`,
			expected: []string{"glow: OKLCH {0.88 0.1 230} is outside the sRGB gamut"},
		},
		{
			name:     "hex mismatch",
			glow:     `{"oklch": [0.85, 0.105, 218], "hex": "#76dffc"}`,
			expected: []string{"glow: hex #76dffc does not match OKLCH {0.85 0.105 218} (#76dffb)"},
		},
		{
			name:     "low contrast",
			glow:     `{"oklch": [0.40, 0.05, 218], "hex": "#254e5a"}`,
			expected: []string{"glow on abyss: contrast 2.11:1 is below 4.50:1"},
		},
		{
			name:     "neighbours too close",
			glow:     `{"oklch": [0.94, 0.004, 250], "hex": "#e9ebee"}`,
			expected: []string{"neighbours white, glow: ΔE 0.010 is below 0.100"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse([]byte(testDocument(tt.glow)))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			var got []string
			for _, err := range p.Check() {
				got = append(got, err.Error())
			}
			if strings.Join(got, "\n") != strings.Join(tt.expected, "\n") {
				t.Errorf("Check() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func testRulesDocument(
	contrast string,
	rules string,
) string {
	return `{
	"name": "Test",
	"space": "oklch",
	"groups": {
		"abyss": {
			"abyss": {"oklch": [0.165, 0.033, 166], "hex": "#01130b"},
			"current": {"oklch": [0.25, 0.050, 165], "hex": "#02291c"}
		},
		"foam": {"white": {"oklch": [0.95, 0.004, 250], "hex": "#eceff1", "contrast": ` + contrast + `}},
		"kelp": {"verdigris": {"oklch": [0.77, 0.105, 135], "hex": "#97c37f"}},
		"reef": {"coral": {"oklch": [0.80, 0.100, 35], "hex": "#f7a792"}}
	},
	"rules": ` + rules + `
}`
}

func TestCheckLightnessContrastAndColourBlindness(t *testing.T) {
	const selectionRules = `{"background": "abyss", "selection": "current", "selectionAllowance": 10}`
	tests := []struct {
		name     string
		contrast string
		rules    string
		expected []string
	}{
		{
			name:     "within range",
			contrast: `{"min": 60, "max": 100}`,
			rules:    selectionRules,
		},
		{
			name:     "below minimum",
			contrast: `{"min": 98}`,
			rules:    selectionRules,
			expected: []string{"white on abyss: Lc 96.6 is below 98.0"},
		},
		{
			name:     "above maximum",
			contrast: `{"min": 60, "max": 92}`,
			rules:    selectionRules,
			expected: []string{"white on abyss: Lc 96.6 is above 92.0"},
		},
		{
			name:     "below minimum on the selection",
			contrast: `{"min": 95}`,
			rules:    `{"background": "abyss", "selection": "current"}`,
			expected: []string{"white on current: Lc 94.3 is below 95.0"},
		},
		{
			name:     "colour-blind neighbours too close",
			contrast: `null`,
			rules:    `{"minDeltaE": 0.10, "minDeltaEColourBlind": 0.07, "neighbours": [["verdigris", "coral"]]}`,
			expected: []string{
				"neighbours verdigris, coral with protanopia: ΔE 0.046 is below 0.070",
				"neighbours verdigris, coral with deuteranopia: ΔE 0.038 is below 0.070",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse([]byte(testRulesDocument(tt.contrast, tt.rules)))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			var got []string
			for _, err := range p.Check() {
				got = append(got, err.Error())
			}
			if strings.Join(got, "\n") != strings.Join(tt.expected, "\n") {
				t.Errorf("Check() = %q, want %q", got, tt.expected)
			}
		})
	}
}
