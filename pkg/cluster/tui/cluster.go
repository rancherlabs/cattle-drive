package tui

import (
	"context"
	"fmt"
	"rancherlabs/cattle-drive/pkg/cluster/tui/constants"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	nav mode = iota
	migrate
	migrated
)

type Model struct {
	mode         mode
	migratingAll bool
	list         list.Model
	progress     progress.Model
	quitting     bool
}

type item struct {
	title   string
	desc    string
	objType string
	obj     interface{}
	status  constants.MigrationStatus
}

var (
	delegateKeys = newDelegateKeyMap()
)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

func InitCluster(msg tea.Msg) (tea.Model, tea.Cmd) {
	prog := progress.New(progress.WithSolidFill("#04B575"))
	items := newClusterList()
	delegate := newItemDelegate(delegateKeys)
	clusterList := list.New(items, delegate, 8, 8)
	clusterList.Styles.Title = constants.TitleStyle

	m := Model{mode: nav, list: clusterList, progress: prog}
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	m.list.Title = "Cluster " + constants.SC.Obj.Spec.DisplayName + " migration"
	if msg != nil {
		return m, func() tea.Msg { return msg }
	}
	return m, func() tea.Msg { return errMsg{nil} }
}

func newClusterList() []list.Item {
	items := []list.Item{
		item{title: "Projects", desc: "Projects, Namespaces, and PRTB", objType: constants.ProjectsType, obj: nil},
		item{title: "Cluster User Permissions", desc: "user permissions for the cluster (CRTB)", objType: constants.CRTBsType, obj: nil},
		item{title: "Catalog Repos", desc: "Cluster apps repos", objType: constants.ReposType, obj: nil},
	}
	if constants.TC.ExternalRancher || constants.SC.ExternalRancher {
		items = append(items, item{title: "Users", desc: "Rancher Users", objType: constants.UsersType, obj: nil})
	}
	return items
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tickMsg:
		m.migratingAll = true
		for {
			select {
			case <-time.After(time.Millisecond * 500):
				if m.progress.Percent() == 1.0 {
					cmd := m.progress.SetPercent(0)
					return m, tea.Batch(tickCmd(), cmd)
				}
				cmd := m.progress.IncrPercent(0.25)
				return m, tea.Batch(tickCmd(), cmd)
			case <-constants.Migratedch:
				return InitCluster(nil)
			}
		}
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, delegateKeys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, delegateKeys.Enter):
			entry := InitObjects(m.list.SelectedItem().(item))
			return entry.Update(constants.WindowSize)
		case key.Matches(msg, delegateKeys.MigrateAll):
			m.mode = migrate
			go m.migrateCluster(context.Background())
			return InitCluster(tickMsg{})
		default:
			m.list, cmd = m.list.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.migratingAll {
		pad := strings.Repeat(" ", 2)
		return "\n\n Migrating all objects.. please wait" + pad + m.progress.View() + "\n\n" + pad
	}
	return constants.DocStyle.Render(m.list.View() + "\n")
}

func (m *Model) migrateCluster(ctx context.Context) {
	fmt.Fprintf(&constants.LogFile, "[%s] initiating cluster objects migrate:\n", time.Now().String())
	cl := constants.TClient
	if cl == nil {
		cl = constants.Lclient
	}
	if err := constants.SC.Migrate(ctx, cl, constants.TC, &constants.LogFile); err != nil {
		fmt.Fprintf(&constants.LogFile, "[%s] [error] %v\n", time.Now().String(), err)
		m.Update(tea.Quit())
	}
	if err := updateClusters(ctx); err != nil {
		fmt.Fprintf(&constants.LogFile, "[%s] [error] %v\n", time.Now().String(), err)
		m.Update(tea.Quit())
	}
	constants.Migratedch <- true
}
