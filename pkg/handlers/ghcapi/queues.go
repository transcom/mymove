package ghcapi

import (
	"fmt"
	"slices"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

// GetMovesQueueHandler returns the moves for the TOO queue user via GET /queues/moves
type GetMovesQueueHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
	services.MoveUnlocker
	services.OfficeUserFetcherPop
}

// FilterOption defines the type for the functional arguments used for private functions in OrderFetcher
type FilterOption func(*pop.Query)

// Handle returns the paginated list of moves for the TOO or HQ user
func (h GetMovesQueueHandler) Handle(params queues.GetMovesQueueParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				(!appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) && !appCtx.Session().Roles.HasRole(roles.RoleTypeHQ)) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with TOO or HQ office role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetMovesQueueForbidden(), forbiddenErr
			}

			ListOrderParams := services.ListOrderParams{
				Branch:                  params.Branch,
				Locator:                 params.Locator,
				Edipi:                   params.Edipi,
				Emplid:                  params.Emplid,
				CustomerName:            params.CustomerName,
				DestinationDutyLocation: params.DestinationDutyLocation,
				OriginDutyLocation:      params.OriginDutyLocation,
				AppearedInTOOAt:         handlers.FmtDateTimePtrToPopPtr(params.AppearedInTooAt),
				RequestedMoveDate:       params.RequestedMoveDate,
				Status:                  params.Status,
				Page:                    params.Page,
				PerPage:                 params.PerPage,
				Sort:                    params.Sort,
				Order:                   params.Order,
				OrderType:               params.OrderType,
				TOOAssignedUser:         params.AssignedTo,
				CounselingOffice:        params.CounselingOffice,
			}

			// When no status filter applied, TOO should only see moves with status of New Move, Service Counseling Completed, or Approvals Requested
			if params.Status == nil {
				ListOrderParams.Status = []string{string(models.MoveStatusServiceCounselingCompleted), string(models.MoveStatusAPPROVALSREQUESTED), string(models.MoveStatusSUBMITTED)}
			}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20 for perPage.
			if params.Page == nil {
				ListOrderParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				ListOrderParams.PerPage = models.Int64Pointer(20)
			}

			if params.ViewAsGBLOC != nil && appCtx.Session().Roles.HasRole(roles.RoleTypeHQ) {
				ListOrderParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			moves, count, err := h.OrderFetcher.ListOrders(
				appCtx,
				appCtx.Session().OfficeUserID,
				roles.RoleTypeTOO,
				&ListOrderParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("error fetching list of moves for office user", zap.Error(err))
				return queues.NewGetMovesQueueInternalServerError(), err
			}

			var officeUser models.OfficeUser
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByID(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetServicesCounselingQueueInternalServerError(), err
				}
			}

			privileges, err := models.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges
			officeUser.User.Roles = appCtx.Session().Roles

			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers
			if privileges.HasPrivilege(models.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(models.PrivilegeTypeSafety) {
					officeUsersSafety, err = h.OfficeUserFetcherPop.FetchSafetyMoveOfficeUsersByRoleAndOffice(
						appCtx,
						roles.RoleTypeTOO,
						officeUser.TransportationOfficeID,
					)
					if err != nil {
						appCtx.Logger().
							Error("error fetching safety move office users", zap.Error(err))
						return queues.NewGetMovesQueueInternalServerError(), err
					}
				}
				officeUsers, err = h.OfficeUserFetcherPop.FetchOfficeUsersByRoleAndOffice(
					appCtx,
					roles.RoleTypeTOO,
					officeUser.TransportationOfficeID,
				)
			} else {
				officeUsers = models.OfficeUsers{officeUser}
			}

			if err != nil {
				appCtx.Logger().
					Error("error fetching office users", zap.Error(err))
				return queues.NewGetMovesQueueInternalServerError(), err
			}

			// if the TOO/office user is accessing the queue, we need to unlock move/moves they have locked
			if appCtx.Session().IsOfficeUser() {
				officeUserID := appCtx.Session().OfficeUserID
				for i, move := range moves {
					lockedOfficeUserID := move.LockedByOfficeUserID
					if lockedOfficeUserID != nil && *lockedOfficeUserID == officeUserID {
						copyOfMove := move
						unlockedMove, err := h.UnlockMove(appCtx, &copyOfMove, officeUserID)
						if err != nil {
							return queues.NewGetMovesQueueInternalServerError(), err
						}
						moves[i] = *unlockedMove
					}
				}
				// checking if moves that are NOT in their queue are locked by the user (using search, etc)
				err := h.CheckForLockedMovesAndUnlock(appCtx, officeUserID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("failed to unlock moves for office user ID: %s", officeUserID), zap.Error(err))
				}
			}

			queueMoves := payloads.QueueMoves(moves, officeUsers, nil, officeUser, officeUsersSafety)

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetMovesQueueOK().WithPayload(result), nil
		})
}

