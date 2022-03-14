package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func MakeAuditHistory(db *pop.Connection, assertions Assertions) models.AuditHistory {

	genericUUID := uuid.Must(uuid.NewV4())
	movesTableUUID := assertions.Move.ID
	fakeEventName := "setFinancialReviewFlag"
	fakeOldData := models.JSONMap{"Locator": "asdf"}
	fakeChangedData := models.JSONMap{"Locator": "fdsa"}
	var fakeTransactionID int64 = 1234
	fakeClientQuery := "select 1;"

	history := models.AuditHistory{
		ID:              genericUUID,
		SchemaName:      "public",
		TableName:       "moves",
		RelID:           16592,
		ObjectID:        &movesTableUUID,
		SessionUserID:   &assertions.User.ID,
		TransactionID:   &fakeTransactionID,
		ClientQuery:     &fakeClientQuery,
		Action:          "I",
		EventName:       &fakeEventName,
		OldData:         &fakeOldData,
		ChangedData:     &fakeChangedData,
		ActionTstampTx:  time.Now(),
		ActionTstampStm: time.Now(),
		ActionTstampClk: time.Now(),
	}

	mergeModels(&history, assertions.AuditHistory)

	mustCreate(db, &history, assertions.Stub)

	return history

}
