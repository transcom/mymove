package adminapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/services/upload"

	"github.com/stretchr/testify/mock"

	"github.com/gofrs/uuid"

	uploadop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/upload"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetUploadHandler() {
	setupTestData := func() (models.UserUpload, models.Move) {
		sm := testdatagen.MakeDefaultServiceMember(suite.DB())
		suite.MustSave(&sm)

		orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
			},
		})
		suite.MustSave(&orders)

		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Orders:   orders,
				OrdersID: orders.ID,
			},
		})
		suite.MustSave(&move)

		document := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMember:   sm,
				ServiceMemberID: sm.ID,
			},
		})
		suite.MustSave(&document)

		uploadID, _ := uuid.NewV4()
		uploadUserAssertions := models.UserUpload{
			Document:   document,
			DocumentID: &document.ID,
			CreatedAt:  time.Now(),
			UploaderID: sm.UserID,
			Upload: models.Upload{
				ID:          uploadID,
				Filename:    "FileName",
				Bytes:       int64(15),
				ContentType: "application/pdf",
				CreatedAt:   time.Now(),
			},
		}

		uploadInstance := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{UserUpload: uploadUserAssertions})
		suite.MustSave(&uploadInstance)

		return uploadInstance, move
	}

	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		uploadInstance, move := setupTestData()
		params := uploadop.GetUploadParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/uploads/%s", uploadInstance.UploadID.String())),
			UploadID:    *handlers.FmtUUID(uploadInstance.UploadID),
		}

		uploadInformationFetcher := upload.NewUploadInformationFetcher()
		handler := GetUploadHandler{
			HandlerConfig:            handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			UploadInformationFetcher: uploadInformationFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&uploadop.GetUploadOK{}, response)
		okResponse := response.(*uploadop.GetUploadOK)
		suite.Equal(*handlers.FmtUUID(uploadInstance.UploadID), okResponse.Payload.ID)
		suite.Equal(move.Locator, *okResponse.Payload.MoveLocator)
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		uploadInstance, _ := setupTestData()
		params := uploadop.GetUploadParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/uploads/%s", uploadInstance.UploadID.String())),
			UploadID:    *handlers.FmtUUID(uploadInstance.UploadID),
		}
		expectedError := models.ErrFetchNotFound
		uploadInformationFetcher := &mocks.UploadInformationFetcher{}
		uploadInformationFetcher.On("FetchUploadInformation",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(services.UploadInformation{}, expectedError).Once()
		handler := GetUploadHandler{
			HandlerConfig:            handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			UploadInformationFetcher: uploadInformationFetcher,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
