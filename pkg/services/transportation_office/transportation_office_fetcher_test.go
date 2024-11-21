package transportationoffice

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TransportationOfficeServiceSuite struct {
	*testingsuite.PopTestSuite
	toFetcher services.TransportationOfficesFetcher
}

func TestTransportationOfficeServiceSuite(t *testing.T) {

	ts := &TransportationOfficeServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *TransportationOfficeServiceSuite) Test_SearchTransportationOffice() {

	transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "LRC Fort Knox",
				ProvidesCloseout: true,
			},
		},
	}, nil)
	office, err := FindTransportationOffice(suite.AppContextForTest(), "LRC Fort Knox", true)

	suite.NoError(err)
	suite.Equal(transportationOffice.Name, office[0].Name)
	suite.Equal(transportationOffice.Address.ID, office[0].Address.ID)
	suite.Equal(transportationOffice.Gbloc, office[0].Gbloc)

}

func (suite *TransportationOfficeServiceSuite) Test_SearchWithNoTransportationOffices() {

	office, err := FindTransportationOffice(suite.AppContextForTest(), "LRC Fort Knox", true)
	suite.NoError(err)
	suite.Len(office, 0)
}

func (suite *TransportationOfficeServiceSuite) Test_SortedTransportationOffices() {

	transportationOffice1 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "JPPSO",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	transportationOffice3 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "SO",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "PPSO",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	office, err := FindTransportationOffice(suite.AppContextForTest(), "JPPSO", true)

	suite.NoError(err)
	suite.Equal(transportationOffice1.Name, office[0].Name)
	suite.Equal(transportationOffice1.ProvidesCloseout, true)
	suite.Equal(transportationOffice2.Name, office[1].Name)
	suite.Equal(transportationOffice2.ProvidesCloseout, true)
	suite.Equal(transportationOffice3.Name, office[2].Name)
	suite.Equal(transportationOffice3.ProvidesCloseout, true)

}

