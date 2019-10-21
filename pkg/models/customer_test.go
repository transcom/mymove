package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestGetCustomerMoveItems() {
	affiliation := models.AffiliationARMY
	testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:     models.StringPointer("Test"),
			LastName:      models.StringPointer("User"),
			Affiliation: &affiliation,
		},
		Move: models.Move{
			Locator: "DFTMVE",
		},
	})

	testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:     models.StringPointer("Test"),
			LastName:      models.StringPointer("User 2"),
			Affiliation: &affiliation,
		},
		Move: models.Move{
			Locator: "TES123",
		},
	})

	customerMoveItems, err := models.GetCustomerMoveItems(suite.DB())
	suite.NoError(err)
	suite.Len(customerMoveItems, 2)
	suite.Equal(customerMoveItems[1].CustomerName, "User 2, Test")
	suite.Equal(customerMoveItems[1].ConfirmationNumber, "TES123")
	suite.Equal(customerMoveItems[0].CustomerName, "User, Test")
	suite.Equal(customerMoveItems[0].ConfirmationNumber, "DFTMVE")
}
