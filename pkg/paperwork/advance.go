package paperwork

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
)

// GenerateAdvancePaperwork generates the advance paperwork for a move.
// TODO We're going to need a logger in here most likely
func GenerateAdvancePaperwork(db *pop.Connection, moveID uuid.UUID) error {
	move, err := models.FetchMoveForAdvancePaperwork(db, moveID)
	if err != nil {
		return err
	}
	fmt.Println(move)
	return nil
}
