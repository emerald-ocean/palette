package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/emerald-ocean/palette"
)

const usage = `usage: emerald-gen <command> [flags]

commands:
  check     check the palette against its rules
  generate  render a template with the palette's colours
`

var errRules = errors.New("palette breaks its rules")

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		slog.Error("emerald-gen failed", "error", err)
		os.Exit(1)
	}
}

func run(
	args []string,
	stdout io.Writer,
) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command\n%s", usage)
	}
	switch args[0] {
	case "check":
		return check(args[1:], stdout)
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
	if path == "" {
		return palette.Default()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return palette.Palette{}, fmt.Errorf("read palette: %w", err)
	}
	return palette.Parse(data)
}
