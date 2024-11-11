package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildUBAllowance() {
	suite.Run("Successful creation of default entitlement", func() {
		// Under test:      FetchOrBuildUBAllowance
		// Mocked:          None
		// Set up:          Create an entitlement with no customizations or traits
		// Expected outcome:User should be created with default values

		// SETUP
		// Create a default entitlement to compare values
		branch := models.AffiliationAIRFORCE
		grade := models.ServiceMemberGradeE1
		orderType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		hasDependents := true
		accompaniedTour := true
		ubAllowanceValue := 2000
		defUBAllowance := models.UBAllowances{
			BranchOfService: (*string)(&branch),
			OrderPayGrade:   (*string)(&grade),
			OrdersType:      (*string)(&orderType),
			HasDependents:   &hasDependents,
			AccompaniedTour: &accompaniedTour,
			UBAllowance:     &ubAllowanceValue,
		}

		// FUNCTION UNDER TEST
		ubAllowance := FetchOrBuildUBAllowance(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defUBAllowance.BranchOfService, ubAllowance.BranchOfService)
		suite.Equal(defUBAllowance.OrderPayGrade, ubAllowance.OrderPayGrade)
		suite.Equal(defUBAllowance.OrdersType, ubAllowance.OrdersType)
		suite.Equal(defUBAllowance.HasDependents, ubAllowance.HasDependents)
		suite.Equal(defUBAllowance.AccompaniedTour, ubAllowance.AccompaniedTour)
		suite.Equal(defUBAllowance.UBAllowance, ubAllowance.UBAllowance)
	})

	suite.Run("Successful creation of customized ubAllowance", func() {
		// Under test:      FetchOrBuildUBAllowance
		// Mocked:          None
		// Set up:          Create ubAllowance with customization
		// Expected outcome:ubAllowance should customized values

		// SETUP
		// Create a default ubAllowance to compare values
		branch := models.AffiliationARMY
		grade := models.ServiceMemberGradeE2
		orderType := internalmessages.OrdersTypeTEMPORARYDUTY
		ubAllowanceValue := 400
		custUBAllowance := models.UBAllowances{
			BranchOfService: (*string)(&branch),
			OrderPayGrade:   (*string)(&grade),
			OrdersType:      (*string)(&orderType),
			UBAllowance:     &ubAllowanceValue,
		}

		// FUNCTION UNDER TEST
		ubAllowance := FetchOrBuildUBAllowance(suite.DB(), []Customization{
			{Model: custUBAllowance},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(custUBAllowance.BranchOfService, ubAllowance.BranchOfService)
		suite.Equal(custUBAllowance.OrderPayGrade, ubAllowance.OrderPayGrade)
		suite.Equal(custUBAllowance.OrdersType, ubAllowance.OrdersType)
		suite.Equal(custUBAllowance.UBAllowance, ubAllowance.UBAllowance)
	})

	suite.Run("Successful return of linkOnly ubAllowance", func() {
		// Under test:       FetchOrBuildUBAllowance
		// Set up:           Create an ubAllowance and pass in a linkOnly ubAllowance
		// Expected outcome: No new ubAllowance should be created.

		// Check num ubAllowances
		precount, err := suite.DB().Count(&models.UBAllowances{})
		suite.NoError(err)

		branch := models.AffiliationARMY
		grade := models.ServiceMemberGradeE2
		orderType := internalmessages.OrdersTypeTEMPORARYDUTY
		hasDependents := true
		accompaniedTour := true
		ubAllowanceValue := 400
		ubAllowance := FetchOrBuildUBAllowance(suite.DB(), []Customization{
			{
				Model: models.UBAllowances{
					ID:              uuid.Must(uuid.NewV4()),
					BranchOfService: (*string)(&branch),
					OrderPayGrade:   (*string)(&grade),
					OrdersType:      (*string)(&orderType),
					HasDependents:   &hasDependents,
					AccompaniedTour: &accompaniedTour,
					UBAllowance:     &ubAllowanceValue,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.Entitlement{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(&ubAllowanceValue, ubAllowance.UBAllowance)

	})
}
