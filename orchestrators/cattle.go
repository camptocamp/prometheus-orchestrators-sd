package orchestrators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/v2"

	"github.com/camptocamp/prometheus-orchestrators-sd/prometheus"
)

// CattleOrchestrator implements a container orchestrator for Cattle
type CattleOrchestrator struct {
	refreshInterval time.Duration
	posdServer      string
	client          *client.RancherClient
}

// NewCattleOrchestrator creates a Cattle client
func NewCattleOrchestrator(posdServer string, refreshInterval string) (o *CattleOrchestrator) {
	var err error

	r, err := time.ParseDuration(refreshInterval)
	if err != nil {
		log.Fatalf("failed to parse refresh interval: %s", err)
	}

	o = &CattleOrchestrator{
		posdServer:      posdServer,
		refreshInterval: r,
	}

	o.client, err = client.NewRancherClient(&client.ClientOpts{
		Url:       os.Getenv("CATTLE_URL"),
		AccessKey: os.Getenv("CATTLE_ACCESS_KEY"),
		SecretKey: os.Getenv("CATTLE_SECRET_KEY"),
		Timeout:   30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create a new Rancher client: %s", err)
	}

	return
}

// GetName returns orchestrator's name
func (o *CattleOrchestrator) GetName() string {
	return "Cattle"
}

// Start discovers targets from the orchestrator
func (o *CattleOrchestrator) Start() (err error) {
	for {
		err = o.discoverTargets()
		if err != nil {
			err = fmt.Errorf("failed to discover targets: %s", err)
			return
		}
		log.Debugf("Sleeping for %s", o.refreshInterval)
		time.Sleep(o.refreshInterval)
	}
	return
}

// sendTarget shares a scrape config with a POSD server
func (o *CattleOrchestrator) sendTarget(p prometheus.ScrapeConfig) (err error) {
	data, err := json.Marshal(p)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal prometheus endpoint: %s", err)
		return
	}

	clientHTTP := &http.Client{}
	req, err := http.NewRequest("POST", o.posdServer, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to build request: %s", err)
	}
	res, err := clientHTTP.Do(req)
	if err != nil {
		return fmt.Errorf("connection to POSD server failed: %s", err)
	}
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	return
}

// buildTarget extracts informations from stacks config to build a scrape config
func (o *CattleOrchestrator) buildTarget(stack client.Stack) (err error) {
	p := prometheus.ScrapeConfig{
		HonorLabels: true,
		MetricsPath: "/federate",
		BasicAuth:   make(map[string]string),
	}

	project, err := o.client.Project.ById(stack.AccountId)
	if err != nil {
		return fmt.Errorf("failed to retrieve project `%s`: %s", stack.AccountId, err)
	}

	p.JobName = fmt.Sprintf("cattle_%s_%s_%s", project.Name, project.Id, stack.Id)

	var promPort string
	if stack.Environment["PROMETHEUS_PORT"] != nil {
		promPort = stack.Environment["PROMETHEUS_PORT"].(string)
	} else {
		promPort = "9443"
	}

	p.StaticConfigs = []prometheus.StaticConfig{
		prometheus.StaticConfig{
			Targets: []string{
				fmt.Sprintf("%s:%s", stack.Environment["PROMETHEUS_FQDN"].(string), promPort),
			},
		},
	}

	if stack.Environment["PROMETHEUS_USERNAME"] != nil {
		p.BasicAuth["username"] = stack.Environment["PROMETHEUS_USERNAME"].(string)
	}

	if stack.Environment["PROMETHEUS_PASSWORD"] != nil {
		p.BasicAuth["password"] = stack.Environment["PROMETHEUS_PASSWORD"].(string)
	}

	if stack.Environment["PROMETHEUS_SCHEME"] != nil {
		p.Scheme = stack.Environment["PROMETHEUS_SCHEME"].(string)
	} else {
		p.Scheme = "https"
	}

	err = o.sendTarget(p)
	if err != nil {
		return fmt.Errorf("failed to export target: %s", err)
	}
	return
}

// discoverTargets is parse all Cattle stacks to detect those which has
// the environment variable "PROMETHEUS_FQDN" defined
func (o *CattleOrchestrator) discoverTargets() (err error) {
	stacks, err := o.client.Stack.List(&client.ListOpts{
		Filters: map[string]interface{}{
			"limit": -2,
			"all":   true,
		},
	})
	if err != nil {
		err = fmt.Errorf("failed to list stacks: %s", err)
		return
	}

	for _, stack := range stacks.Data {
		if stack.Environment["PROMETHEUS_FQDN"] != nil {
			err = o.buildTarget(stack)
			if err != nil {
				log.Errorf("failed to build target: %s", err)
			}
		}
	}
	return
}

func detectCattle() bool {
	_, err := net.LookupHost("rancher-metadata")
	if err != nil {
		return false
	}
	return true
}
