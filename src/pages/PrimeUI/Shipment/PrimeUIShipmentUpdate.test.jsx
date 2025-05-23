import React from 'react';
import { cleanup, render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import PrimeUIShipmentUpdate from './PrimeUIShipmentUpdate';

import { primeSimulatorRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { updatePrimeMTOShipmentV3 } from 'services/primeApi';

const shipmentId = 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee';
const ppmShipmentId = '1b695b60-c3ed-401b-b2e3-808d095eb8cc';
const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';
const moveCode = 'LR4T8V';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = {
  moveCode,
  moveCodeOrID: moveId,
  shipmentId,
};
const ppmRoutingParams = {
  moveCode,
  moveCodeOrID: moveId,
  shipmentId: ppmShipmentId,
};

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  updatePrimeMTOShipmentV2: jest.fn().mockImplementation(() => Promise.resolve()),
  updatePrimeMTOShipmentV3: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));
const approvedMoveTaskOrder = {
  moveTaskOrder: {
    availableToPrimeAt: '2021-10-18T18:24:41.235Z',
    createdAt: '2021-10-18T18:24:41.362Z',
    diversion: false,
    eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNjIxNjRa',
    excessWeightAcknowledgedAt: null,
    excessWeightQualifiedAt: null,
    excessWeightUploadId: null,
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    moveCode: 'LR4T8V',
    mtoShipments: [
      {
        actualPickupDate: '2020-03-16',
        agents: [
          {
            agentType: 'RELEASING_AGENT',
            createdAt: '2021-10-18T18:24:41.521Z',
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41MjE4NzNa',
            email: 'test@test.email.com',
            firstName: 'Test',
            id: 'f2619e1b-7729-4b97-845d-6ae1ebe299f2',
            lastName: 'Agent',
            mtoShipmentID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
            phone: '202-555-9301',
            updatedAt: '2021-10-18T18:24:41.521Z',
          },
        ],
        approvedDate: '2021-10-18',
        createdAt: '2021-10-18T18:24:41.377Z',
        customerRemarks: 'Please treat gently',
        destinationAddress: {
          city: 'Fairfield',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzI3NDJa',
          id: 'bfe61147-5fd7-426e-b473-54ccf77bde35',
          postalCode: '94535',
          state: 'CA',
          streetAddress1: '987 Any Avenue',
          streetAddress2: 'P.O. Box 9876',
          streetAddress3: 'c/o Some Person',
        },
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzc5Nzha',
        firstAvailableDeliveryDate: null,
        id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        pickupAddress: {
          city: 'Beverly Hills',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNjc3Mjda',
          id: 'cf159eca-162c-4131-84a0-795e684416a6',
          postalCode: '90210',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
        },
        primeActualWeight: 2000,
        primeEstimatedWeight: 1400,
        actualProGearWeight: 1500,
        actualSpouseProGearWeight: 250,
        primeEstimatedWeightRecordedDate: null,
        requestedPickupDate: '2020-03-15',
        requiredDeliveryDate: null,
        scheduledPickupDate: '2020-03-16',
        secondaryDeliveryAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        secondaryPickupAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        shipmentType: 'HHG',
        status: 'APPROVED',
        updatedAt: '2021-10-18T18:24:41.377Z',
        mtoServiceItems: null,
      },
      {
        actualPickupDate: null,
        approvedDate: null,
        counselorRemarks: 'These are counselor remarks for a PPM.',
        createdAt: '2022-07-01T13:41:33.261Z',
        destinationAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43MzgzODla',
        firstAvailableDeliveryDate: null,
        id: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        pickupAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        ppmShipment: {
          actualMoveDate: '2022-07-13',
          advanceAmountReceived: 598600,
          advanceAmountRequested: 598700,
          approvedAt: '2022-07-03T14:20:21.620Z',
          createdAt: '2022-06-30T13:41:33.265Z',
          eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43ODA1Mlo=',
          estimatedIncentive: 1000000,
          estimatedWeight: 4000,
          expectedDepartureDate: '2020-03-15',
          hasProGear: true,
          hasReceivedAdvance: true,
          hasRequestedAdvance: true,
          id: 'd733fe2f-b08d-434a-ad8d-551f4d597b03',
          proGearWeight: 1987,
          reviewedAt: '2022-07-02T14:20:14.636Z',
          shipmentId: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
          sitEstimatedCost: 123456,
          sitEstimatedDepartureDate: '2022-07-13',
          sitEstimatedEntryDate: '2022-07-05',
          sitEstimatedWeight: 1100,
          sitExpected: true,
          sitLocation: 'DESTINATION',
          spouseProGearWeight: 498,
          status: 'SUBMITTED',
          submittedAt: '2022-07-01T13:41:33.252Z',
          updatedAt: '2022-07-01T14:23:19.780Z',
        },
        primeEstimatedWeightRecordedDate: null,
        requestedPickupDate: null,
        requiredDeliveryDate: null,
        scheduledPickupDate: null,
        secondaryDeliveryAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        secondaryPickupAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        shipmentType: 'PPM',
        status: 'APPROVED',
        updatedAt: '2022-07-01T14:23:19.738Z',
        mtoServiceItems: [],
      },
    ],
    order: {
      customer: {
        branch: 'AIR_FORCE',
        dodID: '5917531070',
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4xNDIxNTZa',
        email: 'leo_spaceman_sm@example.com',
        firstName: 'Leo',
        id: 'e2de409b-edb9-42af-b50f-564458e08ada',
        lastName: 'Spacemen',
        phone: '555-555-5555',
        userID: 'ae204f8a-6222-45a1-9b79-e2d52441b4f2',
      },
      customerID: 'e2de409b-edb9-42af-b50f-564458e08ada',
      destinationDutyLocation: {
        address: {
          city: 'Augusta',
          country: 'United States',
          eTag: 'MjAyMS0xMC0xOFQxODoyMzoxMi4zMTQzNDZa',
          id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
          postalCode: '30813',
          state: 'GA',
          streetAddress1: 'Fort Gordon',
        },
        addressID: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
        id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
        name: 'Fort Gordon',
      },
      eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4yMzAxMVo=',
      entitlement: {
        authorizedWeight: 8000,
        dependentsAuthorized: true,
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4xNzc0MjZa',
        id: '46ee60c2-9b17-44c7-9202-15a84327fc2f',
        nonTemporaryStorage: true,
        organizationalClothingAndIndividualEquipment: true,
        privatelyOwnedVehicle: true,
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        requiredMedicalEquipmentWeight: 1000,
        storageInTransit: 2,
        totalDependents: 1,
        totalWeight: 5000,
      },
      id: '8cda4825-283c-4910-89f4-1741e2fd9cb7',
      linesOfAccounting: 'F8E1',
      orderNumber: 'ORDER3',
      originDutyLocation: {
        address: {
          city: 'Des Moines',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4yMDgyNjha',
          id: 'dbbee525-9c88-40c1-a549-6330b35972d2',
          postalCode: '50309',
          state: 'IA',
          streetAddress1: '987 Other Avenue',
          streetAddress2: 'P.O. Box 1234',
          streetAddress3: 'c/o Another Person',
        },
        addressID: 'dbbee525-9c88-40c1-a549-6330b35972d2',
        id: '0ecd8fb1-0551-44c8-a15e-83c5f4e3ae0f',
        name: 'XOXhgDSIRS',
      },
      reportByDate: '2018-08-01',
    },
    orderID: '8cda4825-283c-4910-89f4-1741e2fd9cb7',
    paymentRequests: [
      {
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41Nzc2OTha',
        id: '532ec513-8297-44b3-91a8-5167650b2869',
        isFinal: false,
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        paymentRequestNumber: '3301-9920-1',
        paymentServiceItems: [
          {
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS42Mzc5MzJa',
            id: '8fdf0b3a-c102-4084-84fe-22903f20470b',
            mtoServiceItemID: '8829fb28-69c1-45d7-98bc-c724478d5106',
            paymentRequestID: '532ec513-8297-44b3-91a8-5167650b2869',
            referenceID: '3301-9920-8fdf0b3a',
            status: 'REQUESTED',
          },
        ],
        status: 'PENDING',
      },
    ],
    ppmType: 'PARTIAL',
    referenceId: '3301-9920',
    updatedAt: '2021-10-18T18:24:41.362Z',
    mtoServiceItems: [
      {
        reServiceCode: 'STEST',
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41MzE0NjRa',
        id: '8829fb28-69c1-45d7-98bc-c724478d5106',
        modelType: 'MTOServiceItemBasic',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        mtoShipmentID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        reServiceName: 'Test Service',
        status: 'APPROVED',
      },
    ],
  },
};

