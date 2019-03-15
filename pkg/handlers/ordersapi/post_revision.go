package ordersapi

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"

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
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrWriteForbidden, "Not permitted to access this API"))
	}
	if !verifyOrdersWriteAccess(models.Issuer(params.Issuer), clientCert) {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrWriteForbidden, "Not permitted to write Orders from this issuer"))
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
			return handlers.ResponseForError(h.Logger(), err)
		}
		switch matchReasonCode {
		case iws.MatchReasonCodeLimited:
			// limited match: the returned EDIPI matches the provided SSN and maybe first name but DMDC doesn't think the last name matches
			edipi = fmt.Sprintf("%010d", edipiNum)
		case iws.MatchReasonCodeFull:
			// full match means the returned EDIPI matches the provided SSN and last name
			edipi = fmt.Sprintf("%010d", edipiNum)
		case iws.MatchReasonCodeMultiple:
			// more than one EDIPI for this SSN! Uhh... how to choose? FWIW it's unlikely but not impossible to encounter this in the wild
			return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrFetchNotFound, "DMDC IWS matched multiple EDIPIs for this SSN"))
		case iws.MatchReasonCodeNone:
			// No match: fail
			return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrFetchNotFound, "DMDC IWS matched no EDIPI for this SSN"))
		}
	} else if len(params.MemberID) == 10 {
		edipi = params.MemberID
	} else {
		// go-swagger's validation should prevent this from happening
		return ordersoperations.NewPostRevisionBadRequest()
	}

	// Is there already a Revision matching these Orders? (same ordersNum and issuer)
	orders, err := models.FetchElectronicOrderByIssuerAndOrdersNum(h.DB(), params.Issuer, params.OrdersNum)

	var newRevision *models.ElectronicOrdersRevision
	var verrs *validate.Errors
	if err == models.ErrFetchNotFound {
		// New Orders
		orders = &models.ElectronicOrder{
			OrdersNumber: params.OrdersNum,
			Edipi:        edipi,
			Issuer:       models.Issuer(params.Issuer),
			Revisions:    []models.ElectronicOrdersRevision{},
		}
		newRevision = toElectronicOrdersRevision(orders, params.Revision)
		verrs, err = models.CreateElectronicOrderWithRevision(ctx, h.DB(), orders, newRevision)
	} else if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	} else if orders.Edipi != edipi {
		return handlers.ResponseForError(
			h.Logger(),
			errors.WithMessage(
				models.ErrWriteConflict,
				fmt.Sprintf("Cannot POST Revision for EDIPI %s to Electronic Orders with OrdersNum %s from Issuer %s: the existing orders are issued to EDIPI %s", edipi, params.OrdersNum, params.Issuer, orders.Edipi)))
	} else {
		// Amending Orders
		for _, r := range orders.Revisions {
			// SeqNum collision
			if r.SeqNum == int(*params.Revision.SeqNum) {
				return handlers.ResponseForError(
					h.Logger(),
					errors.WithMessage(
						models.ErrWriteConflict,
						fmt.Sprintf("Cannot POST Revision with SeqNum %d for EDIPI %s to Electronic Orders with OrdersNum %s from Issuer %s: a Revision with that SeqNum already exists in those Orders", r.SeqNum, edipi, params.OrdersNum, params.Issuer)))
			}
		}

		newRevision = toElectronicOrdersRevision(orders, params.Revision)
		verrs, err = models.CreateElectronicOrdersRevision(ctx, h.DB(), newRevision)
	}
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

// toElectronicOrdersRevision converts an API Revision to a DB model
// ElectronicOrdersRevision, and sets the association with the provided DB
// model ElectronicOrder
func toElectronicOrdersRevision(orders *models.ElectronicOrder, rev *ordersmessages.Revision) *models.ElectronicOrdersRevision {
	var dateIssued time.Time
	if rev.DateIssued == nil {
		dateIssued = time.Now()
	} else {
		dateIssued = time.Time(*rev.DateIssued)
	}

	var tourType models.TourType
	if rev.TourType == "" {
		tourType = models.TourTypeAccompanied
	} else {
		tourType = models.TourType(rev.TourType)
	}

	newRevision := models.ElectronicOrdersRevision{
		ElectronicOrderID:   orders.ID,
		ElectronicOrder:     *orders,
		SeqNum:              int(*rev.SeqNum),
		GivenName:           rev.Member.GivenName,
		MiddleName:          rev.Member.MiddleName,
		FamilyName:          rev.Member.FamilyName,
		NameSuffix:          rev.Member.Suffix,
		Affiliation:         models.ElectronicOrdersAffiliation(rev.Member.Affiliation),
		Paygrade:            models.Paygrade(rev.Member.Rank),
		Title:               rev.Member.Title,
		Status:              models.ElectronicOrdersStatus(rev.Status),
		DateIssued:          dateIssued,
		NoCostMove:          rev.NoCostMove,
		TdyEnRoute:          rev.TdyEnRoute,
		TourType:            models.TourType(tourType),
		OrdersType:          models.ElectronicOrdersType(rev.OrdersType),
		HasDependents:       *rev.HasDependents,
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
