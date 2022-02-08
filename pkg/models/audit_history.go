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
	ID              uuid.UUID  `json:"id" db:"id"`
	SchemaName      string     `json:"schema_name" db:"schema_name"`
	TableName       string     `json:"table_name" db:"table_name"`
	RelID           int64      `json:"rel_id" db:"relid"`
	ObjectID        *uuid.UUID `json:"object_id" db:"object_id"`
	SessionUserID   *uuid.UUID `json:"session_user_id" db:"session_userid"`
	TransactionID   *int64     `json:"transaction_id" db:"transaction_id"`
	ClientQuery     *string    `json:"client_query" db:"client_query"`
	Action          string     `json:"action" db:"action"`
	EventName       *string    `json:"event_name" db:"event_name"`
	OldData         *JSONMap   `json:"old_data" db:"old_data"`
	ChangedData     *JSONMap   `json:"changed_data" db:"changed_data"`
	StatementOnly   bool       `json:"statement_only" db:"statement_only"`
	ActionTstampTx  time.Time  `json:"action_tstamp_tx" db:"action_tstamp_tx"`
	ActionTstampStm time.Time  `json:"action_tstamp_stm" db:"action_tstamp_stm"`
	ActionTstampClk time.Time  `json:"action_tstamp_clk" db:"action_tstamp_clk"`
}

// AuditHistories is not required by pop and may be deleted
type AuditHistories []AuditHistory

func (jm JSONMap) Value() (driver.Value, error) {
	if jm == nil {
		return nil, nil
	}
	// ba, err := jm.MarshalJSON()
	// return string(ba), err
	return json.Marshal(jm)

	//return json.Marshal(a)
}

//func (a *Attrs) Scan(value interface{}) error {
func (jm *JSONMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &jm)
}

// MoveAuditHistory adapted from AuditHistory for ghc.yaml/GHC API payload
// This struct doesn't have a database table
// NO DATABASE TABLE
/*
type MoveAuditHistory struct {
	ID                           uuid.UUID                 `json:"id"`
	SchemaName                      string                 `json:"schema_name"`
	TableName						string					`json:"table_name"`
	RelID							int						`json:"rel_id"`
	ObjectID						uuid.UUID 				`json:"object_id"`
	SessionUserID						uuid.UUID 			`json:"session_user_id"`
	TransactionID 				uuid.UUID				`json:"transaction_id"`
	ClientQuery  				string					`json:"client_query"`
	Action 						string                 	`json:"action"`
	OldValues     MoveAuditHistoryItems                 `json:"old_values`
	ChangedValues MoveAuditHistoryItems                 `json:"changed_values`
	StatementOnly               bool       	             `json:"statement_only"`
	ActionTstampTx  time.Time `json:"action_tstamp_tx"`
	ActionTstampStm time.Time `json:"action_tstamp_stm"`
	ActionTstampClk time.Time `json:"action_tstamp_clk"`
}

// MoveAuditHistories slice of MoveAuditHistory
// NO DATABASE TABLE
type MoveAuditHistories []MoveAuditHistory


// MoveAuditHistoryItem adapted from AuditHistory.old_data/changed_data for ghc.yaml/GHC API payload
// This struct doesn't have a database table
// NO DATABASE TABLE
type MoveAuditHistoryItem struct {
	ColumnName  string
	ColumnValue string
}

// MoveAuditHistoryItems slice of AuditHistoryItem
// NO DATABASE TABLE
type MoveAuditHistoryItems []MoveAuditHistoryItem

*/
