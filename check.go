package palette

import (
	"fmt"
	"slices"
)

// Check returns every rule the palette breaks, or nil if it meets them all.
func (p Palette) Check() []error {
	var errs []error
	for _, name := range p.Names() {
		c := p.Colours[name]
		if !c.OKLCH.InGamut() {
			errs = append(errs, fmt.Errorf("%s: OKLCH %v is outside the sRGB gamut", name, c.OKLCH))
		}
		if got := c.OKLCH.Hex(); got != c.Hex {
			errs = append(errs, fmt.Errorf("%s: hex %s does not match OKLCH %v (%s)", name, c.Hex, c.OKLCH, got))
		}
	}
	errs = append(errs, p.checkContrast()...)
	errs = append(errs, p.checkNeighbours()...)
	return errs
}

func (p Palette) checkContrast() []error {
	var errs []error
	for _, bgName := range p.Rules.Backgrounds {
		bg, ok := p.Colours[bgName]
		if !ok {
			errs = append(errs, fmt.Errorf("contrast: unknown background %q", bgName))
			continue
		}
		bgRGB, err := ParseHex(bg.Hex)
		if err != nil {
			errs = append(errs, fmt.Errorf("contrast background %s: %w", bgName, err))
			continue
		}
		for _, name := range p.Names() {
			c := p.Colours[name]
			if c.Group == bg.Group || slices.Contains(p.Rules.Backgrounds, name) {
				continue
			}
			fg, err := ParseHex(c.Hex)
			if err != nil {
				errs = append(errs, fmt.Errorf("contrast foreground %s: %w", name, err))
				continue
			}
			if ratio := Contrast(fg, bgRGB); ratio < p.Rules.MinContrast {
				errs = append(errs, fmt.Errorf(
					"%s on %s: contrast %.2f:1 is below %.2f:1",
					name,
					bgName,
					ratio,
					p.Rules.MinContrast,
				))
			}
		}
	}
	return errs
}

func (p Palette) checkNeighbours() []error {
	var errs []error
	for _, pair := range p.Rules.Neighbours {
		a, okA := p.Colours[pair[0]]
		b, okB := p.Colours[pair[1]]
		if !okA || !okB {
			errs = append(errs, fmt.Errorf("neighbours %s, %s: unknown colour", pair[0], pair[1]))
			continue
		}
		if d := a.OKLCH.DeltaE(b.OKLCH); d < p.Rules.MinDeltaE {
			errs = append(errs, fmt.Errorf(
				"neighbours %s, %s: ΔE %.3f is below %.3f",
				pair[0],
				pair[1],
				d,
				p.Rules.MinDeltaE,
			))
		}
	}
	return errs
}
