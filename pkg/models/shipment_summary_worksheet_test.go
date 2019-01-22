package models_test

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) setupTestDutyStations() (currentDutyStation models.DutyStation, newDutyStation models.DutyStation) {
	fortBraggAssertions := testdatagen.Assertions{
		Address: models.Address{
			City:       "Fort Bragg",
			State:      "NC",
			PostalCode: "28310",
		},
		DutyStation: models.DutyStation{
			Name: "Fort Bragg",
		},
	}
	fortBragg := testdatagen.MakeDutyStation(suite.DB(), fortBraggAssertions)

	fortBenningAssertions := testdatagen.Assertions{
		Address: models.Address{
			City:       "Fort Benning",
			State:      "GA",
			PostalCode: "31905",
		},
		DutyStation: models.DutyStation{
			Name: "Fort Benning",
		},
	}
	fortBenning := testdatagen.MakeDutyStation(suite.DB(), fortBenningAssertions)

	return fortBragg, fortBenning
}

func (suite *ModelSuite) TestFetchDataShipmentSummaryWorksFormData() {
	moveID, _ := uuid.NewV4()
	serviceMemberID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	fortBragg, fortBenning := suite.setupTestDutyStations()
	rank := models.ServiceMemberRankE9
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID: moveID,
		},
		Order: models.Order{
			OrdersType:       ordersType,
			NewDutyStationID: fortBenning.ID,
		},
		ServiceMember: models.ServiceMember{
			ID:            serviceMemberID,
			DutyStationID: &fortBragg.ID,
			Rank:          &rank,
		},
	})
	ssd, err := models.FetchDataShipmentSummaryWorksFormData(suite.DB(), move.ID)

	suite.NoError(err)
	suite.Equal(move.Orders.ID, ssd.Order.ID)
	suite.Equal(serviceMemberID, ssd.ServiceMember.ID)
	suite.Equal(fortBragg.ID, ssd.CurrentDutyStation.ID)
	suite.Equal(fortBenning.ID, ssd.NewDutyStation.ID)
	rankWtgAllotment := models.GetWeightAllotment(rank)
	suite.Equal(rankWtgAllotment, ssd.WeightAllotment)
}

func (suite *ModelSuite) TestFormatValuesShipmentSummaryWorksheetFormPage1() {
	fortBragg, fortBenning := suite.setupTestDutyStations()
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
		DutyStationID: &fortBragg.ID,
	}
	orderIssueDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	order := models.Order{
		IssueDate:           orderIssueDate,
		OrdersType:          internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:        models.StringPointer("012345"),
		NewDutyStationID:    fortBenning.ID,
		OrdersIssuingAgency: models.StringPointer(string(serviceBranch)),
	}
	ssd := models.ShipmentSummaryFormData{
		ServiceMember:      serviceMember,
		Order:              order,
		CurrentDutyStation: fortBragg,
		NewDutyStation:     fortBenning,
		WeightAllotment:    wtgEntitlements,
	}
	sswPage1 := models.FormatValuesShipmentSummaryWorksheetFormPage1(ssd)

	suite.Equal("Jenkins Jr., Marcus Joseph", sswPage1.ServiceMemberName)
	suite.Equal("90 days per each shipment", sswPage1.MaxSITStorageEntitlement)
	suite.Equal("444-555-8888", sswPage1.PreferredPhone)
	suite.Equal("michael+ppm-expansion_1@truss.works", sswPage1.PreferredEmail)
	suite.Equal("1234567890", sswPage1.DODId)
	suite.Equal(string(serviceBranch), sswPage1.ServiceBranch)
	suite.Equal(string(rank), sswPage1.Rank)

	suite.Equal("012345", sswPage1.OrdersNumber)
	suite.Equal(string(serviceBranch), sswPage1.IssuingAgency)
	suite.True(orderIssueDate.Equal(sswPage1.OrderIssueDate))
	suite.Equal(internalmessages.OrdersTypePERMANENTCHANGEOFSTATION, sswPage1.OrdersType)

	suite.Equal(fortBragg.ID, sswPage1.AuthorizedOrigin.ID)
	suite.Equal(fortBragg.Address.State, sswPage1.AuthorizedOrigin.Address.State)
	suite.Equal(fortBragg.Address.City, sswPage1.AuthorizedOrigin.Address.City)
	suite.Equal(fortBragg.Address.PostalCode, sswPage1.AuthorizedOrigin.Address.PostalCode)

	suite.Equal(fortBenning.ID, sswPage1.AuthorizedDestination.ID)
	suite.Equal(fortBenning.Address.State, sswPage1.AuthorizedDestination.Address.State)
	suite.Equal(fortBenning.Address.City, sswPage1.AuthorizedDestination.Address.City)
	suite.Equal(fortBenning.Address.PostalCode, sswPage1.AuthorizedDestination.Address.PostalCode)

	suite.Equal("13,000", sswPage1.WeightAllotmentSelf)
	suite.Equal("2,000", sswPage1.WeightAllotmentProgear)
	suite.Equal("500", sswPage1.WeightAllotmentProgearSpouse)
	suite.Equal("15,500", sswPage1.TotalWeightAllotment)
}

func (suite *ModelSuite) TestFormatWeights() {
	suite.Equal("0", models.FormatWeights(0))
	suite.Equal("10", models.FormatWeights(10))
	suite.Equal("1,000", models.FormatWeights(1000))
	suite.Equal("1,000,000", models.FormatWeights(1000000))
}
