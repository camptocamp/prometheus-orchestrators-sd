package orchestrators

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// Orchestrator implements a container Orchestrator interface
type Orchestrator interface {
	GetName() string
	Start() (err error)
}

// GetOrchestrator returns the Orchestrator based on configuration or envionment if not defined
func GetOrchestrator(orchestratorArg string, posdServer string, refreshInterval string) (orch Orchestrator, err error) {
	if orchestratorArg != "" {
		log.Debugf("Choosing orchestrator based on configuration...")
		switch orchestratorArg {
		case "cattle":
			orch = NewCattleOrchestrator(posdServer, refreshInterval)
		default:
			err = fmt.Errorf("'%s' is not a valid orchestrator", orchestratorArg)
			return
		}
	} else {
		log.Debugf("Detecting orchestrator based on environment...")
		if detectCattle() {
			orch = NewCattleOrchestrator(posdServer, refreshInterval)
		} else {
			err = fmt.Errorf("no orchestrator detected")
			return
		}
	}
	log.Debugf("Using orchestrator: %s", orch.GetName())
	return
}
