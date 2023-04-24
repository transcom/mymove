package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestValidationRules() {

	suite.Run("checkMTOIDField", func() {

		suite.Run("success", func() {

			move := factory.BuildMove(suite.DB(), nil, nil)
			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: move.ID,
			}

			err := checkMTOIDField().Validate(suite.AppContextForTest(), paymentRequest, nil)
			suite.NoError(err)
		})

		suite.Run("failure", func() {

			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")),
			}

			err := checkMTOIDField().Validate(suite.AppContextForTest(), paymentRequest, nil)
			switch err.(type) {
			case apperror.InvalidCreateInputError:
				suite.Equal(err.Error(), "Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create")
			default:
				suite.Failf("expected *apperror.InvalidCreateInputError", "%v", err)
			}
		})

	})

}
