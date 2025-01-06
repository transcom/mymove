package report

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ReportServiceSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *ReportServiceSuite) SetupSuite() {
	suite.PreloadData(func() {
		factory.SetupDefaultAllotments(suite.DB())
	})
}

func TestReportServiceSuite(t *testing.T) {
	ts := &ReportServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
