package movetaskorder

import (
	"log"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	mto := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	log.Println(mto)

}
