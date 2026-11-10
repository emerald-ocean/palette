package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/emerald-ocean/palette/pkg/colour"
	"github.com/emerald-ocean/palette/pkg/recolour"
)

func recolourImage(args []string) error {
	fs := flag.NewFlagSet("recolour", flag.ContinueOnError)
	path := fs.String("palette", "", "palette JSON file (default: the embedded palette)")
	inPath := fs.String("in", "", "JPEG or PNG image to recolour")
	outPath := fs.String("out", "", "PNG file to write the result to")
	ramp := fs.String("ramp", "abyss1,abyss2,abyss3,emerald3,emerald4,foam5", "comma-separated colours, darkest first")
	dark := fs.Float64("dark", 0.19, "input lightness (OKLab L) that maps to the darkest ramp colour")
	light := fs.Float64("light", 0.80, "input lightness (OKLab L) that maps to the lightest ramp colour")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse recolour flags: %w", err)
	}
	if *inPath == "" || *outPath == "" {
		return errors.New("recolour needs -in and -out")
	}
	p, err := load(*path)
	if err != nil {
		return err
	}
	var stops []colour.RGB
	for name := range strings.SplitSeq(*ramp, ",") {
		c, ok := p.Colours[name]
		if !ok {
			return fmt.Errorf("recolour ramp: unknown colour %q", name)
		}
		v, err := colour.ParseHex(c.Hex)
		if err != nil {
			return fmt.Errorf("recolour ramp colour %s: %w", name, err)
		}
		stops = append(stops, v)
	}
	r, err := recolour.New(stops, recolour.WithInput(*dark, *light))
	if err != nil {
		return err
	}

	data, err := os.ReadFile(*inPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", *inPath, err)
	}
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode %s: %w", *inPath, err)
	}

	var b bytes.Buffer
	if err := png.Encode(&b, r.Image(src)); err != nil {
		return fmt.Errorf("encode %s: %w", *outPath, err)
	}
	if err := os.WriteFile(*outPath, b.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", *outPath, err)
	}
	return nil
}
