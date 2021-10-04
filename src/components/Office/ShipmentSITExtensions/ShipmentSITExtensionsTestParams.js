import SERVICE_ITEM_STATUSES, { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { SIT_EXTENSION_REASON } from 'constants/sitExtensions';
// import { SHIPMENT_OPTIONS } from 'shared/constants';

const LOCATION_VALUES = {
  ORIGIN: 'ORIGIN',
  DESTINATION: 'DESTINATION',
};

export const SITExtensions = [
  {
    createdAt: '2021-09-13T15:41:59.373Z',
    decisionDate: '2021-09-13T15:41:59.373Z',
    eTag: 'MjAyMS0wOS0xM1QxNTo0MTo1OS4zNzM2NTRa',
    id: '7af5d51a-789c-4f5e-83dd-d905daed0785',
    mtoShipmentID: '8afd043a-8304-4e36-a695-7728e415990d',
    requestReason: SIT_EXTENSION_REASON.SERIOUS_ILLNESS_MEMBER,
    approvedDays: 30,
    status: 'APPROVED',
    updatedAt: '2021-09-13T15:41:59.373Z',
  },
];
export const SITExtensionsWithComments = [
  {
    createdAt: '2021-09-13T15:41:59.373Z',
    decisionDate: '0001-01-01T00:00:00.000Z',
    eTag: 'MjAyMS0wOS0xM1QxNTo0MTo1OS4zNzM2NTRa',
    id: '7af5d51a-789c-4f5e-83dd-d905daed0785',
    mtoShipmentID: '8afd043a-8304-4e36-a695-7728e415990d',
    requestReason: SIT_EXTENSION_REASON.SERIOUS_ILLNESS_MEMBER,
    approvedDays: 30,
    status: 'APPROVED',
    updatedAt: '2021-09-13T15:41:59.373Z',
    officeRemarks: 'The service member is unable to move into their new home at the expected time.',
    contractorRemarks: 'The customer requested an extension.',
  },
];

export const SITExtensionPending = [
  {
    id: '7af5d51a-789c-4f5e-83dd-d905daed0785',
    mtoShipmentID: '8afd043a-8304-4e36-a695-7728e415990d',
    requestReason: SIT_EXTENSION_REASON.SERIOUS_ILLNESS_MEMBER,
    status: 'PENDING',
    contractorRemarks: 'The customer requested an extension.',
  },
];

export const SITStatusOrigin = {
  location: LOCATION_VALUES.ORIGIN,
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  daysInSIT: 15,
  sitEntryDate: '2021-08-13T15:41:59.373Z',
  sitDepartureDate: '2021-08-28T15:41:59.373Z',
};

export const SITStatusDestination = {
  location: LOCATION_VALUES.DESTINATION,
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  daysInSIT: 15,
  sitEntryDate: '2021-08-13T15:41:59.373Z',
  sitDepartureDate: '2021-08-28T15:41:59.373Z',
};

export const SITStatusWithPastSITOriginServiceItem = {
  daysInSIT: 30,
  location: LOCATION_VALUES.DESTINATION,
  sitEntryDate: '2021-08-23T00:00:00.000Z',
  totalDaysRemaining: 210,
  totalSITDaysUsed: 60,
  pastSITServiceItems: [
    {
      SITPostalCode: '90210',
      createdAt: '2021-09-22T14:48:37.610Z',
      deletedAt: '0001-01-01',
      description: null,
      eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy42MTAxODZa',
      id: 'eb3a1983-4961-4be2-bfb6-73ad1720418a',
      moveTaskOrderID: '55c1cdbb-95b0-47f0-ab17-cfe5a0e46ab8',
      mtoShipmentID: 'f39ba92d-7d42-446a-be70-3a97b5f9f081',
      pickupPostalCode: null,
      reServiceCode: SERVICE_ITEM_CODES.DOPSIT,
      reServiceID: 'd1a4f062-0ca3-4387-8f8e-3dd20493d0b7',
      reServiceName: 'Domestic origin SIT pickup',
      reason: 'peak season all trucks in use',
      sitDepartureDate: '2021-08-23T00:00:00.000Z',
      sitEntryDate: '2021-07-24T00:00:00.000Z',
      status: SERVICE_ITEM_STATUSES.APPROVED,
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
};

export const SITStatusWithPastSITServiceItems = {
  daysInSIT: 30,
  location: LOCATION_VALUES.DESTINATION,
  sitEntryDate: '2021-08-23T00:00:00.000Z',
  totalDaysRemaining: 210,
  totalSITDaysUsed: 60,
  pastSITServiceItems: [
    {
      SITPostalCode: '90210',
      createdAt: '2021-09-22T14:48:37.610Z',
      deletedAt: '0001-01-01',
      description: null,
      eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy42MTAxODZa',
      id: 'eb3a1983-4961-4be2-bfb6-73ad1720418a',
      moveTaskOrderID: '55c1cdbb-95b0-47f0-ab17-cfe5a0e46ab8',
      mtoShipmentID: 'f39ba92d-7d42-446a-be70-3a97b5f9f081',
      pickupPostalCode: null,
      reServiceCode: SERVICE_ITEM_CODES.DOPSIT,
      reServiceID: 'd1a4f062-0ca3-4387-8f8e-3dd20493d0b7',
      reServiceName: 'Domestic origin SIT pickup',
      reason: 'peak season all trucks in use',
      sitDepartureDate: '2021-08-23T00:00:00.000Z',
      sitEntryDate: '2021-07-24T00:00:00.000Z',
      status: SERVICE_ITEM_STATUSES.APPROVED,
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      SITPostalCode: '08540',
      createdAt: '2021-09-22T14:48:37.610Z',
      deletedAt: '0001-01-01',
      description: null,
      eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy42MTAxODZb',
      id: 'eb3a1983-4961-4be2-bfb6-73ad1720418b',
      moveTaskOrderID: '55c1cdbb-95b0-47f0-ab17-cfe5a0e46ab8',
      mtoShipmentID: 'f39ba92d-7d42-446a-be70-3a97b5f9f081',
      pickupPostalCode: null,
      reServiceCode: SERVICE_ITEM_CODES.DDDSIT,
      reServiceID: 'd0561c49-e1a9-40b8-a739-3e639a9d77af',
      reServiceName: 'Domestic destination SIT pickup',
      reason: 'peak season all trucks in use',
      sitDepartureDate: '2021-09-24T00:00:00.000Z',
      sitEntryDate: '2021-09-03T00:00:00.000Z',
      status: SERVICE_ITEM_STATUSES.APPROVED,
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
};

export const SITShipment = {
  actualPickupDate: '2020-03-16',
  approvedDate: '2020-03-20T00:00:00.000Z',
  calculatedBillableWeight: 980,
  createdAt: '2021-09-22T14:48:37.546Z',
  customerRemarks: 'Please treat gently',
  destinationAddress: {
    city: 'Fairfield',
    country: 'US',
    eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy41MzYyNDRa',
    id: 'a755279e-462d-46e3-8701-d882f61b1b5c',
    postal_code: '94535',
    state: 'CA',
    street_address_1: '987 Any Avenue',
    street_address_2: 'P.O. Box 9876',
    street_address_3: 'c/o Some Person',
  },
  eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy41NDY2ODVa',
  id: 'f39ba92d-7d42-446a-be70-3a97b5f9f081',
  moveTaskOrderID: '55c1cdbb-95b0-47f0-ab17-cfe5a0e46ab8',
  pickupAddress: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy41MjQyMDla',
    id: '1a42d854-f0b8-428b-a0cc-a43c73a5df8f',
    postal_code: '90210',
    state: 'CA',
    street_address_1: '123 Any Street',
    street_address_2: 'P.O. Box 12345',
    street_address_3: 'c/o Some Person',
  },
  primeActualWeight: 980,
  requestedDeliveryDate: '2020-03-15',
  requestedPickupDate: '2020-03-15',
  scheduledPickupDate: '2020-03-16',
  shipmentType: 'HHG',
  sitDaysAllowance: 270,
  status: 'APPROVED',
  updatedAt: '2021-09-22T14:48:37.546Z',
};
