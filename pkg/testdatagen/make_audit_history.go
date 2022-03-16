package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// TestDataAuditHistory for testing purposes only
// Had to separate this from AuditHistory struct because some calculated fields are not meant to be saved but still
// need to define a "db" map to fetch.
type TestDataAuditHistory struct {
	ID uuid.UUID `json:"id" db:"id"`
	// Database schema audited table for this event is in
	SchemaName string `json:"schema_name" db:"schema_name"`
	// name of database table that was changed
	TableNameDB string `json:"table_name" db:"table_name"`
	// relation OID. Table OID (object identifier). Changes with drop/create
	RelID int64 `json:"rel_id" db:"relid"`
	// id column for the tableName where the data was changed
	ObjectID             *uuid.UUID `json:"object_id" db:"object_id"`
	SessionUserID        *uuid.UUID `json:"session_user_id" db:"session_userid"`
	SessionUserFirstName *string    `json:"session_user_first_name" db:"-"`
	SessionUserLastName  *string    `json:"session_user_last_name" db:"-"`
	SessionUserEmail     *string    `json:"session_user_email" db:"-"`
	SessionUserTelephone *string    `json:"session_user_telephone" db:"-"`
	Context              *string    `json:"context" db:"-"`
	// Identifier of transaction that made the change. May wrap, but unique paired with action_tstamp_tx
	TransactionID *int64 `json:"transaction_id" db:"transaction_id"`
	// Record the text of the client query that triggered the audit event
	ClientQuery *string `json:"client_query" db:"client_query"`
	// Action type; I = insert, D = delete, U = update, T = truncate
	Action string `json:"action" db:"action"`
	// API endpoint name that was called to make the change
	EventName   *string         `json:"event_name" db:"event_name"`
	OldData     *models.JSONMap `json:"old_data" db:"old_data"`
	ChangedData *models.JSONMap `json:"changed_data" db:"changed_data"`
	// true if audit event is from an FOR EACH STATEMENT trigger, false for FOR EACH ROW'
	StatementOnly bool `json:"statement_only" db:"statement_only"`
	// Transaction start timestamp for tx in which audited event occurred
	ActionTstampTx time.Time `json:"action_tstamp_tx" db:"action_tstamp_tx"`
	// Statement start timestamp for tx in which audited event occurred
	ActionTstampStm time.Time `json:"action_tstamp_stm" db:"action_tstamp_stm"`
	// Wall clock time at which audited event's trigger call occurred
	ActionTstampClk time.Time `json:"action_tstamp_clk" db:"action_tstamp_clk"`
}

// TableName overrides the table name used by Pop.
func (t *TestDataAuditHistory) TableName() string {
	return "audit_history"
}

func MakeAuditHistory(db *pop.Connection, assertions Assertions) TestDataAuditHistory {

	genericUUID := uuid.Must(uuid.NewV4())
	movesTableUUID := assertions.Move.ID
	fakeEventName := "setFinancialReviewFlag"
	fakeOldData := models.JSONMap{"Locator": "asdf"}
	fakeChangedData := models.JSONMap{"Locator": "fdsa"}
	var fakeTransactionID int64 = 1234
	fakeClientQuery := "select 1;"

	history := TestDataAuditHistory{
		ID:              genericUUID,
		SchemaName:      "public",
		TableNameDB:     "moves",
		RelID:           16592,
		ObjectID:        &movesTableUUID,
		SessionUserID:   &assertions.User.ID,
		TransactionID:   &fakeTransactionID,
		ClientQuery:     &fakeClientQuery,
		Action:          "UPDATE",
		EventName:       &fakeEventName,
		OldData:         &fakeOldData,
		ChangedData:     &fakeChangedData,
		ActionTstampTx:  time.Now(),
		ActionTstampStm: time.Now(),
		ActionTstampClk: time.Now(),
	}

	mergeModels(&history, assertions.TestDataAuditHistory)

	mustCreate(db, &history, assertions.Stub)

	return history

}
