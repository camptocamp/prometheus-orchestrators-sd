package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/camptocamp/prometheus-orchestrators-sd/orchestrators"
)

var orchestratorArg string
var refreshInterval string

var clientCmd = &cobra.Command{
	Use:   "client [posd server]",
	Short: "Start POSD as client",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		o, err := orchestrators.GetOrchestrator(orchestratorArg, args[0], refreshInterval)
		if err != nil {
			log.Fatalf("failed to retrieve orchestrator: %s", err)
		}
		err = o.Start()
		if err != nil {
			log.Fatalf("%s", err)
		}
	},
}

func init() {
	clientCmd.Flags().StringVarP(&orchestratorArg, "orchestrator", "o", viper.GetString("POSD_ORCHESTRATOR"), "Orchestrator from where the client is running.")
	clientCmd.Flags().StringVarP(&refreshInterval, "interval", "i", viper.GetString("POSD_INTERVAL"), "Seconds between two targets discovery.")
	rootCmd.AddCommand(clientCmd)
}
