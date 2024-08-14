package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildMobileHomeShipment() {
	suite.Run("Successful creation of default MobileHomeShipment", func() {
		defaultMobileHome := models.MobileHome{
			Make:           models.StringPointer("Mobile Home Make"),
			Model:          models.StringPointer("Mobile Home Model"),
			Year:           models.IntPointer(1996),
			LengthInInches: models.IntPointer(300),
			HeightInInches: models.IntPointer(72),
			WidthInInches:  models.IntPointer(108),
		}

		mobileHomeShipment := BuildMobileHomeShipment(suite.DB(), nil, nil)

		suite.Equal(defaultMobileHome.Make, mobileHomeShipment.Make)
		suite.Equal(defaultMobileHome.Model, mobileHomeShipment.Model)
		suite.Equal(defaultMobileHome.Year, mobileHomeShipment.Year)
		suite.Equal(defaultMobileHome.LengthInInches, mobileHomeShipment.LengthInInches)
		suite.Equal(defaultMobileHome.HeightInInches, mobileHomeShipment.HeightInInches)
		suite.Equal(defaultMobileHome.WidthInInches, mobileHomeShipment.WidthInInches)
	})

	suite.Run("Successful creation of customized mobileHomeShipment", func() {
		customMobileHome := models.MobileHome{
			ID:             uuid.Must(uuid.NewV4()),
			LengthInInches: models.IntPointer(3000),
		}

		mobileHomeShipment := BuildMobileHomeShipment(suite.DB(), []Customization{
			{Model: customMobileHome},
		}, nil)

		suite.Equal(customMobileHome.LengthInInches, mobileHomeShipment.LengthInInches)
	})
}
