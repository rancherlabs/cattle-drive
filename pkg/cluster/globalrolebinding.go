package cluster

import (
	"fmt"
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

type GlobalRoleBinding struct {
	Name        string
	Obj         *v3.GlobalRoleBinding
	Description string
	Migrated    bool
	Diff        bool
}

func newGRB(obj v3.GlobalRoleBinding) *GlobalRoleBinding {
	return &GlobalRoleBinding{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (g *GlobalRoleBinding) normalize() {
	for _, or := range g.Obj.OwnerReferences {
		or.UID = ""
	}
}

func (g *GlobalRoleBinding) Mutate() {
	g.Obj.SetName(g.Name)
	g.Obj.SetFinalizers(nil)
	g.Obj.SetResourceVersion("")
	g.Obj.SetLabels(nil)
	for annotation := range g.Obj.Annotations {
		if strings.Contains(annotation, lifeCycleAnnotationPrefix) {
			delete(g.Obj.Annotations, annotation)
		}
		if strings.Contains(annotation, "authz.management.cattle.io/crb-name") {
			delete(g.Obj.Annotations, annotation)
		}
	}
}

func (g *GlobalRoleBinding) SetDescription(user v3.User) error {

	g.Description = fmt.Sprintf("%s permission for user %s", g.Obj.GlobalRoleName, user.Username)
	return nil
}
