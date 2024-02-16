package tui

import (
	"galal-hussein/cattle-drive/pkg/cluster/tui/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/apimachinery/pkg/runtime"
)

type mode int

const (
	nav mode = iota
	migrate
)

type Model struct {
	mode     mode
	list     list.Model
	input    textinput.Model
	quitting bool
}

type item struct {
	title   string
	desc    string
	objType string
	obj     runtime.Object
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

func InitCluster() (tea.Model, tea.Cmd) {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	items := newClusterList()
	m := Model{mode: nav, list: list.New(items, list.NewDefaultDelegate(), 8, 8), input: input}
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
		item{title: "Projects", desc: "Projects objects", objType: "project", obj: nil},
		item{title: "Cluster user permissions", desc: "user permissions for the cluster (CRTB)", objType: "crtb", obj: nil},
		item{title: "Catalog repos", desc: "Cluster apps repos", objType: "repo", obj: nil},
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
		if m.input.Focused() {
			if key.Matches(msg, constants.Keymap.Enter) {
				if m.mode == migrate {
					// cmds = append(cmds, createProjectCmd(m.input.Value(), constants.Pr))
					// todo migrate object
				}
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
			}
			if key.Matches(msg, constants.Keymap.Back) {
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
			}
			// only log keypresses for the input field when it's focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, constants.Keymap.Migrate):
				m.mode = migrate
				m.input.Focus()
				cmd = textinput.Blink
			case key.Matches(msg, constants.Keymap.Quit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, constants.Keymap.Enter):
				// activeProject := m.list.SelectedItem().(project.Project)
				// entry := InitEntry(constants.Er, activeProject.ID, constants.P)
				// return entry.Update(constants.WindowSize)

			default:
				m.list, cmd = m.list.Update(msg)
			}
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.input.Focused() {
		return constants.DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return constants.DocStyle.Render(m.list.View() + "\n")
}

// func entryModel(source, target *Cluster) model {
// items := []list.Item{
// 	item{title: "Projects", desc: "Projects objects", objType: "project", obj: nil},
// 	item{title: "Cluster user permissions", desc: "user permissions for the cluster (CRTB)", objType: "crtb", obj: nil},
// 	item{title: "Catalog repos", desc: "Cluster apps repos", objType: "repo", obj: nil},
// }
// 	m := model{
// 		source: source,
// 		target: target,
// 		list:   list.New(items, list.NewDefaultDelegate(), 0, 0),
// 	}
// 	m.list.Title = "Cluster [" + source.Obj.Name + "] migration"
// 	return m
// }

// func (m MainModel) Init() tea.Cmd {
// 	return nil
// }

// func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	var cmds []tea.Cmd
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c":
// 			return m, tea.Quit
// 		case "enter":
// 			m.state = objectView
// 		}
// 	case tea.WindowSizeMsg:
// 		h, v := docStyle.GetFrameSize()
// 		m.list.SetSize(msg.Width-h, msg.Height-v)
// 	}

// 	switch m.state {
// 	case entryView:

// 	}
// }

// func (m MainModel) View() string {
// 	switch m.state {
// 	case entryView:
// 		return m.entry.View()
// 	default:
// 		return m.objects.View()
// 	}
// }

// func (m model) buildModel(i list.Item, msg tea.Msg) tea.Model {
// 	selectedItem, _ := i.(item)

// 	// list over each item and check its type
// 	switch selectedItem.objType {
// 	case "project":
// 		items := []list.Item{}
// 		for _, project := range m.source.ToMigrate.Projects {
// 			items = append(items, item{title: project.Name, desc: "project description", objType: "projectObj", obj: project.Obj})
// 		}
// 		newModel := model{
// 			source: m.source,
// 			target: m.target,
// 			list:   list.New(items, list.NewDefaultDelegate(), 0, 0),
// 		}
// 		newModel.list.Title = "Projects"

// 		return newModel
// 	}
// 	return m
// }
