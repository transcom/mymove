package primeapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	sitaddressupdateops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/sit_address_update"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	sitaddressupdate "github.com/transcom/mymove/pkg/services/sit_address_update"
)

func (suite *HandlerSuite) TestCreateSITAddressUpdateRequest() {
	mockPlanner := &routemocks.Planner{}
	mockedDistance := 55
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(mockedDistance, nil)

	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(mockedDistance, nil)
	moveRouter := moverouter.NewMoveRouter()
	addressCreator := address.NewAddressCreator()
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(mockPlanner, query.NewQueryBuilder(), moveRouter, mtoshipment.NewMTOShipmentFetcher(), addressCreator)
	sitAddressUpdateCreator := sitaddressupdate.NewSITAddressUpdateRequestCreator(mockPlanner, addressCreator, serviceItemUpdater, moveRouter)

	suite.Run("Success 201 - Create SIT address update request", func() {
		// Testcase:   sitExtension is created
		// Expected:   Success response 201
		handlerConfig := suite.HandlerConfig()
		handler := CreateSITAddressUpdateRequestHandler{
			handlerConfig,
			sitAddressUpdateCreator,
		}

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		req := httptest.NewRequest("POST", "/sit-address-updates", nil)

		contractorRemarks := "This is a contractor remark"
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		createParams := sitaddressupdateops.CreateSITAddressUpdateRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreateSITAddressUpdateRequest{
				ContractorRemarks: &contractorRemarks,
				MtoServiceItemID:  *handlers.FmtUUID(serviceItem.ID),
				NewAddress: &primemessages.Address{
					City:           &newAddress.City,
					Country:        newAddress.Country,
					PostalCode:     &newAddress.PostalCode,
					State:          &newAddress.State,
					StreetAddress1: &newAddress.StreetAddress1,
					StreetAddress2: newAddress.StreetAddress2,
					StreetAddress3: newAddress.StreetAddress3,
				},
			},
		}

		//Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		//Run handler
		handlerResponse := handler.Handle(createParams)

		//Check response type
		suite.IsType(&sitaddressupdateops.CreateSITAddressUpdateRequestCreated{}, handlerResponse)
		successResponse := handlerResponse.(*sitaddressupdateops.CreateSITAddressUpdateRequestCreated).Payload

		//validate outgoing payload
		suite.NoError(successResponse.Validate(strfmt.Default))

		//Check returned values
		suite.Equal(*createParams.Body.ContractorRemarks, *successResponse.ContractorRemarks)
		suite.Equal(serviceItem.ID.String(), successResponse.MtoServiceItemID.String())
		suite.Equal(models.SITAddressUpdateStatusRequested, successResponse.Status)
		suite.Equal(successResponse.Distance, successResponse.Distance)

		suite.NotNil(successResponse.ID)
		suite.NotNil(successResponse.NewAddressID)
		suite.NotNil(successResponse.NewAddress)
		suite.NotNil(successResponse.UpdatedAt)
		suite.NotNil(successResponse.CreatedAt)
		suite.NotNil(successResponse.ETag)
	})

	suite.Run("Returns 422 when attempting to update an unapproved service item", func() {
		handlerConfig := suite.HandlerConfig()
		handler := CreateSITAddressUpdateRequestHandler{
			handlerConfig,
			sitAddressUpdateCreator,
		}

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					// Not allowed to update on an unapproved service item
					Status: models.MTOServiceItemStatusRejected,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		req := httptest.NewRequest("POST", "/sit-address-updates", nil)

		contractorRemarks := "This is a contractor remark"
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		createParams := sitaddressupdateops.CreateSITAddressUpdateRequestParams{
			HTTPRequest: req,
			Body: &primemessages.CreateSITAddressUpdateRequest{
				ContractorRemarks: &contractorRemarks,
				MtoServiceItemID:  *handlers.FmtUUID(serviceItem.ID),
				NewAddress: &primemessages.Address{
					City:           &newAddress.City,
					Country:        newAddress.Country,
					PostalCode:     &newAddress.PostalCode,
					State:          &newAddress.State,
					StreetAddress1: &newAddress.StreetAddress1,
					StreetAddress2: newAddress.StreetAddress2,
					StreetAddress3: newAddress.StreetAddress3,
				},
			},
		}

		//Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		//Run handler
		handlerResponse := handler.Handle(createParams)

		suite.IsType(&sitaddressupdateops.CreateSITAddressUpdateRequestUnprocessableEntity{}, handlerResponse)
		failureResponse := handlerResponse.(*sitaddressupdateops.CreateSITAddressUpdateRequestUnprocessableEntity).Payload

		suite.NoError(failureResponse.Validate(strfmt.Default))
	})
}
