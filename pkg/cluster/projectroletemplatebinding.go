package cluster

import (
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProjectRoleTemplateBinding struct {
	Name     string
	Obj      *v3.ProjectRoleTemplateBinding
	Migrated bool
	Diff     bool
}

func newPRTB(obj v3.ProjectRoleTemplateBinding) *ProjectRoleTemplateBinding {
	return &ProjectRoleTemplateBinding{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
		Diff:     false,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (p *ProjectRoleTemplateBinding) normalize() {
	// removing objectMeta and projectName since prtb has no spec
	p.Obj.ObjectMeta = v1.ObjectMeta{}
	p.Obj.ProjectName = ""
}

func (p *ProjectRoleTemplateBinding) mutate(clusterName, projectName string) {
	p.Obj.ProjectName = clusterName + ":" + projectName
	p.Obj.SetName(p.Name)
	p.Obj.SetNamespace(projectName)
	p.Obj.SetFinalizers(nil)
	p.Obj.SetResourceVersion("")
	p.Obj.SetLabels(nil)
	for annotation := range p.Obj.Annotations {
		if strings.Contains(annotation, lifeCycleAnnotationPrefix) {
			delete(p.Obj.Annotations, annotation)
		}
	}
}
