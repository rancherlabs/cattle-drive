package cluster

import (
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	projectIDLabelAnnotation = "field.cattle.io/projectId"
	statusAnnotation         = "cattle.io/status"
)

type Namespace struct {
	Name               string
	Obj                *v1.Namespace
	ProjectName        string
	ProjectDisplayName string
	Migrated           bool
	Diff               bool
}

func newNamespace(obj v1.Namespace, projectName, projectDisplayName string) *Namespace {
	return &Namespace{
		Name:               obj.Name,
		Obj:                obj.DeepCopy(),
		Migrated:           false,
		Diff:               false,
		ProjectName:        projectName,
		ProjectDisplayName: projectDisplayName,
	}
}

func (n Namespace) normalize() {
	delete(n.Obj.Annotations, statusAnnotation)
	n.Obj.Annotations[projectIDLabelAnnotation] = ""
	n.Obj.Labels[projectIDLabelAnnotation] = ""
	n.Obj.SetManagedFields(nil)
	n.Obj.SetCreationTimestamp(metav1.Time{})
	n.Obj.SetResourceVersion("")
	n.Obj.SetUID("")
}

func (n Namespace) Mutate(clusterName, projectName string) {
	n.Obj.Annotations[projectIDLabelAnnotation] = clusterName + ":" + projectName
	n.Obj.Labels[projectIDLabelAnnotation] = projectName
	n.Obj.SetFinalizers(nil)
	for annotation := range n.Obj.Annotations {
		if strings.Contains(annotation, lifeCycleAnnotationPrefix) {
			delete(n.Obj.Annotations, annotation)
		}
	}
}
