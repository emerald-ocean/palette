package cli

import "runtime/debug"

// version returns the module version emerald-gen was built from: the release tag when installed
// with `go run github.com/emerald-ocean/palette/cmd/emerald-gen@vX.Y.Z`, or "(devel)" otherwise.
func version() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok || bi.Main.Version == "" {
		return "(devel)"
	}
	return bi.Main.Version
}
