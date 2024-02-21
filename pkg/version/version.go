package version

import "strings"

var (
	Program      = "cattle-drive"
	ProgramUpper = strings.ToUpper(Program)
	Version      = "dev"
	GitCommit    = "HEAD"
)
