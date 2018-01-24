package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/genserver/operations/issues"
	"github.com/transcom/mymove/pkg/gen/messages"
	"github.com/transcom/mymove/pkg/models"
)

// func CreateMoveHandler(params issues.CreateMoveParams) middleware.Responder {
// 		fmt.Println("HELLO HANDLER")
// 		fmt.Println(params.Move)
// 		moveParams := params.Move

// 		fmt.Println("AND NOW, MAGIC")
// 		newMove := &messages.Move{}
// 		newMove.URL = "http://localhost:12324/moves/112-ebf"
// 		newMove.Name = moveParams.Name
// 		newMove.Email = moveParams.Email

// 		myMoves = append(myMoves, *newMove)

// 		fmt.Println(newMove.Email)
// 		success := issues.NewCreateMoveCreated().WithPayload(newMove)
// 		return success
// 	}

// Creates a new issue via POST /issue
func CreateIssueHandler(params issues.CreateIssueParams) middleware.Responder {
	fmt.Println("NEW ISSUE TIME")

	payload := *params.CreateIssuePayload
	newIssue := models.Issue{
		Description: *payload.Description,
	}
	var response middleware.Responder
	if err := dbConnection.Create(&newIssue); err != nil {
		zap.L().Error("DB Insertion", zap.Error(err))
		// how do I raise an erorr?
		response = issues.NewCreateIssueBadRequest()
	} else {
		ca := strfmt.DateTime(newIssue.CreatedAt)
		id := strfmt.UUID(newIssue.ID.String())
		ua := strfmt.DateTime(newIssue.UpdatedAt)
		issueResponse := messages.Issue{
			CreatedAt:   &ca,
			Description: &newIssue.Description,
			ID:          &id,
			UpdatedAt:   &ua,
		}
		response = issues.NewCreateIssueCreated().WithPayload(&issueResponse)

	}
	return response
}
