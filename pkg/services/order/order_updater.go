package order

import (
	"database/sql"
	"io"
	"strings"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

type orderUpdater struct {
	moveRouter services.MoveRouter
}

// NewOrderUpdater creates a new struct with the service dependencies
func NewOrderUpdater(moveRouter services.MoveRouter) services.OrderUpdater {
	return &orderUpdater{moveRouter}
}

// UpdateOrderAsTOO updates an order as permitted by a TOO
func (f *orderUpdater) UpdateOrderAsTOO(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.UpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error) {
	const SAC_LIMIT = 80
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, apperror.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	if payload.Sac.Present && payload.Sac.Value != nil && len(*payload.Sac.Value) > SAC_LIMIT {
		return &models.Order{}, uuid.Nil, apperror.NewInvalidInputError(orderID, nil, nil, "SAC cannot be more than 80 characters")
	}

	orderToUpdate, err := orderFromTOOPayload(appCtx, *order, payload)
	if err != nil {
		return nil, uuid.Nil, err
	}

	return f.updateOrderAsTOO(appCtx, orderToUpdate, CheckRequiredFields())
}

// UpdateOrderAsCounselor updates an order as permitted by a service counselor
func (f *orderUpdater) UpdateOrderAsCounselor(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.CounselingUpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error) {
	const SAC_LIMIT = 80
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, apperror.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	if payload.Sac.Present && payload.Sac.Value != nil && len(*payload.Sac.Value) > SAC_LIMIT {
		return &models.Order{}, uuid.Nil, apperror.NewInvalidInputError(orderID, nil, nil, "SAC cannot be more than 80 characters")
	}

	orderToUpdate, err := orderFromCounselingPayload(appCtx, *order, payload)
	if err != nil {
		return &models.Order{}, uuid.Nil, nil
	}

	return f.updateOrder(appCtx, orderToUpdate, CheckRequiredFields())
}

// UpdateAllowanceAsTOO updates an allowance as permitted by a service counselor
func (f *orderUpdater) UpdateAllowanceAsTOO(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.UpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, apperror.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	orderToUpdate, err := allowanceFromTOOPayload(appCtx, *order, payload)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	return f.updateOrder(appCtx, orderToUpdate)
}

// UpdateAllowanceAsCounselor updates an allowance as permitted by a service counselor
func (f *orderUpdater) UpdateAllowanceAsCounselor(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.CounselingUpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, apperror.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	orderToUpdate, err := allowanceFromCounselingPayload(appCtx, *order, payload)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	return f.updateOrder(appCtx, orderToUpdate)
}

// UploadAmendedOrdersAsCustomer add amended order documents to an existing order
func (f *orderUpdater) UploadAmendedOrdersAsCustomer(appCtx appcontext.AppContext, userID uuid.UUID, orderID uuid.UUID, file io.ReadCloser, filename string, storer storage.FileStorer) (models.Upload, string, *validate.Errors, error) {
	orderToUpdate, findErr := f.findOrderWithAmendedOrders(appCtx, orderID)
	if findErr != nil {
		return models.Upload{}, "", nil, findErr
	}

	userUpload, url, verrs, err := f.amendedOrder(appCtx, userID, *orderToUpdate, file, filename, storer, models.UploadTypeUSER)
	if verrs.HasAny() || err != nil {
		return models.Upload{}, "", verrs, err
	}

	return userUpload.Upload, url, nil, nil
}

// UploadAmendedOrdersAsOffice add amended order documents to an existing order
func (f *orderUpdater) UploadAmendedOrdersAsOffice(appCtx appcontext.AppContext, userID uuid.UUID, orderID uuid.UUID, file io.ReadCloser, filename string, storer storage.FileStorer) (models.Upload, string, *validate.Errors, error) {
	orderToUpdate, findErr := f.findOrderWithAmendedOrders(appCtx, orderID)
	if findErr != nil {
		return models.Upload{}, "", nil, findErr
	}

	userUpload, url, verrs, err := f.amendedOrder(appCtx, userID, *orderToUpdate, file, filename, storer, models.UploadTypeOFFICE)
	if verrs.HasAny() || err != nil {
		return models.Upload{}, "", verrs, err
	}

	return userUpload.Upload, url, nil, nil
}

func (f *orderUpdater) findOrder(appCtx appcontext.AppContext, orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := appCtx.DB().Q().EagerPreload("Moves", "ServiceMember", "Entitlement").Find(&order, orderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(orderID, "while looking for order")
		default:
			return nil, apperror.NewQueryError("Order", err, "")
		}
	}

	return &order, nil
}

