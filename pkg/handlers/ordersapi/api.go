package ordersapi

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/ordersapi"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// NewOrdersAPI the Orders API
func NewOrdersAPI(handlerConfig handlers.HandlerConfig) *ordersoperations.MymoveAPI {

	// Wire up the handlers to the ordersAPIMux
	ordersSpec, err := loads.Analyzed(ordersapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	ordersAPI := ordersoperations.NewMymoveAPI(ordersSpec)
	ordersAPI.ServeError = handlers.ServeCustomError
	ordersAPI.GetOrdersHandler = GetOrdersHandler{handlerConfig}
	ordersAPI.GetOrdersByIssuerAndOrdersNumHandler = GetOrdersByIssuerAndOrdersNumHandler{handlerConfig}
	ordersAPI.IndexOrdersForMemberHandler = IndexOrdersForMemberHandler{handlerConfig}
	ordersAPI.PostRevisionHandler = PostRevisionHandler{handlerConfig}
	ordersAPI.PostRevisionToOrdersHandler = PostRevisionToOrdersHandler{handlerConfig}
	return ordersAPI
}

func payloadForElectronicOrderModel(order *models.ElectronicOrder) (*ordersmessages.Orders, error) {
	var revisionPayloads []*ordersmessages.Revision
	for _, revision := range order.Revisions {
		payload, err := payloadForElectronicOrdersRevisionModel(revision)
		if err != nil {
			return nil, err
		}
		revisionPayloads = append(revisionPayloads, payload)
	}

	ordersPayload := &ordersmessages.Orders{
		UUID:      strfmt.UUID(order.ID.String()),
		OrdersNum: order.OrdersNumber,
		Edipi:     order.Edipi,
		Issuer:    ordersmessages.NewIssuer(ordersmessages.Issuer(order.Issuer)),
		Revisions: revisionPayloads,
	}
	return ordersPayload, nil
}

func payloadForElectronicOrdersRevisionModel(revision models.ElectronicOrdersRevision) (*ordersmessages.Revision, error) {
	seqNum := int64(revision.SeqNum)
	affiliation := ordersmessages.NewAffiliation(ordersmessages.Affiliation(revision.Affiliation))
	rank := ordersmessages.NewRank(ordersmessages.Rank(revision.Paygrade))
	status := ordersmessages.NewStatus(ordersmessages.Status(revision.Status))
	ordersType := ordersmessages.NewOrdersType(ordersmessages.OrdersType(revision.OrdersType))

	revisionPayload := &ordersmessages.Revision{
		SeqNum: &seqNum,
		Member: &ordersmessages.Member{
			GivenName:   revision.GivenName,
			MiddleName:  revision.MiddleName,
			FamilyName:  revision.FamilyName,
			Suffix:      revision.NameSuffix,
			Affiliation: affiliation,
			Rank:        rank,
			Title:       revision.Title,
		},
		Status:        status,
		DateIssued:    handlers.FmtDateTimePtr(&revision.DateIssued),
		NoCostMove:    revision.NoCostMove,
		TdyEnRoute:    revision.TdyEnRoute,
		TourType:      ordersmessages.TourType(revision.TourType),
		OrdersType:    ordersType,
		HasDependents: &revision.HasDependents,
		LosingUnit: &ordersmessages.Unit{
			Uic:        revision.LosingUIC,
			Name:       revision.LosingUnitName,
			City:       revision.LosingUnitCity,
			Locality:   revision.LosingUnitLocality,
			Country:    revision.LosingUnitCountry,
			PostalCode: revision.LosingUnitPostalCode,
		},
		GainingUnit: &ordersmessages.Unit{
			Uic:        revision.GainingUIC,
			Name:       revision.GainingUnitName,
			City:       revision.GainingUnitCity,
			Locality:   revision.GainingUnitLocality,
			Country:    revision.GainingUnitCountry,
			PostalCode: revision.GainingUnitPostalCode,
		},
		ReportNoEarlierThan: handlers.FmtDatePtr(revision.ReportNoEarlierThan),
		ReportNoLaterThan:   handlers.FmtDatePtr(revision.ReportNoLaterThan),
		PcsAccounting: &ordersmessages.Accounting{
			Tac: revision.HhgTAC,
			Sdn: revision.HhgSDN,
			Loa: revision.HhgLOA,
		},
		NtsAccounting: &ordersmessages.Accounting{
			Tac: revision.NtsTAC,
			Sdn: revision.NtsSDN,
			Loa: revision.NtsLOA,
		},
		PovShipmentAccounting: &ordersmessages.Accounting{
			Tac: revision.PovShipmentTAC,
			Sdn: revision.PovShipmentSDN,
			Loa: revision.PovShipmentLOA,
		},
		PovStorageAccounting: &ordersmessages.Accounting{
			Tac: revision.PovStorageTAC,
			Sdn: revision.PovStorageSDN,
			Loa: revision.PovStorageLOA,
		},
		UbAccounting: &ordersmessages.Accounting{
			Tac: revision.UbTAC,
			Sdn: revision.UbSDN,
			Loa: revision.UbLOA,
		},
		Comments: revision.Comments,
	}
	return revisionPayload, nil
}

func verifyOrdersReadAccess(issuer models.Issuer, cert *models.ClientCert) bool {
	switch issuer {
	case models.IssuerAirForce:
		return cert.AllowAirForceOrdersRead
	case models.IssuerArmy:
		return cert.AllowArmyOrdersRead
	case models.IssuerCoastGuard:
		return cert.AllowCoastGuardOrdersRead
	case models.IssuerMarineCorps:
		return cert.AllowMarineCorpsOrdersRead
	case models.IssuerNavy:
		return cert.AllowNavyOrdersRead
	default:
		// Unknown issuer
		return false
	}
}

func verifyOrdersWriteAccess(issuer models.Issuer, cert *models.ClientCert) bool {
	switch issuer {
	case models.IssuerAirForce:
		return cert.AllowAirForceOrdersWrite
	case models.IssuerArmy:
		return cert.AllowArmyOrdersWrite
	case models.IssuerCoastGuard:
		return cert.AllowCoastGuardOrdersWrite
	case models.IssuerMarineCorps:
		return cert.AllowMarineCorpsOrdersWrite
	case models.IssuerNavy:
		return cert.AllowNavyOrdersWrite
	default:
		// Unknown issuer
		return false
	}
}
