package cluster

import (
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

const (
	lifeCycleAnnotationPrefix = "lifecycle.cattle.io"
)

type Project struct {
	Name       string
	TargetName string
	Obj        *v3.Project
	Migrated   bool
	Diff       bool
	PRTBs      []*ProjectRoleTemplateBinding
	Namespaces []*Namespace
}

func newProject(obj v3.Project, prtbs []*ProjectRoleTemplateBinding, namespaces []*Namespace) *Project {
	return &Project{
		Name:       obj.Spec.DisplayName,
		Obj:        obj.DeepCopy(),
		Migrated:   false,
		Diff:       false,
		PRTBs:      prtbs,
		Namespaces: namespaces,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (p *Project) normalize() {
	p.Obj.Spec.ClusterName = ""
	p.Obj.Spec.ResourceQuota.UsedLimit = v3.ResourceQuotaLimit{}
}

// mutate will change the project object to be suitable for recreation to the target cluster
func (p *Project) Mutate(c *Cluster) {
	newProjectName := "p-" + generateName(5)
	p.Obj.Spec.ClusterName = c.Obj.Name
	p.Obj.Namespace = c.Obj.Name
	p.Obj.Status = v3.ProjectStatus{}
	p.Obj.SetFinalizers(nil)
	p.Obj.SetResourceVersion("")
	p.Obj.SetName(newProjectName)
	for annotation := range p.Obj.Annotations {
		if strings.Contains(annotation, lifeCycleAnnotationPrefix) {
			delete(p.Obj.Annotations, annotation)
		}
	}
}
