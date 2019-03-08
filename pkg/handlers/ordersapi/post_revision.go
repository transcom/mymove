package ordersapi

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-openapi/runtime/middleware"
	beeline "github.com/honeycombio/beeline-go"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/iws"
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
	if clientCert == nil {
		h.Logger().Info("No client certificate provided")
		return ordersoperations.NewPostRevisionUnauthorized()
	}
	if !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not permitted to access this API")
		return ordersoperations.NewPostRevisionForbidden()
	}
	if params.Issuer == string(ordersmessages.IssuerAirForce) {
		if !clientCert.AllowAirForceOrdersWrite {
			h.Logger().Info("Client certificate is not permitted to write Air Force Orders")
			return ordersoperations.NewPostRevisionForbidden()
		}
	} else if params.Issuer == string(ordersmessages.IssuerArmy) {
		if !clientCert.AllowArmyOrdersWrite {
			h.Logger().Info("Client certificate is not permitted to write Army Orders")
			return ordersoperations.NewPostRevisionForbidden()
		}
	} else if params.Issuer == string(ordersmessages.IssuerCoastGuard) {
		if !clientCert.AllowCoastGuardOrdersWrite {
			h.Logger().Info("Client certificate is not permitted to write Coast Guard Orders")
			return ordersoperations.NewPostRevisionForbidden()
		}
	} else if params.Issuer == string(ordersmessages.IssuerMarineCorps) {
		if !clientCert.AllowMarineCorpsOrdersWrite {
			h.Logger().Info("Client certificate is not permitted to write Marine Corps Orders")
			return ordersoperations.NewPostRevisionForbidden()
		}
	} else if params.Issuer == string(ordersmessages.IssuerNavy) {
		if !clientCert.AllowNavyOrdersWrite {
			h.Logger().Info("Client certificate is not permitted to write Navy Orders")
			return ordersoperations.NewPostRevisionForbidden()
		}
	} else {
		// Unknown issuer
		h.Logger().Info(fmt.Sprint("Unknown issuer ", params.Issuer))
		return ordersoperations.NewPostRevisionBadRequest()
	}

	var edipi string
	if len(params.MemberID) == 9 {
		rbsPersonLookup := h.IWSPersonLookup()

		iwsParams := iws.GetPersonUsingSSNParams{
			Ssn:       params.MemberID,
			LastName:  params.Revision.Member.FamilyName,
			FirstName: params.Revision.Member.GivenName,
		}
		matchReasonCode, edipiNum, _, _, err := rbsPersonLookup.GetPersonUsingSSN(iwsParams)
		if err != nil {
			h.Logger().Warn(fmt.Sprint("Error while retrieving EDIPI from Identity Web Services: ", err.Error()))
			return ordersoperations.NewPostRevisionInternalServerError()
		}
		switch matchReasonCode {
		case iws.MatchReasonCodeLimited:
			// limited match: the returned EDIPI matches the provided SSN and maybe first name but DMDC doesn't think the last name matches
			edipi = string(edipiNum)
		case iws.MatchReasonCodeFull:
			// full match means the returned EDIPI matches the provided SSN and last name
			edipi = string(edipiNum)
		case iws.MatchReasonCodeMultiple:
			// more than one EDIPI for this SSN! Uhh... how to choose? FWIW it's unlikely but not impossible to encounter this in the wild
			return ordersoperations.NewPostRevisionNotFound()
		case iws.MatchReasonCodeNone:
			// No match: fail
			return ordersoperations.NewPostRevisionNotFound()
		}
	} else if len(params.MemberID) != 10 {
		return ordersoperations.NewPostRevisionBadRequest()
	} else {
		edipi = params.MemberID
	}

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
		h.Logger().Warn(fmt.Sprintf("Error fetching electronic orders with OrdersNum %s, EDIPI %s, and Issuer %s: %s", params.OrdersNum, params.MemberID, params.Issuer, err.Error()))
		return ordersoperations.NewPostRevisionInternalServerError()
	}

	for _, r := range orders.Revisions {
		// SeqNum collision
		if r.SeqNum == int(params.Revision.SeqNum) {
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
		ElectronicOrderID:   orders.ID,
		ElectronicOrder:     orders,
		SeqNum:              int(rev.SeqNum),
		GivenName:           rev.Member.GivenName,
		MiddleName:          rev.Member.MiddleName,
		FamilyName:          rev.Member.FamilyName,
		NameSuffix:          rev.Member.Suffix,
		Affiliation:         rev.Member.Affiliation,
		Paygrade:            rev.Member.Rank,
		Title:               rev.Member.Title,
		Status:              rev.Status,
		DateIssued:          dateIssued,
		NoCostMove:          rev.NoCostMove,
		TdyEnRoute:          rev.TdyEnRoute,
		TourType:            tourType,
		OrdersType:          rev.OrdersType,
		HasDependents:       rev.HasDependents,
		ReportNoEarlierThan: (*time.Time)(rev.ReportNoEarlierThan),
		ReportNoLaterThan:   (*time.Time)(rev.ReportNoLaterThan),
		Comments:            rev.Comments,
	}
	if rev.LosingUnit != nil {
		newRevision.LosingUIC = rev.LosingUnit.Uic
		newRevision.LosingUnitName = rev.LosingUnit.Name
		newRevision.LosingUnitCity = rev.LosingUnit.City
		newRevision.LosingUnitLocality = rev.LosingUnit.Locality
		newRevision.LosingUnitCountry = rev.LosingUnit.Country
		newRevision.LosingUnitPostalCode = rev.LosingUnit.PostalCode
	}
	if rev.GainingUnit != nil {
		newRevision.GainingUIC = rev.GainingUnit.Uic
		newRevision.GainingUnitName = rev.GainingUnit.Name
		newRevision.GainingUnitCity = rev.GainingUnit.City
		newRevision.GainingUnitLocality = rev.GainingUnit.Locality
		newRevision.GainingUnitCountry = rev.GainingUnit.Country
		newRevision.GainingUnitPostalCode = rev.GainingUnit.PostalCode
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
