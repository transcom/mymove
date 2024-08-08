package transportationoffice

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TransportationOfficeServiceSuite struct {
	*testingsuite.PopTestSuite
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
	factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "JPPSO",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	transportationOffice1 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "PPPO Holloman AFB - USAF",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "PPPO Hill AFB - USAF",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	customAddress := models.Address{
		ID:         uuid.Must(uuid.NewV4()),
		PostalCode: "59801",
	}
	destDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{Model: customAddress, Type: &factory.Addresses.DutyLocationAddress},
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
	}, nil)

	offices, err := FindCounselingOffice(suite.AppContextForTest(), destDutyLocation.ID)

	suite.NoError(err)

	// return should not include any JPPSO
	suite.Len(offices, 2)

	// return should be ordered by name asc
	suite.Equal(transportationOffice2.Name, offices[0].Name)
	suite.Equal(transportationOffice1.Name, offices[1].Name)
}
