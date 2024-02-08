package cluster

import (
	"context"
	"galal-hussein/cattle-drive/pkg/client"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	Obj      *v3.Cluster
	Projects map[string]*Project
}

type Project struct {
	Obj      *v3.Project
	Migrated bool
}

func (c *Cluster) Status(ctx context.Context, client *client.Clients, target *Cluster) error {
	// get all projects for the source and target cluster
	var (
		sourceProjects v3.ProjectList
		targetProjects v3.ProjectList
	)
	if err := client.Projects.List(ctx, c.Obj.Name, &sourceProjects, v1.ListOptions{}); err != nil {
		return err
	}

	if err := client.Projects.List(ctx, target.Obj.Name, &targetProjects, v1.ListOptions{}); err != nil {
		return err
	}

	for _, i := range sourceProjects.Items {
		logrus.Infof("project %s", i.Spec.DisplayName)
	}

	return nil
}

// func toMap(objects interface{}) (map[string]runtime.Object, error) {
// 	m := make(map[string]runtime.Object)
// 	objs, ok := objects.([]runtime.Object)
// 	if !ok {
// 		return nil, errors.New("error")
// 	}
// 	for _, item := range objs {
// 		objCopy := item.DeepCopyObject()

// 		objMeta, err := meta.Accessor(objCopy)
// 		if err != nil {
// 			return nil, err
// 		}

// 		m[objMeta.GetName()] = objCopy
// 	}
// 	return m, nil
// }
