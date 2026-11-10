package palette

import (
	"fmt"
	"math"
	"slices"

	"github.com/emerald-ocean/palette/pkg/colour"
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
	errs = append(errs, p.checkLightnessContrast()...)
	errs = append(errs, p.checkNeighbours()...)
	return errs
}

func (p Palette) checkLightnessContrast() []error {
	if p.Rules.Background == "" {
		return nil
	}
	bg, err := p.rgb(p.Rules.Background)
	if err != nil {
		return []error{fmt.Errorf("lightness contrast background: %w", err)}
	}
	var sel *colour.RGB
	if p.Rules.Selection != "" {
		s, err := p.rgb(p.Rules.Selection)
		if err != nil {
			return []error{fmt.Errorf("lightness contrast selection: %w", err)}
		}
		sel = &s
	}
	var errs []error
	for _, name := range p.Names() {
		r := p.Colours[name].Contrast
		if r == nil {
			continue
		}
		fg, err := p.rgb(name)
		if err != nil {
			errs = append(errs, fmt.Errorf("lightness contrast foreground: %w", err))
			continue
		}
		lc := math.Abs(colour.APCA(fg, bg))
		if lc < r.Min {
			errs = append(errs, fmt.Errorf("%s on %s: Lc %.1f is below %.1f", name, p.Rules.Background, lc, r.Min))
		}
		if r.Max > 0 && lc > r.Max {
			errs = append(errs, fmt.Errorf("%s on %s: Lc %.1f is above %.1f", name, p.Rules.Background, lc, r.Max))
		}
		if sel == nil {
			continue
		}
		minSel := r.Min - p.Rules.SelectionAllowance
		if lc := math.Abs(colour.APCA(fg, *sel)); lc < minSel {
			errs = append(errs, fmt.Errorf("%s on %s: Lc %.1f is below %.1f", name, p.Rules.Selection, lc, minSel))
		}
	}
	return errs
}

func (p Palette) rgb(name string) (colour.RGB, error) {
	c, ok := p.Colours[name]
	if !ok {
		return colour.RGB{}, fmt.Errorf("unknown colour %q", name)
	}
	v, err := colour.ParseHex(c.Hex)
	if err != nil {
		return colour.RGB{}, fmt.Errorf("colour %s: %w", name, err)
	}
	return v, nil
}

func (p Palette) checkContrast() []error {
	var errs []error
	for _, bgName := range p.Rules.Backgrounds {
		bg, ok := p.Colours[bgName]
		if !ok {
			errs = append(errs, fmt.Errorf("contrast: unknown background %q", bgName))
			continue
		}
		bgRGB, err := colour.ParseHex(bg.Hex)
		if err != nil {
			errs = append(errs, fmt.Errorf("contrast background %s: %w", bgName, err))
			continue
		}
		for _, name := range p.Names() {
			c := p.Colours[name]
			if c.Group == bg.Group || slices.Contains(p.Rules.Backgrounds, name) {
				continue
			}
			fg, err := colour.ParseHex(c.Hex)
			if err != nil {
				errs = append(errs, fmt.Errorf("contrast foreground %s: %w", name, err))
				continue
			}
			if ratio := colour.Contrast(fg, bgRGB); ratio < p.Rules.MinContrast {
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
		if p.Rules.MinDeltaEColourBlind > 0 {
			errs = append(errs, p.checkColourBlindPair(pair)...)
		}
	}
	return errs
}

func (p Palette) checkColourBlindPair(pair [2]string) []error {
	a, err := p.rgb(pair[0])
	if err != nil {
		return []error{fmt.Errorf("neighbours %s, %s: %w", pair[0], pair[1], err)}
	}
	b, err := p.rgb(pair[1])
	if err != nil {
		return []error{fmt.Errorf("neighbours %s, %s: %w", pair[0], pair[1], err)}
	}
	var errs []error
	for _, v := range []colour.Vision{colour.Protanopia, colour.Deuteranopia} {
		if d := colour.Simulate(a, v).DeltaE(colour.Simulate(b, v)); d < p.Rules.MinDeltaEColourBlind {
			errs = append(errs, fmt.Errorf(
				"neighbours %s, %s with %s: ΔE %.3f is below %.3f",
				pair[0],
				pair[1],
				v,
				d,
				p.Rules.MinDeltaEColourBlind,
			))
		}
	}
	return errs
}
