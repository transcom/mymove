package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

func MakeOfficeMoveRemark(db *pop.Connection, assertions Assertions) models.OfficeMoveRemark {
	move := assertions.Move
	officeUser := assertions.OfficeUser

	if isZeroUUID(assertions.OfficeMoveRemark.MoveID) {
		move = MakeMove(db, assertions)
	}

	if isZeroUUID(assertions.OfficeMoveRemark.OfficeUserID) {
		officeUser = MakeOfficeUser(db, assertions)
	}

	officeMoveRemark := models.OfficeMoveRemark{
		Content:      "This is an office remark.",
		OfficeUserID: officeUser.ID,
		MoveID:       move.ID,
	}

	// Overwrite with assertions
	mergeModels(&officeMoveRemark, assertions.OfficeMoveRemark)

	mustCreate(db, &officeMoveRemark, assertions.Stub)

	return officeMoveRemark
}

func MakeDefaultOfficeMoveRemark(db *pop.Connection) models.OfficeMoveRemark {
	officeMoveRemark := MakeOfficeMoveRemark(db, Assertions{})
	return officeMoveRemark
}
