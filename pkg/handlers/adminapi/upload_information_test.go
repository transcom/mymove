package adminapi

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	uploadop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/upload"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetUploadHandler() {
	setupTestData := func() (models.UserUpload, models.Move) {
		sm := factory.BuildServiceMember(suite.DB(), nil, nil)
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

		document := factory.BuildDocumentLinkServiceMember(suite.DB(), sm)
		suite.MustSave(&document)

		uploadInstance := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    document,
				LinkOnly: true,
			},
			{
				Model: models.Upload{
					Filename:    "FileName",
					Bytes:       int64(15),
					ContentType: "application/pdf",
				},
			},
		}, nil)
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
			HandlerConfig:            suite.HandlerConfig(),
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
			HandlerConfig:            suite.HandlerConfig(),
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
