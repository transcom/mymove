package order

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type orderUpdater struct {
	db *pop.Connection
}

// NewOrderUpdater creates a new struct with the service dependencies
func NewOrderUpdater(db *pop.Connection) services.OrderUpdater {
	return &orderUpdater{db}
}

// UpdateOrder updates the Order model
func (s *orderUpdater) UpdateOrder(order models.Order) (*models.Order, error) {
	handleError := func(verrs *validate.Errors, err error) error {
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(order.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		return nil
	}

	transactionError := s.db.Transaction(func(tx *pop.Connection) error {
		var verrs *validate.Errors
		var err error

		// update service member
		if order.Grade != nil {
			// keep grade and rank in sync
			order.ServiceMember.Rank = (*models.ServiceMemberRank)(order.Grade)
		}

		verrs, err = tx.ValidateAndUpdate(&order.ServiceMember)
		if e := handleError(verrs, err); e != nil {
			return e
		}

		// update entitlement
		if order.Entitlement != nil {
			verrs, err = tx.ValidateAndUpdate(order.Entitlement)
			if e := handleError(verrs, err); e != nil {
				return e
			}
		}

		if order.OriginDutyStationID != nil {
			// TODO refactor to use service objects to fetch duty station
			var originDutyStation models.DutyStation
			originDutyStation, err = models.FetchDutyStation(s.db, *order.OriginDutyStationID)
			if e := handleError(verrs, err); e != nil {
				return e
			}
			order.OriginDutyStationID = &originDutyStation.ID
			order.OriginDutyStation = &originDutyStation

			order.ServiceMember.DutyStationID = &originDutyStation.ID
			order.ServiceMember.DutyStation = originDutyStation

			verrs, err = tx.ValidateAndUpdate(&order.ServiceMember)
			if e := handleError(verrs, err); e != nil {
				return e
			}
		}

		if order.NewDutyStationID != uuid.Nil {
			// TODO refactor to use service objects to fetch duty station
			var newDutyStation models.DutyStation
			newDutyStation, err = models.FetchDutyStation(s.db, order.NewDutyStationID)
			if e := handleError(verrs, err); e != nil {
				return e
			}
			order.NewDutyStationID = newDutyStation.ID
			order.NewDutyStation = newDutyStation
		}

		verrs, err = tx.ValidateAndUpdate(&order)
		if e := handleError(verrs, err); e != nil {
			return e
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &order, nil
}
