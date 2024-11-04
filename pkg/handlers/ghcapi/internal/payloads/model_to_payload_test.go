package payloads

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/test"
)

func TestOrder(_ *testing.T) {
	order := &models.Order{}
	Order(order)
}

// TestMove makes sure zero values/optional fields are handled
func TestMove(t *testing.T) {
	_, err := Move(&models.Move{}, &test.FakeS3Storage{})
	if err != nil {
		t.Fail()
	}
}

func (suite *PayloadsSuite) TestFetchPPMShipment() {

	ppmShipmentID, _ := uuid.NewV4()
	streetAddress1 := "MacDill AFB"
	streetAddress2, streetAddress3 := "", ""
	city := "Tampa"
	state := "FL"
	postalcode := "33621"
	county := "HILLSBOROUGH"

	country := models.Country{
		Country: "US",
	}

	expectedAddress := models.Address{
		StreetAddress1: streetAddress1,
		StreetAddress2: &streetAddress2,
		StreetAddress3: &streetAddress3,
		City:           city,
		State:          state,
		PostalCode:     postalcode,
		Country:        &country,
		County:         county,
	}

	isActualExpenseReimbursement := true

	expectedPPMShipment := models.PPMShipment{
		ID:                           ppmShipmentID,
		PickupAddress:                &expectedAddress,
		DestinationAddress:           &expectedAddress,
		IsActualExpenseReimbursement: &isActualExpenseReimbursement,
	}

	suite.Run("Success -", func() {
		returnedPPMShipment := PPMShipment(nil, &expectedPPMShipment)

		suite.IsType(returnedPPMShipment, &ghcmessages.PPMShipment{})
		suite.Equal(&streetAddress1, returnedPPMShipment.PickupAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress2, returnedPPMShipment.PickupAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress3, returnedPPMShipment.PickupAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.PickupAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.PickupAddress.City)
		suite.Equal(&state, returnedPPMShipment.PickupAddress.State)
		suite.Equal(&country.Country, returnedPPMShipment.PickupAddress.Country)
		suite.Equal(&county, returnedPPMShipment.PickupAddress.County)

		suite.Equal(&streetAddress1, returnedPPMShipment.DestinationAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress2, returnedPPMShipment.DestinationAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress3, returnedPPMShipment.DestinationAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.DestinationAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.DestinationAddress.City)
		suite.Equal(&state, returnedPPMShipment.DestinationAddress.State)
		suite.Equal(&country.Country, returnedPPMShipment.DestinationAddress.Country)
		suite.Equal(&county, returnedPPMShipment.DestinationAddress.County)
		suite.True(*returnedPPMShipment.IsActualExpenseReimbursement)
	})

	suite.Run("Destination street address 1 returns empty string to convey OPTIONAL state ", func() {
		expected_street_address_1 := ""
		expectedAddress2 := models.Address{
			StreetAddress1: expected_street_address_1,
			StreetAddress2: &streetAddress2,
			StreetAddress3: &streetAddress3,
			City:           city,
			State:          state,
			PostalCode:     postalcode,
			Country:        &country,
			County:         county,
		}

		expectedPPMShipment2 := models.PPMShipment{
			ID:                 ppmShipmentID,
			PickupAddress:      &expectedAddress,
			DestinationAddress: &expectedAddress2,
		}
		returnedPPMShipment := PPMShipment(nil, &expectedPPMShipment2)

		suite.IsType(returnedPPMShipment, &ghcmessages.PPMShipment{})
		suite.Equal(&streetAddress1, returnedPPMShipment.PickupAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress2, returnedPPMShipment.PickupAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress3, returnedPPMShipment.PickupAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.PickupAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.PickupAddress.City)
		suite.Equal(&state, returnedPPMShipment.PickupAddress.State)
		suite.Equal(&county, returnedPPMShipment.PickupAddress.County)

		suite.Equal(&expected_street_address_1, returnedPPMShipment.DestinationAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress2, returnedPPMShipment.DestinationAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress3, returnedPPMShipment.DestinationAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.DestinationAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.DestinationAddress.City)
		suite.Equal(&state, returnedPPMShipment.DestinationAddress.State)
		suite.Equal(&county, returnedPPMShipment.DestinationAddress.County)
	})
}

