package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/mocks"
)

func (suite *PayloadsSuite) TestFetchPPMShipment() {

	ppmShipmentID, _ := uuid.NewV4()
	streetAddress1 := "MacDill AFB"
	streetAddress2, streetAddress3 := "", ""
	city := "Tampa"
	state := "FL"
	postalcode := "33621"
	country := models.Country{
		Country: "US",
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

	expectedPPMShipment := models.PPMShipment{
		ID:                           ppmShipmentID,
		PickupAddress:                &expectedAddress,
		DestinationAddress:           &expectedAddress,
		IsActualExpenseReimbursement: &isActualExpenseReimbursement,
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

		suite.True(*returnedPPMShipment.IsActualExpenseReimbursement)
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
