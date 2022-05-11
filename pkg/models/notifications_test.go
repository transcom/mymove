package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestNotificationValidations() {
	suite.T().Run("test valid Notification", func(t *testing.T) {
		validNotification := models.Notification{
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validNotification, expErrors)
	})

	suite.T().Run("test empty Notification", func(t *testing.T) {
		emptyNotification := models.Notification{}
		expErrors := map[string][]string{
			"service_member_id": {"ServiceMemberID can not be blank."},
		}
		suite.verifyValidationErrors(&emptyNotification, expErrors)
	})
}
