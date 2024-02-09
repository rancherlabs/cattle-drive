package cluster

import (
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterRoleTemplateBinding struct {
	Name     string
	Obj      *v3.ClusterRoleTemplateBinding
	Migrated bool
	Diff     bool
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (c *ClusterRoleTemplateBinding) normalize() {
	// removing objectMeta and projectName since prtb has no spec
	c.Obj.ObjectMeta = v1.ObjectMeta{}
	c.Obj.ClusterName = ""
}
