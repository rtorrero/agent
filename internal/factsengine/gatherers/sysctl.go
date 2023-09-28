package gatherers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/agent/pkg/factsengine/entities"
	"github.com/trento-project/agent/pkg/utils"
)

const (
	SysctlGathererName = "sysctl"
)

// nolint:gochecknoglobals
var (
	SysctlCommandError = entities.FactGatheringError{
		Type:    "sysctl-cmd-error",
		Message: "error executing sysctl command",
	}

	SysctlMissingArgument = entities.FactGatheringError{
		Type:    "sysctl-missing-argument",
		Message: "missing required argument",
	}
)

type SysctlGatherer struct {
	executor utils.CommandExecutor
}

func NewDefaultSysctlGatherer() *SysctlGatherer {
	return NewSysctlGatherer(utils.Executor{})
}

func NewSysctlGatherer(executor utils.CommandExecutor) *SysctlGatherer {
	return &SysctlGatherer{
		executor: executor,
	}
}

func (s *SysctlGatherer) Gather(factsRequests []entities.FactRequest) ([]entities.Fact, error) {
	facts := []entities.Fact{}
	log.Infof("Starting sysctl facts gathering process")

	for _, factReq := range factsRequests {
		var fact entities.Fact
		if len(factReq.Argument) == 0 {
			log.Error(SysctlMissingArgument.Message)
			fact = entities.NewFactGatheredWithError(factReq, &SysctlMissingArgument)
		} else if factValue, err := runCommand(s.executor, factReq.Argument); err != nil {
			fact = entities.NewFactGatheredWithError(factReq, err)
		} else {
			fact = entities.NewFactGatheredWithRequest(factReq, factValue)
		}

		facts = append(facts, fact)
	}

	log.Infof("Requested %s facts gathered", SysctlGathererName)
	return facts, nil
}

func runCommand(executor utils.CommandExecutor, argument string) (entities.FactValue, error) {
	sysctlOutput, commandError := saptuneRetriever.RunCommandJSON(arguments...)
	if commandError != nil {
		return nil, commandError
	}

	var jsonData interface{}
	if err := json.Unmarshal(saptuneOutput, &jsonData); err != nil {
		return nil, err
	}

	return entities.NewFactValue(jsonData, entities.WithSnakeCaseKeys())
}