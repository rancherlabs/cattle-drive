package tui

import (
	"fmt"
	"galal-hussein/cattle-drive/pkg/cluster"
	"galal-hussein/cattle-drive/pkg/cluster/tui/constants"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func StartTea(sc, tc *cluster.Cluster) error {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	constants.SC = sc
	constants.TC = tc

	m, _ := InitCluster() // TODO: can we acknowledge this error
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		return err
	}
	return nil
}
