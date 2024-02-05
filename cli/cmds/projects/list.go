package projects

import (
	"fmt"
	"galal-hussein/cattle-drive/cli/cmds"
	"strings"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	mgmtClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/urfave/cli"
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
	// ctx := context.Background()

	restConfig, err := clientcmd.BuildConfigFromFlags("", cmds.Kubeconfig)
	if err != nil {
		return err
	}
	client, err := mgmtClient.NewClient(&clientbase.ClientOpts{
		URL:       restConfig.Host,
		AccessKey: strings.Split(restConfig.BearerToken, ":")[0],
		SecretKey: strings.Split(restConfig.BearerToken, ":")[1],
		Insecure:  true,
	})
	if err != nil {
		return err
	}
	projects, _ := client.Project.List(&types.ListOpts{})
	fmt.Printf("%#v\n", projects)
	return nil
}
