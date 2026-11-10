package recolour

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/emerald-ocean/palette/pkg/colour"
)

const (
	testDark  = "#01130b"
	testGreen = "#54f9a3"
	testFoam  = "#dbeaeb"
)

func testStops(t testing.TB, hexes ...string) []colour.RGB {
	t.Helper()
	stops := make([]colour.RGB, len(hexes))
	for i, h := range hexes {
		c, err := colour.ParseHex(h)
		if err != nil {
			t.Fatalf("ParseHex() error = %v", err)
		}
		stops[i] = c
	}
	return stops
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		stops   []string
		opts    []Option
		wantErr string
	}{
		{name: "valid", stops: []string{testDark, testGreen, testFoam}},
		{name: "one stop", stops: []string{testDark}, wantErr: "at least 2 stops"},
		{name: "not lighter", stops: []string{testDark, testFoam, testGreen}, wantErr: "stop 3 (#54f9a3) is not lighter"},
		{
			name:    "input range reversed",
			stops:   []string{testDark, testFoam},
			opts:    []Option{WithInput(0.8, 0.2)},
			wantErr: "must be below light",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(testStops(t, tt.stops...), tt.opts...)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Errorf("New() error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestRampColour(t *testing.T) {
	r, err := New(testStops(t, testDark, testGreen, testFoam), WithInput(0.2, 0.8))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	tests := []struct {
		name     string
		in       string
		expected string
	}{
		{name: "black clamps to the darkest stop", in: "#000000", expected: testDark},
		{name: "white clamps to the lightest stop", in: "#ffffff", expected: testFoam},
		{name: "red lands on the ramp", in: "#ff0000", expected: "#3cbc7b"},
		{name: "grey lands on the ramp", in: "#808080", expected: "#37af72"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := testStops(t, tt.in)[0]
			if got := r.Colour(in).Hex(); got != tt.expected {
				t.Errorf("Colour(%s) = %q, want %q", tt.in, got, tt.expected)
			}
		})
	}
}

func TestRampImage(t *testing.T) {
	r, err := New(testStops(t, testDark, testGreen, testFoam))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	src := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	src.SetNRGBA(0, 0, color.NRGBA{A: 255})
	src.SetNRGBA(1, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 128})
	got := r.Image(src)

	tests := []struct {
		name     string
		x        int
		expected color.NRGBA
	}{
		{name: "black", x: 0, expected: color.NRGBA{R: 0x01, G: 0x13, B: 0x0b, A: 255}},
		{name: "translucent white keeps its alpha", x: 1, expected: color.NRGBA{R: 0xdb, G: 0xea, B: 0xeb, A: 128}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if c := got.NRGBAAt(tt.x, 0); c != tt.expected {
				t.Errorf("Image().NRGBAAt(%d, 0) = %v, want %v", tt.x, c, tt.expected)
			}
		})
	}
}
