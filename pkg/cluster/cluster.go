package cluster

import (
	"context"
	"fmt"
	"galal-hussein/cattle-drive/pkg/client"
	"reflect"

	v1catalog "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	Obj        *v3.Cluster
	ToMigrate  ToMigrate
	SystemUser *v3.User
	Client     *client.Clients
}

type ToMigrate struct {
	Projects []*Project
	CRTBs    []*ClusterRoleTemplateBinding
	// apps related objects
	ClusterRepos []*ClusterRepo
	Apps         []*App
}

// Populate will fill in the objects to be migrated
func (c *Cluster) Populate(ctx context.Context, client *client.Clients) error {
	var (
		projects                    v3.ProjectList
		projectRoleTemplateBindings v3.ProjectRoleTemplateBindingList
		clusterRoleTemplateBindings v3.ClusterRoleTemplateBindingList
		users                       v3.UserList
		repos                       v1catalog.ClusterRepoList
	)
	// systemUsers
	if err := client.Users.List(ctx, "", &users, v1.ListOptions{}); err != nil {
		return err
	}
	for _, user := range users.Items {
		for _, principalID := range user.PrincipalIDs {
			if principalID == "system://"+c.Obj.Name {
				c.SystemUser = user.DeepCopy()
				break
			}
		}
	}
	// namespaces
	namespaces, err := c.Client.Namespace.List(v1.ListOptions{})
	if err != nil {
		return err
	}
	// projects
	if err := client.Projects.List(ctx, c.Obj.Name, &projects, v1.ListOptions{}); err != nil {
		return err
	}
	pList := []*Project{}
	for _, item := range projects.Items {
		// skip default projects before listing their prtb or roles
		if item.Spec.DisplayName == "Default" || item.Spec.DisplayName == "System" {
			continue
		}
		// prtbs
		if err := client.ProjectRoleTemplateBindings.List(ctx, item.Name, &projectRoleTemplateBindings, v1.ListOptions{}); err != nil {
			return err
		}
		prtbList := []*ProjectRoleTemplateBinding{}
		for _, item := range projectRoleTemplateBindings.Items {
			if item.Name == "creator-project-owner" {
				continue
			}
			prtb := newPRTB(item)
			prtb.normalize()
			prtbList = append(prtbList, prtb)
		}
		nsList := []*Namespace{}
		for _, ns := range namespaces.Items {
			if projectID, ok := ns.Labels[projectIDLabelAnnotation]; ok && projectID == item.Name {
				n := newNamespace(ns)
				n.normalize()
				nsList = append(nsList, n)
			}
		}
		p := newProject(item, prtbList, nsList)
		p.normalize()
		pList = append(pList, p)
	}

	crtbList := []*ClusterRoleTemplateBinding{}
	if err := client.ClusterRoleTemplateBindings.List(ctx, c.Obj.Name, &clusterRoleTemplateBindings, v1.ListOptions{}); err != nil {
		return err
	}
	for _, item := range clusterRoleTemplateBindings.Items {
		crtb, isDefault := newCRTB(item, c.SystemUser)
		if isDefault {
			continue
		}
		crtb.normalize()
		crtbList = append(crtbList, crtb)
	}
	// apps
	// cluster repos
	reposList := []*ClusterRepo{}
	if err := c.Client.ClusterRepos.List(ctx, "", &repos, v1.ListOptions{}); err != nil {
		return err
	}
	for _, item := range repos.Items {
		repo, isDefault := newClusterRepo(item)
		if isDefault {
			continue
		}
		repo.normalize()
		reposList = append(reposList, repo)

	}

	c.ToMigrate = ToMigrate{
		Projects:     pList,
		CRTBs:        crtbList,
		ClusterRepos: reposList,
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
				} else {
					// its critical to adjust the project name here because its used in different other objects ns/prtbs
					sProject.Obj.Name = tProject.Obj.Name
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
				// namespaces
				for _, ns := range sProject.Namespaces {
					for _, tns := range tProject.Namespaces {
						if ns.Name == tns.Name {
							ns.Migrated = true
							if !reflect.DeepEqual(ns.Obj, tns.Obj) {
								ns.Diff = true
							}
						}
					}
				}
			}
		}
	}

	// crtbs
	for _, sCrtb := range c.ToMigrate.CRTBs {
		for _, tCrtb := range tc.ToMigrate.CRTBs {
			if sCrtb.Name == tCrtb.Name {
				sCrtb.Migrated = true
				if !reflect.DeepEqual(sCrtb.Obj, tCrtb.Obj) {
					sCrtb.Diff = true
					break
				}
			}
		}
	}

	for _, sRepo := range c.ToMigrate.ClusterRepos {
		for _, tRepo := range tc.ToMigrate.ClusterRepos {
			if sRepo.Name == tRepo.Name {
				sRepo.Migrated = true
				if !reflect.DeepEqual(sRepo.Obj.Spec, tRepo.Obj.Spec) {
					sRepo.Diff = true
					break
				}
			}
		}
	}
	return nil
}

