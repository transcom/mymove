/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { generatePath } from 'react-router-dom';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingMoveDetails from './ServicesCounselingMoveDetails';

import MOVE_STATUSES from 'constants/moves';
import { ERROR_RETURN_VALUE, LOADING_RETURN_VALUE, INACCESSIBLE_RETURN_VALUE } from 'utils/test/api';
import { ORDERS_TYPE, ORDERS_TYPE_DETAILS } from 'constants/orders';
import { servicesCounselingRoutes } from 'constants/routes';
import { permissionTypes } from 'constants/permissions';
import { SHIPMENT_OPTIONS_URL } from 'shared/constants';
import { useMoveDetailsQueries } from 'hooks/queries';
import { formatDateWithUTC } from 'shared/dates';
import { MockProviders } from 'testUtils';
import { updateMTOShipment, updateMoveStatusServiceCounselingCompleted } from 'services/ghcApi';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const mockRequestedMoveCode = 'LR4T8V';
const mockRoutingParams = { moveCode: mockRequestedMoveCode };
const mockRoutingOptions = { path: servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, params: mockRoutingParams };

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('hooks/queries', () => ({
  useMoveDetailsQueries: jest.fn(),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateMTOShipment: jest.fn(),
  updateMoveStatusServiceCounselingCompleted: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const mtoShipments = [
  {
    customerRemarks: 'please treat gently',
    counselorRemarks: 'all good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
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
    secondaryPickupAddress: {
      city: 'Los Angeles',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: 'b941a74a-e77e-4575-bea3-e7e01b226422',
      postalCode: '90222',
      state: 'CA',
      streetAddress1: '456 Any Street',
      streetAddress2: 'P.O. Box 67890',
      streetAddress3: 'c/o A Friendly Person',
    },
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-eeee-c0f467d13c19',
      postalCode: '90215',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    requestedPickupDate: '2020-06-04',
    scheduledPickupDate: '2020-06-05',
    shipmentType: 'HHG',
    status: 'SUBMITTED',
    updatedAt: '2020-05-10T15:58:02.404031Z',
  },
  {
    customerRemarks: 'do not drop!',
    counselorRemarks: 'looks good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      id: '672ff379-f6e3-48b4-a87d-752463f8f997',
      postalCode: '94534',
      state: 'CA',
      streetAddress1: '111 Everywhere',
      streetAddress2: 'Apt #1',
      streetAddress3: '',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-8799-8a8d-edb60f2a4aee',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Austin',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c55',
      postalCode: '78712',
      state: 'TX',
      streetAddress1: '888 Lucky Street',
      streetAddress2: '#4',
      streetAddress3: 'c/o rabbit',
    },
    requestedPickupDate: '2020-06-05',
    scheduledPickupDate: '2020-06-06',
    shipmentType: 'HHG',
    status: 'SUBMITTED',
    updatedAt: '2020-05-15T15:58:02.404031Z',
  },
];

const ntsrShipmentMissingRequiredInfo = {
  shipmentType: 'HHG_OUTOF_NTS',
  ntsRecordedWeight: 2000,
  id: 'ce01a5b8-9b44-8799-8a8d-edb60f2a4aee',
  serviceOrderNumber: '12341234',
  requestedDeliveryDate: '26 Mar 2020',
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  agents: [
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  counselorRemarks:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ' +
    'Morbi porta nibh nibh, ac malesuada tortor egestas.',
  customerRemarks: 'Ut enim ad minima veniam',
  sacType: 'NTS',
};

const orderMissingRequiredInfo = {
  id: '1',
  originDutyLocation: {
    address: {
      streetAddress1: '',
      city: 'Fort Knox',
      state: 'KY',
      postalCode: '40121',
    },
  },
  destinationDutyLocation: {
    address: {
      streetAddress1: '',
      city: 'Fort Irwin',
      state: 'CA',
      postalCode: '92310',
    },
  },
  customer: {
    agency: 'ARMY',
    backup_contact: {
      email: 'email@example.com',
      name: 'name',
      phone: '555-555-5555',
    },
    current_address: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
      id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    dodID: '6833908165',
    eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
    email: 'combo@ppm.hhg',
    first_name: 'Submitted',
    id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
    last_name: 'Ppmhhg',
    phone: '555-555-5555',
  },
  entitlement: {
    authorizedWeight: 8000,
    dependentsAuthorized: true,
    eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
    id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
    nonTemporaryStorage: true,
    privatelyOwnedVehicle: true,
    proGearWeight: 1,
    proGearWeightSpouse: 500,
    storageInTransit: 2,
    totalDependents: 1,
    totalWeight: 8000,
  },
  orderDocuments: undefined,
  order_number: '',
  order_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  order_type_detail: '',
  department_indicator: '',
  tac: '',
};

const newMoveDetailsQuery = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
  },
  closeoutOffice: undefined,
  customerData: {
    id: '2468',
    last_name: 'Kerry',
    first_name: 'Smith',
    dodID: '999999999',
    agency: 'NAVY',
    backupAddress: {
      streetAddress1: '813 S 129th St',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78234',
    },
  },
  order: {
    id: '1',
    originDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
        isOconus: true,
      },
    },
    destinationDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 1,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    order_type_detail: ORDERS_TYPE_DETAILS.HHG_PERMITTED,
    department_indicator: 'ARMY',
    tac: '9999',
  },
  orderDocuments: {
    z: {
      bytes: 2202009,
      contentType: 'application/pdf',
      createdAt: '2024-10-23T16:31:21.085Z',
      filename: 'testFile.pdf',
      id: 'z',
      status: 'PROCESSING',
      updatedAt: '2024-10-23T16:31:21.085Z',
      uploadType: 'USER',
      url: '/storage/USER/uploads/z?contentType=application%2Fpdf',
    },
  },
  mtoShipments,
  mtoServiceItems: [],
  mtoAgents: [],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const zeroIncentiveMoveDetailsQuery = {
  ...newMoveDetailsQuery,
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
  },
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: null,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 0,
        estimatedWeight: 400,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: false,
        hasReceivedAdvance: null,
        hasRequestedAdvance: false,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: null,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: false,
        hasSecondaryDestinationAddress: false,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: null,
        status: 'SUBMITTED',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 200 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '2023-01-10',
      requestedPickupDate: '2023-01-10',
      shipmentType: 'PPM',
      status: 'SUBMITTED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};

