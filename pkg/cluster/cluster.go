package cluster

import (
	"context"
	"galal-hussein/cattle-drive/pkg/client"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

type Cluster struct {
	Obj                         *v3.Cluster
	Projects                    map[string]Project
	ProjectRoleTemplateBindings map[string]ProjectRoleTemplateBinding
}

type Project struct {
	Obj      *v3.Project
	Migrated bool
	Diff     bool
}

type ProjectRoleTemplateBinding struct {
	Obj      *v3.ProjectRoleTemplateBinding
	Migrated bool
	Diff     bool
}

func (c *Cluster) Status(ctx context.Context, client *client.Clients, target *Cluster) error {

	if err := c.projectsStatus(ctx, client, target); err != nil {
		return err
	}

	// if project.Migrated {
	// 	if project.Diff {
	// 		fmt.Printf("- [%s] \u2718 (wrong spec) \n", projectName)
	// 	} else {
	// 		fmt.Printf("- [%s] \u2714 \n", projectName)
	// 	}
	// } else {
	// 	fmt.Printf("- [%s] \u2718 \n", projectName)
	// }

	return nil
}
