package cluster

import (
	v1 "k8s.io/api/rbac/v1"
)

type Role struct {
	Name     string
	Obj      *v1.Role
	Migrated bool
	Diff     bool
}

func newRole(obj v1.Role) *Role {
	return &Role{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
		Diff:     false,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (r *Role) normalize() {
	//
}
