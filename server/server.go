package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/camptocamp/prometheus-orchestrators-sd/prometheus"
)

type formattedError struct {
	Msg string `json:"msg"`
}

type prometheusConfig struct {
	GlobalConfig   interface{}               `json:"global"`
	AlertingConfig interface{}               `json:"alerting,omitempty"`
	RuleFiles      interface{}               `json:"rule_files,omitempty"`
	ScrapeConfigs  []prometheus.ScrapeConfig `json:"scrape_configs"`

	RemoteWriteConfigs interface{} `json:"remote_write,omitempty"`
	RemoteReadConfigs  interface{} `json:"remote_read,omitempty"`
}

// Start is the main function that handles requests from POSD agents
func Start(cmd *cobra.Command, args []string) {
	var pc prometheusConfig
	bindAddress, _ := cmd.Flags().GetString("bind-address")
	outputFile, _ := cmd.Flags().GetString("output-file")
	inputFile, _ := cmd.Flags().GetString("input-file")
	customSCFile, _ := cmd.Flags().GetString("custom-sc-file")

	yamlFile, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("failed to load file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &pc)
	if err != nil {
		log.Fatalf("failed to unmarshal prometheus config: %s", err)
	}

	// Write basic output file to allow prometheus server to start
	ioutil.WriteFile(outputFile, yamlFile, 0644)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/endpoint", func(w http.ResponseWriter, r *http.Request) {
		var promEndpoint prometheus.ScrapeConfig
		err := json.NewDecoder(r.Body).Decode(&promEndpoint)
		if err != nil {
			json.NewEncoder(w).Encode(formattedError{
				Msg: fmt.Sprintf("failed to decode body: %s", err),
			})
			return
		}
		updateConfig(&pc, promEndpoint, outputFile, customSCFile)
		json.NewEncoder(w).Encode(promEndpoint)
	}).Methods("POST")
	log.Infof("Listening on %s", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, router))
}

func updateConfig(pc *prometheusConfig, pe prometheus.ScrapeConfig, outputFile, customSCFile string) (err error) {
	err = formatScrapeConfig(&pe, customSCFile)
	if err != nil {
		err = fmt.Errorf("failed to format scrape config: %s", err)
		return
	}

	exists := false
	for sck, sc := range pc.ScrapeConfigs {
		if sc.JobName == pe.JobName {
			exists = true
			if reflect.DeepEqual(sc, pe) {
				return
			}
			pc.ScrapeConfigs[sck] = pe
			log.WithFields(log.Fields{
				"name": pe.JobName,
			}).Infof("config updated")
		}
	}
	if !exists {
		pc.ScrapeConfigs = append(pc.ScrapeConfigs, pe)
		log.WithFields(log.Fields{
			"name": pe.JobName,
		}).Infof("config added")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"name": pe.JobName,
		}).Errorf("failed to export targets: %s", err)
	}
	d, err := yaml.Marshal(pc)
	if err != nil {
		log.Errorf("Failed to encode prometheus configuration: %s", d)
	}
	ioutil.WriteFile(outputFile, d, 0644)

	return
}

func formatScrapeConfig(pe *prometheus.ScrapeConfig, customSCFile string) (err error) {
	if customSCFile == "" {
		return
	}

	var customSC prometheus.ScrapeConfig
	yamlFile, err := ioutil.ReadFile(customSCFile)
	if err != nil {
		return fmt.Errorf("failed to load file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &customSC)
	if err != nil {
		return fmt.Errorf("failed to unmarshal prometheus config: %s", err)
	}

	if customSC.Params != nil && pe.Params == nil {
		pe.Params = customSC.Params
	}
	return
}
