package models_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReContractYearValidations() {
	suite.T().Run("test valid ReContractYear", func(t *testing.T) {
		validReContractYear := models.ReContractYear{
			ContractID:           uuid.Must(uuid.NewV4()),
			Name:                 "Base Period Year 1",
			StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
			EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
			Escalation:           1.03,
			EscalationCompounded: 1.74,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReContractYear, expErrors)
	})

	suite.T().Run("test empty ReContractYear", func(t *testing.T) {
		emptyReContractYear := &models.ReContractYear{}
		expErrors := map[string][]string{
			"contract_id":           {"ContractID can not be blank."},
			"name":                  {"Name can not be blank."},
			"start_date":            {"StartDate can not be blank."},
			"end_date":              {"EndDate can not be blank."},
			"escalation":            {"Escalation can not be blank.", "0.000000 is not greater than 0.000000."},
			"escalation_compounded": {"EscalationCompounded can not be blank.", "0.000000 is not greater than 0.000000."},
		}
		suite.verifyValidationErrors(emptyReContractYear, expErrors)
	})

	suite.T().Run("test end date after start date, negative escalation, negative escalation compounded for ReContractYear", func(t *testing.T) {
		badDatesReContractYear := models.ReContractYear{
			ContractID:           uuid.Must(uuid.NewV4()),
			Name:                 "Base Period Year 2",
			StartDate:            time.Date(2021, time.September, 30, 0, 0, 0, 0, time.UTC),
			EndDate:              time.Date(2020, time.October, 1, 0, 0, 0, 0, time.UTC),
			Escalation:           -1,
			EscalationCompounded: -1.74,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		expErrors := map[string][]string{
			"end_date":              {"EndDate must be after StartDate."},
			"escalation":            {"-1.000000 is not greater than 0.000000."},
			"escalation_compounded": {"-1.740000 is not greater than 0.000000."},
		}
		suite.verifyValidationErrors(&badDatesReContractYear, expErrors)
	})
}
