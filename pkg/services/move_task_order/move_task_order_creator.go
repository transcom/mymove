package movetaskorder

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

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

	// create default service items as well
	err = o.createDefaultServiceItems(moveTaskOrder)
	if err != nil {
		return nil, nil, err
	}

	return moveTaskOrder, nil, nil
}

// NewMoveTaskOrderCreator returns an new creator
func NewMoveTaskOrderCreator(builder createMoveTaskOrderQueryBuilder, db *pop.Connection) services.MoveTaskOrderCreator {
	return &moveTaskOrderCreator{builder, db}
}

func (o *moveTaskOrderCreator) createDefaultServiceItems(moveTaskOrder *models.MoveTaskOrder) error {
	var reServices []models.ReService
	err := o.db.Where("code in (?)", []string{"MS", "CS"}).All(&reServices)

	if err != nil {
		return err
	}

	defaultServiceItems := make(map[uuid.UUID]models.MTOServiceItem)
	for _, reService := range reServices {
		defaultServiceItems[reService.ID] = models.MTOServiceItem{
			ReServiceID:     reService.ID,
			MoveTaskOrderID: moveTaskOrder.ID,
			Status:          models.MTOServiceItemStatusSubmitted,
		}
	}

	// Remove the ones that exist on the mto
	for _, item := range moveTaskOrder.MTOServiceItems {
		for _, reService := range reServices {
			if item.ReServiceID == reService.ID {
				delete(defaultServiceItems, reService.ID)
			}
		}
	}

	for _, serviceItem := range defaultServiceItems {
		verrs, err := o.db.ValidateAndCreate(&serviceItem)

		if err != nil || (verrs != nil && verrs.HasAny()) {
			return fmt.Errorf("%v %#v", err, verrs)
		}

		moveTaskOrder.MTOServiceItems = append(moveTaskOrder.MTOServiceItems, serviceItem)
	}

	return nil
}
