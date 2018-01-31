package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
	issueop "github.com/transcom/mymove/pkg/gen/restapi/operations/issues"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForIssueModel(issue models.Issue) messages.IssuePayload {
	issuePayload := messages.IssuePayload{
		CreatedAt:    pointerFromSDateTime(strfmt.DateTime(issue.CreatedAt)),
		Description:  pointerFromString(issue.Description),
		ID:           pointerFromSUUID(strfmt.UUID(issue.ID.String())),
		UpdatedAt:    pointerFromSDateTime(strfmt.DateTime(issue.UpdatedAt)),
		ReporterName: issue.ReporterName,
	}
	return issuePayload
}

// CreateIssueHandler creates a new issue via POST /issue
func CreateIssueHandler(params issueop.CreateIssueParams) middleware.Responder {
	newIssue := models.Issue{
		Description:  *params.CreateIssuePayload.Description,
		ReporterName: params.CreateIssuePayload.ReporterName,
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
		issuePayloads := make(messages.IndexIssuesPayload, len(issues))
		for i, issue := range issues {
			issuePayload := payloadForIssueModel(issue)
			issuePayloads[i] = &issuePayload
		}
		response = issueop.NewIndexIssuesOK().WithPayload(issuePayloads)
	}
	return response
}
