package models_test

import (
	"fmt"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"time"
)

func (suite *ModelSuite) TestFetchShipmentSummaryWorksheetFormValues() {
	moveID, _ := uuid.NewV4()
	firstName := "Marcus"
	middleName := "Joseph"
	lastName := "Jenkins"
	suffix := "Jr."
	fullName := fmt.Sprintf("%s %s, %s %s", lastName, suffix, firstName, middleName)
	preferredPhoneNumber := "444-555-8888"
	preferredEmail := "michael+ppm-expansion_1@truss.works"
	maxSITStorageEntitlementDefault := "90 days per each shipment"
	dodID := "1234567890"
	serviceBranch := models.AffiliationAIRFORCE
	rank := models.ServiceMemberRankE9
	wtgEntitlements := models.WeightAllotment{
		TotalWeightSelf:     13000,
		ProGearWeight:       2000,
		ProGearWeightSpouse: 500,
	}
	totalWeightEntitlement := 15500
	ordersDate := time.Date(2018, time.December, 21, 0, 0, 0, 0, time.UTC)
	print(ordersDate.Location())
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	ordersNumber := "012345"
	issuingBranch := string(serviceBranch)

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

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            moveID,
			FirstName:     models.StringPointer(firstName),
			MiddleName:    models.StringPointer(middleName),
			LastName:      models.StringPointer(lastName),
			Suffix:        models.StringPointer(suffix),
			Telephone:     models.StringPointer(preferredPhoneNumber),
			PersonalEmail: models.StringPointer(preferredEmail),
			Edipi:         models.StringPointer(dodID),
			Affiliation:   &serviceBranch,
			Rank:          &rank,
			DutyStationID: &fortBragg.ID,
		},
		Order: models.Order{
			IssueDate:           ordersDate,
			OrdersType:          ordersType,
			OrdersNumber:        models.StringPointer(ordersNumber),
			NewDutyStationID:    fortBenning.ID,
			OrdersIssuingAgency: &issuingBranch,
		},
	})
	sswPage1, _, err := models.FetchShipmentSummaryWorksheetFormValues(suite.DB(), move.ID)

	suite.NoError(err)
	suite.Equal(fullName, sswPage1.ServiceMemberName)
	suite.Equal(maxSITStorageEntitlementDefault, sswPage1.MaxSITStorageEntitlement)
	suite.Equal(preferredPhoneNumber, sswPage1.PreferredPhone)
	suite.Equal(preferredEmail, sswPage1.PreferredEmail)
	suite.Equal(dodID, sswPage1.DODId)
	suite.Equal(string(serviceBranch), sswPage1.ServiceBranch)
	suite.Equal(string(rank), sswPage1.Rank)

	suite.Equal(ordersNumber, sswPage1.OrdersNumber)
	suite.Equal(issuingBranch, sswPage1.IssuingAgency)
	suite.True(ordersDate.Equal(sswPage1.OrderIssueDate))
	suite.Equal(ordersType, sswPage1.OrdersType)

	suite.Equal(fortBragg.ID, sswPage1.AuthorizedOrigin.ID)
	suite.Equal(fortBragg.Address.State, sswPage1.AuthorizedOrigin.Address.State)
	suite.Equal(fortBragg.Address.City, sswPage1.AuthorizedOrigin.Address.City)
	suite.Equal(fortBragg.Address.PostalCode, sswPage1.AuthorizedOrigin.Address.PostalCode)

	suite.Equal(fortBenning.ID, sswPage1.AuthorizedDestination.ID)
	suite.Equal(fortBenning.Address.State, sswPage1.AuthorizedDestination.Address.State)
	suite.Equal(fortBenning.Address.City, sswPage1.AuthorizedDestination.Address.City)
	suite.Equal(fortBenning.Address.PostalCode, sswPage1.AuthorizedDestination.Address.PostalCode)

	suite.Equal(wtgEntitlements.TotalWeightSelf, sswPage1.WeightAllotment.TotalWeightSelf)
	suite.Equal(wtgEntitlements.ProGearWeight, sswPage1.WeightAllotment.ProGearWeight)
	suite.Equal(wtgEntitlements.ProGearWeightSpouse, sswPage1.WeightAllotment.ProGearWeightSpouse)
	suite.Equal(totalWeightEntitlement, sswPage1.TotalWeightAllotment)
}
