package handlers

import (
	"fmt"
	// "time"

	"github.com/go-openapi/runtime/middleware"
	// "github.com/go-openapi/strfmt"
	// "go.uber.org/zap"

	// "github.com/transcom/mymove/pkg/gen/messages"
	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
	// "github.com/transcom/mymove/pkg/models"
)

// func payloadForIssueModel(issue models.Issue) messages.IssuePayload {
// 	issuePayload := messages.IssuePayload{
// 		CreatedAt:    fmtDateTime(issue.CreatedAt),
// 		Description:  stringPointer(issue.Description),
// 		ID:           fmtUUID(issue.ID),
// 		UpdatedAt:    fmtDateTime(issue.UpdatedAt),
// 		ReporterName: issue.ReporterName,
// 		DueDate:      (*strfmt.Date)(issue.DueDate),
// 	}
// 	return issuePayload
// }

// CreateIssueHandler creates a new issue via POST /issue
func ShowForm1299Handler(params form1299op.ShowForm1299Params) middleware.Responder {
	fmt.Println("WEOINWEFOWFNWOEFNN")
	fmt.Println(params.Form1299ID)
	// newIssue := models.Issue{
	// 	Description:  *params.CreateIssuePayload.Description,
	// 	ReporterName: params.CreateIssuePayload.ReporterName,
	// 	DueDate:      (*time.Time)(params.CreateIssuePayload.DueDate),
	// }
	// var response middleware.Responder
	// if _, err := dbConnection.ValidateAndCreate(&newIssue); err != nil {
	// 	zap.L().Error("DB Insertion", zap.Error(err))
	// 	// how do I raise an erorr?
	// 	response = issueop.NewCreateIssueBadRequest()
	// } else {
	// 	issuePayload := payloadForIssueModel(newIssue)
	// 	response = issueop.NewCreateIssueCreated().WithPayload(&issuePayload)

	// }
	return nil
}
