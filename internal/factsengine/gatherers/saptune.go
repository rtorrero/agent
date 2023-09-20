package gatherers

import (
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
var whitelistedArguments = map[string]func(string) (entities.FactValue, *entities.FactGatheringError){
	"status --non-compliance-check":          gatherStatus,
	"solution-verify": gatherSolutionVerify,
	"solution-list":   gatherSolutionList,
	"note-verify":     gatherNoteVerify,
	"note-list":       gatherNoteList,
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

func (s *SaptuneGatherer) Gather(factsRequests []entities.FactRequest) ([]entities.Fact, error) {
	facts := []entities.Fact{}
	log.Infof("Starting %s facts gathering process", SaptuneGathererName)
	saptuneRetriever, _ := saptune.NewSaptune(utils.Executor{})
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
	argumentHandler, ok := whitelistedArguments[argument]

	if !ok {
		gatheringError := SaptuneUnknownArgument.Wrap(argument)
		log.Error(gatheringError)
		return nil, gatheringError
	}

	argList := strings.Split(argument, " ")
	saptuneOutput, commandError := saptuneRetriever.RunCommandJSON(argList...)
	if commandError != nil {
		gatheringError := SaptuneCommandError.Wrap(commandError.Error())
		log.Error(gatheringError)
		return nil, gatheringError
	}

	return argumentHandler(string(saptuneOutput))
}

func gatherStatus(commandOutput string) (entities.FactValue, *entities.FactGatheringError) {
	// profit?
}

func gatherSolutionVerify(commandOutput string) (entities.FactValue, *entities.FactGatheringError) {
	// profit?
}
func gatherSolutionList(commandOutput string) (entities.FactValue, *entities.FactGatheringError) {
	// profit?
}
func gatherNoteVerify(commandOutput string) (entities.FactValue, *entities.FactGatheringError) {
	// profit?
}
func gatherNoteList(commandOutput string) (entities.FactValue, *entities.FactGatheringError) {
	// profit?
}
