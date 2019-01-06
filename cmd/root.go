package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose         bool
	psk             string
	refreshInterval string
	bindAddress     string
)

var rootCmd = &cobra.Command{
	Use:   "posd",
	Short: "Prometheus Orchestrators Service Discovery",
}

// Execute is the main thread, required by Cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&psk, "psk", "p", os.Getenv("POSD_PSK"), "Pre-shared key which allows communication between agent and server.")
	rootCmd.PersistentFlags().StringVarP(&refreshInterval, "refresh-interval", "i", os.Getenv("POSD_INTERVAL"), "Seconds between two targets discovery.")
	rootCmd.PersistentFlags().StringVarP(&bindAddress, "bind-address", "b", "", "Address to bind on.")
}

func initConfig() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if psk == "" {
		log.Fatalf("Flag `psk` is required.")
	}

	if refreshInterval == "" {
		refreshInterval = "1m"
	}

	if bindAddress == "" {
		bindAddress = "0.0.0.0:8000"
	}
}
