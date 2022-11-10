package gatherers

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/agent/internal/factsengine/entities"
	"github.com/trento-project/agent/test/helpers"
)

type HostsFileTestSuite struct {
	suite.Suite
	fileSystem afero.Fs
}

func TestHostsFileTestSuite(t *testing.T) {
	suite.Run(t, new(HostsFileTestSuite))
}

// func (suite *HostsFileTestSuite) SetupTest() {
// 	suite.fileSystem = afero.NewMemMapFs()

// 	err := suite.fileSystem.MkdirAll("/etc", 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func (suite *HostsFileTestSuite) TestHostsFileBasic() {
	c := NewHostsFileGatherer(helpers.GetFixturePath("gatherers/hosts.basic"))

	factRequests := []entities.FactRequest{
		{
			Name:     "hosts_localhost",
			Gatherer: "hosts",
			Argument: "localhost",
			CheckID:  "check1",
		},
		{
			Name:     "hosts_somehost",
			Gatherer: "hosts",
			Argument: "somehost",
			CheckID:  "check2",
		},
		{
			Name:     "hosts_ip6-localhost",
			Gatherer: "hosts",
			Argument: "ip6-localhost",
			CheckID:  "check3",
		},
	}

	factResults, err := c.Gather(factRequests)

	expectedResults := []entities.Fact{
		{
			Name:    "hosts_localhost",
			Value:   []string{"127.0.0.1", "::1"},
			CheckID: "check1",
		},
		{
			Name:    "hosts_somehost",
			Value:   []string{"127.0.1.1"},
			CheckID: "check2",
		},
		{
			Name:    "hosts_ip6-localhost",
			Value:   []string{"::1"},
			CheckID: "check3",
		},
	}

	suite.NoError(err)
	suite.ElementsMatch(expectedResults, factResults)
}

func (suite *HostsFileTestSuite) TestHostsFileNotExists() {
	c := NewHostsFileGatherer("non_existing_file")

	factRequests := []entities.FactRequest{
		{
			Name:     "hosts_somehost",
			Gatherer: "hosts",
			Argument: "somehost",
		},
	}

	_, err := c.Gather(factRequests)

	suite.EqualError(err, "could not open /etc/hosts file: could not open corosync.conf file: open non_existing_file: no such file or directory")
}

func (suite *HostsFileTestSuite) TestHostsFileIgnoresCommentedHosts() {

	c := NewHostsFileGatherer(helpers.GetFixturePath("gatherers/hosts.basic"))

	factRequests := []entities.FactRequest{
		{
			Name:     "hosts_commented-host",
			Gatherer: "hosts",
			Argument: "commented-host",
		},
	}

	factResults, err := c.Gather(factRequests)

	expectedResults := []entities.Fact{}

	suite.NoError(err)
	suite.ElementsMatch(expectedResults, factResults)
}
