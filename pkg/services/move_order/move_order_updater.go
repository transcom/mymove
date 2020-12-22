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

		if moveOrder.ServiceMember.Affiliation != nil {
			existingOrder.ServiceMember.Affiliation = moveOrder.ServiceMember.Affiliation
			err = tx.Save(&existingOrder.ServiceMember)
			if err != nil {
				return err
			}
		}

		if moveOrder.Entitlement.DBAuthorizedWeight != nil || moveOrder.Entitlement.DependentsAuthorized != nil {

			if moveOrder.Entitlement.DBAuthorizedWeight != nil {
				existingOrder.Entitlement.DBAuthorizedWeight = moveOrder.Entitlement.DBAuthorizedWeight
			}

			if moveOrder.Entitlement.DependentsAuthorized != nil {
				existingOrder.Entitlement.DependentsAuthorized = moveOrder.Entitlement.DependentsAuthorized
			}

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

		if moveOrder.Grade != nil {
			existingOrder.Grade = moveOrder.Grade
		}

		if moveOrder.OrdersTypeDetail != nil {
			existingOrder.OrdersTypeDetail = moveOrder.OrdersTypeDetail
		}

		if moveOrder.TAC != nil {
			existingOrder.TAC = moveOrder.TAC
		}

		if moveOrder.SAC != nil {
			existingOrder.SAC = moveOrder.SAC
		}

		if moveOrder.OrdersNumber != nil {
			existingOrder.OrdersNumber = moveOrder.OrdersNumber
		}

		if moveOrder.DepartmentIndicator != nil {
			existingOrder.DepartmentIndicator = moveOrder.DepartmentIndicator
		}

		existingOrder.IssueDate = moveOrder.IssueDate
		existingOrder.ReportByDate = moveOrder.ReportByDate
		existingOrder.OrdersType = moveOrder.OrdersType

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
