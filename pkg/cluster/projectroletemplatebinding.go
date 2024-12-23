package cluster

import (
	"context"
	"fmt"
	"rancherlabs/cattle-drive/pkg/client"
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProjectRoleTemplateBinding struct {
	Name               string
	Obj                *v3.ProjectRoleTemplateBinding
	ProjectName        string
	ProjectDisplayName string
	Migrated           bool
	Diff               bool
	Description        string
}

func newPRTB(obj v3.ProjectRoleTemplateBinding, projectName, projectDisplayName string) *ProjectRoleTemplateBinding {
	return &ProjectRoleTemplateBinding{
		Name:               obj.Name,
		Obj:                obj.DeepCopy(),
		Migrated:           false,
		Diff:               false,
		ProjectName:        projectName,
		ProjectDisplayName: projectDisplayName,
	}
}

// normalize will remove unneeded fields in the spec to make it easier to compare
func (p *ProjectRoleTemplateBinding) normalize() {
	// removing objectMeta and projectName since prtb has no spec
	p.Obj.ObjectMeta = v1.ObjectMeta{}
	p.Obj.ProjectName = ""
}

func (p *ProjectRoleTemplateBinding) Mutate(clusterName, projectName string) {
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

func (p *ProjectRoleTemplateBinding) SetDescription(ctx context.Context, client *client.Clients) error {
	if p.Obj.UserName == "" && p.Obj.UserPrincipalName == "" {
		groupName := p.Obj.GroupName
		if groupName == "" {
			// handling external auth providers
			groupName = p.Obj.GroupPrincipalName
		}
		p.Description = fmt.Sprintf("%s permission for group %s", p.Obj.RoleTemplateName, groupName)
		return nil
	}
	// setting description for external users
	if p.Obj.UserPrincipalName != "" {
		p.Description = fmt.Sprintf("%s permission for user %s", p.Obj.RoleTemplateName, p.Obj.UserPrincipalName)
		return nil
	}
	var user v3.User
	userID := p.Obj.UserName
	if err := client.Users.Get(ctx, "", userID, &user, v1.GetOptions{}); err != nil {
		return err
	}
	name := user.DisplayName
	if name == "" {
		name = user.Username
	}
	p.Description = fmt.Sprintf("%s permission for user %s", p.Obj.RoleTemplateName, name)
	return nil
}
