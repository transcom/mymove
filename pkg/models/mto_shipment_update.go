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

// Helper function for converting an Update model to a Data model
// Created to aid in server testing
func (m *MTOShipmentUpdate) GetShipmentModelFromUpdateModel() (*MTOShipment, error) {
	mtoShipment := MTOShipment{
		ID:                               m.ID,
		MoveTaskOrder:                    m.MoveTaskOrder,
		MoveTaskOrderID:                  m.MoveTaskOrderID,
		ScheduledPickupDate:              m.ScheduledPickupDate,
		RequestedPickupDate:              m.RequestedPickupDate,
		RequestedDeliveryDate:            m.RequestedDeliveryDate,
		ApprovedDate:                     m.ApprovedDate,
		FirstAvailableDeliveryDate:       m.FirstAvailableDeliveryDate,
		ActualPickupDate:                 m.ActualPickupDate,
		RequiredDeliveryDate:             m.RequiredDeliveryDate,
		ScheduledDeliveryDate:            m.ScheduledDeliveryDate,
		ActualDeliveryDate:               m.ActualDeliveryDate,
		CustomerRemarks:                  m.CustomerRemarks,
		CounselorRemarks:                 m.CounselorRemarks.Value,
		PickupAddress:                    m.PickupAddress,
		PickupAddressID:                  m.PickupAddressID,
		DestinationAddress:               m.DestinationAddress,
		DestinationAddressID:             m.DestinationAddressID,
		DestinationType:                  m.DestinationType,
		MTOAgents:                        m.MTOAgents,
		MTOServiceItems:                  m.MTOServiceItems,
		SecondaryPickupAddress:           m.SecondaryPickupAddress,
		SecondaryPickupAddressID:         m.SecondaryPickupAddressID,
		SecondaryDeliveryAddress:         m.SecondaryDeliveryAddress,
		SecondaryDeliveryAddressID:       m.SecondaryDeliveryAddressID,
		SITDaysAllowance:                 m.SITDaysAllowance,
		SITExtensions:                    m.SITExtensions,
		PrimeEstimatedWeight:             m.PrimeEstimatedWeight,
		PrimeEstimatedWeightRecordedDate: m.PrimeEstimatedWeightRecordedDate,
		PrimeActualWeight:                m.PrimeActualWeight,
		BillableWeightCap:                m.BillableWeightCap,
		BillableWeightJustification:      m.BillableWeightJustification,
		NTSRecordedWeight:                m.NTSRecordedWeight,
		ShipmentType:                     m.ShipmentType,
		Status:                           m.Status,
		Diversion:                        m.Diversion,
		RejectionReason:                  m.RejectionReason,
		Distance:                         m.Distance,
		Reweigh:                          m.Reweigh,
		UsesExternalVendor:               m.UsesExternalVendor,
		StorageFacility:                  m.StorageFacility,
		StorageFacilityID:                m.StorageFacilityID,
		ServiceOrderNumber:               m.ServiceOrderNumber,
		TACType:                          m.TACType,
		SACType:                          m.SACType,
		PPMShipment:                      m.PPMShipment,
	}

	return &mtoShipment, nil
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
