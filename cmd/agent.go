package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/camptocamp/prometheus-orchestrators-sd/agent"
)

var (
	orchestratorArg string
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start POSD agent",
	Run:   agent.Start,
}

func init() {
	agentCmd.Flags().StringVarP(&orchestratorArg, "orchestrator", "o", viper.GetString("POSD_ORCHESTRATOR"), "Orchestrator from where the client is running.")
	rootCmd.AddCommand(agentCmd)
}
