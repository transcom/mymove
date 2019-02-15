package ordersapi

import (
	"reflect"
	"time"

	"github.com/go-openapi/runtime/middleware"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// PostRevisionHandler adds a Revision to Orders matching the provided search parameters
type PostRevisionHandler struct {
	handlers.HandlerContext
}

// Handle (params ordersoperations.PostRevisionParams) responds to POST /orders
func (h PostRevisionHandler) Handle(params ordersoperations.PostRevisionParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewPostRevisionUnauthorized()
	}

	// TODO check enumerated values, or are these already checked for me by swagger?

	var edipi string
	if len(params.MemberID) == 9 {
		// TODO EDIPI lookup from DMDC
		return middleware.NotImplemented("EDIPI lookup from SSN not supported yet")
	} else if len(params.MemberID) != 10 {
		return ordersoperations.NewPostRevisionBadRequest()
	}
	edipi = params.MemberID

	// Is there already a Revision matching these Orders? (same ordersNum, edipi, issuer)
	orders, err := models.FetchElectronicOrderByUniqueFeatures(h.DB(), params.OrdersNum, params.MemberID, params.Issuer)

	if err == models.ErrFetchNotFound {
		orders = models.ElectronicOrder{
			OrdersNumber: params.OrdersNum,
			Edipi:        edipi,
			Issuer:       ordersmessages.Issuer(params.Issuer),
			Revisions:    []models.ElectronicOrdersRevision{},
		}
		verrs, err := models.CreateElectronicOrder(ctx, h.DB(), &orders)
		if err != nil || verrs.HasAny() {
			return handlers.ResponseForVErrors(h.Logger(), verrs, err)
		}
	} else if err != nil {
		return ordersoperations.NewPostRevisionInternalServerError()
	}

	for _, r := range orders.Revisions {
		// SeqNum collision
		if r.SeqNum == int(*params.Revision.SeqNum) {
			return ordersoperations.NewPostRevisionConflict()
		}
	}

	newRevision := toElectronicOrdersRevision(orders, params.Revision)
	verrs, err := models.CreateElectronicOrdersRevision(ctx, h.DB(), newRevision)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	orders.Revisions = append(orders.Revisions, *newRevision)

	orderPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ordersoperations.NewPostRevisionCreated().WithPayload(orderPayload)
}

func toElectronicOrdersRevision(orders models.ElectronicOrder, rev *ordersmessages.Revision) *models.ElectronicOrdersRevision {
	var dateIssued time.Time
	if rev.DateIssued == nil {
		dateIssued = time.Now()
	} else {
		dateIssued = time.Time(*rev.DateIssued)
	}

	var tourType ordersmessages.TourType
	if rev.TourType == "" {
		tourType = ordersmessages.TourTypeAccompanied
	} else {
		tourType = rev.TourType
	}

	newRevision := models.ElectronicOrdersRevision{
		ElectronicOrderID:     orders.ID,
		ElectronicOrder:       orders,
		SeqNum:                int(*rev.SeqNum),
		GivenName:             rev.Member.GivenName,
		MiddleName:            rev.Member.MiddleName,
		FamilyName:            rev.Member.FamilyName,
		NameSuffix:            rev.Member.Suffix,
		Affiliation:           rev.Member.Affiliation,
		Paygrade:              rev.Member.Rank,
		Title:                 rev.Member.Title,
		Status:                rev.Status,
		DateIssued:            dateIssued,
		NoCostMove:            rev.NoCostMove,
		TdyEnRoute:            rev.TdyEnRoute,
		TourType:              tourType,
		OrdersType:            rev.OrdersType,
		HasDependents:         *rev.HasDependents,
		LosingUIC:             rev.LosingUnit.Uic,
		LosingUnitName:        rev.LosingUnit.Name,
		LosingUnitCity:        rev.LosingUnit.City,
		LosingUnitLocality:    rev.LosingUnit.Locality,
		LosingUnitCountry:     rev.LosingUnit.Country,
		LosingUnitPostalCode:  rev.LosingUnit.PostalCode,
		GainingUIC:            rev.GainingUnit.Uic,
		GainingUnitName:       rev.GainingUnit.Name,
		GainingUnitCity:       rev.GainingUnit.City,
		GainingUnitLocality:   rev.GainingUnit.Locality,
		GainingUnitCountry:    rev.GainingUnit.Country,
		GainingUnitPostalCode: rev.GainingUnit.PostalCode,
		ReportNoEarlierThan:   (*time.Time)(rev.ReportNoEarlierThan),
		ReportNoLaterThan:     (*time.Time)(rev.ReportNoLaterThan),
		Comments:              rev.Comments,
	}
	if rev.PcsAccounting != nil {
		newRevision.HhgTAC = rev.PcsAccounting.Tac
		newRevision.HhgSDN = rev.PcsAccounting.Sdn
		newRevision.HhgLOA = rev.PcsAccounting.Loa
	}
	if rev.NtsAccounting != nil {
		newRevision.NtsTAC = rev.NtsAccounting.Tac
		newRevision.NtsSDN = rev.NtsAccounting.Sdn
		newRevision.NtsLOA = rev.NtsAccounting.Loa
	}
	if rev.PovShipmentAccounting != nil {
		newRevision.PovShipmentTAC = rev.PovShipmentAccounting.Tac
		newRevision.PovShipmentSDN = rev.PovShipmentAccounting.Sdn
		newRevision.PovShipmentLOA = rev.PovShipmentAccounting.Loa
	}
	if rev.PovStorageAccounting != nil {
		newRevision.PovStorageTAC = rev.PovStorageAccounting.Tac
		newRevision.PovStorageSDN = rev.PovStorageAccounting.Sdn
		newRevision.PovStorageLOA = rev.PovStorageAccounting.Loa
	}
	if rev.UbAccounting != nil {
		newRevision.UbTAC = rev.UbAccounting.Tac
		newRevision.UbSDN = rev.UbAccounting.Sdn
		newRevision.UbLOA = rev.UbAccounting.Loa
	}

	return &newRevision
}
