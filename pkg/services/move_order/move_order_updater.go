package moveorder

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type moveOrderUpdater struct {
	db *pop.Connection
	moveOrderFetcher
	builder UpdateMoveOrderQueryBuilder
}

// NewMoveOrderUpdater creates a new struct with the service dependencies
func NewMoveOrderUpdater(db *pop.Connection, builder UpdateMoveOrderQueryBuilder) services.MoveOrderUpdater {
	return &moveOrderUpdater{db, moveOrderFetcher{db}, builder}
}

// UpdateMoveOrderQueryBuilder interface performs fetch and updates during move order update
type UpdateMoveOrderQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

func (s *moveOrderUpdater) UpdateMoveOrder(moveOrderID uuid.UUID, eTag string, moveOrder models.Order) (*models.Order, error) {

	existingOrder, err := s.moveOrderFetcher.FetchMoveOrder(moveOrder.ID)
	if err != nil {
		return nil, services.NewNotFoundError(moveOrder.ID, "while looking for moveOrder")
	}

	existingETag := etag.GenerateEtag(existingOrder.UpdatedAt)
	if existingETag != eTag {
		return nil, services.NewPreconditionFailedError(moveOrder.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	transactionError := s.db.Transaction(func(tx *pop.Connection) error {

		if moveOrder.Entitlement.DBAuthorizedWeight != nil {
			existingOrder.Entitlement.DBAuthorizedWeight = moveOrder.Entitlement.DBAuthorizedWeight
			err = tx.Save(existingOrder.Entitlement)
			if err != nil {
				return err
			}
		}

		if moveOrder.OriginDutyStationID != existingOrder.OriginDutyStationID {
			originDutyStation, fetchErr := models.FetchDutyStation(s.db, *moveOrder.OriginDutyStationID)
			if fetchErr != nil {
				return services.NewInvalidInputError(moveOrder.ID, fetchErr, nil, "unable to find origin duty station")
			}
			existingOrder.OriginDutyStationID = moveOrder.OriginDutyStationID
			existingOrder.OriginDutyStation = &originDutyStation
		}

		if moveOrder.NewDutyStationID != existingOrder.NewDutyStationID {
			newDutyStation, fetchErr := models.FetchDutyStation(s.db, moveOrder.NewDutyStationID)
			if fetchErr != nil {
				return services.NewInvalidInputError(moveOrder.ID, fetchErr, nil, "unable to find destination duty station")
			}
			existingOrder.NewDutyStationID = moveOrder.NewDutyStationID
			existingOrder.NewDutyStation = newDutyStation
		}

		existingOrder.IssueDate = moveOrder.IssueDate
		existingOrder.ReportByDate = moveOrder.ReportByDate
		existingOrder.OrdersType = moveOrder.OrdersType
		existingOrder.OrdersTypeDetail = moveOrder.OrdersTypeDetail
		existingOrder.OrdersNumber = moveOrder.OrdersNumber
		existingOrder.TAC = moveOrder.TAC
		existingOrder.SAC = moveOrder.SAC
		existingOrder.DepartmentIndicator = moveOrder.DepartmentIndicator

		verrs, updateErr := s.builder.UpdateOne(existingOrder, &eTag)

		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(moveOrder.ID, err, verrs, "")
		}

		if updateErr != nil {
			switch updateErr.(type) {
			case query.StaleIdentifierError:
				return services.NewPreconditionFailedError(moveOrder.ID, err)
			default:
				return updateErr
			}
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return existingOrder, err
}
