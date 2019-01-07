package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/camptocamp/prometheus-orchestrators-sd/prometheus"
)

// Start is the main function that handles requests from POSD agents
func Start(cmd *cobra.Command, args []string) {
	bindAddress, _ := cmd.Flags().GetString("bind-address")
	psk, _ := cmd.Flags().GetString("psk")
	outputFile, _ := cmd.Flags().GetString("output-file")
	inputFile, _ := cmd.Flags().GetString("input-file")
	customSCFile, _ := cmd.Flags().GetString("custom-sc-file")
	refreshInterval, _ := cmd.Flags().GetString("refresh-interval")
	agentsFile, _ := cmd.Flags().GetString("agents-file")

	interval, err := time.ParseDuration(refreshInterval)
	if err != nil {
		log.Fatalf("failed to parse refresh interval: %s", err)
	}

	os.MkdirAll(filepath.Dir(outputFile), 0755)

	var jobs []prometheus.ScrapeConfig
	var pc prometheus.Config

	yamlFile, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("failed to load file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &pc)
	if err != nil {
		log.Fatalf("failed to unmarshal prometheus config: %s", err)
	}

	go httpServer(bindAddress)

	for {
		log.Debugf("Sleeping for %s", interval)
		time.Sleep(interval)

		agents, err := retrieveAgentsList(agentsFile)
		if err != nil {
			log.Errorf("failed to retrieve agents list: %s", err)
			continue
		}

		for agentName, agentEndpoint := range agents {
			jobs, err = retrieveJobsFromAgent(agentEndpoint, psk)
			if err != nil {
				log.WithFields(log.Fields{
					"agent": agentName,
				}).Errorf("failed to retrieve jobs from agent: %s", err)
				continue
			}
			err = updatePromConfig(&pc, jobs, outputFile, customSCFile)
			if err != nil {
				log.WithFields(log.Fields{
					"agent": agentName,
				}).Errorf("failed to update prometheus config: %s", err)
				continue
			}
		}
	}
	return
}

func updatePromConfig(pc *prometheus.Config, jobs []prometheus.ScrapeConfig, outputFile, customSCFile string) (err error) {
	for _, job := range jobs {
		err = formatScrapeConfig(&job, customSCFile)
		if err != nil {
			err = fmt.Errorf("failed to format scrape config: %s", err)
			return
		}

		exists := false
		for sck, sc := range pc.ScrapeConfigs {
			if sc.JobName == job.JobName {
				exists = true
				if reflect.DeepEqual(sc, job) {
					return
				}
				pc.ScrapeConfigs[sck] = job
				log.WithFields(log.Fields{
					"job_name": job.JobName,
				}).Infof("prometheus scrape config updated")
			}
		}
		if !exists {
			pc.ScrapeConfigs = append(pc.ScrapeConfigs, job)
			log.WithFields(log.Fields{
				"job_name": job.JobName,
			}).Infof("prometheus scrape config added")
		}

		if err != nil {
			log.WithFields(log.Fields{
				"job_name": job.JobName,
			}).Errorf("failed to export targets: %s", err)
		}
		d, err := yaml.Marshal(pc)
		if err != nil {
			log.Errorf("Failed to encode prometheus configuration: %s", d)
		}

		err = ioutil.WriteFile(outputFile, d, 0644)
		if err != nil {
			log.Errorf("failed to write output file: %s", err)
		}
	}
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

func retrieveAgentsList(listPath string) (agents map[string]string, err error) {
	raw, err := ioutil.ReadFile(listPath)
	if err != nil {
		err = fmt.Errorf("failed to read agents list file: %s", err)
		return
	}

	err = yaml.Unmarshal(raw, &agents)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal agent list: %s", err)
		return
	}
	return
}

func retrieveJobsFromAgent(endpoint, psk string) (jobs []prometheus.ScrapeConfig, err error) {
	clientHTTP := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		err = fmt.Errorf("failed to build request: %s", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", psk))

	res, err := clientHTTP.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to request endpoint: %s", err)
		return
	}

	defer res.Body.Close()
	rawJobs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("failed to read the response body: %s", err)
		return
	}

	err = yaml.Unmarshal([]byte(rawJobs), &jobs)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal jobs: %s", err)
		return
	}
	return
}

func httpServer(bindAddress string) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	log.Infof("Prometheus endpoint on http://%s/metrics", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, router))
}