func (suite *PayloadsSuite) TestUpload() {
	uploadID, _ := uuid.NewV4()
	testURL := "https://testurl.com"

	basicUpload := models.Upload{
		ID:          uploadID,
		Filename:    "fileName",
		ContentType: "image/png",
		Bytes:       1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedUpload := Upload(suite.storer, basicUpload, testURL)

		suite.IsType(returnedUpload, &ghcmessages.Upload{})
		expectedID := handlers.FmtUUIDValue(basicUpload.ID)
		suite.Equal(expectedID, returnedUpload.ID)
		suite.Equal(basicUpload.Filename, returnedUpload.Filename)
		suite.Equal(basicUpload.ContentType, returnedUpload.ContentType)
		suite.Equal(basicUpload.Bytes, returnedUpload.Bytes)
		suite.Equal(testURL, returnedUpload.URL.String())
	})
}

func (suite *PayloadsSuite) TestShipmentAddressUpdate() {
	id, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()

	newAddress := models.Address{
		StreetAddress1: "123 New St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89503",
		County:         *models.StringPointer("WASHOE"),
	}

	oldAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         *models.StringPointer("WASHOE"),
	}

	sitOriginalAddress := models.Address{
		StreetAddress1: "123 SIT St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89501",
		County:         *models.StringPointer("WASHOE"),
	}
	officeRemarks := "some office remarks"
	newSitDistanceBetween := 0
	oldSitDistanceBetween := 0

	shipmentAddressUpdate := models.ShipmentAddressUpdate{
		ID:                    id,
		ShipmentID:            id2,
		NewAddress:            newAddress,
		OriginalAddress:       oldAddress,
		SitOriginalAddress:    &sitOriginalAddress,
		ContractorRemarks:     "some remarks",
		OfficeRemarks:         &officeRemarks,
		Status:                models.ShipmentAddressUpdateStatusRequested,
		NewSitDistanceBetween: &newSitDistanceBetween,
		OldSitDistanceBetween: &oldSitDistanceBetween,
	}

	emptyShipmentAddressUpdate := models.ShipmentAddressUpdate{ID: uuid.Nil}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedShipmentAddressUpdate := ShipmentAddressUpdate(&shipmentAddressUpdate)

		suite.IsType(returnedShipmentAddressUpdate, &ghcmessages.ShipmentAddressUpdate{})
	})
	suite.Run("Failure - Returns nil", func() {
		returnedShipmentAddressUpdate := ShipmentAddressUpdate(&emptyShipmentAddressUpdate)

		suite.Nil(returnedShipmentAddressUpdate)
	})
}

func (suite *PayloadsSuite) TestWeightTicketUpload() {
	uploadID, _ := uuid.NewV4()
	testURL := "https://testurl.com"
	isWeightTicket := true

	basicUpload := models.Upload{
		ID:          uploadID,
		Filename:    "fileName",
		ContentType: "image/png",
		Bytes:       1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedUpload := WeightTicketUpload(suite.storer, basicUpload, testURL, isWeightTicket)

		suite.IsType(returnedUpload, &ghcmessages.Upload{})
		expectedID := handlers.FmtUUIDValue(basicUpload.ID)
		suite.Equal(expectedID, returnedUpload.ID)
		suite.Equal(basicUpload.Filename, returnedUpload.Filename)
		suite.Equal(basicUpload.ContentType, returnedUpload.ContentType)
		suite.Equal(basicUpload.Bytes, returnedUpload.Bytes)
		suite.Equal(testURL, returnedUpload.URL.String())
		suite.Equal(isWeightTicket, returnedUpload.IsWeightTicket)
	})
}

