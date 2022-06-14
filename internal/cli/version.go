package cli

import "fmt"

// Injected by goreleaser through ldflags (see .goreleaser.yml)
var (
	version = "0.0.0"
	commit  = "local"
	builtBy = "local"
	date    = "local"
)

func Version() string {
	return fmt.Sprintf("mtcli version: %v, commit: %v, builtBy: %v (%v)", version, commit, builtBy, date)
}
