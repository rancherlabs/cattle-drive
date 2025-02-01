package cluster

import (
	"context"
	"errors"
	"fmt"
	"io"
	"rancherlabs/cattle-drive/pkg/client"
	"reflect"

	v1catalog "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	Obj             *v3.Cluster
	ToMigrate       ToMigrate
	SystemUser      *v3.User
	DefaultAdmin    *v3.User
	Client          *client.Clients
	ExternalRancher bool
}

type ToMigrate struct {
	Projects []*Project
	CRTBs    []*ClusterRoleTemplateBinding
	// apps related objects
	ClusterRepos []*ClusterRepo
	Apps         []*App
	Users        []*User
}

// Populate will fill in the objects to be migrated
func (c *Cluster) Populate(ctx context.Context, client *client.Clients) error {
	var (
		projects                    v3.ProjectList
		projectRoleTemplateBindings v3.ProjectRoleTemplateBindingList
		clusterRoleTemplateBindings v3.ClusterRoleTemplateBindingList
		users                       v3.UserList
		repos                       v1catalog.ClusterRepoList
		grbs                        v3.GlobalRoleBindingList
	)
	// systemUsers
	if err := client.Users.List(ctx, "", &users, v1.ListOptions{}); err != nil {
		return err
	}
	usersList := []*User{}
	for _, user := range users.Items {
		for _, principalID := range user.PrincipalIDs {
			if principalID == "system://"+c.Obj.Name {
				c.SystemUser = user.DeepCopy()
				break
			}
		}
		if c.ExternalRancher {
			if user.Name == c.Obj.Annotations["field.cattle.io/creatorId"] {
				// use the cluster creator user as the default admin for any new project
				c.DefaultAdmin = user.DeepCopy()
			}
			if user.Username == "admin" || user.Username == "" {
				continue
			}
			var grbList []*GlobalRoleBinding
			if err := client.GlobalRoleBindings.List(ctx, "", &grbs, v1.ListOptions{}); err != nil {
				return err
			}
			for _, grb := range grbs.Items {
				if grb.UserName == user.Name {
					newGRB := newGRB(grb)
					newGRB.normalize()
					newGRB.SetDescription(user)
					grbList = append(grbList, newGRB)
				}
			}
			// populate users
			u := newUser(user, grbList)
			u.normalize()
			usersList = append(usersList, u)

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
	for _, p := range projects.Items {
		// skip default projects before listing their prtb or roles
		if p.Spec.DisplayName == "Default" || p.Spec.DisplayName == "System" {
			continue
		}
		// prtbs
		if err := client.ProjectRoleTemplateBindings.List(ctx, p.Name, &projectRoleTemplateBindings, v1.ListOptions{}); err != nil {
			return err
		}
		prtbList := []*ProjectRoleTemplateBinding{}
		for _, item := range projectRoleTemplateBindings.Items {
			if item.Name == "creator-project-owner" || item.Name == "creator-project-member" {
				continue
			}
			prtb := newPRTB(item, "", p.Spec.DisplayName)
			prtb.normalize()
			if err := prtb.SetDescription(ctx, client); err != nil {
				return err
			}
			prtbList = append(prtbList, prtb)
		}
		nsList := []*Namespace{}
		for _, ns := range namespaces.Items {
			if projectID, ok := ns.Labels[projectIDLabelAnnotation]; ok && projectID == p.Name {
				n := newNamespace(ns, "", p.Spec.DisplayName)
				n.normalize()
				nsList = append(nsList, n)
			}
		}
		p := newProject(p, prtbList, nsList)
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
		if err := crtb.SetDescription(ctx, client); err != nil {
			return err
		}
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
		Users:        usersList,
	}
	return nil
}

// Compare will compare between objects of downstream source cluster and target cluster
func (c *Cluster) Compare(ctx context.Context, tc *Cluster) error {
	// users
	for _, sUser := range c.ToMigrate.Users {
		for _, tUser := range tc.ToMigrate.Users {
			if sUser.Name == tUser.Name && sUser.Obj.Username == tUser.Obj.Username {
				sUser.Migrated = true
				for _, sGRB := range sUser.GlobalRoleBindings {
					for _, tGRB := range tUser.GlobalRoleBindings {
						if sGRB.Name == tGRB.Name {
							sGRB.Migrated = true
						}
					}
				}
			}
		}
	}
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
					for _, sPRTB := range sProject.PRTBs {
						sPRTB.ProjectName = tProject.Obj.Name
					}
					for _, ns := range sProject.Namespaces {
						ns.ProjectName = tProject.Obj.Name
					}
				}
				// now we check for prtbs related to that project
				for _, sPrtb := range sProject.PRTBs {
					for _, tPrtb := range tProject.PRTBs {
						if sPrtb.Name == tPrtb.Name {
							sPrtb.Migrated = true
							if !sPrtb.Compare(tPrtb) {
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

func (c *Cluster) Status(ctx context.Context) error {
	if c.ExternalRancher {
		fmt.Printf("Users status:\n")
		for _, u := range c.ToMigrate.Users {
			print(u.Obj.Username, u.Migrated, u.Diff, 0)
			if len(u.GlobalRoleBindings) > 0 {
				fmt.Printf("  -> user permissions:\n")
			}
			for _, grb := range u.GlobalRoleBindings {
				print(grb.Name+": "+grb.Description, grb.Migrated, grb.Diff, 1)
			}
		}
	}

	fmt.Printf("Project status:\n")
	for _, p := range c.ToMigrate.Projects {
		print(p.Name, p.Migrated, p.Diff, 0)
		if len(p.PRTBs) > 0 {
			fmt.Printf("  -> users permissions:\n")
		}
		for _, prtb := range p.PRTBs {
			print(prtb.Name+": "+prtb.Description, prtb.Migrated, prtb.Diff, 1)
		}
		if len(p.Namespaces) > 0 {
			fmt.Printf("  -> namespaces:\n")
		}
		for _, ns := range p.Namespaces {
			print(ns.Name, ns.Migrated, ns.Diff, 1)
		}

	}
	fmt.Printf("Cluster users permissions:\n")
	for _, crtb := range c.ToMigrate.CRTBs {
		print(crtb.Name+": "+crtb.Description, crtb.Migrated, crtb.Diff, 0)
	}
	fmt.Printf("Catalog repos:\n")
	for _, repo := range c.ToMigrate.ClusterRepos {
		print(repo.Name, repo.Migrated, repo.Diff, 0)
	}
	return nil
}

func (c *Cluster) Migrate(ctx context.Context, client *client.Clients, tc *Cluster, w io.Writer) error {
	fmt.Fprintf(w, "Migrating Objects from cluster [%s] to cluster [%s]:\n", c.Obj.Spec.DisplayName, tc.Obj.Spec.DisplayName)
	// users
	if c.ExternalRancher {
		for _, u := range c.ToMigrate.Users {
			if !u.Migrated {
				fmt.Fprintf(w, "- migrating User [%s]... ", u.Obj.Username)

				u.Mutate()
				if err := client.Users.Create(ctx, "", u.Obj, nil, v1.CreateOptions{}); err != nil {
					return err
				}
				// migrating all grbs for this user
				for _, grb := range u.GlobalRoleBindings {
					grb.Mutate()
					if err := client.GlobalRoleBindings.Create(ctx, "", grb.Obj, nil, v1.CreateOptions{}); err != nil {
						return err
					}
				}
				fmt.Fprintf(w, "Done.\n")
			}
		}
	}

	for _, p := range c.ToMigrate.Projects {
		if !p.Migrated {
			fmt.Fprintf(w, "- migrating Project [%s]... ", p.Name)
			p.Mutate(tc)
			if err := client.Projects.Create(ctx, tc.Obj.Name, p.Obj, nil, v1.CreateOptions{}); err != nil {
				return err
			}
			// set ProjectName for all ns and prtbs for this project
			for _, sPRTB := range p.PRTBs {
				sPRTB.ProjectName = p.Obj.Name
			}
			for _, ns := range p.Namespaces {
				ns.ProjectName = p.Obj.Name
			}
			fmt.Fprintf(w, "Done.\n")
		}

		for _, prtb := range p.PRTBs {
			if !prtb.Migrated {
				fmt.Fprintf(w, "  - migrating PRTB [%s]... ", prtb.Name)
				// check if the user exists first in case of external rancher
				userID := prtb.Obj.UserName
				var user v3.User
				if err := client.Users.Get(ctx, "", userID, &user, v1.GetOptions{}); err != nil {
					if apierrors.IsNotFound(err) {
						return errors.New("user " + userID + " does not exists, please migrate user first")
					}
				}
				prtb.Mutate(tc.Obj.Name, prtb.ProjectName)
				if err := client.ProjectRoleTemplateBindings.Create(ctx, prtb.ProjectName, prtb.Obj, nil, v1.CreateOptions{}); err != nil {
					return err
				}
				fmt.Fprintf(w, "Done.\n")
			}
		}
		for _, ns := range p.Namespaces {
			if !ns.Migrated {
				fmt.Fprintf(w, "  - migrating Namespace [%s]... ", ns.Name)
				ns.Mutate(tc.Obj.Name, ns.ProjectName)
				if _, err := tc.Client.Namespace.Create(ns.Obj); err != nil {
					return err
				}
				fmt.Fprintf(w, "Done.\n")
			}
		}
	}
	for _, crtb := range c.ToMigrate.CRTBs {
		if !crtb.Migrated {
			fmt.Fprintf(w, "- migrating CRTB [%s]... ", crtb.Name)
			// check if the user exists first in case of external rancher
			userID := crtb.Obj.UserName
			var user v3.User
			if err := client.Users.Get(ctx, "", userID, &user, v1.GetOptions{}); err != nil {
				if apierrors.IsNotFound(err) {
					return errors.New("user " + userID + " does not exists, please migrate user first")
				}
			}

			crtb.Mutate(tc)
			if err := client.ClusterRoleTemplateBindings.Create(ctx, tc.Obj.Name, crtb.Obj, nil, v1.CreateOptions{}); err != nil {
				return err
			}
			fmt.Fprintf(w, "Done.\n")
		}
	}
	// catalog repos
	for _, repo := range c.ToMigrate.ClusterRepos {
		if !repo.Migrated {
			fmt.Fprintf(w, "- migrating catalog repo [%s]... ", repo.Name)
			repo.Mutate()
			if err := tc.Client.ClusterRepos.Create(ctx, tc.Obj.Name, repo.Obj, nil, v1.CreateOptions{}); err != nil {
				return err
			}
			fmt.Fprintf(w, "Done.\n")
		}
	}

	return nil
}

func NewProjectName(ctx context.Context, targetClusterName, oldProjectName string, client *client.Clients) (string, error) {
	var projects v3.ProjectList
	if err := client.Projects.List(ctx, targetClusterName, &projects, v1.ListOptions{}); err != nil {
		return "", err
	}
	for _, project := range projects.Items {
		if oldProjectName == project.Spec.DisplayName {
			return project.Name, nil
		}
	}
	return "", errors.New("failed to find project with the name " + oldProjectName)
}
