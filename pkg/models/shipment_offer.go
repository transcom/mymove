package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// ShipmentOffer maps a Transportation Service Provider to a shipment,
// indicating that the shipment has been offered to that TSP.
type ShipmentOffer struct {
	ID                                         uuid.UUID                                `json:"id" db:"id"`
	CreatedAt                                  time.Time                                `json:"created_at" db:"created_at"`
	UpdatedAt                                  time.Time                                `json:"updated_at" db:"updated_at"`
	ShipmentID                                 uuid.UUID                                `json:"shipment_id" db:"shipment_id"`
	Shipment                                   Shipment                                 `belongs_to:"shipments"`
	TransportationServiceProviderID            uuid.UUID                                `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	TransportationServiceProvider              TransportationServiceProvider            `belongs_to:"transportation_service_providers"`
	TransportationServiceProviderPerformanceID uuid.UUID                                `json:"transportation_service_provider_performance_id" db:"transportation_service_provider_performance_id"`
	TransportationServiceProviderPerformance   TransportationServiceProviderPerformance `belongs_to:"transportation_service_provider_performances"`
	AdministrativeShipment                     bool                                     `json:"administrative_shipment" db:"administrative_shipment"`
	Accepted                                   *bool                                    `json:"accepted" db:"accepted"`
	RejectionReason                            *string                                  `json:"rejection_reason" db:"rejection_reason"`
}

// String is not required by pop and may be deleted
func (so ShipmentOffer) String() string {
	ja, _ := json.Marshal(so)
	return string(ja)
}

// ShipmentOffers is not required by pop and may be deleted
type ShipmentOffers []ShipmentOffer

// String is not required by pop and may be deleted
func (so ShipmentOffers) String() string {
	ja, _ := json.Marshal(so)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (so *ShipmentOffer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: so.ShipmentID, Name: "ShipmentID"},
		&validators.UUIDIsPresent{Field: so.TransportationServiceProviderID, Name: "TransportationServiceProviderID"},
		&validators.UUIDIsPresent{Field: so.TransportationServiceProviderPerformanceID, Name: "TransportationServiceProviderPerformanceID"},
	), nil
}

// State Machinery
// Avoid calling ShipmentOffer.Accepted = ... or ShipmentOffer.RejectionReason = ... ever. Use these methods to change the state.

// Accept marks the Shipment Offer request as Accepted.
func (so *ShipmentOffer) Accept() error {
	if so.Accepted != nil {
		return errors.Wrap(ErrInvalidTransition, "Accept")
	}
	accepted := true
	so.Accepted = &accepted
	return nil
}

// Reject marks the Shipment Offer request as Rejected and sets the Rejection Reason.
func (so *ShipmentOffer) Reject(rejectionReason string) error {
	if so.Accepted != nil {
		return errors.Wrap(ErrInvalidTransition, "Reject")
	}
	notAccepted := false
	so.Accepted = &notAccepted
	so.RejectionReason = &rejectionReason
	return nil
}

// CreateShipmentOffer connects a shipment to a transportation service provider. This
// function assumes that the match has been validated by the caller.
func CreateShipmentOffer(tx *pop.Connection,
	shipmentID uuid.UUID,
	tspID uuid.UUID,
	tsppID uuid.UUID,
	administrativeShipment bool) (*ShipmentOffer, error) {

	shipmentOffer := ShipmentOffer{
		ShipmentID:                                 shipmentID,
		TransportationServiceProviderID:            tspID,
		TransportationServiceProviderPerformanceID: tsppID,
		AdministrativeShipment:                     administrativeShipment,
	}
	_, err := tx.ValidateAndSave(&shipmentOffer)

	return &shipmentOffer, err
}

// FetchShipmentOfferByTSP Fetches a shipment belonging to a TSP ID by Shipment ID
func FetchShipmentOfferByTSP(tx *pop.Connection, tspID uuid.UUID, shipmentID uuid.UUID) (*ShipmentOffer, error) {

	shipmentOffers := []ShipmentOffer{}

	err := tx.
		Where("shipment_offers.transportation_service_provider_id = $1 and shipment_offers.shipment_id = $2", tspID, shipmentID).
		All(&shipmentOffers)

	if err != nil {
		return nil, err
	}

	// Unlikely that we see more than one but to be safe this will error.
	if len(shipmentOffers) != 1 {
		return nil, ErrFetchNotFound
	}

	return &shipmentOffers[0], err
}

// Accepted returns the accepted shipment offers from a slice of shipment offers.
func (so *ShipmentOffers) Accepted() (ShipmentOffers, error) {
	var acceptedOffers ShipmentOffers
	for _, offer := range *so {
		if offer.Accepted != nil && *offer.Accepted {
			acceptedOffers = append(acceptedOffers, offer)

			if offer.TransportationServiceProviderPerformance.TransportationServiceProvider.ID == uuid.Nil {
				return nil, errors.Errorf("Accepted shipment offer is missing Transportation Service Provider")
			}

			if offer.TransportationServiceProviderPerformance.TransportationServiceProvider.StandardCarrierAlphaCode == "" {
				return nil, errors.Errorf("Accepted shipment offer TSP is missing SCAC for TSP ID: %s\n",
					offer.TransportationServiceProviderPerformance.TransportationServiceProvider.ID.String())
			}

			if offer.TransportationServiceProviderPerformance.TransportationServiceProvider.SupplierID == nil ||
				*(offer.TransportationServiceProviderPerformance.TransportationServiceProvider.SupplierID) == "" {
				return nil, errors.Errorf("Accepted shipment offer TSP is missing SupplierID for TSP ID: %s\n",
					offer.TransportationServiceProviderPerformance.TransportationServiceProvider.ID.String())
			}
		}
	}

	return acceptedOffers, nil
}

// SCAC return TransportationServiceProvider's SCAC from the ShipmentOffer
func (so *ShipmentOffer) SCAC() (string, error) {
	scac := so.TransportationServiceProviderPerformance.TransportationServiceProvider.StandardCarrierAlphaCode
	if scac == "" {
		return "", errors.New("SCAC for TSP is missing with TSP ID: " + so.TransportationServiceProviderPerformance.TransportationServiceProviderID.String())
	}
	return scac, nil
}

// SupplierID return TransportationServiceProvider's SupplierID from the ShipmentOffer
func (so *ShipmentOffer) SupplierID() (*string, error) {
	supplierID := so.TransportationServiceProviderPerformance.TransportationServiceProvider.SupplierID
	if supplierID == nil {
		return nil, errors.New("SupplierID for TSP is missing with TSP ID " + so.TransportationServiceProviderPerformance.TransportationServiceProviderID.String())
	}
	return supplierID, nil
}
