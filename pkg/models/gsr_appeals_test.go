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

func (suite *ModelSuite) TestTotalDependentsCalculation() {
	suite.Run("calculates total dependents correctly when both fields are set", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   models.IntPointer(2),
			DependentsTwelveAndOver: models.IntPointer(3),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Equal(2, *fetchedEntitlement.DependentsUnderTwelve)
		suite.Equal(3, *fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(5, *fetchedEntitlement.TotalDependents) // sum of 2 + 3
	})
	suite.Run("calculates total dependents correctly when DependentsUnderTwelve is nil", func() {
		entitlement := models.Entitlement{
			DependentsTwelveAndOver: models.IntPointer(3),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Nil(fetchedEntitlement.DependentsUnderTwelve)
		suite.Equal(3, *fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(3, *fetchedEntitlement.TotalDependents) // sum of 0 + 3
	})
	suite.Run("calculates total dependents correctly when DependentsTwelveAndOver is nil", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve: models.IntPointer(2),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Equal(2, *fetchedEntitlement.DependentsUnderTwelve)
		suite.Nil(fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(2, *fetchedEntitlement.TotalDependents) // sum of 2 + 0
	})
	suite.Run("sets total dependents to nil when both fields are nil", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   nil,
			DependentsTwelveAndOver: nil,
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Nil(fetchedEntitlement.DependentsUnderTwelve)
		suite.Nil(fetchedEntitlement.DependentsTwelveAndOver)
		suite.Nil(fetchedEntitlement.TotalDependents) // NOT 0, NOT A SUM, nil + nil is NULL
	})
}
