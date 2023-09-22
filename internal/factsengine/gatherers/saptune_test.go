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
	suite.mockExecutor.On("Exec", "saptune", "--format", "json", "status", "--non-compliance-check").Return(mockOutput, nil)
	suite.mockExecutor.On("Exec", "rpm", "-q", "--qf", "%{VERSION}", "saptune").Return(
		[]byte("3.1.0"), nil,
	)
	c := gatherers.NewSaptuneGatherer(suite.mockExecutor)

	factRequests := []entities.FactRequest{
		{
			Name:     "saptune_status",
			Gatherer: "saptune",
			Argument: "status",
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

func (suite *SaptuneTestSuite) TestSaptuneGathererNoteVerify() {
	mockOutputFile, _ := os.Open(helpers.GetFixturePath("gatherers/saptune-note-verify.output"))
	mockOutput, _ := io.ReadAll(mockOutputFile)
	suite.mockExecutor.On("Exec", "saptune", "--format", "json", "note", "verify").Return(mockOutput, nil)
	suite.mockExecutor.On("Exec", "rpm", "-q", "--qf", "%{VERSION}", "saptune").Return(
		[]byte("3.1.0"), nil,
	)
	c := gatherers.NewSaptuneGatherer(suite.mockExecutor)

	factRequests := []entities.FactRequest{
		{
			Name:     "saptune_note_verify",
			Gatherer: "saptune",
			Argument: "note-verify",
		},
	}

	factResults, err := c.Gather(factRequests)

	expectedResults := []entities.Fact{
		{
			Name:  "saptune_note_verify",
			Value: &entities.FactValueMap{
				Value: map[string]entities.FactValue{
					"$schema": &entities.FactValueString{
						Value: "file:///usr/share/saptune/schemas/1.0/saptune_note_verify.schema.json",
					},
					"publish time": &entities.FactValueString{
						Value: "2023-04-24 15:49:43.399",
					},
					"argv": &entities.FactValueString{
						Value: "saptune --format json note verify",
					},
					"pid": &entities.FactValueInt{
						Value: 25202,
					},
					"command": &entities.FactValueString{
						Value: "note verify",
					},
					"exit code": &entities.FactValueInt{
						Value: 1,
					},
					"result": &entities.FactValueMap{
						Value: map[string]entities.FactValue{
							"verifications": &entities.FactValueList{
								// Note: Due to the length of the list, only the first verification is represented.
								//       You can follow the same pattern for other items in the list.
								Value: []entities.FactValue{
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "1771258"},
											"Note version":   &entities.FactValueString{Value: "6"},
											"parameter":      &entities.FactValueString{Value: "LIMIT_@dba_hard_nofile"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "@dba hard nofile 1048576"},
											"actual value":   &entities.FactValueString{Value: "@dba hard nofile 1048576"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "1771258"},
											"Note version":   &entities.FactValueString{Value: "6"},
											"parameter":      &entities.FactValueString{Value: "LIMIT_@dba_soft_nofile"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "@dba soft nofile 1048576"},
											"actual value":   &entities.FactValueString{Value: "@dba soft nofile 1048576"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "1771258"},
											"Note version":   &entities.FactValueString{Value: "6"},
											"parameter":      &entities.FactValueString{Value: "LIMIT_@sapsys_hard_nofile"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "@sapsys hard nofile 1048576"},
											"actual value":   &entities.FactValueString{Value: "@sapsys hard nofile 1048576"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "1771258"},
											"Note version":   &entities.FactValueString{Value: "6"},
											"parameter":      &entities.FactValueString{Value: "LIMIT_@sapsys_soft_nofile"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "@sapsys soft nofile 1048576"},
											"actual value":   &entities.FactValueString{Value: "@sapsys soft nofile 1048576"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "1771258"},
											"Note version":   &entities.FactValueString{Value: "6"},
											"parameter":      &entities.FactValueString{Value: "LIMIT_@sdba_hard_nofile"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "@sdba hard nofile 1048576"},
											"actual value":   &entities.FactValueString{Value: "@sdba hard nofile 1048576"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "1771258"},
											"Note version":   &entities.FactValueString{Value: "6"},
											"parameter":      &entities.FactValueString{Value: "LIMIT_@sdba_soft_nofile"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "@sdba soft nofile 1048576"},
											"actual value":   &entities.FactValueString{Value: "@sdba soft nofile 1048576"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "IO_SCHEDULER_sda"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "none"},
											"actual value":   &entities.FactValueString{Value: "none"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "IO_SCHEDULER_sdb"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "none"},
											"actual value":   &entities.FactValueString{Value: "none"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "IO_SCHEDULER_sdc"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "none"},
											"actual value":   &entities.FactValueString{Value: "none"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "IO_SCHEDULER_sdd"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "none"},
											"actual value":   &entities.FactValueString{Value: "none"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "rpm:psmisc"},
											"compliant":      &entities.FactValueBool{Value: false},
											"expected value": &entities.FactValueString{Value: "23.0-6.16.1"},
											"actual value":   &entities.FactValueString{Value: "23.0-6.13.1"},
											"amendments": &entities.FactValueList{
												Value: []entities.FactValue{
													&entities.FactValueMap{
														Value: map[string]entities.FactValue{
															"index": &entities.FactValueInt{Value: 3},
															"amendment": &entities.FactValueString{Value: " [3] value is only checked, but NOT set"},
														},
													},
												},
											},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "rpm:uuidd"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "2.36.2-150300.4.17.1"},
											"actual value":   &entities.FactValueString{Value: "2.36.2-150300.4.23.1"},
											"amendments": &entities.FactValueList{
												Value: []entities.FactValue{
													&entities.FactValueMap{
														Value: map[string]entities.FactValue{
															"index": &entities.FactValueInt{Value: 3},
															"amendment": &entities.FactValueString{Value: " [3] value is only checked, but NOT set"},
														},
													},
												},
											},											
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "systemd:sysstat.service"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "start"},
											"actual value":   &entities.FactValueString{Value: "start, disable"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "systemd:uuidd.socket"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "start"},
											"actual value":   &entities.FactValueString{Value: "start, enable"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "vm.dirty_background_bytes"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "314572800"},
											"actual value":   &entities.FactValueString{Value: "314572800"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "vm.dirty_bytes"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "629145600"},
											"actual value":   &entities.FactValueString{Value: "629145600"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "xfsopt_barrier"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "-barrier"},
											"actual value":   &entities.FactValueString{Value: "-barrier"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "2578899"},
											"Note version":   &entities.FactValueString{Value: "41"},
											"parameter":      &entities.FactValueString{Value: "xfsopt_nobarrier"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "-nobarrier"},
											"actual value":   &entities.FactValueString{Value: "-nobarrier"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "900929"},
											"Note version":   &entities.FactValueString{Value: "7"},
											"parameter":      &entities.FactValueString{Value: "vm.max_map_count"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "2147483647"},
											"actual value":   &entities.FactValueString{Value: "2147483647"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "941735"},
											"Note version":   &entities.FactValueString{Value: "11"},
											"parameter":      &entities.FactValueString{Value: "ShmFileSystemSizeMB"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "725"},
											"actual value":   &entities.FactValueString{Value: "725"},
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "941735"},
											"Note version":   &entities.FactValueString{Value: "11"},
											"parameter":      &entities.FactValueString{Value: "VSZ_TMPFS_PERCENT"},
											"expected value": &entities.FactValueString{Value: "75"},
											"actual value":   &entities.FactValueString{Value: "75"},
											"amendments": &entities.FactValueList{
												Value: []entities.FactValue{
													&entities.FactValueMap{
														Value: map[string]entities.FactValue{
															"index": &entities.FactValueInt{Value: 15},
															"amendment": &entities.FactValueString{Value: "[15] the parameter is only used to calculate the size of tmpfs (/dev/shm)"},
														},
													},
												},
											},												
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "941735"},
											"Note version":   &entities.FactValueString{Value: "11"},
											"parameter":      &entities.FactValueString{Value: "kernel.shmall"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "1152921504606846720"},
											"actual value":   &entities.FactValueString{Value: "1152921504606846720"},
											"amendments": &entities.FactValueList{
												Value: []entities.FactValue{
													&entities.FactValueMap{
														Value: map[string]entities.FactValue{
															"index": &entities.FactValueInt{Value: 11},
															"amendment": &entities.FactValueString{Value: "[11] parameter is additional defined in sysctl config file /boot/sysctl.conf-5.3.18-150300.59.90-default(0x0fffffffffffff00), /boot/sysctl.conf-5.3.18-150300.59.93-default(0x0fffffffffffff00)"},
														},
													},
												},
											},											
										},
									},
									&entities.FactValueMap{
										Value: map[string]entities.FactValue{
											"Note ID":        &entities.FactValueString{Value: "941735"},
											"Note version":   &entities.FactValueString{Value: "11"},
											"parameter":      &entities.FactValueString{Value: "kernel.shmmax"},
											"compliant":      &entities.FactValueBool{Value: true},
											"expected value": &entities.FactValueString{Value: "18446744073709551615"},
											"actual value":   &entities.FactValueString{Value: "18446744073709551615"},
											"amendments": &entities.FactValueList{
												Value: []entities.FactValue{
													&entities.FactValueMap{
														Value: map[string]entities.FactValue{
															"index": &entities.FactValueInt{Value: 11},
															"amendment": &entities.FactValueString{Value: "[11] parameter is additional defined in sysctl config file /boot/sysctl.conf-5.3.18-150300.59.90-default(0x0fffffffffffff00), /boot/sysctl.conf-5.3.18-150300.59.93-default(0x0fffffffffffff00)\n [11] parameter is additional defined in sysctl config file /boot/sysctl.conf-5.3.18-150300.59.90-default(0xffffffffffffffff), /boot/sysctl.conf-5.3.18-150300.59.93-default(0xffffffffffffffff)"},
														},
													},
												},
											},											
										},
									},
									// ... represent other verifications similarly.
								},
							},
							// Other fields under "result" will come here.
							"attentions": &entities.FactValueList{
								Value: []entities.FactValue{},
							},
							"Notes enabled": &entities.FactValueList{
								Value: []entities.FactValue{
									&entities.FactValueString{Value: "941735"},
									&entities.FactValueString{Value: "1771258"},
									&entities.FactValueString{Value: "2578899"},
									&entities.FactValueString{Value: "2993054"},
									&entities.FactValueString{Value: "1656250"},
									&entities.FactValueString{Value: "900929"},
								},
							},
							// Other root-level fields will come here.
							"system compliance": &entities.FactValueBool{Value: false},
						},
					},
					"messages": &entities.FactValueList{
						Value: []entities.FactValue{
							&entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"priority": &entities.FactValueString{Value: "NOTICE"},
									"message": &entities.FactValueString{Value: "actions.go:85: ATTENTION: You are running a test version (3.1.0 from 2022/11/28) of saptune which is not supported for production use\n"},
								},
							},
							&entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"priority": &entities.FactValueString{Value: "WARNING"},
									"message": &entities.FactValueString{Value: "sysctl.go:73: Parameter 'kernel.shmmax' additional defined in the following sysctl config file /boot/sysctl.conf-5.3.18-150300.59.90-default(0xffffffffffffffff), /boot/sysctl.conf-5.3.18-150300.59.93-default(0xffffffffffffffff).\n"},
								},
							},
							&entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"priority": &entities.FactValueString{Value: "WARNING"},
									"message": &entities.FactValueString{Value: "sysctl.go:73: Parameter 'kernel.shmall' additional defined in the following sysctl config file /boot/sysctl.conf-5.3.18-150300.59.90-default(0x0fffffffffffff00), /boot/sysctl.conf-5.3.18-150300.59.93-default(0x0fffffffffffff00).\n"},
								},
							},
							&entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"priority": &entities.FactValueString{Value: "NOTICE"},
									"message": &entities.FactValueString{Value: "ini.go:308: block device related section settings detected: Traversing all block devices can take a considerable amount of time.\n"},
								},
							},
							&entities.FactValueMap{
								Value: map[string]entities.FactValue{
									"priority": &entities.FactValueString{Value: "ERROR"},
									"message": &entities.FactValueString{Value: "system.go:148: The parameters listed above have deviated from SAP/SUSE recommendations.\n\n"},
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
