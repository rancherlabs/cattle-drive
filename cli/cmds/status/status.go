package status

import (
	"context"
	"galal-hussein/cattle-drive/cli/cmds"
	"galal-hussein/cattle-drive/pkg/client"
	"galal-hussein/cattle-drive/pkg/cluster"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	source      string
	target      string
	statusFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "source",
			Usage:       "name of the source cluster",
			Destination: &source,
			Aliases:     []string{"s"},
		},
		&cli.StringFlag{
			Name:        "target",
			Usage:       "name of the target cluster",
			Destination: &target,
			Aliases:     []string{"t"},
		},
	}
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "status",
		Usage:  "status command",
		Action: status,
		Flags:  append(cmds.CommonFlags, statusFlags...),
	}
}

func status(clx *cli.Context) error {
	ctx := context.Background()

	restConfig, err := clientcmd.BuildConfigFromFlags("", cmds.Kubeconfig)
	if err != nil {
		return err
	}
	client, err := client.New(ctx, restConfig)
	if err != nil {
		return err
	}
	if source == "" || target == "" {
		logrus.Fatalf("source or target is not specified")
	}

	var clusters v3.ClusterList
	var sourceCluster, targetCluster *v3.Cluster
	if err := client.Clusters.List(ctx, "", &clusters, v1.ListOptions{}); err != nil {
		return err
	}

	for _, cluster := range clusters.Items {
		logrus.Debugf("check cluster: %s", cluster.Spec.DisplayName)
		if cluster.Spec.DisplayName == source {
			logrus.Debugf("cluster %s found", source)
			sourceCluster = cluster.DeepCopy()
		}
		if cluster.Spec.DisplayName == target {
			logrus.Debugf("cluster %s found", target)
			targetCluster = cluster.DeepCopy()
		}
	}
	if sourceCluster == nil || targetCluster == nil {
		logrus.Fatal("failed to find source or target cluster")
	}
	sc := &cluster.Cluster{
		Obj: sourceCluster,
	}
	tc := &cluster.Cluster{
		Obj: targetCluster,
	}
	if err := sc.Status(ctx, client, tc); err != nil {
		return err
	}
	return nil
}
