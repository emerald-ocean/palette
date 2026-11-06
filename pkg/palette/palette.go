package palette

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/emerald-ocean/palette/pkg/colour"
)

//go:embed palette.json
var embedded []byte

// Colour is a named palette colour.
type Colour struct {
	Name  string
	Group string
	OKLCH colour.OKLCH
	Hex   string
	Use   string
}

// Rules are the readability rules a palette must meet.
type Rules struct {
	MinContrast float64     `json:"minContrast"`
	MinDeltaE   float64     `json:"minDeltaE"`
	Backgrounds []string    `json:"backgrounds"`
	Neighbours  [][2]string `json:"neighbours"`
}

// Palette is a set of named colours and the rules they must meet.
type Palette struct {
	Name    string
	Version string
	Colours map[string]Colour
	Rules   Rules
}

type document struct {
	Name   string                               `json:"name"`
	Space  string                               `json:"space"`
	Groups map[string]map[string]documentColour `json:"groups"`
	Rules  Rules                                `json:"rules"`
}

type documentColour struct {
	OKLCH [3]float64 `json:"oklch"`
	Hex   string     `json:"hex"`
	Use   string     `json:"use"`
}

// Default returns the embedded Emerald Ocean palette.
func Default() (Palette, error) {
	return Parse(embedded)
}

// Parse decodes a palette from its JSON document.
func Parse(data []byte) (Palette, error) {
	var doc document
	if err := json.Unmarshal(data, &doc); err != nil {
		return Palette{}, fmt.Errorf("decode palette: %w", err)
	}
	if doc.Space != "oklch" {
		return Palette{}, fmt.Errorf("decode palette: space %q, want oklch", doc.Space)
	}
	p := Palette{
		Name:    doc.Name,
		Colours: map[string]Colour{},
		Rules:   doc.Rules,
	}
	for group, colours := range doc.Groups {
		for name, c := range colours {
			if _, ok := p.Colours[name]; ok {
				return Palette{}, fmt.Errorf("decode palette: colour %q defined twice", name)
			}
			p.Colours[name] = Colour{
				Name:  name,
				Group: group,
				OKLCH: colour.OKLCH{L: c.OKLCH[0], C: c.OKLCH[1], H: c.OKLCH[2]},
				Hex:   c.Hex,
				Use:   c.Use,
			}
		}
	}
	return p, nil
}

// Names returns the colour names in alphabetical order.
func (p Palette) Names() []string {
	names := make([]string, 0, len(p.Colours))
	for name := range p.Colours {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Hexes returns each colour's hex value by name.
func (p Palette) Hexes() map[string]string {
	hexes := make(map[string]string, len(p.Colours))
	for name, c := range p.Colours {
		hexes[name] = c.Hex
	}
	return hexes
}
