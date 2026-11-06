package main

import (
	"log/slog"
	"os"

	"github.com/emerald-ocean/palette/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout); err != nil {
		slog.Error("emerald-gen failed", "error", err)
		os.Exit(1)
	}
}
