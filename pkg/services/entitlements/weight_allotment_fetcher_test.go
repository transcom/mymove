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
		e5Allowance := factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
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
		e6Allowance := factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
			{
				Model:    e6,
				LinkOnly: true,
			},
		}, nil)

		// Assert both can be fetched
		allotments, err := fetcher.GetAllWeightAllotments(suite.AppContextForTest())
		suite.NoError(err)
		suite.Len(allotments, 2)

		// Check the first allotment (E-5)
		suite.Equal(e5Allowance.TotalWeightSelf, allotments[0].TotalWeightSelf)
		suite.Equal(e5Allowance.TotalWeightSelfPlusDependents, allotments[0].TotalWeightSelfPlusDependents)
		suite.Equal(e5Allowance.ProGearWeight, allotments[0].ProGearWeight)
		suite.Equal(e5Allowance.ProGearWeightSpouse, allotments[0].ProGearWeightSpouse)
		suite.Equal(e5.Grade, allotments[0].PayGrade.Grade)

		// Check the second allotment (E-6)
		suite.Equal(e6Allowance.TotalWeightSelf, allotments[1].TotalWeightSelf)
		suite.Equal(e6Allowance.TotalWeightSelfPlusDependents, allotments[1].TotalWeightSelfPlusDependents)
		suite.Equal(e6Allowance.ProGearWeight, allotments[1].ProGearWeight)
		suite.Equal(e6Allowance.ProGearWeightSpouse, allotments[1].ProGearWeightSpouse)
		suite.Equal(e6.Grade, allotments[1].PayGrade.Grade)
	})
}