func (c *Cluster) Status(ctx context.Context, client *client.Clients) error {
	fmt.Printf("Project status:\n")
	for _, p := range c.ToMigrate.Projects {
		print(p.Name, p.Migrated, p.Diff, 0)
		if p.Migrated && !p.Diff {
			fmt.Printf("  -> users permissions:\n")
			for _, prtb := range p.PRTBs {
				print(prtb.Name, prtb.Migrated, prtb.Diff, 1)
			}
			fmt.Printf("  -> namespaces:\n")
			for _, ns := range p.Namespaces {
				print(ns.Name, ns.Migrated, ns.Diff, 1)
			}
		}
	}
	fmt.Printf("Cluster users permissions:\n")
	for _, crtb := range c.ToMigrate.CRTBs {
		print(crtb.Name, crtb.Migrated, crtb.Diff, 0)
	}
	fmt.Printf("Catalog repos:\n")
	for _, repo := range c.ToMigrate.ClusterRepos {
		print(repo.Name, repo.Migrated, repo.Diff, 0)
	}
	return nil
}

func (c *Cluster) Migrate(ctx context.Context, client *client.Clients, tc *Cluster) error {
	fmt.Printf("Migrating Objects from cluster [%s] to cluster [%s]:\n", c.Obj.Spec.DisplayName, tc.Obj.Spec.DisplayName)
	for _, p := range c.ToMigrate.Projects {
		if !p.Migrated {
			fmt.Printf("- migrating Project [%s]... ", p.Name)
			p.mutate(tc)
			if err := client.Projects.Create(ctx, tc.Obj.Name, p.Obj, nil, v1.CreateOptions{}); err != nil {
				return err
			}
			fmt.Printf("Done.\n")
		}
		for _, prtb := range p.PRTBs {
			if !prtb.Migrated {
				fmt.Printf("  - migrating PRTB [%s]... ", prtb.Name)
				prtb.mutate(tc.Obj.Name, p.Obj.Name)
				if err := client.ProjectRoleTemplateBindings.Create(ctx, p.Obj.Name, prtb.Obj, nil, v1.CreateOptions{}); err != nil {
					return err
				}
				fmt.Printf("Done.\n")
			}
		}
		for _, ns := range p.Namespaces {
			if !ns.Migrated {
				fmt.Printf("  - migrating Namespace [%s]... ", ns.Name)
				ns.mutate(tc.Obj.Name, p.Obj.Name)
				if _, err := tc.Client.Namespace.Create(ns.Obj); err != nil {
					return err
				}
				fmt.Printf("Done.\n")
			}
		}
	}
	for _, crtb := range c.ToMigrate.CRTBs {
		if !crtb.Migrated {
			fmt.Printf("- migrating CRTB [%s]... ", crtb.Name)
			crtb.mutate(tc)
			if err := client.ClusterRoleTemplateBindings.Create(ctx, tc.Obj.Name, crtb.Obj, nil, v1.CreateOptions{}); err != nil {
				return err
			}
			fmt.Printf("Done.\n")
		}
	}
	// catalog repos
	for _, repo := range c.ToMigrate.ClusterRepos {
		if !repo.Migrated {
			fmt.Printf("- migrating catalog repo [%s]... ", repo.Name)
			repo.mutate()
			if err := tc.Client.ClusterRepos.Create(ctx, tc.Obj.Name, repo.Obj, nil, v1.CreateOptions{}); err != nil {
				return err
			}
			fmt.Printf("Done.\n")
		}
	}

	return nil
}

func (c *Cluster) Interactive(ctx context.Context, client *client.Clients, tc *Cluster) error {

	return nil
}
