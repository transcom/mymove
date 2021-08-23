package movetaskorder

import (
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type createMoveTaskOrderQueryBuilder interface {
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
}

type moveTaskOrderCreator struct {
	builder createMoveTaskOrderQueryBuilder
}

// CreateMoveTaskOrder creates a move task order
func (o *moveTaskOrderCreator) CreateMoveTaskOrder(appCtx appcontext.AppContext, moveTaskOrder *models.Move) (*models.Move, *validate.Errors, error) {
	// generate reference id if empty
	if moveTaskOrder.ReferenceID == nil || strings.TrimSpace(*moveTaskOrder.ReferenceID) == "" {
		referenceID, err := models.GenerateReferenceID(appCtx.DB())
		if err != nil {
			return nil, nil, err
		}

		moveTaskOrder.ReferenceID = &referenceID
	}

	moveTaskOrder.Show = swag.Bool(true)

	// TODO: Remove this? Doesn't Pop automatically do this?
	moveTaskOrder.CreatedAt = time.Now()
	moveTaskOrder.UpdatedAt = time.Now()

	verrs, err := o.builder.CreateOne(appCtx, moveTaskOrder)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return moveTaskOrder, nil, nil
}

// NewMoveTaskOrderCreator returns an new creator
func NewMoveTaskOrderCreator(builder createMoveTaskOrderQueryBuilder) services.MoveTaskOrderCreator {
	return &moveTaskOrderCreator{builder}
}
