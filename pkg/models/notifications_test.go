package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestNotificationValidations() {
	suite.Run("test valid Notification", func() {
		validNotification := models.Notification{
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validNotification, expErrors)
	})

	suite.Run("test empty Notification", func() {
		emptyNotification := models.Notification{}
		expErrors := map[string][]string{
			"service_member_id": {"ServiceMemberID can not be blank."},
		}
		suite.verifyValidationErrors(&emptyNotification, expErrors)
	})
}
