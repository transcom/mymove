package movetaskorder

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoservicehelper "github.com/transcom/mymove/pkg/services/move_task_order/shared"
)

type createMoveTaskOrderQueryBuilder interface {
	CreateOne(model interface{}) (*validate.Errors, error)
}

type moveTaskOrderCreator struct {
	builder createMoveTaskOrderQueryBuilder
	db      *pop.Connection
}

// CreateMoveTaskOrder creates a move task order
func (o *moveTaskOrderCreator) CreateMoveTaskOrder(moveTaskOrder *models.MoveTaskOrder) (*models.MoveTaskOrder, *validate.Errors, error) {
	// generate reference id if empty
	if strings.TrimSpace(moveTaskOrder.ReferenceID) == "" {
		referenceID, err := mtoservicehelper.GenerateReferenceID(o.db)
		if err != nil {
			return nil, nil, err
		}

		moveTaskOrder.ReferenceID = referenceID
	}

	moveTaskOrder.CreatedAt = time.Now()
	moveTaskOrder.UpdatedAt = time.Now()

	verrs, err := o.builder.CreateOne(moveTaskOrder)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return moveTaskOrder, nil, nil
}

// NewMoveTaskOrderCreator returns an new creator
func NewMoveTaskOrderCreator(builder createMoveTaskOrderQueryBuilder, db *pop.Connection) services.MoveTaskOrderCreator {
	return &moveTaskOrderCreator{builder, db}
}
