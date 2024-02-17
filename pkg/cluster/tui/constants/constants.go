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
	CheckMark                      = "\u2714"
	WrongMark                      = "\u2718"
	WrongSpec                      = "(Wrong fields)"
	MigratedStatus MigrationStatus = iota
	WrongSpecStatus
	NotMigratedStatus
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

// DocStyle styling for viewports
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

// ErrStyle provides styling for error messages
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render

// AlertStyle provides styling for alert messages
var AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render

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
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
