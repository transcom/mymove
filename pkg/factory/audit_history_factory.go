package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
	ObjectID             *uuid.UUID      `json:"object_id" db:"object_id"`
	SessionUserID        *uuid.UUID      `json:"session_user_id" db:"session_userid"`
	SessionUserFirstName *string         `json:"session_user_first_name" db:"-"`
	SessionUserLastName  *string         `json:"session_user_last_name" db:"-"`
	SessionUserEmail     *string         `json:"session_user_email" db:"-"`
	SessionUserTelephone *string         `json:"session_user_telephone" db:"-"`
	Context              *models.JSONMap `json:"context" db:"-"`
	ContextID            *string         `json:"context_id" db:"-"`
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
func (t TestDataAuditHistory) TableName() string {
	return "audit_history"
}

// BuildAuditHistory creates an AuditHistory
// Notes:
//   - SessionUserID and ObjectID are only set if they're passed in via customizations
//   - This follows a slightly different pattern than other factories.
//   - If necessary caller should create Objects/Users and set ObjectID/SessionUserID in customizations,
//     rather than creating LinkOnly models
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildAuditHistory(db *pop.Connection, customs []Customization, traits []Trait) TestDataAuditHistory {
	customs = setupCustomizations(customs, traits)

	// Find AuditHistory Customization and extract the custom AuditHistory
	var cAuditHistory TestDataAuditHistory
	if result := findValidCustomization(customs, AuditHistory); result != nil {
		cAuditHistory = result.Model.(TestDataAuditHistory)
		if result.LinkOnly {
			return cAuditHistory
		}
	}

	fakeEventName := "setFinancialReviewFlag"
	fakeContext := models.JSONMap{}
	fakeContextID := "1234567"
	fakeOldData := models.JSONMap{"Locator": "asdf"}
	fakeChangedData := models.JSONMap{"Locator": "fdsa"}
	var fakeTransactionID int64 = 1234
	fakeClientQuery := "select 1;"

	// Create default AuditHistory
	auditHistory := TestDataAuditHistory{
		SchemaName:      "public",
		TableNameDB:     "moves",
		RelID:           16592,
		Context:         &fakeContext,
		ContextID:       &fakeContextID,
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

	// Overwrite default values with those from custom AuditHistory
	testdatagen.MergeModels(&auditHistory, cAuditHistory)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &auditHistory)
	}

	return auditHistory
}
