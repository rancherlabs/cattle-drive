package cmds

import (
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
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "cattle-drive"
	app.Usage = "Tool for migrating rancher objects for RKE downstream clusters"
	return app
}
