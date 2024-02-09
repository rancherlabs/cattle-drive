package cluster

import (
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProjectRoleTemplateBindings struct {
	Objs map[string]ProjectRoleTemplateBinding
}

type ProjectRoleTemplateBinding struct {
	Name     string
	Obj      *v3.ProjectRoleTemplateBinding
	Migrated bool
	Diff     bool
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (p *ProjectRoleTemplateBinding) normalize() {
	// removing objectMeta since prtb has no spec
	p.Obj.ObjectMeta = v1.ObjectMeta{}
	p.Obj.ProjectName = ""
}