// ListMovesHandler lists moves with the option to filter since a particular date. Optimized ver.
type ListPrimeMovesHandler struct {
	handlers.HandlerConfig
	services.MoveTaskOrderFetcher
}

// Handle fetches all moves with the option to filter since a particular date. Optimized version.
func (h ListPrimeMovesHandler) Handle(params queues.ListPrimeMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// adding in moveCode and Id params that are sent in from the UI
			// we will use these params to refine the search in the service object
			searchParams := services.MoveTaskOrderFetcherParams{
				Page:     params.Page,
				PerPage:  params.PerPage,
				MoveCode: params.MoveCode,
				ID:       params.ID,
			}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
			// for perPage.
			if params.Page == nil {
				searchParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				searchParams.PerPage = models.Int64Pointer(20)
			}

			mtos, count, err := h.MoveTaskOrderFetcher.ListNewPrimeMoveTaskOrders(appCtx, &searchParams)

			if err != nil {
				appCtx.Logger().Error("Unexpected error while fetching moves:", zap.Error(err))
				return queues.NewListPrimeMovesInternalServerError(), err
			}

			queueMoves := payloads.ListMoves(&mtos)

			result := ghcmessages.ListPrimeMovesResult{
				Page:       *searchParams.Page,
				PerPage:    *searchParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: queueMoves,
			}

			return queues.NewListPrimeMovesOK().WithPayload(&result), nil

		})
}

// GetPaymentRequestsQueueHandler returns the payment requests for the TIO queue user via GET /queues/payment-requests
type GetPaymentRequestsQueueHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestListFetcher
	services.MoveUnlocker
	services.OfficeUserFetcherPop
}

// Handle returns the paginated list of payment requests for the TIO user
func (h GetPaymentRequestsQueueHandler) Handle(
	params queues.GetPaymentRequestsQueueParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				(!appCtx.Session().Roles.HasRole(roles.RoleTypeTIO) && !appCtx.Session().Roles.HasRole(roles.RoleTypeHQ)) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with TIO or HQ office role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetPaymentRequestsQueueForbidden(), forbiddenErr
			}

			listPaymentRequestParams := services.FetchPaymentRequestListParams{
				Branch:                  params.Branch,
				Locator:                 params.Locator,
				Edipi:                   params.Edipi,
				Emplid:                  params.Emplid,
				CustomerName:            params.CustomerName,
				DestinationDutyLocation: params.DestinationDutyLocation,
				Status:                  params.Status,
				Page:                    params.Page,
				PerPage:                 params.PerPage,
				SubmittedAt:             handlers.FmtDateTimePtrToPopPtr(params.SubmittedAt),
				Sort:                    params.Sort,
				Order:                   params.Order,
				OriginDutyLocation:      params.OriginDutyLocation,
				OrderType:               params.OrderType,
				TIOAssignedUser:         params.AssignedTo,
				CounselingOffice:        params.CounselingOffice,
			}

			listPaymentRequestParams.Status = []string{string(models.QueuePaymentRequestPaymentRequested)}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
			// for perPage.
			if params.Page == nil {
				listPaymentRequestParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				listPaymentRequestParams.PerPage = models.Int64Pointer(20)
			}

			if params.ViewAsGBLOC != nil && appCtx.Session().Roles.HasRole(roles.RoleTypeHQ) {
				listPaymentRequestParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			paymentRequests, count, err := h.FetchPaymentRequestList(
				appCtx,
				appCtx.Session().OfficeUserID,
				&listPaymentRequestParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("payment requests queue", zap.String("office_user_id", appCtx.Session().OfficeUserID.String()), zap.Error(err))
				return queues.NewGetPaymentRequestsQueueInternalServerError(), err
			}

			var officeUser models.OfficeUser
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByID(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetServicesCounselingQueueInternalServerError(), err
				}
			}

			privileges, err := models.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges
			officeUser.User.Roles = appCtx.Session().Roles

			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers

			if privileges.HasPrivilege(models.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(models.PrivilegeTypeSafety) {
					officeUsersSafety, err = h.OfficeUserFetcherPop.FetchSafetyMoveOfficeUsersByRoleAndOffice(
						appCtx,
						roles.RoleTypeTIO,
						officeUser.TransportationOfficeID,
					)
					if err != nil {
						appCtx.Logger().
							Error("error fetching safety move office users", zap.Error(err))
						return queues.NewGetMovesQueueInternalServerError(), err
					}
				}
				officeUsers, err = h.OfficeUserFetcherPop.FetchOfficeUsersByRoleAndOffice(
					appCtx,
					roles.RoleTypeTIO,
					officeUser.TransportationOfficeID,
				)
			}

			if err != nil {
				appCtx.Logger().
					Error("error fetching office users", zap.Error(err))
				return queues.NewGetPaymentRequestsQueueInternalServerError(), err
			}

			// if this TIO/office user is accessing the queue, we need to unlock move/moves they have locked
			if appCtx.Session().IsOfficeUser() {
				officeUserID := appCtx.Session().OfficeUserID
				for i, pr := range *paymentRequests {
					move := pr.MoveTaskOrder
					lockedOfficeUserID := move.LockedByOfficeUserID
					if lockedOfficeUserID != nil && *lockedOfficeUserID == officeUserID {
						unlockedMove, err := h.UnlockMove(appCtx, &move, officeUserID)
						if err != nil {
							return queues.NewGetMovesQueueInternalServerError(), err
						}
						(*paymentRequests)[i].MoveTaskOrder = *unlockedMove
					}
				}
				// checking if moves that are NOT in their queue are locked by the user (using search, etc)
				err := h.CheckForLockedMovesAndUnlock(appCtx, officeUserID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("failed to unlock moves for office user ID: %s", officeUserID), zap.Error(err))
				}
			}

			queuePaymentRequests := payloads.QueuePaymentRequests(paymentRequests, officeUsers, officeUser, officeUsersSafety)

			result := &ghcmessages.QueuePaymentRequestsResult{
				TotalCount:           int64(count),
				Page:                 int64(*listPaymentRequestParams.Page),
				PerPage:              int64(*listPaymentRequestParams.PerPage),
				QueuePaymentRequests: *queuePaymentRequests,
			}

			return queues.NewGetPaymentRequestsQueueOK().WithPayload(result), nil
		})
}

