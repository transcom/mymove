package payloads

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func TestOrder(_ *testing.T) {
	order := &models.Order{}
	Order(order)
}

// TestMove makes sure zero values/optional fields are handled
func TestMove(_ *testing.T) {
	Move(&models.Move{})
}
