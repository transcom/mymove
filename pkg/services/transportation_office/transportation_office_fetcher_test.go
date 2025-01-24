package transportationoffice

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
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
	office, err := FindTransportationOffice(suite.AppContextForTest(), "LRC Fort Knox", true, false)

	suite.NoError(err)
	suite.Equal(transportationOffice.Name, office[0].Name)
	suite.Equal(transportationOffice.Address.ID, office[0].Address.ID)
	suite.Equal(transportationOffice.Gbloc, office[0].Gbloc)

}

func (suite *TransportationOfficeServiceSuite) Test_SearchWithNoTransportationOffices() {

	office, err := FindTransportationOffice(suite.AppContextForTest(), "LRC Fort Knox", true, false)
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

	office, err := FindTransportationOffice(suite.AppContextForTest(), "JPPSO", true, false)

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

	setupDataForOconusSearchCounselingOffice := func(postalCode string, gbloc string) (models.ReRateArea, models.OconusRateArea, models.UsPostRegionCity, models.DutyLocation) {
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: models.ReRateArea{
				ContractID: contract.ID,
				IsOconus:   true,
				Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
				Contract:   contract,
			},
		})
		suite.NotNil(rateArea)

		us_country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), postalCode)
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		oconusRateArea, err := models.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
		suite.NotNil(oconusRateArea)
		suite.Nil(err)

		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &usprc.ID,
				},
			},
		}, nil)

		origDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					AddressID:                  address.ID,
					ProvidesServicesCounseling: true,
				},
			},
			{
				Model: models.TransportationOffice{
					Name:             "TEST - PPO",
					Gbloc:            gbloc,
					ProvidesCloseout: true,
				},
			},
		}, nil)
		suite.MustSave(&origDutyLocation)

		found_duty_location, _ := models.FetchDutyLocation(suite.DB(), origDutyLocation.ID)

		return rateArea, *oconusRateArea, *usprc, found_duty_location
	}

	suite.Run("success - findOconusGblocDepartmentIndicator - returns default GBLOC for departmentAffiliation if no specific departmentAffilation mapping is defined", func() {
		const fairbanksAlaskaPostalCode = "99790"
		_, oconusRateArea, _, dutylocation := setupDataForOconusSearchCounselingOffice(fairbanksAlaskaPostalCode, "JEAT")

		// setup department affiliation to GBLOC mappings
		jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "JEAT")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := models.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, models.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		serviceMember := setupServiceMember(models.AffiliationARMY)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: serviceMember.ID,
		})
		departmentIndictor, err := findOconusGblocDepartmentIndicator(appCtx, dutylocation)
		suite.NotNil(departmentIndictor)
		suite.Nil(err)
		suite.Nil(departmentIndictor.DepartmentIndicator)
		suite.Equal("JEAT", departmentIndictor.Gbloc)
	})

	suite.Run("success - findOconusGblocDepartmentIndicator - returns specific GBLOC for departmentAffiliation when a specific departmentAffilation mapping is defined -- simulate Zone 2 scenerio", func() {
		const fairbanksAlaskaPostalCode = "99790"
		_, oconusRateArea, _, dutylocation := setupDataForOconusSearchCounselingOffice(fairbanksAlaskaPostalCode, "MBFL")

		// setup department affiliation to GBLOC mappings
		jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "MBFL")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := models.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, models.DepartmentIndicatorAIRANDSPACEFORCE.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		// loop through and make sure all branches are using it's own dedicated GBLOC and not default
		serviceMember := setupServiceMember(models.AffiliationAIRFORCE)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: serviceMember.ID,
		})
		suite.Nil(err)
		departmentIndictor, err := findOconusGblocDepartmentIndicator(appCtx, dutylocation)
		suite.NotNil(departmentIndictor)
		suite.Nil(err)
		suite.NotNil(departmentIndictor.DepartmentIndicator)
		suite.Equal(models.DepartmentIndicatorAIRANDSPACEFORCE.String(), *departmentIndictor.DepartmentIndicator)
		suite.Equal("MBFL", departmentIndictor.Gbloc)
	})

	suite.Run("failure - findOconusGblocDepartmentIndicator - returns error when find service member ID fails", func() {
		_, _, _, dutylocation := setupDataForOconusSearchCounselingOffice("99714", "JEAT")

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			// create fake service member ID to raise NOT found error
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})

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
		_, oconusRateArea, _, dutylocation := setupDataForOconusSearchCounselingOffice("99619", "MAPS")

		// setup department affiliation to GBLOC mappings
		jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "MAPS")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := models.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, models.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		serviceMember := setupServiceMember(models.AffiliationAIRFORCE)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: serviceMember.ID,
		})

		suite.Nil(err)
		offices, err := findCounselingOffice(appCtx, dutylocation.ID)
		suite.NotNil(offices)
		suite.Nil(err)
		suite.Equal(1, len(offices))
		suite.Equal("TEST - PPO", offices[0].Name)

		// add another transportation office
		factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:             "TEST - PPO2",
					ProvidesCloseout: true,
					Gbloc:            "MAPS",
				},
			},
		}, nil)
		offices, err = findCounselingOffice(appCtx, dutylocation.ID)
		suite.NotNil(offices)
		suite.Nil(err)
		suite.Equal(2, len(offices))
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
