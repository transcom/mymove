package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/mocks"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestFetchPPMShipment() {

	ppmShipmentID, _ := uuid.NewV4()
	streetAddress1 := "MacDill AFB"
	streetAddress2, streetAddress3 := "", ""
	city := "Tampa"
	state := "FL"
	postalcode := "33621"
	country := models.Country{
		Country:     "US",
		CountryName: "United States",
	}
	county := "HILLSBOROUGH"

	expectedAddress := models.Address{
		StreetAddress1: streetAddress1,
		StreetAddress2: &streetAddress2,
		StreetAddress3: &streetAddress3,
		City:           city,
		State:          state,
		PostalCode:     postalcode,
		Country:        &country,
		County:         &county,
	}

	isActualExpenseReimbursement := true
	hasGunSafe := models.BoolPointer(true)
	gunSafeWeight := models.PoundPointer(333)

	expectedPPMShipment := models.PPMShipment{
		ID:                           ppmShipmentID,
		PPMType:                      models.PPMTypeActualExpense,
		PickupAddress:                &expectedAddress,
		DestinationAddress:           &expectedAddress,
		IsActualExpenseReimbursement: &isActualExpenseReimbursement,
		HasGunSafe:                   hasGunSafe,
		GunSafeWeight:                gunSafeWeight,
	}

	suite.Run("Success -", func() {
		returnedPPMShipment := PPMShipment(nil, &expectedPPMShipment)

		suite.IsType(&internalmessages.PPMShipment{}, returnedPPMShipment)
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

		suite.Equal(internalmessages.PPMType(models.PPMTypeActualExpense), returnedPPMShipment.PpmType)
		suite.True(*returnedPPMShipment.IsActualExpenseReimbursement)
		suite.Equal(handlers.FmtBool(*hasGunSafe), returnedPPMShipment.HasGunSafe)
		suite.Equal(handlers.FmtPoundPtr(gunSafeWeight), returnedPPMShipment.GunSafeWeight)
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

func (suite *PayloadsSuite) TestVLocation() {
	suite.Run("correctly maps VLocation with all fields populated", func() {
		city := "LOS ANGELES"
		state := "CA"
		postalCode := "90210"
		county := "LOS ANGELES"
		usPostRegionCityID := uuid.Must(uuid.NewV4())

		vLocation := &models.VLocation{
			CityName:             city,
			StateName:            state,
			UsprZipID:            postalCode,
			UsprcCountyNm:        county,
			UsPostRegionCitiesID: &usPostRegionCityID,
		}

		payload := VLocation(vLocation)

		suite.IsType(payload, &internalmessages.VLocation{})
		suite.Equal(handlers.FmtUUID(usPostRegionCityID), &payload.UsPostRegionCitiesID, "Expected UsPostRegionCitiesID to match")
		suite.Equal(city, payload.City, "Expected City to match")
		suite.Equal(state, payload.State, "Expected State to match")
		suite.Equal(postalCode, payload.PostalCode, "Expected PostalCode to match")
		suite.Equal(county, *(payload.County), "Expected County to match")
	})
}

func (suite *PayloadsSuite) TestSignedCertification() {
	suite.Run("Certification model", func() {
		uuid, _ := uuid.NewV4()
		certType := models.SignedCertificationTypeHHG
		model := models.SignedCertification{
			ID:                uuid,
			SubmittingUserID:  uuid,
			MoveID:            uuid,
			PpmID:             &uuid,
			CertificationText: "dummy",
			CertificationType: &certType,
		}
		parsedSignedCert := SignedCertification(&model)
		suite.NotNil(parsedSignedCert)
		suite.Equal(uuid.String(), parsedSignedCert.ID.String())
		suite.Equal(uuid.String(), parsedSignedCert.SubmittingUserID.String())
		suite.Equal(uuid.String(), parsedSignedCert.MoveID.String())
		suite.Equal(uuid.String(), parsedSignedCert.PpmID.String())
		suite.Equal("dummy", *parsedSignedCert.CertificationText)
		suite.Equal(string(certType), string(parsedSignedCert.CertificationType))
	})
}

func (suite *PayloadsSuite) TestWeightTicket() {
	suite.Run("WeightTicket model", func() {
		mockStorer := &mocks.FileStorer{}

		weightTicketID := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		documentID := uuid.Must(uuid.NewV4())
		now := time.Now()

		weightTicket := &models.WeightTicket{
			ID:                                weightTicketID,
			PPMShipmentID:                     ppmShipmentID,
			EmptyWeight:                       models.PoundPointer(4000),
			SubmittedEmptyWeight:              models.PoundPointer(4200),
			EmptyDocumentID:                   documentID,
			FullWeight:                        models.PoundPointer(6000),
			SubmittedFullWeight:               models.PoundPointer(6200),
			FullDocumentID:                    documentID,
			ProofOfTrailerOwnershipDocumentID: documentID,
			AdjustedNetWeight:                 models.PoundPointer(2000),
			NetWeightRemarks:                  models.StringPointer("Test remarks"),
			CreatedAt:                         now,
			UpdatedAt:                         now,
		}
		parsedWeightTicket := WeightTicket(mockStorer, weightTicket)
		suite.NotNil(parsedWeightTicket)
		suite.Equal(weightTicketID.String(), parsedWeightTicket.ID.String())
		suite.Equal(ppmShipmentID.String(), parsedWeightTicket.PpmShipmentID.String())
		suite.Equal(handlers.FmtPoundPtr(weightTicket.EmptyWeight), parsedWeightTicket.EmptyWeight)
		suite.Equal(handlers.FmtPoundPtr(weightTicket.SubmittedEmptyWeight), parsedWeightTicket.SubmittedEmptyWeight)
		suite.Equal(handlers.FmtUUID(weightTicket.EmptyDocumentID), &parsedWeightTicket.EmptyDocumentID)
		suite.Equal(handlers.FmtPoundPtr(weightTicket.FullWeight), parsedWeightTicket.FullWeight)
		suite.Equal(handlers.FmtPoundPtr(weightTicket.SubmittedFullWeight), parsedWeightTicket.SubmittedFullWeight)
		suite.Equal(handlers.FmtUUID(weightTicket.FullDocumentID), &parsedWeightTicket.FullDocumentID)
		suite.Equal(handlers.FmtUUID(weightTicket.ProofOfTrailerOwnershipDocumentID), &parsedWeightTicket.ProofOfTrailerOwnershipDocumentID)
		suite.Equal(handlers.FmtPoundPtr(weightTicket.AdjustedNetWeight), parsedWeightTicket.AdjustedNetWeight)
		suite.Equal("Test remarks", *parsedWeightTicket.NetWeightRemarks)
		suite.Equal(etag.GenerateEtag(weightTicket.UpdatedAt), parsedWeightTicket.ETag)
	})
}

func (suite *PayloadsSuite) TestCounselingOffices() {
	suite.Run("correctly maps transportaion offices to counseling offices payload", func() {
		office1 := factory.BuildTransportationOffice(nil, []factory.Customization{
			{
				Model: models.TransportationOffice{
					ID:   uuid.Must(uuid.NewV4()),
					Name: "PPPO Fort Liberty",
				},
			},
		}, nil)

		office2 := factory.BuildTransportationOffice(nil, []factory.Customization{
			{
				Model: models.TransportationOffice{
					ID:   uuid.Must(uuid.NewV4()),
					Name: "PPPO Fort Walker",
				},
			},
		}, nil)

		offices := models.TransportationOffices{office1, office2}

		payload := CounselingOffices(offices)

		suite.IsType(payload, internalmessages.CounselingOffices{})
		suite.Equal(2, len(payload))
		suite.Equal(office1.ID.String(), payload[0].ID.String())
		suite.Equal(office2.ID.String(), payload[1].ID.String())
	})
}

func (suite *PayloadsSuite) TestMovingExpense() {
	mockStorer := &mocks.FileStorer{}

	suite.Run("successfully converts a fully populated MovingExpense", func() {
		document := factory.BuildDocument(suite.DB(), nil, nil)
		id := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		documentID := document.ID
		now := time.Now()
		description := "Test description"
		submittedDescription := "Submitted description"
		paidWithGTCC := true
		amount := unit.Cents(1000)
		submittedAmount := unit.Cents(1100)
		missingReceipt := false
		movingExpenseType := models.MovingExpenseReceiptTypeSmallPackage
		submittedMovingExpenseType := models.MovingExpenseReceiptTypeSmallPackage
		status := models.PPMDocumentStatusApproved
		reason := "Some reason"
		sitStartDate := now.AddDate(0, -1, 0)
		submittedSitStartDate := now.AddDate(0, -1, 1)
		sitEndDate := now.AddDate(0, 0, -10)
		submittedSitEndDate := now.AddDate(0, 0, -9)
		weightStored := unit.Pound(150)
		sitLocation := models.SITLocationTypeOrigin
		sitReimburseableAmount := unit.Cents(2000)
		trackingNumber := "tracking123"
		weightShipped := unit.Pound(100)
		isProGear := true
		proGearBelongsToSelf := false
		proGearDescription := "Pro gear desc"

		expense := &models.MovingExpense{
			ID:                         id,
			PPMShipmentID:              ppmShipmentID,
			Document:                   document,
			DocumentID:                 documentID,
			CreatedAt:                  now,
			UpdatedAt:                  now,
			Description:                &description,
			SubmittedDescription:       &submittedDescription,
			PaidWithGTCC:               &paidWithGTCC,
			Amount:                     &amount,
			SubmittedAmount:            &submittedAmount,
			MissingReceipt:             &missingReceipt,
			MovingExpenseType:          &movingExpenseType,
			SubmittedMovingExpenseType: &submittedMovingExpenseType,
			Status:                     &status,
			Reason:                     &reason,
			SITStartDate:               &sitStartDate,
			SubmittedSITStartDate:      &submittedSitStartDate,
			SITEndDate:                 &sitEndDate,
			SubmittedSITEndDate:        &submittedSitEndDate,
			WeightStored:               &weightStored,
			SITLocation:                &sitLocation,
			SITReimburseableAmount:     &sitReimburseableAmount,
			TrackingNumber:             &trackingNumber,
			WeightShipped:              &weightShipped,
			IsProGear:                  &isProGear,
			ProGearBelongsToSelf:       &proGearBelongsToSelf,
			ProGearDescription:         &proGearDescription,
		}

		result := MovingExpense(mockStorer, expense)
		suite.NotNil(result, "Expected non-nil payload for valid input")

		// Check required fields.
		suite.Equal(*handlers.FmtUUID(id), result.ID, "ID should match")
		suite.Equal(*handlers.FmtUUID(ppmShipmentID), result.PpmShipmentID, "PPMShipmentID should match")
		suite.Equal(*handlers.FmtUUID(documentID), result.DocumentID, "DocumentID should match")
		suite.NotNil(result.Document)
		suite.Equal(strfmt.DateTime(now), result.CreatedAt, "CreatedAt should match")
		suite.Equal(strfmt.DateTime(now), result.UpdatedAt, "UpdatedAt should match")
		suite.Equal(description, *result.Description, "Description should match")
		suite.Equal(submittedDescription, *result.SubmittedDescription, "SubmittedDescription should match")
		suite.Equal(paidWithGTCC, *result.PaidWithGtcc, "PaidWithGTCC should match")
		suite.Equal(handlers.FmtCost(&amount), result.Amount, "Amount should match")
		suite.Equal(handlers.FmtCost(&submittedAmount), result.SubmittedAmount, "SubmittedAmount should match")
		suite.Equal(missingReceipt, *result.MissingReceipt, "MissingReceipt should match")
		suite.Equal(etag.GenerateEtag(now), result.ETag, "ETag should be generated from UpdatedAt")

		// Check optional fields.
		if expense.MovingExpenseType != nil {
			expectedType := internalmessages.OmittableMovingExpenseType(*expense.MovingExpenseType)
			suite.Equal(&expectedType, result.MovingExpenseType, "MovingExpenseType should match")
		}
		if expense.SubmittedMovingExpenseType != nil {
			expectedSubmittedType := internalmessages.SubmittedMovingExpenseType(*expense.SubmittedMovingExpenseType)
			suite.Equal(expectedSubmittedType, *result.SubmittedMovingExpenseType, "SubmittedMovingExpenseType should match")
		}
		if expense.Status != nil {
			expectedStatus := internalmessages.OmittablePPMDocumentStatus(*expense.Status)
			suite.Equal(&expectedStatus, result.Status, "Status should match")
		}
		if expense.Reason != nil {
			expectedReason := internalmessages.PPMDocumentStatusReason(*expense.Reason)
			suite.Equal(&expectedReason, result.Reason, "Reason should match")
		}
		suite.Equal(handlers.FmtDatePtr(&sitStartDate), result.SitStartDate, "SITStartDate should match")
		suite.Equal(handlers.FmtDatePtr(&submittedSitStartDate), result.SubmittedSitStartDate, "SubmittedSITStartDate should match")
		suite.Equal(handlers.FmtDatePtr(&sitEndDate), result.SitEndDate, "SITEndDate should match")
		suite.Equal(handlers.FmtDatePtr(&submittedSitEndDate), result.SubmittedSitEndDate, "SubmittedSITEndDate should match")
		suite.Equal(handlers.FmtPoundPtr(&weightStored), result.WeightStored, "WeightStored should match")
		if expense.SITLocation != nil {
			expectedSitLocation := internalmessages.SITLocationType(*expense.SITLocation)
			suite.Equal(&expectedSitLocation, result.SitLocation, "SITLocation should match")
		}
		suite.Equal(handlers.FmtCost(&sitReimburseableAmount), result.SitReimburseableAmount, "SITReimburseableAmount should match")
		suite.Equal(&trackingNumber, result.TrackingNumber, "TrackingNumber should match")
		suite.Equal(handlers.FmtPoundPtr(&weightShipped), result.WeightShipped, "WeightShipped should match")
		suite.Equal(expense.IsProGear, result.IsProGear, "IsProGear should match")
		suite.Equal(expense.ProGearBelongsToSelf, result.ProGearBelongsToSelf, "ProGearBelongsToSelf should match")
		suite.Equal(proGearDescription, result.ProGearDescription, "ProGearDescription should match")
	})
}

func (suite *PayloadsSuite) TestPayGradesForInternalPayloadToModel() {
	payGrades := models.PayGrades{
		{Grade: string(models.ServiceMemberGradeE1), GradeDescription: models.StringPointer(string(models.ServiceMemberGradeE1))},
		{Grade: string(models.ServiceMemberGradeO3), GradeDescription: models.StringPointer(string(models.ServiceMemberGradeO3))},
		{Grade: string(models.ServiceMemberGradeW2), GradeDescription: models.StringPointer(string(models.ServiceMemberGradeW2))},
	}
	for _, payGrade := range payGrades {
		suite.Run(payGrade.Grade, func() {
			grades := models.PayGrades{payGrade}
			result := PayGrades(grades)

			suite.Require().Len(result, 1)
			actual := result[0]

			suite.Equal(payGrade.Grade, actual.Grade)
			suite.Equal(*payGrade.GradeDescription, actual.Description)
		})
	}
}

func (suite *PayloadsSuite) TestCountriesPayload() {
	suite.Run("Correctly transform array of countries into payload", func() {
		countries := make([]models.Country, 0)
		countries = append(countries, models.Country{Country: "US", CountryName: "UNITED STATES"})
		payload := Countries(countries)
		suite.True(len(payload) == 1)
		suite.Equal(payload[0].Code, "US")
		suite.Equal(payload[0].Name, "UNITED STATES")
	})

	suite.Run("empty array of countries into payload", func() {
		countries := make([]models.Country, 0)
		payload := Countries(countries)
		suite.True(len(payload) == 0)
	})

	suite.Run("nil countries into payload", func() {
		payload := Countries(nil)
		suite.True(len(payload) == 0)
	})
}

func (suite *PayloadsSuite) TestGunSafe() {
	mockStorer := &mocks.FileStorer{}

	suite.Run("successfully converts a fully populated GunSafeWeightTicket", func() {
		document := factory.BuildDocument(suite.DB(), nil, nil)
		id := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		documentID := document.ID
		now := time.Now()
		description := "Test description"
		status := models.PPMDocumentStatusApproved
		reason := "Some reason"
		weight := unit.Pound(150)
		submittedWeight := unit.Pound(200)

		gunSafe := &models.GunSafeWeightTicket{
			ID:              id,
			PPMShipmentID:   ppmShipmentID,
			Document:        document,
			DocumentID:      documentID,
			CreatedAt:       now,
			UpdatedAt:       now,
			Description:     &description,
			Status:          &status,
			Weight:          &weight,
			SubmittedWeight: &submittedWeight,
			Reason:          &reason,
		}

		result := GunSafeWeightTicket(mockStorer, gunSafe)
		suite.NotNil(result, "Expected non-nil payload for valid input")

		// Check required fields.
		suite.Equal(*handlers.FmtUUID(id), result.ID, "ID should match")
		suite.Equal(*handlers.FmtUUID(ppmShipmentID), result.PpmShipmentID, "PPMShipmentID should match")
		suite.Equal(*handlers.FmtUUID(documentID), result.DocumentID, "DocumentID should match")
		suite.NotNil(result.Document)
		suite.Equal(strfmt.DateTime(now), result.CreatedAt, "CreatedAt should match")
		suite.Equal(strfmt.DateTime(now), result.UpdatedAt, "UpdatedAt should match")
		suite.Equal(description, *result.Description, "Description should match")
		suite.Equal(etag.GenerateEtag(now), result.ETag, "ETag should be generated from UpdatedAt")

		// Check optional fields.
		suite.Equal(handlers.FmtPoundPtr(&weight), result.Weight, "Weight should match")
		suite.Equal(handlers.FmtPoundPtr(&submittedWeight), result.SubmittedWeight, "SubmittedWeight should match")
		if gunSafe.Status != nil {
			expectedStatus := internalmessages.OmittablePPMDocumentStatus(*gunSafe.Status)
			suite.Equal(&expectedStatus, result.Status, "Status should match")
		}
		if gunSafe.Reason != nil {
			expectedReason := internalmessages.PPMDocumentStatusReason(*gunSafe.Reason)
			suite.Equal(&expectedReason, result.Reason, "Reason should match")
		}
	})

	suite.Run("successfully converts an array of fully populated GunSafeWeightTickets", func() {
		document := factory.BuildDocument(suite.DB(), nil, nil)
		id := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		documentID := document.ID
		now := time.Now()
		description := "Test description"
		status := models.PPMDocumentStatusApproved
		reason := "Some reason"
		weight := unit.Pound(150)
		submittedWeight := unit.Pound(200)

		secondDocument := factory.BuildDocument(suite.DB(), nil, nil)
		secondID := uuid.Must(uuid.NewV4())
		secondDocumentID := document.ID
		secondDescription := "Another Test description"
		secondStatus := models.PPMDocumentStatusRejected
		secondReason := "Some other reason"
		secondWeight := unit.Pound(225)
		secondSubmittedWeight := unit.Pound(190)

		gunSafeWeightTickets := &models.GunSafeWeightTickets{
			models.GunSafeWeightTicket{
				ID:              id,
				PPMShipmentID:   ppmShipmentID,
				Document:        document,
				DocumentID:      documentID,
				CreatedAt:       now,
				UpdatedAt:       now,
				Description:     &description,
				Status:          &status,
				Weight:          &weight,
				SubmittedWeight: &submittedWeight,
				Reason:          &reason,
			},
			models.GunSafeWeightTicket{
				ID:              secondID,
				PPMShipmentID:   ppmShipmentID,
				Document:        secondDocument,
				DocumentID:      secondDocumentID,
				CreatedAt:       now,
				UpdatedAt:       now,
				Description:     &secondDescription,
				Status:          &secondStatus,
				Weight:          &secondWeight,
				SubmittedWeight: &secondSubmittedWeight,
				Reason:          &secondReason,
			},
		}

		result := GunSafeWeightTickets(mockStorer, *gunSafeWeightTickets)
		suite.NotNil(result, "Expected non-nil payload for valid input")

		// Check required fields.
		suite.Equal(*handlers.FmtUUID(id), result[0].ID, "ID should match")
		suite.Equal(*handlers.FmtUUID(ppmShipmentID), result[0].PpmShipmentID, "PPMShipmentID should match")
		suite.Equal(*handlers.FmtUUID(documentID), result[0].DocumentID, "DocumentID should match")
		suite.NotNil(result[0].Document)
		suite.Equal(strfmt.DateTime(now), result[0].CreatedAt, "CreatedAt should match")
		suite.Equal(strfmt.DateTime(now), result[0].UpdatedAt, "UpdatedAt should match")
		suite.Equal(description, *result[0].Description, "Description should match")
		suite.Equal(etag.GenerateEtag(now), result[0].ETag, "ETag should be generated from UpdatedAt")

		// Check optional fields.
		suite.Equal(handlers.FmtPoundPtr(&weight), result[0].Weight, "Weight should match")
		suite.Equal(handlers.FmtPoundPtr(&submittedWeight), result[0].SubmittedWeight, "SubmittedWeight should match")
		if result[0].Status != nil {
			expectedStatus := internalmessages.OmittablePPMDocumentStatus(*result[0].Status)
			suite.Equal(&expectedStatus, result[0].Status, "Status should match")
		}
		if result[0].Reason != nil {
			expectedReason := internalmessages.PPMDocumentStatusReason(*result[0].Reason)
			suite.Equal(&expectedReason, result[0].Reason, "Reason should match")
		}

		// Second gun safe ticket
		suite.Equal(*handlers.FmtUUID(secondID), result[1].ID, "ID should match")
		suite.Equal(*handlers.FmtUUID(ppmShipmentID), result[1].PpmShipmentID, "PPMShipmentID should match")
		suite.Equal(*handlers.FmtUUID(secondDocumentID), result[1].DocumentID, "DocumentID should match")
		suite.NotNil(result[1].Document)
		suite.Equal(strfmt.DateTime(now), result[1].CreatedAt, "CreatedAt should match")
		suite.Equal(strfmt.DateTime(now), result[1].UpdatedAt, "UpdatedAt should match")
		suite.Equal(secondDescription, *result[1].Description, "Description should match")
		suite.Equal(etag.GenerateEtag(now), result[1].ETag, "ETag should be generated from UpdatedAt")

		suite.Equal(handlers.FmtPoundPtr(&secondWeight), result[1].Weight, "Weight should match")
		suite.Equal(handlers.FmtPoundPtr(&secondSubmittedWeight), result[1].SubmittedWeight, "SubmittedWeight should match")
		if result[1].Status != nil {
			expectedStatus := internalmessages.OmittablePPMDocumentStatus(*result[1].Status)
			suite.Equal(&expectedStatus, result[1].Status, "Status should match")
		}
		if result[1].Reason != nil {
			expectedReason := internalmessages.PPMDocumentStatusReason(*result[1].Reason)
			suite.Equal(&expectedReason, result[1].Reason, "Reason should match")
		}
	})
}

func (suite *PayloadsSuite) TestPayGrades() {
	payGrades := models.PayGrades{
		{Grade: "E-1", GradeDescription: models.StringPointer("E-1")},
		{Grade: "O-3", GradeDescription: models.StringPointer("O-3")},
		{Grade: "W-2", GradeDescription: models.StringPointer("W-2")},
	}
	for _, payGrade := range payGrades {
		suite.Run(payGrade.Grade, func() {
			grades := models.PayGrades{payGrade}
			result := PayGrades(grades)

			suite.Require().Len(result, 1)
			actual := result[0]

			suite.Equal(payGrade.Grade, actual.Grade)
			suite.Equal(*payGrade.GradeDescription, actual.Description)
		})
	}
}
