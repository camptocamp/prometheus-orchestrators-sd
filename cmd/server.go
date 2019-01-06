package cmd

import (
	"github.com/spf13/cobra"

	"github.com/camptocamp/prometheus-orchestrators-sd/server"
)

var (
	outputFile   string
	inputFile    string
	customSCFile string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start POSD as server",
	PreRun: func(cmd *cobra.Command, args []string) {
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
	serverCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Output file (default: prometheus.yml)")
	serverCmd.Flags().StringVarP(&inputFile, "input-file", "k", "", "Input file (default: prometheus.yml)")
	serverCmd.Flags().StringVarP(&customSCFile, "custom-sc-file", "c", "", "Custom scrape config fields file")
	rootCmd.AddCommand(serverCmd)
}
