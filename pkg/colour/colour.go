package colour

import (
	"fmt"
	"math"
	"strconv"
)

const gamutTolerance = 1e-4

// OKLCH is a colour in the OKLCH space, with the hue in degrees.
type OKLCH struct {
	L float64
	C float64
	H float64
}

// RGB is a gamma-encoded sRGB colour with channels in [0, 1].
type RGB struct {
	R float64
	G float64
	B float64
}

func (c OKLCH) lab() (float64, float64, float64) {
	h := c.H * math.Pi / 180
	return c.L, c.C * math.Cos(h), c.C * math.Sin(h)
}

func (c OKLCH) linear() (float64, float64, float64) {
	l, a, b := c.lab()
	lc := cube(l + 0.3963377774*a + 0.2158037573*b)
	mc := cube(l - 0.1055613458*a - 0.0638541728*b)
	sc := cube(l - 0.0894841775*a - 1.2914855480*b)
	return 4.0767416621*lc - 3.3077115913*mc + 0.2309699292*sc,
		-1.2684380046*lc + 2.6097574011*mc - 0.3413193965*sc,
		-0.0041960863*lc - 0.7034186147*mc + 1.7076147010*sc
}

// InGamut reports whether c can be shown in sRGB without clipping.
func (c OKLCH) InGamut() bool {
	r, g, b := c.linear()
	for _, v := range []float64{r, g, b} {
		if v < -gamutTolerance || v > 1+gamutTolerance {
			return false
		}
	}
	return true
}

// RGB converts c to sRGB, clipping it to the gamut.
func (c OKLCH) RGB() RGB {
	r, g, b := c.linear()
	return RGB{R: encode(r), G: encode(g), B: encode(b)}
}

// Hex returns c as a #rrggbb string, clipping it to the gamut.
func (c OKLCH) Hex() string {
	return c.RGB().Hex()
}

// DeltaE returns the perceptual distance between c and o in OKLab.
func (c OKLCH) DeltaE(o OKLCH) float64 {
	l1, a1, b1 := c.lab()
	l2, a2, b2 := o.lab()
	return math.Sqrt((l1-l2)*(l1-l2) + (a1-a2)*(a1-a2) + (b1-b2)*(b1-b2))
}

// Hex returns c as a #rrggbb string.
func (c RGB) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", channel(c.R), channel(c.G), channel(c.B))
}

func (c RGB) luminance() float64 {
	return 0.2126*decode(c.R) + 0.7152*decode(c.G) + 0.0722*decode(c.B)
}

// ParseHex parses a #rrggbb string.
func ParseHex(s string) (RGB, error) {
	if len(s) != 7 || s[0] != '#' {
		return RGB{}, fmt.Errorf("parse hex %q: want #rrggbb", s)
	}
	v, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return RGB{}, fmt.Errorf("parse hex %q: %w", s, err)
	}
	return RGB{
		R: float64(v>>16&0xff) / 255,
		G: float64(v>>8&0xff) / 255,
		B: float64(v&0xff) / 255,
	}, nil
}

// Contrast returns the WCAG contrast ratio between a and b.
func Contrast(
	a RGB,
	b RGB,
) float64 {
	la, lb := a.luminance(), b.luminance()
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func encode(x float64) float64 {
	x = math.Max(0, math.Min(1, x))
	if x <= 0.0031308 {
		return 12.92 * x
	}
	return 1.055*math.Pow(x, 1/2.4) - 0.055
}

func decode(x float64) float64 {
	if x <= 0.04045 {
		return x / 12.92
	}
	return math.Pow((x+0.055)/1.055, 2.4)
}

func channel(x float64) int {
	return int(math.Round(math.Max(0, math.Min(1, x)) * 255))
}

func cube(x float64) float64 {
	return x * x * x
}
