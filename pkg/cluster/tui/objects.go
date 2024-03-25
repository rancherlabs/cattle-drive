package tui

import (
	"context"
	"fmt"
	"galal-hussein/cattle-drive/pkg/cluster"
	"galal-hussein/cattle-drive/pkg/cluster/tui/constants"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	errMsg  struct{ error }
	tickMsg time.Time
)

type Objects struct {
	mode         mode
	list         list.Model
	quitting     bool
	progress     progress.Model
	activeObject string
}

// Init run any intial IO on program start
func (m Objects) Init() tea.Cmd {
	return tickCmd()
}

func InitObjects(i item) *Objects {
	prog := progress.New(progress.WithSolidFill("#04B575"))
	var title string
	items := []list.Item{}
	m := Objects{mode: nav, progress: prog}
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
			i := item{title: t, desc: prtb.Description, status: status, objType: constants.PRTBType, obj: prtb}
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
			i := item{title: title, desc: crtb.Description, status: status, objType: constants.CRTBType, obj: crtb}
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
	delegateObjKeys := *delegateKeys
	delegateObjKeys.MigrateAll = key.NewBinding()
	delegate := newItemDelegate(&delegateObjKeys)
	objList := list.New(items, delegate, 8, 8)
	objList.Styles.Title = constants.TitleStyle
	m.list = objList
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	m.list.Title = title
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
	case tickMsg:
		m.mode = migrate
		cmd := m.progress.IncrPercent(0.25)
		if m.progress.Percent() == 1.0 || m.mode == migrated {
			return InitCluster(nil)
		}
		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, delegateKeys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, delegateKeys.Back):
			return InitCluster(nil)
		case key.Matches(msg, delegateKeys.Enter):
			if m.list.SelectedItem() == nil {
				return m, tea.Batch(cmds...)
			}
			item := m.list.SelectedItem().(item)
			if item.objType == constants.ProjectType && item.status == constants.MigratedStatus {
				entry := InitObjects(item)
				return entry.Update(constants.WindowSize)
			}
			if item.objType == constants.PRTBsType || item.objType == constants.NamespacesType {
				entry := InitObjects(item)
				return entry.Update(constants.WindowSize)
			}
		case key.Matches(msg, delegateKeys.Migrate):
			item := m.list.SelectedItem().(item)
			if item.status == constants.NotMigratedStatus {
				m.mode = migrate
				m.activeObject = item.objType + "/" + item.title
				go m.migrateObject(context.Background(), item)
				return m, tickCmd()
			}
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
	if m.mode == migrate {
		pad := strings.Repeat(" ", 2)
		return "\n\n Waiting for object [" + m.activeObject + "] to be migrated\n\n" + pad + m.progress.View() + "\n\n" + pad
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

func (m *Objects) migrateObject(ctx context.Context, i item) (tea.Msg, error) {
	var (
		objectMigrated bool
		msg            string
	)
	cl := constants.TClient
	if cl == nil {
		cl = constants.Lclient
	}

	switch i.objType {
	case constants.ProjectType:
		if i.status == constants.NotMigratedStatus {
			p := i.obj.(*cluster.Project)
			p.Mutate(constants.TC)
			if err := cl.Projects.Create(ctx, constants.TC.Obj.Name, p.Obj, nil, v1.CreateOptions{}); err != nil {
				return nil, err
			}
			objectMigrated = true
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
			objectMigrated = true
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = ns.Name
		}

	case constants.PRTBType:
		if i.status == constants.NotMigratedStatus {
			prtb := i.obj.(*cluster.ProjectRoleTemplateBinding)
			prtb.Mutate(constants.TC.Obj.Name, prtb.ProjectName)
			if err := cl.ProjectRoleTemplateBindings.Create(ctx, prtb.ProjectName, prtb.Obj, nil, v1.CreateOptions{}); err != nil {
				return nil, err
			}
			objectMigrated = true
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = prtb.Name
		}
	case constants.CRTBType:
		if i.status == constants.NotMigratedStatus {
			crtb := i.obj.(*cluster.ClusterRoleTemplateBinding)
			crtb.Mutate(constants.TC)
			if err := cl.ClusterRoleTemplateBindings.Create(ctx, constants.TC.Obj.Name, crtb.Obj, nil, v1.CreateOptions{}); err != nil {
				return nil, err
			}
			objectMigrated = true
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
			objectMigrated = true
			if err := updateClusters(ctx); err != nil {
				return nil, err
			}
			msg = repo.Name
		}

	}
	m.mode = migrated
	if objectMigrated {
		fmt.Fprintf(&constants.LogFile, "[%s] migrated object [%s/%s]\n", time.Now().String(), i.objType, i.title)
	}
	return msg, nil
}

func updateClusters(ctx context.Context) error {
	if err := constants.SC.Populate(ctx, constants.Lclient); err != nil {
		return err
	}
	if constants.TClient != nil {
		if err := constants.TC.Populate(ctx, constants.TClient); err != nil {
			return err
		}
	} else {
		if err := constants.TC.Populate(ctx, constants.Lclient); err != nil {
			return err
		}
	}
	if err := constants.SC.Compare(ctx, constants.TC); err != nil {
		return err
	}
	fmt.Fprintf(&constants.LogFile, "[%s] successfully updated cluster [%s]\n", time.Now().String(), constants.SC.Obj.Spec.DisplayName)
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
