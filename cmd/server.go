package cmd

import (
	"github.com/spf13/cobra"

	"github.com/camptocamp/prometheus-orchestrators-sd/server"
)

var bindAddress string
var outputFile string
var inputFile string
var customSCFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start POSD as server",
	PreRun: func(cmd *cobra.Command, args []string) {
		if bindAddress == "" {
			bindAddress = "0.0.0.0:8000"
		}

		if outputFile == "" {
			outputFile = "prometheus.yml"
		}

		if inputFile == "" {
			inputFile = "prometheus.yml"
		}
	},
	Run: server.Start,
}

func init() {
	serverCmd.Flags().StringVarP(&bindAddress, "bind-address", "b", "", "Address to bind on.")
	serverCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Output file (default: prometheus.yml)")
	serverCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "Input file (default: prometheus.yml)")
	serverCmd.Flags().StringVarP(&customSCFile, "custom-sc-file", "c", "", "Custom scrape config fields file")
	rootCmd.AddCommand(serverCmd)
}
