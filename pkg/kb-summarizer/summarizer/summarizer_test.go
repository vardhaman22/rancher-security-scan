package summarizer

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testAvNodeMap1 = map[string]string{
		"node1": "permission=644",
		"node2": "permission=640",
		"node3": "permission=600",
	}

	testAvNodeMap2 = map[string]string{
		"node1": "testvalue",
		"node2": "testvalue1",
		"node3": "testvalue2",
	}

	testAvNodeMap3 = map[string]string{
		"node1": "true",
		"node2": "false",
		"node3": "true",
	}

	groupWrappersTestData = []*GroupWrapper{
		&GroupWrapper{
			ID:   "1.1",
			Text: "Checks for group 1.1",
			CheckWrappers: []*CheckWrapper{
				&CheckWrapper{
					ID:                 "1.1.1",
					Text:               "Check 1.1.1",
					ActualValueNodeMap: testAvNodeMap1,
				},
				&CheckWrapper{
					ID:                 "1.1.2",
					Text:               "Check 1.1.2",
					ActualValueNodeMap: testAvNodeMap2,
				},
			},
		},
		&GroupWrapper{
			ID:   "1.2",
			Text: "Checks for group 1.2",
			CheckWrappers: []*CheckWrapper{
				&CheckWrapper{
					ID:                 "1.2.1",
					Text:               "Check 1.2.1",
					ActualValueNodeMap: testAvNodeMap2,
				},
			},
		},
		&GroupWrapper{
			ID:   "2.1",
			Text: "Checks for group 2.1",
			CheckWrappers: []*CheckWrapper{
				&CheckWrapper{
					ID:                 "2.1.1",
					Text:               "Check 2.1.1",
					ActualValueNodeMap: testAvNodeMap2,
				},
			},
		},
		&GroupWrapper{
			ID:   "3.1",
			Text: "Checks for group 2.2",
			CheckWrappers: []*CheckWrapper{
				&CheckWrapper{
					ID:                 "3.2",
					Text:               "Check 3.1",
					ActualValueNodeMap: testAvNodeMap3,
				},
			},
		},
	}

	avGroupsTestData = []*ActualValueGroup{
		&ActualValueGroup{
			ID:   "1.1",
			Text: "Checks for group 1.1",
			ActualValueChecks: []*ActualValueCheck{
				&ActualValueCheck{
					ID:                 "1.1.1",
					Text:               "Check 1.1.1",
					ActualValueNodeMap: testAvNodeMap1,
				},
				&ActualValueCheck{
					ID:                 "1.1.2",
					Text:               "Check 1.1.2",
					ActualValueNodeMap: testAvNodeMap2,
				},
			},
		},
		&ActualValueGroup{
			ID:   "1.2",
			Text: "Checks for group 1.2",
			ActualValueChecks: []*ActualValueCheck{
				&ActualValueCheck{
					ID:                 "1.2.1",
					Text:               "Check 1.2.1",
					ActualValueNodeMap: testAvNodeMap2,
				},
			},
		},
		&ActualValueGroup{
			ID:   "2.1",
			Text: "Checks for group 2.1",
			ActualValueChecks: []*ActualValueCheck{
				&ActualValueCheck{
					ID:                 "2.1.1",
					Text:               "Check 2.1.1",
					ActualValueNodeMap: testAvNodeMap2,
				},
			},
		},
		&ActualValueGroup{
			ID:   "3.1",
			Text: "Checks for group 2.2",
			ActualValueChecks: []*ActualValueCheck{
				&ActualValueCheck{
					ID:                 "3.2",
					Text:               "Check 3.1",
					ActualValueNodeMap: testAvNodeMap3,
				},
			},
		},
	}
)

func TestSummarizer_SetfullReportActualValueMapData(t *testing.T) {

	s := Summarizer{
		fullReport: &SummarizedReport{
			GroupWrappers: groupWrappersTestData,
		},
	}

	err := s.setfullReportActualValueMapData()
	require.Nil(t, err, "failed to ActualValueMapData to fullReport")

	require.NotEmpty(t, s.fullReport.ActualValueMapData, "empty actualValueMapData found")

	compressedAvMapData, err := base64.StdEncoding.DecodeString(s.fullReport.ActualValueMapData)
	require.Nil(t, err, "error while decoding ActualValueMapData")
	require.NotNil(t, compressedAvMapData, "compressed ActualValueMapData should not be empty")

	r, err := gzip.NewReader(bytes.NewBuffer(compressedAvMapData))
	require.Nil(t, err, "error while reading compressed avMapData")
	defer r.Close()

	avgroupsData, err := io.ReadAll(r)
	require.Nil(t, err, "error while reading avMapData json")

	expectedAvGroupsData, err := json.Marshal(avGroupsTestData)
	require.Nil(t, err, "error while encoding expected avgroups data")

	assert.EqualValues(t, string(expectedAvGroupsData), string(avgroupsData), "avmapData is not correctly encoded")
}

func TestMapGroupWrappersToActualValueGroups(t *testing.T) {

	avgroupsData := mapGroupWrappersToActualValueGroups(groupWrappersTestData)
	expectedAvGroups := avGroupsTestData

	assert.EqualValues(t, expectedAvGroups, avgroupsData, "GroupWrappers are not correctly mapped to ActualValueGroups")
}
