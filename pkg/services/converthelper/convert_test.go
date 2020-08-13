package converthelper_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/converthelper"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ConvertSuite) TestConvertProfileOrdersToGHC() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	suite.NotNil(move)

	contractor := testdatagen.MakeContractor(suite.DB(), testdatagen.Assertions{
		Contractor: models.Contractor{
			ContractNumber: "HTC111-11-1-1111",
		},
	})
	suite.NotNil(contractor)

	miramar := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name: "USMC Miramar",
		},
	})
	suite.NotNil(miramar)

	sanDiego := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name: "USMC San Diego",
		},
	})
	suite.NotNil(sanDiego)

	moID, conversionErr := converthelper.ConvertProfileOrdersToGHC(suite.DB(), move.ID)
	suite.FatalNoError(conversionErr)

	var mo models.Order
	suite.FatalNoError(suite.DB().Eager("ServiceMember", "Entitlement").Find(&mo, moID))

	suite.NotNil(mo.ReportByDate)
	suite.NotNil(mo.IssueDate)
	suite.NotNil(mo.OrdersType)
	suite.Equal(mo.Grade, (*string)(mo.ServiceMember.Rank))

	suite.NotEqual(uuid.Nil, mo.NewDutyStationID)
	suite.Equal(move.Orders.NewDutyStationID, mo.NewDutyStationID)

	suite.NotEqual(uuid.Nil, mo.OriginDutyStationID)
	suite.Equal(move.Orders.ServiceMember.DutyStationID, mo.OriginDutyStationID)

	suite.NotEqual(uuid.Nil, mo.EntitlementID)
	suite.Equal(false, *mo.Entitlement.DependentsAuthorized)
	suite.Equal(7000, *mo.Entitlement.DBAuthorizedWeight)

	customer := mo.ServiceMember
	suite.Equal(*move.Orders.ServiceMember.Edipi, *customer.Edipi)
	suite.NotEqual(uuid.Nil, customer.UserID)
	suite.Equal(move.Orders.ServiceMember.UserID, customer.UserID)

	suite.NotEqual(uuid.Nil, customer.UserID)
	suite.Equal(customer.ID, mo.ServiceMemberID)

	var mto models.Move
	suite.FatalNoError(suite.DB().Eager().Where("orders_id = ?", mo.ID).First(&mto))
}

func (suite *ConvertSuite) TestConvertFromPPMToGHCMoveOrdersExist() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	suite.NotNil(move)

	sm := move.Orders.ServiceMember
	contractor := testdatagen.MakeContractor(suite.DB(), testdatagen.Assertions{
		Contractor: models.Contractor{
			ContractNumber: "HTC111-11-1-1111",
		},
	})
	suite.NotNil(contractor)

	miramar := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name: "USMC Miramar",
		},
	})
	suite.NotNil(miramar)

	sanDiego := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name: "USMC San Diego",
		},
	})
	suite.NotNil(sanDiego)

	_, conversionErr := converthelper.ConvertProfileOrdersToGHC(suite.DB(), move.ID)
	suite.FatalNoError(conversionErr)

	var orders []models.Order

	err := suite.DB().Where("service_member_id = $1", sm.ID).All(&orders)
	suite.FatalNoError(err)
	suite.Equal(1, len(orders))

	var moveTaskOrders []models.Move
	err = suite.DB().Where("orders_id = $1", orders[0].ID).All(&moveTaskOrders)
	suite.FatalNoError(err)
	suite.Equal(1, len(moveTaskOrders))
}
