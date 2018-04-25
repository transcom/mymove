package models

import (
	"github.com/gobuffalo/pop"
)

// GetQueueMoves gets all queueMove models for a specific lifecycleState
func GetQueueMoves(db *pop.Connection, lifecycleState string) (QueueMoves, error) {
	var queueMoves QueueMoves

	query := `
		SELECT sm.id, sm.edipi, sm.rank, sm.first_name, sm.last_name
		FROM service_members AS sm
		LEFT JOIN moves ON moves.user_id = sm.user_id
		LEFT JOIN personally_procured_moves ON moves.id = personally_procured_moves.move_id
		WHERE
			moves.lifecycle_state = $1
	`

	err = db.RawQuery(query, lifecycleState).All(&queueMoves)
	return queueMoves, err
}
