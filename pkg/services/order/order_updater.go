package order

import (
	"io"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

type orderUpdater struct {
	db *pop.Connection
}

// NewOrderUpdater creates a new struct with the service dependencies
func NewOrderUpdater(db *pop.Connection) services.OrderUpdater {
	return &orderUpdater{db}
}

// UpdateOrderAsTOO updates an order as permitted by a TOO
func (f *orderUpdater) UpdateOrderAsTOO(orderID uuid.UUID, payload ghcmessages.UpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, services.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	orderToUpdate := orderFromTOOPayload(*order, payload)

	return f.updateOrder(orderToUpdate, CheckRequiredFields())
}

// UpdateOrderAsCounselor updates an order as permitted by a service counselor
func (f *orderUpdater) UpdateOrderAsCounselor(orderID uuid.UUID, payload ghcmessages.CounselingUpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, services.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	orderToUpdate := orderFromCounselingPayload(*order, payload)

	return f.updateOrder(orderToUpdate)
}

// UpdateAllowanceAsTOO updates an allowance as permitted by a service counselor
func (f *orderUpdater) UpdateAllowanceAsTOO(orderID uuid.UUID, payload ghcmessages.UpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, services.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	orderToUpdate := allowanceFromTOOPayload(*order, payload)

	return f.updateOrder(orderToUpdate)
}

// UpdateAllowanceAsCounselor updates an allowance as permitted by a service counselor
func (f *orderUpdater) UpdateAllowanceAsCounselor(orderID uuid.UUID, payload ghcmessages.CounselingUpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, services.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	orderToUpdate := allowanceFromCounselingPayload(*order, payload)

	return f.updateOrder(orderToUpdate)
}

// UploadAmendedOrdersAsCustomer add amended order documents to an existing order
func (f *orderUpdater) UploadAmendedOrdersAsCustomer(logger *zap.Logger, userID uuid.UUID, orderID uuid.UUID, file io.ReadCloser, filename string, storer storage.FileStorer) (models.Upload, string, *validate.Errors, error) {
	orderToUpdate, findErr := f.findOrderWithAmendedOrders(orderID)
	if findErr != nil {
		return models.Upload{}, "", nil, findErr
	}

	userUpload, url, verrs, err := f.amendedOrder(f.db, logger, userID, *orderToUpdate, file, filename, storer)
	if verrs.HasAny() || err != nil {
		return models.Upload{}, "", verrs, err
	}

	return userUpload.Upload, url, nil, nil
}

func (f *orderUpdater) findOrder(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := f.db.Q().EagerPreload("Moves", "ServiceMember", "Entitlement").Find(&order, orderID)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(orderID, "while looking for order")
		}
	}

	return &order, nil
}

func (f *orderUpdater) findOrderWithAmendedOrders(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := f.db.Q().EagerPreload("ServiceMember", "UploadedAmendedOrders").Find(&order, orderID)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(orderID, "while looking for order")
		}
	}

	return &order, nil
}

func orderFromTOOPayload(existingOrder models.Order, payload ghcmessages.UpdateOrderPayload) models.Order {
	order := existingOrder

	// update both order origin duty station and service member duty station
	if payload.OriginDutyStationID != nil {
		originDutyStationID := uuid.FromStringOrNil(payload.OriginDutyStationID.String())
		order.OriginDutyStationID = &originDutyStationID
		order.ServiceMember.DutyStationID = &originDutyStationID
	}

	if payload.NewDutyStationID != nil {
		newDutyStationID := uuid.FromStringOrNil(payload.NewDutyStationID.String())
		order.NewDutyStationID = newDutyStationID
	}

	if payload.DepartmentIndicator != nil {
		departmentIndicator := (*string)(payload.DepartmentIndicator)
		order.DepartmentIndicator = departmentIndicator
	}

	if payload.IssueDate != nil {
		order.IssueDate = time.Time(*payload.IssueDate)
	}

	if payload.OrdersNumber != nil {
		order.OrdersNumber = payload.OrdersNumber
	}

	if payload.OrdersTypeDetail != nil {
		orderTypeDetail := internalmessages.OrdersTypeDetail(*payload.OrdersTypeDetail)
		order.OrdersTypeDetail = &orderTypeDetail
	}

	if payload.ReportByDate != nil {
		order.ReportByDate = time.Time(*payload.ReportByDate)
	}

	if payload.Sac != nil {
		order.SAC = payload.Sac
	}

	if payload.Tac != nil {
		order.TAC = payload.Tac
	}

	if payload.OrdersType != nil {
		order.OrdersType = internalmessages.OrdersType(*payload.OrdersType)
	}

	// if the order has amended order documents and it has not been previously acknowledged record the current timestamp
	if payload.OrdersAcknowledgement != nil && *payload.OrdersAcknowledgement && existingOrder.UploadedAmendedOrdersID != nil && existingOrder.AmendedOrdersAcknowledgedAt == nil {
		acknowledgedAt := time.Now()
		order.AmendedOrdersAcknowledgedAt = &acknowledgedAt
	}

	return order
}