const approvedPPMMoveTaskOrder = {
  moveTaskOrder: {
    availableToPrimeAt: '2021-10-18T18:24:41.235Z',
    createdAt: '2021-10-18T18:24:41.362Z',
    diversion: false,
    eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNjIxNjRa',
    excessWeightAcknowledgedAt: null,
    excessWeightQualifiedAt: null,
    excessWeightUploadId: null,
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    moveCode: 'LR4T8V',
    mtoShipments: [
      {
        actualPickupDate: null,
        approvedDate: null,
        counselorRemarks: 'These are counselor remarks for a PPM.',
        createdAt: '2022-07-01T13:41:33.261Z',
        destinationAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43MzgzODla',
        firstAvailableDeliveryDate: null,
        id: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        pickupAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        ppmShipment: {
          actualMoveDate: '2022-07-13',
          advanceAmountReceived: 598600,
          advanceAmountRequested: 598700,
          approvedAt: '2022-07-03T14:20:21.620Z',
          createdAt: '2022-06-30T13:41:33.265Z',
          eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43ODA1Mlo=',
          estimatedIncentive: 1000000,
          estimatedWeight: 4000,
          expectedDepartureDate: '2020-03-15',
          hasProGear: true,
          hasReceivedAdvance: true,
          hasRequestedAdvance: true,
          id: 'd733fe2f-b08d-434a-ad8d-551f4d597b03',
          proGearWeight: 1987,
          reviewedAt: '2022-07-02T14:20:14.636Z',
          shipmentId: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
          sitEstimatedCost: 123456,
          sitEstimatedDepartureDate: '2022-07-13',
          sitEstimatedEntryDate: '2022-07-05',
          sitEstimatedWeight: 1100,
          sitExpected: true,
          sitLocation: 'DESTINATION',
          spouseProGearWeight: 498,
          status: 'SUBMITTED',
          submittedAt: '2022-07-01T13:41:33.252Z',
          updatedAt: '2022-07-01T14:23:19.780Z',
          pickupAddress: {
            city: 'Beverly Hills',
            county: 'Los Angeles',
            country: 'US',
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNjc3Mjda',
            id: 'cf159eca-162c-4131-84a0-795e684416a6',
            postalCode: '90210',
            state: 'CA',
            streetAddress1: '123 Any Street',
            streetAddress2: 'P.O. Box 12345',
            streetAddress3: 'c/o Some Person',
          },
          destinationAddress: {
            city: 'Fairfield',
            county: 'Solano',
            country: 'US',
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzI3NDJa',
            id: 'bfe61147-5fd7-426e-b473-54ccf77bde35',
            postalCode: '94535',
            state: 'CO',
            streetAddress1: '987 Any Avenue',
            streetAddress2: 'P.O. Box 9876',
            streetAddress3: 'c/o Some Person',
          },
          secondaryDeliveryAddress: {
            city: null,
            postalCode: null,
            state: null,
            streetAddress1: null,
          },
          secondaryPickupAddress: {
            city: null,
            postalCode: null,
            state: null,
            streetAddress1: null,
          },
          tertiaryDeliveryAddress: {
            city: null,
            postalCode: null,
            state: null,
            streetAddress1: null,
          },
          tertiaryPickupAddress: {
            city: null,
            postalCode: null,
            state: null,
            streetAddress1: null,
          },
        },
        primeEstimatedWeightRecordedDate: null,
        requestedPickupDate: null,
        requiredDeliveryDate: null,
        scheduledPickupDate: null,
        shipmentType: 'PPM',
        status: 'APPROVED',
        updatedAt: '2022-07-01T14:23:19.738Z',
        mtoServiceItems: [],
      },
    ],
    order: {
      customer: {
        branch: 'AIR_FORCE',
        dodID: '5917531070',
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4xNDIxNTZa',
        email: 'leo_spaceman_sm@example.com',
        firstName: 'Leo',
        id: 'e2de409b-edb9-42af-b50f-564458e08ada',
        lastName: 'Spacemen',
        phone: '555-555-5555',
        userID: 'ae204f8a-6222-45a1-9b79-e2d52441b4f2',
      },
      customerID: 'e2de409b-edb9-42af-b50f-564458e08ada',
      destinationDutyLocation: {
        address: {
          city: 'Augusta',
          country: 'United States',
          eTag: 'MjAyMS0xMC0xOFQxODoyMzoxMi4zMTQzNDZa',
          id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
          postalCode: '30813',
          state: 'GA',
          streetAddress1: 'Fort Gordon',
        },
        addressID: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
        id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
        name: 'Fort Gordon',
      },
      eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4yMzAxMVo=',
      entitlement: {
        authorizedWeight: 8000,
        dependentsAuthorized: true,
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4xNzc0MjZa',
        id: '46ee60c2-9b17-44c7-9202-15a84327fc2f',
        nonTemporaryStorage: true,
        organizationalClothingAndIndividualEquipment: true,
        privatelyOwnedVehicle: true,
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        requiredMedicalEquipmentWeight: 1000,
        storageInTransit: 2,
        totalDependents: 1,
        totalWeight: 5000,
      },
      id: '8cda4825-283c-4910-89f4-1741e2fd9cb7',
      linesOfAccounting: 'F8E1',
      orderNumber: 'ORDER3',
      originDutyLocation: {
        address: {
          city: 'Des Moines',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4yMDgyNjha',
          id: 'dbbee525-9c88-40c1-a549-6330b35972d2',
          postalCode: '50309',
          state: 'IA',
          streetAddress1: '987 Other Avenue',
          streetAddress2: 'P.O. Box 1234',
          streetAddress3: 'c/o Another Person',
        },
        addressID: 'dbbee525-9c88-40c1-a549-6330b35972d2',
        id: '0ecd8fb1-0551-44c8-a15e-83c5f4e3ae0f',
        name: 'XOXhgDSIRS',
      },
      reportByDate: '2018-08-01',
    },
    orderID: '8cda4825-283c-4910-89f4-1741e2fd9cb7',
    paymentRequests: [
      {
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41Nzc2OTha',
        id: '532ec513-8297-44b3-91a8-5167650b2869',
        isFinal: false,
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        paymentRequestNumber: '3301-9920-1',
        paymentServiceItems: [
          {
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS42Mzc5MzJa',
            id: '8fdf0b3a-c102-4084-84fe-22903f20470b',
            mtoServiceItemID: '8829fb28-69c1-45d7-98bc-c724478d5106',
            paymentRequestID: '532ec513-8297-44b3-91a8-5167650b2869',
            referenceID: '3301-9920-8fdf0b3a',
            status: 'REQUESTED',
          },
        ],
        status: 'PENDING',
      },
    ],
    ppmType: 'PARTIAL',
    referenceId: '3301-9920',
    updatedAt: '2021-10-18T18:24:41.362Z',
    mtoServiceItems: [
      {
        reServiceCode: 'STEST',
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41MzE0NjRa',
        id: '8829fb28-69c1-45d7-98bc-c724478d5106',
        modelType: 'MTOServiceItemBasic',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        mtoShipmentID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        reServiceName: 'Test Service',
        status: 'APPROVED',
      },
    ],
  },
};

