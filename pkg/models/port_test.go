package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPortValidation() {
	suite.Run("test valid Port", func() {
		validPort := models.Port{
			ID:        uuid.Must(uuid.NewV4()),
			PortCode:  "1234",
			PortType:  models.PortTypeA,
			PortName:  "Valid port name",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPort, expErrors)
	})

	suite.Run("test missing required fields", func() {
		invalidPort := models.Port{
			ID:        uuid.Must(uuid.NewV4()),
			PortType:  models.PortTypeA,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		expErrors := map[string][]string{
			"port_code": {"PortCode can not be blank."},
			"port_name": {"PortName can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPort, expErrors)
	})

	suite.Run("test invalid port type", func() {
		invalidPort := models.Port{
			ID:        uuid.Must(uuid.NewV4()),
			PortType:  "I", //'I' for invalid
			PortCode:  "1234",
			PortName:  "Invalid port type",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expErrors := map[string][]string{
			"port_type": {"PortType is not in the list [A, S, B]."},
		}
		suite.verifyValidationErrors(&invalidPort, expErrors)
	})
}
