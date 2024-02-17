package interactive

import (
	"context"
	"galal-hussein/cattle-drive/cli/cmds"
	"galal-hussein/cattle-drive/pkg/client"
	"galal-hussein/cattle-drive/pkg/cluster"
	"galal-hussein/cattle-drive/pkg/cluster/tui"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	source           string
	target           string
	interactiveFlags = []cli.Flag{
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
		Name:   "interactive",
		Usage:  "Interactive command",
		Action: migrate,
		Flags:  append(cmds.CommonFlags, interactiveFlags...),
	}
}

func migrate(clx *cli.Context) error {
	ctx := context.Background()

	restConfig, err := clientcmd.BuildConfigFromFlags("", cmds.Kubeconfig)
	if err != nil {
		return err
	}
	cl, err := client.New(ctx, restConfig)
	if err != nil {
		return err
	}
	if source == "" || target == "" {
		logrus.Fatalf("source or target is not specified")
	}

	var clusters v3.ClusterList
	var sourceCluster, targetCluster *v3.Cluster
	if err := cl.Clusters.List(ctx, "", &clusters, v1.ListOptions{}); err != nil {
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
	// initiate client for the cluster
	scConfig := *restConfig
	scConfig.Host = restConfig.Host + "/k8s/clusters/" + sourceCluster.Name
	scClient, err := client.New(ctx, &scConfig)
	if err != nil {
		return err
	}
	sc := &cluster.Cluster{
		Obj:    sourceCluster,
		Client: scClient,
	}
	tcConfig := *restConfig
	tcConfig.Host = restConfig.Host + "/k8s/clusters/" + targetCluster.Name
	tcClient, err := client.New(ctx, &tcConfig)
	if err != nil {
		return err
	}
	tc := &cluster.Cluster{
		Obj:    targetCluster,
		Client: tcClient,
	}
	if err := sc.Populate(ctx, cl); err != nil {
		return err
	}
	if err := tc.Populate(ctx, cl); err != nil {
		return err
	}
	if err := sc.Compare(ctx, cl, tc); err != nil {
		return err
	}

	return tui.StartTea(sc, tc, cl)
}
