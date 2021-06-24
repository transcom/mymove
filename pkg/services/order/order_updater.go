package order

import (
	"fmt"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
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

	return f.updateOrder(orderToUpdate)
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

// UploadAmendedOrders add amended order documents to an existing order
func (f *orderUpdater) UploadAmendedOrders(orderToUpdate models.Order, payload *internalmessages.UserUploadPayload, eTag string) (*models.Order, uuid.UUID, error) {
	existingETag := etag.GenerateEtag(orderToUpdate.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, services.NewPreconditionFailedError(orderToUpdate.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	amendedOrder := f.amendedOrderFromUserUploadPayload(orderToUpdate, payload)

	return f.updateOrder(amendedOrder)
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

	order.OrdersType = internalmessages.OrdersType(payload.OrdersType)

	return order
}

func (f *orderUpdater) amendedOrderFromUserUploadPayload(existingOrder models.Order, payload *internalmessages.UserUploadPayload) models.Order {
	order := existingOrder
	// ultimately want to take UserUploadPayload and attach it to order.UploadedAmendedOrders
	// UserUploadPayload doesn't have a document attached
	// this service will check if orders already has an amendedOrders doc and
	// add these uploads to them if that doc already exists
	// however, if it doesn't exist, this service will create a new doc to attach the upload to
	// there's a foreign key relationship between these things
	//	- orders.uploaded_amended_orders_id must exist in the document table
	//	- userUpload.document_id must exist in the document table
	//	- userUpload.upload_id must exist in the uploads table

	if order.UploadedAmendedOrdersID == nil {
		// create an amended orders document ID
		// this will be used for
		// - orders.uploaded_amended_orders_id
		// - userUpload.document_id
		// - and to create and save a document
		// 1) create and save the base document
		// 	order.UploadedAmendedOrdersID = &amendedOrdersID
		amendedOrdersDocumentID := uuid.Must(uuid.NewV4())
		amendedOrdersDoc := models.Document{
			ID:              amendedOrdersDocumentID,
			ServiceMemberID: order.UploadedOrders.ServiceMemberID,
			CreatedAt:       time.Now(),
		}
		// amendedOrdersDoc.UserUploads = append(amendedOrdersDoc.UserUploads, *savedUserUpload)

		savedAmendedDoc, err := f.saveDocumentForAmendedOrder(amendedOrdersDoc)
		if err != nil {
			fmt.Printf("========================== ERRORR SAVING DOCUMENT ===========================")
		}
		// 2) massage userupload payload into a userUpload model linked to the newly created document

	}

	// var userUpload models.UserUpload
	// userUpload.ID = uuid.FromStringOrNil(payload.ID.String())
	// if order.UploadedAmendedOrdersID != nil {
	// 	userUpload.Document = *order.UploadedAmendedOrders
	// 	userUpload.DocumentID = order.UploadedAmendedOrdersID
	// }

	// upload := models.Upload{
	// 	ID:          uuid.FromStringOrNil(payload.Upload.ID.String()),
	// 	UploadType:  models.UploadTypeUSER,
	// 	ContentType: *payload.Upload.ContentType,
	// 	Checksum:    payload.Upload.Checksum,
	// }
	// if payload.Upload.Filename != nil {
	// 	userUpload.Upload.Filename = *payload.Upload.Filename
	// 	upload.Filename = *payload.Upload.Filename
	// }
	// if payload.Upload.Bytes != nil {
	// 	userUpload.Upload.Bytes = int64(*payload.Upload.Bytes)
	// 	upload.Bytes = int64(*payload.Upload.Bytes)
	// }
	// if payload.Upload.ContentType != nil {
	// 	userUpload.Upload.ContentType = *payload.Upload.ContentType
	// 	upload.ContentType = *payload.Upload.ContentType
	// }
	// savedUpload, err := f.saveUploadForAmendedOrder(upload)
	// if err != nil {
	// 	fmt.Printf("========================== ERRORR SAVING UPLOAD: %v ===========================", err)
	// }
	// if savedUpload != nil {
	// 	userUpload.Upload = models.Upload{
	// 		ID: savedUpload.ID,
	// 	}
	// }

	// if payload.Upload.CreatedAt != nil {
	// 	userUpload.Upload.CreatedAt = time.Time(*payload.Upload.CreatedAt)
	// }
	// if payload.Upload.UpdatedAt != nil {
	// 	userUpload.Upload.UpdatedAt = time.Time(*payload.Upload.UpdatedAt)
	// }

	// if payload.CreatedAt != nil {
	// 	userUpload.CreatedAt = time.Time(*payload.CreatedAt)
	// }
	// if payload.UpdatedAt != nil {
	// 	userUpload.UpdatedAt = time.Time(*payload.UpdatedAt)
	// }
	// amendedOrdersID := uuid.Must(uuid.NewV4())

	// userUpload.UploadID = uuid.FromStringOrNil(payload.UploadID.String())
	// userUpload.UploaderID = uuid.FromStringOrNil(payload.UploaderID.String())
	// userUpload.DocumentID = &amendedOrdersID
	// savedUserUpload, err := f.saveUserUploadForAmendedOrder(userUpload)

	// if order.UploadedAmendedOrdersID == nil {
	// 	order.UploadedAmendedOrdersID = &amendedOrdersID

	// 	amendedOrdersDoc := models.Document{
	// 		ID:              amendedOrdersID,
	// 		ServiceMemberID: order.UploadedOrders.ServiceMemberID,
	// 		CreatedAt:       time.Now(),
	// 	}
	// 	amendedOrdersDoc.UserUploads = append(amendedOrdersDoc.UserUploads, *savedUserUpload)

	// 	savedAmendedDoc, err := f.saveDocumentForAmendedOrder(amendedOrdersDoc)
	// 	if err != nil {
	// 		fmt.Printf("========================== ERRORR SAVING DOCUMENT ===========================")
	// 	}
	// 	savedUserUpload.DocumentID = &savedAmendedDoc.ID
	// 	savedUserUpload.Document = *savedAmendedDoc
	// 	order.UploadedAmendedOrders = savedAmendedDoc
	// 	fmt.Printf("========================== amendedOrdersDoc.ID: %v ===========================", savedAmendedDoc.ID)
	// } else {
	// 	if savedUserUpload.DocumentID == nil {
	// 		savedUserUpload.DocumentID = order.UploadedAmendedOrdersID
	// 		savedUserUpload.Document = *order.UploadedAmendedOrders
	// 	}
	// }

	// if err != nil {
	// 	fmt.Printf("========================== ERRORR SAVING USER_ UPLOAD ===========================")
	// }
	// order.UploadedAmendedOrders.UserUploads = append(order.UploadedAmendedOrders.UserUploads, *savedUserUpload)
	// fmt.Printf("========================== orders.UploadedAmendedOrdersID: %v ===========================", order.UploadedAmendedOrdersID)

	// return order
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

	order.OrdersType = internalmessages.OrdersType(payload.OrdersType)

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

func (f *orderUpdater) saveUserUploadForAmendedOrder(userUpload models.UserUpload) (*models.UserUpload, error) {
	handleError := func(verrs *validate.Errors, err error) error {
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(userUpload.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		return nil
	}

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {
		var verrs *validate.Errors
		var err error

		verrs, err = tx.ValidateAndSave(&userUpload)
		if e := handleError(verrs, err); e != nil {
			return e
		}

		return nil
	})

	if transactionError != nil {

		return nil, transactionError
	}
	return &userUpload, nil
}

func (f *orderUpdater) saveUploadForAmendedOrder(upload models.Upload) (*models.Upload, error) {
	handleError := func(verrs *validate.Errors, err error) error {
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(upload.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		return nil
	}

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {
		var verrs *validate.Errors
		var err error

		verrs, err = tx.ValidateAndSave(&upload)
		if e := handleError(verrs, err); e != nil {
			return e
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}
	return &upload, nil
}

func (f *orderUpdater) saveDocumentForAmendedOrder(doc models.Document) (*models.Document, error) {
	handleError := func(verrs *validate.Errors, err error) error {
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(doc.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		return nil
	}

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {
		var verrs *validate.Errors
		var err error
		fmt.Printf("==================================================")
		fmt.Printf("=====================IN SAVE DOCUMENT TRANSACTION =============================")
		fmt.Printf("==================================================")
		verrs, err = tx.ValidateAndSave(&doc)
		if e := handleError(verrs, err); e != nil {
			return e
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}
	return &doc, nil
}

func (f *orderUpdater) updateOrder(order models.Order) (*models.Order, uuid.UUID, error) {
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

		// update uploaded amended orders
		if order.UploadedAmendedOrdersID != nil {
			verrs, err = tx.ValidateAndUpdate(order.UploadedAmendedOrders)
			if e := handleError(verrs, err); e != nil {
				fmt.Printf("=============================================")
				fmt.Printf("=====================ERROR: %v ========================", e)
				fmt.Printf("=============================================")
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

	return &order, order.Moves[0].ID, nil
}
