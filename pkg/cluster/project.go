package cluster

import (
	"context"
	"galal-hussein/cattle-drive/pkg/client"
	"reflect"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Cluster) projectsStatus(ctx context.Context, client *client.Clients, target *Cluster) error {
	// get all projects for the source and target cluster
	var (
		sourceProjects v3.ProjectList
		targetProjects v3.ProjectList
	)
	if err := client.Projects.List(ctx, c.Obj.Name, &sourceProjects, v1.ListOptions{}); err != nil {
		return err
	}
	c.projectsToMap(sourceProjects)

	if err := client.Projects.List(ctx, target.Obj.Name, &targetProjects, v1.ListOptions{}); err != nil {
		return err
	}
	target.projectsToMap(targetProjects)

	for projectName, project := range c.Projects {
		if projectName == "Default" || projectName == "System" {
			continue
		}
		if targetProject, ok := target.Projects[projectName]; ok {
			// project exists in target cluster, comparing specs
			project.Obj.Spec.ClusterName = ""
			targetProject.Obj.Spec.ClusterName = ""
			if reflect.DeepEqual(project.Obj.Spec, targetProject.Obj.Spec) {
				project.Migrated = true
			} else {
				// project migrated but with different spec than
				project.Migrated = true
				project.Diff = true
			}
		}
	}
	return nil
}

func (c *Cluster) projectsToMap(projectList v3.ProjectList) {
	c.Projects = make(map[string]Project, len(projectList.Items))
	for _, project := range projectList.Items {
		c.Projects[project.Spec.DisplayName] = Project{
			Obj:      project.DeepCopy(),
			Migrated: false,
		}
	}
}
