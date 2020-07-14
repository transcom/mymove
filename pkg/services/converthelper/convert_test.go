package converthelper_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/services/converthelper"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ConvertSuite) TestConvertFromPPMToGHC() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			HasDependents: false,
		},
	})
	moveOrder := move.Orders
	sm := moveOrder.ServiceMember
	hasDependents := moveOrder.HasDependents
	spouseHasProGear := moveOrder.SpouseHasProGear
	weight, _ := models.GetEntitlement(*sm.Rank, hasDependents, spouseHasProGear)
	entitlement := models.Entitlement{
		DependentsAuthorized: &hasDependents,
		DBAuthorizedWeight:   models.IntPointer(weight),
	}
	suite.MustSave(&entitlement)
	moveOrder.Entitlement = &entitlement
	moveOrder.EntitlementID = &entitlement.ID
	moveOrder.OriginDutyStation = &sm.DutyStation
	moveOrder.OriginDutyStationID = &sm.DutyStation.ID
	orderTypeDetail := internalmessages.OrdersTypeDetailHHGPERMITTED
	moveOrder.OrdersTypeDetail = &orderTypeDetail
	moveOrder.Grade = (*string)(sm.Rank)
	suite.MustSave(&moveOrder)

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

	moID, conversionErr := converthelper.ConvertFromPPMToGHC(suite.DB(), move.ID)
	suite.FatalNoError(conversionErr)

	var mo models.Order
	suite.FatalNoError(suite.DB().Eager("ServiceMember", "Entitlement").Find(&mo, moID))

	suite.NotNil(mo.ReportByDate)
	suite.NotNil(mo.IssueDate)
	suite.NotNil(mo.OrdersType)
	suite.NotNil(mo.OrdersTypeDetail)
	suite.NotNil(mo.Grade)

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

	var mto models.MoveTaskOrder
	suite.FatalNoError(suite.DB().Eager().Where("move_order_id = ?", mo.ID).First(&mto))

	var mtoShipmentHHG models.MTOShipment
	suite.FatalNoError(suite.DB().Eager().Where("move_task_order_id = ? and shipment_type = ?", mto.ID, models.MTOShipmentTypeHHGLongHaulDom).First(&mtoShipmentHHG))

	suite.NotNil(mtoShipmentHHG.ScheduledPickupDate)
	suite.Equal(unit.Pound(4096), *mtoShipmentHHG.PrimeEstimatedWeight)

	expectedNilTime := time.Time{}
	suite.NotEqual(expectedNilTime, *mtoShipmentHHG.ScheduledPickupDate)
	suite.NotNil(*mtoShipmentHHG.PrimeEstimatedWeightRecordedDate)
	suite.NotEqual(expectedNilTime, *mtoShipmentHHG.PrimeEstimatedWeightRecordedDate)

	var mtoShipmentHHGDomShortHaul models.MTOShipment
	suite.FatalNoError(suite.DB().Eager().Where("move_task_order_id = ? and shipment_type = ?", mto.ID, models.MTOShipmentTypeHHGShortHaulDom).First(&mtoShipmentHHGDomShortHaul))

	suite.NotNil(mtoShipmentHHGDomShortHaul.ScheduledPickupDate)
	suite.Equal(unit.Pound(4096), *mtoShipmentHHGDomShortHaul.PrimeEstimatedWeight)
	suite.NotNil(*mtoShipmentHHGDomShortHaul.PrimeEstimatedWeightRecordedDate)
	suite.NotEqual(expectedNilTime, *mtoShipmentHHGDomShortHaul.PrimeEstimatedWeightRecordedDate)

	suite.NotEqual(expectedNilTime, *mtoShipmentHHGDomShortHaul.ScheduledPickupDate)

	suite.Equal(sanDiego.Address.ID, *mtoShipmentHHGDomShortHaul.PickupAddressID)
	suite.Equal(miramar.Address.ID, *mtoShipmentHHGDomShortHaul.DestinationAddressID)
}

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
	suite.NotNil(mo.OrdersTypeDetail)
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

	var mto models.MoveTaskOrder
	suite.FatalNoError(suite.DB().Eager().Where("move_order_id = ?", mo.ID).First(&mto))
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
	_, conversionErr = converthelper.ConvertFromPPMToGHC(suite.DB(), move.ID)
	suite.FatalNoError(conversionErr)
	var moveOrders []models.Order

	err := suite.DB().Where("service_member_id = $1", sm.ID).All(&moveOrders)
	suite.FatalNoError(err)
	suite.Equal(2, len(moveOrders))
}
