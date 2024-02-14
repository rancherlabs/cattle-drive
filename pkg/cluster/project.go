package cluster

import (
	"galal-hussein/cattle-drive/pkg/util"
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

const (
	lifeCycleAnnotationPrefix = "lifecycle.cattle.io"
)

type Project struct {
	Name         string
	Obj          *v3.Project
	Migrated     bool
	Diff         bool
	PRTBs        []*ProjectRoleTemplateBinding
	Roles        []*Role
	RoleBindings []*RoleBinding
}

func newProject(obj v3.Project, prtbs []*ProjectRoleTemplateBinding, roles []*Role, roleBindings []*RoleBinding) *Project {
	return &Project{
		Name:         obj.Spec.DisplayName,
		Obj:          obj.DeepCopy(),
		Migrated:     false,
		Diff:         false,
		PRTBs:        prtbs,
		Roles:        roles,
		RoleBindings: roleBindings,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (p *Project) normalize() {
	p.Obj.Spec.ClusterName = ""
	p.Obj.Spec.ResourceQuota.UsedLimit = v3.ResourceQuotaLimit{}
}

// mutate will change the project object to be suitable for recreation to the target cluster
func (p *Project) mutate(c *Cluster) {
	p.Obj.Spec.ClusterName = c.Obj.Name
	p.Obj.Namespace = c.Obj.Name
	p.Obj.Status = v3.ProjectStatus{}
	p.Obj.SetFinalizers(nil)
	p.Obj.SetResourceVersion("")
	p.Obj.SetName("p-" + util.GenerateName(5))
	for annotation := range p.Obj.Annotations {
		if strings.Contains(annotation, lifeCycleAnnotationPrefix) {
			delete(p.Obj.Annotations, annotation)
		}
	}
}
