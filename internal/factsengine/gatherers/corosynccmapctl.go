package gatherers

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/agent/internal/factsengine/entities"
	"github.com/trento-project/agent/internal/utils"
)

const (
	CorosyncCmapCtlGathererName = "corosync-cmapctl"
)

// nolint:gochecknoglobals
var (
	CorosyncCmapCtlValueNotFound = entities.FactGatheringError{
		Type:    "corosync-cmapctl-value-not-found",
		Message: "requested value not found in corosync-cmapctl output",
	}

	CorosyncCmapCtlCommandError = entities.FactGatheringError{
		Type:    "corosync-cmapctl-command-error",
		Message: "error while executing corosynccmap-ctl",
	}
)

type CorosyncCmapctlGatherer struct {
	executor utils.CommandExecutor
}

func NewDefaultCorosyncCmapctlGatherer() *CorosyncCmapctlGatherer {
	return NewCorosyncCmapctlGatherer(utils.Executor{})
}

func NewCorosyncCmapctlGatherer(executor utils.CommandExecutor) *CorosyncCmapctlGatherer {
	return &CorosyncCmapctlGatherer{
		executor: executor,
	}
}

func (s *CorosyncCmapctlGatherer) Gather(factsRequests []entities.FactRequest) ([]entities.Fact, error) {
	facts := []entities.Fact{}
	log.Infof("Starting %s facts gathering process", CorosyncCmapCtlGathererName)

	corosyncCmapctl, err := s.executor.Exec(
		"corosync-cmapctl", "-b")
	if err != nil {
		return nil, CorosyncCmapCtlCommandError.Wrap(err.Error())
	}

	corosyncCmapctlMap := utils.FindMatches(`(?m)^(\S*)\s\(\S*\)\s=\s(.*)$`, corosyncCmapctl)
	outputBytes := readCmapCtlOutputByLines(corosyncCmapctl)
	alternativeMap, _ := cmapCtlOutputToMap(outputBytes)

	log.Info("If this works im spiderman on top of a horse: ", alternativeMap)
	for _, factReq := range factsRequests {
		var fact entities.Fact

		if value, ok := corosyncCmapctlMap[factReq.Argument]; ok {
			fact = entities.NewFactGatheredWithRequest(factReq, &entities.FactValueString{Value: fmt.Sprint(value)})
		} else {
			gatheringError := CorosyncCmapCtlValueNotFound.Wrap(factReq.Argument)
			log.Error(gatheringError)
			fact = entities.NewFactGatheredWithError(factReq, gatheringError)
		}
		facts = append(facts, fact)
	}

	log.Infof("Requested %s facts gathered", CorosyncCmapCtlGathererName)
	return facts, nil
}

func readCmapCtlOutputByLines(data []byte) []string {
	reader := bytes.NewReader(data)
	outputScanner := bufio.NewScanner(reader)

	outputScanner.Split(bufio.ScanLines)
	var fileLines []string

	for outputScanner.Scan() {
		scannedLine := outputScanner.Text()
		if strings.HasPrefix(scannedLine, "#") || scannedLine == "" {
			continue
		}
		fileLines = append(fileLines, scannedLine)
	}

	return fileLines
}

func cmapCtlOutputToMap(lines []string) (*entities.FactValueMap, error) {
	//var corosyncCmapCtlMap = make(map[string]entities.FactValue)

	var corosyncCmapCtlMap *entities.FactValueMap

	for _, line := range lines {
		keys := strings.Split(line, ".")
		corosyncCmapCtlMap, _ = parseValue(keys)
	}
	return corosyncCmapCtlMap, nil

	// for _, line := range lines {
	// 	keys := strings.Split(line, ".")
	// 	var p *FactValue
	// 	if len(keys) > 0 {
	// 		corosyncCmapCtlMap[keys[0]] =
	// 		p = &corosyncCmapCtlMap[keys[0]]
	// 	}
	// 	for i, key := range keys {
	// 		corosyncCmapCtlMap[] :=
	// 		if i >= len(keys) {

	// 		}
	// 	}
	// }

}

func parseValue(keys []string) (*entities.FactValueMap, error) {
	var outputMap = make(map[string]entities.FactValue)
	var resultMap entities.FactValueMap

	if len(keys) < 3 {
		innerMostKeyValue := strings.Split(keys[len(keys)-1], " = ")
		innerMostKey := strings.Split(innerMostKeyValue[0], " ")[0]
		outputMap = make(map[string]entities.FactValue)
		outputMap[innerMostKey] = entities.ParseStringToFactValue(innerMostKeyValue[1])

		resultMap.Value = outputMap
	} else {
		//factValueMap.Value = make(map[string]entities.FactValue)
		reducedKeys := keys[1:]
		var cacota *entities.FactValueMap
		cacota, _ = parseValue(reducedKeys)
		resultMap.Value = cacota.Value
	}

	return &resultMap, nil
}
