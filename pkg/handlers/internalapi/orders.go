package internalapi

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

func payloadForUploadModelFromAmendedOrdersUpload(storer storage.FileStorer, upload models.Upload, url string) (*internalmessages.Upload, error) {
	uploadPayload := &internalmessages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil {
		uploadPayload.Status = string(models.AVStatusPROCESSING)
	} else {
		uploadPayload.Status = string(models.GetAVStatusFromTags(tags))
	}
	return uploadPayload, nil
}

func payloadForOrdersModel(storer storage.FileStorer, order models.Order) (*internalmessages.Orders, error) {
	orderPayload, err := payloads.PayloadForDocumentModel(storer, order.UploadedOrders)
	if err != nil {
		return nil, err
	}

	var amendedOrderPayload *internalmessages.Document
	if order.UploadedAmendedOrders != nil {
		amendedOrderPayload, err = payloads.PayloadForDocumentModel(storer, *order.UploadedAmendedOrders)
		if err != nil {
			return nil, err
		}
	}

	var moves internalmessages.IndexMovesPayload
	for _, move := range order.Moves {
		payload, err := payloadForMoveModel(storer, order, move)
		if err != nil {
			return nil, err
		}
		moves = append(moves, payload)
	}

	var dBAuthorizedWeight *int64
	dBAuthorizedWeight = nil
	var entitlement internalmessages.Entitlement
	if order.Entitlement != nil {
		dBAuthorizedWeight = models.Int64Pointer(int64(*order.Entitlement.AuthorizedWeight()))
		entitlement.ProGear = models.Int64Pointer(int64(order.Entitlement.ProGearWeight))
		entitlement.ProGearSpouse = models.Int64Pointer(int64(order.Entitlement.ProGearWeightSpouse))
		if order.Entitlement.AccompaniedTour != nil {
			entitlement.AccompaniedTour = models.BoolPointer(*order.Entitlement.AccompaniedTour)
		}
		if order.Entitlement.DependentsUnderTwelve != nil {
			entitlement.DependentsUnderTwelve = models.Int64Pointer(int64(*order.Entitlement.DependentsUnderTwelve))
		}
		if order.Entitlement.DependentsTwelveAndOver != nil {
			entitlement.DependentsTwelveAndOver = models.Int64Pointer(int64(*order.Entitlement.DependentsTwelveAndOver))
		}
		if order.Entitlement.UBAllowance != nil {
			entitlement.UbAllowance = models.Int64Pointer(int64(*order.Entitlement.UBAllowance))
		}
		if order.Entitlement.WeightRestriction != nil {
			entitlement.WeightRestriction = models.Int64Pointer(int64(*order.Entitlement.WeightRestriction))
		}
		if order.Entitlement.UBWeightRestriction != nil {
			entitlement.UbWeightRestriction = models.Int64Pointer(int64(*order.Entitlement.UBWeightRestriction))
		}
	}
	var originDutyLocation models.DutyLocation
	originDutyLocation = models.DutyLocation{}
	if order.OriginDutyLocation != nil {
		originDutyLocation = *order.OriginDutyLocation
	}

	var grade internalmessages.OrderPayGrade
	if order.Grade != nil {
		grade = internalmessages.OrderPayGrade(*order.Grade)
	}

	ordersType := order.OrdersType
	payload := &internalmessages.Orders{
		ID:                         handlers.FmtUUID(order.ID),
		CreatedAt:                  handlers.FmtDateTime(order.CreatedAt),
		UpdatedAt:                  handlers.FmtDateTime(order.UpdatedAt),
		ServiceMemberID:            handlers.FmtUUID(order.ServiceMemberID),
		IssueDate:                  handlers.FmtDate(order.IssueDate),
		ReportByDate:               handlers.FmtDate(order.ReportByDate),
		OrdersType:                 &ordersType,
		OrdersTypeDetail:           order.OrdersTypeDetail,
		OriginDutyLocation:         payloadForDutyLocationModel(originDutyLocation),
		OriginDutyLocationGbloc:    handlers.FmtStringPtr(order.OriginDutyLocationGBLOC),
		Grade:                      &grade,
		NewDutyLocation:            payloadForDutyLocationModel(order.NewDutyLocation),
		HasDependents:              handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:           handlers.FmtBool(order.SpouseHasProGear),
		UploadedOrders:             orderPayload,
		UploadedAmendedOrders:      amendedOrderPayload,
		OrdersNumber:               order.OrdersNumber,
		Moves:                      moves,
		Tac:                        order.TAC,
		Sac:                        order.SAC,
		DepartmentIndicator:        (*internalmessages.DeptIndicator)(order.DepartmentIndicator),
		Status:                     internalmessages.OrdersStatus(order.Status),
		AuthorizedWeight:           dBAuthorizedWeight,
		Entitlement:                &entitlement,
		ProvidesServicesCounseling: originDutyLocation.ProvidesServicesCounseling,
	}

	return payload, nil
}