// GetServicesCounselingQueueHandler returns the moves for the Service Counselor queue user via GET /queues/counselor
type GetServicesCounselingQueueHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
	services.MoveUnlocker
	services.OfficeUserFetcherPop
}

// Handle returns the paginated list of moves for the services counselor
func (h GetServicesCounselingQueueHandler) Handle(
	params queues.GetServicesCounselingQueueParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				(!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) && !appCtx.Session().Roles.HasRole(roles.RoleTypeHQ)) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with Services Counselor or HQ office role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetServicesCounselingQueueForbidden(), forbiddenErr
			}

			ListOrderParams := services.ListOrderParams{
				Branch:                  params.Branch,
				Locator:                 params.Locator,
				Edipi:                   params.Edipi,
				Emplid:                  params.Emplid,
				CustomerName:            params.CustomerName,
				OriginDutyLocation:      params.OriginDutyLocation,
				DestinationDutyLocation: params.DestinationDutyLocation,
				OriginGBLOC:             params.OriginGBLOC,
				SubmittedAt:             handlers.FmtDateTimePtrToPopPtr(params.SubmittedAt),
				RequestedMoveDate:       params.RequestedMoveDate,
				Page:                    params.Page,
				PerPage:                 params.PerPage,
				Sort:                    params.Sort,
				Order:                   params.Order,
				NeedsPPMCloseout:        params.NeedsPPMCloseout,
				PPMType:                 params.PpmType,
				CloseoutInitiated:       handlers.FmtDateTimePtrToPopPtr(params.CloseoutInitiated),
				CloseoutLocation:        params.CloseoutLocation,
				OrderType:               params.OrderType,
				PPMStatus:               params.PpmStatus,
				CounselingOffice:        params.CounselingOffice,
				SCAssignedUser:          params.AssignedTo,
			}

			var requestedPpmStatus models.PPMShipmentStatus
			if params.NeedsPPMCloseout != nil && *params.NeedsPPMCloseout {
				requestedPpmStatus = models.PPMShipmentStatusNeedsCloseout
				ListOrderParams.Status = []string{string(models.MoveStatusAPPROVED), string(models.MoveStatusServiceCounselingCompleted)}
			} else if len(params.Status) == 0 {
				ListOrderParams.Status = []string{string(models.MoveStatusNeedsServiceCounseling)}
			} else {
				ListOrderParams.Status = params.Status
			}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
			// for perPage.
			if params.Page == nil {
				ListOrderParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				ListOrderParams.PerPage = models.Int64Pointer(20)
			}

			if params.ViewAsGBLOC != nil && appCtx.Session().Roles.HasRole(roles.RoleTypeHQ) {
				ListOrderParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			moves, count, err := h.OrderFetcher.ListOrders(
				appCtx,
				appCtx.Session().OfficeUserID,
				roles.RoleTypeServicesCounselor,
				&ListOrderParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("error fetching list of moves for office user", zap.Error(err))
				return queues.NewGetServicesCounselingQueueInternalServerError(), err
			}

			var officeUser models.OfficeUser
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByID(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetServicesCounselingQueueInternalServerError(), err
				}
			}

			privileges, err := models.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges
			officeUser.User.Roles = appCtx.Session().Roles

			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers

			if privileges.HasPrivilege(models.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(models.PrivilegeTypeSafety) {
					officeUsersSafety, err = h.OfficeUserFetcherPop.FetchSafetyMoveOfficeUsersByRoleAndOffice(
						appCtx,
						roles.RoleTypeServicesCounselor,
						officeUser.TransportationOfficeID,
					)
					if err != nil {
						appCtx.Logger().
							Error("error fetching safety move office users", zap.Error(err))
						return queues.NewGetMovesQueueInternalServerError(), err
					}
				}
				officeUsers, err = h.OfficeUserFetcherPop.FetchOfficeUsersByRoleAndOffice(
					appCtx,
					roles.RoleTypeServicesCounselor,
					officeUser.TransportationOfficeID,
				)
			} else {
				officeUsers = models.OfficeUsers{officeUser}
			}

			if err != nil {
				appCtx.Logger().
					Error("error fetching office users", zap.Error(err))
				return queues.NewGetServicesCounselingQueueInternalServerError(), err
			}

			// if the SC/office user is accessing the queue, we need to unlock move/moves they have locked
			if appCtx.Session().IsOfficeUser() {
				officeUserID := appCtx.Session().OfficeUserID
				for i, move := range moves {
					lockedOfficeUserID := move.LockedByOfficeUserID
					if lockedOfficeUserID != nil && *lockedOfficeUserID == officeUserID {
						copyOfMove := move
						unlockedMove, err := h.UnlockMove(appCtx, &copyOfMove, officeUserID)
						if err != nil {
							return queues.NewGetMovesQueueInternalServerError(), err
						}
						moves[i] = *unlockedMove
					}
				}
				// checking if moves that are NOT in their queue are locked by the user (using search, etc)
				err := h.CheckForLockedMovesAndUnlock(appCtx, officeUserID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("failed to unlock moves for office user ID: %s", officeUserID), zap.Error(err))
				}
			}

			queueMoves := payloads.QueueMoves(moves, officeUsers, &requestedPpmStatus, officeUser, officeUsersSafety)

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetServicesCounselingQueueOK().WithPayload(result), nil
		})
}

