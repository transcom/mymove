/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { v4 } from 'uuid';
import { mount } from 'enzyme';
import { act, waitFor } from '@testing-library/react';

import MoveHome from './MoveHome';

import { MockProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';
import { cancelMove, downloadPPMAOAPacket } from 'services/internalApi';
import { ORDERS_TYPE } from 'constants/orders';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('containers/FlashMessage/FlashMessage', () => {
  const MockFlash = () => <div>Flash message</div>;
  MockFlash.displayName = 'ConnectedFlashMessage';
  return MockFlash;
});

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('store/entities/actions', () => ({
  updateMTOShipments: jest.fn(),
  updateAllMoves: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  deleteMTOShipment: jest.fn(),
  getMTOShipmentsForMove: jest.fn(),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
  downloadPPMAOAPacket: jest.fn().mockImplementation(() => Promise.resolve()),
  cancelMove: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const props = {
  serviceMember: {
    id: v4(),
    current_location: {
      transportation_office: {
        name: 'Test Transportation Office Name',
        phone_lines: ['555-555-5555'],
      },
    },
  },
  showLoggedInUser: jest.fn(),
  createServiceMember: jest.fn(),
  getSignedCertification: jest.fn(),
  updateAllMoves: jest.fn(),
  mtoShipments: [],
  mtoShipment: {},
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  loadMTOShipments: jest.fn(),
  updateShipmentList: jest.fn(),
};

const defaultPropsNoOrders = {
  ...props,
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-16T15:55:20.639Z',
        eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
        id: '6dad799c-4567-4a7d-9419-1a686797768f',
        moveCode: '4H8VCD',
        orders: {},
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const defaultPropsOrdersWithUBAllowance = {
  ...props,
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-16T15:55:20.639Z',
        eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
        id: '6dad799c-4567-4a7d-9419-1a686797768f',
        moveCode: '4H8VCD',
        orders: {
          authorizedWeight: 11000,
          created_at: '2024-02-16T15:55:20.634Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
            ub_allowance: 2000,
          },
          grade: 'E_7',
          has_dependents: false,
          id: '667b1ca7-f904-43c4-8f2d-a2ea2375d7d3',
          issue_date: '2024-02-22',
          new_duty_location: {
            address: {
              city: 'Fort Knox',
              country: 'United States',
              id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
              postalCode: '40121',
              state: 'KY',
              streetAddress1: 'n/a',
            },
            address_id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
            affiliation: 'ARMY',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '866ac8f6-94f5-4fa0-b7d1-be7fcf9d51e9',
            name: 'Fort Knox, KY 40121',
            transportation_office: {
              address: {
                city: 'Fort Knox',
                country: 'United States',
                id: 'ca758d13-b3b7-48a5-93bd-64912f0e2434',
                postalCode: '40121',
                state: 'KY',
                streetAddress1: 'LRC 25 W. Chaffee Ave',
                streetAddress2: 'Bldg 1384, 2nd Floor',
              },
              created_at: '2018-05-28T14:27:36.193Z',
              gbloc: 'BGAC',
              id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
              name: 'PPPO Fort Knox - USA',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:36.193Z',
            },
            transportation_office_id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          originDutyLocationGbloc: 'HAFC',
          origin_duty_location: {
            address: {
              city: 'Tinker AFB',
              country: 'United States',
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              postalCode: '73145',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            affiliation: 'AIR_FORCE',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            name: 'Tinker AFB, OK 73145',
            transportation_office: {
              address: {
                city: 'Tinker AFB',
                country: 'United States',
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                postalCode: '73145',
                state: 'OK',
                streetAddress1: '7330 Century Blvd',
                streetAddress2: 'Bldg 469',
              },
              created_at: '2018-05-28T14:27:40.605Z',
              gbloc: 'HAFC',
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              name: 'PPPO Tinker AFB - USAF',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:40.605Z',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          report_by_date: '2024-02-29',
          service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-16T15:55:20.634Z',
          uploaded_orders: {
            id: '573a2d22-8edf-467c-90dc-3885de10e2d2',
            service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
            uploads: [
              {
                bytes: 84847,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:12:56.328Z',
                filename: 'myUpload.png',
                id: '99fab296-ad63-4e34-8724-a8b73e357480',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:12:56.328Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/99fab296-ad63-4e34-8724-a8b73e357480?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const defaultPropsOrdersWithUploads = {
  ...props,
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-16T15:55:20.639Z',
        eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
        id: '6dad799c-4567-4a7d-9419-1a686797768f',
        moveCode: '4H8VCD',
        orders: {
          authorizedWeight: 11000,
          created_at: '2024-02-16T15:55:20.634Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_7',
          has_dependents: false,
          id: '667b1ca7-f904-43c4-8f2d-a2ea2375d7d3',
          issue_date: '2024-02-22',
          new_duty_location: {
            address: {
              city: 'Fort Knox',
              country: 'United States',
              id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
              postalCode: '40121',
              state: 'KY',
              streetAddress1: 'n/a',
            },
            address_id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
            affiliation: 'ARMY',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '866ac8f6-94f5-4fa0-b7d1-be7fcf9d51e9',
            name: 'Fort Knox, KY 40121',
            transportation_office: {
              address: {
                city: 'Fort Knox',
                country: 'United States',
                id: 'ca758d13-b3b7-48a5-93bd-64912f0e2434',
                postalCode: '40121',
                state: 'KY',
                streetAddress1: 'LRC 25 W. Chaffee Ave',
                streetAddress2: 'Bldg 1384, 2nd Floor',
              },
              created_at: '2018-05-28T14:27:36.193Z',
              gbloc: 'BGAC',
              id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
              name: 'PPPO Fort Knox - USA',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:36.193Z',
            },
            transportation_office_id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          originDutyLocationGbloc: 'HAFC',
          origin_duty_location: {
            address: {
              city: 'Tinker AFB',
              country: 'United States',
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              postalCode: '73145',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            affiliation: 'AIR_FORCE',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            name: 'Tinker AFB, OK 73145',
            transportation_office: {
              address: {
                city: 'Tinker AFB',
                country: 'United States',
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                postalCode: '73145',
                state: 'OK',
                streetAddress1: '7330 Century Blvd',
                streetAddress2: 'Bldg 469',
              },
              created_at: '2018-05-28T14:27:40.605Z',
              gbloc: 'HAFC',
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              name: 'PPPO Tinker AFB - USAF',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:40.605Z',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          report_by_date: '2024-02-29',
          service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-16T15:55:20.634Z',
          uploaded_orders: {
            id: '573a2d22-8edf-467c-90dc-3885de10e2d2',
            service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
            uploads: [
              {
                bytes: 84847,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:12:56.328Z',
                filename: 'myUpload.png',
                id: '99fab296-ad63-4e34-8724-a8b73e357480',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:12:56.328Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/99fab296-ad63-4e34-8724-a8b73e357480?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const defaultPropsOrdersWithUnsubmittedShipments = {
  ...props,
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-16T15:55:20.639Z',
        eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
        id: '6dad799c-4567-4a7d-9419-1a686797768f',
        moveCode: '4H8VCD',
        mtoShipments: [
          {
            createdAt: '2024-02-20T17:21:05.318Z',
            customerRemarks: 'some remarks',
            destinationAddress: {
              city: 'Fort Sill',
              country: 'United States',
              id: '7787c25e-fe15-4e13-8e38-23397e5dbfb3',
              postalCode: '73503',
              state: 'OK',
              streetAddress1: 'N/A',
            },
            eTag: 'MjAyNC0wMi0yMFQxNzoyMTowNS4zMTgwODNa',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: 'be807bb2-572b-4677-9896-c7f670ac72fa',
            moveTaskOrderID: 'cf2508aa-2b0a-47e9-8688-37b41623837d',
            pickupAddress: {
              city: 'Oklahoma City',
              id: 'c8ef1288-1588-44ee-b9fb-b38c703d2ca5',
              postalCode: '74133',
              state: 'OK',
              streetAddress1: '1234 S Somewhere Street',
              streetAddress2: '',
            },
            requestedDeliveryDate: '2024-03-15',
            requestedPickupDate: '2024-02-29',
            shipmentType: 'HHG',
            status: 'SUBMITTED',
            updatedAt: '2024-02-20T17:21:05.318Z',
          },
          {
            createdAt: '2024-02-20T17:21:48.242Z',
            eTag: 'MjAyNC0wMi0yMFQxNzoyMjowMy4wMzk2Njla',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: '0c7f88b8-75a9-41fe-b884-ea39e6024f24',
            moveTaskOrderID: 'cf2508aa-2b0a-47e9-8688-37b41623837d',
            ppmShipment: {
              actualDestinationPostalCode: null,
              actualMoveDate: null,
              actualPickupPostalCode: null,
              advanceAmountReceived: null,
              advanceAmountRequested: null,
              approvedAt: null,
              createdAt: '2024-02-20T17:21:48.248Z',
              eTag: 'MjAyNC0wMi0yMFQxNzoyMjowMy4wODU5Mzda',
              estimatedIncentive: 339123,
              estimatedWeight: 2000,
              expectedDepartureDate: '2024-02-23',
              finalIncentive: null,
              hasProGear: false,
              hasReceivedAdvance: null,
              hasRequestedAdvance: false,
              id: '5f1f0b88-9cb9-4b48-a9ad-2af6c1113ca2',
              movingExpenses: [],
              proGearWeight: null,
              proGearWeightTickets: [],
              reviewedAt: null,
              shipmentId: '0c7f88b8-75a9-41fe-b884-ea39e6024f24',
              sitEstimatedCost: null,
              sitEstimatedDepartureDate: null,
              sitEstimatedEntryDate: null,
              sitEstimatedWeight: null,
              sitExpected: false,
              spouseProGearWeight: null,
              status: 'DRAFT',
              submittedAt: null,
              updatedAt: '2024-02-20T17:22:03.085Z',
              weightTickets: [],
              pickupAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Pickup Test City',
                state: 'NY',
                postalCode: '10001',
              },
              destinationAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Destination Test City',
                state: 'NY',
                postalCode: '11111',
              },
            },
            shipmentType: 'PPM',
            status: 'DRAFT',
            updatedAt: '2024-02-20T17:22:03.039Z',
          },
        ],
        orders: {
          authorizedWeight: 11000,
          created_at: '2024-02-16T15:55:20.634Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_7',
          has_dependents: false,
          id: '667b1ca7-f904-43c4-8f2d-a2ea2375d7d3',
          issue_date: '2024-02-22',
          new_duty_location: {
            address: {
              city: 'Fort Knox',
              country: 'United States',
              id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
              postalCode: '40121',
              state: 'KY',
              streetAddress1: 'n/a',
            },
            address_id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
            affiliation: 'ARMY',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '866ac8f6-94f5-4fa0-b7d1-be7fcf9d51e9',
            name: 'Fort Knox, KY 40121',
            transportation_office: {
              address: {
                city: 'Fort Knox',
                country: 'United States',
                id: 'ca758d13-b3b7-48a5-93bd-64912f0e2434',
                postalCode: '40121',
                state: 'KY',
                streetAddress1: 'LRC 25 W. Chaffee Ave',
                streetAddress2: 'Bldg 1384, 2nd Floor',
              },
              created_at: '2018-05-28T14:27:36.193Z',
              gbloc: 'BGAC',
              id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
              name: 'PPPO Fort Knox - USA',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:36.193Z',
            },
            transportation_office_id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          originDutyLocationGbloc: 'HAFC',
          origin_duty_location: {
            address: {
              city: 'Tinker AFB',
              country: 'United States',
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              postalCode: '73145',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            affiliation: 'AIR_FORCE',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            name: 'Tinker AFB, OK 73145',
            transportation_office: {
              address: {
                city: 'Tinker AFB',
                country: 'United States',
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                postalCode: '73145',
                state: 'OK',
                streetAddress1: '7330 Century Blvd',
                streetAddress2: 'Bldg 469',
              },
              created_at: '2018-05-28T14:27:40.605Z',
              gbloc: 'HAFC',
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              name: 'PPPO Tinker AFB - USAF',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:40.605Z',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          report_by_date: '2024-02-29',
          service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-16T15:55:20.634Z',
          uploaded_orders: {
            id: '573a2d22-8edf-467c-90dc-3885de10e2d2',
            service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
            uploads: [
              {
                bytes: 84847,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:12:56.328Z',
                filename: 'myUpload.png',
                id: '99fab296-ad63-4e34-8724-a8b73e357480',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:12:56.328Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/99fab296-ad63-4e34-8724-a8b73e357480?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const defaultPropsOrdersWithSubmittedShipments = {
  ...props,
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-16T15:55:20.639Z',
        eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
        id: '6dad799c-4567-4a7d-9419-1a686797768f',
        moveCode: '4H8VCD',
        mtoShipments: [
          {
            createdAt: '2024-02-20T17:21:05.318Z',
            customerRemarks: 'some remarks',
            destinationAddress: {
              city: 'Fort Sill',
              country: 'United States',
              id: '7787c25e-fe15-4e13-8e38-23397e5dbfb3',
              postalCode: '73503',
              state: 'OK',
              streetAddress1: 'N/A',
            },
            eTag: 'MjAyNC0wMi0yMFQxNzoyMTowNS4zMTgwODNa',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: 'be807bb2-572b-4677-9896-c7f670ac72fa',
            moveTaskOrderID: 'cf2508aa-2b0a-47e9-8688-37b41623837d',
            pickupAddress: {
              city: 'Oklahoma City',
              id: 'c8ef1288-1588-44ee-b9fb-b38c703d2ca5',
              postalCode: '74133',
              state: 'OK',
              streetAddress1: '1234 S Somewhere Street',
              streetAddress2: '',
            },
            requestedDeliveryDate: '2024-03-15',
            requestedPickupDate: '2024-02-29',
            shipmentType: 'HHG',
            status: 'SUBMITTED',
            updatedAt: '2024-02-20T17:21:05.318Z',
          },
          {
            createdAt: '2024-02-20T17:21:48.242Z',
            eTag: 'MjAyNC0wMi0yMFQxNzoyMjowMy4wMzk2Njla',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: '0c7f88b8-75a9-41fe-b884-ea39e6024f24',
            moveTaskOrderID: 'cf2508aa-2b0a-47e9-8688-37b41623837d',
            ppmShipment: {
              actualDestinationPostalCode: null,
              actualMoveDate: null,
              actualPickupPostalCode: null,
              advanceAmountReceived: null,
              advanceAmountRequested: null,
              approvedAt: null,
              createdAt: '2024-02-20T17:21:48.248Z',
              eTag: 'MjAyNC0wMi0yMFQxNzoyMjowMy4wODU5Mzda',
              estimatedIncentive: 339123,
              estimatedWeight: 2000,
              expectedDepartureDate: '2024-02-23',
              finalIncentive: null,
              hasProGear: false,
              hasReceivedAdvance: null,
              hasRequestedAdvance: false,
              id: '5f1f0b88-9cb9-4b48-a9ad-2af6c1113ca2',
              movingExpenses: [],
              proGearWeight: null,
              proGearWeightTickets: [],
              reviewedAt: null,
              shipmentId: '0c7f88b8-75a9-41fe-b884-ea39e6024f24',
              sitEstimatedCost: null,
              sitEstimatedDepartureDate: null,
              sitEstimatedEntryDate: null,
              sitEstimatedWeight: null,
              sitExpected: false,
              spouseProGearWeight: null,
              status: 'DRAFT',
              submittedAt: null,
              updatedAt: '2024-02-20T17:22:03.085Z',
              weightTickets: [],
              pickupAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Pickup Test City',
                state: 'NY',
                postalCode: '10001',
              },
              destinationAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Destination Test City',
                state: 'NY',
                postalCode: '11111',
              },
            },
            shipmentType: 'PPM',
            status: 'DRAFT',
            updatedAt: '2024-02-20T17:22:03.039Z',
          },
        ],
        orders: {
          authorizedWeight: 11000,
          created_at: '2024-02-16T15:55:20.634Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_7',
          has_dependents: false,
          id: '667b1ca7-f904-43c4-8f2d-a2ea2375d7d3',
          issue_date: '2024-02-22',
          new_duty_location: {
            address: {
              city: 'Fort Knox',
              country: 'United States',
              id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
              postalCode: '40121',
              state: 'KY',
              streetAddress1: 'n/a',
            },
            address_id: '31ed530d-4b59-42d7-9ea9-88ccc2978723',
            affiliation: 'ARMY',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '866ac8f6-94f5-4fa0-b7d1-be7fcf9d51e9',
            name: 'Fort Knox, KY 40121',
            transportation_office: {
              address: {
                city: 'Fort Knox',
                country: 'United States',
                id: 'ca758d13-b3b7-48a5-93bd-64912f0e2434',
                postalCode: '40121',
                state: 'KY',
                streetAddress1: 'LRC 25 W. Chaffee Ave',
                streetAddress2: 'Bldg 1384, 2nd Floor',
              },
              created_at: '2018-05-28T14:27:36.193Z',
              gbloc: 'BGAC',
              id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
              name: 'PPPO Fort Knox - USA',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:36.193Z',
            },
            transportation_office_id: '0357f830-2f32-41f3-9ca2-268dd70df5cb',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          originDutyLocationGbloc: 'HAFC',
          origin_duty_location: {
            address: {
              city: 'Tinker AFB',
              country: 'United States',
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              postalCode: '73145',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            affiliation: 'AIR_FORCE',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            name: 'Tinker AFB, OK 73145',
            transportation_office: {
              address: {
                city: 'Tinker AFB',
                country: 'United States',
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                postalCode: '73145',
                state: 'OK',
                streetAddress1: '7330 Century Blvd',
                streetAddress2: 'Bldg 469',
              },
              created_at: '2018-05-28T14:27:40.605Z',
              gbloc: 'HAFC',
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              name: 'PPPO Tinker AFB - USAF',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:40.605Z',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          report_by_date: '2024-02-29',
          service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-16T15:55:20.634Z',
          uploaded_orders: {
            id: '573a2d22-8edf-467c-90dc-3885de10e2d2',
            service_member_id: '856fec24-a70b-4860-9ba8-98d25676317e',
            uploads: [
              {
                bytes: 84847,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:12:56.328Z',
                filename: 'myUpload.png',
                id: '99fab296-ad63-4e34-8724-a8b73e357480',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:12:56.328Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/99fab296-ad63-4e34-8724-a8b73e357480?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'NEEDS_SERVICE_COUNSELING',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const defaultPropsAmendedOrdersWithAdvanceRequested = {
  ...props,
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-16T15:55:20.639Z',
        eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
        id: '6dad799c-4567-4a7d-9419-1a686797768f',
        moveCode: '4H8VCD',
        mtoShipments: [
          {
            createdAt: '2024-02-20T17:40:25.836Z',
            eTag: 'MjAyNC0wMi0yMFQxNzo0MDo0Ny43NzA5NjVa',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: '322ebc9f-0ca8-4943-a7a8-39235f4e680b',
            moveTaskOrderID: '4918b8c9-5e0a-4d65-a6b8-6a7a6ce265d4',
            ppmShipment: {
              actualDestinationPostalCode: null,
              actualMoveDate: null,
              actualPickupPostalCode: null,
              advanceAmountReceived: null,
              advanceAmountRequested: 400000,
              approvedAt: null,
              createdAt: '2024-02-20T17:40:25.842Z',
              eTag: 'MjAyNC0wMi0yMFQxNzo0MDo0Ny43NzI5MzNa',
              estimatedIncentive: 678255,
              estimatedWeight: 4000,
              expectedDepartureDate: '2024-02-24',
              finalIncentive: null,
              hasProGear: false,
              hasReceivedAdvance: null,
              hasRequestedAdvance: true,
              id: 'd18b865f-fd12-495d-91fa-65b53d72705a',
              movingExpenses: [],
              proGearWeight: null,
              proGearWeightTickets: [],
              reviewedAt: null,
              shipmentId: '322ebc9f-0ca8-4943-a7a8-39235f4e680b',
              sitEstimatedCost: null,
              sitEstimatedDepartureDate: null,
              sitEstimatedEntryDate: null,
              sitEstimatedWeight: null,
              sitExpected: false,
              spouseProGearWeight: null,
              status: 'SUBMITTED',
              submittedAt: null,
              updatedAt: '2024-02-20T17:40:47.772Z',
              weightTickets: [],
              pickupAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Pickup Test City',
                state: 'NY',
                postalCode: '10001',
              },
              destinationAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Destination Test City',
                state: 'NY',
                postalCode: '11111',
              },
            },
            shipmentType: 'PPM',
            status: 'SUBMITTED',
            updatedAt: '2024-02-20T17:40:47.770Z',
          },
        ],
        orders: {
          authorizedWeight: 11000,
          created_at: '2024-02-20T17:11:08.815Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_7',
          has_dependents: false,
          id: '9db91886-40eb-4910-9c87-968fecd44d4b',
          issue_date: '2024-02-22',
          new_duty_location: {
            address: {
              city: 'Fort Sill',
              country: 'United States',
              id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
              postalCode: '73503',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
            affiliation: 'ARMY',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '5c182566-0e6e-46f2-9eef-f07963783575',
            name: 'Fort Sill, OK 73503',
            transportation_office: {
              address: {
                city: 'Fort Sill',
                country: 'United States',
                id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
                postalCode: '73503',
                state: 'OK',
                streetAddress1: '4700 Mow Way Rd',
                streetAddress2: 'Room 110',
              },
              created_at: '2018-05-28T14:27:35.547Z',
              gbloc: 'JEAT',
              id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
              name: 'PPPO Fort Sill - USA',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:35.547Z',
            },
            transportation_office_id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          originDutyLocationGbloc: 'HAFC',
          origin_duty_location: {
            address: {
              city: 'Tinker AFB',
              country: 'United States',
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              postalCode: '73145',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            affiliation: 'AIR_FORCE',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            name: 'Tinker AFB, OK 73145',
            transportation_office: {
              address: {
                city: 'Tinker AFB',
                country: 'United States',
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                postalCode: '73145',
                state: 'OK',
                streetAddress1: '7330 Century Blvd',
                streetAddress2: 'Bldg 469',
              },
              created_at: '2018-05-28T14:27:40.605Z',
              gbloc: 'HAFC',
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              name: 'PPPO Tinker AFB - USAF',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:40.605Z',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          report_by_date: '2024-02-24',
          service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-20T17:40:58.221Z',
          uploaded_amended_orders: {
            id: '33c8773e-3409-457f-b94e-b8683514cbcd',
            service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
            uploads: [
              {
                bytes: 1578588,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:40:58.233Z',
                filename: 'Screenshot 2024-02-15 at 12.22.53 PM (2).png',
                id: 'f26f3427-a289-4faf-90da-2d02f3094a00',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:40:58.233Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/f26f3427-a289-4faf-90da-2d02f3094a00?contentType=image%2Fpng',
              },
            ],
          },
          uploaded_orders: {
            id: 'fa2e5695-8c95-4460-91d5-e7d29dafa0b0',
            service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
            uploads: [
              {
                bytes: 84847,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:12:56.328Z',
                filename: 'Screenshot 2024-02-12 at 8.26.20 AM.png',
                id: '99fab296-ad63-4e34-8724-a8b73e357480',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:12:56.328Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/99fab296-ad63-4e34-8724-a8b73e357480?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'NEEDS SERVICE COUNSELING',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const expectedPpmShipmentID = 'd18b865f-fd12-495d-91fa-65b53d72705a';

const defaultPropsWithAdvanceAndPPMApproved = {
  ...props,
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-16T15:55:20.639Z',
        eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
        id: '6dad799c-4567-4a7d-9419-1a686797768f',
        moveCode: '4H8VCD',
        mtoShipments: [
          {
            createdAt: '2024-02-20T17:40:25.836Z',
            eTag: 'MjAyNC0wMi0yMFQxODowMToxNC43NTY1MTJa',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: '322ebc9f-0ca8-4943-a7a8-39235f4e680b',
            moveTaskOrderID: '4918b8c9-5e0a-4d65-a6b8-6a7a6ce265d4',
            ppmShipment: {
              actualDestinationPostalCode: null,
              actualMoveDate: null,
              actualPickupPostalCode: null,
              advanceAmountReceived: null,
              advanceAmountRequested: 400000,
              advanceStatus: 'APPROVED',
              approvedAt: '2024-02-20T18:01:14.760Z',
              createdAt: '2024-02-20T17:40:25.842Z',
              eTag: 'MjAyNC0wMi0yMFQxODowMToxNC43NjAyNTha',
              estimatedIncentive: 678255,
              estimatedWeight: 4000,
              expectedDepartureDate: '2024-02-24',
              finalIncentive: null,
              hasProGear: false,
              hasReceivedAdvance: null,
              hasRequestedAdvance: true,
              id: expectedPpmShipmentID,
              movingExpenses: [],
              proGearWeight: null,
              proGearWeightTickets: [],
              reviewedAt: null,
              shipmentId: '322ebc9f-0ca8-4943-a7a8-39235f4e680b',
              sitEstimatedCost: null,
              sitEstimatedDepartureDate: null,
              sitEstimatedEntryDate: null,
              sitEstimatedWeight: null,
              sitExpected: false,
              spouseProGearWeight: null,
              status: 'WAITING_ON_CUSTOMER',
              submittedAt: null,
              updatedAt: '2024-02-20T18:01:14.760Z',
              weightTickets: [],
              pickupAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Pickup Test City',
                state: 'NY',
                postalCode: '10001',
              },
              destinationAddress: {
                streetAddress1: '1 Test Street',
                streetAddress2: '2 Test Street',
                streetAddress3: '3 Test Street',
                city: 'Destination Test City',
                state: 'NY',
                postalCode: '11111',
              },
            },
            shipmentType: 'PPM',
            status: 'APPROVED',
            updatedAt: '2024-02-20T18:01:14.756Z',
          },
        ],
        orders: {
          authorizedWeight: 11000,
          created_at: '2024-02-20T17:11:08.815Z',
          department_indicator: 'ARMY',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_7',
          has_dependents: false,
          id: '9db91886-40eb-4910-9c87-968fecd44d4b',
          issue_date: '2024-02-22',
          new_duty_location: {
            address: {
              city: 'Fort Sill',
              country: 'United States',
              id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
              postalCode: '73503',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
            affiliation: 'ARMY',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '5c182566-0e6e-46f2-9eef-f07963783575',
            name: 'Fort Sill, OK 73503',
            transportation_office: {
              address: {
                city: 'Fort Sill',
                country: 'United States',
                id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
                postalCode: '73503',
                state: 'OK',
                streetAddress1: '4700 Mow Way Rd',
                streetAddress2: 'Room 110',
              },
              created_at: '2018-05-28T14:27:35.547Z',
              gbloc: 'JEAT',
              id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
              name: 'PPPO Fort Sill - USA',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:35.547Z',
            },
            transportation_office_id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          orders_number: '12345678901234',
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          orders_type_detail: 'PCS_TDY',
          originDutyLocationGbloc: 'HAFC',
          origin_duty_location: {
            address: {
              city: 'Tinker AFB',
              country: 'United States',
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              postalCode: '73145',
              state: 'OK',
              streetAddress1: 'n/a',
            },
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            affiliation: 'AIR_FORCE',
            created_at: '2024-02-15T14:42:58.875Z',
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            name: 'Tinker AFB, OK 73145',
            transportation_office: {
              address: {
                city: 'Tinker AFB',
                country: 'United States',
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                postalCode: '73145',
                state: 'OK',
                streetAddress1: '7330 Century Blvd',
                streetAddress2: 'Bldg 469',
              },
              created_at: '2018-05-28T14:27:40.605Z',
              gbloc: 'HAFC',
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              name: 'PPPO Tinker AFB - USAF',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:40.605Z',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            updated_at: '2024-02-15T14:42:58.875Z',
          },
          report_by_date: '2024-02-24',
          service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          tac: '1111',
          updated_at: '2024-02-20T18:01:06.825Z',
          uploaded_amended_orders: {
            id: '33c8773e-3409-457f-b94e-b8683514cbcd',
            service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
            uploads: [
              {
                bytes: 1578588,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:40:58.233Z',
                filename: 'Screenshot 2024-02-15 at 12.22.53 PM (2).png',
                id: 'f26f3427-a289-4faf-90da-2d02f3094a00',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:40:58.233Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/f26f3427-a289-4faf-90da-2d02f3094a00?contentType=image%2Fpng',
              },
            ],
          },
          uploaded_orders: {
            id: 'fa2e5695-8c95-4460-91d5-e7d29dafa0b0',
            service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
            uploads: [
              {
                bytes: 84847,
                contentType: 'image/png',
                createdAt: '2024-02-20T17:12:56.328Z',
                filename: 'Screenshot 2024-02-12 at 8.26.20 AM.png',
                id: '99fab296-ad63-4e34-8724-a8b73e357480',
                status: 'PROCESSING',
                updatedAt: '2024-02-20T17:12:56.328Z',
                url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/99fab296-ad63-4e34-8724-a8b73e357480?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'APPROVED',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

const mountMoveHomeWithProviders = (defaultProps) => {
  const moveId = defaultProps.serviceMemberMoves.currentMove[0].id;
  return mount(
    <MockProviders path={customerRoutes.MOVE_HOME_PATH} params={{ moveId }}>
      <MoveHome {...defaultProps} />
    </MockProviders>,
  );
};

afterEach(() => {
  jest.resetAllMocks();
});

describe('Home component', () => {
  describe('with default props, renders the right allowances', () => {
    it('renders Home with the right amount of components', async () => {
      isBooleanFlagEnabled.mockResolvedValue(true);
      let wrapper;
      // wrapping rendering in act to ensure all state updates are complete
      await act(async () => {
        wrapper = mountMoveHomeWithProviders(defaultPropsOrdersWithUBAllowance);
      });
      await waitFor(() => {
        expect(wrapper.text()).toContain('Weight allowance');
        expect(wrapper.text()).toContain('11,000 lbs');
        expect(wrapper.text()).toContain('UB allowance');
        expect(wrapper.text()).toContain('2,000 lbs');
      });

      const ubToolTip = wrapper.find('ToolTip');
      expect(ubToolTip.exists()).toBe(true);

      ubToolTip.simulate('click');
      const toolTipText = 'The weight of your UB shipment is also part of your overall authorized weight allowance.';
      expect(ubToolTip.text()).toBe(toolTipText);
    });
  });

  describe('with default props, orders but no uploads', () => {
    const wrapper = mountMoveHomeWithProviders(defaultPropsNoOrders);

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(4);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('profile step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
    });

    it('has appropriate step headers for no orders', () => {
      expect(wrapper.text()).toContain('Next step: Add your orders');
      expect(wrapper.text()).toContain('Profile complete');
      expect(wrapper.text()).toContain('Make sure to keep your personal information up to date during your move.');
      expect(wrapper.text()).toContain('Upload orders');
      expect(wrapper.text()).toContain('Upload photos of each page, or upload a PDF.');
      expect(wrapper.text()).toContain('Set up shipments');
      expect(wrapper.text()).toContain(
        'We will collect addresses, dates, and how you want to move your personal property.',
      );
      expect(wrapper.text()).toContain(
        'Note: You can change these details later by talking to a move counselor or customer care representative.',
      );
      expect(wrapper.text()).toContain('Confirm move request');
      expect(wrapper.text()).toContain(
        'Review your move details and sign the legal paperwork, then send the info on to your move counselor.',
      );
    });

    it('has enabled and disabled buttons based on step', () => {
      // shipment step button should have a disabled button
      const shipmentStep = wrapper.find('Step[step="3"]');
      expect(shipmentStep.prop('actionBtnDisabled')).toBeTruthy();
      // confirm move request step should have a disabled button
      const confirmMoveRequest = wrapper.find('Step[step="4"]');
      expect(confirmMoveRequest.prop('actionBtnDisabled')).toBeTruthy();
    });
  });

  describe('with default props, orders with uploads', () => {
    const wrapper = mountMoveHomeWithProviders(defaultPropsOrdersWithUploads);

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(4);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('profile and order step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      const orderStep = wrapper.find('Step[step="2"]');
      expect(orderStep.prop('editBtnLabel')).toEqual('Edit');
    });

    it('has appropriate step headers for orders with uploads', () => {
      expect(wrapper.text()).toContain('Time for step 3: Set up your shipments');
      expect(wrapper.text()).toContain(
        "Share where and when you're moving, and how you want your things to be shipped.",
      );
      expect(wrapper.text()).toContain('Profile complete');
      expect(wrapper.text()).toContain('Orders uploaded');
      expect(wrapper.find('DocsUploaded').length).toBe(1);
      expect(wrapper.text()).toContain('Set up shipments');
      expect(wrapper.text()).toContain('Confirm move request');
    });

    it('has enabled and disabled buttons based on step', () => {
      // shipment step button should now be enabled
      const shipmentStep = wrapper.find('Step[step="3"]');
      expect(shipmentStep.prop('actionBtnDisabled')).toBeFalsy();
      // confirm move request step should still be disabled
      const confirmMoveRequest = wrapper.find('Step[step="4"]');
      expect(confirmMoveRequest.prop('actionBtnDisabled')).toBeTruthy();
    });
  });

  describe('with default props, orders and unsubmitted HHG & PPM shipments', () => {
    const wrapper = mountMoveHomeWithProviders(defaultPropsOrdersWithUnsubmittedShipments);

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(4);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('profile and order step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      const orderStep = wrapper.find('Step[step="2"]');
      expect(orderStep.prop('editBtnLabel')).toEqual('Edit');
    });

    it('has appropriate step headers for orders with shipments', () => {
      expect(wrapper.text()).toContain('Time to submit your move');
      expect(wrapper.text()).toContain('Profile complete');
      expect(wrapper.text()).toContain('Orders uploaded');
      expect(wrapper.find('DocsUploaded').length).toBe(1);
      expect(wrapper.text()).toContain('Shipments');
      expect(wrapper.find('ShipmentList').length).toBe(1);
      expect(wrapper.text()).toContain('Confirm move request');
    });

    it('has enabled and disabled buttons based on step', () => {
      // shipment step button should now be "Add another shipment"
      const shipmentStep = wrapper.find('Step[step="3"]');
      expect(shipmentStep.prop('actionBtnLabel')).toBe('Add another shipment');
      // confirm move request step should now be enabled
      const confirmMoveRequest = wrapper.find('Step[step="4"]');
      expect(confirmMoveRequest.prop('actionBtnDisabled')).toBeFalsy();
    });

    it('cancel move button is visible', async () => {
      const cancelMoveButtonId = `button[data-testid="cancel-move-button"]`;
      expect(wrapper.find(cancelMoveButtonId).length).toBe(1);

      const mockResponse = {
        status: 'CANCELED',
      };
      cancelMove.mockImplementation(() => Promise.resolve(mockResponse));
      await wrapper.find(cancelMoveButtonId).simulate('click');
      await waitFor(() => {
        wrapper.find(`button[data-testid="modalSubmitButton"]`).simulate('click');
        expect(cancelMove).toHaveBeenCalledTimes(1);
      });
    });
  });

  describe('with default props, orders with HHG & PPM shipments and NEEDS_SERVICE_COUNSELING move status', () => {
    const wrapper = mountMoveHomeWithProviders(defaultPropsOrdersWithSubmittedShipments);

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(5);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('profile and order step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      const orderStep = wrapper.find('Step[step="2"]');
      expect(orderStep.prop('editBtnLabel')).toEqual('Upload/Manage Orders Documentation');
    });

    it('has appropriate step headers for orders with shipments', () => {
      expect(wrapper.text()).toContain('Next step: Your move gets approved');
      expect(wrapper.text()).toContain('Profile complete');
      expect(wrapper.text()).toContain('Orders');
      expect(wrapper.text()).toContain('If you receive amended orders');
      expect(wrapper.text()).toContain('Shipments');
      expect(wrapper.find('ShipmentList').length).toBe(1);
      expect(wrapper.text()).toContain(
        'If you need to change, add, or cancel shipments, talk to your move counselor or Customer Care Representative',
      );
      expect(wrapper.text()).toContain('Move request confirmed');
      expect(wrapper.text()).toContain('Manage your PPM');
    });

    it('has enabled and disabled buttons based on step', () => {
      // confirm move request step should now be enabled
      const confirmMoveRequest = wrapper.find('Step[step="4"]');
      expect(confirmMoveRequest.prop('actionBtnDisabled')).toBeFalsy();
      expect(confirmMoveRequest.prop('actionBtnLabel')).toBe('Review your request');
    });

    it('cancel move button is not visible', () => {
      const cancelMoveButton = wrapper.find('button[data-testid="cancel-move-button"]');
      expect(cancelMoveButton.length).toBe(0);
    });
  });

  describe('with default props, with amended orders and advance requested', () => {
    const wrapper = mountMoveHomeWithProviders(defaultPropsAmendedOrdersWithAdvanceRequested);

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(6);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('profile and order step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      const orderStep = wrapper.find('Step[step="2"]');
      expect(orderStep.prop('editBtnLabel')).toEqual('Upload/Manage Orders Documentation');
    });

    it('has appropriate step headers for orders with shipments', () => {
      expect(wrapper.text()).toContain(
        'The transportation office will review your new documents and update your move info. Contact your movers to coordinate any changes to your move.',
      );
      expect(wrapper.text()).toContain('Next step: Contact your movers (if you have them)');
      expect(wrapper.text()).toContain('Profile complete');
      expect(wrapper.text()).toContain('Orders');
      expect(wrapper.text()).toContain('If you receive amended orders');
      expect(wrapper.text()).toContain('Shipments');
      expect(wrapper.find('ShipmentList').length).toBe(1);
      expect(wrapper.text()).toContain(
        'If you need to change, add, or cancel shipments, talk to your move counselor or Customer Care Representative',
      );
      expect(wrapper.text()).toContain('Move request confirmed');
      expect(wrapper.text()).toContain('Advance request submitted');
      expect(wrapper.text()).toContain('Manage your PPM');
      expect(wrapper.find('PPMSummaryList').length).toBe(1);
    });

    it('has enabled and disabled buttons based on step', () => {
      // confirm move request step should now be enabled
      const confirmMoveRequest = wrapper.find('Step[step="4"]');
      expect(confirmMoveRequest.prop('actionBtnDisabled')).toBeFalsy();
      expect(confirmMoveRequest.prop('actionBtnLabel')).toBe('Review your request');
    });
  });

  describe('with default props, with approved PPM and advance', () => {
    const wrapper = mountMoveHomeWithProviders(defaultPropsWithAdvanceAndPPMApproved);

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(6);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('profile and order step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      const orderStep = wrapper.find('Step[step="2"]');
      expect(orderStep.prop('editBtnLabel')).toEqual('Upload/Manage Orders Documentation');
      const advanceStep = wrapper.find('Step[step="5"]');
      expect(advanceStep.prop('completedHeaderText')).toEqual('Advance request reviewed');
    });

    it('has appropriate step headers for orders with shipments', () => {
      expect(wrapper.text()).toContain('Your move is in progress.');
      expect(wrapper.text()).toContain('Profile complete');
      expect(wrapper.text()).toContain('Orders');
      expect(wrapper.text()).toContain('If you receive amended orders');
      expect(wrapper.text()).toContain('Shipments');
      expect(wrapper.find('ShipmentList').length).toBe(1);
      expect(wrapper.text()).toContain(
        'If you need to change, add, or cancel shipments, talk to your move counselor or Customer Care Representative',
      );
      expect(wrapper.text()).toContain('Move request confirmed');
      expect(wrapper.text()).toContain(
        'Your Advance Operating Allowance (AOA) request has been reviewed. Download the paperwork for approved requests and submit it to your Finance Office to receive your advance.',
      );
      expect(wrapper.text()).toContain('Manage your PPM');
      expect(wrapper.find('PPMSummaryList').length).toBe(1);
    });

    it('has enabled and disabled buttons based on step', () => {
      // confirm move request step should be enabled
      const confirmMoveRequest = wrapper.find('Step[step="4"]');
      expect(confirmMoveRequest.prop('actionBtnDisabled')).toBeFalsy();
      expect(confirmMoveRequest.prop('actionBtnLabel')).toBe('Review your request');
    });

    it('Download AOA Paperwork - success', async () => {
      const buttonId = `button[data-testid="asyncPacketDownloadLink${expectedPpmShipmentID}"]`;
      expect(wrapper.find(buttonId).length).toBe(1);
      const mockResponse = {
        ok: true,
        headers: {
          'content-disposition': 'filename="test.pdf"',
        },
        status: 200,
        data: null,
      };
      downloadPPMAOAPacket.mockImplementation(() => Promise.resolve(mockResponse));
      await wrapper.find(buttonId).simulate('click');
      await waitFor(() => {
        expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
      });
    });

    it('Download AOA Paperwork - failure', async () => {
      const buttonId = `button[data-testid="asyncPacketDownloadLink${expectedPpmShipmentID}"]`;
      expect(wrapper.find(buttonId).length).toBe(1);
      downloadPPMAOAPacket.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });
      await wrapper.find(buttonId).simulate('click');
      await waitFor(() => {
        // scrape text from error modal
        expect(wrapper.text()).toContain('Something went wrong downloading PPM paperwork.');
        expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
      });
    });
  });

  const defaultPropsWithEditedAdvanceAndPPMApproved = {
    ...props,
    serviceMemberMoves: {
      currentMove: [
        {
          createdAt: '2024-02-16T15:55:20.639Z',
          eTag: 'MjAyNC0wMi0xNlQxNTo1NToyMC42Mzk5MDRa',
          id: '6dad799c-4567-4a7d-9419-1a686797768f',
          moveCode: '4H8VCD',
          mtoShipments: [
            {
              createdAt: '2024-02-20T17:40:25.836Z',
              eTag: 'MjAyNC0wMi0yMFQxODowMToxNC43NTY1MTJa',
              hasSecondaryDeliveryAddress: false,
              hasSecondaryPickupAddress: false,
              id: '322ebc9f-0ca8-4943-a7a8-39235f4e680b',
              moveTaskOrderID: '4918b8c9-5e0a-4d65-a6b8-6a7a6ce265d4',
              ppmShipment: {
                actualDestinationPostalCode: null,
                actualMoveDate: null,
                actualPickupPostalCode: null,
                advanceAmountReceived: null,
                advanceAmountRequested: 400000,
                advanceStatus: 'APPROVED',
                approvedAt: '2024-02-20T18:01:14.760Z',
                createdAt: '2024-02-20T17:40:25.842Z',
                eTag: 'MjAyNC0wMi0yMFQxODowMToxNC43NjAyNTha',
                estimatedIncentive: 678255,
                estimatedWeight: 4000,
                expectedDepartureDate: '2024-02-24',
                finalIncentive: null,
                hasProGear: false,
                hasReceivedAdvance: null,
                hasRequestedAdvance: true,
                id: expectedPpmShipmentID,
                movingExpenses: [],
                proGearWeight: null,
                proGearWeightTickets: [],
                reviewedAt: null,
                shipmentId: '322ebc9f-0ca8-4943-a7a8-39235f4e680b',
                sitEstimatedCost: null,
                sitEstimatedDepartureDate: null,
                sitEstimatedEntryDate: null,
                sitEstimatedWeight: null,
                sitExpected: false,
                spouseProGearWeight: null,
                status: 'WAITING_ON_CUSTOMER',
                submittedAt: null,
                updatedAt: '2024-02-20T18:01:14.760Z',
                weightTickets: [],
                pickupAddress: {
                  streetAddress1: '1 Test Street',
                  streetAddress2: '2 Test Street',
                  streetAddress3: '3 Test Street',
                  city: 'Pickup Test City',
                  state: 'NY',
                  postalCode: '10001',
                },
                destinationAddress: {
                  streetAddress1: '1 Test Street',
                  streetAddress2: '2 Test Street',
                  streetAddress3: '3 Test Street',
                  city: 'Destination Test City',
                  state: 'NY',
                  postalCode: '11111',
                },
              },
              shipmentType: 'PPM',
              status: 'APPROVED',
              updatedAt: '2024-02-20T18:01:14.756Z',
            },
          ],
          orders: {
            authorizedWeight: 11000,
            created_at: '2024-02-20T17:11:08.815Z',
            department_indicator: 'ARMY',
            entitlement: {
              proGear: 2000,
              proGearSpouse: 500,
            },
            grade: 'E_7',
            has_dependents: false,
            id: '9db91886-40eb-4910-9c87-968fecd44d4b',
            issue_date: '2024-02-22',
            new_duty_location: {
              address: {
                city: 'Fort Sill',
                country: 'United States',
                id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
                postalCode: '73503',
                state: 'OK',
                streetAddress1: 'n/a',
              },
              address_id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
              affiliation: 'ARMY',
              created_at: '2024-02-15T14:42:58.875Z',
              id: '5c182566-0e6e-46f2-9eef-f07963783575',
              name: 'Fort Sill, OK 73503',
              transportation_office: {
                address: {
                  city: 'Fort Sill',
                  country: 'United States',
                  id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
                  postalCode: '73503',
                  state: 'OK',
                  streetAddress1: '4700 Mow Way Rd',
                  streetAddress2: 'Room 110',
                },
                created_at: '2018-05-28T14:27:35.547Z',
                gbloc: 'JEAT',
                id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
                name: 'PPPO Fort Sill - USA',
                phone_lines: [],
                updated_at: '2018-05-28T14:27:35.547Z',
              },
              transportation_office_id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
              updated_at: '2024-02-15T14:42:58.875Z',
            },
            orders_number: '12345678901234',
            orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
            orders_type_detail: 'PCS_TDY',
            originDutyLocationGbloc: 'HAFC',
            origin_duty_location: {
              address: {
                city: 'Tinker AFB',
                country: 'United States',
                id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
                postalCode: '73145',
                state: 'OK',
                streetAddress1: 'n/a',
              },
              address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              affiliation: 'AIR_FORCE',
              created_at: '2024-02-15T14:42:58.875Z',
              id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
              name: 'Tinker AFB, OK 73145',
              transportation_office: {
                address: {
                  city: 'Tinker AFB',
                  country: 'United States',
                  id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                  postalCode: '73145',
                  state: 'OK',
                  streetAddress1: '7330 Century Blvd',
                  streetAddress2: 'Bldg 469',
                },
                created_at: '2018-05-28T14:27:40.605Z',
                gbloc: 'HAFC',
                id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
                name: 'PPPO Tinker AFB - USAF',
                phone_lines: [],
                updated_at: '2018-05-28T14:27:40.605Z',
              },
              transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              updated_at: '2024-02-15T14:42:58.875Z',
            },
            report_by_date: '2024-02-24',
            service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
            spouse_has_pro_gear: false,
            status: 'DRAFT',
            tac: '1111',
            updated_at: '2024-02-20T18:01:06.825Z',
            uploaded_amended_orders: {
              id: '33c8773e-3409-457f-b94e-b8683514cbcd',
              service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
              uploads: [
                {
                  bytes: 1578588,
                  contentType: 'image/png',
                  createdAt: '2024-02-20T17:40:58.233Z',
                  filename: 'Screenshot 2024-02-15 at 12.22.53 PM (2).png',
                  id: 'f26f3427-a289-4faf-90da-2d02f3094a00',
                  status: 'PROCESSING',
                  updatedAt: '2024-02-20T17:40:58.233Z',
                  url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/f26f3427-a289-4faf-90da-2d02f3094a00?contentType=image%2Fpng',
                },
              ],
            },
            uploaded_orders: {
              id: 'fa2e5695-8c95-4460-91d5-e7d29dafa0b0',
              service_member_id: 'd6d26f51-a8f2-4294-aba4-2f38a759afe2',
              uploads: [
                {
                  bytes: 84847,
                  contentType: 'image/png',
                  createdAt: '2024-02-20T17:12:56.328Z',
                  filename: 'Screenshot 2024-02-12 at 8.26.20 AM.png',
                  id: '99fab296-ad63-4e34-8724-a8b73e357480',
                  status: 'PROCESSING',
                  updatedAt: '2024-02-20T17:12:56.328Z',
                  url: '/storage/user/9e16e5d7-4548-4f70-8a2a-b87d34ab3067/uploads/99fab296-ad63-4e34-8724-a8b73e357480?contentType=image%2Fpng',
                },
              ],
            },
          },
          status: 'APPROVED',
          submittedAt: '0001-01-01T00:00:00.000Z',
          updatedAt: '0001-01-01T00:00:00.000Z',
        },
      ],
      previousMoves: [],
    },
    uploadedOrderDocuments: [],
    uploadedAmendedOrderDocuments: [],
  };

  describe('with default props, with approved PPM and edited advance', () => {
    const wrapper = mountMoveHomeWithProviders(defaultPropsWithEditedAdvanceAndPPMApproved);

    it('renders Home with the right amount of components', () => {
      expect(wrapper.find('ConnectedFlashMessage').length).toBe(1);
      expect(wrapper.find('Step').length).toBe(6);
      expect(wrapper.find('Helper').length).toBe(1);
      expect(wrapper.find('Contact').length).toBe(1);
    });

    it('profile and order step is editable', () => {
      const profileStep = wrapper.find('Step[step="1"]');
      expect(profileStep.prop('editBtnLabel')).toEqual('Edit');
      const orderStep = wrapper.find('Step[step="2"]');
      expect(orderStep.prop('editBtnLabel')).toEqual('Upload/Manage Orders Documentation');
      const advanceStep = wrapper.find('Step[step="5"]');
      expect(advanceStep.prop('completedHeaderText')).toEqual('Advance request reviewed');
    });

    it('has appropriate step headers for orders with shipments', () => {
      expect(wrapper.text()).toContain('Your move is in progress.');
      expect(wrapper.text()).toContain('Profile complete');
      expect(wrapper.text()).toContain('Orders');
      expect(wrapper.text()).toContain('If you receive amended orders');
      expect(wrapper.text()).toContain('Shipments');
      expect(wrapper.find('ShipmentList').length).toBe(1);
      expect(wrapper.text()).toContain(
        'If you need to change, add, or cancel shipments, talk to your move counselor or Customer Care Representative',
      );
      expect(wrapper.text()).toContain('Move request confirmed');
      expect(wrapper.text()).toContain(
        'Your Advance Operating Allowance (AOA) request has been reviewed. Download the paperwork for approved requests and submit it to your Finance Office to receive your advance.',
      );
      expect(wrapper.text()).toContain('Manage your PPM');
      expect(wrapper.find('PPMSummaryList').length).toBe(1);
    });

    it('has enabled and disabled buttons based on step', () => {
      // confirm move request step should be enabled
      const confirmMoveRequest = wrapper.find('Step[step="4"]');
      expect(confirmMoveRequest.prop('actionBtnDisabled')).toBeFalsy();
      expect(confirmMoveRequest.prop('actionBtnLabel')).toBe('Review your request');
    });

    it('Download AOA Paperwork - success', async () => {
      const buttonId = `button[data-testid="asyncPacketDownloadLink${expectedPpmShipmentID}"]`;
      expect(wrapper.find(buttonId).length).toBe(1);
      const mockResponse = {
        ok: true,
        headers: {
          'content-disposition': 'filename="test.pdf"',
        },
        status: 200,
        data: null,
      };
      downloadPPMAOAPacket.mockImplementation(() => Promise.resolve(mockResponse));
      await wrapper.find(buttonId).simulate('click');
      await waitFor(() => {
        expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
      });
    });

    it('Download AOA Paperwork - failure', async () => {
      const buttonId = `button[data-testid="asyncPacketDownloadLink${expectedPpmShipmentID}"]`;
      expect(wrapper.find(buttonId).length).toBe(1);
      downloadPPMAOAPacket.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });
      await wrapper.find(buttonId).simulate('click');
      await waitFor(() => {
        // scrape text from error modal
        expect(wrapper.text()).toContain('Something went wrong downloading PPM paperwork.');
        expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
      });
    });
  });
});
