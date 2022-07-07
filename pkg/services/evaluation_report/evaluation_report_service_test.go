package evaluationreport

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type EvaluationReportSuite struct {
	testingsuite.PopTestSuite
}

func TestEvaluationReportServiceSuite(t *testing.T) {
	ts := &EvaluationReportSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