func (suite *TransportationOfficeServiceSuite) Test_FindCounselingOffices() {
	// duty location in KKFA with provies services counseling false
	customAddress1 := models.Address{
		ID:         uuid.Must(uuid.NewV4()),
		PostalCode: "59801",
	}
	factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{Model: customAddress1, Type: &factory.Addresses.DutyLocationAddress},
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: false,
			},
		},
		{
			Model: models.TransportationOffice{
				Name: "PPPO Holloman AFB - USAF",
			},
		},
	}, nil)

	// duty location in KKFA with provides services counseling true
	customAddress2 := models.Address{
		ID:         uuid.Must(uuid.NewV4()),
		PostalCode: "59801",
	}
	factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{Model: customAddress2, Type: &factory.Addresses.DutyLocationAddress},
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
		{
			Model: models.TransportationOffice{
				Name: "PPPO Hill AFB - USAF",
			},
		},
	}, nil)

	// duty location in KKFA with provides services counseling true
	customAddress3 := models.Address{
		ID:         uuid.Must(uuid.NewV4()),
		PostalCode: "59801",
	}
	origDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{Model: customAddress3, Type: &factory.Addresses.DutyLocationAddress},
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
		{
			Model: models.TransportationOffice{
				Name:             "PPPO Travis AFB - USAF",
				Gbloc:            "KKFA",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	// duty location NOT in KKFA with provides services counseling true
	customAddress4 := models.Address{
		ID:         uuid.Must(uuid.NewV4()),
		PostalCode: "20906",
	}
	factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{Model: customAddress4, Type: &factory.Addresses.DutyLocationAddress},
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
		{
			Model: models.TransportationOffice{
				Name:             "PPPO Fort Meade - USA",
				Gbloc:            "BGCA",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	offices, err := findCounselingOffice(suite.AppContextForTest(), origDutyLocation.ID)

	suite.NoError(err)
	suite.Len(offices, 2)
	suite.Equal(offices[0].Name, "PPPO Hill AFB - USAF")
	suite.Equal(offices[1].Name, "PPPO Travis AFB - USAF")
}

func (suite *TransportationOfficeServiceSuite) Test_Oconus_AK_FindCounselingOffices() {
	testContractName := "Test_findOconusGblocDepartmentIndicator"
	testContractCode := "Test_findOconusGblocDepartmentIndicator_Code"
	testPostalCode := "99790"
	testPostalCode2 := "99701"
	testGbloc := "ABCD"
	testGbloc2 := "EFGH"
	testTransportationName := "TEST - PPO"
	testTransportationName2 := "TEST - PPO2"

	serviceAffiliations := []models.ServiceMemberAffiliation{models.AffiliationARMY,
		models.AffiliationNAVY, models.AffiliationMARINES, models.AffiliationAIRFORCE, models.AffiliationCOASTGUARD,
		models.AffiliationSPACEFORCE}

	setupServiceMember := func(serviceMemberAffiliation models.ServiceMemberAffiliation) models.ServiceMember {
		customServiceMember := models.ServiceMember{
			FirstName:          models.StringPointer("Gregory"),
			LastName:           models.StringPointer("Van der Heide"),
			Telephone:          models.StringPointer("999-999-9999"),
			SecondaryTelephone: models.StringPointer("123-555-9999"),
			PersonalEmail:      models.StringPointer("peyton@example.com"),
			Edipi:              models.StringPointer("1000011111"),
			Affiliation:        &serviceMemberAffiliation,
			Suffix:             models.StringPointer("Random suffix string"),
			PhoneIsPreferred:   models.BoolPointer(false),
			EmailIsPreferred:   models.BoolPointer(false),
		}

		customAddress := models.Address{
			StreetAddress1: "987 Another Street",
		}

		customUser := models.User{
			OktaEmail: "test_email@email.com",
		}

		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{Model: customServiceMember},
			{Model: customAddress},
			{Model: customUser},
		}, nil)

		return serviceMember
	}

	createContract := func(appCtx appcontext.AppContext, contractCode string, contractName string) (*models.ReContract, error) {
		// See if contract code already exists.
		exists, err := appCtx.DB().Where("code = ?", contractCode).Exists(&models.ReContract{})
		if err != nil {
			return nil, fmt.Errorf("could not determine if contract code [%s] existed: %w", contractCode, err)
		}
		if exists {
			return nil, fmt.Errorf("the provided contract code [%s] already exists", contractCode)
		}

		// Contract code is new; insert it.
		contract := models.ReContract{
			Code: contractCode,
			Name: contractName,
		}
		verrs, err := appCtx.DB().ValidateAndSave(&contract)
		if verrs.HasAny() {
			return nil, fmt.Errorf("validation errors when saving contract [%+v]: %w", contract, verrs)
		}
		if err != nil {
			return nil, fmt.Errorf("could not save contract [%+v]: %w", contract, err)
		}
		return &contract, nil
	}

	setupDataForOconusSearchCounselingOffice := func(contract models.ReContract, postalCode string, gbloc string, transportationName string) (models.ReRateArea, models.OconusRateArea, models.UsPostRegionCity, models.DutyLocation) {
		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := models.ReRateArea{
			ID:         uuid.Must(uuid.NewV4()),
			ContractID: contract.ID,
			IsOconus:   true,
			Code:       rateAreaCode,
			Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
			Contract:   contract,
		}
		verrs, err := suite.DB().ValidateAndCreate(&rateArea)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		us_country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), postalCode)
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		oconusRateArea := models.OconusRateArea{
			ID:                 uuid.Must(uuid.NewV4()),
			RateAreaId:         rateArea.ID,
			CountryId:          us_country.ID,
			UsPostRegionCityId: usprc.ID,
			Active:             true,
		}
		verrs, err = suite.DB().ValidateAndCreate(&oconusRateArea)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		address := models.Address{
			StreetAddress1:     "n/a",
			City:               "SomeCity",
			State:              "AK",
			PostalCode:         postalCode,
			County:             "SomeCounty",
			IsOconus:           models.BoolPointer(true),
			UsPostRegionCityId: &usprc.ID,
			CountryId:          models.UUIDPointer(us_country.ID),
		}
		suite.MustSave(&address)

		origDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					AddressID:                  address.ID,
					ProvidesServicesCounseling: true,
				},
			},
			{
				Model: models.TransportationOffice{
					Name:             transportationName,
					Gbloc:            gbloc,
					ProvidesCloseout: true,
				},
			},
		}, nil)
		suite.MustSave(&origDutyLocation)

		found_duty_location, _ := models.FetchDutyLocation(suite.DB(), origDutyLocation.ID)

		return rateArea, oconusRateArea, *usprc, found_duty_location
	}

	suite.Run("success - findOconusGblocDepartmentIndicator - returns default GLOC for departmentAffiliation if no specific departmentAffilation mapping is defined", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		const fairbanksAlaskaPostalCode = "99790"
		_, oconusRateArea, _, dutylocation := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testTransportationName)

		// setup department affiliation to GBLOC mappings
		expected_gbloc := "TEST-GBLOC"
		jppsoRegion := models.JppsoRegions{
			Code: expected_gbloc,
			Name: "TEST PPM",
		}
		suite.MustSave(&jppsoRegion)

		gblocAors := models.GblocAors{
			JppsoRegionID:    jppsoRegion.ID,
			OconusRateAreaID: oconusRateArea.ID,
			// DepartmentIndicator is nil,
		}
		suite.MustSave(&gblocAors)

		serviceAffiliations := []models.ServiceMemberAffiliation{models.AffiliationARMY,
			models.AffiliationNAVY, models.AffiliationMARINES, models.AffiliationAIRFORCE, models.AffiliationCOASTGUARD,
			models.AffiliationSPACEFORCE}

		// loop through and make sure all branches are using expected default GBLOC
		for _, affiliation := range serviceAffiliations {
			serviceMember := setupServiceMember(affiliation)
			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ServiceMemberID: serviceMember.ID,
			})
			suite.Nil(err)
			departmentIndictor, err := findOconusGblocDepartmentIndicator(appCtx, dutylocation)
			suite.NotNil(departmentIndictor)
			suite.Nil(err)
			suite.Nil(departmentIndictor.DepartmentIndicator)
			suite.Equal(expected_gbloc, departmentIndictor.Gbloc)
		}
	})

	suite.Run("success - findOconusGblocDepartmentIndicator - returns specific GLOC for departmentAffiliation when a specific departmentAffilation mapping is defined", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		_, oconusRateArea, _, dutylocation := setupDataForOconusSearchCounselingOffice(*contract, testPostalCode, testGbloc, testTransportationName)

		departmentIndicators := []models.DepartmentIndicator{models.DepartmentIndicatorARMY,
			models.DepartmentIndicatorARMYCORPSOFENGINEERS, models.DepartmentIndicatorCOASTGUARD,
			models.DepartmentIndicatorNAVYANDMARINES, models.DepartmentIndicatorAIRANDSPACEFORCE}

		expectedAffiliationToDepartmentIndicatorMap := make(map[string]string, 0)
		expectedAffiliationToDepartmentIndicatorMap[models.AffiliationARMY.String()] = models.DepartmentIndicatorARMY.String()
		expectedAffiliationToDepartmentIndicatorMap[models.AffiliationNAVY.String()] = models.DepartmentIndicatorNAVYANDMARINES.String()
		expectedAffiliationToDepartmentIndicatorMap[models.AffiliationMARINES.String()] = models.DepartmentIndicatorNAVYANDMARINES.String()
		expectedAffiliationToDepartmentIndicatorMap[models.AffiliationAIRFORCE.String()] = models.DepartmentIndicatorAIRANDSPACEFORCE.String()
		expectedAffiliationToDepartmentIndicatorMap[models.AffiliationCOASTGUARD.String()] = models.DepartmentIndicatorCOASTGUARD.String()
		expectedAffiliationToDepartmentIndicatorMap[models.AffiliationSPACEFORCE.String()] = models.DepartmentIndicatorAIRANDSPACEFORCE.String()

		// setup department affiliation to GBLOC mappings
		expected_gbloc := "TEST-GBLOC"
		jppsoRegion := models.JppsoRegions{
			Code: expected_gbloc,
			Name: "TEST PPM",
		}
		suite.MustSave(&jppsoRegion)

		defaultGblocAors := models.GblocAors{
			JppsoRegionID:    jppsoRegion.ID,
			OconusRateAreaID: oconusRateArea.ID,
			//DepartmentIndicator is nil,
		}
		suite.MustSave(&defaultGblocAors)

		// setup specific departmentAffiliation mapping for each branch
		for _, departmentIndicator := range departmentIndicators {
			gblocAors := models.GblocAors{
				JppsoRegionID:       jppsoRegion.ID,
				OconusRateAreaID:    oconusRateArea.ID,
				DepartmentIndicator: models.StringPointer(departmentIndicator.String()),
			}
			suite.MustSave(&gblocAors)
		}

		// loop through and make sure all branches are using it's own dedicated GBLOC and not default
		for _, affiliation := range serviceAffiliations {
			serviceMember := setupServiceMember(affiliation)
			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ServiceMemberID: serviceMember.ID,
			})
			suite.Nil(err)
			departmentIndictor, err := findOconusGblocDepartmentIndicator(appCtx, dutylocation)
			suite.NotNil(departmentIndictor)
			suite.Nil(err)
			suite.NotNil(departmentIndictor.DepartmentIndicator)
			if match, ok := expectedAffiliationToDepartmentIndicatorMap[affiliation.String()]; ok {
				// verify service member's affiliation matches on specific departmentIndicator mapping record
				suite.Equal(match, *departmentIndictor.DepartmentIndicator)
			} else {
				suite.Fail(fmt.Sprintf("key does not exist for %s", affiliation.String()))
			}
			suite.Equal(expected_gbloc, departmentIndictor.Gbloc)
		}
	})

	suite.Run("failure -- returns error when there are default and no department specific GBLOC", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		_, _, _, dutylocation := setupDataForOconusSearchCounselingOffice(*contract, testPostalCode, testGbloc, testTransportationName)

		// No specific departmentAffiliation mapping or default were created.
		// Expect an error response when nothing is found.
		serviceMember := setupServiceMember(models.AffiliationAIRFORCE)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: serviceMember.ID,
		})
		suite.Nil(err)
		departmentIndictor, err := findOconusGblocDepartmentIndicator(appCtx, dutylocation)
		suite.Nil(departmentIndictor)
		suite.NotNil(err)
	})

	suite.Run("failure - findOconusGblocDepartmentIndicator - returns error when find service member ID fails", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		_, _, _, dutylocation := setupDataForOconusSearchCounselingOffice(*contract, testPostalCode, testGbloc, testTransportationName)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			// create fake service member ID to raise NOT found error
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})

		suite.Nil(err)
		departmentIndictor, err := findOconusGblocDepartmentIndicator(appCtx, dutylocation)
		suite.Nil(departmentIndictor)
		suite.NotNil(err)
	})

	suite.Run("failure - not found duty location returns error", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})
		unknown_duty_location_id := uuid.Must(uuid.NewV4())
		offices, err := findCounselingOffice(appCtx, unknown_duty_location_id)
		suite.Nil(offices)
		suite.NotNil(err)
	})

	suite.Run("success - offices using default departmentIndicator mapping", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		_, oconusRateArea, _, dutylocation := setupDataForOconusSearchCounselingOffice(*contract, testPostalCode, testGbloc, testTransportationName)

		// setup department affiliation to GBLOC mappings
		jppsoRegion := models.JppsoRegions{
			Code: testGbloc,
			Name: "TEST PPM",
		}
		suite.MustSave(&jppsoRegion)

		gblocAors := models.GblocAors{
			JppsoRegionID:    jppsoRegion.ID,
			OconusRateAreaID: oconusRateArea.ID,
			// DepartmentIndicator is nil,
		}
		suite.MustSave(&gblocAors)

		postalCodeToGBLOC := models.PostalCodeToGBLOC{
			PostalCode: testPostalCode,
			GBLOC:      testGbloc,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		suite.MustSave(&postalCodeToGBLOC)

		serviceMember := setupServiceMember(models.AffiliationAIRFORCE)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: serviceMember.ID,
		})

		suite.Nil(err)
		offices, err := findCounselingOffice(appCtx, dutylocation.ID)
		suite.NotNil(offices)
		suite.Nil(err)
		suite.Equal(1, len(offices))
		suite.Equal(testTransportationName, offices[0].Name)

		// add another transportation office
		factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:             testTransportationName2,
					ProvidesCloseout: true,
					Gbloc:            testGbloc,
				},
			},
		}, nil)
		offices, err = findCounselingOffice(appCtx, dutylocation.ID)
		suite.NotNil(offices)
		suite.Nil(err)
		suite.Equal(2, len(offices))
	})

	suite.Run("success - returns correct office based on service affiliation -- simulate Zone 2 scenerio", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		_, oconusRateArea, _, dutylocation := setupDataForOconusSearchCounselingOffice(*contract, testPostalCode, testGbloc, testTransportationName)

		// ******************************************************************************
		// setup department affiliation to GBLOC mappings for AF/SF
		// ******************************************************************************
		jppsoRegion_AFSF := models.JppsoRegions{
			Code: testGbloc,
			Name: "TEST PPO",
		}
		suite.MustSave(&jppsoRegion_AFSF)

		gblocAors_AFSF := models.GblocAors{
			JppsoRegionID:       jppsoRegion_AFSF.ID,
			OconusRateAreaID:    oconusRateArea.ID,
			DepartmentIndicator: models.StringPointer(models.DepartmentIndicatorAIRANDSPACEFORCE.String()),
		}
		suite.MustSave(&gblocAors_AFSF)

		postalCodeToGBLOC_AFSF := models.PostalCodeToGBLOC{
			PostalCode: testPostalCode,
			GBLOC:      testGbloc,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		suite.MustSave(&postalCodeToGBLOC_AFSF)
		// ******************************************************************************

		// ******************************************************************************
		// setup department affiliation to GBLOC mappings for other branches NOT AF/SF
		// ******************************************************************************
		jppsoRegion_not_AFSF := models.JppsoRegions{
			Code: testGbloc2,
			Name: "TEST PPO 2",
		}
		suite.MustSave(&jppsoRegion_not_AFSF)

		gblocAors_not_AFSF := models.GblocAors{
			JppsoRegionID:    jppsoRegion_not_AFSF.ID,
			OconusRateAreaID: oconusRateArea.ID,
		}
		suite.MustSave(&gblocAors_not_AFSF)

		postalCodeToGBLOC_not_AFSF := models.PostalCodeToGBLOC{
			PostalCode: testPostalCode2,
			GBLOC:      testGbloc2,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		suite.MustSave(&postalCodeToGBLOC_not_AFSF)

		// add transportation office for other branches not AF/SF
		factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:             testTransportationName2,
					ProvidesCloseout: true,
					Gbloc:            testGbloc2,
				},
			},
		}, nil)
		// ******************************************************************************

		for _, affiliation := range serviceAffiliations {
			serviceMember := setupServiceMember(affiliation)
			if affiliation == models.AffiliationAIRFORCE || affiliation == models.AffiliationSPACEFORCE {
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ServiceMemberID: serviceMember.ID,
				})
				offices, err := findCounselingOffice(appCtx, dutylocation.ID)
				suite.NotNil(offices)
				suite.Nil(err)
				suite.Equal(1, len(offices))
				// verify expected office is for AF/SF and not for Navy ..etc..
				suite.Equal(testTransportationName, offices[0].Name)
				suite.NotEqual(testTransportationName2, offices[0].Name)
			} else {
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ServiceMemberID: serviceMember.ID,
				})
				offices, err := findCounselingOffice(appCtx, dutylocation.ID)
				suite.NotNil(offices)
				suite.Nil(err)
				suite.Equal(1, len(offices))
				// verify expected office is for Navy ..etc.. and not AF/SF
				suite.Equal(testTransportationName2, offices[0].Name)
				suite.NotEqual(testTransportationName, offices[0].Name)
			}
		}
	})

}

func (suite *TransportationOfficeServiceSuite) Test_GetTransportationOffice() {
	suite.toFetcher = NewTransportationOfficesFetcher()
	transportationOffice1 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "OFFICE ONE",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "OFFICE TWO",
				ProvidesCloseout: false,
			},
		},
	}, nil)

	office1t, err1t := suite.toFetcher.GetTransportationOffice(suite.AppContextForTest(), transportationOffice1.ID, true)
	office1f, err1f := suite.toFetcher.GetTransportationOffice(suite.AppContextForTest(), transportationOffice1.ID, false)

	_, err2t := suite.toFetcher.GetTransportationOffice(suite.AppContextForTest(), transportationOffice2.ID, true)
	office2f, err2f := suite.toFetcher.GetTransportationOffice(suite.AppContextForTest(), transportationOffice2.ID, false)

	suite.NoError(err1t)
	suite.NoError(err1f)
	// Should return an error since no office matches the ID and provides closeout
	suite.Error(err2t)
	suite.NoError(err2f)

	suite.Equal("OFFICE ONE", office1t.Name)
	suite.Equal("OFFICE ONE", office1f.Name)
	suite.Equal("OFFICE TWO", office2f.Name)
}
