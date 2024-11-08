package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildUBAllowance() {
	suite.Run("Successful creation of default entitlement", func() {
		// Under test:      BuildUBAllowance
		// Mocked:          None
		// Set up:          Create an entitlement with no customizations or traits
		// Expected outcome:User should be created with default values

		// SETUP
		// Create a default entitlement to compare values
		defUBAllowance := models.UBAllowances{
			BranchOfService: string(models.AffiliationAIRFORCE),
			OrderPayGrade:   string(models.ServiceMemberGradeE1),
			OrdersType:      string(internalmessages.OrdersTypePERMANENTCHANGEOFSTATION),
			HasDependents:   true,
			AccompaniedTour: true,
			UBAllowance:     2000,
		}

		// FUNCTION UNDER TEST
		ubAllowance := BuildUBAllowance(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defUBAllowance.BranchOfService, ubAllowance.BranchOfService)
		suite.Equal(defUBAllowance.OrderPayGrade, ubAllowance.OrderPayGrade)
		suite.Equal(defUBAllowance.OrdersType, ubAllowance.OrdersType)
		suite.Equal(defUBAllowance.HasDependents, ubAllowance.HasDependents)
		suite.Equal(defUBAllowance.AccompaniedTour, ubAllowance.AccompaniedTour)
		suite.Equal(defUBAllowance.UBAllowance, ubAllowance.UBAllowance)
	})

	suite.Run("Successful creation of customized ubAllowance", func() {
		// Under test:      BuildUBAllowance
		// Mocked:          None
		// Set up:          Create ubAllowance with customization
		// Expected outcome:ubAllowance should customized values

		// SETUP
		// Create a default ubAllowance to compare values
		custUBAllowance := models.UBAllowances{
			BranchOfService: string(models.AffiliationARMY),
			OrderPayGrade:   string(models.ServiceMemberGradeE2),
			OrdersType:      string(internalmessages.OrdersTypeTEMPORARYDUTY),
			UBAllowance:     400,
		}

		// FUNCTION UNDER TEST
		ubAllowance := BuildUBAllowance(suite.DB(), []Customization{
			{Model: custUBAllowance},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(custUBAllowance.BranchOfService, ubAllowance.BranchOfService)
		suite.Equal(custUBAllowance.OrderPayGrade, ubAllowance.OrderPayGrade)
		suite.Equal(custUBAllowance.OrdersType, ubAllowance.OrdersType)
		suite.Equal(custUBAllowance.UBAllowance, ubAllowance.UBAllowance)
	})

	suite.Run("Successful return of linkOnly ubAllowance", func() {
		// Under test:       BuildUBAllowance
		// Set up:           Create an ubAllowance and pass in a linkOnly ubAllowance
		// Expected outcome: No new ubAllowance should be created.

		// Check num ubAllowances
		precount, err := suite.DB().Count(&models.UBAllowances{})
		suite.NoError(err)

		ubAllowance := BuildUBAllowance(suite.DB(), []Customization{
			{
				Model: models.UBAllowances{
					ID:              uuid.Must(uuid.NewV4()),
					BranchOfService: string(models.AffiliationARMY),
					OrderPayGrade:   string(models.ServiceMemberGradeE2),
					OrdersType:      string(internalmessages.OrdersTypeTEMPORARYDUTY),
					HasDependents:   false,
					AccompaniedTour: false,
					UBAllowance:     400,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.Entitlement{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(400, ubAllowance.UBAllowance)

	})
}
