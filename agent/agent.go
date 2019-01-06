package agent

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/camptocamp/prometheus-orchestrators-sd/orchestrators"
	"github.com/camptocamp/prometheus-orchestrators-sd/prometheus"
)

func Start(cmd *cobra.Command, args []string) {
	bindAddress, _ := cmd.Flags().GetString("bind-address")
	orchArg, _ := cmd.Flags().GetString("orchestrator")
	psk, _ := cmd.Flags().GetString("psk")
	refreshInterval, _ := cmd.Flags().GetString("refresh-interval")

	interval, err := time.ParseDuration(refreshInterval)
	if err != nil {
		log.Fatalf("failed to parse refresh interval: %s", err)
	}

	o, err := orchestrators.GetOrchestrator(orchArg)
	if err != nil {
		log.Fatalf("failed to retrieve orchestrator: %s", err)
	}

	var jobs []prometheus.ScrapeConfig
	go apiServer(bindAddress, psk, &jobs)

	for {
		jobs, err = o.DiscoverTargets()
		if err != nil {
			err = fmt.Errorf("failed to discover targets: %s", err)
			return
		}
		log.Debugf("Sleeping for %s", interval)
		time.Sleep(interval)
	}
}

func apiServer(bindAddress string, psk string, jobs *[]prometheus.ScrapeConfig) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	router.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", psk) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("403 - Unauthorized"))
			return
		}

		d, err := yaml.Marshal(jobs)
		if err != nil {
			log.Errorf("failed to marshal jobs: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
			return
		}
		w.Write(d)
	}).Methods("GET")

	log.Infof("Listening on %s", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, router))
}
