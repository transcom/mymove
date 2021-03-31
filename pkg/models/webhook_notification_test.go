package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestWebhookNotification() {
	now := time.Now()
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
	})
	paymentRequestID := paymentRequest.ID
	mtoID := paymentRequest.MoveTaskOrderID

	suite.T().Run("test full notification", func(t *testing.T) {
		// Normal notification with Object in payload
		trace := uuid.Must(uuid.NewV4())
		newNotification := models.WebhookNotification{
			EventKey:        "PaymentRequest.Update",
			TraceID:         &trace,
			MoveTaskOrderID: &mtoID,
			ObjectID:        &paymentRequestID,
			Payload:         "{\"msg\": \"This is the payload\"}",
			Status:          "PENDING",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&newNotification, expErrors)
		if verrs, err := suite.DB().ValidateAndCreate(&newNotification); err != nil || verrs.HasAny() {
			if verrs.HasAny() {
				fmt.Println(verrs)
			}
			fmt.Println(err)
			suite.T().Fatal("Didn't create notification in db.")
		}

		if newNotification.ID == uuid.Nil {
			t.Error("Didn't get an id back for Notification.")
		} else {
			fmt.Println("got id ", newNotification.ID.String())
		}

	})

	suite.T().Run("test simple notification", func(t *testing.T) {
		// Allowing for a simple message notification, with an eventkey and payload
		newNotification := models.WebhookNotification{
			EventKey: "PaymentRequest.Update",
			Payload:  "{\"msg\": \"This is the payload\"}",
			Status:   "SKIPPED",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&newNotification, expErrors)
	})

	suite.T().Run("test notification with validation errors", func(t *testing.T) {
		trace := uuid.Must(uuid.NewV4())
		newNotification := models.WebhookNotification{
			EventKey: "",
			TraceID:  &trace,
			Payload:  "",
			Status:   "NEW",
		}
		expErrors := map[string][]string{}
		expErrors["status"] = []string{"Status is not in the list [PENDING, SENT, SKIPPED, FAILING, FAILED]."}
		expErrors["event_key"] = []string{"Eventkey should be in Subject.Action format."}
		expErrors["payload"] = []string{"Payload can not be blank."}
		suite.verifyValidationErrors(&newNotification, expErrors)

	})
}
