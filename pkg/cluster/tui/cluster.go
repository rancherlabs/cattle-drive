package tui

import (
	"galal-hussein/cattle-drive/pkg/cluster/tui/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	nav mode = iota
	migrate
	migrated
)

type Model struct {
	mode     mode
	list     list.Model
	quitting bool
}

type item struct {
	title   string
	desc    string
	objType string
	obj     interface{}
	status  constants.MigrationStatus
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

func InitCluster() (tea.Model, tea.Cmd) {
	items := newClusterList()
	m := Model{mode: nav, list: list.New(items, list.NewDefaultDelegate(), 8, 8)}
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	m.list.Title = "Cluster " + constants.SC.Obj.Spec.DisplayName + " migration"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Migrate,
			constants.Keymap.Back,
		}
	}
	return m, func() tea.Msg { return errMsg{nil} }
}

func newClusterList() []list.Item {
	items := []list.Item{
		item{title: "Projects", desc: "Projects, Namespaces, and PRTB", objType: constants.ProjectsType, obj: nil},
		item{title: "Cluster User Permissions", desc: "user permissions for the cluster (CRTB)", objType: constants.CRTBsType, obj: nil},
		item{title: "Catalog Repos", desc: "Cluster apps repos", objType: constants.ReposType, obj: nil},
	}
	return items
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case key.Matches(msg, constants.Keymap.Enter):
			entry := InitObjects(m.list.SelectedItem().(item))
			return entry.Update(constants.WindowSize)

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
	return constants.DocStyle.Render(m.list.View() + "\n")
}
