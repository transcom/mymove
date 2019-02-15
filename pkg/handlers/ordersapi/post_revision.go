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

	r := params.Revision
	var dateIssued time.Time
	if r.DateIssued == nil {
		dateIssued = time.Now()
	} else {
		dateIssued = time.Time(*r.DateIssued)
	}

	var tourType ordersmessages.TourType
	if r.TourType == "" {
		tourType = ordersmessages.TourTypeAccompanied
	} else {
		tourType = r.TourType
	}

	newRevision := models.ElectronicOrdersRevision{
		ElectronicOrderID:     orders.ID,
		ElectronicOrder:       orders,
		SeqNum:                int(*r.SeqNum),
		GivenName:             r.Member.GivenName,
		MiddleName:            r.Member.MiddleName,
		FamilyName:            r.Member.FamilyName,
		NameSuffix:            r.Member.Suffix,
		Affiliation:           r.Member.Affiliation,
		Paygrade:              r.Member.Rank,
		Title:                 r.Member.Title,
		Status:                r.Status,
		DateIssued:            dateIssued,
		NoCostMove:            r.NoCostMove,
		TdyEnRoute:            r.TdyEnRoute,
		TourType:              tourType,
		OrdersType:            r.OrdersType,
		HasDependents:         *r.HasDependents,
		LosingUIC:             r.LosingUnit.Uic,
		LosingUnitName:        r.LosingUnit.Name,
		LosingUnitCity:        r.LosingUnit.City,
		LosingUnitLocality:    r.LosingUnit.Locality,
		LosingUnitCountry:     r.LosingUnit.Country,
		LosingUnitPostalCode:  r.LosingUnit.PostalCode,
		GainingUIC:            r.GainingUnit.Uic,
		GainingUnitName:       r.GainingUnit.Name,
		GainingUnitCity:       r.GainingUnit.City,
		GainingUnitLocality:   r.GainingUnit.Locality,
		GainingUnitCountry:    r.GainingUnit.Country,
		GainingUnitPostalCode: r.GainingUnit.PostalCode,
		ReportNoEarlierThan:   (*time.Time)(r.ReportNoEarlierThan),
		ReportNoLaterThan:     (*time.Time)(r.ReportNoLaterThan),
		Comments:              r.Comments,
	}
	if r.PcsAccounting != nil {
		newRevision.HhgTAC = r.PcsAccounting.Tac
		newRevision.HhgSDN = r.PcsAccounting.Sdn
		newRevision.HhgLOA = r.PcsAccounting.Loa
	}
	if r.NtsAccounting != nil {
		newRevision.NtsTAC = r.NtsAccounting.Tac
		newRevision.NtsSDN = r.NtsAccounting.Sdn
		newRevision.NtsLOA = r.NtsAccounting.Loa
	}
	if r.PovShipmentAccounting != nil {
		newRevision.PovShipmentTAC = r.PovShipmentAccounting.Tac
		newRevision.PovShipmentSDN = r.PovShipmentAccounting.Sdn
		newRevision.PovShipmentLOA = r.PovShipmentAccounting.Loa
	}
	if r.PovStorageAccounting != nil {
		newRevision.PovStorageTAC = r.PovStorageAccounting.Tac
		newRevision.PovStorageSDN = r.PovStorageAccounting.Sdn
		newRevision.PovStorageLOA = r.PovStorageAccounting.Loa
	}
	if r.UbAccounting != nil {
		newRevision.UbTAC = r.UbAccounting.Tac
		newRevision.UbSDN = r.UbAccounting.Sdn
		newRevision.UbLOA = r.UbAccounting.Loa
	}

	verrs, err := models.CreateElectronicOrdersRevision(ctx, h.DB(), &newRevision)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	orders.Revisions = append(orders.Revisions, newRevision)

	orderPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ordersoperations.NewPostRevisionCreated().WithPayload(orderPayload)
}
