package gatherers

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/trento-project/agent/internal/utils"
)

const (
	HostsFileGathererName = "hosts"
	hostsFilePath         = "/etc/hosts"
)

type HostsFileGatherer struct {
	fileSystem afero.Fs
}

func NewHostsFileGatherer(fileSystem afero.Fs) *HostsFileGatherer {
	return &HostsFileGatherer{fileSystem: fileSystem}
}

func (s *HostsFileGatherer) Gather(factsRequests []FactRequest) ([]Fact, error) {
	facts := []Fact{}
	log.Infof("Starting /etc/hosts file facts gathering process")

	hostsFile, err := s.fileSystem.Open(hostsFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not open /etc/hosts file")
	}

	defer hostsFile.Close()

	hostsFileContents, err := ioutil.ReadAll(hostsFile)
	if err != nil {
		return facts, err
	}

	hostsMap := utils.FindMatches(`(?m)^([^#|\s]+)\s([^#\n]+)#?.*$`, hostsFileContents)

	for _, factReq := range factsRequests {
		factValue := getIPAddressesFromHostname(hostsMap, factReq.Argument)
		if len(factValue) < 1 {
			continue
		}
		fact := NewFactWithRequest(factReq, factValue)
		facts = append(facts, fact)
	}

	log.Infof("Requested /etc/hosts file facts gathered")
	return facts, nil
}

func getIPAddressesFromHostname(hostsMap map[string]interface{}, requestedHostname string) []string {
	var matchingValues []string
	for key, value := range hostsMap {
		valueString := fmt.Sprintf("%v", value)
		valueString = strings.TrimLeft(valueString, " ")
		hostnames := strings.Split(valueString, " ")
		for _, hostname := range hostnames {
			if hostname == requestedHostname {
				matchingValues = append(matchingValues, key)
			}
		}
	}
	sort.Strings(matchingValues)

	return matchingValues
}
