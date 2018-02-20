package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	issueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/issues"
	"github.com/transcom/mymove/pkg/gen/internalmodel"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForIssueModel(issue models.Issue) internalmodel.IssuePayload {
	issuePayload := internalmodel.IssuePayload{
		CreatedAt:    fmtDateTime(issue.CreatedAt),
		Description:  swag.String(issue.Description),
		ID:           fmtUUID(issue.ID),
		UpdatedAt:    fmtDateTime(issue.UpdatedAt),
		ReporterName: issue.ReporterName,
		DueDate:      (*strfmt.Date)(issue.DueDate),
	}
	return issuePayload
}

// CreateIssueHandler creates a new issue via POST /issue
func CreateIssueHandler(params issueop.CreateIssueParams) middleware.Responder {
	newIssue := models.Issue{
		Description:  *params.CreateIssuePayload.Description,
		ReporterName: params.CreateIssuePayload.ReporterName,
		DueDate:      (*time.Time)(params.CreateIssuePayload.DueDate),
	}
	var response middleware.Responder
	if _, err := dbConnection.ValidateAndCreate(&newIssue); err != nil {
		zap.L().Error("DB Insertion", zap.Error(err))
		// how do I raise an erorr?
		response = issueop.NewCreateIssueBadRequest()
	} else {
		issuePayload := payloadForIssueModel(newIssue)
		response = issueop.NewCreateIssueCreated().WithPayload(&issuePayload)

	}
	return response
}

// IndexIssuesHandler returns a list of all issues
func IndexIssuesHandler(params issueop.IndexIssuesParams) middleware.Responder {
	var issues models.Issues
	var response middleware.Responder
	if err := dbConnection.All(&issues); err != nil {
		zap.L().Error("DB Query", zap.Error(err))
		response = issueop.NewIndexIssuesBadRequest()
	} else {
		issuePayloads := make(internalmodel.IndexIssuesPayload, len(issues))
		for i, issue := range issues {
			issuePayload := payloadForIssueModel(issue)
			issuePayloads[i] = &issuePayload
		}
		response = issueop.NewIndexIssuesOK().WithPayload(issuePayloads)
	}
	return response
}
