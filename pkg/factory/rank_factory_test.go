package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildRank() {
	suite.Run("with no customization nor db generate a generic", func() {
		expected := models.Rank{
			Affiliation: string(models.DepartmentIndicatorARMY),
			RankName:    "Senior Airman",
			RankAbbv:    "SrA"}
		rank := FetchOrBuildRank(nil, nil, nil)

		suite.Equal(expected, rank)
	})
	suite.Run("with link only customization", func() {
		customs := []Customization{
			{
				Model: models.Rank{
					ID:       uuid.Must(uuid.NewV4()),
					RankName: "Linked Rank"},
				LinkOnly: true,
			},
		}
		rank := FetchOrBuildRank(nil, customs, nil)

		suite.Equal("Linked Rank", rank.RankName)
	})

	suite.Run("with valid database connection", func() {
		rank := FetchOrBuildRank(suite.DB(), nil, nil)

		suite.Equal(string(models.DepartmentIndicatorARMY), rank.Affiliation)
		suite.Equal("SrA", rank.RankAbbv)
		suite.Equal("Senior Airman", rank.RankName)
	})

	suite.Run("with database and missing pay grade", func() {
		rank := FetchOrBuildRank(suite.DB(), nil, nil)

		suite.Equal(string(models.DepartmentIndicatorARMY), rank.Affiliation)
		suite.Equal("SrA", rank.RankAbbv)
		suite.Equal("Senior Airman", rank.RankName)

	})

}
