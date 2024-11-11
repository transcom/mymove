package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestFetchReService() {

	suite.Run("Successful creation of default reService (customer)", func() {
		// Under test:      FetchReService
		// Mocked:          None
		// Set up:          Create a ReService with no customizations or traits
		// Expected outcome:ReService should be created with default values
		defaultReServiceCode := models.ReServiceCode("DLH")
		reService := FetchReService(suite.DB(), nil, nil)
		suite.Equal(defaultReServiceCode, reService.Code)
	})

	suite.Run("Successful retrieval of ReService", func() {
		// Under test:      FetchReService
		// Set up:          Create a ReService with a customized email and no trait
		// Expected outcome:ReService should be returned
		customReService := models.ReService{
			Code: models.ReServiceCodeCS,
			Name: "Counseling",
		}
		reService := FetchReService(suite.DB(), []Customization{
			{
				Model: customReService,
			},
		}, nil)
		suite.Equal(models.ReServiceCodeCS, reService.Code)
		suite.Equal("Counseling", reService.Name)
	})

	suite.Run("Successful creation of stubbed reService", func() {
		// Under test:      FetchReService
		// Set up:          Create a customized reService, but don't pass in a db
		// Expected outcome:ReService should be created with email and active status
		//                  No reService should be created in database
		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchReService(nil, []Customization{
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

func (suite *FactorySuite) TestReServiceHelpers() {
	suite.Run("FetchReService - reService exists", func() {
		// Under test:      FetchReService
		// Set up:          Create a reService, then call FetchReService
		// Expected outcome:Existing ReService should be returned
		//                  No new reService should be created in database

		ServicesCounselorReService := FetchReService(suite.DB(), []Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
		},
			nil)

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchReService(suite.DB(), []Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
		}, nil)
		suite.NoError(err)
		suite.Equal(ServicesCounselorReService.Code, reService.Code)
		suite.Equal(ServicesCounselorReService.ID, reService.ID)

		// Count how many reServices are in the DB, no new reServices should have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("FetchReServiceByCode - reService exists", func() {
		// Under test:      FetchReServiceByCode
		// Set up:          Create a reService, then call FetchReServiceByCode
		// Expected outcome:Existing ReService should be returned
		//                  No new reService should be created in database

		ServicesCounselorReService := FetchReService(suite.DB(), []Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
		},
			nil)

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeCS)
		suite.NoError(err)
		suite.Equal(ServicesCounselorReService.Code, reService.Code)
		suite.Equal(ServicesCounselorReService.ID, reService.ID)

		// Count how many reServices are in the DB, no new reServices should have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("FetchReService - reService does not exists", func() {
		// Under test:      FetchReService
		// Set up:          Call FetchReService with a non-existent reService
		// Expected outcome:new reService is created

		customReService := models.ReService{
			Name: "custom name",
			Code: "Not a real service",
		}
		reService := FetchReService(suite.DB(), []Customization{
			{
				Model: customReService,
			},
		}, nil)

		suite.Equal(models.ReServiceCodeDLH, reService.Code)
		suite.Equal("Domestic linehaul", reService.Name)

		// find by ID
		customReServiceByID := models.ReService{
			ID: reService.ID,
		}
		reService = FetchReService(suite.DB(), []Customization{
			{
				Model: customReServiceByID,
			},
		}, nil)

		suite.Equal(reService.ID, reService.ID)
		suite.Equal(reService.Code, reService.Code)
		suite.Equal(reService.Name, reService.Name)
	})

	suite.Run("FetchReServiceByCode - reService does not exists", func() {
		// Under test:      FetchReServiceByCode
		// Set up:          Call FetchReServiceByCode with a non-existent reService
		// Expected outcome:new reService is NOT created

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchReServiceByCode(suite.DB(), "Not a real service")
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeDLH, reService.Code)

		// Count how many reServices are in the DB, new reService should NOT have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("FetchReServiceByCode - stubbed reService", func() {
		// Under test:      FetchReServiceByCode
		// Set up:          Call FetchReServiceByCode without a db
		// Expected outcome:ReService is created but not saved to db

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := FetchReServiceByCode(nil, models.ReServiceCodeCS)
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeCS, reService.Code)

		// Count how many reServices are in the DB, no new reServices should have been created.
		count, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("FetchReServiceByCode", func() {
		// Under test:      FetchReServiceByCode
		// Set up:          Call FetchReServiceByCode with ReServiceCodeCS
		// Expected outcome:ReService is returned

		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeCS)
		suite.Equal(models.ReServiceCodeCS, reService.Code)
	})

	suite.Run("FetchReServiceByCode", func() {
		// Under test:      FetchReServiceByCode
		// Set up:          Call FetchReServiceByCode
		// Expected outcome:DDFSIT reservice is returned. DDASIT, DDDSIT, and DDSFSC are also returned.

		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
		suite.Equal(models.ReServiceCodeDDFSIT, reService.Code)

		reService = FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDASIT)
		suite.Equal(models.ReServiceCodeDDASIT, reService.Code)

		reService = FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDDSIT)
		suite.Equal(models.ReServiceCodeDDDSIT, reService.Code)

		reService = FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDSFSC)
		suite.Equal(models.ReServiceCodeDDSFSC, reService.Code)
	})

	suite.Run("FetchReServiceByCode", func() {
		// Under test:      FetchReServiceByCode
		// Expected outcome:DOFSIT, DOASIT, DOPSIT, DOSFSC reservices are returned

		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)
		suite.Equal(models.ReServiceCodeDOFSIT, reService.Code)

		reService = FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOASIT)
		suite.Equal(models.ReServiceCodeDOASIT, reService.Code)

		reService = FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOPSIT)
		suite.Equal(models.ReServiceCodeDOPSIT, reService.Code)

		reService = FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOSFSC)
		suite.Equal(models.ReServiceCodeDOSFSC, reService.Code)
	})

	suite.Run("BuildIDFSITReService", func() {
		// Under test:      BuildIDFSITReService
		// Set up:          Call BuildIDFSITReService
		// Expected outcome:IDFSIT reservice is returned. IDASIT, IDDSIT, and IDSFSC are also created

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := BuildIDFSITReService(suite.DB())
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeIDFSIT, reService.Code)

		// Count how many reServices are in the DB, 3 new reServices should have been created.
		var reServices []models.ReService
		var hasIDASIT, hasIDDSIT, hasIDSFSC bool
		err = suite.DB().All(&reServices)
		suite.NoError(err)
		suite.Equal(precount+4, len(reServices))
		for _, service := range reServices {
			if service.Code == models.ReServiceCodeIDASIT {
				hasIDASIT = true
				continue
			}
			if service.Code == models.ReServiceCodeIDDSIT {
				hasIDDSIT = true
				continue
			}
			if service.Code == models.ReServiceCodeIDSFSC {
				hasIDSFSC = true
			}
		}
		suite.True(hasIDASIT)
		suite.True(hasIDDSIT)
		suite.True(hasIDSFSC)
	})

	suite.Run("BuildIOFSITReService", func() {
		// Under test:      BuildIOFSITReService
		// Set up:          Call BuildIOFSITReService
		// Expected outcome:IOFSIT reservice is returned. IOPSIT and IOASIT are also created

		precount, err := suite.DB().Count(&models.ReService{})
		suite.NoError(err)

		reService := BuildIOFSITReService(suite.DB())
		suite.NoError(err)

		suite.Equal(models.ReServiceCodeIOFSIT, reService.Code)

		// Count how many reServices are in the DB, 3 new reServices should have been created.
		var reServices []models.ReService
		var hasIOPSIT, hasIOASIT, hasIOSFSC bool
		err = suite.DB().All(&reServices)
		suite.NoError(err)
		suite.Equal(precount+4, len(reServices))
		for _, service := range reServices {
			if service.Code == models.ReServiceCodeIOPSIT {
				hasIOPSIT = true
				continue
			}
			if service.Code == models.ReServiceCodeIOASIT {
				hasIOASIT = true
				continue
			}
			if service.Code == models.ReServiceCodeIOSFSC {
				hasIOSFSC = true
			}
		}
		suite.True(hasIOPSIT)
		suite.True(hasIOASIT)
		suite.True(hasIOSFSC)
	})
}