func (suite *PayloadsSuite) TestProofOfServiceDoc() {
	uploadID1, _ := uuid.NewV4()
	uploadID2, _ := uuid.NewV4()
	isWeightTicket := true

	// Create sample ProofOfServiceDoc
	proofOfServiceDoc := models.ProofOfServiceDoc{
		ID:               uuid.Must(uuid.NewV4()),
		PaymentRequestID: uuid.Must(uuid.NewV4()),
		IsWeightTicket:   isWeightTicket,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Create sample PrimeUploads
	primeUpload1 := models.PrimeUpload{
		ID:                  uuid.Must(uuid.NewV4()),
		ProofOfServiceDocID: uuid.Must(uuid.NewV4()),
		ContractorID:        uuid.Must(uuid.NewV4()),
		UploadID:            uploadID1,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	primeUpload2 := models.PrimeUpload{
		ID:                  uuid.Must(uuid.NewV4()),
		ProofOfServiceDocID: uuid.Must(uuid.NewV4()),
		ContractorID:        uuid.Must(uuid.NewV4()),
		UploadID:            uploadID2,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	proofOfServiceDoc.PrimeUploads = []models.PrimeUpload{primeUpload1, primeUpload2}

	suite.Run("Success - Returns a ghcmessages Proof of Service payload from a Struct", func() {
		returnedProofOfServiceDoc, _ := ProofOfServiceDoc(proofOfServiceDoc, suite.storer)

		suite.IsType(returnedProofOfServiceDoc, &ghcmessages.ProofOfServiceDoc{})
	})
}

func (suite *PayloadsSuite) TestCustomer() {
	id, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()

	residentialAddress := models.Address{
		StreetAddress1: "123 New St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89503",
		County:         *models.StringPointer("WASHOE"),
	}

	backupAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         *models.StringPointer("WASHOE"),
	}

	phone := "444-555-6677"

	firstName := "First"
	lastName := "Last"
	affiliation := models.AffiliationARMY
	email := "dontEmailMe@gmail.com"
	cacValidated := true
	customer := models.ServiceMember{
		ID:                   id,
		UserID:               id2,
		FirstName:            &firstName,
		LastName:             &lastName,
		Affiliation:          &affiliation,
		PersonalEmail:        &email,
		Telephone:            &phone,
		ResidentialAddress:   &residentialAddress,
		BackupMailingAddress: &backupAddress,
		CacValidated:         cacValidated,
	}

	suite.Run("Success - Returns a ghcmessages Customer payload from Customer Struct", func() {
		customer := Customer(&customer)

		suite.IsType(customer, &ghcmessages.Customer{})
	})
}

func (suite *PayloadsSuite) TestEntitlement() {
	entitlementID, _ := uuid.NewV4()
	dependentsAuthorized := true
	nonTemporaryStorage := true
	privatelyOwnedVehicle := true
	proGearWeight := 1000
	proGearWeightSpouse := 500
	storageInTransit := 90
	totalDependents := 2
	requiredMedicalEquipmentWeight := 200
	accompaniedTour := true
	dependentsUnderTwelve := 1
	dependentsTwelveAndOver := 1
	authorizedWeight := 8000

	entitlement := &models.Entitlement{
		ID:                             entitlementID,
		DBAuthorizedWeight:             &authorizedWeight,
		DependentsAuthorized:           &dependentsAuthorized,
		NonTemporaryStorage:            &nonTemporaryStorage,
		PrivatelyOwnedVehicle:          &privatelyOwnedVehicle,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		StorageInTransit:               &storageInTransit,
		TotalDependents:                &totalDependents,
		RequiredMedicalEquipmentWeight: requiredMedicalEquipmentWeight,
		AccompaniedTour:                &accompaniedTour,
		DependentsUnderTwelve:          &dependentsUnderTwelve,
		DependentsTwelveAndOver:        &dependentsTwelveAndOver,
		UpdatedAt:                      time.Now(),
	}

	returnedEntitlement := Entitlement(entitlement)

	suite.IsType(&ghcmessages.Entitlements{}, returnedEntitlement)

	suite.Equal(strfmt.UUID(entitlementID.String()), returnedEntitlement.ID)
	suite.Equal(authorizedWeight, int(*returnedEntitlement.AuthorizedWeight))
	suite.Equal(entitlement.DependentsAuthorized, returnedEntitlement.DependentsAuthorized)
	suite.Equal(entitlement.NonTemporaryStorage, returnedEntitlement.NonTemporaryStorage)
	suite.Equal(entitlement.PrivatelyOwnedVehicle, returnedEntitlement.PrivatelyOwnedVehicle)
	suite.Equal(int64(proGearWeight), returnedEntitlement.ProGearWeight)
	suite.Equal(int64(proGearWeightSpouse), returnedEntitlement.ProGearWeightSpouse)
	suite.Equal(storageInTransit, int(*returnedEntitlement.StorageInTransit))
	suite.Equal(totalDependents, int(returnedEntitlement.TotalDependents))
	suite.Equal(int64(requiredMedicalEquipmentWeight), returnedEntitlement.RequiredMedicalEquipmentWeight)
	suite.Equal(models.BoolPointer(accompaniedTour), returnedEntitlement.AccompaniedTour)
	suite.Equal(dependentsUnderTwelve, int(*returnedEntitlement.DependentsUnderTwelve))
	suite.Equal(dependentsTwelveAndOver, int(*returnedEntitlement.DependentsTwelveAndOver))
}

func (suite *PayloadsSuite) TestCreateCustomer() {
	id, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	oktaID := "thisIsNotARealID"

	oktaUser := models.CreatedOktaUser{
		ID: oktaID,
		Profile: struct {
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			MobilePhone string `json:"mobilePhone"`
			SecondEmail string `json:"secondEmail"`
			Login       string `json:"login"`
			Email       string `json:"email"`
		}{
			Email: "john.doe@example.com",
		},
	}

	residentialAddress := models.Address{
		StreetAddress1: "123 New St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89503",
		County:         *models.StringPointer("WASHOE"),
	}

	backupAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         *models.StringPointer("WASHOE"),
	}

	phone := "444-555-6677"
	backupContact := models.BackupContact{
		Name:  "Billy Bob",
		Email: "billBob@mail.mil",
		Phone: &phone,
	}

	firstName := "First"
	lastName := "Last"
	affiliation := models.AffiliationARMY
	email := "dontEmailMe@gmail.com"
	sm := models.ServiceMember{
		ID:                   id,
		UserID:               id2,
		FirstName:            &firstName,
		LastName:             &lastName,
		Affiliation:          &affiliation,
		PersonalEmail:        &email,
		Telephone:            &phone,
		ResidentialAddress:   &residentialAddress,
		BackupMailingAddress: &backupAddress,
	}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedShipmentAddressUpdate := CreatedCustomer(&sm, &oktaUser, &backupContact)

		suite.IsType(returnedShipmentAddressUpdate, &ghcmessages.CreatedCustomer{})
	})
}

