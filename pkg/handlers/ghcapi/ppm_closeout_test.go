package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	ppmcloseoutops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/models"
	paymentrequest "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	ppmcloseout "github.com/transcom/mymove/pkg/services/ppm_closeout"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestGetPPMCloseoutHandler() {

	setUpMockCloseoutFetcher := func(returnValues ...interface{}) services.PPMCloseoutFetcher {
		mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}

		mockPPMCloseoutFetcher.On("GetPPMCloseout",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(returnValues...)

		return mockPPMCloseoutFetcher
	}

	setUpHandler := func(ppmCloseoutFetcher services.PPMCloseoutFetcher) GetPPMCloseoutHandler {
		return GetPPMCloseoutHandler{
			suite.HandlerConfig(),
			ppmCloseoutFetcher,
		}
	}

	// Success integration test
	suite.Run("Successful fetch (integration) test", func() {
		// Create mock object for return from Handler
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		SITReimbursement := unit.Cents(100000)
		ActualWeight := unit.Pound(980)
		DDP := unit.Cents(33280)
		DOP := unit.Cents(16160)
		EstimatedWeight := unit.Pound(4000)
		GrossIncentive := unit.Cents(5000000)
		Miles := 2082
		PackPrice := unit.Cents(295800)
		ProGearWeightCustomer := unit.Pound(1978)
		ProGearWeightSpouse := unit.Pound(280)
		UnpackPrice := unit.Cents(23800)
		aoa := unit.Cents(50)
		remainingIncentive := GrossIncentive - aoa
		haulPrice := unit.Cents(2300)
		haulFSC := unit.Cents(23)
		gcc := unit.Cents(500)
		actualMoveDate := time.Now()
		plannedMoveDate := time.Now()
		ppmCloseoutObj := models.PPMCloseout{
			ID:                    &ppmShipment.ID,
			SITReimbursement:      &SITReimbursement,
			ActualMoveDate:        &actualMoveDate,
			ActualWeight:          &ActualWeight,
			AOA:                   &aoa,
			DDP:                   &DDP,
			DOP:                   &DOP,
			EstimatedWeight:       &EstimatedWeight,
			GCC:                   &gcc,
			GrossIncentive:        &GrossIncentive,
			HaulFSC:               &haulFSC,
			HaulPrice:             &haulPrice,
			Miles:                 &Miles,
			PackPrice:             &PackPrice,
			PlannedMoveDate:       &plannedMoveDate,
			ProGearWeightCustomer: &ProGearWeightCustomer,
			ProGearWeightSpouse:   &ProGearWeightSpouse,
			RemainingIncentive:    &remainingIncentive,
			UnpackPrice:           &UnpackPrice,
		}
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		fetcher := setUpMockCloseoutFetcher(&ppmCloseoutObj, nil)
		handler := setUpHandler(fetcher)
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/closeout", ppmShipment.ID.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := ppmcloseoutops.GetPPMCloseoutParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(ppmShipment.ID.String()),
		}

		response := handler.Handle(params)

		suite.IsType(&ppmcloseoutops.GetPPMCloseoutOK{}, response)
		payload := response.(*ppmcloseoutops.GetPPMCloseoutOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// 404 response
	suite.Run("404 response when the service returns not found", func() {
		uuidForShipment, _ := uuid.NewV4()
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		handlerConfig := suite.HandlerConfig()
		fetcher := ppmcloseout.NewPPMCloseoutFetcher(suite.HandlerConfig().DTODPlanner(), &paymentrequest.RequestPaymentHelper{}, &mocks.PPMEstimator{})
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/closeout", uuidForShipment.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := ppmcloseoutops.GetPPMCloseoutParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(uuidForShipment.String()),
		}

		handler := GetPPMCloseoutHandler{
			handlerConfig,
			fetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&ppmcloseoutops.GetPPMCloseoutNotFound{}, response)
		payload := response.(*ppmcloseoutops.GetPPMCloseoutNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
