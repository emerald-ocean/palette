package palette

import (
	"fmt"
	"io"
	"math"
	"text/template"

	"github.com/emerald-ocean/palette/pkg/colour"
)

// Render executes the template text with the palette's colours by name and writes the result to w.
// Templates can use {{.emerald}} for a hex value, {{alpha .glow 0.16}} for a hex value with an alpha
// channel and {{version}} for the palette version.
func (p Palette) Render(
	w io.Writer,
	name string,
	text string,
) error {
	funcs := template.FuncMap{
		"alpha":   alpha,
		"version": func() string { return p.Version },
	}
	t, err := template.New(name).Option("missingkey=error").Funcs(funcs).Parse(text)
	if err != nil {
		return fmt.Errorf("parse template %s: %w", name, err)
	}
	if err := t.Execute(w, p.Hexes()); err != nil {
		return fmt.Errorf("render template %s: %w", name, err)
	}
	return nil
}

func alpha(
	hex string,
	a float64,
) (string, error) {
	if _, err := colour.ParseHex(hex); err != nil {
		return "", fmt.Errorf("alpha: %w", err)
	}
	if a < 0 || a > 1 {
		return "", fmt.Errorf("alpha %v: want a value from 0 to 1", a)
	}
	return fmt.Sprintf("%s%02x", hex, int(math.Round(a*255))), nil
}
