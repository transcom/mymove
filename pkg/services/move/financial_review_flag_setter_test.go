package move

import (
	"errors"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MoveServiceSuite) TestFinancialReviewFlagSetter() {
	flagCreator := NewFinancialReviewFlagSetter()
	defaultFlagReason := "destination address is far from duty location"

	suite.Run("flag can be set", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(move.UpdatedAt)

		suite.Require().Equal(false, move.FinancialReviewFlag)
		suite.Require().Nil(move.FinancialReviewFlagSetAt)
		suite.Require().Nil(move.FinancialReviewRemarks)
		m, err := flagCreator.SetFinancialReviewFlag(suite.AppContextForTest(), move.ID, eTag, true, &defaultFlagReason)
		suite.NoError(suite.DB().Reload(&move))
		suite.Require().NotNil(m)
		suite.Require().NoError(err)
		suite.Require().Equal(true, move.FinancialReviewFlag)
		suite.Require().NotNil(move.FinancialReviewFlagSetAt)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)
	})

	suite.Run("Wrong moveID should result in error", func() {
		wrongUUID := uuid.Must(uuid.NewV4())

		_, err := flagCreator.SetFinancialReviewFlag(suite.AppContextForTest(), wrongUUID, "", true, &defaultFlagReason)
		suite.Error(err)
		suite.Require().True(errors.As(err, &apperror.NotFoundError{}))
	})

	suite.Run("Empty remarks param should result in error", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(move.UpdatedAt)

		_, err := flagCreator.SetFinancialReviewFlag(suite.AppContextForTest(), move.ID, eTag, true, models.StringPointer(""))
		suite.Error(err)
		suite.Require().True(errors.As(err, &apperror.InvalidInputError{}))
	})

	suite.Run("setting flag after it has already been set should have no effect", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(move.UpdatedAt)
		// Make sure move starts out as we expect it to
		suite.Require().False(move.FinancialReviewFlag)
		suite.Require().Nil(move.FinancialReviewFlagSetAt)
		suite.Require().Nil(move.FinancialReviewRemarks)

		// Set the flag
		_, err := flagCreator.SetFinancialReviewFlag(suite.AppContextForTest(), move.ID, eTag, true, &defaultFlagReason)
		suite.Require().NoError(err)
		suite.Require().NoError(suite.DB().Reload(&move))
		suite.Require().True(move.FinancialReviewFlag)
		suite.Require().NotNil(move.FinancialReviewFlagSetAt)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)

	})
	// If we set the flag to false, the timestamp and remarks fields should be nilled out
	suite.Run("when flag is set to false we nil out FinancialReviewFlagSetAt and FinancialReviewRemarks", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(move.UpdatedAt)

		suite.Require().Equal(false, move.FinancialReviewFlag)
		m, err := flagCreator.SetFinancialReviewFlag(suite.AppContextForTest(), move.ID, eTag, false, nil)
		suite.NoError(suite.DB().Reload(&move))
		suite.Require().NotNil(m)
		suite.Require().NoError(err)
		suite.Require().Equal(false, move.FinancialReviewFlag)
		suite.Nil(move.FinancialReviewFlagSetAt)
		suite.Nil(move.FinancialReviewRemarks)
	})
}
