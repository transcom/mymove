package primeapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/models"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
)

func (suite *HandlerSuite) CreateServiceRequestDocumentUploadHandler() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	setupTestData := func() (CreateServiceRequestDocumentUploadHandler, models.MTOServiceItem) {
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handler := CreateServiceRequestDocumentUploadHandler{
			handlerConfig,
			mtoserviceitem.NewServiceRequestDocumentUploadCreator(handlerConfig.FileStorer()),
		}
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)

		factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
		return handler, mtoServiceItem
	}

	suite.Run("successful create service request document upload", func() {
		primeUser := factory.BuildUser(nil, nil, nil)
		handler, mtoServiceItem := setupTestData()
		req := httptest.NewRequest("POST", fmt.Sprintf("/mto_service_items/%s/uploads", mtoServiceItem.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)

		file := suite.Fixture("test.pdf")

		params := mtoserviceitemops.CreateServiceRequestDocumentUploadParams{
			HTTPRequest:      req,
			File:             file,
			MtoServiceItemID: mtoServiceItem.ID.String(),
		}

		response := handler.Handle(params)

		suite.IsType(&mtoserviceitemops.CreateServiceRequestDocumentUploadCreated{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateServiceRequestDocumentUploadCreated).Payload

		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("create service request document upload fail - invalid service item ID format", func() {
		primeUser := factory.BuildUser(nil, nil, nil)
		handler, mtoServiceItem := setupTestData()

		badFormatID := strfmt.UUID("gb7b134a-7c44-45f2-9114-bb0831cc5db3")
		file := suite.Fixture("test.pdf")

		req := httptest.NewRequest("POST", fmt.Sprintf("/mto_service_items/%s/uploads", mtoServiceItem.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)
		params := mtoserviceitemops.CreateServiceRequestDocumentUploadParams{
			HTTPRequest:      req,
			File:             file,
			MtoServiceItemID: badFormatID.String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&mtoserviceitemops.CreateServiceRequestDocumentUploadUnprocessableEntity{}, response)

		// Validate outgoing payload
		// TODO: Can't validate the response because of the issue noted below. Figure out a way to
		//   either alter the service or relax the swagger requirements.
		// responsePayload := response.(*mtoserviceitemops.CreateServiceRequestDocumentUploadCreated).Payload
		// suite.NoError(responsePayload.Validate(strfmt.Default))
		// Handler is not setting any validation errors so InvalidFields won't be added to the payload.
	})

	suite.Run("create service request document upload fail - service item not found", func() {
		primeUser := factory.BuildUser(nil, nil, nil)
		handler, mtoServiceItem := setupTestData()

		badFormatID, _ := uuid.NewV4()
		file := suite.Fixture("test.pdf")

		req := httptest.NewRequest("POST", fmt.Sprintf("/mto_service_items/%s/uploads", mtoServiceItem.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)
		params := mtoserviceitemops.CreateServiceRequestDocumentUploadParams{
			HTTPRequest:      req,
			File:             file,
			MtoServiceItemID: badFormatID.String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&mtoserviceitemops.CreateServiceRequestDocumentUploadNotFound{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateServiceRequestDocumentUploadNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}
