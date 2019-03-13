package ordersapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/ordersapi"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// NewOrdersAPIHandler returns a handler for the Orders API
func NewOrdersAPIHandler(context handlers.HandlerContext) http.Handler {

	// Wire up the handlers to the ordersAPIMux
	ordersSpec, err := loads.Analyzed(ordersapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	ordersAPI := ordersoperations.NewMymoveAPI(ordersSpec)
	ordersAPI.GetOrdersHandler = GetOrdersHandler{context}
	ordersAPI.GetOrdersByIssuerAndOrdersNumHandler = GetOrdersByIssuerAndOrdersNumHandler{context}
	ordersAPI.IndexOrdersForMemberHandler = IndexOrdersForMemberHandler{context}
	ordersAPI.PostRevisionHandler = PostRevisionHandler{context}
	ordersAPI.PostRevisionToOrdersHandler = PostRevisionToOrdersHandler{context}
	return ordersAPI.Serve(nil)
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
		Issuer:    ordersmessages.Issuer(order.Issuer),
		Revisions: revisionPayloads,
	}
	return ordersPayload, nil
}

func payloadForElectronicOrdersRevisionModel(revision models.ElectronicOrdersRevision) (*ordersmessages.Revision, error) {
	seqNum := int64(revision.SeqNum)
	revisionPayload := &ordersmessages.Revision{
		SeqNum: &seqNum,
		Member: &ordersmessages.Member{
			GivenName:   revision.GivenName,
			MiddleName:  revision.MiddleName,
			FamilyName:  revision.FamilyName,
			Suffix:      revision.NameSuffix,
			Affiliation: ordersmessages.Affiliation(revision.Affiliation),
			Rank:        ordersmessages.Rank(revision.Paygrade),
			Title:       revision.Title,
		},
		Status:        ordersmessages.Status(revision.Status),
		DateIssued:    (*strfmt.DateTime)(&revision.DateIssued),
		NoCostMove:    revision.NoCostMove,
		TdyEnRoute:    revision.TdyEnRoute,
		TourType:      ordersmessages.TourType(revision.TourType),
		OrdersType:    ordersmessages.OrdersType(revision.OrdersType),
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
		ReportNoEarlierThan: (*strfmt.Date)(revision.ReportNoEarlierThan),
		ReportNoLaterThan:   (*strfmt.Date)(revision.ReportNoLaterThan),
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

func verifyOrdersReadAccess(issuer models.Issuer, cert *models.ClientCert, logger *zap.Logger, logFailure bool) bool {
	switch issuer {
	case models.IssuerAirForce:
		if !cert.AllowAirForceOrdersRead {
			if logFailure {
				logger.Info("Client certificate is not permitted to read Air Force Orders")
			}
			return false
		}
	case models.IssuerArmy:
		if !cert.AllowArmyOrdersRead {
			if logFailure {
				logger.Info("Client certificate is not permitted to read Army Orders")
			}
			return false
		}
	case models.IssuerCoastGuard:
		if !cert.AllowCoastGuardOrdersRead {
			if logFailure {
				logger.Info("Client certificate is not permitted to read Coast Guard Orders")
			}
			return false
		}
	case models.IssuerMarineCorps:
		if !cert.AllowMarineCorpsOrdersRead {
			if logFailure {
				logger.Info("Client certificate is not permitted to read Marine Corps Orders")
			}
			return false
		}
	case models.IssuerNavy:
		if !cert.AllowNavyOrdersRead {
			if logFailure {
				logger.Info("Client certificate is not permitted to read Navy Orders")
			}
			return false
		}
	default:
		// Unknown issuer
		logger.Error(fmt.Sprint("Unknown issuer ", issuer))
		return false
	}
	return true
}

func verifyOrdersWriteAccess(issuer models.Issuer, cert *models.ClientCert, logger *zap.Logger) bool {
	switch issuer {
	case models.IssuerAirForce:
		if !cert.AllowAirForceOrdersWrite {
			logger.Info("Client certificate is not permitted to write Air Force Orders")
			return false
		}
	case models.IssuerArmy:
		if !cert.AllowArmyOrdersWrite {
			logger.Info("Client certificate is not permitted to write Army Orders")
			return false
		}
	case models.IssuerCoastGuard:
		if !cert.AllowCoastGuardOrdersWrite {
			logger.Info("Client certificate is not permitted to write Coast Guard Orders")
			return false
		}
	case models.IssuerMarineCorps:
		if !cert.AllowMarineCorpsOrdersWrite {
			logger.Info("Client certificate is not permitted to write Marine Corps Orders")
			return false
		}
	case models.IssuerNavy:
		if !cert.AllowNavyOrdersWrite {
			logger.Info("Client certificate is not permitted to write Navy Orders")
			return false
		}
	default:
		// Unknown issuer
		logger.Error(fmt.Sprint("Unknown issuer ", issuer))
		return false
	}
	return true
}
