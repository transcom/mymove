package officemoveremarks

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OfficeMoveRemarksSuite) setupTestData() models.OfficeMoveRemarks {

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())

	var officeMoveRemarks models.OfficeMoveRemarks
	for i := 0; i < 5; i++ {
		remark := testdatagen.MakeOfficeMoveRemark(suite.DB(),
			testdatagen.Assertions{
				OfficeMoveRemark: models.OfficeMoveRemark{
					Content:      fmt.Sprintln("This is remark number: %i", i),
					OfficeUserID: officeUser.ID,
					MoveID:       move.ID,
				}})
		officeMoveRemarks = append(officeMoveRemarks, remark)
	}
	return officeMoveRemarks
}

func (suite *OfficeMoveRemarksSuite) setupTestDataMultipleUsers() models.OfficeMoveRemarks {
	move := testdatagen.MakeDefaultMove(suite.DB())

	var officeUsers models.OfficeUsers
	var officeMoveRemarks models.OfficeMoveRemarks
	for i := 0; i < 3; i++ {
		officeUsers = append(officeUsers, testdatagen.MakeDefaultOfficeUser(suite.DB()))
		for x := 0; x < 2; x++ {
			remark := testdatagen.MakeOfficeMoveRemark(suite.DB(),
				testdatagen.Assertions{
					OfficeMoveRemark: models.OfficeMoveRemark{
						Content:      fmt.Sprintln("This is remark number: %i", i),
						OfficeUserID: officeUsers[i].ID,
						MoveID:       move.ID,
					}})
			officeMoveRemarks = append(officeMoveRemarks, remark)
		}
	}
	return officeMoveRemarks
}

func (suite *OfficeMoveRemarksSuite) TestOfficeRemarksListFetcher() {
	fetcher := NewOfficeMoveRemarksFetcher()

	suite.Run("Can fetch office move remarks successfully", func() {
		createdMoveRemarks := suite.setupTestData()
		officeMoveRemarks, err := fetcher.ListOfficeMoveRemarks(suite.AppContextForTest(), createdMoveRemarks[0].MoveID)
		suite.NoError(err)
		suite.NotNil(officeMoveRemarks)

		officeMoveRemarksValues := *officeMoveRemarks
		suite.Len(officeMoveRemarksValues, 5)
	})

	suite.Run("Can fetch office move remarks involving multiple users properly", func() {
		createdMoveRemarks := suite.setupTestDataMultipleUsers()
		officeMoveRemarks, err := fetcher.ListOfficeMoveRemarks(suite.AppContextForTest(), createdMoveRemarks[0].MoveID)
		suite.NoError(err)
		suite.NotNil(officeMoveRemarks)

		officeMoveRemarksValues := *officeMoveRemarks
		suite.Len(createdMoveRemarks, len(officeMoveRemarksValues))
	})

	suite.Run("Office move remarks aren't found", func() {
		_ = suite.setupTestData()
		randomUUID, _ := uuid.NewV4()
		_, err := fetcher.ListOfficeMoveRemarks(suite.AppContextForTest(), randomUUID)
		suite.Error(models.ErrFetchNotFound, err)
	})

}
