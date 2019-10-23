package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetCustomerInfoHandler fetches the information of a specific customer
type GetCustomerInfoHandler struct {
	handlers.HandlerContext
}

// Handle getting the information of a specific customer
func (h GetCustomerInfoHandler) Handle(params customercodeop.GetCustomerInfoParams) middleware.Responder {
	// for now just return static data
	customer := &ghcmessages.Customer{
		FirstName:              models.StringPointer("First"),
		MiddleName:             models.StringPointer("Middle"),
		LastName:               models.StringPointer("Last"),
		Agency:                 models.StringPointer("Agency"),
		Grade:                  models.StringPointer("Grade"),
		Email:                  models.StringPointer("Example@example.com"),
		Telephone:              models.StringPointer("213-213-3232"),
		OriginDutyStation:      models.StringPointer("Origin Station"),
		DestinationDutyStation: models.StringPointer("Destination Station"),
		DependentsAuthorized:   true,
	}
	return customercodeop.NewGetCustomerInfoOK().WithPayload(customer)
}
