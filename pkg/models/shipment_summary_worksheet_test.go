package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchShipmentSummaryWorksheetFormValues() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:  models.StringPointer("Marcus"),
			MiddleName: models.StringPointer("Joseph"),
			LastName:   models.StringPointer("Jenkins"),
			Suffix:     models.StringPointer("Jr."),
		},
	})
	sswPage1, _, err := models.FetchShipmentSummaryWorksheetFormValues(suite.DB(), move.ID)
	suite.NoError(err)

	suite.Equal("Jenkins Jr., Marcus Joseph", sswPage1.ServiceMemberName)
	suite.Equal("90 days per each shipment", sswPage1.MaxSITStorageEntitlement)
}
