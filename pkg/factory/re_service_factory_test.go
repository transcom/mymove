package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildReService() {

	suite.Run("Successful creation of default reService (customer)", func() {
		// Under test:      BuildReService
		// Mocked:          None
		// Set up:          Create a ReService with no customizations or traits
		// Expected outcome:ReService should be created with default values
		defaultReServiceCode := models.ReServiceCode("STEST")
		reService := BuildReService(suite.DB(), nil, nil)
		suite.Equal(defaultReServiceCode, reService.Code)
	})

	suite.Run("Successful creation of reService with customization", func() {
		// Under test:      BuildReService
		// Set up:          Create a ReService with a customized email and no trait
		// Expected outcome:ReService should be created with email and inactive status
		customReService := models.ReService{
			ID:       uuid.Must(uuid.NewV4()),
			Code:     models.ReServiceCodeCS,
			Name:     "Counseling",
			Priority: 2,
		}
		reService := BuildReService(suite.DB(), []Customization{
			{
				Model: customReService,
			},
		}, nil)
		suite.Equal(customReService.ID, reService.ID)
		suite.Equal(customReService.Code, reService.Code)
		suite.Equal(customReService.Name, reService.Name)
		suite.Equal(customReService.Priority, reService.Priority)
	})

	suite.Run("Successful creation of stubbed reService", func() {
		// Under test:      BuildReService
		// Set up:          Create a customized reService, but don't pass in a db
		// Expected outcome:ReService should be created with email and active status
		//                  No reService should be created in database
		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := BuildReService(nil, []Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
		}, nil)

		suite.Equal(models.ReServiceCodeCS, reService.Code)
		// Count how many reServices are in the DB, no new reServices should have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}

func (suite *FactorySuite) TestBuildReServiceHelpers() {
	suite.Run("FetchOrBuildReServiceByCode - reService exists", func() {
		// Under test:      FetchOrBuildReServiceByCode
		// Set up:          Create a reService, then call FetchOrBuildReServiceByCode
		// Expected outcome:Existing ReService should be returned
		//                  No new reService should be created in database

		ServicesCounselorReService := BuildReService(suite.DB(), []Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
		},
			nil)

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeCS)
		suite.NoError(err)
		suite.Equal(ServicesCounselorReService.Code, reService.Code)
		suite.Equal(ServicesCounselorReService.ID, reService.ID)

		// Count how many reServices are in the DB, no new reServices should have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("FetchOrBuildReServiceByCode - reService does not exists", func() {
		// Under test:      FetchOrBuildReServiceByCode
		// Set up:          Call FetchOrBuildReServiceByCode with a non-existent reService
		// Expected outcome:new reService is created

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeCS)
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeCS, reService.Code)

		// Count how many reServices are in the DB, new reService should have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount+1, count)
	})

	suite.Run("FetchOrBuildReServiceByCode - stubbed reService", func() {
		// Under test:      FetchOrBuildReServiceByCode
		// Set up:          Call FetchOrBuildReServiceByCode without a db
		// Expected outcome:ReService is created but not saved to db

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchOrBuildReServiceByCode(nil, models.ReServiceCodeCS)
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeCS, reService.Code)

		// Count how many reServices are in the DB, no new reServices should have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("BuildReServiceByCode", func() {
		// Under test:      BuildDDFSITReService
		// Set up:          Call BuildDDFSITReService with ReServiceCodeCS
		// Expected outcome:ReService is created

		reService := BuildReServiceByCode(suite.DB(), models.ReServiceCodeCS)
		suite.Equal(models.ReServiceCodeCS, reService.Code)
	})

	suite.Run("BuildDDFSITReService", func() {
		// Under test:      BuildDDFSITReService
		// Set up:          Call BuildDDFSITReService
		// Expected outcome:DDFSIT reservice is returned. DDASIT and DDDSIT are also created

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := BuildDDFSITReService(suite.DB())
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeDDFSIT, reService.Code)

		// Count how many reServices are in the DB, 3 new reServices should have been created.
		var reServices []models.ReService
		var hasDDASIT, hasDDDSIT bool
		err = suite.DB().All(&reServices)
		suite.NoError(err)
		suite.Equal(precount+3, len(reServices))
		for _, service := range reServices {
			if service.Code == models.ReServiceCodeDDASIT {
				hasDDASIT = true
				continue
			}
			if service.Code == models.ReServiceCodeDDDSIT {
				hasDDDSIT = true
			}
		}
		suite.True(hasDDASIT)
		suite.True(hasDDDSIT)
	})

	suite.Run("BuildDOFSITReService", func() {
		// Under test:      BuildDOFSITReService
		// Set up:          Call BuildDOFSITReService
		// Expected outcome:DOFSIT reservice is returned. DOPSIT and DOASIT are also created

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := BuildDOFSITReService(suite.DB())
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeDOFSIT, reService.Code)

		// Count how many reServices are in the DB, 3 new reServices should have been created.
		var reServices []models.ReService
		var hasDOPSIT, hasDOASIT bool
		err = suite.DB().All(&reServices)
		suite.NoError(err)
		suite.Equal(precount+3, len(reServices))
		for _, service := range reServices {
			if service.Code == models.ReServiceCodeDOPSIT {
				hasDOPSIT = true
				continue
			}
			if service.Code == models.ReServiceCodeDOASIT {
				hasDOASIT = true
			}
		}
		suite.True(hasDOPSIT)
		suite.True(hasDOASIT)
	})
}
