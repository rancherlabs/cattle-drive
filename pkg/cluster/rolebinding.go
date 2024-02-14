package cluster

import v1 "k8s.io/api/rbac/v1"

type RoleBinding struct {
	Name     string
	Obj      *v1.RoleBinding
	Migrated bool
	Diff     bool
}

func newRoleBinding(obj v1.RoleBinding) *RoleBinding {
	return &RoleBinding{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
		Diff:     false,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (r *RoleBinding) normalize() {
	//
}
