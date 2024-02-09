package cluster

import (
	"context"
	"fmt"
	"galal-hussein/cattle-drive/pkg/client"
	"reflect"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	Obj       *v3.Cluster
	ToMigrate ToMigrate
}

type ToMigrate struct {
	Projects []*Project
}

// Populate will fill in the objects to be migrated
func (c *Cluster) Populate(ctx context.Context, client *client.Clients) error {
	var (
		projects                    v3.ProjectList
		projectRoleTemplateBindings v3.ProjectRoleTemplateBindingList
	)
	if err := client.Projects.List(ctx, c.Obj.Name, &projects, v1.ListOptions{}); err != nil {
		return err
	}
	pList := []*Project{}
	for _, item := range projects.Items {
		if err := client.ProjectRoleTemplateBindings.List(ctx, item.Name, &projectRoleTemplateBindings, v1.ListOptions{}); err != nil {
			return err
		}
		prtbList := []*ProjectRoleTemplateBinding{}
		for _, item := range projectRoleTemplateBindings.Items {
			prtb := &ProjectRoleTemplateBinding{
				Name:     item.Name,
				Obj:      item.DeepCopy(),
				Migrated: false,
				Diff:     false,
			}
			prtb.normalize()
			prtbList = append(prtbList, prtb)
		}
		p := &Project{
			Name:     item.Spec.DisplayName,
			Obj:      item.DeepCopy(),
			PRTBs:    prtbList,
			Migrated: false,
			Diff:     false,
		}
		p.normalize()
		pList = append(pList, p)
	}

	c.ToMigrate = ToMigrate{
		Projects: pList,
	}
	return nil
}

// Compare will compare between objects of downstream source cluster and target cluster
func (c *Cluster) Compare(ctx context.Context, client *client.Clients, tc *Cluster) error {
	// projects
	for _, sProject := range c.ToMigrate.Projects {
		for _, tProject := range tc.ToMigrate.Projects {
			if sProject.Name == tProject.Name {
				sProject.Migrated = true
				if !reflect.DeepEqual(sProject.Obj.Spec, tProject.Obj.Spec) {
					sProject.Diff = true
					break
				}
				// now we check for prtbs related to that project

				for _, sPrtb := range sProject.PRTBs {
					for _, tPrtb := range tProject.PRTBs {
						if sPrtb.Name == tPrtb.Name {
							sPrtb.Migrated = true
							if !reflect.DeepEqual(sPrtb.Obj, tPrtb.Obj) {
								sPrtb.Diff = true
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (c *Cluster) Status(ctx context.Context, client *client.Clients) error {
	fmt.Printf("Project status:\n")
	for _, p := range c.ToMigrate.Projects {
		if p.Migrated {
			if p.Diff {
				fmt.Printf("- [%s] \u2718 (wrong spec) \n", p.Name)
			} else {
				fmt.Printf("- [%s] \u2714 \n", p.Name)
				fmt.Printf("  prtbs:\n")
				for _, prtb := range p.PRTBs {
					if prtb.Migrated {
						if prtb.Diff {
							fmt.Printf("  - [%s] \u2718 (wrong fields) \n", prtb.Name)
						} else {
							fmt.Printf("  - [%s] \u2714 \n", prtb.Name)
						}
					} else {
						fmt.Printf("  - [%s] \u2718 \n", prtb.Name)
					}
				}
			}
		} else {
			fmt.Printf("- [%s] \u2718 \n", p.Name)
		}
	}
	return nil
}
