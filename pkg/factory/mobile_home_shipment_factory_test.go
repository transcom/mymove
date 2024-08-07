package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildMobileHomeShipment() {
	suite.Run("Successful creation of default MobileHomeShipment", func() {
		defaultBoat := models.MobileHome{
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Mobile Home Make"),
			Model:          models.StringPointer("Mobile Home Model"),
			LengthInInches: models.IntPointer(300),
			WidthInInches:  models.IntPointer(108),
			HeightInInches: models.IntPointer(72),
			CreatedAt:      time.Now().Local(),
		}

		mobileHomeShipment := BuildMobileHomeShipment(suite.DB(), nil, nil)

		suite.Equal(defaultBoat.Year, mobileHomeShipment.Year)
		suite.Equal(defaultBoat.Make, mobileHomeShipment.Make)
		suite.Equal(defaultBoat.Model, mobileHomeShipment.Model)
		suite.Equal(defaultBoat.LengthInInches, mobileHomeShipment.LengthInInches)
		suite.Equal(defaultBoat.WidthInInches, mobileHomeShipment.WidthInInches)
		suite.Equal(defaultBoat.HeightInInches, mobileHomeShipment.HeightInInches)
		suite.Equal(defaultBoat.CreatedAt, mobileHomeShipment.CreatedAt)
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
