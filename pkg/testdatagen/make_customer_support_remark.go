package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

func MakeCustomerSupportRemark(db *pop.Connection, assertions Assertions) models.CustomerSupportRemark {
	move := assertions.Move
	officeUser := assertions.OfficeUser

	if isZeroUUID(assertions.CustomerSupportRemark.MoveID) {
		move = MakeMove(db, assertions)
	}

	if isZeroUUID(assertions.CustomerSupportRemark.OfficeUserID) {
		officeUser = MakeOfficeUser(db, assertions)
	}

	officeMoveRemark := models.CustomerSupportRemark{
		Content:      "This is an office remark.",
		OfficeUserID: officeUser.ID,
		MoveID:       move.ID,
	}

	// Overwrite with assertions
	mergeModels(&officeMoveRemark, assertions.CustomerSupportRemark)

	mustCreate(db, &officeMoveRemark, assertions.Stub)

	return officeMoveRemark
}

func MakeDefaultCustomerSupportRemark(db *pop.Connection) models.CustomerSupportRemark {
	officeMoveRemark := MakeCustomerSupportRemark(db, Assertions{})
	return officeMoveRemark
}
