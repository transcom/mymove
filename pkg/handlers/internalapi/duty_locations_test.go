package internalapi

import (
	"context"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	locationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_locations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *HandlerSuite) TestSearchDutyLocationHandler() {
	t := suite.T()

	// Need a logged in user
	lgu := uuid.Must(uuid.NewV4()).String()
	user := models.User{
		OktaID:    lgu,
		OktaEmail: "email@example.com",
	}
	suite.MustSave(&user)

	newAKAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "AK",
		PostalCode:     "12345",
		County:         models.StringPointer("County"),
	}

	newHIAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "HI",
		PostalCode:     "12345",
		County:         models.StringPointer("County"),
	}
	factory.FetchOrBuildCountry(suite.AppContextForTest().DB(), nil, nil)
	addressCreator := address.NewAddressCreator()
	createdAKAddress, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newAKAddress)
	suite.NoError(err)

	createdHIAddress, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newHIAddress)
	suite.NoError(err)

	dutylocationAK := models.DutyLocation{
		Name:        "HELLOWORLD 1",
		AddressID:   createdAKAddress.ID,
		Affiliation: internalmessages.NewAffiliation(internalmessages.AffiliationAIRFORCE),
	}
	suite.MustSave(&dutylocationAK)

	dutylocationHI := models.DutyLocation{
		Name:        "HELLOWORLD 2",
		AddressID:   createdHIAddress.ID,
		Affiliation: internalmessages.NewAffiliation(internalmessages.AffiliationAIRFORCE),
	}
	suite.MustSave(&dutylocationHI)

	setupTestHandler := func(isAlaskaEnabled bool, isSimulateFeatureFlagError bool) SearchDutyLocationsHandler {
		handlerConfig := suite.HandlerConfig()
		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}

		mockGetFlagFunc := func(_ context.Context, _ *zap.Logger, entityID string, key string, _ map[string]string, mockVariant string) (services.FeatureFlag, error) {
			return services.FeatureFlag{
				Entity:    entityID,
				Key:       key,
				Match:     isAlaskaEnabled,
				Variant:   mockVariant,
				Namespace: "test",
			}, nil
		}
		if isSimulateFeatureFlagError {
			mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
				mock.Anything,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(services.FeatureFlag{}, errors.New("Some error"))
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		} else {
			mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
				mock.Anything,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(func(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {
				return mockGetFlagFunc(ctx, appCtx.Logger(), "user@example.com", key, flagContext, "")
			})
		}
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		handler := SearchDutyLocationsHandler{handlerConfig}
		return handler
	}

	req := httptest.NewRequest("GET", "/duty_locations", nil)
	// Make sure the context contains the auth values
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, session)
	newSearchParams := locationop.SearchDutyLocationsParams{
		HTTPRequest: req.WithContext(ctx),
		Search:      "helloworld",
	}

	//////////////////////////////////////////////////////////////
	// test when alaska is enabled
	//////////////////////////////////////////////////////////////
	var handlerAlaska = setupTestHandler(true, false)

	var responseAlaska = handlerAlaska.Handle(newSearchParams)

	// Assert we got back the 201 response
	var searchResponseAlaska = responseAlaska.(*locationop.SearchDutyLocationsOK)
	var locationPayloadsAlaska = searchResponseAlaska.Payload

	suite.NoError(locationPayloadsAlaska.Validate(strfmt.Default))

	if len(locationPayloadsAlaska) != 2 {
		t.Errorf("Should have 2 responses, got %v", len(locationPayloadsAlaska))
	}

	//////////////////////////////////////////////////////////////
	// test when alaska is not enabled
	//////////////////////////////////////////////////////////////
	handlerAlaska = setupTestHandler(false, false)

	responseAlaska = handlerAlaska.Handle(newSearchParams)

	searchResponseAlaska = responseAlaska.(*locationop.SearchDutyLocationsOK)
	locationPayloadsAlaska = searchResponseAlaska.Payload

	suite.NoError(locationPayloadsAlaska.Validate(strfmt.Default))

	// should return zero matches
	if len(locationPayloadsAlaska) != 0 {
		t.Errorf("Should have 0 responses, got %v", len(locationPayloadsAlaska))
	}

	//////////////////////////////////////////////////////////////
	// test when FF retrieval throws an error
	//////////////////////////////////////////////////////////////
	handlerAlaska = setupTestHandler(true, true)

	responseAlaska = handlerAlaska.Handle(newSearchParams)

	// Assert we got back the 201 response
	searchResponseAlaska = responseAlaska.(*locationop.SearchDutyLocationsOK)
	locationPayloadsAlaska = searchResponseAlaska.Payload

	suite.NoError(locationPayloadsAlaska.Validate(strfmt.Default))

	// simulating FeatureFlagFetcher().GetBooleanFlagForUser
	// throws error, defaults to FALSE returning 0 matches.
	if len(locationPayloadsAlaska) != 0 {
		t.Errorf("Should have 0 responses because error sets flag to false, got %v", len(locationPayloadsAlaska))
	}
}
