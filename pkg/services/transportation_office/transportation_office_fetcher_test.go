package transportationoffice

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

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
	suite.toFetcher = NewTransportationOfficesFetcher()
	// duty location in KKFA with provides_services_counseling = false
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

	// duty locations in KKFA with provides_services_counseling = true
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

	// this one will not show in the return since it is not KKFA
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

	armyAffliation := models.AffiliationARMY
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &armyAffliation,
			},
		},
	}, nil)

	offices, err := suite.toFetcher.GetCounselingOffices(suite.AppContextForTest(), origDutyLocation.ID, serviceMember.ID)
	suite.NoError(err)
	suite.Len(*offices, 2)
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