// CreateOrdersHandler creates new orders via POST /orders
type CreateOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.CreateOrders

			serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			originDutyLocationID, err := uuid.FromString(payload.OriginDutyLocationID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			originDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), originDutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			newDutyLocationID, err := uuid.FromString(payload.NewDutyLocationID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			newDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), newDutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			var newDutyLocationGBLOC *string
			if *newDutyLocation.Address.IsOconus {
				newDutyLocationGBLOCOconus, err := models.FetchAddressGbloc(appCtx.DB(), newDutyLocation.Address, serviceMember)
				if err != nil {
					return nil, apperror.NewNotFoundError(newDutyLocation.ID, "while looking for New Duty Location Oconus GBLOC")
				}
				newDutyLocationGBLOC = newDutyLocationGBLOCOconus
			} else {
				newDutyLocationGBLOCConus, err := models.FetchGBLOCForPostalCode(appCtx.DB(), newDutyLocation.Address.PostalCode)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						return nil, apperror.NewNotFoundError(newDutyLocation.ID, "while looking for New Duty Location PostalCodeToGBLOC")
					default:
						err = apperror.NewBadDataError("New duty location GBLOC cannot be verified")
						appCtx.Logger().Error(err.Error())
						return handlers.ResponseForError(appCtx.Logger(), err), err
					}
				}
				newDutyLocationGBLOC = &newDutyLocationGBLOCConus.GBLOC
			}

			var dependentsTwelveAndOver *int
			var dependentsUnderTwelve *int
			if payload.DependentsTwelveAndOver != nil {
				// Convert from int64 to int
				dependentsTwelveAndOver = models.IntPointer(int(*payload.DependentsTwelveAndOver))
			}
			if payload.DependentsUnderTwelve != nil {
				// Convert from int64 to int
				dependentsUnderTwelve = models.IntPointer(int(*payload.DependentsUnderTwelve))
			}

			var originDutyLocationGBLOC *string
			if *originDutyLocation.Address.IsOconus {
				originDutyLocationGBLOCOconus, err := models.FetchAddressGbloc(appCtx.DB(), originDutyLocation.Address, serviceMember)
				if err != nil {
					return nil, apperror.NewNotFoundError(originDutyLocation.ID, "while looking for Origin Duty Location Oconus GBLOC")
				}
				originDutyLocationGBLOC = originDutyLocationGBLOCOconus
			} else {
				originDutyLocationGBLOCConus, err := models.FetchGBLOCForPostalCode(appCtx.DB(), originDutyLocation.Address.PostalCode)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						return nil, apperror.NewNotFoundError(originDutyLocation.ID, "while looking for Origin Duty Location PostalCodeToGBLOC")
					default:
						return nil, apperror.NewQueryError("PostalCodeToGBLOC", err, "")
					}
				}
				originDutyLocationGBLOC = &originDutyLocationGBLOCConus.GBLOC
			}

			grade := payload.Grade

			if payload.OrdersType == nil {
				errMsg := "missing required field: OrdersType"
				return handlers.ResponseForError(appCtx.Logger(), errors.New(errMsg)), apperror.NewBadDataError("missing required field: OrdersType")
			}

			// Calculate the entitlement for the order
			ordersType := payload.OrdersType
			waf := entitlements.NewWeightAllotmentFetcher()
			weightAllotment, err := waf.GetWeightAllotment(appCtx, string(*grade), *ordersType)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			weight := weightAllotment.TotalWeightSelf
			if *payload.HasDependents {
				weight = weightAllotment.TotalWeightSelfPlusDependents
			}

			civilianTDYUBAllowance := 0
			if payload.CivilianTdyUbAllowance != nil {
				civilianTDYUBAllowance = int(*payload.CivilianTdyUbAllowance)
			}
			// Calculate UB allowance for the order entitlement
			unaccompaniedBaggageAllowance, err := models.GetUBWeightAllowance(appCtx, originDutyLocation.Address.IsOconus, newDutyLocation.Address.IsOconus, serviceMember.Affiliation, grade, payload.OrdersType, payload.HasDependents, payload.AccompaniedTour, dependentsUnderTwelve, dependentsTwelveAndOver, &civilianTDYUBAllowance)
			if err == nil {
				weightAllotment.UnaccompaniedBaggageAllowance = unaccompaniedBaggageAllowance
			}

			maxGunSafeWeightAllowance, err := models.GetMaxGunSafeAllowance(appCtx)
			if err == nil {
				weightAllotment.GunSafeWeight = maxGunSafeWeightAllowance
			}

			// Assign default SIT allowance based on customer type.
			// We only have service members right now, but once we introduce more, this logic will have to change.
			sitDaysAllowance := models.DefaultServiceMemberSITDaysAllowance

			entitlement := models.Entitlement{
				DependentsAuthorized:    payload.HasDependents,
				AccompaniedTour:         payload.AccompaniedTour,
				DependentsUnderTwelve:   dependentsUnderTwelve,
				DependentsTwelveAndOver: dependentsTwelveAndOver,
				DBAuthorizedWeight:      models.IntPointer(weight),
				StorageInTransit:        models.IntPointer(sitDaysAllowance),
				ProGearWeight:           weightAllotment.ProGearWeight,
				ProGearWeightSpouse:     weightAllotment.ProGearWeightSpouse,
				UBAllowance:             &weightAllotment.UnaccompaniedBaggageAllowance,
				GunSafeWeight:           weightAllotment.GunSafeWeight,
			}

			/*
				IF you get that to work you'll still have to add conditionals for all the places the entitlement is used because it
				isn't inheritly clear if it's using the spouse weight or not. So you'll be creating new variables and conditionals
				in move_dats.go, move_weights, and move_submitted, etc
			*/

			verrs, err := appCtx.DB().ValidateAndSave(&entitlement)
			if err != nil {
				appCtx.Logger().Error("Error saving customer entitlement", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			if verrs.HasAny() {
				appCtx.Logger().Error("Validation error saving customer entitlement", zap.Any("errors", verrs.Errors))
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, nil), nil
			}

			var deptIndicator *string
			if payload.DepartmentIndicator != nil {
				converted := string(*payload.DepartmentIndicator)
				deptIndicator = &converted
			}

			contractor, err := models.FetchGHCPrimeContractor(appCtx.DB())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			packingAndShippingInstructions := models.InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + models.InstructionsAfterContractNumber
			newOrder, verrs, err := serviceMember.CreateOrder(
				appCtx,
				time.Time(*payload.IssueDate),
				time.Time(*payload.ReportByDate),
				*payload.OrdersType,
				*payload.HasDependents,
				*payload.SpouseHasProGear,
				newDutyLocation,
				payload.OrdersNumber,
				payload.Tac,
				payload.Sac,
				deptIndicator,
				&originDutyLocation,
				grade,
				&entitlement,
				originDutyLocationGBLOC,
				packingAndShippingInstructions,
				newDutyLocationGBLOC,
			)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			moveOptions := models.MoveOptions{
				Show: models.BoolPointer(true),
			}

			if payload.CounselingOfficeID != nil {
				counselingOffice, err := uuid.FromString(payload.CounselingOfficeID.String())
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
				moveOptions.CounselingOfficeID = &counselingOffice
			}

			newMove, verrs, err := newOrder.CreateNewMove(appCtx.DB(), moveOptions)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}
			newOrder.Moves = append(newOrder.Moves, *newMove)

			orderPayload, err := payloadForOrdersModel(h.FileStorer(), newOrder)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload), nil
		})
}

