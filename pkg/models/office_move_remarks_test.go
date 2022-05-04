package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestOfficeMoveRemarkCreation() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	suite.NotNil(move)

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	suite.NotNil(officeUser)

	suite.T().Run("test valid office remark", func(t *testing.T) {
		officeMoveRemarkContent := "This is a note that's saying something about the move."
		validOfficeMoveRemark := models.OfficeMoveRemark{
			Content:      officeMoveRemarkContent,
			OfficeUser:   officeUser,
			OfficeUserID: officeUser.ID,
			Move:         move,
			MoveID:       move.ID,
		}

		suite.MustSave(&validOfficeMoveRemark)
		suite.NotNil(validOfficeMoveRemark.ID)
		suite.NotEqual(uuid.Nil, validOfficeMoveRemark.ID)
		suite.Equal(move.ID, validOfficeMoveRemark.MoveID)
		suite.Equal(officeMoveRemarkContent, validOfficeMoveRemark.Content)
		suite.Equal(officeUser.ID, validOfficeMoveRemark.OfficeUserID)
	})
}
