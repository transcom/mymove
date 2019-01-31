package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchDataShipmentSummaryWorksFormData() {
	moveID, _ := uuid.NewV4()
	serviceMemberID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	rank := models.ServiceMemberRankE9

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID: moveID,
		},
		Order: models.Order{
			OrdersType:       ordersType,
			NewDutyStationID: fortGordon.ID,
		},
		ServiceMember: models.ServiceMember{
			ID:            serviceMemberID,
			DutyStationID: &yuma.ID,
			Rank:          &rank,
		},
	})
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			ServiceMemberID: serviceMemberID,
			MoveID:          moveID,
		},
	})
	session := auth.Session{
		UserID:          move.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MyApp,
	}
	ssd, err := models.FetchDataShipmentSummaryWorksheetFormData(suite.DB(), &session, moveID)

	suite.NoError(err)
	suite.Equal(move.Orders.ID, ssd.Order.ID)
	suite.Require().Len(ssd.Shipments, 1)
	suite.Equal(shipment.ID, ssd.Shipments[0].ID)
	suite.Equal(serviceMemberID, ssd.ServiceMember.ID)
	suite.Equal(yuma.ID, ssd.CurrentDutyStation.ID)
	suite.Equal(yuma.Address.ID, ssd.CurrentDutyStation.Address.ID)
	suite.Equal(fortGordon.ID, ssd.NewDutyStation.ID)
	suite.Equal(fortGordon.Address.ID, ssd.NewDutyStation.Address.ID)
	rankWtgAllotment := models.GetWeightAllotment(rank)
	suite.Equal(rankWtgAllotment, ssd.WeightAllotment)
}

func (suite *ModelSuite) TestFormatValuesShipmentSummaryWorksheetFormPage1() {
	yuma := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	wtgEntitlements := models.WeightAllotment{
		TotalWeightSelf:     13000,
		ProGearWeight:       2000,
		ProGearWeightSpouse: 500,
	}
	serviceMemberID, _ := uuid.NewV4()
	serviceBranch := models.AffiliationAIRFORCE
	rank := models.ServiceMemberRankE9
	serviceMember := models.ServiceMember{
		ID:            serviceMemberID,
		FirstName:     models.StringPointer("Marcus"),
		MiddleName:    models.StringPointer("Joseph"),
		LastName:      models.StringPointer("Jenkins"),
		Suffix:        models.StringPointer("Jr."),
		Telephone:     models.StringPointer("444-555-8888"),
		PersonalEmail: models.StringPointer("michael+ppm-expansion_1@truss.works"),
		Edipi:         models.StringPointer("1234567890"),
		Affiliation:   &serviceBranch,
		Rank:          &rank,
		DutyStationID: &yuma.ID,
	}

	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	order := models.Order{
		IssueDate:           orderIssueDate,
		OrdersType:          internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:        models.StringPointer("012345"),
		NewDutyStationID:    fortGordon.ID,
		OrdersIssuingAgency: models.StringPointer(string(serviceBranch)),
		TAC:                 models.StringPointer("NTA4"),
	}
	pickupDate := time.Date(2019, time.January, 11, 0, 0, 0, 0, time.UTC)
	weight := unit.Pound(5000)
	shipments := []models.Shipment{
		{
			ActualPickupDate: &pickupDate,
			NetWeight:        &weight,
			Status:           models.ShipmentStatusDELIVERED,
		},
	}

	ssd := models.ShipmentSummaryFormData{
		ServiceMember:      serviceMember,
		Order:              order,
		CurrentDutyStation: yuma,
		NewDutyStation:     fortGordon,
		WeightAllotment:    wtgEntitlements,
		Shipments:          shipments,
		PreparationDate:    time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
	}
	sswPage1 := models.FormatValuesShipmentSummaryWorksheetFormPage1(ssd)

	suite.Equal("01-Jan-2019", sswPage1.PreparationDate)

	suite.Equal("Jenkins Jr., Marcus Joseph", sswPage1.ServiceMemberName)
	suite.Equal("90 days per each shipment", sswPage1.MaxSITStorageEntitlement)
	suite.Equal("Yuma AFB, IA 50309", sswPage1.AuthorizedOrigin)
	suite.Equal("Fort Gordon, GA 30813", sswPage1.AuthorizedDestination)
	suite.Equal("NO", sswPage1.POVAuthorized)
	suite.Equal("444-555-8888", sswPage1.PreferredPhone)
	suite.Equal("michael+ppm-expansion_1@truss.works", sswPage1.PreferredEmail)
	suite.Equal("1234567890", sswPage1.DODId)

	suite.Equal("Air Force", sswPage1.IssuingBranchOrAgency)
	suite.Equal("21-Dec-2018", sswPage1.OrdersIssueDate)
	suite.Equal("PCS/012345", sswPage1.OrdersTypeAndOrdersNumber)
	suite.Equal("NTA4", sswPage1.TAC)

	suite.Equal("Fort Gordon, GA", sswPage1.NewDutyAssignment)

	suite.Equal("13,000", sswPage1.WeightAllotmentSelf)
	suite.Equal("2,000", sswPage1.WeightAllotmentProgear)
	suite.Equal("500", sswPage1.WeightAllotmentProgearSpouse)
	suite.Equal("15,500", sswPage1.TotalWeightAllotment)

	suite.Equal("01 - PPM", sswPage1.Shipment1NumberAndType)
	suite.Equal("11-Jan-2019", sswPage1.Shipment1PickUpDate)
	suite.Equal("5,000 lbs - FINAL", sswPage1.Shipment1Weight)
	suite.Equal("Delivered", sswPage1.Shipment1CurrentShipmentStatus)

}