const counselingCompletedMoveDetailsQuery = {
  ...newMoveDetailsQuery,
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED,
  },
};

const ppmShipmentQuery = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'NEEDS_CLOSEOUT',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'APPROVED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: 'e33a1a7b-530f-4df4-b947-d3d719786385',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        shipmentId: 'e33a1a7b-530f-4df4-b947-d3d719786385',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'NEEDS_CLOSEOUT',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: null,
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'APPROVED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
    {
      actualPickupDate: '2020-03-16',
      createdAt: '2022-11-08T23:44:58.306Z',
      customerRemarks: 'Please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4zMDQxOTRa',
        id: '290f7c0a-867f-4d33-83e2-a465dcd83423',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4zMDY2Mzda',
      id: 'a335b359-96cd-4123-8d07-d2270ebaa95c',
      moveTaskOrderID: 'bfd1e5ad-bcbe-4a67-a8e5-4436396cc946',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4zMDE5Njha',
        id: 'da8852a6-9a77-4e8b-b943-f65e616fbc7a',
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
      status: 'SUBMITTED',
      updatedAt: '2022-11-08T23:44:58.306Z',
    },
  ],
};

const ppmShipmentQueryNeedsCloseout = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'NEEDS_CLOSEOUT',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'APPROVED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};

const ppmShipmentQueryWaitingOnCustomer = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'WAITING_ON_CUSTOMER',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'APPROVED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};

const ppmShipmentQuerySubmitted = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'SUBMITTED',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'SUBMITTED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};

const ppmShipmentQueryCloseoutComplete = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'CLOSEOUT_COMPLETE',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'APPROVED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};

const ppmShipmentQueryCancelled = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'CANCELED',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'CANCELED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};

const ppmShipmentQueryDraft = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'DRAFT',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'DRAFT',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};

