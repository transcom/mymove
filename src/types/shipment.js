/* eslint-disable import/prefer-default-export */
import { arrayOf, bool, oneOf, number, shape, string } from 'prop-types';

import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape, ResidentialAddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { LOCATION_TYPES_ONE_OF } from 'types/sitStatusShape';

export const ShipmentOptionsOneOf = oneOf([
  SHIPMENT_OPTIONS.HHG,
  SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.NTS,
  SHIPMENT_OPTIONS.NTSR,
  SHIPMENT_OPTIONS.PPM,
]);

export const ShipmentStatusesOneOf = oneOf([
  shipmentStatuses.DRAFT,
  shipmentStatuses.SUBMITTED,
  shipmentStatuses.APPROVED,
  shipmentStatuses.CANCELLATION_REQUESTED,
  shipmentStatuses.DIVERSION_REQUESTED,
  shipmentStatuses.CANCELED,
  shipmentStatuses.REJECTED,
]);

export const PPMShipmentStatus = oneOf([
  ppmShipmentStatuses.DRAFT,
  ppmShipmentStatuses.SUBMITTED,
  ppmShipmentStatuses.WAITING_ON_CUSTOMER,
  ppmShipmentStatuses.NEEDS_ADVANCE_APPROVAL,
  ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL,
  ppmShipmentStatuses.PAYMENT_APPROVED,
]);

export const PPMShipmentShape = shape({
  id: string,
  shipmentId: string,
  createdAt: string,
  status: PPMShipmentStatus,
  expectedDepartureDate: string,
  actualMoveDate: string,
  submittedAt: string,
  reviewedAt: string,
  approvedAt: string,
  pickupPostalCode: string,
  secondaryPickupPostalCode: string,
  actualPickupPostalCode: string,
  destinationPostalCode: string,
  secondaryDestinationPostalCode: string,
  actualDestinationPostalCode: string,
  sitExpected: bool,
  estimatedWeight: number,
  netWeight: number,
  hasProGear: bool,
  proGearWeight: number,
  spouseProGearWeight: number,
  estimatedIncentive: number,
  hasRequestedAdvance: bool,
  advanceAmountRequested: number,
  hasReceivedAdvance: bool,
  advanceAmountReceived: number,
  sitLocation: LOCATION_TYPES_ONE_OF,
  sitEstimatedWeight: number,
  sitEstimatedEntryDate: string,
  sitEstimatedDepartureDate: string,
  sitEstimatedCost: number,
  eTag: string,
});

export const ShipmentShape = shape({
  id: string,
  shipmentType: ShipmentOptionsOneOf,
  requestedPickupDate: string,
  scheduledPickupDate: string,
  actualPickupDate: string,
  requestedDeliveryDate: string,
  pickupAddress: AddressShape,
  secondaryPickupAddress: AddressShape,
  destinationAddress: AddressShape,
  secondaryDeliveryAddress: AddressShape,
  agents: arrayOf(AgentShape),
  primeEstimatedWeight: number,
  primeActualWeight: number,
  ntsRecordedWeight: number,
  diversion: bool,
  counselorRemarks: string,
  customerRemarks: string,
  status: ShipmentStatusesOneOf,
  reweigh: shape({
    id: string,
    weight: number,
  }),
  storageFacility: shape({
    address: AddressShape.isRequired,
    facilityName: string.isRequired,
    lotNumber: string,
    phone: string,
    email: string,
  }),
  ppmShipment: PPMShipmentShape,
  eTag: string,
});

export const StorageFacilityShape = shape({
  facilityName: string,
  phone: string,
  email: string,
  address: ResidentialAddressShape,
  lotNumber: string,
});
