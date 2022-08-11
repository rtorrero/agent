package gatherers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HostsFileTestSuite struct {
	suite.Suite
	fileSystem afero.Fs
}

func TestHostsFileTestSuite(t *testing.T) {
	suite.Run(t, new(HostsFileTestSuite))
}

func (suite *HostsFileTestSuite) SetupTest() {
	suite.fileSystem = afero.NewMemMapFs()

	err := suite.fileSystem.MkdirAll("/etc", 0644)
	if err != nil {
		panic(err)
	}
}

func (suite *HostsFileTestSuite) TestHostsFileBasic() {
	testFile, _ := os.Open("../../../test/fixtures/gatherers/hosts.basic")
	confFile, _ := ioutil.ReadAll(testFile)
	err := afero.WriteFile(suite.fileSystem, "/etc/hosts", confFile, 0644)
	assert.NoError(suite.T(), err)
	c := NewHostsFileGatherer(suite.fileSystem)

	factRequests := []FactRequest{
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

	expectedResults := []Fact{
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
	c := NewHostsFileGatherer(suite.fileSystem)

	factRequests := []FactRequest{
		{
			Name:     "hosts_somehost",
			Gatherer: "hosts",
			Argument: "somehost",
		},
	}

	_, err := c.Gather(factRequests)

	suite.EqualError(err, "could not open /etc/hosts file: open /etc/hosts: file does not exist")
}

func (suite *HostsFileTestSuite) TestHostsFileIgnoresCommentedHosts() {
	testFile, _ := os.Open("../../../test/fixtures/gatherers/hosts.basic")
	confFile, _ := ioutil.ReadAll(testFile)
	err := afero.WriteFile(suite.fileSystem, "/etc/hosts", confFile, 0644)
	assert.NoError(suite.T(), err)
	c := NewHostsFileGatherer(suite.fileSystem)

	factRequests := []FactRequest{
		{
			Name:     "hosts_commented-host",
			Gatherer: "hosts",
			Argument: "commented-host",
		},
	}

	factResults, err := c.Gather(factRequests)

	expectedResults := []Fact{}

	suite.NoError(err)
	suite.ElementsMatch(expectedResults, factResults)
}