// ShowOrdersHandler returns orders for a user and order ID
type ShowOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowOrdersHandler) Handle(params ordersop.ShowOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			orderID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			order, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), orderID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewShowOrdersOK().WithPayload(orderPayload), nil
		})
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			orderID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			order, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), orderID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := params.UpdateOrders
			dutyLocationID, err := uuid.FromString(payload.NewDutyLocationID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			dutyLocation, err := models.FetchDutyLocation(appCtx.DB(), dutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			var newDutyLocationGBLOC *string
			if *dutyLocation.Address.IsOconus {
				newDutyLocationGBLOCOconus, err := models.FetchAddressGbloc(appCtx.DB(), dutyLocation.Address, order.ServiceMember)
				if err != nil {
					return nil, apperror.NewNotFoundError(dutyLocation.ID, "while looking for New Duty Location Oconus GBLOC")
				}
				newDutyLocationGBLOC = newDutyLocationGBLOCOconus
			} else {
				newDutyLocationGBLOCConus, err := models.FetchGBLOCForPostalCode(appCtx.DB(), dutyLocation.Address.PostalCode)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						return nil, apperror.NewNotFoundError(dutyLocation.ID, "while looking for New Duty Location PostalCodeToGBLOC")
					default:
						err = apperror.NewBadDataError("New duty location GBLOC cannot be verified")
						appCtx.Logger().Error(err.Error())
						return handlers.ResponseForError(appCtx.Logger(), err), err
					}
				}
				newDutyLocationGBLOC = &newDutyLocationGBLOCConus.GBLOC
			}

			if payload.OriginDutyLocationID != "" {
				originDutyLocationID, errorOrigin := uuid.FromString(payload.OriginDutyLocationID.String())
				if errorOrigin != nil {
					return handlers.ResponseForError(appCtx.Logger(), errorOrigin), errorOrigin
				}
				originDutyLocation, errorOrigin := models.FetchDutyLocation(appCtx.DB(), originDutyLocationID)
				if errorOrigin != nil {
					return handlers.ResponseForError(appCtx.Logger(), errorOrigin), errorOrigin
				}
				order.OriginDutyLocation = &originDutyLocation
				order.OriginDutyLocationID = &originDutyLocationID

				var originDutyLocationGBLOC *string
				if *originDutyLocation.Address.IsOconus {
					originDutyLocationGBLOCOconus, err := models.FetchAddressGbloc(appCtx.DB(), originDutyLocation.Address, order.ServiceMember)
					if err != nil {
						return handlers.ResponseForError(appCtx.Logger(), err), err
					}
					originDutyLocationGBLOC = originDutyLocationGBLOCOconus
				} else {
					originDutyLocationGBLOCConus, err := models.FetchGBLOCForPostalCode(appCtx.DB(), originDutyLocation.Address.PostalCode)
					if err != nil {
						switch err {
						case sql.ErrNoRows:
							return nil, apperror.NewNotFoundError(originDutyLocation.ID, "while looking for Origin Duty Location PostalCodeToGBLOC")
						default:
							return handlers.ResponseForError(appCtx.Logger(), err), err
						}
					}
					originDutyLocationGBLOC = &originDutyLocationGBLOCConus.GBLOC
				}
				order.OriginDutyLocationGBLOC = originDutyLocationGBLOC

				if payload.MoveID != "" {

					moveID, err := uuid.FromString(payload.MoveID.String())
					if err != nil {
						return handlers.ResponseForError(appCtx.Logger(), err), err
					}
					move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
					if err != nil {
						return handlers.ResponseForError(appCtx.Logger(), err), err
					}
					if originDutyLocation.ProvidesServicesCounseling {
						counselingOfficeID, err := uuid.FromString(payload.CounselingOfficeID.String())
						if err != nil {
							return handlers.ResponseForError(appCtx.Logger(), err), err
						}
						move.CounselingOfficeID = &counselingOfficeID
					} else {
						move.CounselingOfficeID = nil
					}
					verrs, err := models.SaveMoveDependencies(appCtx.DB(), move)
					if err != nil || verrs.HasAny() {
						return handlers.ResponseForError(appCtx.Logger(), err), err
					}
				}
			}

			if payload.OrdersType == nil {
				errMsg := "missing required field: OrdersType"
				return handlers.ResponseForError(appCtx.Logger(), errors.New(errMsg)), apperror.NewBadDataError("missing required field: OrdersType")
			}

			order.OrdersNumber = payload.OrdersNumber
			order.IssueDate = time.Time(*payload.IssueDate)
			order.ReportByDate = time.Time(*payload.ReportByDate)
			order.OrdersType = *payload.OrdersType
			order.OrdersTypeDetail = payload.OrdersTypeDetail
			order.HasDependents = *payload.HasDependents
			order.SpouseHasProGear = *payload.SpouseHasProGear
			order.NewDutyLocationID = dutyLocation.ID
			order.NewDutyLocation = dutyLocation
			order.DestinationGBLOC = newDutyLocationGBLOC
			order.TAC = payload.Tac
			order.SAC = payload.Sac

			serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// Check if the grade or dependents are receiving an update
			if hasEntitlementChanged(order, payload.OrdersType, payload.Grade, payload.DependentsUnderTwelve, payload.DependentsTwelveAndOver, payload.AccompaniedTour) {
				waf := entitlements.NewWeightAllotmentFetcher()
				weightAllotment, err := waf.GetWeightAllotment(appCtx, string(*payload.Grade), *payload.OrdersType)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
				weight := weightAllotment.TotalWeightSelf
				if *payload.HasDependents {
					weight = weightAllotment.TotalWeightSelfPlusDependents
				}

				// Assign default SIT allowance based on customer type.
				// We only have service members right now, but once we introduce more, this logic will have to change.
				sitDaysAllowance := models.DefaultServiceMemberSITDaysAllowance
				var dependentsTwelveAndOver *int
				var dependentsUnderTwelve *int
				if payload.DependentsTwelveAndOver != nil {
					// Convert from int64 to int
					dependentsTwelveAndOver = models.IntPointer(int(*payload.DependentsTwelveAndOver))
				}
				if payload.DependentsUnderTwelve != nil {
					// Convert from int64 to int
					dependentsUnderTwelve = models.IntPointer(int(*payload.DependentsUnderTwelve))
				}
				var grade *internalmessages.OrderPayGrade
				if payload.Grade != nil {
					grade = payload.Grade
				} else {
					grade = order.Grade
				}

				civilianTDYUBAllowance := 0
				if payload.CivilianTdyUbAllowance != nil {
					civilianTDYUBAllowance = int(*payload.CivilianTdyUbAllowance)
				}
				// Calculate UB allowance for the order entitlement
				if order.Entitlement != nil {
					unaccompaniedBaggageAllowance, err := models.GetUBWeightAllowance(appCtx, order.OriginDutyLocation.Address.IsOconus, order.NewDutyLocation.Address.IsOconus, serviceMember.Affiliation, grade, payload.OrdersType, payload.HasDependents, payload.AccompaniedTour, dependentsUnderTwelve, dependentsTwelveAndOver, &civilianTDYUBAllowance)
					if err == nil {
						weightAllotment.UnaccompaniedBaggageAllowance = unaccompaniedBaggageAllowance
					}
				}

				entitlement := models.Entitlement{
					DependentsAuthorized:    payload.HasDependents,
					DBAuthorizedWeight:      models.IntPointer(weight),
					StorageInTransit:        models.IntPointer(sitDaysAllowance),
					ProGearWeight:           weightAllotment.ProGearWeight,
					ProGearWeightSpouse:     weightAllotment.ProGearWeightSpouse,
					DependentsUnderTwelve:   dependentsUnderTwelve,
					DependentsTwelveAndOver: dependentsTwelveAndOver,
					AccompaniedTour:         payload.AccompaniedTour,
					UBAllowance:             &weightAllotment.UnaccompaniedBaggageAllowance,
					GunSafe:                 order.Entitlement.GunSafe,
					GunSafeWeight:           order.Entitlement.GunSafeWeight,
				}

				/*
					IF you get that to work you'll still have to add conditionals for all the places the entitlement is used because it
					isn't inheritly clear if it's using the spouse weight or not. So you'll be creating new variables and conditionals
					in move_dats.go, move_weights, and move_submitted, etc
				*/

				if saveEntitlementErr := appCtx.DB().Save(&entitlement); saveEntitlementErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), saveEntitlementErr), saveEntitlementErr
				}

				order.EntitlementID = &entitlement.ID
				order.Entitlement = &entitlement

				// change actual expense reimbursement to 'true' for all PPM shipments if pay grade is civilian
				// if not, do the opposite and make the PPM type INCENTIVE_BASED
				if payload.Grade != nil && *payload.Grade != *order.Grade {
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
							isCivilian := *payload.Grade == models.ServiceMemberGradeCIVILIANEMPLOYEE
							reimbursementVal := isCivilian
							var ppmType models.PPMType
							// setting the default ppmType
							if isCivilian {
								ppmType = models.PPMTypeActualExpense
							} else {
								ppmType = models.PPMTypeIncentiveBased
							}

							for i := range move.MTOShipments {
								shipment := &move.MTOShipments[i]
								if shipment.ShipmentType == models.MTOShipmentTypePPM {
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

			}
			order.Grade = payload.Grade

			if payload.DepartmentIndicator != nil {
				order.DepartmentIndicator = handlers.FmtString(string(*payload.DepartmentIndicator))
			}

			verrs, err := models.SaveOrder(appCtx.DB(), &order)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload), nil
		})
}

