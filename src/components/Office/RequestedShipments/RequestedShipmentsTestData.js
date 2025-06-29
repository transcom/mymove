import { ORDERS_TYPE, ORDERS_BRANCH_OPTIONS, ORDERS_PAY_GRADE_TYPE } from '../../../constants/orders';
import { DEPARTMENT_INDICATOR_OPTIONS } from '../../../constants/departmentIndicators';

import { SHIPMENT_OPTIONS, MTOAgentType, PPM_TYPES } from 'shared/constants';

export const shipments = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'looks good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'looks good',
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'cd01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.431993Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'looks good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTcyMjZa',
      id: '00a5dfeb-c6a0-4ed8-965c-89943163fee4',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MzE5OTVa',
    id: 'c2f68d97-b960-4c86-a418-c70a0aeba04e',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTMyNDha',
      id: '14b1d10d-b34b-4ec5-80e6-69d885206a2a',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjYxODVa',
      id: '1a4f6fec-42b9-4dd2-b205-c6770ac7ea27',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjIwNzVa',
      id: 'e188f33f-f84d-4f86-954a-938b52e38741',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.NTS,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
  },
];

export const shipmentsNoApprovedDate = [
  {
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'looks good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
];

export const ntsExternalVendorShipments = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.431993Z',
    customerRemarks: 'please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTcyMjZa',
      id: '00a5dfeb-c6a0-4ed8-965c-89943163fee4',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MzE5OTVa',
    id: 'c2f68d97-b960-4c86-a418-c70a0aeba04e',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTMyNDha',
      id: '14b1d10d-b34b-4ec5-80e6-69d885206a2a',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjYxODVa',
      id: '1a4f6fec-42b9-4dd2-b205-c6770ac7ea27',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjIwNzVa',
      id: 'e188f33f-f84d-4f86-954a-938b52e38741',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.NTS,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
    usesExternalVendor: true,
  },
];

export const ppmOnlyShipments = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'looks good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.PPM,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
];

export const zeroIncentivePPM = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'looks good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    ppmShipment: {
      ppmType: PPM_TYPES.INCENTIVE_BASED,
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '813 S 129th St',
        streetAddress2: '#124',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
      secondaryPickupAddress: {
        streetAddress1: '814 S 129th St',
        streetAddress2: '#125',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      secondaryDestinationAddress: {
        streetAddress1: '815 S 129th St',
        streetAddress2: '#126',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
      hasSecondaryDestinationAddress: true,
      hasSecondaryPickupAddress: true,
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 500,
      estimatedIncentive: 0,
    },
    shipmentType: SHIPMENT_OPTIONS.PPM,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
];

export const ordersInfo = {
  newDutyLocation: {
    address: {
      city: 'Augusta',
      country: 'United States',
      eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
      id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
      postalCode: '30813',
      state: 'GA',
      streetAddress1: 'Fort Gordon',
    },
    address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
    id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
    name: 'Fort Gordon',
  },
  currentDutyLocation: {
    address: {
      city: 'Des Moines',
      country: 'US',
      eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42NjEwODFa',
      id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
      postalCode: '50309',
      state: 'IA',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
    },
    address_id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
    eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42Njg5MDFa',
    id: '07282a8f-a496-4648-ae24-119775eef57d',
    name: 'vC6w22RPYC',
  },
  issuedDate: '2018-03-15',
  reportByDate: '2018-08-01',
  departmentIndicator: DEPARTMENT_INDICATOR_OPTIONS.COAST_GUARD,
  ordersNumber: 'ORDER3',
  ordersType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  ordersTypeDetail: 'TBD',
  ordersDocuments: [
    {
      'c0a22a98-a806-47a2-ab54-2dac938667b3': {
        bytes: 2202009,
        contentType: 'application/pdf',
        createdAt: '2024-10-23T16:31:21.085Z',
        filename: 'testFile.pdf',
        id: 'c0a22a98-a806-47a2-ab54-2dac938667b3',
        status: 'PROCESSING',
        updatedAt: '2024-10-23T16:31:21.085Z',
        uploadType: 'USER',
        url: '/storage/USER/uploads/c0a22a98-a806-47a2-ab54-2dac938667b3?contentType=application%2Fpdf',
      },
    },
  ],
  tacMDC: '',
  sacSDN: '',
};

