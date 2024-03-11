package payloads

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func TestOrder(_ *testing.T) {
	order := &models.Order{}
	Order(order)
}

// TestMove makes sure zero values/optional fields are handled
func TestMove(_ *testing.T) {
	Move(&models.Move{})
}

func (suite *PayloadsSuite) TestFetchPPMShipment() {

	ppmShipmentID, _ := uuid.NewV4()
	streetAddress1 := "MacDill AFB"
	streetAddress2, streetAddress3 := "", ""
	city := "Tampa"
	state := "FL"
	postalcode := "33621"
	country := "US"

	expectedAddress := models.Address{
		StreetAddress1: streetAddress1,
		StreetAddress2: &streetAddress2,
		StreetAddress3: &streetAddress3,
		City:           city,
		State:          "FL",
		PostalCode:     "33621",
		Country:        &country,
	}

	expectedPPMShipment := models.PPMShipment{
		ID:                 ppmShipmentID,
		PickupAddress:      &expectedAddress,
		DestinationAddress: &expectedAddress,
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
		suite.Equal(&country, returnedPPMShipment.PickupAddress.Country)

		suite.Equal(&streetAddress1, returnedPPMShipment.DestinationAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress2, returnedPPMShipment.DestinationAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress3, returnedPPMShipment.DestinationAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.DestinationAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.DestinationAddress.City)
		suite.Equal(&state, returnedPPMShipment.DestinationAddress.State)
		suite.Equal(&country, returnedPPMShipment.DestinationAddress.Country)
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
		Country:        models.StringPointer("United States"),
	}

	oldAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		Country:        models.StringPointer("United States"),
	}

	sitOriginalAddress := models.Address{
		StreetAddress1: "123 SIT St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89501",
		Country:        models.StringPointer("United States"),
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
