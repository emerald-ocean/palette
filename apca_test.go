package palette

import (
	"math"
	"testing"
)

func TestAPCA(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		bg       string
		expected float64
	}{
		{name: "grey on white", text: "#888888", bg: "#ffffff", expected: 63.06},
		{name: "white on grey", text: "#ffffff", bg: "#888888", expected: -68.54},
		{name: "black on white", text: "#000000", bg: "#ffffff", expected: 106.04},
		{name: "white on black", text: "#ffffff", bg: "#000000", expected: -107.88},
		{name: "same colour", text: "#01130b", bg: "#01130b", expected: 0},
		{name: "white on abyss", text: "#eceff1", bg: "#01130b", expected: -96.6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := ParseHex(tt.text)
			if err != nil {
				t.Fatalf("ParseHex() error = %v", err)
			}
			bg, err := ParseHex(tt.bg)
			if err != nil {
				t.Fatalf("ParseHex() error = %v", err)
			}
			if got := APCA(text, bg); math.Abs(got-tt.expected) > 0.05 {
				t.Errorf("APCA() = %.2f, want %.2f", got, tt.expected)
			}
		})
	}
}
