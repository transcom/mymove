package payloads

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func TestMoveOrder(t *testing.T) {
	moveOrder := &models.Order{}
	MoveOrder(moveOrder)
}

// TestMove makes sure zero values/optional fields are handled
func TestMove(t *testing.T) {
	Move(&models.Move{})
}
