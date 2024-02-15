package cluster

import (
	"strings"

	v1 "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
)

type ClusterRepo struct {
	Name     string
	Obj      *v1.ClusterRepo
	Migrated bool
	Diff     bool
}

func newClusterRepo(obj v1.ClusterRepo) (*ClusterRepo, bool) {
	if strings.HasPrefix(obj.Name, "rancher-") {
		return nil, true
	}
	return &ClusterRepo{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
		Diff:     false,
	}, false
}

func (c *ClusterRepo) normalize() {
}

func (c *ClusterRepo) mutate() {
	c.Obj.SetFinalizers(nil)
	c.Obj.SetResourceVersion("")
	c.Obj.Status = v1.RepoStatus{}
}