const approvedMoveTaskOrderWithOneHHG = {
  moveTaskOrder: {
    availableToPrimeAt: '2021-10-18T18:24:41.235Z',
    createdAt: '2021-10-18T18:24:41.362Z',
    diversion: false,
    eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNjIxNjRa',
    excessWeightAcknowledgedAt: null,
    excessWeightQualifiedAt: null,
    excessWeightUploadId: null,
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    moveCode: 'AR5K8V',
    mtoShipments: [
      {
        actualPickupDate: '2020-03-16',
        agents: [
          {
            agentType: 'RELEASING_AGENT',
            createdAt: '2021-10-18T18:24:41.521Z',
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41MjE4NzNa',
            email: 'test@test.email.com',
            firstName: 'Test',
            id: 'f2619e1b-7729-4b97-845d-6ae1ebe299f2',
            lastName: 'Agent',
            mtoShipmentID: 'ce01a5b7-9b44-4513-8a8d-edb60f2a4aee',
            phone: '202-555-9301',
            updatedAt: '2021-10-18T18:24:41.521Z',
          },
        ],
        approvedDate: '2021-10-18',
        createdAt: '2021-10-18T18:24:41.377Z',
        customerRemarks: 'Please treat gently',
        destinationAddress: {
          city: 'Fairfield',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzI3NDJa',
          id: 'bfe61147-5fd7-426e-b473-54ccf77bde35',
          postalCode: '94535',
          state: 'CA',
          streetAddress1: '987 Any Avenue',
          streetAddress2: 'P.O. Box 9876',
          streetAddress3: 'c/o Some Person',
        },
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzc5Nzha',
        firstAvailableDeliveryDate: null,
        id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        pickupAddress: {
          city: 'Beverly Hills',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNjc3Mjda',
          id: 'cf159eca-162c-4131-84a0-795e684416a6',
          postalCode: '90210',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
        },
        primeActualWeight: null,
        primeEstimatedWeight: null,
        actualProGearWeight: 1500,
        actualSpouseProGearWeight: 250,
        primeEstimatedWeightRecordedDate: null,
        requestedPickupDate: '2020-03-15',
        requiredDeliveryDate: null,
        scheduledPickupDate: '2020-03-16',
        secondaryDeliveryAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        secondaryPickupAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        shipmentType: 'HHG',
        status: 'APPROVED',
        updatedAt: '2021-10-18T18:24:41.377Z',
        mtoServiceItems: null,
      },
    ],
    order: {
      customer: {
        branch: 'AIR_FORCE',
        dodID: '5917531070',
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4xNDIxNTZa',
        email: 'leo_spaceman_sm@example.com',
        firstName: 'Leo',
        id: 'e2de409b-edb9-42af-b50f-564458e08ada',
        lastName: 'Spacemen',
        phone: '555-555-5555',
        userID: 'ae204f8a-6222-45a1-9b79-e2d52441b4f2',
      },
      customerID: 'e2de409b-edb9-42af-b50f-564458e08ada',
      destinationDutyLocation: {
        address: {
          city: 'Augusta',
          country: 'United States',
          eTag: 'MjAyMS0xMC0xOFQxODoyMzoxMi4zMTQzNDZa',
          id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
          postalCode: '30813',
          state: 'GA',
          streetAddress1: 'Fort Gordon',
        },
        addressID: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
        id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
        name: 'Fort Gordon',
      },
      eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4yMzAxMVo=',
      entitlement: {
        authorizedWeight: 8000,
        dependentsAuthorized: true,
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4xNzc0MjZa',
        id: '46ee60c2-9b17-44c7-9202-15a84327fc2f',
        nonTemporaryStorage: true,
        organizationalClothingAndIndividualEquipment: true,
        privatelyOwnedVehicle: true,
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        requiredMedicalEquipmentWeight: 1000,
        storageInTransit: 2,
        totalDependents: 1,
        totalWeight: 5000,
      },
      id: '8cda4825-283c-4910-89f4-1741e2fd9cb7',
      linesOfAccounting: 'F8E1',
      orderNumber: 'ORDER3',
      originDutyLocation: {
        address: {
          city: 'Des Moines',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4yMDgyNjha',
          id: 'dbbee525-9c88-40c1-a549-6330b35972d2',
          postalCode: '50309',
          state: 'IA',
          streetAddress1: '987 Other Avenue',
          streetAddress2: 'P.O. Box 1234',
          streetAddress3: 'c/o Another Person',
        },
        addressID: 'dbbee525-9c88-40c1-a549-6330b35972d2',
        id: '0ecd8fb1-0551-44c8-a15e-83c5f4e3ae0f',
        name: 'XOXhgDSIRS',
      },
      reportByDate: '2018-08-01',
    },
    orderID: '8cda4825-283c-4910-89f4-1741e2fd9cb7',
    paymentRequests: [
      {
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41Nzc2OTha',
        id: '532ec513-8297-44b3-91a8-5167650b2869',
        isFinal: false,
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        paymentRequestNumber: '3301-9920-1',
        paymentServiceItems: [
          {
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS42Mzc5MzJa',
            id: '8fdf0b3a-c102-4084-84fe-22903f20470b',
            mtoServiceItemID: '8829fb28-69c1-45d7-98bc-c724478d5106',
            paymentRequestID: '532ec513-8297-44b3-91a8-5167650b2869',
            referenceID: '3301-9920-8fdf0b3a',
            status: 'REQUESTED',
          },
        ],
        status: 'PENDING',
      },
    ],
    ppmType: 'PARTIAL',
    referenceId: '3301-9920',
    updatedAt: '2021-10-18T18:24:41.362Z',
    mtoServiceItems: [
      {
        reServiceCode: 'STEST',
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41MzE0NjRa',
        id: '8829fb28-69c1-45d7-98bc-c724478d5106',
        modelType: 'MTOServiceItemBasic',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        mtoShipmentID: 'ce01a5b7-9b44-4513-8a8d-edb60f2a4aee',
        reServiceName: 'Test Service',
        status: 'APPROVED',
      },
    ],
  },
};

