package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGsrAppealValidation() {
	suite.Run("test valid GsrAppeal", func() {
		validGsrAppeal := models.GsrAppeal{
			ID:                      uuid.Must(uuid.NewV4()),
			EvaluationReportID:      models.UUIDPointer(uuid.Must(uuid.NewV4())),
			ReportViolationID:       models.UUIDPointer(uuid.Must(uuid.NewV4())),
			OfficeUserID:            uuid.Must(uuid.NewV4()),
			IsSeriousIncidentAppeal: models.BoolPointer(true),
			AppealStatus:            models.AppealStatusSustained,
			Remarks:                 "Valid appeal remarks",
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
			DeletedAt:               nil,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validGsrAppeal, expErrors)
	})

	suite.Run("test missing required fields", func() {
		invalidGsrAppeal := models.GsrAppeal{
			ID:           uuid.Must(uuid.NewV4()),
			AppealStatus: models.AppealStatusRejected,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		expErrors := map[string][]string{
			"office_user_id": {"OfficeUserID can not be blank."},
			"remarks":        {"Remarks can not be blank."},
		}

		suite.verifyValidationErrors(&invalidGsrAppeal, expErrors)
	})

	suite.Run("test invalid appeal status", func() {
		invalidGsrAppeal := models.GsrAppeal{
			ID:                      uuid.Must(uuid.NewV4()),
			EvaluationReportID:      models.UUIDPointer(uuid.Must(uuid.NewV4())),
			ReportViolationID:       models.UUIDPointer(uuid.Must(uuid.NewV4())),
			OfficeUserID:            uuid.Must(uuid.NewV4()),
			IsSeriousIncidentAppeal: models.BoolPointer(true),
			AppealStatus:            "INVALID_STATUS", // Invalid status
			Remarks:                 "Invalid appeal status",
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
		}
		expErrors := map[string][]string{
			"appeal_status": {"AppealStatus is not in the list [SUSTAINED, REJECTED]."},
		}
		suite.verifyValidationErrors(&invalidGsrAppeal, expErrors)
	})
}
