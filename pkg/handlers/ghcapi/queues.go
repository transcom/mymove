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
				(!(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeTOO) && !(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ)) {
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
				AssignedTo:              params.AssignedTo,
				CounselingOffice:        params.CounselingOffice,
			}

			var activeRole string
			if params.ActiveRole != nil {
				activeRole = *params.ActiveRole
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

			var officeUser models.OfficeUser
			var assignedGblocs []string
			var err error
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetMovesQueueInternalServerError(), err
				}

				assignedGblocs = models.GetAssignedGBLOCs(officeUser)
			}

			if params.ViewAsGBLOC != nil && ((appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ) || slices.Contains(assignedGblocs, *params.ViewAsGBLOC)) {
				ListOrderParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges

			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers
			if privileges.HasPrivilege(roles.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
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

			moves, count, err := h.OrderFetcher.ListOriginRequestsOrders(
				appCtx,
				appCtx.Session().OfficeUserID,
				&ListOrderParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("error fetching list of moves for office user", zap.Error(err))
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

			queueMoves := payloads.QueueMoves(moves, officeUsers, nil, officeUser, officeUsersSafety, activeRole, string(models.QueueTypeTaskOrder))

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetMovesQueueOK().WithPayload(result), nil
		})
}

// GetDestinationRequestsQueueHandler returns the moves for the TOO queue user via GET /queues/destination-requests
type GetDestinationRequestsQueueHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
	services.MoveUnlocker
	services.OfficeUserFetcherPop
}

