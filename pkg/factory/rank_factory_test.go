package factory

import "github.com/transcom/mymove/pkg/models"

func (suite *FactorySuite) TestBuildRank() {

	suite.Run("Successful creation of stubbed rank", func() {
		// Under test:      BuildRank
		// Set up:          Create a customized rank, but don't pass in a db
		// Expected outcome:Rank should be created with default values
		//                  No rank should be created in database
		customRankAbbv := "CRA"

		rank := BuildRank(suite.DB(), []Customization{
			{
				Model: models.Rank{
					RankAbbv: customRankAbbv,
				},
			},
		}, nil)

		suite.Equal(customRankAbbv, rank.RankAbbv)
		suite.Equal(string(models.AffiliationAIRFORCE), rank.Affiliation)
		suite.Equal("Senior Airman", rank.RankName)
	})

}

func (suite *FactorySuite) TestBuildRankHelpers() {
	suite.Run("FetchOrBuildRankByPayGradeAndAffiliation - rank exists", func() {
		// Under test:      FetchOrBuildRankByPayGradeAndAffiliation
		// Set up:          Try to create a rank, then call FetchOrBuildRankByPayGradeAndAffiliation
		// Expected outcome:Existing Rank should be returned
		//                  No new rank should be created in database

		existingRank := FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), string(models.AffiliationAIRFORCE))

		precount, err := suite.DB().Count(&models.Rank{})
		suite.NoError(err)

		rank := FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), string(models.AffiliationAIRFORCE))
		suite.NoError(err)
		suite.Equal(existingRank.RankAbbv, rank.RankAbbv)
		suite.Equal(existingRank.Affiliation, rank.Affiliation)
		suite.Equal(existingRank.ID, rank.ID)

		// Count how many ranks are in the DB, no new ranks should have been created.
		count, err := suite.DB().Count(&models.Rank{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("FetchOrBuildRankByPayGradeAndAffiliation - stubbed rank", func() {
		// Under test:      FetchOrBuildRankByPayGradeAndAffiliation
		// Set up:          Call FetchOrBuildRankByPayGradeAndAffiliation without a db
		// Expected outcome:Rank is created but not saved to db

		precount, err := suite.DB().Count(&models.Rank{})
		suite.NoError(err)

		rank := FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), string(models.AffiliationARMY))
		suite.NoError(err)

		suite.Equal("CPL", rank.RankAbbv)
		suite.Equal(string(models.AffiliationARMY), rank.Affiliation)
		suite.Equal("Corporal", rank.RankName)

		// Count how many ranks are in the DB, no new ranks should have been created.
		count, err := suite.DB().Count(&models.Rank{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})
}
