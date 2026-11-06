package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/emerald-ocean/palette/pkg/colour"
	"github.com/emerald-ocean/palette/pkg/palette"
)

const usage = `usage: emerald-gen <command> [flags]

commands:
  check     check the palette against its rules
  generate  render a template with the palette's colours
  report    print APCA contrast and colour-blindness separation
`

var errRules = errors.New("palette breaks its rules")

// Run executes the emerald-gen command in args and writes its output to stdout.
func Run(
	args []string,
	stdout io.Writer,
) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command\n%s", usage)
	}
	switch args[0] {
	case "check":
		return check(args[1:], stdout)
	case "report":
		return report(args[1:], stdout)
	case "generate":
		return generate(args[1:])
	default:
		return fmt.Errorf("unknown command %q\n%s", args[0], usage)
	}
}

func check(
	args []string,
	stdout io.Writer,
) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	path := fs.String("palette", "", "palette JSON file (default: the embedded palette)")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse check flags: %w", err)
	}
	p, err := load(*path)
	if err != nil {
		return err
	}
	errs := p.Check()
	for _, err := range errs {
		if _, werr := fmt.Fprintln(stdout, err); werr != nil {
			return fmt.Errorf("write check result: %w", werr)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s %s: %d problems: %w", p.Name, p.Version, len(errs), errRules)
	}
	if _, err := fmt.Fprintf(stdout, "%s %s: all %d colours meet the rules\n", p.Name, p.Version, len(p.Colours)); err != nil {
		return fmt.Errorf("write check summary: %w", err)
	}
	return nil
}

func generate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	path := fs.String("palette", "", "palette JSON file (default: the embedded palette)")
	tmplPath := fs.String("template", "", "template file to render")
	outPath := fs.String("out", "", "file to write the result to")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse generate flags: %w", err)
	}
	if *tmplPath == "" || *outPath == "" {
		return errors.New("generate needs -template and -out")
	}
	p, err := load(*path)
	if err != nil {
		return err
	}
	text, err := os.ReadFile(*tmplPath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}
	var b bytes.Buffer
	if err := p.Render(&b, *tmplPath, string(text)); err != nil {
		return err
	}
	if err := os.WriteFile(*outPath, b.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", *outPath, err)
	}
	return nil
}

func load(path string) (palette.Palette, error) {
	p, err := parse(path)
	if err != nil {
		return palette.Palette{}, err
	}
	p.Version = version()
	return p, nil
}

func parse(path string) (palette.Palette, error) {
	if path == "" {
		return palette.Default()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return palette.Palette{}, fmt.Errorf("read palette: %w", err)
	}
	return palette.Parse(data)
}

func report(
	args []string,
	stdout io.Writer,
) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	path := fs.String("palette", "", "palette JSON file (default: the embedded palette)")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse report flags: %w", err)
	}
	p, err := load(*path)
	if err != nil {
		return err
	}
	rgb := make(map[string]colour.RGB, len(p.Colours))
	for name, c := range p.Colours {
		v, err := colour.ParseHex(c.Hex)
		if err != nil {
			return fmt.Errorf("report colour %s: %w", name, err)
		}
		rgb[name] = v
	}

	rows := [][]string{append([]string{"APCA Lc"}, p.Rules.Backgrounds...)}
	for _, name := range p.Names() {
		if slices.Contains(p.Rules.Backgrounds, name) {
			continue
		}
		row := []string{name}
		for _, bg := range p.Rules.Backgrounds {
			row = append(row, fmt.Sprintf("%.1f", math.Abs(colour.APCA(rgb[name], rgb[bg]))))
		}
		rows = append(rows, row)
	}
	header := []string{"Neighbour ΔE"}
	for _, v := range colour.Visions {
		header = append(header, v.String())
	}
	rows = append(rows, nil, header)
	for _, pair := range p.Rules.Neighbours {
		row := []string{pair[0] + " " + pair[1]}
		for _, v := range colour.Visions {
			a, b := colour.Simulate(rgb[pair[0]], v), colour.Simulate(rgb[pair[1]], v)
			row = append(row, fmt.Sprintf("%.3f", a.DeltaE(b)))
		}
		rows = append(rows, row)
	}

	w := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	for _, row := range rows {
		if _, err := io.WriteString(w, strings.Join(row, "\t")+"\t\n"); err != nil {
			return fmt.Errorf("write report row: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}