const moveDetailsURL = generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID: moveId });

const mockedComponent = (
  <MockProviders path={primeSimulatorRoutes.UPDATE_SHIPMENT_PATH} params={routingParams}>
    <PrimeUIShipmentUpdate setFlashMessage={jest.fn()} />
  </MockProviders>
);

const mockedPPMComponent = (
  <MockProviders path={primeSimulatorRoutes.UPDATE_SHIPMENT_PATH} params={ppmRoutingParams}>
    <PrimeUIShipmentUpdate setFlashMessage={jest.fn()} />
  </MockProviders>
);

const ppmMockedComponent = (
  <MockProviders
    path={primeSimulatorRoutes.UPDATE_SHIPMENT_PATH}
    params={{
      ...routingParams,
      shipmentId: ppmShipmentId,
    }}
  >
    <PrimeUIShipmentUpdate setFlashMessage={jest.fn()} />
  </MockProviders>
);

const readyReturnValue = {
  ...approvedMoveTaskOrder,
  isLoading: false,
  isError: false,
  isSuccess: true,
};
const readyReturnValueWithOneHHG = {
  ...approvedMoveTaskOrderWithOneHHG,
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const readyPPMReturnValue = {
  ...approvedPPMMoveTaskOrder,
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const loadingReturnValue = {
  ...approvedMoveTaskOrder,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...approvedMoveTaskOrder,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('Update Shipment Page', () => {
  it('renders the page without errors', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);

    render(mockedComponent);

    expect(await screen.findByText('Shipment Dates')).toBeInTheDocument();
    expect(await screen.findByText('Shipment Weights')).toBeInTheDocument();
    expect(await screen.findByText('Pro Gear Weights')).toBeInTheDocument();
    expect(await screen.findByText('Shipment Addresses')).toBeInTheDocument();
  });

  it('renders the Loading Placeholder when the query is still loading', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(loadingReturnValue);

    render(mockedComponent);

    const h2 = await screen.findByRole('heading', { name: 'Loading, please wait...', level: 2 });
    expect(h2).toBeInTheDocument();
  });

  it('renders the Something Went Wrong component when the query errors', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(errorReturnValue);

    render(mockedComponent);

    const errorMessage = await screen.findByText(/Something went wrong./);
    expect(errorMessage).toBeInTheDocument();
  });

  it('navigates the user to the home page when the cancel button is clicked', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);

    render(mockedComponent);

    const cancel = screen.getByRole('button', { name: 'Cancel' });
    await userEvent.click(cancel);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});

