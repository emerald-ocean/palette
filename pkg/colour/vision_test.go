package colour

import (
	"math"
	"testing"
)

const (
	testVerdigrisHex = "#97c37f"
	testCoralHex     = "#f7a792"
)

func mustHex(
	t *testing.T,
	s string,
) RGB {
	t.Helper()
	c, err := ParseHex(s)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}
	return c
}

func TestSimulateKeepsGreys(t *testing.T) {
	for _, v := range Visions {
		t.Run(v.String(), func(t *testing.T) {
			grey := mustHex(t, "#808080")
			if got := Simulate(grey, v).Hex(); got != "#808080" {
				t.Errorf("Simulate() = %q, want %q", got, "#808080")
			}
		})
	}
}

func TestSimulatedDeltaE(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		vision   Vision
		expected float64
	}{
		{name: "keyword and constant, normal", a: testVerdigrisHex, b: testCoralHex, vision: Normal, expected: 0.159},
		{name: "keyword and constant, protanopia", a: testVerdigrisHex, b: testCoralHex, vision: Protanopia, expected: 0.046},
		{name: "keyword and constant, deuteranopia", a: testVerdigrisHex, b: testCoralHex, vision: Deuteranopia, expected: 0.038},
		{name: "keyword and constant, tritanopia", a: testVerdigrisHex, b: testCoralHex, vision: Tritanopia, expected: 0.164},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Simulate(mustHex(t, tt.a), tt.vision)
			b := Simulate(mustHex(t, tt.b), tt.vision)
			if got := a.DeltaE(b); math.Abs(got-tt.expected) > 0.0015 {
				t.Errorf("DeltaE() = %.3f, want %.3f", got, tt.expected)
			}
		})
	}
}

func TestVisionString(t *testing.T) {
	tests := []struct {
		name     string
		vision   Vision
		expected string
	}{
		{name: "normal", vision: Normal, expected: "normal"},
		{name: "deuteranopia", vision: Deuteranopia, expected: "deuteranopia"},
		{name: "unknown", vision: Vision(99), expected: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vision.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}
