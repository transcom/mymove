package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestIntlTransitTimesModel() {
	suite.Run("test that FetchInternationalTransitTime returns the re_intl_transit_times given an originRateAreaId and a destinationRateAreaId", func() {
		originRateAreaId := uuid.FromStringOrNil("6e802149-7e46-4d7a-ab57-6c4df832085d")
		destinationRateAreaId := uuid.FromStringOrNil("c18e25f9-ec34-41ca-8c1b-05558c8d6364")

		fetchedIntlTransitTime, err := models.FetchInternationalTransitTime(suite.DB(), originRateAreaId, destinationRateAreaId)

		suite.Nil(err)
		suite.NotNil(fetchedIntlTransitTime)
		suite.Equal(originRateAreaId, fetchedIntlTransitTime.OriginRateAreaId)
		suite.Equal(destinationRateAreaId, fetchedIntlTransitTime.DestinationRateAreaId)
	})

	suite.Run("test that FetchInternationalTransitTime is bi-directional", func() {
		originRateAreaId := uuid.FromStringOrNil("6e802149-7e46-4d7a-ab57-6c4df832085d")
		destinationRateAreaId := uuid.FromStringOrNil("c18e25f9-ec34-41ca-8c1b-05558c8d6364")

		fetchedIntlTransitTime, err := models.FetchInternationalTransitTime(suite.DB(), destinationRateAreaId, originRateAreaId)

		suite.Nil(err)
		suite.NotNil(fetchedIntlTransitTime)
		suite.Equal(originRateAreaId, fetchedIntlTransitTime.OriginRateAreaId)
		suite.Equal(destinationRateAreaId, fetchedIntlTransitTime.DestinationRateAreaId)
	})

	suite.Run("test that FetchInternationalTransitTime errors when a transit time is not found", func() {
		originRateAreaId := uuid.FromStringOrNil("6e956185-7e46-4d7a-ab57-6c4df832085d")
		destinationRateAreaId := uuid.FromStringOrNil("c18e895f9-ec34-41ca-4c1b-05558c8d6364")

		fetchedIntlTransitTime, err := models.FetchInternationalTransitTime(suite.DB(), destinationRateAreaId, originRateAreaId)

		suite.Error(err)
		suite.Equal(uuid.Nil, fetchedIntlTransitTime.ID)
		suite.Equal("could not look up intl transit time", err.Error())
	})
}