describe('Displays the shipment information to update', () => {
  it('displays the shipment information', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);

    render(mockedComponent);

    const shipmentDatesHeader = screen.getByRole('heading', { name: 'Shipment Dates', level: 2 });
    expect(shipmentDatesHeader).toBeInTheDocument();
    const updateShipmentContainer = shipmentDatesHeader.parentElement;

    expect(
      await within(updateShipmentContainer).findByRole('heading', {
        name: 'Shipment Weights',
        level: 2,
      }),
    ).toBeInTheDocument();
    expect(
      await within(updateShipmentContainer).findByRole('heading', {
        name: 'Pro Gear Weights',
        level: 2,
      }),
    ).toBeInTheDocument();
    expect(
      within(updateShipmentContainer).getByRole('heading', {
        name: 'Shipment Addresses',
        level: 2,
      }),
    ).toBeInTheDocument();
  });
  /*
it('displays the submit button disabled', async () => {

usePrimeSimulatorGetMove.mockReturnValue(missingPrimeUpdates);

render(<PrimeUIShipmentUpdate />);

expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
expect(
  screen.getByText(
    'At least one basic service item or shipment service item is required to create a payment request',
  ),
).toBeInTheDocument();

  });
   */
  it('displays the submit button active', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);

    render(mockedComponent);

    expect(await screen.findByRole('button', { name: 'Save' })).toBeEnabled();
    expect(screen.getByText(/123 Any Street/)).toBeInTheDocument();
  });

  // Tested here rather than on the form because this component parses the data.
  // Tests on the form were missing a parsing issue since they were passed fully formed test data.
  it('displays all location lookup values', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyPPMReturnValue);

    render(mockedPPMComponent);

    expect(screen.getAllByText(/Beverly Hills/).length).toBe(2);
    expect(screen.getAllByText(/CA/).length).toBe(2);
    expect(screen.getAllByText(/Los Angeles/).length).toBe(2);
    expect(screen.getAllByText(/90210/).length).toBe(2);

    expect(screen.getAllByText(/Fairfield/).length).toBe(2);
    expect(screen.getAllByText(/CO/).length).toBe(2);
    expect(screen.getAllByText(/Solano/).length).toBe(2);
    expect(screen.getAllByText(/94535/).length).toBe(2);
  });
});

