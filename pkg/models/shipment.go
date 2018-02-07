package models

import (
	"time"

	"github.com/satori/go.uuid"
)

// A Shipment represents a transfer of items form one location to another.
type Shipment struct {
	ID                              uuid.UUID     `json:"id" db:"id"`
	CreatedAt                       time.Time     `json:"created_at"`
	UpdatedAt                       time.Time     `json:"updated_at"`
	Name                            string        `json:"name"`
	PickupDate                      time.Time     `json:"pickup_date"`
	DeliveryDate                    time.Time     `json:"delivery_date"`
	TrafficDistributionListID       uuid.UUID     `json:"traffic_distribution_list_id"`
	TransportationServiceProviderID uuid.NullUUID `json:"transportation_service_provider_id"`
	AdministrativeShipment          bool
}

// Shipments is a slice of individual Shipment values.
type Shipments []Shipment
