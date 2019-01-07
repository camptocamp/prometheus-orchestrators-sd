package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/camptocamp/prometheus-orchestrators-sd/server"
)

var (
	outputFile   string
	inputFile    string
	customSCFile string
	agentsFile   string
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

		if agentsFile == "" {
			agentsFile = "agents.yml"
		}
	},
	Run: server.Start,
}

func init() {
	serverCmd.Flags().StringVarP(&outputFile, "output-file", "o", os.Getenv("POSD_OUTPUT_FILE"), "Output file (default: prometheus.yml)")
	serverCmd.Flags().StringVarP(&inputFile, "input-file", "i", os.Getenv("POSD_INPUT_FILE"), "Input file (default: prometheus.yml)")
	serverCmd.Flags().StringVarP(&customSCFile, "custom-sc-file", "c", os.Getenv("POSD_CUSTOM_SC_FILE"), "Custom scrape config fields file")
	serverCmd.Flags().StringVarP(&agentsFile, "agents-file", "a", os.Getenv("POSD_AGENTS_FILE"), "Agents list file path (default: agents.yml)")
	rootCmd.AddCommand(serverCmd)
}
