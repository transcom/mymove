package entitlements

import (
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
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

func (suite *EntitlementsServiceSuite) TestGetAllWeightAllotments() {
	suite.Run("Successfully fetch all weight allotments", func() {
		fetcher := NewWeightAllotmentFetcher()

		// Make the default E-5
		e5 := factory.BuildPayGrade(suite.DB(), nil, nil)
		factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
			{
				Model:    e5,
				LinkOnly: true,
			},
		}, nil)

		// Make an E-6
		e6 := factory.BuildPayGrade(suite.DB(), []factory.Customization{
			{
				Model: models.PayGrade{
					Grade: "E-6",
				},
			},
		}, nil)
		factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
			{
				Model:    e6,
				LinkOnly: true,
			},
		}, nil)

		// Assert both can be fetched
		allotments, err := fetcher.GetAllWeightAllotments(suite.AppContextForTest())
		suite.NoError(err)
		suite.Len(allotments, 2)
	})
}
