package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildContractor() {
	suite.Run("Successful creation of default contractor", func() {
		// Under test:      BuildContractor
		// Mocked:          None
		// Set up:          Create a contractor with no customizations or traits
		// Expected outcome:Contractor should be created with default values

		// SETUP
		// Create a default contractor to compare values
		defContractor := models.Contractor{
			Name:           DefaultContractName,
			ContractNumber: DefaultContractNumber,
			Type:           DefaultContractType,
		}

		// FUNCTION UNDER TEST
		contractor := BuildContractor(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defContractor.Name, contractor.Name)
		suite.Equal(defContractor.ContractNumber, contractor.ContractNumber)
		suite.Equal(defContractor.Type, contractor.Type)
	})

	suite.Run("Successful creation of customized contractor", func() {
		// Under test:      BuildContractor
		// Mocked:          None
		// Set up:          Create a contractor with customization
		// Expected outcome:Contractor should be created with customized values

		// SETUP
		// Create a custom contractor to compare values
		custContractor := models.Contractor{
			Name:           "Custom Contract Name",
			ContractNumber: "11111-2222-3333-4444",
			Type:           "Super Prime",
		}

		// FUNCTION UNDER TEST
		contractor := BuildContractor(suite.DB(), []Customization{
			{Model: custContractor},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(custContractor.Name, contractor.Name)
		suite.Equal(custContractor.ContractNumber, contractor.ContractNumber)
		suite.Equal(custContractor.Type, contractor.Type)
	})

	suite.Run("Successful return of linkOnly contractor", func() {
		// Under test:      BuildContractor
		// Set up:          Create a contractor and pass in a linkOnly contractor
		// Expected outcome:No new contractor should be created

		// Check num contractors
		precount, err := suite.DB().Count(&models.Contractor{})
		suite.NoError(err)

		contractor := BuildContractor(suite.DB(), []Customization{
			{
				Model: models.Contractor{
					ID:   uuid.Must(uuid.NewV4()),
					Type: "Super Prime",
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.Contractor{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal("Super Prime", contractor.Type)
	})

	suite.Run("Two contractors should not be created", func() {
		// Under test:      FetchOrBuildDefaultContractor
		// Set up:          Create a contractor with no customized state or traits
		// Expected outcome:Only 1 contractor should be created
		count, potentialErr := suite.DB().Where(`contract_number=$1`, DefaultContractNumber).Count(&models.Contractor{})
		suite.NoError(potentialErr)
		suite.Zero(count)

		firstContractor := FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

		secondContractor := FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

		suite.Equal(firstContractor.ID, secondContractor.ID)

		existingContractorCount, err := suite.DB().Where(`contract_number=$1`, DefaultContractNumber).Count(&models.Contractor{})
		suite.NoError(err)
		suite.Equal(1, existingContractorCount)
	})
}
