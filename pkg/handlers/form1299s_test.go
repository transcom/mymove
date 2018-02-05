package handlers

import (
	"fmt"
	"testing"

	"github.com/go-openapi/strfmt"
	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
)

func TestShowFormHandler(t *testing.T) {

	goodID := strfmt.UUID("240ec3c5-019d-4031-9c27-8a553e022297")
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: goodID}

	response := ShowForm1299Handler(showFormParams)

	// assert we got back the 201 response
	createdResponse := response.(*form1299op.ShowForm1299OK)
	fmt.Println(createdResponse)
	// createdIssuePayload := createdResponse.Payload

	// if *createdIssuePayload.Description != testDescription {
	// 	t.Error("Didn't get the same description back")
	// }

	// if createdIssuePayload.ReporterName != nil {
	// 	t.Error("We should not have sent anything back for the reporter name")
	// }

}
