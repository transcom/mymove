import moment from 'moment';

import SERVICE_ITEM_STATUSES, { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { SIT_EXTENSION_REASON, SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { swaggerDateFormat } from 'shared/dates';

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

export const SITExtensionDenied = [
  {
    id: '7af5d51a-789c-4f5e-83dd-d905daed0785',
    decisionDate: '2021-09-13T15:41:59.373Z',
    mtoShipmentID: '8afd043a-8304-4e36-a695-7728e415990d',
    requestReason: SIT_EXTENSION_REASON.SERIOUS_ILLNESS_MEMBER,
    status: 'DENIED',
    contractorRemarks: 'The customer requested an extension.',
  },
];

export const SITStatusOrigin = {
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  calculatedTotalDaysInSIT: 45,
  currentSIT: {
    location: LOCATION_VALUES.ORIGIN,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
    sitCustomerContacted: '2021-08-26',
    sitRequestedDelivery: '2021-08-30',
  },
};

export const SITStatusOriginAuthorized = {
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  calculatedTotalDaysInSIT: 45,
  currentSIT: {
    location: LOCATION_VALUES.ORIGIN,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
    sitCustomerContacted: '2021-08-26',
    sitRequestedDelivery: '2021-08-30',
  },
};

export const SITStatusShowConvert = {
  totalSITDaysUsed: 45,
  totalDaysRemaining: 30,
  currentSIT: {
    location: LOCATION_VALUES.ORIGIN,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
  },
};

export const SITStatusDontShowConvert = {
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  currentSIT: {
    location: LOCATION_VALUES.ORIGIN,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
  },
};

export const SITStatusDestination = {
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  calculatedTotalDaysInSIT: 45,
  currentSIT: {
    location: LOCATION_VALUES.DESTINATION,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
    sitCustomerContacted: '2021-08-26',
    sitRequestedDelivery: '2021-08-30',
  },
};

export const SITStatusDestinationWithoutCustomerDeliveryInfo = {
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  calculatedTotalDaysInSIT: 45,
  currentSIT: {
    location: LOCATION_VALUES.DESTINATION,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
  },
};
export const SITStatusOriginWithoutCustomerDeliveryInfo = {
  totalSITDaysUsed: 45,
  totalDaysRemaining: 60,
  calculatedTotalDaysInSIT: 45,
  currentSIT: {
    location: LOCATION_VALUES.ORIGIN,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
  },
};

export const futureSITStatus = {
  totalDaysRemaining: 365,
  totalSITDaysUsed: 0,
  calculatedTotalDaysInSIT: 0,
  currentSIT: {
    location: LOCATION_VALUES.ORIGIN,
    daysInSIT: 0,
    sitEntryDate: moment().add(2, 'years').format(swaggerDateFormat),
    sitAuthorizedEndDate: moment().add(3, 'years').format(swaggerDateFormat),
  },
};

export const SITStatusWithPastSITOriginServiceItem = {
  daysInSIT: 30,
  location: LOCATION_VALUES.DESTINATION,
  sitEntryDate: '2021-08-23',
  totalDaysRemaining: 210,
  totalSITDaysUsed: 60,
  calculatedTotalDaysInSIT: 60,
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
      reServiceCode: SERVICE_ITEM_CODES.DOFSIT,
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
  sitEntryDate: '2021-08-23',
  totalDaysRemaining: 210,
  totalSITDaysUsed: 60,
  calculatedTotalDaysInSIT: 60,
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
      reServiceCode: SERVICE_ITEM_CODES.DOFSIT,
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
      sitDepartureDate: '2021-09-24',
      sitEntryDate: '2021-09-03',
      status: SERVICE_ITEM_STATUSES.APPROVED,
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
};

export const SITStatusWithPastSITServiceItemsDeparted = {
  totalDaysRemaining: 210,
  totalSITDaysUsed: 60,
  calculatedTotalDaysInSIT: 60,
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
      reServiceCode: SERVICE_ITEM_CODES.DOFSIT,
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
      sitDepartureDate: '2021-09-24',
      sitEntryDate: '2021-09-03',
      status: SERVICE_ITEM_STATUSES.APPROVED,
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
};

export const noSITShipment = {
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
    postalCode: '94535',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
    streetAddress2: 'P.O. Box 9876',
    streetAddress3: 'c/o Some Person',
  },
  eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy41NDY2ODVa',
  id: 'f39ba92d-7d42-446a-be70-3a97b5f9f081',
  moveTaskOrderID: '55c1cdbb-95b0-47f0-ab17-cfe5a0e46ab8',
  pickupAddress: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMS0wOS0yMlQxNDo0ODozNy41MjQyMDla',
    id: '1a42d854-f0b8-428b-a0cc-a43c73a5df8f',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
  primeActualWeight: 980,
  requestedDeliveryDate: '2020-03-15',
  requestedPickupDate: '2020-03-15',
  scheduledPickupDate: '2020-03-16',
  shipmentType: 'HHG',
  status: 'APPROVED',
  updatedAt: '2021-09-22T14:48:37.546Z',
};

const mtoServiceItemsWithSIT = [
  {
    SITPostalCode: '12345',
    createdAt: '2023-05-01T21:09:47.239Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0wNS0wMVQyMTowOTo0Ny4yMzk0ODJa',
    id: 'def9090f-7943-4c10-8383-f25faaff5835',
    moveTaskOrderID: '6fa2b2f6-ef29-4e09-8c0b-e30b574667ec',
    mtoShipmentID: 'b91abd5d-9893-46ff-97e9-75923ec8c750',
    reServiceCode: 'DOASIT',
    reServiceID: '05eb6ff1-5cf6-4918-b887-8260dda6b9fe',
    reServiceName: "Domestic origin add'l SIT",
    reason: 'I need my stuff delivered',
    sitEntryDate: '2023-04-24T00:00:00.000Z',
    status: 'SUBMITTED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
  {
    SITPostalCode: '12345',
    createdAt: '2023-05-01T21:09:47.243Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0wNS0wMVQyMTowOTo0Ny4yNDMxOTZa',
    id: 'ebf82a87-a7ae-40dd-969e-0664b76ccbea',
    moveTaskOrderID: '6fa2b2f6-ef29-4e09-8c0b-e30b574667ec',
    mtoShipmentID: 'b91abd5d-9893-46ff-97e9-75923ec8c750',
    reServiceCode: 'DOPSIT',
    reServiceID: 'd1a4f062-0ca3-4387-8f8e-3dd20493d0b7',
    reServiceName: 'Domestic origin SIT pickup',
    reason: 'I need my stuff delivered',
    sitEntryDate: '2023-04-24T00:00:00.000Z',
    status: 'SUBMITTED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
  {
    SITPostalCode: '12345',
    createdAt: '2023-05-01T21:09:47.180Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0wNS0wMVQyMTowOTo0Ny4yNDU4NzRa',
    id: 'ea9f1d19-438f-4bb2-a9d5-0328a9474433',
    moveTaskOrderID: '6fa2b2f6-ef29-4e09-8c0b-e30b574667ec',
    mtoShipmentID: 'b91abd5d-9893-46ff-97e9-75923ec8c750',
    reServiceCode: 'DOFSIT',
    reServiceID: '998beda7-e390-4a83-b15e-578a24326937',
    reServiceName: 'Domestic origin 1st day SIT',
    reason: 'I need my stuff delivered',
    sitEntryDate: '2023-04-24T00:00:00.000Z',
    status: 'SUBMITTED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
];

const mtoServiceItemsWithFutureSIT = [
  {
    SITPostalCode: '12345',
    createdAt: '2023-05-01T21:09:47.239Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0wNS0wMVQyMTowOTo0Ny4yMzk0ODJa',
    id: 'def9090f-7943-4c10-8383-f25faaff5835',
    moveTaskOrderID: '6fa2b2f6-ef29-4e09-8c0b-e30b574667ec',
    mtoShipmentID: 'b91abd5d-9893-46ff-97e9-75923ec8c750',
    reServiceCode: 'DOASIT',
    reServiceID: '05eb6ff1-5cf6-4918-b887-8260dda6b9fe',
    reServiceName: "Domestic origin add'l SIT",
    reason: 'I need my stuff delivered',
    sitEntryDate: '2025-02-25T00:00:00.000Z',
    status: 'SUBMITTED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
  {
    SITPostalCode: '12345',
    createdAt: '2023-05-01T21:09:47.243Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0wNS0wMVQyMTowOTo0Ny4yNDMxOTZa',
    id: 'ebf82a87-a7ae-40dd-969e-0664b76ccbea',
    moveTaskOrderID: '6fa2b2f6-ef29-4e09-8c0b-e30b574667ec',
    mtoShipmentID: 'b91abd5d-9893-46ff-97e9-75923ec8c750',
    reServiceCode: 'DOPSIT',
    reServiceID: 'd1a4f062-0ca3-4387-8f8e-3dd20493d0b7',
    reServiceName: 'Domestic origin SIT pickup',
    reason: 'I need my stuff delivered',
    sitEntryDate: '2025-02-25T00:00:00.000Z',
    status: 'SUBMITTED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
  {
    SITPostalCode: '12345',
    createdAt: '2023-05-01T21:09:47.180Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0wNS0wMVQyMTowOTo0Ny4yNDU4NzRa',
    id: 'ea9f1d19-438f-4bb2-a9d5-0328a9474433',
    moveTaskOrderID: '6fa2b2f6-ef29-4e09-8c0b-e30b574667ec',
    mtoShipmentID: 'b91abd5d-9893-46ff-97e9-75923ec8c750',
    reServiceCode: 'DOFSIT',
    reServiceID: '998beda7-e390-4a83-b15e-578a24326937',
    reServiceName: 'Domestic origin 1st day SIT',
    reason: 'I need my stuff delivered',
    sitEntryDate: '2025-02-25T00:00:00.000Z',
    status: 'SUBMITTED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
];

export const SITShipment = {
  ...noSITShipment,
  sitStatus: {
    daysInSIT: 15,
    location: LOCATION_VALUES.DESTINATION,
    sitDepartureDate: '2023-03-11',
    sitEntryDate: '2023-04-24',
    totalDaysRemaining: 210,
    totalSITDaysUsed: 270,
    calculatedTotalDaysInSIT: 270,
  },
  sitDaysAllowance: 270,
  mtoServiceItems: mtoServiceItemsWithSIT,
};

export const futureSITShipment = {
  ...noSITShipment,
  sitDaysAllowance: 15,
  mtoServiceItems: mtoServiceItemsWithFutureSIT,
  sitStatus: futureSITStatus,
};

export const futureSITShipmentSITExtension = {
  ...noSITShipment,
  sitDaysAllowance: 15,
  mtoServiceItems: mtoServiceItemsWithFutureSIT,
  sitStatus: futureSITStatus,
  sitExtensions: [
    {
      status: SIT_EXTENSION_STATUS.PENDING,
    },
  ],
};

export const SITStatusExpired = {
  totalSITDaysUsed: 270,
  totalDaysRemaining: -2,
  calculatedTotalDaysInSIT: 270,
  currentSIT: {
    location: LOCATION_VALUES.DESTINATION,
    daysInSIT: 15,
    sitEntryDate: '2021-08-13',
    sitAuthorizedEndDate: '2021-08-28',
  },
};
