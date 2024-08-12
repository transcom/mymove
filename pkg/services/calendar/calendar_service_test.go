package calendar

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type CalendarSuite struct {
	*testingsuite.PopTestSuite
}

func TestCalendarServiceSuite(t *testing.T) {
	ts := &CalendarSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
