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
	var corosyncCmapCtlMap entities.FactValueMap

	for _, line := range lines {
		keyValue := strings.Split(line, " = ")
		keys := strings.Split(keyValue[0], ".")
		parseValue(&corosyncCmapCtlMap.Value, keys, keyValue[1])
		//corosyncCmapCtlMap.Value = parseValue(keys, keyValue[1])
	}

	return &corosyncCmapCtlMap, nil
}

func cleanKey(key string) string {
	return strings.Split(key, " ")[0]
}

func parseValue(fvm *map[string]entities.FactValue, keys []string, value string) {
	currentKey := keys[0]

	if len(keys) == 1 {
		(*fvm)[currentKey] = entities.ParseStringToFactValue(value)
		return
	}

	// DOES EXIST
	if currentValue, ok := (*fvm)[currentKey]; ok {
		var tacua2 = map[string]entities.FactValue{
			currentKey: currentValue,
		}

		parseValue(&currentValue, keys[1:], value)
		return
	}

	// DOES NOT EXIST
	fvmNew := &entities.FactValueMap{}
	(*fvm)[currentKey] = fvmNew
	parseValue(&fvmNew.Value, keys[1:], value)
}

// func parseValue(fvm *entities.FactValueMap, keys []string, value string) {
// 	if len(keys) == 1 {
// 		var newMap = make(map[string]entities.FactValue)
// 		newMap[cleanKey(keys[0])] = entities.ParseStringToFactValue(value)
// 		fvm.Value = newMap
// 		return
// 	}
// 	if _, ok := fvm.Value[keys[0]]; ok {
// 		// EXISTS!
// 	}

// 	// DOES NOT EXIST
// 	newFVM := &entities.FactValueMap{
// 		Value: map[string]entities.FactValue{
// 			keys[0]: &entities.FactValueMap{},
// 		},
// 	}
// 	fvm.Value = newFVM.Value
// 	parseValue(&fvm.Value[keys[0]], keys[1:], value)

// 	// fvm.Value[keys[0]]

// 	// if _, ok := fvm.Value[keys[0]]; ok {
// 	// 	cosita := &entities.FactValue{fvm.Value[keys[0]]
// 	// 	if fvmnew, ok := fvm.Value[keys[0]].(entities.FactValueMap); ok {
// 	// 		parseValue(&fvm.Value[keys[0]], keys[1:], value)
// 	// 	}
// 	// }

// 	return
// }

// func parseValue(fvm *entities.FactValueMap, keys []string, value string) {
// 	var outputMap = make(map[string]entities.FactValue)
// 	//var fvmChild *entities.FactValueMap

// 	if len(keys) < 2 {
// 		innerMostKey := strings.Split(keys[0], " ")[0]
// 		outputMap = make(map[string]entities.FactValue)
// 		outputMap[innerMostKey] = entities.ParseStringToFactValue(value)
// 		fvm.Value = outputMap
// 	} else {
// 		reducedKeys := keys[1:]
// 		var fvmChild entities.FactValue
// 		fvm.Value = map[string]entities.FactValue{keys[0]: fvmChild}

// 		parseValue(fvm.Value[keys[0]], reducedKeys, value)
// 		//outputMap[keys[0]] = &entities.FactValueMap{Value: parseValue(fvmChild, reducedKeys, value)}
// 	}

// 	return
// }