// Helper func for the UpdateOrdersHandler to check if the entitlement has changed from the new payload
// This handles the nil checks and value comparisons outside of the handler func for organization
func hasEntitlementChanged(order models.Order, payloadOrderType *internalmessages.OrdersType, payloadPayGrade *internalmessages.OrderPayGrade, payloadDependentsUnderTwelve *int64, payloadDependentsTwelveAndOver *int64, payloadAccompaniedTour *bool) bool {
	// Check pay grade
	if (order.Grade == nil && payloadPayGrade != nil) || (order.Grade != nil && payloadPayGrade == nil) || (order.Grade != nil && payloadPayGrade != nil && *order.Grade != *payloadPayGrade) {
		return true
	}
	// check orders type
	if (order.OrdersType == "" && payloadOrderType != nil) || (order.OrdersType != "" && payloadPayGrade == nil) || (order.OrdersType != "" && payloadPayGrade != nil && internalmessages.OrderPayGrade(order.OrdersType) != *payloadPayGrade) {
		return true
	}
	// Check entitlement
	if order.Entitlement != nil {
		// Check dependents under twelve
		if (order.Entitlement.DependentsUnderTwelve == nil && payloadDependentsUnderTwelve != nil) ||
			(order.Entitlement.DependentsUnderTwelve != nil && payloadDependentsUnderTwelve == nil) ||
			(order.Entitlement.DependentsUnderTwelve != nil && payloadDependentsUnderTwelve != nil && *order.Entitlement.DependentsUnderTwelve != int(*payloadDependentsUnderTwelve)) {
			return true
		}

		// Check dependents twelve and over
		if (order.Entitlement.DependentsTwelveAndOver == nil && payloadDependentsTwelveAndOver != nil) ||
			(order.Entitlement.DependentsTwelveAndOver != nil && payloadDependentsTwelveAndOver == nil) ||
			(order.Entitlement.DependentsTwelveAndOver != nil && payloadDependentsTwelveAndOver != nil && *order.Entitlement.DependentsTwelveAndOver != int(*payloadDependentsTwelveAndOver)) {
			return true
		}

		// Check accompanied tour
		if (order.Entitlement.AccompaniedTour == nil && payloadAccompaniedTour != nil) ||
			(order.Entitlement.AccompaniedTour != nil && payloadAccompaniedTour == nil) ||
			(order.Entitlement.AccompaniedTour != nil && payloadAccompaniedTour != nil && *order.Entitlement.AccompaniedTour != *payloadAccompaniedTour) {
			return true
		}

	}

	return false
}

