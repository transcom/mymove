package models

import (
	"github.com/gobuffalo/pop"
)

// GetMoveQueueItems gets all moveQueueItems for a specific lifecycleState
func GetMoveQueueItems(db *pop.Connection, lifecycleState string) (MoveQueueItems, error) {
	var moveQueueItems MoveQueueItems
	query := `
		SELECT moves.*, sm.id AS service_member_id, sm.edipi, sm.rank, sm.first_name, sm.last_name, ppm.id AS ppm_id
		FROM moves
		JOIN service_members AS sm ON moves.user_id = sm.user_id
		JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
		WHERE
			moves.lifecycle_state = $1
	`

	err = db.RawQuery(query, lifecycleState).All(&moveQueueItems)
	return moveQueueItems, err
}