describe('successful submission of form', () => {
  it('calls history router back to move details', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);
    updatePrimeMTOShipmentV3.mockReturnValue({});

    render(mockedComponent);

    const actualPickupDateInput = await screen.findByLabelText('Actual pickup');
    await userEvent.clear(actualPickupDateInput);
    await userEvent.type(actualPickupDateInput, '20 Oct 2021');

    const actualWeightInput = screen.getByLabelText(/Actual weight/);
    await userEvent.type(actualWeightInput, '10000');

    const actualProGearWeightInput = screen.getByLabelText(/Actual pro gear weight/);
    await userEvent.type(actualProGearWeightInput, '2000');

    const actualSpouseProGearWeightInput = screen.getByLabelText(/Actual spouse pro gear weight/);
    await userEvent.type(actualSpouseProGearWeightInput, '500');

    const saveButton = await screen.getByRole('button', { name: 'Save' });

    await waitFor(() => {
      expect(saveButton).toBeEnabled();
    });
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  }, 50000);

  it('successful submission of form when updating a shipments actual weight but not estimated weight', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValueWithOneHHG);
    updatePrimeMTOShipmentV3.mockReturnValue({});

    render(mockedComponent);

    const actualPickupDateInput = await screen.findByLabelText('Actual pickup');
    await userEvent.clear(actualPickupDateInput);
    await userEvent.type(actualPickupDateInput, '20 Oct 2021');

    const actualWeightInput = screen.getByLabelText(/Actual weight/);
    await userEvent.type(actualWeightInput, '500');

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await waitFor(() => {
      expect(saveButton).toBeEnabled();
    });
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  });

  it('successful submission of form when updating a shipments estimated weight but not actual weight', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValueWithOneHHG);
    updatePrimeMTOShipmentV3.mockReturnValue({});

    render(mockedComponent);

    const actualPickupDateInput = await screen.findByLabelText('Actual pickup');
    await userEvent.clear(actualPickupDateInput);
    await userEvent.type(actualPickupDateInput, '20 Oct 2021');

    const estimatedWeightInput = screen.getByLabelText(/Estimated weight/);
    await userEvent.type(estimatedWeightInput, '500');

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await waitFor(() => {
      expect(saveButton).toBeEnabled();
    });
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});

