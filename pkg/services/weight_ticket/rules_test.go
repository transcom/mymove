package weightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *WeightTicketSuite) TestCheckID() {
	suite.Run("Success", func() {
		suite.Run("Create WeightTicket without an ID", func() {
			err := checkID().Validate(suite.AppContextForTest(), nil, nil)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("Update WeightTicket with matching ID", func() {
			id := uuid.Must(uuid.NewV4())

			err := checkID().Validate(suite.AppContextForTest(), &models.WeightTicket{ID: id}, &models.WeightTicket{ID: id})

			suite.NilOrNoVerrs(err)
		})
	})

	suite.Run("Failure", func() {
		suite.Run("Return an error if the IDs don't match", func() {
			err := checkID().Validate(suite.AppContextForTest(), &models.WeightTicket{ID: uuid.Must(uuid.NewV4())}, &models.WeightTicket{ID: uuid.Must(uuid.NewV4())})

			suite.Error(err)
			suite.IsType(&validate.Errors{}, err)
			suite.Contains(err.Error(), "new WeightTicket ID must match original WeightTicket ID")
		})
	})
}