func (f *orderUpdater) amendedOrder(db *pop.Connection, logger *zap.Logger, userID uuid.UUID, order models.Order, file io.ReadCloser, filename string, storer storage.FileStorer) (models.UserUpload, string, *validate.Errors, error) {

	// If Order does not have a Document for amended orders uploads, then create a new one
	var err error
	savedAmendedOrdersDoc := order.UploadedAmendedOrders
	if order.UploadedAmendedOrders == nil {
		amendedOrdersDoc := &models.Document{
			ServiceMemberID: order.ServiceMemberID,
		}
		savedAmendedOrdersDoc, err = f.saveDocumentForAmendedOrder(amendedOrdersDoc)
		if err != nil {
			return models.UserUpload{}, "", nil, err
		}

		// save new UploadedAmendedOrdersID (document ID) to orders
		order.UploadedAmendedOrders = savedAmendedOrdersDoc
		order.UploadedAmendedOrdersID = &savedAmendedOrdersDoc.ID
		_, _, err = f.updateOrder(order)
		if err != nil {
			return models.UserUpload{}, "", nil, err
		}
	}

	// Create new user upload for amended order
	var userUpload *models.UserUpload
	var verrs *validate.Errors
	var url string
	userUpload, url, verrs, err = uploader.CreateUserUploadForDocumentWrapper(
		db,
		logger,
		userID,
		storer,
		file,
		filename,
		uploader.MaxCustomerUserUploadFileSizeLimit,
		&savedAmendedOrdersDoc.ID,
	)

	if verrs.HasAny() || err != nil {
		return models.UserUpload{}, "", verrs, err
	}

	order.UploadedAmendedOrders.UserUploads = append(order.UploadedAmendedOrders.UserUploads, *userUpload)

	return *userUpload, url, nil, nil
}

func orderFromCounselingPayload(existingOrder models.Order, payload ghcmessages.CounselingUpdateOrderPayload) models.Order {
	order := existingOrder

	// update both order origin duty station and service member duty station
	if payload.OriginDutyStationID != nil {
		originDutyStationID := uuid.FromStringOrNil(payload.OriginDutyStationID.String())
		order.OriginDutyStationID = &originDutyStationID
		order.ServiceMember.DutyStationID = &originDutyStationID
	}

	if payload.NewDutyStationID != nil {
		newDutyStationID := uuid.FromStringOrNil(payload.NewDutyStationID.String())
		order.NewDutyStationID = newDutyStationID
	}

	if payload.IssueDate != nil {
		order.IssueDate = time.Time(*payload.IssueDate)
	}

	if payload.ReportByDate != nil {
		order.ReportByDate = time.Time(*payload.ReportByDate)
	}

	if payload.OrdersType != nil {
		order.OrdersType = internalmessages.OrdersType(*payload.OrdersType)
	}

	return order
}

func allowanceFromTOOPayload(existingOrder models.Order, payload ghcmessages.UpdateAllowancePayload) models.Order {
	order := existingOrder

	if payload.ProGearWeight != nil {
		order.Entitlement.ProGearWeight = int(*payload.ProGearWeight)
	}

	if payload.ProGearWeightSpouse != nil {
		order.Entitlement.ProGearWeightSpouse = int(*payload.ProGearWeightSpouse)
	}

	if payload.RequiredMedicalEquipmentWeight != nil {
		order.Entitlement.RequiredMedicalEquipmentWeight = int(*payload.RequiredMedicalEquipmentWeight)
	}

	// branch for service member
	if payload.Agency != "" {
		serviceMemberAffiliation := models.ServiceMemberAffiliation(payload.Agency)
		order.ServiceMember.Affiliation = &serviceMemberAffiliation
	}

	// rank
	if payload.Grade != nil {
		grade := (*string)(payload.Grade)
		order.Grade = grade
	}

	if payload.OrganizationalClothingAndIndividualEquipment != nil {
		order.Entitlement.OrganizationalClothingAndIndividualEquipment = *payload.OrganizationalClothingAndIndividualEquipment
	}

	if payload.AuthorizedWeight != nil {
		dbAuthorizedWeight := swag.Int(int(*payload.AuthorizedWeight))
		order.Entitlement.DBAuthorizedWeight = dbAuthorizedWeight
	}

	if payload.DependentsAuthorized != nil {
		order.Entitlement.DependentsAuthorized = payload.DependentsAuthorized
	}

	return order
}

