package tui

import (
	"fmt"
	"log"
	"os"
	"rancherlabs/cattle-drive/pkg/client"
	"rancherlabs/cattle-drive/pkg/cluster"
	"rancherlabs/cattle-drive/pkg/cluster/tui/constants"

	tea "github.com/charmbracelet/bubbletea"
)

func StartTea(sc, tc *cluster.Cluster, client, tClient *client.Clients, logFilePath string) error {
	if f, err := tea.LogToFile(logFilePath, "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		constants.LogFile = *f
		defer func() {
			err = constants.LogFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	constants.SC = sc
	constants.TC = tc
	constants.Lclient = client
	constants.TClient = tClient

	m, _ := InitCluster(nil)
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		return err
	}
	return nil
}
