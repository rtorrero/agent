package gatherers

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/agent/internal/core/saptune"
	"github.com/trento-project/agent/pkg/factsengine/entities"
	"github.com/trento-project/agent/pkg/utils"
)

const (
	SaptuneGathererName = "saptune"
)

// nolint:gochecknoglobals
var whitelistedArguments = map[string]string{
	"status":          "status --non-compliance-check",
	"solution-verify": "solution verify",
	"solution-list":   "solution list",
	"note-verify":     "note verify",
	"note-list":       "note list",
}

// nolint:gochecknoglobals
var (
	SaptuneVersionUnsupported = entities.FactGatheringError{
		Type:    "saptune-version-not-supported",
		Message: "currently installed version of saptune is not supported",
	}

	SaptuneUnknownArgument = entities.FactGatheringError{
		Type:    "saptune-unknown-error",
		Message: "the requested argument is not currently supported",
	}

	SaptuneMissingArgument = entities.FactGatheringError{
		Type:    "saptune-missing-argument",
		Message: "missing required argument",
	}

	SaptuneCommandError = entities.FactGatheringError{
		Type:    "saptune-cmd-error",
		Message: "error executing saptune command",
	}

)

type SaptuneGatherer struct {
	executor utils.CommandExecutor
}

func NewDefaultSaptuneGatherer() *SaptuneGatherer {
	return NewSaptuneGatherer(utils.Executor{})
}

func NewSaptuneGatherer(executor utils.CommandExecutor) *SaptuneGatherer {
	return &SaptuneGatherer{
		executor: executor,
	}
}

func parseJSONToFactValue(jsonStr string) (entities.FactValue, error) {
	// Unmarshal the JSON into an interface{} type.
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, err
	}

	// Convert the parsed jsonData into a FactValue using NewFactValue.
	return entities.NewFactValueMod(jsonData)
}

func (s *SaptuneGatherer) Gather(factsRequests []entities.FactRequest) ([]entities.Fact, error) {
	facts := []entities.Fact{}
	log.Infof("Starting %s facts gathering process", SaptuneGathererName)
	saptuneRetriever, _ := saptune.NewSaptune(s.executor)
	for _, factReq := range factsRequests {
			var fact entities.Fact
		if len(factReq.Argument) == 0 {
			log.Error(SaptuneMissingArgument.Message)
			fact = entities.NewFactGatheredWithError(factReq, &SaptuneMissingArgument)
		} else if factValue, err := handleArgument(&saptuneRetriever, factReq.Argument); err != nil {
			fact = entities.NewFactGatheredWithError(factReq, err)
		} else {
			fact = entities.NewFactGatheredWithRequest(factReq, factValue)
		}
		facts = append(facts, fact)
	}

	log.Infof("Requested %s facts gathered", SaptuneGathererName)
	return facts, nil
}

func handleArgument(
	saptuneRetriever *saptune.Saptune,
	argument string,
) (entities.FactValue, *entities.FactGatheringError) { 
	internalArguments, ok := whitelistedArguments[argument]

	if !ok {
		gatheringError := SaptuneUnknownArgument.Wrap(internalArguments)
		log.Error(gatheringError)
		return nil, gatheringError
	}

	argList := strings.Split(internalArguments, " ")
	saptuneOutput, commandError := saptuneRetriever.RunCommandJSON(argList...)
	if commandError != nil {
		gatheringError := SaptuneCommandError.Wrap(commandError.Error())
		log.Error(gatheringError)
		return nil, gatheringError
	}

	return gatherFactsFromOutput(saptuneOutput)
}

func gatherFactsFromOutput(commandOutput []byte) (entities.FactValue, *entities.FactGatheringError) {
	status, err := parseJSONToFactValue(string(commandOutput))
	if err != nil {
		gatheringError := SaptuneCommandError.Wrap(err.Error())
		log.Error(gatheringError)
		return nil, gatheringError
	}

	return status, nil
}
