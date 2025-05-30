package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// AuditHistoryJSONData
type JSONMap map[string]interface{}

// AuditHistory is a record from the audit_history table
type AuditHistory struct {
	ID uuid.UUID `json:"id" db:"id"`
	// Database schema audited table for this event is in
	SchemaName string `json:"schema_name" db:"schema_name"`
	// name of database table that was changed
	// Note: naming this "AuditedTable" instead of "TableName" in the struct so we can have a "TableName"
	// method as well on this struct to override the default Pop table name algorithm.  Go won't let
	// you have both.  Leaving the database column name and swagger field name the same for now.
	AuditedTable string `json:"table_name" db:"table_name"`
	// relation OID. Table OID (object identifier). Changes with drop/create
	RelID int64 `json:"rel_id" db:"relid"`
	// id column for the tableName where the data was changed
	ObjectID             *uuid.UUID `json:"object_id" db:"object_id"`
	SessionUserID        *uuid.UUID `json:"session_user_id" db:"session_userid"`
	SessionUserFirstName *string    `json:"session_user_first_name" db:"session_user_first_name"`
	SessionUserLastName  *string    `json:"session_user_last_name" db:"session_user_last_name"`
	SessionUserEmail     *string    `json:"session_user_email" db:"session_user_email"`
	SessionUserTelephone *string    `json:"session_user_telephone" db:"session_user_telephone"`
	Context              *string    `json:"context" db:"context"`
	ContextID            *string    `json:"context_id" db:"context_id"`
	// Identifier of transaction that made the change. May wrap, but unique paired with action_tstamp_tx
	TransactionID *int64 `json:"transaction_id" db:"transaction_id"`
	// Record the text of the client query that triggered the audit event
	ClientQuery *string `json:"client_query" db:"client_query"`
	// Action type; I = insert, D = delete, U = update, T = truncate
	Action string `json:"action" db:"action"`
	// API endpoint name that was called to make the change
	EventName   *string `json:"event_name" db:"event_name"`
	OldData     *string `json:"old_data" db:"old_data"`
	ChangedData *string `json:"changed_data" db:"changed_data"`
	// true if audit event is from an FOR EACH STATEMENT trigger, false for FOR EACH ROW'
	StatementOnly bool `json:"statement_only" db:"statement_only"`
	// Transaction start timestamp for tx in which audited event occurred
	ActionTstampTx time.Time `json:"action_tstamp_tx" db:"action_tstamp_tx"`
	// Statement start timestamp for tx in which audited event occurred
	ActionTstampStm time.Time `json:"action_tstamp_stm" db:"action_tstamp_stm"`
	// Wall clock time at which audited event's trigger call occurred
	ActionTstampClk time.Time  `json:"action_tstamp_clk" db:"action_tstamp_clk"`
	SeqNum          *int       `json:"seq_num" db:"seq_num" rw:"r"`
	MoveID          *uuid.UUID `json:"move_id" db:"move_id"`
	ShipmentID      *uuid.UUID `json:"shipment_id" db:"shipment_id"`
}

// TableName overrides the table name used by Pop.
func (a AuditHistory) TableName() string {
	return "audit_history"
}

type AuditHistories []AuditHistory

// Value returns a JSON value (from JSONMap to JSON string)
func (jm JSONMap) Value() (driver.Value, error) {
	if jm == nil {
		return nil, nil
	}

	return json.Marshal(jm)
}

// Scan reads a data type and update the JSONMap to represent the value read from JSON
func (jm *JSONMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &jm)
}
