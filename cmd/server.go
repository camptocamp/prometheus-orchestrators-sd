package cmd

import (
	"github.com/spf13/cobra"

	"github.com/camptocamp/prometheus-orchestrators-sd/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start POSD as server",
	Run:   server.Start,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
