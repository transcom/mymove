package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReContractYearValidations() {
	suite.Run("test valid ReContractYear", func() {
		validReContractYear := models.ReContractYear{
			ContractID:           uuid.Must(uuid.NewV4()),
			Name:                 "Base Period Year 1",
			StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
			EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
			Escalation:           1.03,
			EscalationCompounded: 1.74,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReContractYear, expErrors, nil)
	})

	suite.Run("test empty ReContractYear", func() {
		emptyReContractYear := models.ReContractYear{}
		expErrors := map[string][]string{
			"contract_id":           {"ContractID can not be blank."},
			"name":                  {"Name can not be blank."},
			"start_date":            {"StartDate can not be blank."},
			"end_date":              {"EndDate can not be blank."},
			"escalation":            {"Escalation can not be blank.", "0.000000 is not greater than 0.000000."},
			"escalation_compounded": {"EscalationCompounded can not be blank.", "0.000000 is not greater than 0.000000."},
		}
		suite.verifyValidationErrors(&emptyReContractYear, expErrors, nil)
	})

	suite.Run("test end date after start date, negative escalation, negative escalation compounded for ReContractYear", func() {
		badDatesReContractYear := models.ReContractYear{
			ContractID:           uuid.Must(uuid.NewV4()),
			Name:                 "Base Period Year 2",
			StartDate:            time.Date(2021, time.September, 30, 0, 0, 0, 0, time.UTC),
			EndDate:              time.Date(2020, time.October, 1, 0, 0, 0, 0, time.UTC),
			Escalation:           -1,
			EscalationCompounded: -1.74,
		}
		expErrors := map[string][]string{
			"end_date":              {"EndDate must be after StartDate."},
			"escalation":            {"-1.000000 is not greater than 0.000000."},
			"escalation_compounded": {"-1.740000 is not greater than 0.000000."},
		}
		suite.verifyValidationErrors(&badDatesReContractYear, expErrors, nil)
	})
}

func (suite *ModelSuite) TestReContractYearModel() {
	suite.Run("test that FetchContractId returns the contractId given a requestedPickupDate", func() {
		validReContractYear := testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 "Base Period Year 1",
				StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
				Escalation:           1.03,
				EscalationCompounded: 1.74,
			},
		})

		requestedPickupDate := time.Date(2019, time.October, 25, 0, 0, 0, 0, time.UTC)
		contractYearId, err := models.FetchContractId(suite.DB(), requestedPickupDate)

		suite.Nil(err)
		suite.NotNil(contractYearId)
		suite.Equal(contractYearId, validReContractYear.ContractID)
	})

	suite.Run("test that FetchContractId returns error when no requestedPickupDate is given", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 "Base Period Year 1",
				StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
				Escalation:           1.03,
				EscalationCompounded: 1.74,
			},
		})

		var time time.Time
		contractYearId, err := models.FetchContractId(suite.DB(), time)

		suite.NotNil(err)
		suite.Contains(err.Error(), "error fetching contract ID - required parameters not provided")
		suite.Equal(contractYearId, uuid.Nil)
	})

	suite.Run("test that FetchContractId returns error when no contract is found for given requestedPickupDate", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 "Base Period Year 1",
				StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
				Escalation:           1.03,
				EscalationCompounded: 1.74,
			},
		})
		requestedPickupDate := time.Date(2019, time.September, 1, 0, 0, 0, 0, time.UTC)

		contractYearId, err := models.FetchContractId(suite.DB(), requestedPickupDate)

		suite.NotNil(err)
		suite.Contains(err.Error(), "error fetching contract year id")
		suite.Equal(contractYearId, uuid.Nil)
	})
}
