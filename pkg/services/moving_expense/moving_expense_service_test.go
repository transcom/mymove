package movingexpense

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MovingExpenseSuite struct {
	*testingsuite.PopTestSuite
}

func TestMovingExpenseServiceSuite(t *testing.T) {
	ts := &MovingExpenseSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