// Handle returns the paginated list of moves with destination requests for a TOO user
func (h GetDestinationRequestsQueueHandler) Handle(params queues.GetDestinationRequestsQueueParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				(!(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeTOO)) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with TOO role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetDestinationRequestsQueueForbidden(), forbiddenErr
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
				AssignedTo:              params.AssignedTo,
				CounselingOffice:        params.CounselingOffice,
			}

			var activeRole string
			if params.ActiveRole != nil {
				activeRole = *params.ActiveRole
			}
			// we only care about moves in APPROVALS REQUESTED status
			if params.Status == nil {
				ListOrderParams.Status = []string{string(models.MoveStatusAPPROVALSREQUESTED)}
			}

			// default pagination values
			if params.Page == nil {
				ListOrderParams.Page = models.Int64Pointer(1)
			}
			if params.PerPage == nil {
				ListOrderParams.PerPage = models.Int64Pointer(20)
			}

			var officeUser models.OfficeUser
			var assignedGblocs []string
			var err error
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetMovesQueueInternalServerError(), err
				}

				assignedGblocs = models.GetAssignedGBLOCs(officeUser)
			}

			if params.ViewAsGBLOC != nil && ((appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ) || slices.Contains(assignedGblocs, *params.ViewAsGBLOC)) {
				ListOrderParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			moves, count, err := h.OrderFetcher.ListDestinationRequestsOrders(
				appCtx,
				appCtx.Session().OfficeUserID,
				roles.RoleTypeTOO,
				&ListOrderParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("error fetching destinaton queue for office user", zap.Error(err))
				return queues.NewGetDestinationRequestsQueueInternalServerError(), err
			}

			privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges
			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers
			if privileges.HasPrivilege(roles.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
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
				return queues.NewGetDestinationRequestsQueueInternalServerError(), err
			}

			// if the TOO is accessing the queue, we need to unlock move/moves they have locked
			if appCtx.Session().IsOfficeUser() {
				officeUserID := appCtx.Session().OfficeUserID
				for i, move := range moves {
					lockedOfficeUserID := move.LockedByOfficeUserID
					if lockedOfficeUserID != nil && *lockedOfficeUserID == officeUserID {
						copyOfMove := move
						unlockedMove, err := h.UnlockMove(appCtx, &copyOfMove, officeUserID)
						if err != nil {
							return queues.NewGetDestinationRequestsQueueInternalServerError(), err
						}
						moves[i] = *unlockedMove
					}
				}
				err := h.CheckForLockedMovesAndUnlock(appCtx, officeUserID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("failed to unlock moves for office user ID: %s", officeUserID), zap.Error(err))
				}
			}

			queueMoves := payloads.QueueMoves(moves, officeUsers, nil, officeUser, officeUsersSafety, activeRole, string(models.QueueTypeDestinationRequest))

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetDestinationRequestsQueueOK().WithPayload(result), nil
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
				(!(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeTIO) && !(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ)) {
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
				AssignedTo:              params.AssignedTo,
				CounselingOffice:        params.CounselingOffice,
			}

			var activeRole string
			if params.ActiveRole != nil {
				activeRole = *params.ActiveRole
			}

			listPaymentRequestParams.Status = []string{string(models.PaymentRequestStatusPending)}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
			// for perPage.
			if params.Page == nil {
				listPaymentRequestParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				listPaymentRequestParams.PerPage = models.Int64Pointer(20)
			}

			var officeUser models.OfficeUser
			var assignedGblocs []string
			var err error
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetPaymentRequestsQueueInternalServerError(), err
				}

				assignedGblocs = models.GetAssignedGBLOCs(officeUser)
			}

			if params.ViewAsGBLOC != nil && ((appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ) || slices.Contains(assignedGblocs, *params.ViewAsGBLOC)) {
				listPaymentRequestParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges

			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers

			if privileges.HasPrivilege(roles.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
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

			queuePaymentRequests := payloads.QueuePaymentRequests(paymentRequests, officeUsers, officeUser, officeUsersSafety, activeRole)

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

type GetPPMCloseoutQueueHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
	services.MoveUnlocker
	services.OfficeUserFetcherPop
}

// Handle returns the paginated list of moves for the services counselor
func (h GetPPMCloseoutQueueHandler) Handle(
	params queues.GetPPMCloseoutQueueParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			ListOrderParams := services.ListOrderParams{
				Branch:                  params.Branch,
				Locator:                 params.Locator,
				Edipi:                   params.Edipi,
				Emplid:                  params.Emplid,
				CustomerName:            params.CustomerName,
				OriginDutyLocation:      params.OriginDutyLocation,
				DestinationDutyLocation: params.DestinationDutyLocation,
				OriginGBLOC:             params.OriginGBLOC,
				RequestedMoveDate:       params.RequestedMoveDate,
				Page:                    params.Page,
				PerPage:                 params.PerPage,
				Sort:                    params.Sort,
				Order:                   params.Order,
				NeedsPPMCloseout:        params.NeedsPPMCloseout,
				PPMType:                 params.PpmType,
				CloseoutInitiated:       handlers.FmtDateTimePtrToPopPtr(params.CloseoutInitiated), // ListOrderParam SubmittedAt is for filtering moves, closeout is for filtering ppms
				CloseoutLocation:        params.CloseoutLocation,
				OrderType:               params.OrderType,
				PPMStatus:               params.PpmStatus,
				CounselingOffice:        params.CounselingOffice,
				AssignedTo:              params.AssignedTo,
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

			var officeUser models.OfficeUser
			var assignedGblocs []string
			var err error
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetPPMCloseoutQueueInternalServerError(), err
				}

				assignedGblocs = models.GetAssignedGBLOCs(officeUser)
			}

			if params.ViewAsGBLOC != nil && (appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ) || params.ViewAsGBLOC != nil && slices.Contains(assignedGblocs, *params.ViewAsGBLOC) {
				ListOrderParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges

			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers

			if privileges.HasPrivilege(roles.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
					officeUsersSafety, err = h.OfficeUserFetcherPop.FetchSafetyMoveOfficeUsersByRoleAndOffice(
						appCtx,
						roles.RoleTypeServicesCounselor,
						officeUser.TransportationOfficeID,
					)
					if err != nil {
						appCtx.Logger().
							Error("error fetching safety move office users", zap.Error(err))
						return queues.NewGetPPMCloseoutQueueInternalServerError(), err
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
				return queues.NewGetPPMCloseoutQueueInternalServerError(), err
			}

			// Fetch the moves
			moves, count, err := h.ListPPMCloseoutOrders(appCtx, appCtx.Session().OfficeUserID, &ListOrderParams)
			if err != nil {
				appCtx.Logger().
					Error("error fetching list of moves for office user", zap.Error(err))
				return queues.NewGetPPMCloseoutQueueInternalServerError(), err
			}

			// Convert payload
			queueMoves := payloads.QueueMoves(moves, officeUsers, &requestedPpmStatus, officeUser, officeUsersSafety, string(appCtx.Session().ActiveRole.RoleType), string(models.QueueTypeCloseout))

			// if the SC/office user is accessing the queue, we need to unlock move/moves they have locked
			if appCtx.Session().IsOfficeUser() {
				officeUserID := appCtx.Session().OfficeUserID
				for i, move := range moves {
					lockedOfficeUserID := move.LockedByOfficeUserID
					if lockedOfficeUserID != nil && *lockedOfficeUserID == officeUserID {
						copyOfMove := move
						unlockedMove, err := h.UnlockMove(appCtx, &copyOfMove, officeUserID)
						if err != nil {
							return queues.NewGetPPMCloseoutQueueInternalServerError(), err
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

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetPPMCloseoutQueueOK().WithPayload(result), nil
		})
}

// Handle returns the paginated list of moves for the services counselor
func (h GetServicesCounselingQueueHandler) Handle(
	params queues.GetServicesCounselingQueueParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				(!(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeServicesCounselor) && !(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ)) {
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
				AssignedTo:              params.AssignedTo,
			}

			var activeRole string
			if params.ActiveRole != nil {
				activeRole = *params.ActiveRole
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

			var officeUser models.OfficeUser
			var assignedGblocs []string
			var err error
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetServicesCounselingQueueInternalServerError(), err
				}

				assignedGblocs = models.GetAssignedGBLOCs(officeUser)
			}

			if params.ViewAsGBLOC != nil && ((appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ) || slices.Contains(assignedGblocs, *params.ViewAsGBLOC)) {
				ListOrderParams.ViewAsGBLOC = params.ViewAsGBLOC
			}

			privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}
			officeUser.User.Privileges = privileges

			var officeUsers models.OfficeUsers
			var officeUsersSafety models.OfficeUsers

			if privileges.HasPrivilege(roles.PrivilegeTypeSupervisor) {
				if privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
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

			queueType := string(models.QueueTypeCounseling)
			if params.NeedsPPMCloseout != nil && *params.NeedsPPMCloseout {
				queueType = string(models.QueueTypeCloseout)
			}

			queueMoves := payloads.QueueMoves(moves, officeUsers, &requestedPpmStatus, officeUser, officeUsersSafety, activeRole, queueType)

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetServicesCounselingQueueOK().WithPayload(result), nil
		})
}

// GetBulkAssignmentDataHandler returns moves that the supervisor can assign, along with the office users they are able to assign to
type GetBulkAssignmentDataHandler struct {
	handlers.HandlerConfig
	services.OfficeUserFetcherPop
	services.MoveFetcherBulkAssignment
	services.MoveLocker
}

func (h GetBulkAssignmentDataHandler) Handle(
	params queues.GetBulkAssignmentDataParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() {
				err := apperror.NewForbiddenError("not an office user")
				appCtx.Logger().Error("Must be an office user", zap.Error(err))
				return queues.NewGetBulkAssignmentDataUnauthorized(), err
			}

			officeUser, err := h.OfficeUserFetcherPop.FetchOfficeUserByID(appCtx, appCtx.Session().OfficeUserID)
			if err != nil {
				appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
				return queues.NewGetBulkAssignmentDataNotFound(), err
			}

			privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), *officeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
				return queues.NewGetBulkAssignmentDataNotFound(), err
			}

			isSupervisor := privileges.HasPrivilege(roles.PrivilegeTypeSupervisor)
			if !isSupervisor {
				appCtx.Logger().Error("Unauthorized", zap.Error(err))
				return queues.NewGetBulkAssignmentDataUnauthorized(), err
			}

			queueType := params.QueueType
			var officeUserData ghcmessages.BulkAssignmentData

			switch *queueType {
			case string(models.QueueTypeCounseling):
				// fetch the Services Counselors who work at their office
				officeUsers, err := h.OfficeUserFetcherPop.FetchOfficeUsersWithWorkloadByRoleAndOffice(
					appCtx,
					roles.RoleTypeServicesCounselor,
					officeUser.TransportationOfficeID,
					*queueType,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving office users", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}
				// fetch the moves available to be assigned to their office users
				moves, err := h.MoveFetcherBulkAssignment.FetchMovesForBulkAssignmentCounseling(
					appCtx, officeUser.TransportationOffice.Gbloc, officeUser.TransportationOffice.ID,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving moves", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}
				moveIdsToLock := make([]uuid.UUID, len(moves))
				for i, move := range moves {
					moveIdsToLock[i] = move.ID
				}
				err = h.LockMoves(appCtx, moveIdsToLock, officeUser.ID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Failed to lock Services Counseling Queue moves for office user ID: %s", officeUser.ID), zap.Error(err))
				}

				officeUserData = payloads.BulkAssignmentData(appCtx, moves, officeUsers, officeUser.TransportationOffice.ID)
			case string(models.QueueTypeCloseout):
				// fetch the Services Counselors who work at their office
				officeUsers, err := h.OfficeUserFetcherPop.FetchOfficeUsersWithWorkloadByRoleAndOffice(
					appCtx,
					roles.RoleTypeServicesCounselor,
					officeUser.TransportationOfficeID,
					*queueType,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving office users", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}
				// fetch the moves available to be assigned to their office users
				moves, err := h.MoveFetcherBulkAssignment.FetchMovesForBulkAssignmentCloseout(
					appCtx, officeUser.TransportationOffice.Gbloc, officeUser.TransportationOffice.ID,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving moves", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}

				moveIdsToLock := make([]uuid.UUID, len(moves))
				for i, move := range moves {
					moveIdsToLock[i] = move.ID
				}
				err = h.LockMoves(appCtx, moveIdsToLock, officeUser.ID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Failed to lock PPM Closeout Queue moves for office user ID: %s", officeUser.ID), zap.Error(err))
				}

				officeUserData = payloads.BulkAssignmentData(appCtx, moves, officeUsers, officeUser.TransportationOffice.ID)
			case string(models.QueueTypeTaskOrder):
				// fetch the TOOs who work at their office
				officeUsers, err := h.OfficeUserFetcherPop.FetchOfficeUsersWithWorkloadByRoleAndOffice(
					appCtx,
					roles.RoleTypeTOO,
					officeUser.TransportationOfficeID,
					*queueType,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving office users", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}
				// fetch the moves available to be assigned to their office users
				moves, err := h.MoveFetcherBulkAssignment.FetchMovesForBulkAssignmentTaskOrder(
					appCtx, officeUser.TransportationOffice.Gbloc, officeUser.TransportationOffice.ID,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving moves", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}

				moveIdsToLock := make([]uuid.UUID, len(moves))
				for i, move := range moves {
					moveIdsToLock[i] = move.ID
				}
				err = h.LockMoves(appCtx, moveIdsToLock, officeUser.ID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Failed to lock Task Order Queue moves for office user ID: %s", officeUser.ID), zap.Error(err))
				}

				officeUserData = payloads.BulkAssignmentData(appCtx, moves, officeUsers, officeUser.TransportationOffice.ID)
			case string(models.QueueTypePaymentRequest):
				// fetch the TIOs who work at their office
				officeUsers, err := h.OfficeUserFetcherPop.FetchOfficeUsersWithWorkloadByRoleAndOffice(
					appCtx,
					roles.RoleTypeTIO,
					officeUser.TransportationOfficeID,
					*queueType,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving office users", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}
				// fetch the moves available to be assigned to their office users
				moves, err := h.MoveFetcherBulkAssignment.FetchMovesForBulkAssignmentPaymentRequest(
					appCtx, officeUser.TransportationOffice.Gbloc, officeUser.TransportationOffice.ID,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving moves", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}

				moveIdsToLock := make([]uuid.UUID, len(moves))
				for i, move := range moves {
					moveIdsToLock[i] = move.ID
				}
				err = h.LockMoves(appCtx, moveIdsToLock, officeUser.ID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Failed to lock Payment Request Queue moves for office user ID: %s", officeUser.ID), zap.Error(err))
				}

				officeUserData = payloads.BulkAssignmentData(appCtx, moves, officeUsers, officeUser.TransportationOffice.ID)
			case string(models.QueueTypeDestinationRequest):
				// fetch the TOOs who work at their office
				officeUsers, err := h.OfficeUserFetcherPop.FetchOfficeUsersWithWorkloadByRoleAndOffice(
					appCtx,
					roles.RoleTypeTOO,
					officeUser.TransportationOfficeID,
					*queueType,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving office users", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}
				// fetch the moves available to be assigned to their office users
				moves, err := h.MoveFetcherBulkAssignment.FetchMovesForBulkAssignmentDestination(
					appCtx, officeUser.TransportationOffice.Gbloc, officeUser.TransportationOffice.ID,
				)
				if err != nil {
					appCtx.Logger().Error("Error retreiving moves", zap.Error(err))
					return queues.NewGetBulkAssignmentDataInternalServerError(), err
				}

				moveIdsToLock := make([]uuid.UUID, len(moves))
				for i, move := range moves {
					moveIdsToLock[i] = move.ID
				}
				err = h.LockMoves(appCtx, moveIdsToLock, officeUser.ID)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Failed to lock Destination Requests Queue moves for office user ID: %s", officeUser.ID), zap.Error(err))
				}

				officeUserData = payloads.BulkAssignmentData(appCtx, moves, officeUsers, officeUser.TransportationOffice.ID)
			}

			return queues.NewGetBulkAssignmentDataOK().WithPayload(&officeUserData), nil
		})
}

// SaveBulkAssignmentDataHandler saves the bulk assignment data
type SaveBulkAssignmentDataHandler struct {
	handlers.HandlerConfig
	services.OfficeUserFetcherPop
	services.MoveFetcher
	services.MoveAssigner
	services.MoveUnlocker
}

func (h SaveBulkAssignmentDataHandler) Handle(
	params queues.SaveBulkAssignmentDataParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() {
				err := apperror.NewForbiddenError("not an office user")
				appCtx.Logger().Error("Must be an office user", zap.Error(err))
				return queues.NewSaveBulkAssignmentDataUnauthorized(), err
			}

			officeUser, err := h.OfficeUserFetcherPop.FetchOfficeUserByID(appCtx, appCtx.Session().OfficeUserID)
			if err != nil {
				appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
				return queues.NewSaveBulkAssignmentDataNotFound(), err
			}

			privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), *officeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
				return queues.NewSaveBulkAssignmentDataNotFound(), err
			}

			isSupervisor := privileges.HasPrivilege(roles.PrivilegeTypeSupervisor)
			if !isSupervisor {
				appCtx.Logger().Error("Unauthorized", zap.Error(err))
				return queues.NewSaveBulkAssignmentDataUnauthorized(), err
			}

			queueType := params.BulkAssignmentSavePayload.QueueType
			moveData := params.BulkAssignmentSavePayload.MoveData
			userData := params.BulkAssignmentSavePayload.UserData

			// unlock moves that were locked when the bulk assignment modal was opened
			err = h.MoveUnlocker.CheckForLockedMovesAndUnlock(appCtx, officeUser.ID)
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("Failed to unlock moves for office user ID: %s", officeUser.ID), zap.Error(err))
			}

			// fetch the moves available to be assigned to their office users
			movesForAssignment, err := h.MoveFetcher.FetchMovesByIdArray(appCtx, moveData)
			if err != nil {
				appCtx.Logger().Error("Error retreiving moves for assignment", zap.Error(err))
				return queues.NewSaveBulkAssignmentDataInternalServerError(), err
			}

			_, err = h.MoveAssigner.BulkMoveAssignment(appCtx, queueType, userData, movesForAssignment)
			if err != nil {
				appCtx.Logger().Error("Error assigning moves", zap.Error(err))
				return queues.NewGetBulkAssignmentDataInternalServerError(), err
			}

			return queues.NewSaveBulkAssignmentDataNoContent(), nil
		})
}

// GetServicesCounselingOriginListHandler returns the origin list for the Service Counselor user via GET /queues/counselor/origin-list
type GetServicesCounselingOriginListHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
	services.OfficeUserFetcherPop
}

// Handle returns the list of origin list for the services counselor
func (h GetServicesCounselingOriginListHandler) Handle(
	params queues.GetServicesCounselingOriginListParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				!(appCtx.Session().ActiveRole.RoleType == roles.RoleTypeServicesCounselor) {
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

			var officeUser models.OfficeUser
			var assignedGblocs []string
			var err error
			if appCtx.Session().OfficeUserID != uuid.Nil {
				officeUser, err = h.OfficeUserFetcherPop.FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx, appCtx.Session().OfficeUserID)
				if err != nil {
					appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
					return queues.NewGetServicesCounselingOriginListInternalServerError(), err
				}

				assignedGblocs = models.GetAssignedGBLOCs(officeUser)
			}

			if params.ViewAsGBLOC != nil && ((appCtx.Session().ActiveRole.RoleType == roles.RoleTypeHQ) || slices.Contains(assignedGblocs, *params.ViewAsGBLOC)) {
				ListOrderParams.ViewAsGBLOC = params.ViewAsGBLOC
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
