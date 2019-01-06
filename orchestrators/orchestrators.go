package orchestrators

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/camptocamp/prometheus-orchestrators-sd/prometheus"
)

// Orchestrator implements a container Orchestrator interface
type Orchestrator interface {
	GetName() string
	DiscoverTargets() (jobs []prometheus.ScrapeConfig, err error)
}

// GetOrchestrator returns the Orchestrator based on configuration or envionment if not defined
func GetOrchestrator(orchestratorArg string) (orch Orchestrator, err error) {
	if orchestratorArg != "" {
		log.Debugf("Choosing orchestrator based on configuration...")
		switch orchestratorArg {
		case "cattle":
			orch = NewCattleOrchestrator()
		default:
			err = fmt.Errorf("'%s' is not a valid orchestrator", orchestratorArg)
			return
		}
	} else {
		log.Debugf("Detecting orchestrator based on environment...")
		if detectCattle() {
			orch = NewCattleOrchestrator()
		} else {
			err = fmt.Errorf("no orchestrator detected")
			return
		}
	}
	log.Debugf("Using orchestrator: %s", orch.GetName())
	return
}
