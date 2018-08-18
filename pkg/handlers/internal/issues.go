package internal

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	issueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/issues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForIssueModel(issue models.Issue) internalmessages.IssuePayload {
	issuePayload := internalmessages.IssuePayload{
		CreatedAt:    utils.FmtDateTime(issue.CreatedAt),
		Description:  swag.String(issue.Description),
		ID:           utils.FmtUUID(issue.ID),
		UpdatedAt:    utils.FmtDateTime(issue.UpdatedAt),
		ReporterName: issue.ReporterName,
		DueDate:      (*strfmt.Date)(issue.DueDate),
	}
	return issuePayload
}

// CreateIssueHandler creates a new issue via POST /issue
type CreateIssueHandler HandlerContext

// Handle creates a new Issue from a request payload
func (h CreateIssueHandler) Handle(params issueop.CreateIssueParams) middleware.Responder {
	newIssue := models.Issue{
		Description:  *params.CreateIssuePayload.Description,
		ReporterName: params.CreateIssuePayload.ReporterName,
		DueDate:      (*time.Time)(params.CreateIssuePayload.DueDate),
	}
	var response middleware.Responder
	if _, err := h.db.ValidateAndCreate(&newIssue); err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		// how do I raise an erorr?
		response = issueop.NewCreateIssueBadRequest()
	} else {
		issuePayload := payloadForIssueModel(newIssue)
		response = issueop.NewCreateIssueCreated().WithPayload(&issuePayload)

	}
	return response
}

// IndexIssuesHandler returns a list of all issues
type IndexIssuesHandler HandlerContext

// Handle retrieves a list of all issues in the system
func (h IndexIssuesHandler) Handle(params issueop.IndexIssuesParams) middleware.Responder {
	var issues models.Issues
	var response middleware.Responder
	if err := h.db.All(&issues); err != nil {
		h.logger.Error("DB Query", zap.Error(err))
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
