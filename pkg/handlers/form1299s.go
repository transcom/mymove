package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForForm1299Model(form1299 models.Form1299) messages.Form1299Payload {

	form1299Payload := messages.Form1299Payload{
		CreatedAt:              fmtDateTime(form1299.CreatedAt),
		ID:                     fmtUUID(form1299.ID),
		UpdatedAt:              fmtDateTime(form1299.UpdatedAt),
		DatePrepared:           (*strfmt.Date)(form1299.DatePrepared),
		ServiceMemberFirstName: form1299.ServiceMemberFirstName,
		HhgTotalPounds:         form1299.HhgTotalPounds,
	}
	return form1299Payload
}

type myInt int

// CreateForm1299Handler creates a new form1299 via POST /form1299
func CreateForm1299Handler(params form1299op.CreateForm1299Params) middleware.Responder {
	newForm1299 := models.Form1299{
		DatePrepared:           (*time.Time)(params.CreateForm1299Payload.DatePrepared),
		ServiceMemberFirstName: params.CreateForm1299Payload.ServiceMemberFirstName,
		HhgTotalPounds:         params.CreateForm1299Payload.HhgTotalPounds,
	}
	var response middleware.Responder
	if err := dbConnection.Create(&newForm1299); err != nil {
		zap.L().Error("DB Insertion", zap.Error(err))
		// how do I raise an error?
		response = form1299op.NewCreateForm1299BadRequest()
	} else {
		form1299Payload := payloadForForm1299Model(newForm1299)
		response = form1299op.NewCreateForm1299Created().WithPayload(&form1299Payload)

	}
	return response
}

// IndexForm1299sHandler returns a list of all form1299s
func IndexForm1299sHandler(params form1299op.IndexForm1299sParams) middleware.Responder {
	var form1299s models.Form1299s
	var response middleware.Responder
	if err := dbConnection.All(&form1299s); err != nil {
		zap.L().Error("DB Query", zap.Error(err))
		response = form1299op.NewIndexForm1299sBadRequest()
	} else {
		form1299Payloads := make(messages.IndexForm1299sPayload, len(form1299s))
		for i, form1299 := range form1299s {
			form1299Payload := payloadForForm1299Model(form1299)
			form1299Payloads[i] = &form1299Payload
		}
		response = form1299op.NewIndexForm1299sOK().WithPayload(form1299Payloads)
	}
	return response
}
