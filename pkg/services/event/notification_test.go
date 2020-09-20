package event

import (
	"fmt"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EventServiceSuite) Test_MTOServiceItemPayload() {
	now := time.Now()

	mtoServiceItemDOFSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
		ReService: models.ReService{
			Code: "DOFSIT",
		},
	})
	// #TODO: Customer Contacts are not be created in testdatagen?
	// mtoServiceItemDDFSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
	// 	Move: models.Move{
	// 		AvailableToPrimeAt: &now,
	// 	},
	// 	ReService: models.ReService{
	// 		Code: "DDFSIT",
	// 	},
	// })

	// #TODO: Description is not being created in testdatagen
	// mtoServiceItemDOSHUT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
	// 	Move: models.Move{
	// 		AvailableToPrimeAt: &now,
	// 	},
	// 	ReService: models.ReService{
	// 		Code: "DOSHUT",
	// 	},
	// })

	// #TODO: Dimensions are not being created in testdatagen
	// mtoServiceItemDCRT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
	// 	Move: models.Move{
	// 		AvailableToPrimeAt: &now,
	// 	},
	// 	ReService: models.ReService{
	// 		Code: "DCRT",
	// 	},
	// })

	suite.T().Run("Success with MTOServiceItemDOFSIT", func(t *testing.T) {
		data := &primemessages.MTOServiceItemDOFSIT{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.DB(), mtoServiceItemDOFSIT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDOFSIT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDOFSIT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(fmt.Sprintf("%s", mtoServiceItemDOFSIT.ReService.Code), *data.ReServiceCode)
		suite.Equal(mtoServiceItemDOFSIT.Reason, data.Reason)
	})

}
