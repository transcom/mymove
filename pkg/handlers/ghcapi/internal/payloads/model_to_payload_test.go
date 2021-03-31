package payloads

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func TestOrder(t *testing.T) {
	order := &models.Order{}
	Order(order)
}

// TestMove makes sure zero values/optional fields are handled
func TestMove(t *testing.T) {
	Move(&models.Move{})
}
