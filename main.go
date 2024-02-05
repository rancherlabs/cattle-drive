package main

import (
	"galal-hussein/cattle-drive/cli/cmds"
	"galal-hussein/cattle-drive/cli/cmds/projects"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	program   = "cattle-drive"
	version   = "dev"
	gitCommit = "HEAD"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []cli.Command{
		projects.NewCommand(),
	}
	app.Version = version + " (" + gitCommit + ")"

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
