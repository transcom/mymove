package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

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
	handlers.HandlerConfig
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var notificationsFilterConverters = map[string]func(string) []services.QueryFilter{
	"service_member_id": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("service_member_id", "=", content)}
	},
}

// Handle does the index notification
func (h IndexNotificationsHandler) Handle(params notificationsop.IndexNotificationsParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, notificationsFilterConverters)
	pagination := h.NewPagination(params.Page, params.PerPage)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember.User"),
	}
	associations := query.NewQueryAssociationsPreload(queryAssociations)
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	var notifications []models.Notification
	err := h.ListFetcher.FetchRecordList(appCtx, &notifications, queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	totalNotificationsCount, err := h.ListFetcher.FetchRecordCount(appCtx, &notifications, queryFilters)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	queriedNotificationsCount := len(notifications)

	payload := make(adminmessages.Notifications, queriedNotificationsCount)

	for i, s := range notifications {
		payload[i] = payloadForNotificationModel(s)
	}

	return notificationsop.NewIndexNotificationsOK().WithContentRange(fmt.Sprintf("notifications %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedNotificationsCount, totalNotificationsCount)).WithPayload(payload)
}