export const ordersInfoOCONUS = {
  newDutyLocation: {
    address: {
      city: 'Augusta',
      country: 'United States',
      eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
      id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
      postalCode: '30813',
      state: 'GA',
      streetAddress1: 'Fort Gordon',
      isOconus: false,
    },
    address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
    id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
    name: 'Fort Gordon',
  },
  currentDutyLocation: {
    address: {
      city: 'JBER',
      country: 'AK',
      eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42NjEwODFa',
      id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
      postalCode: '99702',
      state: 'AK',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
      isOconus: true,
    },
    address_id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
    eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42Njg5MDFa',
    id: '07282a8f-a496-4648-ae24-119775eef57d',
    name: 'vC6w22RPYC',
  },
  issuedDate: '2018-03-15',
  reportByDate: '2018-08-01',
  departmentIndicator: DEPARTMENT_INDICATOR_OPTIONS.COAST_GUARD,
  ordersNumber: 'ORDER3',
  ordersType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  ordersTypeDetail: 'TBD',
  ordersDocuments: [
    {
      'c0a22a98-a806-47a2-ab54-2dac938667b3': {
        bytes: 2202009,
        contentType: 'application/pdf',
        createdAt: '2024-10-23T16:31:21.085Z',
        filename: 'testFile.pdf',
        id: 'c0a22a98-a806-47a2-ab54-2dac938667b3',
        status: 'PROCESSING',
        updatedAt: '2024-10-23T16:31:21.085Z',
        uploadType: 'USER',
        url: '/storage/USER/uploads/c0a22a98-a806-47a2-ab54-2dac938667b3?contentType=application%2Fpdf',
      },
    },
  ],
  tacMDC: '',
  sacSDN: '',
};

export const ordersInfoOCONUSLocalMove = {
  newDutyLocation: {
    address: {
      city: 'Augusta',
      country: 'United States',
      eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
      id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
      postalCode: '30813',
      state: 'GA',
      streetAddress1: 'Fort Gordon',
      isOconus: false,
    },
    address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
    id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
    name: 'Fort Gordon',
  },
  currentDutyLocation: {
    address: {
      city: 'JBER',
      country: 'AK',
      eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42NjEwODFa',
      id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
      postalCode: '99702',
      state: 'AK',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
      isOconus: true,
    },
    address_id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
    eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42Njg5MDFa',
    id: '07282a8f-a496-4648-ae24-119775eef57d',
    name: 'vC6w22RPYC',
  },
  issuedDate: '2018-03-15',
  reportByDate: '2018-08-01',
  departmentIndicator: DEPARTMENT_INDICATOR_OPTIONS.COAST_GUARD,
  ordersNumber: 'ORDER3',
  ordersType: ORDERS_TYPE.LOCAL_MOVE,
  ordersTypeDetail: 'TBD',
  ordersDocuments: [
    {
      'c0a22a98-a806-47a2-ab54-2dac938667b3': {
        bytes: 2202009,
        contentType: 'application/pdf',
        createdAt: '2024-10-23T16:31:21.085Z',
        filename: 'testFile.pdf',
        id: 'c0a22a98-a806-47a2-ab54-2dac938667b3',
        status: 'PROCESSING',
        updatedAt: '2024-10-23T16:31:21.085Z',
        uploadType: 'USER',
        url: '/storage/USER/uploads/c0a22a98-a806-47a2-ab54-2dac938667b3?contentType=application%2Fpdf',
      },
    },
  ],
  tacMDC: '',
  sacSDN: '',
};

export const allowancesInfo = {
  branch: ORDERS_BRANCH_OPTIONS.NAVY,
  grade: ORDERS_PAY_GRADE_TYPE.E_6,
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
};

export const customerInfo = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  backupContactName: 'Quinn Ocampo',
  backupContactPhone: '+1 999-999-9999',
  backupContactEmail: 'quinnocampo@myemail.com',
};

export const agents = [
  {
    type: MTOAgentType.RELEASING_AGENT,
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
  {
    type: MTOAgentType.RECEIVING_AGENT,
    firstName: 'Dorothy Lagomarsino',
    lastName: 'Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
];

export const serviceItemsMSandCS = [
  {
    reServiceName: 'Move management',
    approvedAt: '2020-01-01',
    id: '76055c99-0990-410c-a7c9-69373b0b53eb',
    status: 'APPROVED',
    reServiceCode: 'MS',
  },
  {
    reServiceName: 'Counseling fee',
    id: '76055c99-0990-410c-a7c9-69373b0b5322',
    status: 'APPROVED',
    reServiceCode: 'CS',
    approvedAt: '2020-01-01',
  },
];

export const serviceItemsMS = [
  {
    reServiceName: 'Move management',
    approvedAt: '2020-01-01',
    id: '76055c99-0990-410c-a7c9-69373b0b53eb',
    status: 'APPROVED',
    reServiceCode: 'MS',
  },
];

export const serviceItemsCS = [
  {
    reServiceName: 'Counseling fee',
    id: '76055c99-0990-410c-a7c9-69373b0b5322',
    status: 'APPROVED',
    reServiceCode: 'CS',
    approvedAt: '2020-01-01',
  },
];

export const serviceItemsEmpty = [];

export const moveTaskOrders = [{}, { serviceCounselingCompletedAt: '2020-10-02T19:20:08.481139Z' }];

export const closeoutOffice = 'Office of Closeout';
