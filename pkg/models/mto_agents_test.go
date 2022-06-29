package models_test

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMTOAgentValidation() {
	suite.Run("test valid MTOAgent", func() {
		mtoShipmentID := uuid.Must(uuid.NewV4())
		mtoAgentID := uuid.Must(uuid.NewV4())

		validMTOAgent := models.MTOAgent{
			ID:            mtoAgentID,
			MTOShipmentID: mtoShipmentID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@testagent.agent"),
			Phone:         nil,
			MTOAgentType:  models.MTOAgentReleasing,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOAgent, expErrors)
	})
}
