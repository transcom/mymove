package internalapi

import (
	"github.com/transcom/mymove/pkg/auth"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	issueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/issues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForIssueModel(issue models.Issue) internalmessages.IssuePayload {
	issuePayload := internalmessages.IssuePayload{
		CreatedAt:    handlers.FmtDateTime(issue.CreatedAt),
		Description:  swag.String(issue.Description),
		ID:           handlers.FmtUUID(issue.ID),
		UpdatedAt:    handlers.FmtDateTime(issue.UpdatedAt),
		ReporterName: issue.ReporterName,
		DueDate:      (*strfmt.Date)(issue.DueDate),
	}
	return issuePayload
}

// CreateIssueHandler creates a new issue via POST /issue
type CreateIssueHandler struct {
	handlers.HandlerContext
}

// Handle creates a new Issue from a request payload
func (h CreateIssueHandler) Handle(params issueop.CreateIssueParams) middleware.Responder {
	newIssue := models.Issue{
		Description:  *params.CreateIssuePayload.Description,
		ReporterName: params.CreateIssuePayload.ReporterName,
		DueDate:      (*time.Time)(params.CreateIssuePayload.DueDate),
	}
	var response middleware.Responder
	if _, err := h.DB().ValidateAndCreate(&newIssue); err != nil {
		h.Logger().Error("DB Insertion", zap.Error(err))
		// how do I raise an erorr?
		response = issueop.NewCreateIssueBadRequest()
	} else {
		issuePayload := payloadForIssueModel(newIssue)
		response = issueop.NewCreateIssueCreated().WithPayload(&issuePayload)

	}
	return response
}

// IndexIssuesHandler returns a list of all issues
type IndexIssuesHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a list of all issues in the system
func (h IndexIssuesHandler) Handle(params issueop.IndexIssuesParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	var issues models.Issues
	var response middleware.Responder
	if session == nil {
		response = issueop.NewIndexIssuesUnauthorized()
	} else if !session.IsOfficeUser() {
		response = issueop.NewIndexIssuesForbidden()
	} else if err := h.DB().All(&issues); err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		response = issueop.NewIndexIssuesBadRequest()
	} else {
		issuePayloads := make(internalmessages.IndexIssuesPayload, len(issues))
		for i, issue := range issues {
			issuePayload := payloadForIssueModel(issue)
			issuePayloads[i] = &issuePayload
		}
		response = issueop.NewIndexIssuesOK().WithPayload(issuePayloads)
	}
	return response
}
