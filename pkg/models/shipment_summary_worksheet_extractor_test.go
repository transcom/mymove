package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchShipmentSummaryWorksheetExtractor() {
	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{},
		ServiceMember: models.ServiceMember{
			FirstName:  models.StringPointer("Marcus"),
			MiddleName: models.StringPointer("Joseph"),
			LastName:   models.StringPointer("Jenkins"),
			Suffix:     models.StringPointer("Jr."),
		},
	})
	sswe, err := models.FetchShipmentSummaryWorksheetExtractor(suite.db, shipment.ID)
	suite.NoError(err)

	suite.Equal("Jenkins Jr., Marcus Joseph", sswe.ServiceMemberName)
}