const ppmShipmentQueryNeedsAdvanceApproval = {
  ...newMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMTc4MVo=',
      id: '167985a7-6d47-4412-b620-d4b7f98a09ed',
      moveTaskOrderID: 'ddf94b4f-db77-4916-83ff-0d6bc68c8b42',
      ppmShipment: {
        actualDestinationPostalCode: null,
        actualMoveDate: null,
        actualPickupPostalCode: null,
        advanceAmountReceived: null,
        advanceAmountRequested: 598700,
        approvedAt: null,
        createdAt: '2022-11-08T23:44:58.226Z',
        eTag: 'MjAyMi0xMS0wOFQyMzo0NDo1OC4yMjY0NTNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: null,
        hasRequestedAdvance: true,
        id: '79b98a71-158d-4b04-9a6c-25543c52183d',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        hasSecondaryPickupAddress: true,
        hasSecondaryDestinationAddress: true,
        pickupAddress: {
          streetAddress1: '111 Test Street',
          streetAddress2: '222 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        secondaryPickupAddress: {
          streetAddress1: '777 Test Street',
          streetAddress2: '888 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42702',
        },
        destinationAddress: {
          streetAddress1: '222 Test Street',
          streetAddress2: '333 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42703',
        },
        secondaryDestinationAddress: {
          streetAddress1: '444 Test Street',
          streetAddress2: '555 Test Street',
          streetAddress3: 'Test Man',
          city: 'Test City',
          state: 'KY',
          postalCode: '42701',
        },
        shipmentId: '167985a7-6d47-4412-b620-d4b7f98a09ed',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'NEEDS_ADVANCE_APPROVAL',
        submittedAt: null,
        updatedAt: '2022-11-08T23:44:58.226Z',
        weightTickets: [{ emptyWeight: 0, fullWeight: 20000 }],
      },
      primeActualWeight: 980,
      requestedDeliveryDate: '0001-01-01',
      requestedPickupDate: '0001-01-01',
      shipmentType: 'PPM',
      status: 'SUBMITTED',
      updatedAt: '2022-11-08T23:44:58.217Z',
    },
  ],
};
const disabledMoveStatuses = [
  MOVE_STATUSES.DRAFT,
  MOVE_STATUSES.SUBMITTED,
  MOVE_STATUSES.APPROVED,
  MOVE_STATUSES.CANCELED,
  MOVE_STATUSES.APPROVALS_REQUESTED,
];

const renderComponent = (props, permissions = [permissionTypes.updateShipment, permissionTypes.updateCustomer]) => {
  return render(
    <MockProviders permissions={permissions} {...mockRoutingOptions}>
      <ServicesCounselingMoveDetails
        setUnapprovedShipmentCount={jest.fn()}
        setMissingOrdersInfoCount={jest.fn()}
        setShipmentWarnConcernCount={jest.fn()}
        setShipmentErrorConcernCount={jest.fn()}
        {...props}
      />
    </MockProviders>,
  );
};

