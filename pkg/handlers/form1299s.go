package handlers

import (
	"fmt"
	// "time"

	"github.com/go-openapi/runtime/middleware"
	// "github.com/go-openapi/strfmt"
	// "go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
	"github.com/transcom/mymove/pkg/models"
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
	formID := params.Form1299ID
	form := models.Form1299{}
	var response middleware.Responder
	if err := dbConnection.Find(&form, formID); err != nil {
		fmt.Println(err)
		response = form1299op.NewShowForm1299NotFound()
	} else {
		formPayload := messages.Form1299Payload{}
		formPayload.ID = fmtUUID(form.ID)
		formPayload.NameOfPreparingOffice = form.NameOfPreparingOffice
		response = form1299op.NewShowForm1299OK().WithPayload(&formPayload)
	}

	return response

}
