package main

import (
	"fmt"
	"io"
	"os"
	"rancherlabs/cattle-drive/cli/cmds"
	"rancherlabs/cattle-drive/cli/cmds/interactive"
	"rancherlabs/cattle-drive/cli/cmds/migrate"
	"rancherlabs/cattle-drive/cli/cmds/status"
	"rancherlabs/cattle-drive/pkg/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []*cli.Command{
		status.NewCommand(),
		migrate.NewCommand(),
		interactive.NewCommand(),
	}
	app.Version = fmt.Sprintf("%s (%s)", version.Version, version.GitCommit)

	logrus.SetOutput(io.Discard)
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("exiting tool: %v", err)
		os.Exit(1)
	}
}
