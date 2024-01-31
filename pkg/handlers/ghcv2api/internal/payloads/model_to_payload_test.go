package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcv2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

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

		suite.IsType(returnedUpload, &ghcv2messages.Upload{})
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

		suite.IsType(returnedUpload, &ghcv2messages.Upload{})
		expectedID := handlers.FmtUUIDValue(basicUpload.ID)
		suite.Equal(expectedID, returnedUpload.ID)
		suite.Equal(basicUpload.Filename, returnedUpload.Filename)
		suite.Equal(basicUpload.ContentType, returnedUpload.ContentType)
		suite.Equal(basicUpload.Bytes, returnedUpload.Bytes)
		suite.Equal(testURL, returnedUpload.URL.String())
		suite.Equal(isWeightTicket, returnedUpload.IsWeightTicket)
	})
}
