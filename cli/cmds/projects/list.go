package projects

import (
	"context"
	"galal-hussein/cattle-drive/cli/cmds"
	"galal-hussein/cattle-drive/pkg/client"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/urfave/cli"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	name string

	projectsListFlags = []cli.Flag{
		cli.StringFlag{
			Name:        "name",
			Usage:       "name of the cluster",
			Destination: &name,
		},
	}
)

func list(clx *cli.Context) error {
	ctx := context.Background()

	restConfig, err := clientcmd.BuildConfigFromFlags("", cmds.Kubeconfig)
	if err != nil {
		return err
	}
	client, err := client.New(ctx, restConfig)
	if err != nil {
		return err
	}
	var projects v3.ProjectList

	if err := client.Projects.List(ctx, "", &projects, v1.ListOptions{}); err != nil {
		return err
	}
	return nil
}
