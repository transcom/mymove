// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	m "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
)

func (suite *ModelSuite) Test_TransportationOfficeInstantiation() {
	office := &m.TransportationOffice{}
	expErrors := map[string][]string{
		"name":       {"Name can not be blank."},
		"address_id": {"AddressID can not be blank."},
	}
	suite.verifyValidationErrors(office, expErrors, nil)
}

func CreateTestShippingOffice(suite *ModelSuite) m.TransportationOffice {
	addressCreator := address.NewAddressCreator()
	newAddress := &m.Address{
		StreetAddress1: "123 washington Ave",
		City:           "ANCHORAGE",
		State:          "AK",
		PostalCode:     "99515"}
	newAddress, err := addressCreator.CreateAddress(suite.AppContextForTest(), newAddress)
	suite.NoError(err)

	office := m.TransportationOffice{
		Name:      "JPSO Supreme",
		AddressID: newAddress.ID,
		Gbloc:     "BMAF",
		Latitude:  61.1262383,
		Longitude: -149.9212882,
		Hours:     m.StringPointer("0900-1800 Mon-Sat"),
	}
	suite.MustSave(&office)
	return office
}

func (suite *ModelSuite) Test_BasicShippingOffice() {
	office := CreateTestShippingOffice(suite)
	var loadedOffice m.TransportationOffice
	suite.DB().Eager().Find(&loadedOffice, office.ID)
	suite.Equal(office.ID, loadedOffice.ID)
	suite.Equal(office.AddressID, loadedOffice.Address.ID)
}

func (suite *ModelSuite) Test_TransportationOffice() {
	jppso := CreateTestShippingOffice(suite)
	addressCreator := address.NewAddressCreator()
	ppoAddress := &m.Address{
		StreetAddress1: "456 Lincoln St",
		City:           "Sitka",
		State:          "AK",
		PostalCode:     "99835"}
	ppoAddress, err := addressCreator.CreateAddress(suite.AppContextForTest(), ppoAddress)
	suite.NoError(err)
	ppo := m.TransportationOffice{
		Name:             "Best PPO of the North",
		ShippingOfficeID: &jppso.ID,
		AddressID:        ppoAddress.ID,
		Gbloc:            "ACQR",
		Latitude:         57.0512403,
		Longitude:        -135.332707,
		Services:         m.StringPointer("Moose Shipping, Personal Goods"),
	}
	suite.MustSave(&ppo)
	var loadedOffice m.TransportationOffice
	suite.DB().Eager().Find(&loadedOffice, ppo.ID)
	suite.Equal(ppo.ID, loadedOffice.ID)
	suite.Equal(jppso.ID, loadedOffice.ShippingOffice.ID)
}

func (suite *ModelSuite) TestGetCounselingOffices() {
	customAddress1 := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: m.Address{
				PostalCode: "59801",
				City:       "MISSOULA",
			},
		},
	}, nil)
	factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: m.DutyLocation{
				ProvidesServicesCounseling: false,
			},
		},
		{
			Model: m.TransportationOffice{
				Name: "PPPO Holloman AFB - USAF",
			},
		},
		{
			Model:    customAddress1,
			LinkOnly: true,
			Type:     &factory.Addresses.DutyLocationAddress,
		},
	}, nil)

	// duty locations in KKFA with provides_services_counseling = true
	customAddress2 := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: m.Address{
				PostalCode: "59801",
				City:       "MISSOULA",
			},
		},
	}, nil)
	factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: m.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
		{
			Model: m.TransportationOffice{
				Name: "PPPO Hill AFB - USAF",
			},
		},
		{
			Model:    customAddress2,
			LinkOnly: true,
			Type:     &factory.Addresses.DutyLocationAddress,
		},
	}, nil)

	customAddress3 := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: m.Address{
				PostalCode: "59801",
				City:       "MISSOULA",
			},
		},
	}, nil)
	origDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: m.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
		{
			Model: m.TransportationOffice{
				Name:             "PPPO Travis AFB - USAF",
				Gbloc:            "KKFA",
				ProvidesCloseout: true,
			},
		},
		{
			Model:    customAddress3,
			LinkOnly: true,
			Type:     &factory.Addresses.DutyLocationAddress,
		},
	}, nil)

	// this one will not show in the return since it is not KKFA
	customAddress4 := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: m.Address{
				PostalCode: "20906",
				City:       "ASPEN HILL",
			},
		},
	}, nil)
	factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: m.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
		{
			Model: m.TransportationOffice{
				Name:             "PPPO Fort Meade - USA",
				Gbloc:            "BGCA",
				ProvidesCloseout: true,
			},
		},
		{
			Model:    customAddress4,
			LinkOnly: true,
			Type:     &factory.Addresses.DutyLocationAddress,
		},
	}, nil)

	armyAffliation := m.AffiliationARMY
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: m.ServiceMember{
				Affiliation: &armyAffliation,
			},
		},
	}, nil)
	offices, err := m.GetCounselingOffices(suite.DB(), origDutyLocation.ID, serviceMember.ID)
	suite.NoError(err)
	suite.Len(offices, 2)
}
