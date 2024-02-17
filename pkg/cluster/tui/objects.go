package tui

import (
	"context"
	"galal-hussein/cattle-drive/pkg/cluster"
	"galal-hussein/cattle-drive/pkg/cluster/tui/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	errMsg struct{ error }
)

type Objects struct {
	mode     mode
	list     list.Model
	quitting bool
	spinner  spinner.Model
	progress progress.Model
}

// Init run any intial IO on program start
func (m Objects) Init() tea.Cmd {
	return nil
}

func InitObjects(i item, p *tea.Program) *Objects {
	var title string
	items := []list.Item{}
	pr := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	m := Objects{mode: nav, progress: pr, spinner: s}
	switch i.objType {
	case "project":
		// in case of individual project then we will list namespaces and prtbs
		project := i.obj.(*cluster.Project)
		title = "Project [" + project.Name + "]"
		items = []list.Item{
			item{title: "Project User Permissions", objType: constants.PRTBsType, obj: project.PRTBs},
			item{title: "Namespaces", objType: constants.NamespacesType, obj: project.Namespaces},
		}
	case "namespaces":
		namespaces := i.obj.([]*cluster.Namespace)
		title = "Namespaces for Project "
		if len(namespaces) > 0 {
			title = title + "[" + namespaces[0].ProjectDisplayName + "]"
		}
		for _, ns := range namespaces {
			t, status := status(ns.Name, ns.Migrated, ns.Diff)
			i := item{title: t, status: status, objType: constants.NamespaceType, obj: ns}
			items = append(items, i)
		}
	case "prtbs":
		prtbs := i.obj.([]*cluster.ProjectRoleTemplateBinding)
		title = "User Permissions for Project "
		if len(prtbs) > 0 {
			title = title + "[" + prtbs[0].ProjectDisplayName + "]"
		}
		for _, prtb := range prtbs {
			t, status := status(prtb.Name, prtb.Migrated, prtb.Diff)
			i := item{title: t, status: status, objType: constants.PRTBType, obj: prtb}
			items = append(items, i)
		}
	case "projects":
		title = "Projects"
		for _, project := range constants.SC.ToMigrate.Projects {
			title, status := status(project.Name, project.Migrated, project.Diff)
			i := item{title: title, status: status, desc: project.Obj.Spec.Description, objType: constants.ProjectType, obj: project}
			items = append(items, i)
		}
	case "crtbs":
		title = "Cluster Users Permissions"
		for _, crtb := range constants.SC.ToMigrate.CRTBs {
			title, status := status(crtb.Name, crtb.Migrated, crtb.Diff)
			i := item{title: title, status: status, objType: constants.CRTBType, obj: crtb}
			items = append(items, i)
		}
	case "repos":
		title = "Cluster Catalog Repos"
		for _, repo := range constants.SC.ToMigrate.ClusterRepos {
			title, status := status(repo.Name, repo.Migrated, repo.Diff)
			i := item{title: title, status: status, objType: constants.RepoType, obj: repo}
			items = append(items, i)
		}
	}
	m.list = list.New(items, list.NewDefaultDelegate(), 8, 8)
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	m.list.Title = title
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Migrate,
			constants.Keymap.Back,
		}
	}
	return &m
}

// Update handle IO and commands
func (m Objects) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case tea.KeyMsg:

		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			return InitCluster()
		case key.Matches(msg, constants.Keymap.Enter):
			item := m.list.SelectedItem().(item)
			if item.objType == constants.ProjectType && item.status == constants.MigratedStatus {
				entry := InitObjects(item, constants.P)
				return entry.Update(constants.WindowSize)
			}
			if item.objType == constants.PRTBsType || item.objType == constants.NamespacesType {
				entry := InitObjects(item, constants.P)
				return entry.Update(constants.WindowSize)
			}
		case key.Matches(msg, constants.Keymap.Migrate):
			var err error
			item := m.list.SelectedItem().(item)
			if _, err = migrateObject(context.Background(), item); err != nil {
				panic(err)
			}
			item.objType = item.objType + "s"
			entry := InitObjects(item, constants.P)
			return entry.Update(constants.WindowSize)

		default:
			m.list, cmd = m.list.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m Objects) View() string {
	if m.quitting {
		return ""
	}
	return constants.DocStyle.Render(m.list.View() + "\n")
}

func status(name string, migrated, diff bool) (string, constants.MigrationStatus) {
	if migrated {
		if !diff {
			return name + " " + constants.CheckMark, constants.MigratedStatus
		} else {
			return name + " " + constants.WrongMark + " " + constants.WrongSpec, constants.WrongSpecStatus
		}
	}
	return name + " " + constants.WrongMark, constants.NotMigratedStatus
}

func migrateObject(ctx context.Context, i item) (tea.Msg, error) {
	var msg string
	switch i.objType {
	case constants.ProjectType:
		if i.status == constants.NotMigratedStatus {
			p := i.obj.(*cluster.Project)
			p.Mutate(constants.TC)
			if err := constants.Lclient.Projects.Create(ctx, constants.TC.Obj.Name, p.Obj, nil, v1.CreateOptions{}); err != nil {
				return nil, err
			}
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = p.Name
		}
	case constants.NamespaceType:
		if i.status == constants.NotMigratedStatus {
			ns := i.obj.(*cluster.Namespace)
			ns.Mutate(constants.TC.Obj.Name, ns.ProjectName)
			if _, err := constants.TC.Client.Namespace.Create(ns.Obj); err != nil {
				return nil, err
			}
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = ns.Name
		}

	case constants.PRTBType:
		if i.status == constants.NotMigratedStatus {
			prtb := i.obj.(*cluster.ProjectRoleTemplateBinding)
			prtb.Mutate(constants.TC.Obj.Name, prtb.ProjectName)
			if err := constants.Lclient.ProjectRoleTemplateBindings.Create(ctx, prtb.ProjectName, prtb.Obj, nil, v1.CreateOptions{}); err != nil {
				return nil, err
			}
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = prtb.Name
		}
	case constants.CRTBType:
		if i.status == constants.NotMigratedStatus {
			crtb := i.obj.(*cluster.ClusterRoleTemplateBinding)
			crtb.Mutate(constants.TC)
			if err := constants.Lclient.ClusterRoleTemplateBindings.Create(ctx, constants.TC.Obj.Name, crtb.Obj, nil, v1.CreateOptions{}); err != nil {
				return nil, err
			}
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = crtb.Name
		}

	case constants.RepoType:
		if i.status == constants.NotMigratedStatus {
			repo := i.obj.(*cluster.ClusterRepo)
			repo.Mutate()
			if err := constants.TC.Client.ClusterRepos.Create(ctx, constants.TC.Obj.Name, repo.Obj, nil, v1.CreateOptions{}); err != nil {
				return nil, err
			}
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = repo.Name
		}

	}
	return msg, nil
}

func updateClusters(ctx context.Context) error {
	if err := constants.SC.Populate(ctx, constants.Lclient); err != nil {
		return err
	}
	if err := constants.TC.Populate(ctx, constants.Lclient); err != nil {
		return err
	}
	if err := constants.SC.Compare(ctx, constants.Lclient, constants.TC); err != nil {
		return err
	}
	return nil
}
