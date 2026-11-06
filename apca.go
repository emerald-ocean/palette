package palette

import "math"

// APCA returns the APCA lightness contrast (Lc) of text on bg, using APCA-W3 0.0.98G-4g.
// It is positive for dark text on a light background and negative for light text on a dark background.
func APCA(
	text RGB,
	bg RGB,
) float64 {
	yText, yBg := screenLuminance(text), screenLuminance(bg)
	if math.Abs(yBg-yText) < 0.0005 {
		return 0
	}
	if yBg > yText {
		sapc := (math.Pow(yBg, 0.56) - math.Pow(yText, 0.57)) * 1.14
		if sapc < 0.1 {
			return 0
		}
		return (sapc - 0.027) * 100
	}
	sapc := (math.Pow(yBg, 0.65) - math.Pow(yText, 0.62)) * 1.14
	if sapc > -0.1 {
		return 0
	}
	return (sapc + 0.027) * 100
}

func screenLuminance(c RGB) float64 {
	y := 0.2126729*math.Pow(c.R, 2.4) + 0.7151522*math.Pow(c.G, 2.4) + 0.0721750*math.Pow(c.B, 2.4)
	if y < 0.022 {
		y += math.Pow(0.022-y, 1.414)
	}
	return y
}
