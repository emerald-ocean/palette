package palette

import (
	"math"
	"testing"
)

var (
	testAbyss     = OKLCH{L: 0.165, C: 0.033, H: 166}
	testWhite     = OKLCH{L: 0.95, C: 0.004, H: 250}
	testVerdigris = OKLCH{L: 0.77, C: 0.105, H: 135}
	testSeaglass  = OKLCH{L: 0.81, C: 0.085, H: 192}
)

func TestOKLCHHex(t *testing.T) {
	tests := []struct {
		name     string
		colour   OKLCH
		expected string
	}{
		{name: "abyss", colour: testAbyss, expected: "#01130b"},
		{name: "white", colour: testWhite, expected: "#eceff1"},
		{name: "verdigris", colour: testVerdigris, expected: "#97c37f"},
		{name: "seaglass", colour: testSeaglass, expected: "#7bd3cf"},
		{name: "black", colour: OKLCH{}, expected: "#000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.colour.Hex(); got != tt.expected {
				t.Errorf("Hex() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestOKLCHInGamut(t *testing.T) {
	tests := []struct {
		name     string
		colour   OKLCH
		expected bool
	}{
		{name: "glow", colour: OKLCH{L: 0.85, C: 0.105, H: 218}, expected: true},
		{name: "clipped v5 glow", colour: OKLCH{L: 0.88, C: 0.10, H: 230}, expected: false},
		{name: "white", colour: testWhite, expected: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.colour.InGamut(); got != tt.expected {
				t.Errorf("InGamut() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOKLCHDeltaE(t *testing.T) {
	tests := []struct {
		name     string
		a        OKLCH
		b        OKLCH
		expected float64
	}{
		{name: "same colour", a: testSeaglass, b: testSeaglass, expected: 0},
		{name: "seaglass to verdigris", a: testSeaglass, b: testVerdigris, expected: 0.101},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.DeltaE(tt.b); math.Abs(got-tt.expected) > 0.0005 {
				t.Errorf("DeltaE() = %.3f, want %.3f", got, tt.expected)
			}
		})
	}
}

func TestContrast(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected float64
	}{
		{name: "white on abyss", a: "#eceff1", b: "#01130b", expected: 16.54},
		{name: "symmetric", a: "#01130b", b: "#eceff1", expected: 16.54},
		{name: "same colour", a: "#01130b", b: "#01130b", expected: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := ParseHex(tt.a)
			if err != nil {
				t.Fatalf("ParseHex() error = %v", err)
			}
			b, err := ParseHex(tt.b)
			if err != nil {
				t.Fatalf("ParseHex() error = %v", err)
			}
			if got := Contrast(a, b); math.Abs(got-tt.expected) > 0.01 {
				t.Errorf("Contrast() = %.2f, want %.2f", got, tt.expected)
			}
		})
	}
}

func TestParseHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{name: "valid", input: "#76dffb", expected: "#76dffb"},
		{name: "missing hash", input: "76dffb", wantErr: true},
		{name: "short", input: "#fff", wantErr: true},
		{name: "not hex", input: "#zzzzzz", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseHex() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Hex() != tt.expected {
				t.Errorf("ParseHex() = %q, want %q", got.Hex(), tt.expected)
			}
		})
	}
}