describe('Date input validation', () => {
  const testCases = [
    { label: 'Scheduled pickup', validInput: '03 Jan 2025', invalidInput: 'X' },
    { label: 'Actual pickup', validInput: '24 Dec 2025', invalidInput: 'Invalid' },
    { label: 'Scheduled delivery', validInput: '15 Jun 2025', invalidInput: '2025-02-29' },
    { label: 'Actual delivery', validInput: '25 Dec 2025', invalidInput: '0' },
  ];

  beforeEach(async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);
  });

  afterEach(() => {
    cleanup();
  });

  testCases.forEach(({ label, validInput, invalidInput }) => {
    describe(`${label} input`, () => {
      it('displays an error when an invalid date is entered', async () => {
        render(mockedComponent);
        const input = screen.getByLabelText(label);
        await userEvent.clear(input);
        await userEvent.type(input, invalidInput);
        await userEvent.click(document.body);

        expect(screen.getByText('Invalid date. Must be in the format: DD MMM YYYY')).toBeInTheDocument();
      });

      it('removes an error when an input is cleared', async () => {
        render(mockedComponent);
        const input = screen.getByLabelText(label);
        await userEvent.clear(input);
        await userEvent.type(input, invalidInput);
        await userEvent.click(document.body);

        expect(screen.getByText('Invalid date. Must be in the format: DD MMM YYYY')).toBeInTheDocument();
        await userEvent.clear(input);
        expect(screen.queryByText('Invalid date. Must be in the format: DD MMM YYYY')).not.toBeInTheDocument();
      });

      it('does not display an error when a valid date is entered', async () => {
        render(mockedComponent);
        const input = screen.getByLabelText(label);
        await userEvent.clear(input);
        await userEvent.type(input, validInput);
        await userEvent.click(document.body);

        expect(screen.queryByText('Invalid date. Must be in the format: DD MMM YYYY')).not.toBeInTheDocument();
      });
    });
  });
});

