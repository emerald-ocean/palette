package palette

import "math"

// Vision is a type of colour vision.
type Vision int

// Colour vision types, simulated at full severity.
const (
	Normal Vision = iota
	Protanopia
	Deuteranopia
	Tritanopia
)

// Visions lists every colour vision type.
var Visions = []Vision{Normal, Protanopia, Deuteranopia, Tritanopia}

var visionMatrices = map[Vision][9]float64{
	Protanopia: {
		0.152286, 1.052583, -0.204868,
		0.114503, 0.786281, 0.099216,
		-0.003882, -0.048116, 1.051998,
	},
	Deuteranopia: {
		0.367322, 0.860646, -0.227968,
		0.280085, 0.672501, 0.047413,
		-0.011820, 0.042940, 0.968881,
	},
	Tritanopia: {
		1.255528, -0.076749, -0.178779,
		-0.078411, 0.930809, 0.147602,
		0.004733, 0.691367, 0.303900,
	},
}

func (v Vision) String() string {
	switch v {
	case Normal:
		return "normal"
	case Protanopia:
		return "protanopia"
	case Deuteranopia:
		return "deuteranopia"
	case Tritanopia:
		return "tritanopia"
	default:
		return "unknown"
	}
}

// Simulate returns c as seen with vision v, using the Machado, Oliveira and Fernandes (2009) model.
func Simulate(
	c RGB,
	v Vision,
) RGB {
	m, ok := visionMatrices[v]
	if !ok {
		return c
	}
	r, g, b := decode(c.R), decode(c.G), decode(c.B)
	return RGB{
		R: encode(m[0]*r + m[1]*g + m[2]*b),
		G: encode(m[3]*r + m[4]*g + m[5]*b),
		B: encode(m[6]*r + m[7]*g + m[8]*b),
	}
}

// DeltaE returns the perceptual distance between c and o in OKLab.
func (c RGB) DeltaE(o RGB) float64 {
	l1, a1, b1 := c.lab()
	l2, a2, b2 := o.lab()
	return math.Sqrt((l1-l2)*(l1-l2) + (a1-a2)*(a1-a2) + (b1-b2)*(b1-b2))
}

func (c RGB) lab() (float64, float64, float64) {
	r, g, b := decode(c.R), decode(c.G), decode(c.B)
	l := math.Cbrt(0.4122214708*r + 0.5363325363*g + 0.0514459929*b)
	m := math.Cbrt(0.2119034982*r + 0.6806995451*g + 0.1073969566*b)
	s := math.Cbrt(0.0883024619*r + 0.2817188376*g + 0.6299787005*b)
	return 0.2104542553*l + 0.7936177850*m - 0.0040720468*s,
		1.9779984951*l - 2.4285922050*m + 0.4505937099*s,
		0.0259040371*l + 0.7827717662*m - 0.8086757660*s
}
