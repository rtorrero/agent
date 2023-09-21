package gatherers_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/agent/internal/factsengine/gatherers"
	"github.com/trento-project/agent/pkg/factsengine/entities"
	utilsMocks "github.com/trento-project/agent/pkg/utils/mocks"
	"github.com/trento-project/agent/test/helpers"
)

type SaptuneTestSuite struct {
	suite.Suite
	mockExecutor *utilsMocks.CommandExecutor
}

func TestSaptuneTestSuite(t *testing.T) {
	suite.Run(t, new(SaptuneTestSuite))
}

func (suite *SaptuneTestSuite) SetupTest() {
	suite.mockExecutor = new(utilsMocks.CommandExecutor)
}

func (suite *SaptuneTestSuite) TestSaptuneGathererStatus() {
	mockOutputFile, _ := os.Open(helpers.GetFixturePath("gatherers/saptune-status.output"))
	mockOutput, _ := io.ReadAll(mockOutputFile)
	suite.mockExecutor.On("Exec", "saptune", "status", "--no-compliance-check").Return(mockOutput, nil)
	c := gatherers.NewSaptuneGatherer(suite.mockExecutor)

	factRequests := []entities.FactRequest{
		{
			Name:     "saptune_status",
			Gatherer: "saptune",
			Argument: "status --non-compliance-check",
		},
	}

	factResults, err := c.Gather(factRequests)

	expectedResults := []entities.Fact{
		{
			Name:  "saptune_status",
			Value: &entities.FactValueMap{
				Value: map[string]entities.FactValue{
					"$schema": &entities.FactValueString{Value: "file:///usr/share/saptune/schemas/1.0/saptune_status.schema.json"},
					"publish time": &entities.FactValueString{Value: "2023-09-15 15:15:14.599"},
					"argv": &entities.FactValueString{Value: "saptune --format json status"},
					"pid": &entities.FactValueInt{Value: 6593},
					"command": &entities.FactValueString{Value: "status"},
					"exit code": &entities.FactValueInt{Value: 1},
					"result": &entities.FactValueMap{
						Value: map[string]entities.FactValue{
							"services": &entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"saptune": &entities.FactValueList{
										Value: []entities.FactValue{
											&entities.FactValueString{Value: "disabled"},
											&entities.FactValueString{Value: "inactive"},
										},
									},
									"sapconf": &entities.FactValueList{Value: []entities.FactValue{}},
									"tuned": &entities.FactValueList{Value: []entities.FactValue{}},
								},
							},
							"systemd system state": &entities.FactValueString{Value: "degraded"},
							"tuning state": &entities.FactValueString{Value: "compliant"},
							"virtualization": &entities.FactValueString{Value: "kvm"},
							"configured version": &entities.FactValueString{Value: "3"},
							"package version": &entities.FactValueString{Value: "3.1.0"},
							"Solution enabled": &entities.FactValueList{Value: []entities.FactValue{}},
							"Notes enabled by Solution": &entities.FactValueList{Value: []entities.FactValue{}},
							"Solution applied": &entities.FactValueList{Value: []entities.FactValue{}},
							"Notes applied by Solution": &entities.FactValueList{Value: []entities.FactValue{}},
							"Notes enabled additionally": &entities.FactValueList{
								Value: []entities.FactValue{
									&entities.FactValueString{Value: "1410736"},
								},
							},
							"Notes enabled": &entities.FactValueList{
								Value: []entities.FactValue{
									&entities.FactValueString{Value: "1410736"},
								},
							},
							"Notes applied": &entities.FactValueList{
								Value: []entities.FactValue{
									&entities.FactValueString{Value: "1410736"},
								},
							},
							"staging": &entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"staging enabled": &entities.FactValueBool{Value: false},
									"Notes staged": &entities.FactValueList{Value: []entities.FactValue{}},
									"Solutions staged": &entities.FactValueList{Value: []entities.FactValue{}},
								},
							},
							"remember message": &entities.FactValueString{Value: "\nRemember: if you wish to automatically activate the note's and solution's tuning options after a reboot, you must enable and start saptune.service by running:\n 'saptune service enablestart'.\nThe systemd system state is NOT ok.\nPlease call 'saptune check' to get guidance to resolve the issues!\n\n"},
						},
					},
					"messages": &entities.FactValueList{
						Value: []entities.FactValue{
							&entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"priority": &entities.FactValueString{Value: "NOTICE"},
									"message": &entities.FactValueString{Value: "actions.go:85: ATTENTION: You are running a test version (3.1.0 from 2023/08/03) of saptune which is not supported for production use\n"},
								},
							},
						},
					},
				},
			},
		},
	}

	suite.NoError(err)
	suite.ElementsMatch(expectedResults, factResults)
}
func (suite *SaptuneTestSuite) TestSaptuneNoArgumentProvided() {
	c := gatherers.NewSaptuneGatherer(suite.mockExecutor)

	factRequests := []entities.FactRequest{
		{
			Name:     "no_argument_fact",
			Gatherer: "saptune",
		},
		{
			Name:     "empty_argument_fact",
			Gatherer: "saptune",
			Argument: "",
		},
	}

	factResults, err := c.Gather(factRequests)

	expectedResults := []entities.Fact{
		{
			Name:  "no_argument_fact",
			Value: nil,
			Error: &entities.FactGatheringError{
				Message: "missing required argument",
				Type:    "saptune-missing-argument",
			},
		},
		{
			Name:  "empty_argument_fact",
			Value: nil,
			Error: &entities.FactGatheringError{
				Message: "missing required argument",
				Type:    "saptune-missing-argument",
			},
		},
	}

	suite.NoError(err)
	suite.ElementsMatch(expectedResults, factResults)
}
