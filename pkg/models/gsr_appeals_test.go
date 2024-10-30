package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGsrAppealValidation() {
	suite.Run("test valid GsrAppeal", func() {
		rejected := models.AppealStatusRejected
		validGsrAppeal := models.GsrAppeal{
			ID:                      uuid.Must(uuid.NewV4()),
			EvaluationReportID:      uuid.Must(uuid.NewV4()),
			ReportViolationID:       models.UUIDPointer(uuid.Must(uuid.NewV4())),
			OfficeUserID:            uuid.Must(uuid.NewV4()),
			IsSeriousIncidentAppeal: models.BoolPointer(true),
			AppealStatus:            rejected,
			Remarks:                 "Valid appeal remarks",
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
			DeletedAt:               nil,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validGsrAppeal, expErrors)
	})

	suite.Run("test missing required fields", func() {
		rejected := models.AppealStatusRejected
		invalidGsrAppeal := models.GsrAppeal{
			ID:           uuid.Must(uuid.NewV4()),
			AppealStatus: rejected,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		expErrors := map[string][]string{
			"office_user_id": {"OfficeUserID can not be blank."},
			"remarks":        {"Remarks can not be blank."},
		}

		suite.verifyValidationErrors(&invalidGsrAppeal, expErrors)
	})
}
