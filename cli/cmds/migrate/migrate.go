package migrate

import (
	"context"
	"errors"
	"fmt"
	"galal-hussein/cattle-drive/cli/cmds"
	"galal-hussein/cattle-drive/pkg/client"
	"galal-hussein/cattle-drive/pkg/cluster"
	"os"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	source       string
	target       string
	migrateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "source",
			Usage:       "name of the source cluster",
			Destination: &source,
			Aliases:     []string{"s"},
		},
		&cli.StringFlag{https://github.com/rancherlabs/cattle-drive/pull/9/conflict?name=cli%252Fcmds%252Fmigrate%252Fmigrate.go&ancestor_oid=38dfd762b78fafcff245ce52f24be4266a5b7830&base_oid=505ab5277a8713ea812d2ab04429496d3b10d096&head_oid=691928f200247163bfa64d6991e3d7dbf93445bf
			Name:        "target",
			Usage:       "name of the target cluster",
			Destination: &target,
			Aliases:     []string{"t"},
		},
	}
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "migrate",
		Usage:  "Migrate command",
		Action: migrate,
		Flags:  append(cmds.CommonFlags, migrateFlags...),
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
	// check for target rancher
	var (
		targetRancherConfig *rest.Config
		targetRancherClient *client.Clients
	)
	if cmds.TargetRancherConfig != "" {
		targetRancherConfig, err = clientcmd.BuildConfigFromFlags("", cmds.TargetRancherConfig)
		if err != nil {
			return err
		}
		targetRancherClient, err = client.New(ctx, targetRancherConfig)
		if err != nil {
			return err
		}
	}
	prefixMsg := fmt.Sprintf("initiating source [%s] and target [%s] clusters objects.. ", source, target)
	if cmds.TargetRancherConfig != "" {
		prefixMsg = fmt.Sprintf("initiating source cluster [%s] on rancher host [%s] and target cluster [%s] on rancher host [%s] ", source, restConfig.Host, target, targetRancherConfig.Host)
	}
	cmds.Spinner.Prefix = prefixMsg
	cmds.Spinner.Start()

	if source == "" || target == "" {
		return errors.New("source or target is not specified")
	}

	var clusters v3.ClusterList
	var sourceCluster, targetCluster *v3.Cluster
	if err := cl.Clusters.List(ctx, "", &clusters, v1.ListOptions{}); err != nil {
		return err
	}

	for _, cluster := range clusters.Items {
		if cluster.Spec.DisplayName == source {
			sourceCluster = cluster.DeepCopy()
		}
		if targetRancherClient == nil {
			if cluster.Spec.DisplayName == target {
				targetCluster = cluster.DeepCopy()
			}
		}
	}
	if targetRancherClient != nil {
		// find the target cluster on the target rancher environment
		if err := targetRancherClient.Clusters.List(ctx, "", &clusters, v1.ListOptions{}); err != nil {
			return err
		}
		for _, cluster := range clusters.Items {
			if cluster.Spec.DisplayName == target {
				targetCluster = cluster.DeepCopy()
			}
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
	if targetRancherClient != nil {
		tcConfig = *targetRancherConfig
		tcConfig.Host = targetRancherConfig.Host + "/k8s/clusters/" + targetCluster.Name
	}
	tcClient, err := client.New(ctx, &tcConfig)
	if err != nil {
		return err
	}
	tc := &cluster.Cluster{
		Obj:    targetCluster,
		Client: tcClient,
	}
	// check if the target cluster is in external rancher environment then set external rancher to true
	if targetRancherClient != nil {
		tc.ExternalRancher = true
		sc.ExternalRancher = true
	}
	if err := sc.Populate(ctx, cl); err != nil {
		return err
	}
	if targetRancherClient != nil {
		if err := tc.Populate(ctx, targetRancherClient); err != nil {
			return err
		}
	} else {
		if err := tc.Populate(ctx, cl); err != nil {
			return err
		}
	}
	if err := sc.Compare(ctx, tc); err != nil {
		return err
	}
	cmds.Spinner.Stop()
	if targetRancherClient != nil {
		return sc.Migrate(ctx, targetRancherClient, tc, os.Stdout)
	}
	return sc.Migrate(ctx, cl, tc, os.Stdout)
}
