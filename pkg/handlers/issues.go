package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
	issueop "github.com/transcom/mymove/pkg/gen/restapi/operations/issues"
	"github.com/transcom/mymove/pkg/models"
)

func responseForIssueModel(issue models.Issue) messages.IssueResponse {
	createdAt := strfmt.DateTime(issue.CreatedAt)
	id := strfmt.UUID(issue.ID.String())
	updatedAt := strfmt.DateTime(issue.UpdatedAt)
	issueResponse := messages.IssueResponse{
		CreatedAt:   &createdAt,
		Description: &issue.Description,
		ID:          &id,
		UpdatedAt:   &updatedAt,
	}
	return issueResponse
}

// CreateIssueHandler creates a new issue via POST /issue
func CreateIssueHandler(params issueop.CreateIssueParams) middleware.Responder {
	newIssue := models.Issue{
		Description: *params.CreateIssuePayload.Description,
	}
	var response middleware.Responder
	if err := dbConnection.Create(&newIssue); err != nil {
		zap.L().Error("DB Insertion", zap.Error(err))
		// how do I raise an erorr?
		response = issueop.NewCreateIssueBadRequest()
	} else {
		issueResponse := responseForIssueModel(newIssue)
		response = issueop.NewCreateIssueCreated().WithPayload(&issueResponse)

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
		issueResponses := make(messages.IndexIssuesResponse, len(issues))
		for i, issue := range issues {
			issueResponse := responseForIssueModel(issue)
			issueResponses[i] = &issueResponse
		}
		response = issueop.NewIndexIssuesOK().WithPayload(issueResponses)
	}
	return response
}