describe('MoveDetails page', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useMoveDetailsQueries.mockReturnValue({ ...newMoveDetailsQuery, ...LOADING_RETURN_VALUE });

      renderComponent();

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMoveDetailsQueries.mockReturnValue({ ...newMoveDetailsQuery, ...ERROR_RETURN_VALUE });

      renderComponent();

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });

    it('renders the Inaccessible component when the query returns an inaccessible response', async () => {
      useMoveDetailsQueries.mockReturnValue({ ...newMoveDetailsQuery, ...INACCESSIBLE_RETURN_VALUE });

      renderComponent();

      const errorMessage = await screen.getByText(/Page is not accessible./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('renders the h1', async () => {
      useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

      renderComponent();

      expect(await screen.findByRole('heading', { name: 'Move Details', level: 1 })).toBeInTheDocument();
    });

    it.each([['Shipments'], ['Orders'], ['Allowances'], ['Customer info']])(
      'renders side navigation for section %s',
      async (sectionName) => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

        renderComponent();

        expect(await screen.findByRole('link', { name: sectionName })).toBeInTheDocument();
      },
    );

    it('renders warnings and errors on left nav bar for all shipments in a section', async () => {
      const moveDetailsQuery = {
        ...newMoveDetailsQuery,
        mtoShipments: [ntsrShipmentMissingRequiredInfo],
      };

      useMoveDetailsQueries.mockReturnValue(moveDetailsQuery);

      renderComponent();

      expect(await screen.findByTestId('requestedShipmentsTag')).toBeInTheDocument();
      expect(await screen.findByTestId('shipment-missing-info-alert')).toBeInTheDocument();
    });

    it('shares the number of missing orders information', () => {
      const moveDetailsQuery = {
        ...newMoveDetailsQuery,
        order: orderMissingRequiredInfo,
      };

      useMoveDetailsQueries.mockReturnValue(moveDetailsQuery);

      const mockSetMissingOrdersInfoCount = jest.fn();
      renderComponent({ setMissingOrdersInfoCount: mockSetMissingOrdersInfoCount });

      // Should have called `setMissingOrdersInfoCount` with 4 missing fields
      expect(mockSetMissingOrdersInfoCount).toHaveBeenCalledTimes(1);
      expect(mockSetMissingOrdersInfoCount).toHaveBeenCalledWith(4);
    });

    /* eslint-disable camelcase */
    it('renders shipments info', async () => {
      useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

      renderComponent();

      expect(await screen.findByRole('heading', { name: 'Shipments', level: 2 })).toBeInTheDocument();

      expect(screen.getAllByRole('heading', { name: 'HHG', level: 3 }).length).toBe(2);

      const moveDateTerms = screen.getAllByText('Requested pickup date');

      expect(moveDateTerms.length).toBe(2);

      for (let i = 0; i < moveDateTerms.length; i += 1) {
        expect(moveDateTerms[i].nextElementSibling.textContent).toBe(
          formatDateWithUTC(newMoveDetailsQuery.mtoShipments[i].requestedPickupDate, 'DD MMM YYYY'),
        );
      }

      const originAddressTerms = screen.getAllByText('Pickup Address');

      expect(originAddressTerms.length).toBe(3);

      for (let i = 0; i < 2; i += 1) {
        const { streetAddress1, city, state, postalCode } = newMoveDetailsQuery.mtoShipments[i].pickupAddress;

        const addressText = originAddressTerms[i].nextElementSibling.textContent;

        expect(addressText).toContain(streetAddress1);
        expect(addressText).toContain(city);
        expect(addressText).toContain(state);
        expect(addressText).toContain(postalCode);
      }

      const destinationAddressTerms = screen.getAllByText('Delivery Address');

      expect(destinationAddressTerms.length).toBe(2);

      for (let i = 0; i < destinationAddressTerms.length; i += 1) {
        const { streetAddress1, city, state, postalCode } = newMoveDetailsQuery.mtoShipments[i].destinationAddress;

        const addressText = destinationAddressTerms[i].nextElementSibling.textContent;

        expect(addressText).toContain(streetAddress1);
        expect(addressText).toContain(city);
        expect(addressText).toContain(state);
        expect(addressText).toContain(postalCode);
      }

      const counselorRemarksTerms = screen.getAllByText('Counselor remarks');

      expect(counselorRemarksTerms.length).toBe(2);

      for (let i = 0; i < counselorRemarksTerms.length; i += 1) {
        expect(counselorRemarksTerms[i].nextElementSibling.textContent).toBe(
          newMoveDetailsQuery.mtoShipments[i].counselorRemarks || '—',
        );
      }
    });

    it('renders review documents button', async () => {
      useMoveDetailsQueries.mockReturnValue(ppmShipmentQuery);
      renderComponent();
      expect(screen.getAllByRole('button', { name: 'Review documents' }).length).toBe(2);
    });

    it('renders review shipment weights button with correct path', async () => {
      useMoveDetailsQueries.mockReturnValue(ppmShipmentQuery);
      const path = generatePath(servicesCounselingRoutes.BASE_REVIEW_SHIPMENT_WEIGHTS_PATH, {
        moveCode: mockRequestedMoveCode,
      });
      renderComponent();

      const reviewShipmentWeightsBtn = screen.getByRole('button', { name: 'Review shipment weights' });

      expect(reviewShipmentWeightsBtn).toBeInTheDocument();
      expect(reviewShipmentWeightsBtn.getAttribute('data-testid')).toBe(path);
    });

    it('shows an error if there is an advance requested and no advance status for a PPM shipment', async () => {
      useMoveDetailsQueries.mockReturnValue(ppmShipmentQuery);
      renderComponent();

      const advanceStatusElement = screen.getAllByTestId('advanceRequestStatus')[0];
      expect(advanceStatusElement.parentElement).toHaveClass('missingInfoError');
    });

    it('renders the excess weight alert and additional shipment concern if there is excess weight', async () => {
      useMoveDetailsQueries.mockReturnValue(ppmShipmentQuery);
      renderComponent();
      const excessWeightAlert = screen.getByText(
        'This move has excess weight. Review PPM weight ticket documents to resolve.',
      );
      expect(excessWeightAlert).toBeInTheDocument();

      expect(await screen.findByTestId('requestedShipmentsTag')).toBeInTheDocument();
    });

    it('renders the allowances error message when allowances are less than moves values', async () => {
      useMoveDetailsQueries.mockReturnValue(ppmShipmentQuery);
      renderComponent();
      const allowanceError = screen.getByTestId('allowanceError');
      expect(allowanceError).toBeInTheDocument();
    });

    it('renders shipments info even if delivery address is missing', async () => {
      const moveDetailsQuery = {
        ...newMoveDetailsQuery,
        mtoShipments: [
          // Want to create a "new" mtoShipment to be able to delete things without messing up existing tests
          { ...newMoveDetailsQuery.mtoShipments[0] },
          newMoveDetailsQuery.mtoShipments[1],
        ],
      };

      delete moveDetailsQuery.mtoShipments[0].destinationAddress;

      useMoveDetailsQueries.mockReturnValue(moveDetailsQuery);

      renderComponent();

      const destinationAddressTerms = screen.getAllByText('Delivery Address');

      expect(destinationAddressTerms.length).toBe(2);

      expect(destinationAddressTerms[0].nextElementSibling.textContent).toBe(
        moveDetailsQuery.order.destinationDutyLocation.address.postalCode,
      );

      const { streetAddress1, city, state, postalCode } = moveDetailsQuery.mtoShipments[1].destinationAddress;

      const addressText = destinationAddressTerms[1].nextElementSibling.textContent;

      expect(addressText).toContain(streetAddress1);
      expect(addressText).toContain(city);
      expect(addressText).toContain(state);
      expect(addressText).toContain(postalCode);
    });
    /* eslint-enable camelcase */

    it('renders customer info', async () => {
      useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

      renderComponent();

      expect(await screen.findByRole('heading', { name: 'Customer info', level: 2 })).toBeInTheDocument();
    });

    it('renders info saved alert', () => {
      renderComponent({
        infoSavedAlert: { alertType: 'success', message: 'great success!' },
        setUnapprovedShipmentCount: jest.fn(),
      });
      expect(screen.getByText('great success!')).toBeInTheDocument();
    });

    describe('new move - needs service counseling', () => {
      it('submit move details button is on page', async () => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

        renderComponent();

        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
      });

      it('submit move details button is disabled when there are no shipments', async () => {
        useMoveDetailsQueries.mockReturnValue({ ...newMoveDetailsQuery, mtoShipments: [] });

        renderComponent();

        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeDisabled();
      });

      it('submit move details button is disabled when all shipments are deleted', async () => {
        const deletedMtoShipments = mtoShipments.map((shipment) => ({ ...shipment, deletedAt: new Date() }));
        useMoveDetailsQueries.mockReturnValue({
          ...newMoveDetailsQuery,
          mtoShipments: deletedMtoShipments,
        });

        renderComponent();

        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeDisabled();
      });

      it('submit move details button is disabled when required orders information is missing', async () => {
        useMoveDetailsQueries.mockReturnValue({
          ...newMoveDetailsQuery,
          order: {
            ...newMoveDetailsQuery.order,
            department_indicator: undefined,
          },
        });

        renderComponent();

        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeDisabled();
      });
      it('should disable submit move details button when order documents is missing', async () => {
        useMoveDetailsQueries.mockReturnValue({
          ...newMoveDetailsQuery,
          orderDocuments: undefined,
        });
        renderComponent();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).not.toBeEnabled();
      });
      it('submit move details button is not disabled when some shipments are deleted', async () => {
        const deletedMtoShipments = mtoShipments.map((shipment, index) => {
          if (index > 0) {
            return { ...shipment, deletedAt: new Date() };
          }
          return shipment;
        });
        useMoveDetailsQueries.mockReturnValue({
          ...newMoveDetailsQuery,
          mtoShipments: deletedMtoShipments,
        });

        renderComponent();

        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).not.toBeDisabled();
      });

      it('buttons are disabled and links are not rendered when move is locked', async () => {
        const deletedMtoShipments = mtoShipments.map((shipment, index) => {
          if (index > 0) {
            return { ...shipment, deletedAt: new Date() };
          }
          return shipment;
        });
        const isMoveLocked = true;
        useMoveDetailsQueries.mockReturnValue({
          ...newMoveDetailsQuery,
          mtoShipments: deletedMtoShipments,
        });

        render(
          <MockProviders
            permissions={[permissionTypes.updateShipment, permissionTypes.updateCustomer]}
            {...mockRoutingOptions}
          >
            <ServicesCounselingMoveDetails
              setUnapprovedShipmentCount={jest.fn()}
              setMissingOrdersInfoCount={jest.fn()}
              setShipmentWarnConcernCount={jest.fn()}
              setShipmentErrorConcernCount={jest.fn()}
              isMoveLocked={isMoveLocked}
            />
          </MockProviders>,
        );

        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeDisabled();
        expect(screen.queryByRole('combobox')).not.toBeInTheDocument(); // Add a new shipment ButtonDropdown

        expect(screen.queryByRole('link', { name: 'View and edit orders' })).not.toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).not.toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).not.toBeInTheDocument();
      });

      it('submit move details button is disabled when a shipment has missing information', async () => {
        const moveDetailsQuery = {
          ...newMoveDetailsQuery,
          mtoShipments: [ntsrShipmentMissingRequiredInfo],
        };
        useMoveDetailsQueries.mockReturnValue(moveDetailsQuery);

        renderComponent();

        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(await screen.findByRole('button', { name: 'Submit move details' })).toBeDisabled();
      });

      it('renders the Orders Definition List', async () => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

        renderComponent();

        expect(await screen.findByRole('heading', { name: 'Orders', level: 2 })).toBeInTheDocument();
        expect(screen.getByText('Current duty location')).toBeInTheDocument();
      });

      it('renders the Allowances Table', async () => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

        renderComponent();

        expect(await screen.findByRole('heading', { name: 'Allowances', level: 2 })).toBeInTheDocument();
        expect(screen.getByText('Branch')).toBeInTheDocument();
      });

      it('allows the service counselor to use the modal as expected', async () => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);
        updateMoveStatusServiceCounselingCompleted.mockImplementation(() => Promise.resolve({}));

        renderComponent();

        const submitButton = await screen.findByRole('button', { name: 'Submit move details' });

        await userEvent.click(submitButton);

        expect(await screen.findByRole('heading', { name: 'Are you sure?', level: 2 })).toBeInTheDocument();

        const modalSubmitButton = screen.getByRole('button', { name: 'Yes, submit' });

        await userEvent.click(modalSubmitButton);

        await waitFor(() => {
          expect(screen.queryByRole('heading', { name: 'Are you sure?', level: 2 })).not.toBeInTheDocument();
        });
      });

      it('allows the service counselor to submit details for ppm with zero incentive', async () => {
        useMoveDetailsQueries.mockReturnValue(zeroIncentiveMoveDetailsQuery);
        updateMTOShipment.mockImplementation(() => Promise.resolve({}));
        updateMoveStatusServiceCounselingCompleted.mockImplementation(() => Promise.resolve({}));

        renderComponent();

        const submitButton = await screen.findByRole('button', { name: 'Submit move details' });

        await userEvent.click(submitButton);

        expect(await screen.findByRole('heading', { name: 'Are you sure?', level: 2 })).toBeInTheDocument();

        const modalSubmitButton = screen.getByRole('button', { name: 'Yes, submit' });

        await userEvent.click(modalSubmitButton);

        await waitFor(() => {
          expect(screen.queryByRole('heading', { name: 'Are you sure?', level: 2 })).not.toBeInTheDocument();
        });
      });

      it.each([
        ['View and edit orders', servicesCounselingRoutes.ORDERS_EDIT_PATH],
        ['Edit allowances', servicesCounselingRoutes.ALLOWANCES_EDIT_PATH],
        ['Edit customer info', servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH],
      ])('shows the "%s" link as expected: %s', async (linkText, route) => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

        renderComponent();

        const link = await screen.findByRole('link', { name: linkText });
        expect(link).toBeInTheDocument();

        const path = `/${generatePath(route, {
          moveCode: mockRequestedMoveCode,
        })}`;
        expect(link).toHaveAttribute('href', path);
      });

      describe('shows the dropdown and navigates to each option', () => {
        it.each([
          [SHIPMENT_OPTIONS_URL.HHG],
          [SHIPMENT_OPTIONS_URL.NTS],
          [SHIPMENT_OPTIONS_URL.NTSrelease],
          [SHIPMENT_OPTIONS_URL.UNACCOMPANIED_BAGGAGE],
        ])('selects the %s option and navigates to the matching form for that shipment type', async (shipmentType) => {
          isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
          render(
            <MockProviders
              path={servicesCounselingRoutes.BASE_SHIPMENT_ADD_PATH}
              params={{ moveCode: mockRequestedMoveCode, shipmentType }}
            >
              <ServicesCounselingMoveDetails
                setUnapprovedShipmentCount={jest.fn()}
                setMissingOrdersInfoCount={jest.fn()}
                setShipmentWarnConcernCount={jest.fn()}
                setShipmentErrorConcernCount={jest.fn()}
              />
            </MockProviders>,
          );

          const path = `../${generatePath(servicesCounselingRoutes.SHIPMENT_ADD_PATH, {
            moveCode: mockRequestedMoveCode,
            shipmentType,
          })}`;

          const buttonDropdown = await screen.findByRole('combobox');

          expect(buttonDropdown).toBeInTheDocument();

          await userEvent.selectOptions(buttonDropdown, shipmentType);

          await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith(path);
          });
        });
      });

      it('does not allow user to add UB shipment for local move orders', async () => {
        isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
        const localMoveDetailsQuery = {
          ...newMoveDetailsQuery,
          order: {
            ...newMoveDetailsQuery.order,
            order_type: ORDERS_TYPE.LOCAL_MOVE,
          },
        };
        useMoveDetailsQueries.mockReturnValue(localMoveDetailsQuery);

        renderComponent();

        const combobox = await screen.getByRole('combobox', { name: 'Add a new shipment' });
        expect(combobox).toBeInTheDocument();

        await userEvent.click(combobox);

        // Check if all expected options appear
        await waitFor(() => {
          expect(screen.getByRole('option', { name: 'HHG' })).toBeInTheDocument();
          expect(screen.getByRole('option', { name: 'PPM' })).toBeInTheDocument();
          expect(screen.getByRole('option', { name: 'NTS' })).toBeInTheDocument();
          expect(screen.getByRole('option', { name: 'NTS-release' })).toBeInTheDocument();
          expect(screen.getByRole('option', { name: 'Boat' })).toBeInTheDocument();
          expect(screen.getByRole('option', { name: 'Mobile Home' })).toBeInTheDocument();
        });
        expect(screen.queryByRole('option', { name: 'UB' })).not.toBeInTheDocument();
      });

      it('shows the edit shipment buttons', async () => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

        renderComponent();

        const shipmentEditButtons = await screen.findAllByRole('button', { name: 'Edit shipment' });

        expect(shipmentEditButtons.length).toBe(2);

        for (let i = 0; i < shipmentEditButtons.length; i += 1) {
          expect(shipmentEditButtons[i].getAttribute('data-testid')).toBe(
            `../${generatePath(servicesCounselingRoutes.SHIPMENT_EDIT_PATH, {
              moveCode: mockRequestedMoveCode,
              shipmentId: newMoveDetailsQuery.mtoShipments[i].id,
            })}`,
          );
        }
      });

      it('shows the customer and counselor remarks', async () => {
        useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

        renderComponent();

        const customerRemarks1 = await screen.findByText('please treat gently');
        const customerRemarks2 = await screen.findByText('do not drop!');

        const counselorRemarks1 = await screen.findByText('all good');
        const counselorRemarks2 = await screen.findByText('looks good');

        expect(customerRemarks1).toBeInTheDocument();
        expect(customerRemarks2).toBeInTheDocument();
        expect(counselorRemarks1).toBeInTheDocument();
        expect(counselorRemarks2).toBeInTheDocument();
      });
    });

    describe('service counseling completed', () => {
      it('hides submit and view/edit buttons/links', async () => {
        useMoveDetailsQueries.mockReturnValue(counselingCompletedMoveDetailsQuery);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).not.toBeInTheDocument();
        // Add a new shipment ButtonDropdown is now used for PPMs at any point in move process
        expect(screen.queryByRole('combobox')).toBeInTheDocument();
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });
    });

    describe('ppm specific statuses for service counselor', () => {
      // PPM Shipments should show the edit button on shipments where PPM Shipment status is DRAFT, SUBMITTED, or NEEDS_ADVANCE_APPROVAL
      it('hides submit and view/edit buttons/links for NEEDS_CLOSEOUT status', async () => {
        useMoveDetailsQueries.mockReturnValue(ppmShipmentQueryNeedsCloseout);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).not.toBeInTheDocument();
        expect(screen.queryByRole('combobox')).toBeInTheDocument(); // Add a new shipment ButtonDropdown
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });

      it('hides submit and view/edit buttons/links for WAITING_ON_CUSTOMER status', async () => {
        useMoveDetailsQueries.mockReturnValue(ppmShipmentQueryWaitingOnCustomer);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(screen.queryByRole('combobox')).toBeInTheDocument(); // Add a new shipment ButtonDropdown
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });

      it('hides submit and view/edit buttons/links for CLOSEOUT_COMPLETE status', async () => {
        useMoveDetailsQueries.mockReturnValue(ppmShipmentQueryCloseoutComplete);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(screen.queryByRole('combobox')).toBeInTheDocument(); // Add a new shipment ButtonDropdown
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });

      it('hides submit and view/edit buttons/links for CANCELLED status', async () => {
        useMoveDetailsQueries.mockReturnValue(ppmShipmentQueryCancelled);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(screen.queryByRole('combobox')).toBeInTheDocument(); // Add a new shipment ButtonDropdown
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });

      it('shows submit and view/edit buttons/links for SUBMITTED status', async () => {
        useMoveDetailsQueries.mockReturnValue(ppmShipmentQuerySubmitted);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(screen.queryByRole('combobox')).toBeInTheDocument(); // Add a new shipment ButtonDropdown
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });

      it('shows submit and view/edit buttons/links for DRAFT status', async () => {
        useMoveDetailsQueries.mockReturnValue(ppmShipmentQueryDraft);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(screen.queryByRole('combobox')).toBeInTheDocument(); // Add a new shipment ButtonDropdown
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });

      it('shows submit and view/edit buttons/links for NEEDS ADVANCE APPROVAL status', async () => {
        useMoveDetailsQueries.mockReturnValue(ppmShipmentQueryNeedsAdvanceApproval);

        renderComponent();

        expect(screen.queryByRole('button', { name: 'Submit move details' })).toBeInTheDocument();
        expect(screen.queryByRole('combobox')).toBeInTheDocument(); // Add a new shipment ButtonDropdown
        expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'View and edit orders' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit allowances' })).toBeInTheDocument();
        expect(screen.queryByRole('link', { name: 'Edit customer info' })).toBeInTheDocument();
      });
    });
    describe('for view only orders and allowance move statuses', () => {
      it('shows view buttons instead of edit', async () => {
        for (let i = 0; i < disabledMoveStatuses.length; i += 1) {
          const counselingCompletedMoveDetailsQueryValues = JSON.parse(
            JSON.stringify(counselingCompletedMoveDetailsQuery),
          );
          counselingCompletedMoveDetailsQueryValues.move = {
            id: 123,
            moveCode: 'GLOBAL123',
            ordersId: 1,
            status: disabledMoveStatuses[i],
          };

          // set return values for mocked functions
          useMoveDetailsQueries.mockReturnValue(counselingCompletedMoveDetailsQueryValues);
          renderComponent();

          const viewOrders = screen.queryAllByRole('link', { name: 'View orders' });
          expect(viewOrders[0]).toBeInTheDocument();
          const viewAllowances = screen.queryAllByRole('link', { name: 'View allowances' });
          expect(viewAllowances[0]).toBeInTheDocument();
        }
      });
    });

    describe('permission dependent rendering', () => {
      useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);

      it('renders the financial review flag button when user has permission', async () => {
        render(
          <MockProviders permissions={[permissionTypes.updateFinancialReviewFlag]} {...mockRoutingOptions}>
            <ServicesCounselingMoveDetails
              setUnapprovedShipmentCount={jest.fn()}
              setMissingOrdersInfoCount={jest.fn()}
              setShipmentWarnConcernCount={jest.fn()}
              setShipmentErrorConcernCount={jest.fn()}
            />
          </MockProviders>,
        );

        expect(await screen.getByText('Flag move for financial review')).toBeInTheDocument();
      });

      it('does not show the financial review flag button if user does not have permission', () => {
        render(
          <MockProviders {...mockRoutingOptions}>
            <ServicesCounselingMoveDetails setUnapprovedShipmentCount={jest.fn()} />
          </MockProviders>,
        );

        expect(screen.queryByText('Flag move for financial review')).not.toBeInTheDocument();
      });

      it('does not show the edit customer info button if user does not have permission', () => {
        render(
          <MockProviders {...mockRoutingOptions}>
            <ServicesCounselingMoveDetails setUnapprovedShipmentCount={jest.fn()} />
          </MockProviders>,
        );

        expect(screen.queryByText('Edit customer info')).not.toBeInTheDocument();
      });

      it('renders the cancel move button when user has permission', async () => {
        render(
          <MockProviders permissions={[permissionTypes.cancelMoveFlag]} {...mockRoutingOptions}>
            <ServicesCounselingMoveDetails
              setUnapprovedShipmentCount={jest.fn()}
              setMissingOrdersInfoCount={jest.fn()}
              setShipmentWarnConcernCount={jest.fn()}
              setShipmentErrorConcernCount={jest.fn()}
            />
          </MockProviders>,
        );

        expect(await screen.getByText('Cancel move')).toBeInTheDocument();
      });

      it('does not show the cancel move button if user does not have permission', () => {
        render(
          <MockProviders {...mockRoutingOptions}>
            <ServicesCounselingMoveDetails setUnapprovedShipmentCount={jest.fn()} />
          </MockProviders>,
        );

        expect(screen.queryByText('Cancel move')).not.toBeInTheDocument();
      });
    });
  });
});
