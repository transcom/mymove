package entitlements

import (
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *EntitlementsServiceSuite) TestGetWeightAllotment() {
	suite.Run("If a weight allotment is fetched by grade, it should be returned", func() {
		fetcher := NewWeightAllotmentFetcher()

		pg := factory.BuildPayGrade(suite.DB(), nil, nil)
		hhgAllowance := factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
			{
				Model:    pg,
				LinkOnly: true,
			},
		}, nil)

		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), pg.Grade)

		suite.NoError(err)
		suite.Equal(hhgAllowance.TotalWeightSelf, allotment.TotalWeightSelf)
		suite.Equal(hhgAllowance.TotalWeightSelfPlusDependents, allotment.TotalWeightSelfPlusDependents)
		suite.Equal(hhgAllowance.ProGearWeight, allotment.ProGearWeight)
		suite.Equal(hhgAllowance.ProGearWeightSpouse, allotment.ProGearWeightSpouse)
	})

	suite.Run("If pay grade does not exist, return an error", func() {
		fetcher := NewWeightAllotmentFetcher()

		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), "X-1")
		suite.Error(err)
		suite.IsType(apperror.QueryError{}, err)
		suite.Nil(allotment)
	})
}
