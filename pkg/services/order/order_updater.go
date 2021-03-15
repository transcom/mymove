package order

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type orderUpdater struct {
	db *pop.Connection
	orderFetcher
}

// NewOrderUpdater creates a new struct with the service dependencies
func NewOrderUpdater(db *pop.Connection) services.OrderUpdater {
	return &orderUpdater{db, orderFetcher{db}}
}

// UpdateOrder updates the Order model
func (s *orderUpdater) UpdateOrder(eTag string, order models.Order) (*models.Order, error) {
	existingOrder, err := s.orderFetcher.FetchOrder(order.ID)
	if err != nil {
		return nil, services.NewNotFoundError(order.ID, "while looking for order")
	}

	existingETag := etag.GenerateEtag(existingOrder.UpdatedAt)
	if existingETag != eTag {
		return nil, services.NewPreconditionFailedError(order.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	transactionError := s.db.Transaction(func(tx *pop.Connection) error {

		if order.ServiceMember.Affiliation != nil {
			existingOrder.ServiceMember.Affiliation = order.ServiceMember.Affiliation
			err = tx.Save(&existingOrder.ServiceMember)
			if err != nil {
				return err
			}
		}

		if entitlement := order.Entitlement; entitlement != nil && (entitlement.DBAuthorizedWeight != nil || entitlement.DependentsAuthorized != nil) {

			if entitlement.DBAuthorizedWeight != nil {
				existingOrder.Entitlement.DBAuthorizedWeight = entitlement.DBAuthorizedWeight
			}

			if entitlement.DependentsAuthorized != nil {
				existingOrder.Entitlement.DependentsAuthorized = entitlement.DependentsAuthorized
			}

			err = tx.Save(existingOrder.Entitlement)
			if err != nil {
				return err
			}
		}

		if order.OriginDutyStationID != existingOrder.OriginDutyStationID {
			originDutyStation, fetchErr := models.FetchDutyStation(s.db, *order.OriginDutyStationID)
			if fetchErr != nil {
				return services.NewInvalidInputError(order.ID, fetchErr, nil, "unable to find origin duty station")
			}
			existingOrder.OriginDutyStationID = order.OriginDutyStationID
			existingOrder.OriginDutyStation = &originDutyStation
		}

		if order.NewDutyStationID != existingOrder.NewDutyStationID {
			newDutyStation, fetchErr := models.FetchDutyStation(s.db, order.NewDutyStationID)
			if fetchErr != nil {
				return services.NewInvalidInputError(order.ID, fetchErr, nil, "unable to find destination duty station")
			}
			existingOrder.NewDutyStationID = order.NewDutyStationID
			existingOrder.NewDutyStation = newDutyStation
		}

		if order.Grade != nil {
			existingOrder.Grade = order.Grade
		}

		if order.OrdersTypeDetail != nil {
			existingOrder.OrdersTypeDetail = order.OrdersTypeDetail
		}

		if order.TAC != nil {
			existingOrder.TAC = order.TAC
		}

		if order.SAC != nil {
			existingOrder.SAC = order.SAC
		}

		if order.OrdersNumber != nil {
			existingOrder.OrdersNumber = order.OrdersNumber
		}

		if order.DepartmentIndicator != nil {
			existingOrder.DepartmentIndicator = order.DepartmentIndicator
		}

		existingOrder.IssueDate = order.IssueDate
		existingOrder.ReportByDate = order.ReportByDate
		existingOrder.OrdersType = order.OrdersType

		// optimistic locking handled before transaction block
		verrs, updateErr := tx.ValidateAndUpdate(existingOrder)

		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(order.ID, err, verrs, "")
		}

		if updateErr != nil {
			return updateErr
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return existingOrder, err
}