func (f *orderUpdater) findOrderWithAmendedOrders(appCtx appcontext.AppContext, orderID uuid.UUID) (*models.Order, error) {
	var order models.Order

	query := appCtx.DB().Q().EagerPreload("ServiceMember", "UploadedAmendedOrders")

	if appCtx.Session().IsMilApp() {
		query = query.Where("orders.service_member_id = ?", appCtx.Session().ServiceMemberID)
	}

	err := query.Find(&order, orderID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(orderID, "while looking for order")
		default:
			return nil, apperror.NewQueryError("Order", err, "")
		}
	}

	return &order, nil
}

func orderFromTOOPayload(appCtx appcontext.AppContext, existingOrder models.Order, payload ghcmessages.UpdateOrderPayload) (models.Order, error) {
	order := existingOrder
	waf := entitlements.NewWeightAllotmentFetcher()

	// update order origin duty location
	if payload.OriginDutyLocationID != nil {
		originDutyLocationID := uuid.FromStringOrNil(payload.OriginDutyLocationID.String())
		order.OriginDutyLocationID = &originDutyLocationID
	}

	if payload.NewDutyLocationID != nil {
		newDutyLocationID := uuid.FromStringOrNil(payload.NewDutyLocationID.String())
		order.NewDutyLocationID = newDutyLocationID
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

	if payload.Sac.Present {
		if payload.Sac.Value != nil && *payload.Sac.Value == "" {
			order.SAC = nil
		} else {
			order.SAC = payload.Sac.Value
		}

	}

	if payload.Tac != nil {
		normalizedTac := strings.ToUpper(*payload.Tac)
		order.TAC = &normalizedTac
	}

	if payload.NtsSac.Present {
		if payload.NtsSac.Value != nil && *payload.NtsSac.Value == "" {
			order.NtsSAC = nil
		} else {
			order.NtsSAC = payload.NtsSac.Value
		}
	}

	if payload.NtsTac.Present {
		if payload.NtsTac.Value != nil && *payload.NtsTac.Value != "" {
			normalizedNtsTac := strings.ToUpper(*payload.NtsTac.Value)
			order.NtsTAC = &normalizedNtsTac
		} else {
			order.NtsTAC = nil
		}
	}

	if payload.OrdersType != nil {
		order.OrdersType = internalmessages.OrdersType(*payload.OrdersType)
	}

	// if the order has amended order documents and it has not been previously acknowledged record the current timestamp
	if payload.OrdersAcknowledgement != nil && *payload.OrdersAcknowledgement && existingOrder.UploadedAmendedOrdersID != nil && existingOrder.AmendedOrdersAcknowledgedAt == nil {
		acknowledgedAt := time.Now()
		order.AmendedOrdersAcknowledgedAt = &acknowledgedAt
	}

	if payload.DependentsAuthorized != nil {
		order.Entitlement.DependentsAuthorized = payload.DependentsAuthorized
	}

	if payload.Grade != nil {
		order.Grade = (*internalmessages.OrderPayGrade)(payload.Grade)
		// Calculate new DBWeightAuthorized based on the new grade
		weightAllotment, err := waf.GetWeightAllotment(appCtx, string(*order.Grade), order.OrdersType)
		if err != nil {
			return models.Order{}, err
		}
		weight := weightAllotment.TotalWeightSelf
		// Payload does not have this information, retrieve dependents from the existing order
		if existingOrder.HasDependents && *order.Entitlement.DependentsAuthorized {
			// Only utilize dependent weight authorized if dependents are both present and authorized
			weight = weightAllotment.TotalWeightSelfPlusDependents
		}
		order.Entitlement.DBAuthorizedWeight = &weight
	}

	return order, nil
}

func (f *orderUpdater) amendedOrder(appCtx appcontext.AppContext, userID uuid.UUID, order models.Order, file io.ReadCloser, filename string, storer storage.FileStorer, uploadType models.UploadType) (models.UserUpload, string, *validate.Errors, error) {
	// If Order does not have a Document for amended orders uploads, then create a new one
	var err error
	savedAmendedOrdersDoc := order.UploadedAmendedOrders
	if order.UploadedAmendedOrders == nil {
		amendedOrdersDoc := &models.Document{
			ServiceMemberID: order.ServiceMemberID,
		}
		savedAmendedOrdersDoc, err = f.saveDocumentForAmendedOrder(appCtx, amendedOrdersDoc)
		if err != nil {
			return models.UserUpload{}, "", nil, err
		}

		// save new UploadedAmendedOrdersID (document ID) to orders
		order.UploadedAmendedOrders = savedAmendedOrdersDoc
		order.UploadedAmendedOrdersID = &savedAmendedOrdersDoc.ID
		_, _, err = f.updateOrder(appCtx, order)
		if err != nil {
			return models.UserUpload{}, "", nil, err
		}
	}

	// Create new user upload for amended order
	var userUpload *models.UserUpload
	var verrs *validate.Errors
	var url string
	userUpload, url, verrs, err = uploader.CreateUserUploadForDocumentWrapper(
		appCtx,
		userID,
		storer,
		file,
		filename,
		uploader.MaxCustomerUserUploadFileSizeLimit,
		uploader.AllowedTypesServiceMember,
		&savedAmendedOrdersDoc.ID,
		uploadType,
	)

	if verrs.HasAny() || err != nil {
		return models.UserUpload{}, "", verrs, err
	}

	order.UploadedAmendedOrders.UserUploads = append(order.UploadedAmendedOrders.UserUploads, *userUpload)

	return *userUpload, url, nil, nil
}

func orderFromCounselingPayload(appCtx appcontext.AppContext, existingOrder models.Order, payload ghcmessages.CounselingUpdateOrderPayload) (models.Order, error) {
	order := existingOrder
	waf := entitlements.NewWeightAllotmentFetcher()

	// update order origin duty location
	if payload.OriginDutyLocationID != nil {
		originDutyLocationID := uuid.FromStringOrNil(payload.OriginDutyLocationID.String())
		order.OriginDutyLocationID = &originDutyLocationID
	}

	if payload.NewDutyLocationID != nil {
		newDutyLocationID := uuid.FromStringOrNil(payload.NewDutyLocationID.String())
		order.NewDutyLocationID = newDutyLocationID
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

	if payload.Sac.Present {
		if payload.Sac.Value != nil && *payload.Sac.Value == "" {
			order.SAC = nil
		} else {
			order.SAC = payload.Sac.Value
		}
	}

	if payload.Tac != nil {
		normalizedTac := strings.ToUpper(*payload.Tac)
		order.TAC = &normalizedTac
	}

	if payload.NtsSac.Present {
		if payload.NtsSac.Value != nil && *payload.NtsSac.Value == "" {
			order.NtsSAC = nil
		} else {
			order.NtsSAC = payload.NtsSac.Value
		}
	}

	if payload.NtsTac.Present {
		if payload.NtsTac.Value != nil && *payload.NtsTac.Value != "" {
			normalizedNtsTac := strings.ToUpper(*payload.NtsTac.Value)
			order.NtsTAC = &normalizedNtsTac
		} else {
			order.NtsTAC = nil
		}
	}

	if payload.OrdersType != nil {
		order.OrdersType = internalmessages.OrdersType(*payload.OrdersType)
	}

	if payload.DependentsAuthorized != nil {
		order.Entitlement.DependentsAuthorized = payload.DependentsAuthorized
	}

	if payload.Grade != nil {
		order.Grade = (*internalmessages.OrderPayGrade)(payload.Grade)
		// Calculate new DBWeightAuthorized based on the new grade
		weightAllotment, err := waf.GetWeightAllotment(appCtx, string(*order.Grade), order.OrdersType)
		if err != nil {
			return models.Order{}, err
		}
		weight := weightAllotment.TotalWeightSelf
		// Payload does not have this information, retrieve dependents from the existing order
		if existingOrder.HasDependents && *order.Entitlement.DependentsAuthorized {
			// Only utilize dependent weight authorized if dependents are both present and authorized
			weight = weightAllotment.TotalWeightSelfPlusDependents
		}
		order.Entitlement.DBAuthorizedWeight = &weight
	}

	if payload.HasDependents != nil {
		order.HasDependents = *payload.HasDependents
	}

	return order, nil
}

func allowanceFromTOOPayload(appCtx appcontext.AppContext, existingOrder models.Order, payload ghcmessages.UpdateAllowancePayload) (models.Order, error) {
	order := existingOrder
	waf := entitlements.NewWeightAllotmentFetcher()

	if payload.ProGearWeight != nil {
		order.Entitlement.ProGearWeight = int(*payload.ProGearWeight)
	}

	if payload.ProGearWeightSpouse != nil {
		order.Entitlement.ProGearWeightSpouse = int(*payload.ProGearWeightSpouse)
	}

	if payload.GunSafeWeight != nil {
		order.Entitlement.GunSafeWeight = int(*payload.GunSafeWeight)
	}

	if payload.RequiredMedicalEquipmentWeight != nil {
		order.Entitlement.RequiredMedicalEquipmentWeight = int(*payload.RequiredMedicalEquipmentWeight)
	}

	// branch for service member
	if payload.Agency != nil {
		order.ServiceMember.Affiliation = (*models.ServiceMemberAffiliation)(payload.Agency)
	}

	// grade
	if payload.Grade != nil {
		grade := internalmessages.OrderPayGrade(*payload.Grade)
		order.Grade = &grade
	}

	// Calculate new DBWeightAuthorized based on the new grade
	weightAllotment, err := waf.GetWeightAllotment(appCtx, string(*order.Grade), order.OrdersType)
	if err != nil {
		return models.Order{}, err
	}
	weight := weightAllotment.TotalWeightSelf
	// Payload does not have this information, retrieve dependents from the existing order
	if existingOrder.HasDependents && *order.Entitlement.DependentsAuthorized {
		// Only utilize dependent weight authorized if dependents are both present and authorized
		weight = weightAllotment.TotalWeightSelfPlusDependents
	}
	order.Entitlement.DBAuthorizedWeight = &weight

	if payload.OrganizationalClothingAndIndividualEquipment != nil {
		order.Entitlement.OrganizationalClothingAndIndividualEquipment = *payload.OrganizationalClothingAndIndividualEquipment
	}

	if payload.StorageInTransit != nil {
		newSITAllowance := int(*payload.StorageInTransit)
		order.Entitlement.StorageInTransit = &newSITAllowance
	}

	if payload.GunSafe != nil {
		order.Entitlement.GunSafe = *payload.GunSafe
	}

	if payload.WeightRestriction != nil {
		weightRestriction := int(*payload.WeightRestriction)
		order.Entitlement.WeightRestriction = &weightRestriction
	} else {
		order.Entitlement.WeightRestriction = nil
	}

	if payload.UbWeightRestriction != nil {
		ubWeightRestriction := int(*payload.UbWeightRestriction)
		order.Entitlement.UBWeightRestriction = &ubWeightRestriction
	} else {
		order.Entitlement.UBWeightRestriction = nil
	}

	if payload.AccompaniedTour != nil {
		order.Entitlement.AccompaniedTour = payload.AccompaniedTour
	}

	if payload.DependentsUnderTwelve != nil {
		order.Entitlement.DependentsUnderTwelve = models.IntPointer(int(*payload.DependentsUnderTwelve))
	}

	if payload.DependentsTwelveAndOver != nil {
		order.Entitlement.DependentsTwelveAndOver = models.IntPointer(int(*payload.DependentsTwelveAndOver))
	}

	if order.OriginDutyLocationID != nil {
		var originDutyLocation models.DutyLocation
		originDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), *order.OriginDutyLocationID)
		if err == nil {
			order.OriginDutyLocationID = &originDutyLocation.ID
			order.OriginDutyLocation = &originDutyLocation
		}
	}

	if order.NewDutyLocationID != uuid.Nil {
		var newDutyLocation models.DutyLocation
		newDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), order.NewDutyLocationID)
		if err == nil {
			order.NewDutyLocationID = newDutyLocation.ID
			order.NewDutyLocation = newDutyLocation
		}
	}

	// Recalculate UB allowance of order entitlement
	hasDepedents := false
	if order.Entitlement != nil && order.Entitlement.DependentsAuthorized != nil && *order.Entitlement.DependentsAuthorized {
		hasDepedents = true
	} else if order.HasDependents {
		hasDepedents = true
	}
	civilianTDYUBAllowance := 0
	if order.Entitlement.UBAllowance != nil {
		civilianTDYUBAllowance = *order.Entitlement.UBAllowance
	}
	if payload.UbAllowance != nil {
		civilianTDYUBAllowance = int(*payload.UbAllowance)
	}
	if order.Entitlement != nil {
		unaccompaniedBaggageAllowance, err := models.GetUBWeightAllowance(appCtx, order.OriginDutyLocation.Address.IsOconus, order.NewDutyLocation.Address.IsOconus, order.ServiceMember.Affiliation, order.Grade, &order.OrdersType, &hasDepedents, order.Entitlement.AccompaniedTour, order.Entitlement.DependentsUnderTwelve, order.Entitlement.DependentsTwelveAndOver, &civilianTDYUBAllowance)
		if err == nil {
			order.Entitlement.UBAllowance = &unaccompaniedBaggageAllowance
		}
	}

	return order, nil
}
func allowanceFromCounselingPayload(appCtx appcontext.AppContext, existingOrder models.Order, payload ghcmessages.CounselingUpdateAllowancePayload) (models.Order, error) {
	order := existingOrder
	waf := entitlements.NewWeightAllotmentFetcher()

	if payload.ProGearWeight != nil {
		order.Entitlement.ProGearWeight = int(*payload.ProGearWeight)
	}

	if payload.ProGearWeightSpouse != nil {
		order.Entitlement.ProGearWeightSpouse = int(*payload.ProGearWeightSpouse)
	}

	if payload.GunSafeWeight != nil {
		order.Entitlement.GunSafeWeight = int(*payload.GunSafeWeight)
	}

	if payload.RequiredMedicalEquipmentWeight != nil {
		order.Entitlement.RequiredMedicalEquipmentWeight = int(*payload.RequiredMedicalEquipmentWeight)
	}

	// branch for service member
	if payload.Agency != nil {
		order.ServiceMember.Affiliation = (*models.ServiceMemberAffiliation)(payload.Agency)
	}

	// grade
	if payload.Grade != nil {
		grade := internalmessages.OrderPayGrade(*payload.Grade)
		order.Grade = &grade
	}

	// Calculate new DBWeightAuthorized based on the new grade
	weightAllotment, err := waf.GetWeightAllotment(appCtx, string(*order.Grade), order.OrdersType)
	if err != nil {
		return models.Order{}, err
	}
	weight := weightAllotment.TotalWeightSelf
	// Payload does not have this information, retrieve dependents from the existing order
	if existingOrder.HasDependents && *order.Entitlement.DependentsAuthorized {
		// Only utilize dependent weight authorized if dependents are both present and authorized
		weight = weightAllotment.TotalWeightSelfPlusDependents
	}
	order.Entitlement.DBAuthorizedWeight = &weight

	if payload.OrganizationalClothingAndIndividualEquipment != nil {
		order.Entitlement.OrganizationalClothingAndIndividualEquipment = *payload.OrganizationalClothingAndIndividualEquipment
	}

	if payload.StorageInTransit != nil {
		newSITAllowance := int(*payload.StorageInTransit)
		order.Entitlement.StorageInTransit = &newSITAllowance
	}

	if payload.GunSafe != nil {
		order.Entitlement.GunSafe = *payload.GunSafe
	}

	if payload.WeightRestriction != nil {
		weightRestriction := int(*payload.WeightRestriction)
		order.Entitlement.WeightRestriction = &weightRestriction
	} else {
		order.Entitlement.WeightRestriction = nil
	}

	if payload.UbWeightRestriction != nil {
		ubWeightRestriction := int(*payload.UbWeightRestriction)
		order.Entitlement.UBWeightRestriction = &ubWeightRestriction
	} else {
		order.Entitlement.UBWeightRestriction = nil
	}

	if payload.AccompaniedTour != nil {
		order.Entitlement.AccompaniedTour = payload.AccompaniedTour
	}

	if payload.DependentsUnderTwelve != nil {
		order.Entitlement.DependentsUnderTwelve = models.IntPointer(int(*payload.DependentsUnderTwelve))
	}

	if payload.DependentsTwelveAndOver != nil {
		order.Entitlement.DependentsTwelveAndOver = models.IntPointer(int(*payload.DependentsTwelveAndOver))
	}

	if order.OriginDutyLocationID != nil {
		var originDutyLocation models.DutyLocation
		originDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), *order.OriginDutyLocationID)
		if err == nil {
			order.OriginDutyLocationID = &originDutyLocation.ID
			order.OriginDutyLocation = &originDutyLocation
		}
	}

	if order.NewDutyLocationID != uuid.Nil {
		var newDutyLocation models.DutyLocation
		newDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), order.NewDutyLocationID)
		if err != nil {
			return models.Order{}, err
		}
		order.NewDutyLocationID = newDutyLocation.ID
		order.NewDutyLocation = newDutyLocation
	}

	// Recalculate UB allowance of order entitlement
	civilianTDYUBAllowance := 0
	if order.Entitlement.UBAllowance != nil {
		civilianTDYUBAllowance = *order.Entitlement.UBAllowance
	}
	if payload.UbAllowance != nil {
		civilianTDYUBAllowance = int(*payload.UbAllowance)
	}
	if order.Entitlement != nil {
		unaccompaniedBaggageAllowance, err := models.GetUBWeightAllowance(appCtx, order.OriginDutyLocation.Address.IsOconus, order.NewDutyLocation.Address.IsOconus, order.ServiceMember.Affiliation, order.Grade, &order.OrdersType, order.Entitlement.DependentsAuthorized, order.Entitlement.AccompaniedTour, order.Entitlement.DependentsUnderTwelve, order.Entitlement.DependentsTwelveAndOver, &civilianTDYUBAllowance)
		if err != nil {
			return models.Order{}, err
		}
		order.Entitlement.UBAllowance = &unaccompaniedBaggageAllowance
	}

	return order, nil
}
func (f *orderUpdater) saveDocumentForAmendedOrder(appCtx appcontext.AppContext, doc *models.Document) (*models.Document, error) {
	var docID uuid.UUID
	if doc != nil {
		docID = doc.ID
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var verrs *validate.Errors
		var err error
		verrs, err = txnAppCtx.DB().ValidateAndSave(doc)
		return handleError(docID, verrs, err)
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return doc, nil
}

func (f *orderUpdater) updateOrder(appCtx appcontext.AppContext, order models.Order, checks ...Validator) (*models.Order, uuid.UUID, error) {
	var returnedOrder *models.Order
	var err error

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		returnedOrder, err = updateOrderInTx(txnAppCtx, order, checks...)
		if err != nil {
			return err
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
	return returnedOrder, moveID, nil
}

func (f *orderUpdater) updateOrderAsTOO(appCtx appcontext.AppContext, order models.Order, checks ...Validator) (*models.Order, uuid.UUID, error) {
	move := order.Moves[0]
	var returnedOrder *models.Order
	var err error

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		returnedOrder, err = updateOrderInTx(txnAppCtx, order, checks...)
		if err != nil {
			return err
		}

		return f.updateMoveInTx(txnAppCtx, move)
	})

	if transactionError != nil {
		return nil, uuid.Nil, transactionError
	}

	return returnedOrder, move.ID, nil
}

