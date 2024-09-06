/* eslint-disable import/prefer-default-export */
import { arrayOf, bool, oneOf, number, shape, string } from 'prop-types';

import {
  ppmShipmentStatuses,
  shipmentDestinationTypes,
  shipmentStatuses,
  boatShipmentTypes,
} from 'constants/shipments';
import { LOA_TYPE, SHIPMENT_TYPES, SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape, ResidentialAddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { LOCATION_TYPES_ONE_OF, SitStatusShape } from 'types/sitStatusShape';
import { SITExtensionShape } from 'types/sitExtensions';
import { ExistingUploadsShape } from 'types/uploads';
import { expenseTypesArr } from 'constants/ppmExpenseTypes';

export const ShipmentOptionsOneOf = oneOf(Object.values(SHIPMENT_OPTIONS));

export const ShipmentTypesOneOf = oneOf(Object.values(SHIPMENT_TYPES));

export const ShipmentStatusesOneOf = oneOf(Object.values(shipmentStatuses));

export const PPMShipmentStatusOneOf = oneOf(Object.values(ppmShipmentStatuses));

export const BoatShipmentTypeOneOf = oneOf(Object.values(boatShipmentTypes));

export const GCCFactorsShape = shape({
  baseLinehaul: number,
  originLinehaulFactor: number,
  destinationLinehaulFactor: number,
  linehaulAdjustment: number,
  shorthaulCharge: number,
  transportationCost: number,
  linehaulFuelSurcharge: number,
  fuelSurchargePercent: number,
  originServiceAreaFee: number,
  originFactor: number,
  destinationServiceAreaFee: number,
  destinationFactor: number,
  packPrice: number,
  unpackPrice: number,
  ppmFactor: number,
});

export const incentivesShape = shape({
  grossIncentive: number,
  gcc: number,
  remi: number,
});

export const PPMShipmentShape = shape({
  id: string,
  shipmentId: string,
  shipmentLocator: string,
  createdAt: string,
  status: PPMShipmentStatusOneOf,
  expectedDepartureDate: string,
  actualMoveDate: string,
  submittedAt: string,
  reviewedAt: string,
  approvedAt: string,
  actualPickupPostalCode: string,
  actualDestinationPostalCode: string,
  sitExpected: bool,
  estimatedWeight: number,
  actualWeight: number,
  hasProGear: bool,
  proGearWeight: number,
  spouseProGearWeight: number,
  estimatedIncentive: number,
  finalEstimatedIncentive: number,
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
  incentives: incentivesShape,
  gcc: GCCFactorsShape,
  miles: number,
});

export const BoatShipmentShape = shape({
  id: string,
  shipmentId: string,
  shipmentLocator: string,
  createdAt: string,
  type: BoatShipmentTypeOneOf,
  year: number,
  make: string,
  model: string,
  lengthInInches: number,
  widthInInches: number,
  heightInInches: number,
  hasTrailer: bool,
  isRoadworthy: bool,
  eTag: string,
});

export const MobileHomeShipmentShape = shape({
  id: string,
  shipmentId: string,
  shipmentLocator: string,
  createdAt: string,
  year: number,
  make: string,
  model: string,
  lengthInInches: number,
  heightInInches: number,
  widthInInches: number,
  eTag: string,
});

export const PPMCloseoutShape = shape({
  id: string,
  plannedMoveDate: string,
  actualMoveDate: string,
  miles: number,
  estimatedWeight: number,
  actualWeight: number,
  proGearWeightCustomer: number,
  proGearWeightSpouse: number,
  grossIncentive: number,
  gcc: number,
  aoa: number,
  RemainingIncentive: number,
  haulPrice: number,
  haulFSC: number,
  dop: number,
  ddp: number,
  packPrice: number,
  unpackPrice: number,
  sitReimbursement: number,
});

// This type is badly defined because we have code that overloads the destinationType field on the shipment object as
// it is passed around with the display value instead of the value that we get from the API and instead of putting it on
// as a separate attribute.
export const ShipmentDestinationTypeOneOf = oneOf(
  Object.keys(shipmentDestinationTypes).concat(Object.values(shipmentDestinationTypes)),
);

export const ShipmentAddressUpdateShape = shape({
  contractorRemarks: string,
  id: string,
  newAddress: AddressShape,
  originalAddress: AddressShape,
  shipmentID: string,
  status: string,
  officeRemarks: string,
});

export const ShipmentShape = shape({
  moveTaskOrderID: string,
  id: string,
  shipmentLocator: string,
  createdAt: string,
  updatedAt: string,
  deletedAt: string,
  primeEstimatedWeight: number,
  primeActualWeight: number,
  calculatedBillableWeight: number,
  ntsRecordedWeight: number,
  scheduledPickupDate: string,
  scheduledDeliveryDate: string,
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
  tertiaryPickupAddress: AddressShape,
  tertiaryDeliveryAddress: AddressShape,
  customerRemarks: string,
  counselorRemarks: string,
  shipmentType: ShipmentTypesOneOf,
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
  boatShipment: BoatShipmentShape,
  mobileHomeShipment: MobileHomeShipmentShape,
  deliveryAddressUpdate: ShipmentAddressUpdateShape,
  actual_pro_gear_weight: number,
  actual_spouse_pro_gear_weight: number,
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
  reason: string,
  status: string,
  adjustedNetWeight: number,
  netWeightRemarks: string,
});

export const ExpenseShape = shape({
  id: string,
  ppmShipmentId: string,
  description: string,
  movingExpenseType: oneOf(expenseTypesArr),
  missingReceipt: bool,
  documentId: string,
  document: DocumentShape,
  amount: number,
  paidWithGtcc: bool,
  sitStartDate: string,
  sitEndDate: string,
  sitLocation: string,
  sitWeight: number,
});

export const StorageFacilityShape = shape({
  facilityName: string,
  phone: string,
  email: string,
  address: ResidentialAddressShape,
  lotNumber: string,
});

export const ProGearTicketShape = shape({
  belongsToSelf: bool,
  proGearWeight: number,
  description: string,
  missingWeightTicket: bool,
  reason: string,
  status: string,
});
