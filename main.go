package main

import (
	"fmt"
	"galal-hussein/cattle-drive/cli/cmds"
	"galal-hussein/cattle-drive/cli/cmds/interactive"
	"galal-hussein/cattle-drive/cli/cmds/migrate"
	"galal-hussein/cattle-drive/cli/cmds/status"
	"galal-hussein/cattle-drive/pkg/version"
	"io"
	"os"

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