// UploadAmendedOrdersHandler uploads amended orders to an order via PATCH /orders/{orderId}/upload_amended_orders
type UploadAmendedOrdersHandler struct {
	handlers.HandlerConfig
	services.OrderUpdater
}

// Handle updates an order to attach amended orders from a request payload
func (h UploadAmendedOrdersHandler) Handle(params ordersop.UploadAmendedOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			file, ok := params.File.(*runtime.File)
			if !ok {
				errMsg := "This should always be a runtime.File, something has changed in go-swagger."

				appCtx.Logger().Error(errMsg)

				return ordersop.NewUploadAmendedOrdersInternalServerError(), nil
			}

			appCtx.Logger().Info(
				"File uploader and size",
				zap.String("userID", appCtx.Session().UserID.String()),
				zap.String("serviceMemberID", appCtx.Session().ServiceMemberID.String()),
				zap.String("officeUserID", appCtx.Session().OfficeUserID.String()),
				zap.String("AdminUserID", appCtx.Session().AdminUserID.String()),
				zap.Int64("size", file.Header.Size),
			)

			orderID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			upload, url, verrs, err := h.OrderUpdater.UploadAmendedOrdersAsCustomer(appCtx, appCtx.Session().UserID, orderID, file.Data, file.Header.Filename, h.FileStorer())

			if verrs.HasAny() || err != nil {
				switch err.(type) {
				case uploader.ErrTooLarge:
					return ordersop.NewUploadAmendedOrdersRequestEntityTooLarge(), err
				case uploader.ErrFile:
					return ordersop.NewUploadAmendedOrdersInternalServerError(), err
				case uploader.ErrFailedToInitUploader:
					return ordersop.NewUploadAmendedOrdersInternalServerError(), err
				case apperror.NotFoundError:
					return ordersop.NewUploadAmendedOrdersNotFound(), err
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
				}
			}

			uploadPayload, err := payloadForUploadModelFromAmendedOrdersUpload(h.FileStorer(), upload, url)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewUploadAmendedOrdersCreated().WithPayload(uploadPayload), nil
		})
}
