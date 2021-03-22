package payloads

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func TestOrder(t *testing.T) {
	order := &models.Order{}
	Order(order)
}
