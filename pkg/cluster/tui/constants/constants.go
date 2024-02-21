package constants

import (
	"galal-hussein/cattle-drive/pkg/client"
	"galal-hussein/cattle-drive/pkg/cluster"

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
	Migratedch = make(chan bool)
)

/* STYLING */
var (
	DocStyle = lipgloss.NewStyle().Margin(0, 2).Foreground(lipgloss.Color("241"))

	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)