func allowanceFromCounselingPayload(existingOrder models.Order, payload ghcmessages.CounselingUpdateAllowancePayload) models.Order {
	order := existingOrder

	if payload.ProGearWeight != nil {
		order.Entitlement.ProGearWeight = int(*payload.ProGearWeight)
	}

	if payload.ProGearWeightSpouse != nil {
		order.Entitlement.ProGearWeightSpouse = int(*payload.ProGearWeightSpouse)
	}

	if payload.RequiredMedicalEquipmentWeight != nil {
		order.Entitlement.RequiredMedicalEquipmentWeight = int(*payload.RequiredMedicalEquipmentWeight)
	}

	// branch for service member
	if payload.Agency != "" {
		serviceMemberAffiliation := models.ServiceMemberAffiliation(payload.Agency)
		order.ServiceMember.Affiliation = &serviceMemberAffiliation
	}

	// rank
	if payload.Grade != nil {
		grade := (*string)(payload.Grade)
		order.Grade = grade
	}

	if payload.OrganizationalClothingAndIndividualEquipment != nil {
		order.Entitlement.OrganizationalClothingAndIndividualEquipment = *payload.OrganizationalClothingAndIndividualEquipment
	}

	if payload.DependentsAuthorized != nil {
		order.Entitlement.DependentsAuthorized = payload.DependentsAuthorized
	}

	return order
}

func (f *orderUpdater) saveDocumentForAmendedOrder(doc *models.Document) (*models.Document, error) {

	var docID uuid.UUID
	if doc != nil {
		docID = doc.ID
	}

	handleError := func(verrs *validate.Errors, err error) error {
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(docID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		return nil
	}

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {
		var verrs *validate.Errors
		var err error
		verrs, err = tx.ValidateAndSave(doc)
		if e := handleError(verrs, err); e != nil {
			return e
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return doc, nil
}

func (f *orderUpdater) updateOrder(order models.Order, checks ...Validator) (*models.Order, uuid.UUID, error) {
	handleError := func(verrs *validate.Errors, err error) error {
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(order.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		return nil
	}

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {
		var verrs *validate.Errors
		var err error

		if verr := ValidateOrder(&order, checks...); verr != nil {
			return verr
		}

		// update service member
		if order.Grade != nil {
			// keep grade and rank in sync
			order.ServiceMember.Rank = (*models.ServiceMemberRank)(order.Grade)
		}

		if order.OriginDutyStationID != nil {
			// TODO refactor to use service objects to fetch duty station
			var originDutyStation models.DutyStation
			originDutyStation, err = models.FetchDutyStation(f.db, *order.OriginDutyStationID)
			if e := handleError(verrs, err); e != nil {
				if errors.Cause(e).Error() == models.RecordNotFoundErrorString {
					return services.NewNotFoundError(*order.OriginDutyStationID, "while looking for OriginDutyStation")
				}
			}
			order.OriginDutyStationID = &originDutyStation.ID
			order.OriginDutyStation = &originDutyStation

			order.ServiceMember.DutyStationID = &originDutyStation.ID
			order.ServiceMember.DutyStation = originDutyStation
		}

		if order.Grade != nil || order.OriginDutyStationID != nil {
			verrs, err = tx.ValidateAndUpdate(&order.ServiceMember)
			if e := handleError(verrs, err); e != nil {
				return e
			}
		}

		// update entitlement
		if order.Entitlement != nil {
			verrs, err = tx.ValidateAndUpdate(order.Entitlement)
			if e := handleError(verrs, err); e != nil {
				return e
			}
		}

		if order.NewDutyStationID != uuid.Nil {
			// TODO refactor to use service objects to fetch duty station
			var newDutyStation models.DutyStation
			newDutyStation, err = models.FetchDutyStation(f.db, order.NewDutyStationID)
			if e := handleError(verrs, err); e != nil {
				if errors.Cause(e).Error() == models.RecordNotFoundErrorString {
					return services.NewNotFoundError(order.NewDutyStationID, "while looking for NewDutyStation")
				}
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
		return nil, uuid.Nil, transactionError
	}

	var moveID uuid.UUID
	if len(order.Moves) > 0 {
		moveID = order.Moves[0].ID
	}
	return &order, moveID, nil
}
