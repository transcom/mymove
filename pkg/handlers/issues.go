package handlers

import (
	"fmt"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/genserver/operations/issues"
	"github.com/transcom/mymove/pkg/messages"
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

	payload := params.CreateIssuePayload
	fmt.Println(*payload.Description)

	desc := "THis is not hte issue you are looking for"
	var id strfmt.UUID = "c56a4180-65aa-42ec-a945-5fd21dec0538"
	dt := strfmt.NewDateTime()

	issueResponse := messages.Issue{
		CreatedAt:   &dt,
		Description: &desc,
		ID:          &id,
		UpdatedAt:   &dt,
	}

	response := issues.NewCreateIssueCreated().WithPayload(&issueResponse)
	return response

}
