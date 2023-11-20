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
