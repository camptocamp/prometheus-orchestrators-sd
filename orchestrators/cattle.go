package orchestrators

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/prometheus-orchestrators-sd/prometheus"
	"github.com/rancher/go-rancher/v2"
)

// CattleOrchestrator implements a container orchestrator for Cattle
type CattleOrchestrator struct {
	client *client.RancherClient
}

// NewCattleOrchestrator creates a Cattle client
func NewCattleOrchestrator() (o *CattleOrchestrator) {
	var err error

	o = &CattleOrchestrator{}

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

// buildTarget extracts informations from stacks config to build a scrape config
func (o *CattleOrchestrator) buildTarget(stack client.Stack) (p prometheus.ScrapeConfig, err error) {
	p = prometheus.ScrapeConfig{
		HonorLabels: true,
		MetricsPath: "/federate",
		BasicAuth:   make(map[string]string),
	}

	project, err := o.client.Project.ById(stack.AccountId)
	if err != nil {
		err = fmt.Errorf("failed to retrieve project `%s`: %s", stack.AccountId, err)
		return
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
			Labels: map[string]string{
				"rancher_site": strings.Split(project.Links["self"], "/")[2],
				"rancher_url":  project.Links["self"],
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
	return
}

// DiscoverTargets is parse all Cattle stacks to detect those which has
// the environment variable "PROMETHEUS_FQDN" defined
func (o *CattleOrchestrator) DiscoverTargets() (jobs []prometheus.ScrapeConfig, err error) {
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

	jobs = []prometheus.ScrapeConfig{}
	for _, stack := range stacks.Data {
		if stack.Environment["PROMETHEUS_FQDN"] != nil {
			target, err := o.buildTarget(stack)
			if err != nil {
				log.Errorf("failed to build target: %s", err)
			}
			jobs = append(jobs, target)
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
