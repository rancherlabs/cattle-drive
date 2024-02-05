package projects

import (
	"galal-hussein/cattle-drive/cli/cmds"

	"github.com/urfave/cli"
)

var subcommands = []cli.Command{
	{
		Name:            "list",
		Usage:           "Create new cluster",
		SkipFlagParsing: false,
		SkipArgReorder:  true,
		Action:          list,
		Flags:           append(cmds.CommonFlags, projectsListFlags...),
	},
}

func NewCommand() cli.Command {
	return cli.Command{
		Name:        "projects",
		Usage:       "projects command",
		Subcommands: subcommands,
	}
}
