package gatherers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/agent/internal/factsengine/entities"
)

const (
	HostsFileFactKey = "hosts"
	HostsFilePath    = "/etc/hosts"
)

var (
	LineIPv4Compiled = regexp.MustCompile(`(?P<ip>[^\s]+)\s+(?P<hostname>[^\s]+)`) //subs := r.FindStringSubmatch(s)
)

// nolint:gochecknoglobals
var (
	HostsFileError = entities.FactGatheringError{
		Type:    "hosts-file-error",
		Message: "error reading /etc/hosts file",
	}

	HostsFileDecodingError = entities.FactGatheringError{
		Type:    "hosts-file-decoding-error",
		Message: "error decoding /etc/hosts file",
	}

	HostsFileNotFoundError = entities.FactGatheringError{
		Type:    "hosts-file-value-not-found",
		Message: "requested field value not found in /etc/hosts file",
	}
)

type HostsFileGatherer struct {
	hostsFile string
}

func NewDefaultHostsFileGatherer() *HostsFileGatherer {
	return NewHostsFileGatherer(HostsFilePath)
}

func NewHostsFileGatherer(hostsFile string) *HostsFileGatherer {
	return &HostsFileGatherer{hostsFile: hostsFile}
}

func (s *HostsFileGatherer) Gather(factsRequests []entities.FactRequest) ([]entities.Fact, error) {
	facts := []entities.Fact{}
	log.Infof("Starting /etc/hosts file facts gathering process")

	hostsFile, err := readHostsFileByLines(s.hostsFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not open /etc/hosts file")
	}

	hostsFileMap, err := hostsFileToMap(hostsFile)
	if err != nil {
		return nil, HostsFileDecodingError.Wrap(err.Error())
	}

	for _, factReq := range factsRequests {
		var fact entities.Fact

		for hostname, ip := range hostsFileMap {
			if hostname == factReq.Argument {
				fact = entities.NewFactGatheredWithRequest(factReq, ip)
				facts = append(facts, fact)
				break
			}
		}
	}

	log.Infof("Requested /etc/hosts file facts gathered")
	return facts, nil
}

func readHostsFileByLines(filePath string) ([]string, error) {
	hostsFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not open corosync.conf file")
	}

	defer hostsFile.Close()

	fileScanner := bufio.NewScanner(hostsFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		if strings.HasPrefix(fileScanner.Text(), "#") || fileScanner.Text() == "" {
			continue
		}
		fileLines = append(fileLines, fileScanner.Text())
	}

	return fileLines, nil
}

func hostsFileToMap(lines []string) (map[string]string, error) {
	var hostsFileMap = make(map[string]string)
	var paramsMap = make(map[string]string)

	for _, line := range lines {
		match := LineIPv4Compiled.FindStringSubmatch(line)

		if match == nil {
			return nil, fmt.Errorf("invalid hosts file structure")
		}
		for i, name := range LineIPv4Compiled.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
		hostsFileMap[paramsMap["hostname"]] = paramsMap["ip"]
	}

	return hostsFileMap, nil
}
