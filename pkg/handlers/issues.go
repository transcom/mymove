package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	issueop "github.com/transcom/mymove/pkg/gen/genserver/operations/issues"
	"github.com/transcom/mymove/pkg/gen/messages"
	"github.com/transcom/mymove/pkg/models"
)

func responseForIssueModel(issue models.Issue) messages.Issue {
	ca := strfmt.DateTime(issue.CreatedAt)
	id := strfmt.UUID(issue.ID.String())
	ua := strfmt.DateTime(issue.UpdatedAt)
	issueResponse := messages.Issue{
		CreatedAt:   &ca,
		Description: &issue.Description,
		ID:          &id,
		UpdatedAt:   &ua,
	}
	return issueResponse
}

// Creates a new issue via POST /issue
func CreateIssueHandler(params issueop.CreateIssueParams) middleware.Responder {
	fmt.Println("NEW ISSUE TIME")

	payload := *params.CreateIssuePayload
	newIssue := models.Issue{
		Description: *payload.Description,
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

// Returns a list of all issues
func IndexIssuesHandler(params issueop.IndexIssuesParams) middleware.Responder {
	fmt.Println("INDEXISSUES TIME")

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
