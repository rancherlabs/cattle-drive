package cluster

import (
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

type Project struct {
	Name     string
	Obj      *v3.Project
	Migrated bool
	Diff     bool
	PRTBs    []*ProjectRoleTemplateBinding
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (p *Project) normalize() {
	p.Obj.Spec.ClusterName = ""
}