func (suite *ModelSuite) FormatAuthorizedLocation() {
	fortGordon := models.DutyStation{Name: "Fort Gordon", Address: models.Address{State: "GA", PostalCode: "30813"}}
	yuma := models.DutyStation{Name: "Yuma AFB", Address: models.Address{State: "IA", PostalCode: "50309"}}

	suite.Equal("Fort Gordon, GA 30813", models.FormatDutyStation(fortGordon))
	suite.Equal("Yuma AFB, IA 50309", models.FormatDutyStation(yuma))
}

func (suite *ModelSuite) TestFormatServiceMemberFullName() {
	sm1 := models.ServiceMember{
		Suffix:     models.StringPointer("Jr."),
		FirstName:  models.StringPointer("Tom"),
		MiddleName: models.StringPointer("James"),
		LastName:   models.StringPointer("Smith"),
	}
	sm2 := models.ServiceMember{
		FirstName: models.StringPointer("Tom"),
		LastName:  models.StringPointer("Smith"),
	}

	suite.Equal("Smith Jr., Tom James", models.FormatServiceMemberFullName(sm1))
	suite.Equal("Smith, Tom", models.FormatServiceMemberFullName(sm2))
}

func (suite *ModelSuite) TestFormatCurrentShipmentStatus() {
	completed := models.Shipment{Status: models.ShipmentStatusDELIVERED}
	inTransit := models.Shipment{Status: models.ShipmentStatusINTRANSIT}

	suite.Equal("Delivered", models.FormatCurrentShipmentStatus(completed))
	suite.Equal("In Transit", models.FormatCurrentShipmentStatus(inTransit))
}

func (suite *ModelSuite) TestFormatShipmentNumberAndType() {
	singleShipment := models.Shipments{models.Shipment{}}
	multipleShipments := models.Shipments{models.Shipment{}, models.Shipment{}}

	multipleShipmentsFormatted := models.FormatShipments(multipleShipments)

	suite.Equal("01 - PPM", models.FormatShipments(singleShipment)[0].ShipmentNumberAndType)
	suite.Require().Len(multipleShipmentsFormatted, 2)
	suite.Equal("01 - PPM", multipleShipmentsFormatted[0].ShipmentNumberAndType)
	suite.Equal("02 - PPM", multipleShipmentsFormatted[1].ShipmentNumberAndType)

}

func (suite *ModelSuite) TestFormatShipmentWeight() {
	pounds := unit.Pound(1000)
	shipment := models.Shipment{NetWeight: &pounds}

	suite.Equal("1,000 lbs - FINAL", models.FormatShipmentWeight(shipment))
}

func (suite *ModelSuite) TestFormatPickupDate() {
	pickupDate := time.Date(2018, time.December, 1, 0, 0, 0, 0, time.UTC)
	shipment := models.Shipment{ActualPickupDate: &pickupDate}

	suite.Equal("01-Dec-2018", models.FormatShipmentPickupDate(shipment))
}

func (suite *ModelSuite) TestFormatWeights() {
	suite.Equal("0", models.FormatWeights(0))
	suite.Equal("10", models.FormatWeights(10))
	suite.Equal("1,000", models.FormatWeights(1000))
	suite.Equal("1,000,000", models.FormatWeights(1000000))
}

func (suite *ModelSuite) TestFormatDutyStation() {
	fortBenning := models.DutyStation{Name: "Fort Benning", Address: models.Address{State: "GA"}}
	yuma := models.DutyStation{Name: "Yuma AFB", Address: models.Address{State: "AZ"}}

	suite.Equal("Fort Benning, GA", models.FormatDutyStation(fortBenning))
	suite.Equal("Yuma AFB, AZ", models.FormatDutyStation(yuma))
}

func (suite *ModelSuite) TestFormatOrdersIssueDate() {
	dec212018 := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	jan012019 := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	suite.Equal("21-Dec-2018", models.FormatDate(dec212018))
	suite.Equal("01-Jan-2019", models.FormatDate(jan012019))
}

func (suite *ModelSuite) TestFormatOrdersType() {
	pcsOrder := models.Order{OrdersType: internalmessages.OrdersTypePERMANENTCHANGEOFSTATION}
	var unknownOrdersType internalmessages.OrdersType = "UNKNOWN_ORDERS_TYPE"
	localMoveOrder := models.Order{OrdersType: unknownOrdersType}

	suite.Equal("PCS", models.FormatOrdersType(pcsOrder))
	suite.Equal("", models.FormatOrdersType(localMoveOrder))
}

func (suite *ModelSuite) TestFormatServiceMemberAffiliation() {
	airForce := models.AffiliationAIRFORCE
	marines := models.AffiliationMARINES

	suite.Equal("Air Force", models.FormatServiceMemberAffiliation(&airForce))
	suite.Equal("Marines", models.FormatServiceMemberAffiliation(&marines))
}
