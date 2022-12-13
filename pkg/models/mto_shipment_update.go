package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

type MTOShipmentUpdate struct {
	ID                               uuid.UUID
	MoveTaskOrder                    Move
	MoveTaskOrderID                  uuid.UUID
	ScheduledPickupDate              *time.Time
	RequestedPickupDate              *time.Time
	RequestedDeliveryDate            *time.Time
	ApprovedDate                     *time.Time
	FirstAvailableDeliveryDate       *time.Time
	ActualPickupDate                 *time.Time
	RequiredDeliveryDate             *time.Time
	ScheduledDeliveryDate            *time.Time
	ActualDeliveryDate               *time.Time
	CustomerRemarks                  *string
	CounselorRemarks                 Nullable[*string]
	PickupAddress                    *Address
	PickupAddressID                  *uuid.UUID
	DestinationAddress               *Address
	DestinationAddressID             *uuid.UUID
	DestinationType                  *DestinationType
	MTOAgents                        MTOAgents
	MTOServiceItems                  MTOServiceItems
	SecondaryPickupAddress           *Address
	SecondaryPickupAddressID         *uuid.UUID
	SecondaryDeliveryAddress         *Address
	SecondaryDeliveryAddressID       *uuid.UUID
	SITDaysAllowance                 *int
	SITExtensions                    SITExtensions
	PrimeEstimatedWeight             *unit.Pound
	PrimeEstimatedWeightRecordedDate *time.Time
	PrimeActualWeight                *unit.Pound
	BillableWeightCap                *unit.Pound
	BillableWeightJustification      *string
	NTSRecordedWeight                *unit.Pound
	ShipmentType                     MTOShipmentType
	Status                           MTOShipmentStatus
	Diversion                        bool
	RejectionReason                  *string
	Distance                         *unit.Miles
	Reweigh                          *Reweigh
	UsesExternalVendor               bool
	StorageFacility                  *StorageFacility
	StorageFacilityID                *uuid.UUID
	ServiceOrderNumber               *string
	TACType                          *LOAType
	SACType                          *LOAType
	PPMShipment                      *PPMShipment
	CreatedAt                        time.Time
	UpdatedAt                        time.Time
	DeletedAt                        *time.Time
}
type Nullable[T any] struct {
	Present bool
	Value   T
}

// func NewNullable(T any) Nullable[T] {
// 	return Nullable[T]{Present: true, Value: T}
// }

// type NullableI[T any] interface {
// 	Present() bool
// 	Value() T
// }
