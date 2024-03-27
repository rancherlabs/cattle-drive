package cluster

import (
	"context"
	"fmt"
	"galal-hussein/cattle-drive/pkg/client"
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	creatorClusterOwner = "creator-cluster-owner"
	fleetDefaultOwner   = "fleet-default-owner"
)

type ClusterRoleTemplateBinding struct {
	Name     string
	Obj      *v3.ClusterRoleTemplateBinding
	Migrated bool
	Diff     bool
	// Description only exists for PRTB and CRTB
	Description string
}

func newCRTB(obj v3.ClusterRoleTemplateBinding, systemUser *v3.User) (*ClusterRoleTemplateBinding, bool) {
	if obj.Name == creatorClusterOwner || strings.Contains(obj.Name, fleetDefaultOwner) || strings.Contains(obj.Name, systemUser.Name) {
		// skipping crtb if its one of the default crtbs created for each cluster
		return nil, true
	}
	return &ClusterRoleTemplateBinding{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
		Diff:     false,
	}, false
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (c *ClusterRoleTemplateBinding) normalize() {
	// removing objectMeta and clusterName since crtb has no spec
	c.Obj.ObjectMeta = v1.ObjectMeta{}
	c.Obj.ClusterName = ""
}

func (c *ClusterRoleTemplateBinding) Mutate(tc *Cluster) {
	c.Obj.ClusterName = tc.Obj.Name
	c.Obj.SetName(c.Name)
	c.Obj.SetNamespace(tc.Obj.Name)
	c.Obj.SetFinalizers(nil)
	c.Obj.SetResourceVersion("")
	c.Obj.SetLabels(nil)
	for annotation := range c.Obj.Annotations {
		if strings.Contains(annotation, lifeCycleAnnotationPrefix) {
			delete(c.Obj.Annotations, annotation)
		}
	}
}
func (c *ClusterRoleTemplateBinding) SetDescription(ctx context.Context, client *client.Clients) error {
	var user v3.User
	userID := c.Obj.UserName
	if err := client.Users.Get(ctx, "", userID, &user, v1.GetOptions{}); err != nil {
		return err
	}
	name := user.DisplayName
	if name == "" {
		name = user.Username
	}
	c.Description = fmt.Sprintf("%s permission for user %s", c.Obj.RoleTemplateName, name)
	return nil
}
