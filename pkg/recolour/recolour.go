package recolour

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/emerald-ocean/palette/pkg/colour"
)

const (
	defaultDark  = 0.19
	defaultLight = 0.80
	steps        = 4096
)

// Ramp recolours by lightness: every colour lands on a line through the ramp's stops in OKLab.
type Ramp struct {
	dark  float64
	light float64
	out   *[steps][3]uint8
}

// Option configures a Ramp.
type Option func(*Ramp)

// WithInput sets the input lightness (OKLab L) that maps to the darkest stop and the one that maps to the lightest.
func WithInput(
	dark float64,
	light float64,
) Option {
	return func(r *Ramp) {
		r.dark = dark
		r.light = light
	}
}

// New returns a Ramp through stops, which must get lighter from first to last.
func New(
	stops []colour.RGB,
	opts ...Option,
) (Ramp, error) {
	if len(stops) < 2 {
		return Ramp{}, errors.New("new ramp: need at least 2 stops")
	}
	r := Ramp{dark: defaultDark, light: defaultLight}
	for _, opt := range opts {
		opt(&r)
	}
	if r.dark >= r.light {
		return Ramp{}, fmt.Errorf("new ramp: input dark %.2f must be below light %.2f", r.dark, r.light)
	}
	labs := make([]colour.OKLab, len(stops))
	for i, s := range stops {
		labs[i] = s.OKLab()
		if i > 0 && labs[i].L <= labs[i-1].L {
			return Ramp{}, fmt.Errorf("new ramp: stop %d (%s) is not lighter than stop %d", i+1, s.Hex(), i)
		}
	}
	r.out = new([steps][3]uint8)
	for i := range r.out {
		c := along(labs, float64(i)/(steps-1)).RGB()
		r.out[i] = [3]uint8{channel(c.R), channel(c.G), channel(c.B)}
	}
	return r, nil
}

// along returns the colour at t (0 to 1) of the way through stops by lightness.
func along(
	stops []colour.OKLab,
	t float64,
) colour.OKLab {
	first, last := stops[0], stops[len(stops)-1]
	x := first.L + t*(last.L-first.L)
	i := 1
	for i < len(stops)-1 && x > stops[i].L {
		i++
	}
	a, b := stops[i-1], stops[i]
	f := min(max((x-a.L)/(b.L-a.L), 0), 1)
	return colour.OKLab{
		L: a.L + f*(b.L-a.L),
		A: a.A + f*(b.A-a.A),
		B: a.B + f*(b.B-a.B),
	}
}

func (r Ramp) lookup(c colour.RGB) [3]uint8 {
	t := (c.OKLab().L - r.dark) / (r.light - r.dark)
	return r.out[int(math.Round(min(max(t, 0), 1)*(steps-1)))]
}

// Colour returns c moved onto the ramp.
func (r Ramp) Colour(c colour.RGB) colour.RGB {
	v := r.lookup(c)
	return colour.RGB{R: float64(v[0]) / 255, G: float64(v[1]) / 255, B: float64(v[2]) / 255}
}

// Image returns a copy of src with every pixel moved onto the ramp, keeping its alpha.
func (r Ramp) Image(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	for i := 0; i < len(dst.Pix); i += 4 {
		px := dst.Pix[i : i+3 : i+3]
		v := r.lookup(colour.RGB{
			R: float64(px[0]) / 255,
			G: float64(px[1]) / 255,
			B: float64(px[2]) / 255,
		})
		px[0], px[1], px[2] = v[0], v[1], v[2]
	}
	return dst
}

func channel(x float64) uint8 {
	return uint8(min(max(x, 0), 1)*255 + 0.5)
}