// GetServicesCounselingOriginListHandler returns the origin list for the Service Counselor user via GET /queues/counselor/origin-list
type GetServicesCounselingOriginListHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
}

// Handle returns the list of origin list for the services counselor
func (h GetServicesCounselingOriginListHandler) Handle(
	params queues.GetServicesCounselingOriginListParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with an office role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetServicesCounselingQueueForbidden(), forbiddenErr
			}

			ListOrderParams := services.ListOrderParams{
				NeedsPPMCloseout: params.NeedsPPMCloseout,
			}

			if params.NeedsPPMCloseout != nil && *params.NeedsPPMCloseout {
				ListOrderParams.Status = []string{string(models.MoveStatusAPPROVED), string(models.MoveStatusServiceCounselingCompleted)}
			} else {
				ListOrderParams.Status = []string{string(models.MoveStatusNeedsServiceCounseling)}
			}

			moves, err := h.OrderFetcher.ListAllOrderLocations(
				appCtx,
				appCtx.Session().OfficeUserID,
				&ListOrderParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("error fetching list of moves for office user", zap.Error(err))
				return queues.NewGetServicesCounselingQueueInternalServerError(), err
			}

			var originLocationList []*ghcmessages.Location
			for _, value := range moves {
				locationString := value.Orders.OriginDutyLocation.Name
				location := ghcmessages.Location{Label: &locationString, Value: &locationString}
				if !slices.Contains(originLocationList, &location) {
					originLocationList = append(originLocationList, &location)
				}
			}

			return queues.NewGetServicesCounselingOriginListOK().WithPayload(originLocationList), nil
		})
}
