package adminapi

import (
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"

	officeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexOfficesHandler() {
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			ID: id,
		},
	}
	testdatagen.MakeTransportationOffice(suite.DB(), assertions)

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", "/offices", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := officeop.IndexOfficesParams{
			HTTPRequest: req,
		}

		handler := IndexOfficesHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&officeop.IndexOfficesOK{}, response)
		okResponse := response.(*officeop.IndexOfficesOK)
		suite.Len(okResponse.Payload, 0)
	})
}
