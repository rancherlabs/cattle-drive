package cmds

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/urfave/cli/v2"
)

var (
	Kubeconfig  string
	CommonFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "kubeconfig",
			EnvVars:     []string{"KUBECONFIG"},
			Usage:       "Kubeconfig path",
			Destination: &Kubeconfig,
		},
	}
	Spinner *spinner.Spinner
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "cattle-drive"
	app.Usage = "Tool for migrating rancher objects for RKE downstream clusters"
	Spinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithColor("green"))
	return app
}
