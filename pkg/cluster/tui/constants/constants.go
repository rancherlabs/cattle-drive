package constants

import (
	"galal-hussein/cattle-drive/pkg/client"
	"galal-hussein/cattle-drive/pkg/cluster"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* CONSTANTS */
const (
	/* Migration Types */
	MigratedStatus MigrationStatus = iota
	WrongSpecStatus
	NotMigratedStatus
	CheckMark = "\u2714"
	WrongMark = "\u2718"
	WrongSpec = "(Wrong fields)"

	/* Object Types */
	ProjectsType   = "projects"
	CRTBsType      = "crtbs"
	NamespacesType = "namespaces"
	PRTBsType      = "prtbs"
	ReposType      = "repos"
	ProjectType    = "project"
	CRTBType       = "crtb"
	NamespaceType  = "namespace"
	PRTBType       = "prtb"
	RepoType       = "repo"
)

type MigrationStatus int

var (
	// P the current tea program
	P *tea.Program
	// SC the source cluster
	SC *cluster.Cluster
	// TC the target cluster
	TC *cluster.Cluster
	// Local Client
	Lclient *client.Clients
	// WindowSize store the size of the terminal window
	WindowSize tea.WindowSizeMsg
)

/* STYLING */
var DocStyle = lipgloss.NewStyle().Margin(0, 2).Foreground(lipgloss.Color("241"))

type keymap struct {
	Enter   key.Binding
	Migrate key.Binding
	Back    key.Binding
	Quit    key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Migrate: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "migrate"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "main menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
