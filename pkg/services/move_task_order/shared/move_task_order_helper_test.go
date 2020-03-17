package movetaskordershared_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	movetaskordershared "github.com/transcom/mymove/pkg/services/move_task_order/shared"
)

func (suite *MoveTaskOrderHelperSuite) TestMoveTaskOrderGenerateReferenceID() {

	refID, err := movetaskordershared.GenerateReferenceID(suite.DB())
	suite.T().Run("reference id is properly created", func(t *testing.T) {
		// testing reference id
		suite.NoError(err)
		suite.NotZero(refID)
		firstNum, _ := strconv.Atoi(strings.Split(refID, "-")[0])
		secondNum, _ := strconv.Atoi(strings.Split(refID, "-")[1])
		suite.Equal(reflect.TypeOf(refID).String(), "string")
		suite.Equal(firstNum >= 0 && firstNum <= 9999, true)
		suite.Equal(secondNum >= 0 && secondNum <= 9999, true)
		suite.Equal(string(refID[4]), "-")
	})
}
