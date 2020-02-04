package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestConvert() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	suite.NotNil(move)

	moID, conversionErr := models.ConvertFromPPMToGHC(suite.DB(), move.ID)
	suite.NoError(conversionErr)

	var mo models.MoveOrder
	suite.NoError(suite.DB().Eager("Customer").Find(&mo, moID))

	suite.NotEqual(uuid.Nil, mo.DestinationDutyStationID)
	suite.Equal(&move.Orders.NewDutyStationID, mo.DestinationDutyStationID)

	suite.NotEqual(uuid.Nil, mo.OriginDutyStationID)
	suite.Equal(move.Orders.ServiceMember.DutyStationID, mo.OriginDutyStationID)

	suite.NotEqual(uuid.Nil, mo.EntitlementID)

	customer := mo.Customer
	suite.Equal(*move.Orders.ServiceMember.Edipi, customer.DODID)
	suite.NotEqual(uuid.Nil, customer.UserID)
	suite.Equal(move.Orders.ServiceMember.UserID, customer.UserID)

	suite.NotEqual(uuid.Nil, customer.UserID)
	suite.Equal(&customer.ID, mo.CustomerID)
}
