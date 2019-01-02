package cmd

import (
	"github.com/spf13/cobra"

	"github.com/camptocamp/prometheus-orchestrators-sd/server"
)

var bindAddress string
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start POSD as server",
	PreRun: func(cmd *cobra.Command, args []string) {
		if bindAddress == "" {
			bindAddress = "0.0.0.0:8000"
		}
	},
	Run: server.Start,
}

func init() {
	serverCmd.Flags().StringVarP(&bindAddress, "bind-address", "b", "", "Address to bind on.")
	rootCmd.AddCommand(serverCmd)
}
