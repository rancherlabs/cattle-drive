package cluster

import (
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	creatorClusterOwner = "creator-cluster-owner"
	fleetDefaultOwner   = "fleet-default-owner"
)

type ClusterRoleTemplateBinding struct {
	Name     string
	Obj      *v3.ClusterRoleTemplateBinding
	Migrated bool
	Diff     bool
}

func newCRTB(obj v3.ClusterRoleTemplateBinding, systemUser *v3.User) (*ClusterRoleTemplateBinding, bool) {
	if obj.Name == creatorClusterOwner || strings.Contains(obj.Name, fleetDefaultOwner) || strings.Contains(obj.Name, systemUser.Name) {
		// skipping crtb if its one of the default crtbs created for each cluster
		return nil, true
	}
	return &ClusterRoleTemplateBinding{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
		Diff:     false,
	}, false
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (c *ClusterRoleTemplateBinding) normalize() {
	// removing objectMeta and clusterName since crtb has no spec
	c.Obj.ObjectMeta = v1.ObjectMeta{}
	c.Obj.ClusterName = ""
}
