package payloads

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func TestMoveOrder(t *testing.T) {
	moveOrder := &models.MoveOrder{}
	MoveOrder(moveOrder)
}
