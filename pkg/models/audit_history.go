package models

// AuditHistory is a record from the audit_history table
type AuditHistory struct {
}

// AuditHistories is not required by pop and may be deleted
type AuditHistories []AuditHistory

// MoveAuditHistory adapted from AuditHistory for ghc.yaml/GHC API payload
// This struct doesn't have a database table
type MoveAuditHistory struct {
	OldValues     MoveAuditHistoryItems
	ChangedValues MoveAuditHistoryItems
}

// MoveAuditHistories slice of MoveAuditHistory
type MoveAuditHistories []MoveAuditHistory

// MoveAuditHistoryItem adapted from AuditHistory.old_data/changed_data for ghc.yaml/GHC API payload
// This struct doesn't have a database table
type MoveAuditHistoryItem struct {
	ColumnName  string
	ColumnValue string
}

// MoveAuditHistoryItems slice of AuditHistoryItem
type MoveAuditHistoryItems []MoveAuditHistoryItem
