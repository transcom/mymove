package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetOrdersHandler returns Orders by uuid
type GetOrdersHandler struct {
	handlers.HandlerContext
}

// Handle (GetOrdersHandler) responds to GET /orders/{uuid}
func (h GetOrdersHandler) Handle(params ordersoperations.GetOrdersParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewGetOrdersUnauthorized()
	}

	var err error

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		h.Logger().Info("Not a valid UUID")
		return ordersoperations.NewGetOrdersBadRequest()
	}

	orders, err := models.FetchElectronicOrderByID(h.DB(), id)
	if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by ID")
		return ordersoperations.NewGetOrdersInternalServerError()
	}

	apiOrders := ordersmessages.Orders{}
	apiOrders.UUID = params.UUID
	apiOrders.Edipi = &orders.Edipi
	apiOrders.OrdersNum = &orders.OrdersNumber
	apiOrders.Issuer = orders.Issuer
	// TODO check permission to retrieve orders by this issuer
	apiOrders.Revisions = make([]*ordersmessages.Revision, len(orders.Revisions))

	for i, o := range orders.Revisions {
		rev := ordersmessages.Revision{}
		seqNum := int64(o.SeqNum)
		rev.SeqNum = &seqNum
		member := ordersmessages.Member{
			Affiliation: o.Affiliation,
			FamilyName:  &o.FamilyName,
			GivenName:   &o.GivenName,
			Rank:        o.Paygrade,
		}
		if o.MiddleName != nil {
			member.MiddleName = *o.MiddleName
		}
		if o.NameSuffix != nil {
			member.Suffix = *o.NameSuffix
		}
		if o.Title != nil {
			member.Title = *o.Title
		}
		rev.Member = &member
		rev.Status = o.Status
		rev.DateIssued = strfmt.DateTime(o.DateIssued)
		rev.NoCostMove = o.NoCostMove
		rev.TdyEnRoute = o.TdyEnRoute
		rev.TourType = o.TourType
		rev.OrdersType = o.OrdersType
		rev.HasDependents = &o.HasDependents
		rev.LosingUnit = new(ordersmessages.Unit)
		if o.LosingUnitName != nil {
			rev.LosingUnit.Name = *o.LosingUnitName
		}
		if o.LosingUnitCity != nil {
			rev.LosingUnit.City = *o.LosingUnitCity
		}
		if o.LosingUnitLocality != nil {
			rev.LosingUnit.Locality = *o.LosingUnitLocality
		}
		if o.LosingUnitCountry != nil {
			rev.LosingUnit.Country = *o.LosingUnitCountry
		}
		if o.LosingUnitPostalCode != nil {
			rev.LosingUnit.PostalCode = *o.LosingUnitPostalCode
		}
		if o.LosingUIC != nil {
			rev.LosingUnit.Uic = *o.LosingUIC
		}
		rev.GainingUnit = new(ordersmessages.Unit)
		if o.GainingUnitName != nil {
			rev.GainingUnit.Name = *o.GainingUnitName
		}
		if o.GainingUnitCity != nil {
			rev.GainingUnit.City = *o.GainingUnitCity
		}
		if o.GainingUnitLocality != nil {
			rev.GainingUnit.Locality = *o.GainingUnitLocality
		}
		if o.GainingUnitCountry != nil {
			rev.GainingUnit.Country = *o.GainingUnitCountry
		}
		if o.GainingUnitPostalCode != nil {
			rev.GainingUnit.PostalCode = *o.GainingUnitPostalCode
		}
		if o.GainingUIC != nil {
			rev.GainingUnit.Uic = *o.GainingUIC
		}
		rev.ReportNoEarlierThan = (*strfmt.Date)(o.ReportNoEarlierThan)
		rev.ReportNoLaterThan = (*strfmt.Date)(o.ReportNoLaterThan)
		rev.PcsAccounting = new(ordersmessages.Accounting)
		if o.HhgTAC != nil {
			rev.PcsAccounting.Tac = *o.HhgTAC
		}
		if o.HhgSDN != nil {
			rev.PcsAccounting.Sdn = *o.HhgSDN
		}
		if o.HhgLOA != nil {
			rev.PcsAccounting.Loa = *o.HhgLOA
		}
		rev.NtsAccounting = new(ordersmessages.Accounting)
		if o.NtsTAC != nil {
			rev.NtsAccounting.Tac = *o.NtsTAC
		}
		if o.NtsSDN != nil {
			rev.NtsAccounting.Sdn = *o.NtsSDN
		}
		if o.NtsLOA != nil {
			rev.NtsAccounting.Loa = *o.NtsLOA
		}
		rev.PovShipmentAccounting = new(ordersmessages.Accounting)
		if o.PovShipmentTAC != nil {
			rev.PovShipmentAccounting.Tac = *o.PovShipmentTAC
		}
		if o.PovShipmentSDN != nil {
			rev.PovShipmentAccounting.Sdn = *o.PovShipmentSDN
		}
		if o.PovShipmentLOA != nil {
			rev.PovShipmentAccounting.Loa = *o.PovShipmentLOA
		}
		rev.PovStorageAccounting = new(ordersmessages.Accounting)
		if o.PovStorageTAC != nil {
			rev.PovStorageAccounting.Tac = *o.PovStorageTAC
		}
		if o.PovStorageSDN != nil {
			rev.PovStorageAccounting.Sdn = *o.PovStorageSDN
		}
		if o.PovStorageLOA != nil {
			rev.PovStorageAccounting.Loa = *o.PovStorageLOA
		}
		rev.UbAccounting = new(ordersmessages.Accounting)
		if o.UbTAC != nil {
			rev.UbAccounting.Tac = *o.UbTAC
		}
		if o.HhgSDN != nil {
			rev.UbAccounting.Sdn = *o.UbSDN
		}
		if o.UbLOA != nil {
			rev.UbAccounting.Loa = *o.UbLOA
		}
		if o.Comments != nil {
			rev.Comments = *o.Comments
		}

		apiOrders.Revisions[i] = &rev
	}

	return ordersoperations.NewGetOrdersOK().WithPayload(&apiOrders)
}
