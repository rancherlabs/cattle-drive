package main

import (
	"fmt"
	"galal-hussein/cattle-drive/cli/cmds"
	"galal-hussein/cattle-drive/cli/cmds/migrate"
	"galal-hussein/cattle-drive/cli/cmds/status"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	program   = "cattle-drive"
	version   = "dev"
	gitCommit = "HEAD"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []*cli.Command{
		status.NewCommand(),
		migrate.NewCommand(),
	}
	app.Version = version + " (" + gitCommit + ")"

	logrus.SetOutput(io.Discard)
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("exiting tool: %v", err)
		os.Exit(1)
	}
}