func (suite *PayloadsSuite) TestSearchMoves() {
	appCtx := suite.AppContextForTest()

	marines := models.AffiliationMARINES
	moveUSMC := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &marines,
			},
		},
	}, nil)

	moves := models.Moves{moveUSMC}
	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		payload := SearchMoves(appCtx, moves)

		suite.IsType(payload, &ghcmessages.SearchMoves{})
		suite.NotNil(payload)
	})
}

func (suite *PayloadsSuite) TestMarketCode() {
	suite.Run("returns nil when marketCode is nil", func() {
		var marketCode *models.MarketCode = nil
		result := MarketCode(marketCode)
		suite.Equal(result, "")
	})

	suite.Run("returns string when marketCode is not nil", func() {
		marketCodeDomestic := models.MarketCodeDomestic
		result := MarketCode(&marketCodeDomestic)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("d", result, "Expected result to be 'd' for domestic market code")
	})

	suite.Run("returns string when marketCode is international", func() {
		marketCodeInternational := models.MarketCodeInternational
		result := MarketCode(&marketCodeInternational)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("i", result, "Expected result to be 'i' for international market code")
	})
}

func (suite *PayloadsSuite) TestGsrAppeal() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

	suite.Run("returns nil when gsrAppeal is nil", func() {
		var gsrAppeal *models.GsrAppeal = nil
		result := GsrAppeal(gsrAppeal)
		suite.Nil(result, "Expected result to be nil when gsrAppeal is nil")
	})

	suite.Run("correctly maps GsrAppeal with all fields populated", func() {
		gsrAppealID := uuid.Must(uuid.NewV4())
		reportViolationID := uuid.Must(uuid.NewV4())
		evaluationReportID := uuid.Must(uuid.NewV4())
		appealStatus := models.AppealStatusSustained
		isSeriousIncident := true
		remarks := "Sample remarks"
		createdAt := time.Now()

		gsrAppeal := &models.GsrAppeal{
			ID:                      gsrAppealID,
			ReportViolationID:       &reportViolationID,
			EvaluationReportID:      evaluationReportID,
			OfficeUser:              &officeUser,
			OfficeUserID:            officeUser.ID,
			IsSeriousIncidentAppeal: &isSeriousIncident,
			AppealStatus:            appealStatus,
			Remarks:                 remarks,
			CreatedAt:               createdAt,
		}

		result := GsrAppeal(gsrAppeal)

		suite.NotNil(result, "Expected result to not be nil when gsrAppeal has values")
		suite.Equal(handlers.FmtUUID(gsrAppealID), &result.ID, "Expected ID to match")
		suite.Equal(handlers.FmtUUID(reportViolationID), &result.ViolationID, "Expected ViolationID to match")
		suite.Equal(handlers.FmtUUID(evaluationReportID), &result.ReportID, "Expected ReportID to match")
		suite.Equal(handlers.FmtUUID(officeUser.ID), &result.OfficeUserID, "Expected OfficeUserID to match")
		suite.Equal(ghcmessages.GSRAppealStatusType(appealStatus), result.AppealStatus, "Expected AppealStatus to match")
		suite.Equal(remarks, result.Remarks, "Expected Remarks to match")
		suite.Equal(strfmt.DateTime(createdAt), result.CreatedAt, "Expected CreatedAt to match")
		suite.True(result.IsSeriousIncident, "Expected IsSeriousIncident to be true")
	})

	suite.Run("handles nil ReportViolationID without panic", func() {
		gsrAppealID := uuid.Must(uuid.NewV4())
		evaluationReportID := uuid.Must(uuid.NewV4())
		isSeriousIncident := false
		appealStatus := models.AppealStatusRejected
		remarks := "Sample remarks"
		createdAt := time.Now()

		gsrAppeal := &models.GsrAppeal{
			ID:                      gsrAppealID,
			ReportViolationID:       nil,
			EvaluationReportID:      evaluationReportID,
			OfficeUser:              &officeUser,
			OfficeUserID:            officeUser.ID,
			IsSeriousIncidentAppeal: &isSeriousIncident,
			AppealStatus:            appealStatus,
			Remarks:                 remarks,
			CreatedAt:               createdAt,
		}

		result := GsrAppeal(gsrAppeal)

		suite.NotNil(result, "Expected result to not be nil when gsrAppeal has values")
		suite.Equal(handlers.FmtUUID(gsrAppealID), &result.ID, "Expected ID to match")
		suite.Equal(strfmt.UUID(""), result.ViolationID, "Expected ViolationID to be nil when ReportViolationID is nil")
		suite.Equal(handlers.FmtUUID(evaluationReportID), &result.ReportID, "Expected ReportID to match")
		suite.Equal(handlers.FmtUUID(officeUser.ID), &result.OfficeUserID, "Expected OfficeUserID to match")
		suite.Equal(ghcmessages.GSRAppealStatusType(appealStatus), result.AppealStatus, "Expected AppealStatus to match")
		suite.Equal(remarks, result.Remarks, "Expected Remarks to match")
		suite.Equal(strfmt.DateTime(createdAt), result.CreatedAt, "Expected CreatedAt to match")
		suite.False(result.IsSeriousIncident, "Expected IsSeriousIncident to be false")
	})
}