describe('Update Shipment Page for PPM', () => {
  it('renders the page without errors', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);

    render(ppmMockedComponent);

    expect(await screen.findByText('Dates')).toBeInTheDocument();
    expect(await screen.findByText('Origin Info')).toBeInTheDocument();
    expect(await screen.findByText('Destination Info')).toBeInTheDocument();
    expect(await screen.findByText('Storage In Transit (SIT)')).toBeInTheDocument();
    expect(await screen.findByText('Weights')).toBeInTheDocument();
    expect(await screen.findByText('Remarks')).toBeInTheDocument();
  });

  it('test destination address street 1 is OPTIONAL', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(readyReturnValue);

    render(ppmMockedComponent);

    waitFor(async () => {
      // Start controlled test case to verify everything is working.
      let input = await document.querySelector('input[name="ppmShipment.pickupAddress.streetAddress1"]');
      expect(input).toBeInTheDocument();
      // enter required street 1 for pickup
      await userEvent.type(input, '123 Street');
      // clear
      await userEvent.clear(input);
      await userEvent.tab();
      // verify Required alert is displayed
      const requiredAlerts = screen.getByRole('alert');
      expect(requiredAlerts).toHaveTextContent('Required');
      // make valid again to clear alert
      await userEvent.type(input, '123 Street');

      // Verify destination address street 1 is OPTIONAL.
      input = await document.querySelector('input[name="ppmShipment.destinationAddress.streetAddress1"]');
      expect(input).toBeInTheDocument();
      // enter something
      await userEvent.type(input, '123 Street');
      // clear
      await userEvent.clear(input);
      await userEvent.tab();
      // verify no validation is displayed after clearing destination address street 1 because it's OPTIONAL
      expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    });
  });
});

describe('Error Handling', () => {
  it('Correctly displays a specific error message when an error response is returned', async () => {
    updatePrimeMTOShipmentV3.mockRejectedValue({ body: { title: 'Error', detail: 'The data entered no good.' } });
    render(mockedComponent);

    waitFor(async () => {
      await userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'HHG');

      const saveButton = await screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);
      expect(screen.getByText('Prime API: Error')).toBeInTheDocument();
      expect(screen.getByText('The data entered no good.')).toBeInTheDocument();
    });
  });
});
