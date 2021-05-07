package adminapi

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	notificationsop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/notification"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForNotificationModel(n models.Notification) *adminmessages.Notification {
	return &adminmessages.Notification{
		ID:               handlers.FmtUUID(n.ID),
		ServiceMemberID:  handlers.FmtUUID(n.ServiceMemberID),
		SesMessageID:     handlers.FmtString(n.SESMessageID),
		NotificationType: handlers.FmtString(string(n.NotificationType)),
		CreatedAt:        handlers.FmtDateTime(n.CreatedAt),
		Email:            handlers.FmtString(n.ServiceMember.User.LoginGovEmail),
	}
}

// IndexNotificationsHandler is the index notification handler
type IndexNotificationsHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle does the index notification
func (h IndexNotificationsHandler) Handle(params notificationsop.IndexNotificationsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	queryFilters := h.generateQueryFilters(params.Filter, logger)
	pagination := h.NewPagination(params.Page, params.PerPage)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember.User"),
	}
	associations := query.NewQueryAssociationsPreload(queryAssociations)
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	var notifications []models.Notification
	err := h.ListFetcher.FetchRecordList(&notifications, queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalNotificationsCount, err := h.ListFetcher.FetchRecordCount(&notifications, queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedNotificationsCount := len(notifications)

	payload := make(adminmessages.Notifications, queriedNotificationsCount)

	for i, s := range notifications {
		payload[i] = payloadForNotificationModel(s)
	}

	return notificationsop.NewIndexNotificationsOK().WithContentRange(fmt.Sprintf("notifications %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedNotificationsCount, totalNotificationsCount)).WithPayload(payload)
}

// generateQueryFilters is helper to convert filter params from a json string
// of the form `{"search": "example1@example.com"}` to an array of services.QueryFilter
func (h IndexNotificationsHandler) generateQueryFilters(filters *string, logger handlers.Logger) []services.QueryFilter {
	type Filter struct {
		ServiceMemberID string `json:"service_member_id"`
	}
	f := Filter{}
	var queryFilters []services.QueryFilter
	if filters == nil {
		return queryFilters
	}
	b := []byte(*filters)
	err := json.Unmarshal(b, &f)
	if err != nil {
		fs := fmt.Sprintf("%v", filters)
		logger.Warn("unable to decode param", zap.Error(err),
			zap.String("filters", fs))
	}
	if f.ServiceMemberID != "" {
		queryFilters = append(queryFilters, query.NewQueryFilter("service_member_id", "=", f.ServiceMemberID))
	}

	return queryFilters
}
