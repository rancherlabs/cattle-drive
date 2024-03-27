package cluster

import (
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

type User struct {
	Name               string
	Obj                *v3.User
	Migrated           bool
	Diff               bool
	GlobalRoleBindings []*GlobalRoleBinding
}

func newUser(obj v3.User, grbList []*GlobalRoleBinding) *User {
	return &User{
		Name:               obj.Name,
		Obj:                obj.DeepCopy(),
		GlobalRoleBindings: grbList,
		Migrated:           false,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (u *User) normalize() {
}

func (u *User) Mutate() {
	u.Obj.SetName(u.Name)
	u.Obj.SetFinalizers(nil)
	u.Obj.SetResourceVersion("")
	u.Obj.SetLabels(nil)
	for annotation := range u.Obj.Annotations {
		if strings.Contains(annotation, lifeCycleAnnotationPrefix) {
			delete(u.Obj.Annotations, annotation)
		}
	}
}
