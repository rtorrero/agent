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
var whitelistedArguments = map[string]func([]byte) (entities.FactValue, *entities.FactGatheringError){
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

type SaptuneOutput struct {
	Schema      string    `json:"$schema"`
	PublishTime string    `json:"publish time"`
	Argv        string    `json:"argv"`
	Pid         int       `json:"pid"`
	Command     string    `json:"command"`
	ExitCode    int       `json:"exit code"`
	Result      Result    `json:"result"`
	Messages    []Message `json:"messages"`
}

type Result struct {
	Services                 Services `json:"services"`
	SystemdSystemState       string   `json:"systemd system state"`
	TuningState              string   `json:"tuning state"`
	Virtualization           string   `json:"virtualization"`
	ConfiguredVersion        string   `json:"configured version"`
	PackageVersion           string   `json:"package version"`
	SolutionEnabled          []string `json:"Solution enabled"`
	NotesEnabledBySolution   []string `json:"Notes enabled by Solution"`
	SolutionApplied          []string `json:"Solution applied"`
	NotesAppliedBySolution   []string `json:"Notes applied by Solution"`
	NotesEnabledAdditionally []string `json:"Notes enabled additionally"`
	NotesEnabled             []string `json:"Notes enabled"`
	NotesApplied             []string `json:"Notes applied"`
	Staging                  Staging  `json:"staging"`
	RememberMessage          string   `json:"remember message"`
}

type Services struct {
	Saptune []string `json:"saptune"`
	Sapconf []string `json:"sapconf"`
	Tuned   []string `json:"tuned"`
}

type Staging struct {
	StagingEnabled  bool     `json:"staging enabled"`
	NotesStaged     []string `json:"Notes staged"`
	SolutionsStaged []string `json:"Solutions staged"`
}

type Message struct {
	Priority string `json:"priority"`
	Message  string `json:"message"`
}

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

	return argumentHandler(saptuneOutput)
}

func gatherStatus(commandOutput []byte) (entities.FactValue, *entities.FactGatheringError) {
	status, err := parseJSONToFactValue(string(commandOutput))
	if err != nil {
		gatheringError := SaptuneCommandError.Wrap(err.Error())
		log.Error(gatheringError)
		return nil, gatheringError
	}

	return status, nil
}

func gatherSolutionVerify(commandOutput []byte) (entities.FactValue, *entities.FactGatheringError) {
	result := &entities.FactValueMap{}
	return result, nil
}
func gatherSolutionList(commandOutput []byte) (entities.FactValue, *entities.FactGatheringError) {
	result := &entities.FactValueMap{}
	return result, nil
}
func gatherNoteVerify(commandOutput []byte) (entities.FactValue, *entities.FactGatheringError) {
	result := &entities.FactValueMap{}
	return result, nil
}
func gatherNoteList(commandOutput []byte) (entities.FactValue, *entities.FactGatheringError) {
	result := &entities.FactValueMap{}
	return result, nil
}
