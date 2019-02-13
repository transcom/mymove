package ordersapi

import (
	"sort"

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

	sharedID, err := uuid.FromString(params.UUID.String())
	if err != nil {
		h.Logger().Info("Not a valid UUID")
		return ordersoperations.NewGetOrdersBadRequest()
	}

	orders, err := models.FetchElectronicOrdersBySharedID(h.DB(), sharedID)
	if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by shared ID")
		return ordersoperations.NewGetOrdersInternalServerError()
	}

	if len(orders) == 0 {
		return ordersoperations.NewGetOrdersNotFound()
	}

	apiOrders := ordersmessages.Orders{}
	apiOrders.UUID = params.UUID
	apiOrders.Revisions = make([]*ordersmessages.Revision, len(orders))

	// sort orders by sequence number ascending
	sort.Slice(orders, func(i, j int) bool { return orders[i].SeqNum < orders[j].SeqNum })

	// use highest sequence number (i.e., the latest) for ordersNum, service member edipi, issuer
	// Although these SHOULD be the same for each revision, it's possible for them to differ
	latestOrders := orders[len(orders)-1]
	apiOrders.OrdersNum = latestOrders.OrdersNumber
	apiOrders.Edipi = latestOrders.ServiceMember.Edipi
	if latestOrders.OrdersIssuingAgency == nil {
		return ordersoperations.NewGetOrdersInternalServerError()
	}
	apiOrders.Issuer = ordersmessages.Issuer(*latestOrders.OrdersIssuingAgency)
	if err != nil {
		h.Logger().Info(err.Error())
		return ordersoperations.NewGetOrdersInternalServerError()
	}
	// TODO check permission to retrieve orders by this issuer

	for i, o := range orders {
		rev := ordersmessages.Revision{}
		seqNum := int64(o.SeqNum)
		rev.SeqNum = &seqNum
		member := ordersmessages.Member{
			Affiliation: o.EOrdersAffiliation,
			FamilyName:  o.EOrdersFamilyName,
			GivenName:   o.EOrdersGivenName,
			Rank:        o.EOrdersPaygrade,
		}
		if o.EOrdersMiddleName != nil {
			member.MiddleName = *o.EOrdersMiddleName
		}
		if o.EOrdersNameSuffix != nil {
			member.Suffix = *o.EOrdersNameSuffix
		}
		if o.EOrdersTitle != nil {
			member.Title = *o.EOrdersTitle
		}
		rev.Member = &member
		rev.Status = o.Impact
		rev.DateIssued = strfmt.DateTime(o.IssueDate)
		rev.NoCostMove = o.NoCostMove
		rev.TdyEnRoute = o.TdyEnRoute
		rev.TourType, err = tourTypeToAPITourType(o.TourType)
		if err != nil {
			h.Logger().Info(err.Error())
			return ordersoperations.NewGetOrdersInternalServerError()
		}
		rev.OrdersType, err = toAPIOrdersType(o.OrdersType)
		if err != nil {
			h.Logger().Info(err.Error())
			return ordersoperations.NewGetOrdersInternalServerError()
		}
		hasDependents := o.HasDependents
		rev.HasDependents = &hasDependents
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
		if o.LosingUnitPostCode != nil {
			rev.LosingUnit.PostalCode = *o.LosingUnitPostCode
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
		if o.GainingUnitPostCode != nil {
			rev.GainingUnit.PostalCode = *o.GainingUnitPostCode
		}
		if o.GainingUIC != nil {
			rev.GainingUnit.Uic = *o.GainingUIC
		}
		rnetDate := strfmt.Date(o.ReportNoEarlierThan)
		rev.ReportNoEarlierThan = &rnetDate
		rnltDate := strfmt.Date(o.ReportByDate)
		rev.ReportNoLaterThan = &rnltDate
		rev.PcsAccounting = new(ordersmessages.Accounting)
		rev.PcsAccounting.Tac = *o.TAC
		rev.PcsAccounting.Sdn = *o.HhgSDN
		rev.PcsAccounting.Loa = *o.HhgLOA
		rev.NtsAccounting = new(ordersmessages.Accounting)
		rev.NtsAccounting.Tac = *o.NtsTAC
		rev.NtsAccounting.Sdn = *o.NtsSDN
		rev.NtsAccounting.Loa = *o.NtsLOA
		rev.PovShipmentAccounting = new(ordersmessages.Accounting)
		rev.PovShipmentAccounting.Tac = *o.PovShipmentTAC
		rev.PovShipmentAccounting.Sdn = *o.PovShipmentSDN
		rev.PovShipmentAccounting.Loa = *o.PovShipmentLOA
		rev.PovStorageAccounting = new(ordersmessages.Accounting)
		rev.PovStorageAccounting.Tac = *o.PovStorageTAC
		rev.PovStorageAccounting.Sdn = *o.PovStorageSDN
		rev.PovStorageAccounting.Loa = *o.PovStorageLOA
		rev.UbAccounting = new(ordersmessages.Accounting)
		rev.UbAccounting.Tac = *o.UbTAC
		rev.UbAccounting.Sdn = *o.UbSDN
		rev.UbAccounting.Loa = *o.UbLOA
		rev.Comments = *o.Comments

		apiOrders.Revisions[i] = &rev
	}

	return ordersoperations.NewGetOrdersOK().WithPayload(&apiOrders)
}
