/* eslint-disable import/prefer-default-export */
import { arrayOf, bool, oneOf, number, shape, string } from 'prop-types';

import { ppmShipmentStatuses, shipmentDestinationTypes, shipmentStatuses } from 'constants/shipments';
import { LOA_TYPE, SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape, ResidentialAddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { LOCATION_TYPES_ONE_OF, SitStatusShape } from 'types/sitStatusShape';
import { SITExtensionShape } from 'types/sitExtensions';
import { ExistingUploadsShape } from 'types/uploads';

export const ShipmentOptionsOneOf = oneOf(Object.values(SHIPMENT_OPTIONS));

export const ShipmentStatusesOneOf = oneOf(Object.values(shipmentStatuses));

export const PPMShipmentStatusOneOf = oneOf(Object.values(ppmShipmentStatuses));

export const PPMShipmentShape = shape({
  id: string,
  shipmentId: string,
  createdAt: string,
  status: PPMShipmentStatusOneOf,
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

// This type is badly defined because we have code that overloads the destinationType field on the shipment object as
// it is passed around with the display value instead of the value that we get from the API and instead of putting it on
// as a separate attribute.
export const ShipmentDestinationTypeOneOf = oneOf(
  Object.keys(shipmentDestinationTypes).concat(Object.values(shipmentDestinationTypes)),
);

export const ShipmentShape = shape({
  moveTaskOrderID: string,
  id: string,
  createdAt: string,
  updatedAt: string,
  deletedAt: string,
  primeEstimatedWeight: number,
  primeActualWeight: number,
  calculatedBillableWeight: number,
  ntsRecordedWeight: number,
  scheduledPickupDate: string,
  requestedPickupDate: string,
  actualPickupDate: string,
  requestedDeliveryDate: string,
  requiredDeliveryDate: string,
  approvedDate: string,
  diversion: bool,
  pickupAddress: AddressShape,
  destinationAddress: AddressShape,
  destinationType: ShipmentDestinationTypeOneOf,
  secondaryPickupAddress: AddressShape,
  secondaryDeliveryAddress: AddressShape,
  customerRemarks: string,
  counselorRemarks: string,
  shipmentType: ShipmentOptionsOneOf,
  status: ShipmentStatusesOneOf,
  rejectionReason: string,
  reweigh: shape({
    id: string,
    weight: number,
  }),
  agents: arrayOf(AgentShape), // We have different API definitions for a shipment and they name this field different things...
  mtoAgents: arrayOf(AgentShape), // We have different API definitions for a shipment and they name this field different things...
  sitDaysAllowance: number,
  sitExtensions: arrayOf(SITExtensionShape),
  sitStatus: SitStatusShape,
  eTag: string,
  billableWeightCap: number,
  billableWeightJustification: string,
  tacType: oneOf(Object.values(LOA_TYPE)),
  sacType: oneOf(Object.values(LOA_TYPE)),
  usesExternalVendor: bool,
  serviceOrderNumber: string,
  storageFacility: shape({
    address: AddressShape.isRequired,
    facilityName: string.isRequired,
    lotNumber: string,
    phone: string,
    email: string,
  }),
  ppmShipment: PPMShipmentShape,
});

const DocumentShape = shape({
  id: string,
  serviceMemberId: string,
  uploads: ExistingUploadsShape,
});

export const WeightTicketShape = shape({
  id: string,
  ppmShipmentId: string,
  vehicleDescription: string,
  missingEmptyWeightTicket: bool,
  emptyWeight: number,
  emptyWeightDocumentId: string,
  emptyDocument: DocumentShape,
  fullWeight: number,
  missingFullWeightTicket: bool,
  fullWeightDocumentId: string,
  fullDocument: DocumentShape,
  ownsTrailer: bool,
  trailerMeetsCriteria: bool,
  trailerOwnershipDocumentId: string,
  proofOfTrailerOwnershipDocument: DocumentShape,
});

export const ExpenseShape = shape({
  id: string,
  ppmShipmentId: string,
  description: string,
  missingReceipt: bool,
  receiptDocumentId: string,
  receiptDocument: DocumentShape,
  amount: number,
  paidWithGTCC: bool,
  sitStartDate: string,
  sitEndDate: string,
});

export const StorageFacilityShape = shape({
  facilityName: string,
  phone: string,
  email: string,
  address: ResidentialAddressShape,
  lotNumber: string,
});
