package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPort() {
	suite.Run("test Port functions", func() {
		port := models.Port{
			ID:        uuid.Must(uuid.NewV4()),
			PortCode:  "PortCode",
			PortType:  "Both",
			PortName:  "PortName",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		suite.Equal(port.TableName(), "ports")
		suite.Equal(port.PortName, "PortName")
		suite.Equal(port.PortType.String(), "Both")
	})

}