func (f *orderUpdater) updateMoveInTx(appCtx appcontext.AppContext, move models.Move) error {
	if move.Status == models.MoveStatusAPPROVALSREQUESTED {
		if _, err := f.moveRouter.ApproveOrRequestApproval(appCtx, move); err != nil {
			return err
		}
	}

	return nil
}

func updateOrderInTx(appCtx appcontext.AppContext, order models.Order, checks ...Validator) (*models.Order, error) {
	var verrs *validate.Errors
	var err error

	if verr := ValidateOrder(&order, checks...); verr != nil {
		return nil, verr
	}

	if order.OriginDutyLocationID != nil {
		// TODO refactor to use service objects to fetch duty location
		var originDutyLocation models.DutyLocation
		originDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), *order.OriginDutyLocationID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, apperror.NewNotFoundError(*order.OriginDutyLocationID, "while looking for OriginDutyLocation")
			default:
				return nil, apperror.NewQueryError("DutyLocation", err, "")
			}
		}
		order.OriginDutyLocationID = &originDutyLocation.ID
		order.OriginDutyLocation = &originDutyLocation

		var originDutyLocationGBLOC *string
		if *originDutyLocation.Address.IsOconus {
			originDutyLocationGBLOCOconus, err := models.FetchAddressGbloc(appCtx.DB(), originDutyLocation.Address, order.ServiceMember)
			if err != nil {
				return nil, apperror.NewNotFoundError(originDutyLocation.ID, "while looking for Duty Location Oconus GBLOC")
			}
			originDutyLocationGBLOC = originDutyLocationGBLOCOconus
		} else {
			originDutyLocationGBLOCConus, err2 := models.FetchGBLOCForPostalCode(appCtx.DB(), originDutyLocation.Address.PostalCode)
			if err2 != nil {
				switch err2 {
				case sql.ErrNoRows:
					return nil, apperror.NewNotFoundError(originDutyLocation.ID, "while looking for Duty Location PostalCodeToGBLOC")
				default:
					return nil, apperror.NewQueryError("PostalCodeToGBLOC", err, "")
				}
			}
			originDutyLocationGBLOC = &originDutyLocationGBLOCConus.GBLOC
		}

		order.OriginDutyLocationGBLOC = originDutyLocationGBLOC
	}

	if order.Grade != nil || order.OriginDutyLocationID != nil {
		verrs, err = appCtx.DB().ValidateAndUpdate(&order.ServiceMember)
		if e := handleError(order.ID, verrs, err); e != nil {
			return nil, e
		}
	}

	if order.NewDutyLocationID != uuid.Nil {
		// TODO refactor to use service objects to fetch duty location
		var newDutyLocation models.DutyLocation
		newDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), order.NewDutyLocationID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, apperror.NewNotFoundError(order.NewDutyLocationID, "while looking for NewDutyLocation")
			default:
				return nil, apperror.NewQueryError("DutyLocation", err, "")
			}
		}

		var newDestinationGBLOC *string
		if *newDutyLocation.Address.IsOconus {
			newDestinationGBLOCOconus, err := models.FetchAddressGbloc(appCtx.DB(), newDutyLocation.Address, order.ServiceMember)
			if err != nil {
				return nil, apperror.NewNotFoundError(newDutyLocation.ID, "while looking for DestinationGBLOC Oconus")
			}
			newDestinationGBLOC = newDestinationGBLOCOconus
		} else {
			newDestinationGBLOCConus, err2 := models.FetchGBLOCForPostalCode(appCtx.DB(), newDutyLocation.Address.PostalCode)
			if err2 != nil {
				switch err {
				case sql.ErrNoRows:
					return nil, apperror.NewNotFoundError(order.NewDutyLocationID, "while looking for DestinationGBLOC")
				default:
					return nil, apperror.NewQueryError("DestinationGBLOC", err, "")
				}
			}
			newDestinationGBLOC = &newDestinationGBLOCConus.GBLOC
		}

		order.NewDutyLocationID = newDutyLocation.ID
		order.NewDutyLocation = newDutyLocation
		order.DestinationGBLOC = newDestinationGBLOC
	}

	civilianTDYUBAllowance := 0
	// Recalculate UB allowance of order entitlement
	if order.Entitlement != nil {
		if order.Entitlement.UBAllowance != nil {
			civilianTDYUBAllowance = *order.Entitlement.UBAllowance
		}
		unaccompaniedBaggageAllowance, err := models.GetUBWeightAllowance(appCtx, order.OriginDutyLocation.Address.IsOconus, order.NewDutyLocation.Address.IsOconus, order.ServiceMember.Affiliation, order.Grade, &order.OrdersType, order.Entitlement.DependentsAuthorized, order.Entitlement.AccompaniedTour, order.Entitlement.DependentsUnderTwelve, order.Entitlement.DependentsTwelveAndOver, &civilianTDYUBAllowance)
		if err == nil {
			order.Entitlement.UBAllowance = &unaccompaniedBaggageAllowance
		}
	}

	// update entitlement
	if order.Entitlement != nil {
		verrs, err = appCtx.DB().ValidateAndUpdate(order.Entitlement)
		if e := handleError(order.ID, verrs, err); e != nil {
			return nil, e
		}
	}

	verrs, err = appCtx.DB().ValidateAndUpdate(&order)
	if e := handleError(order.ID, verrs, err); e != nil {
		return nil, e
	}

	// change actual expense reimbursement to 'true' for all PPM shipments if pay grade is civilian
	// if not, do the opposite and make the PPM type INCENTIVE_BASED
	if order.Grade != nil && *order.Grade == models.ServiceMemberGradeCIVILIANEMPLOYEE {
		moves, fetchErr := models.FetchMovesByOrderID(appCtx.DB(), order.ID)
		if fetchErr != nil {
			appCtx.Logger().Error("failure encountered querying for move associated with the order", zap.Error(fetchErr))
		} else {
			var move *models.Move
			for i := range moves {
				if moves[i].OrdersID == order.ID {
					move = &moves[i]
					break
				}
			}
			if move == nil {
				appCtx.Logger().Error("no move found matching order ID", zap.String("orderID", order.ID.String()))
			} else {
				// look at the values and see if the grade is CIVILIAN_EMPLOYEE
				isCivilian := *order.Grade == models.ServiceMemberGradeCIVILIANEMPLOYEE
				reimbursementVal := isCivilian
				var ppmType models.PPMType
				if isCivilian {
					ppmType = models.PPMTypeActualExpense
				} else {
					ppmType = models.PPMTypeIncentiveBased
				}

				for i := range move.MTOShipments {
					shipment := &move.MTOShipments[i]
					if shipment.ShipmentType == models.MTOShipmentTypePPM && shipment.Status == models.MTOShipmentStatusSubmitted && shipment.PPMShipment.PPMType == models.PPMTypeIncentiveBased {
						if shipment.PPMShipment == nil {
							appCtx.Logger().Warn("PPM shipment not found for MTO shipment", zap.String("shipmentID", shipment.ID.String()))
							continue
						}
						shipment.PPMShipment.IsActualExpenseReimbursement = models.BoolPointer(reimbursementVal)
						shipment.PPMShipment.PPMType = ppmType

						if verrs, err := appCtx.DB().ValidateAndUpdate(shipment.PPMShipment); verrs.HasAny() || err != nil {
							msg := "failure saving PPM shipment when updating orders"
							appCtx.Logger().Error(msg, zap.Error(err))
						}
					}
				}
			}
		}
	}

	return &order, nil
}
