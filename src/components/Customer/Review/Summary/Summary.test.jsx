import React from 'react';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { FEATURE_FLAG_KEYS, MOVE_STATUSES } from 'shared/constants';
import { Summary } from 'components/Customer/Review/Summary/Summary';
import { renderWithRouterProp } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { customerRoutes } from 'constants/routes';
import { selectCurrentMoveFromAllMoves } from 'store/entities/selectors';
import { ORDERS_TYPE } from 'constants/orders';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectServiceMemberFromLoggedInUser: jest.fn(),
  selectCurrentMoveFromAllMoves: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const testMove = {
  createdAt: '2024-02-27T19:16:32.850Z',
  eTag: 'MjAyNC0wMi0yN1QxOToxNjozMi44NTAwNTda',
  id: '123',
  moveCode: 'WWYFP6',
  mtoShipments: [
    {
      createdAt: '2024-02-27T19:27:39.150Z',
      customerRemarks: '',
      destinationAddress: {
        city: 'Flagstaff',
        country: 'United States',
        id: '112e0d7b-90eb-44c4-80b1-5c1214fca0a7',
        postalCode: '86003',
        state: 'AZ',
        streetAddress1: 'N/A',
      },
      eTag: 'MjAyNC0wMi0yN1QxOToyNzozOS4xNTA3MjRa',
      hasSecondaryDeliveryAddress: false,
      hasSecondaryPickupAddress: false,
      id: 'f0082986-8e2f-443b-8411-191b3796e572',
      moveTaskOrderID: 'e23d629e-2a73-4b42-886b-fa60cb3db957',
      pickupAddress: {
        city: 'Tulsa',
        id: 'dac5e64d-1a69-4e83-a154-5fca04384544',
        postalCode: '74133',
        state: 'OK',
        streetAddress1: '8711 S 76th E Ave',
        streetAddress2: '',
      },
      requestedDeliveryDate: '2024-02-29',
      requestedPickupDate: '2024-03-01',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2024-02-27T19:27:39.150Z',
    },
  ],
  orders: {
    authorizedWeight: 8000,
    created_at: '2024-02-27T19:16:32.845Z',
    entitlement: {
      proGear: 2000,
      proGearSpouse: 500,
    },
    grade: 'E_6',
    has_dependents: false,
    id: '786e60ec-1bbd-48dd-bc12-b6ffcf212c54',
    issue_date: '2024-02-29',
    new_duty_location: {
      address: {
        city: 'Flagstaff',
        country: 'United States',
        id: 'cd51f4db-6195-473a-86cd-c3e5e07640e4',
        postalCode: '86003',
        state: 'AZ',
        streetAddress1: 'n/a',
      },
      address_id: 'cd51f4db-6195-473a-86cd-c3e5e07640e4',
      affiliation: null,
      created_at: '2024-02-27T18:22:12.471Z',
      id: '6ea57f62-2995-4b0c-a0a8-f11782cc8a3b',
      name: 'Flagstaff, AZ 86003',
      updated_at: '2024-02-27T18:22:12.471Z',
    },
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    originDutyLocationGbloc: 'BGAC',
    origin_duty_location: {
      address: {
        city: 'Aberdeen Proving Ground',
        country: 'United States',
        id: 'b6ca003e-1528-4e7c-a43e-830222ca7fb3',
        postalCode: '21005',
        state: 'MD',
        streetAddress1: 'n/a',
      },
      address_id: 'b6ca003e-1528-4e7c-a43e-830222ca7fb3',
      affiliation: 'ARMY',
      created_at: '2024-02-27T18:22:12.471Z',
      id: '61e092c4-575c-458a-9c3f-b93ad373c454',
      name: 'Aberdeen Proving Ground, MD 21005',
      transportation_office: {
        address: {
          city: 'Aberdeen Proving Ground',
          country: 'United States',
          id: 'ac4dbfa5-3068-4f8f-99d1-3cd850412bb9',
          postalCode: '21005',
          state: 'MD',
          streetAddress1: '4305 Susquehanna Ave',
          streetAddress2: 'Room 307',
        },
        created_at: '2018-05-28T14:27:41.772Z',
        gbloc: 'BGAC',
        id: '6a27dfbd-2a49-485f-86dd-49475d5facef',
        name: 'PPPO Aberdeen Proving Ground - USA',
        phone_lines: [],
        updated_at: '2018-05-28T14:27:41.772Z',
      },
      transportation_office_id: '6a27dfbd-2a49-485f-86dd-49475d5facef',
      updated_at: '2024-02-27T18:22:12.471Z',
    },
    report_by_date: '2024-02-29',
    service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
    spouse_has_pro_gear: false,
    status: 'DRAFT',
    updated_at: '2024-02-27T19:16:32.845Z',
    uploaded_orders: {
      id: 'b392f96f-20a0-43d2-9ca3-643cfd3b4182',
      service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
      uploads: [
        {
          bytes: 1137126,
          contentType: 'image/png',
          createdAt: '2024-02-27T19:16:38.998Z',
          filename: 'Screenshot 2024-02-15 at 12.22.53 PM.png',
          id: 'bc6c0e2d-fbef-4c32-8487-92c14b613040',
          status: 'PROCESSING',
          updatedAt: '2024-02-27T19:16:38.998Z',
          url: '/storage/user/f94c8fa6-89de-4594-be6a-ca6f1b4e22d0/uploads/bc6c0e2d-fbef-4c32-8487-92c14b613040?contentType=image%2Fpng',
        },
      ],
    },
  },
  status: 'DRAFT',
  submittedAt: '0001-01-01T00:00:00.000Z',
  updatedAt: '0001-01-01T00:00:00.000Z',
};

const testProps = {
  serviceMember: {
    id: '666',
    current_location: {
      name: 'Test Duty Location',
    },
    residential_address: {
      city: 'New York',
      postalCode: '10001',
      state: 'NY',
      streetAddress1: '123 Main St',
    },
    affiliation: 'Navy',
    edipi: '123567890',
    personal_email: 'test@email.com',
    first_name: 'Tester',
    last_name: 'Testing',
    telephone: '123-555-7890',
  },
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-27T19:17:00.321Z',
        eTag: 'MjAyNC0wMi0yN1QxOToxNzowMC4zMjE3MzFa',
        id: '456',
        moveCode: 'PV96MG',
        orders: {
          authorizedWeight: 13000,
          created_at: '2024-02-27T19:17:00.311Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_9',
          has_dependents: false,
          id: '6d1406d6-152e-475c-9365-2c20b6016541',
          issue_date: '2024-03-01',
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
            created_at: '2024-02-27T18:22:12.471Z',
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
            updated_at: '2024-02-27T18:22:12.471Z',
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
            created_at: '2024-02-27T18:22:12.471Z',
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
            updated_at: '2024-02-27T18:22:12.471Z',
          },
          report_by_date: '2024-03-01',
          service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-27T19:17:00.311Z',
          uploaded_orders: {
            id: 'f2a842f2-4708-442c-87cb-67ff388abf92',
            service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
            uploads: [
              {
                bytes: 1792102,
                contentType: 'image/png',
                createdAt: '2024-02-27T19:17:05.565Z',
                filename: 'Screenshot 2024-02-15 at 12.22.53 PM (3).png',
                id: '2b450af2-a6aa-4143-9990-54cddfa80106',
                status: 'PROCESSING',
                updatedAt: '2024-02-27T19:17:05.565Z',
                url: '/storage/user/f94c8fa6-89de-4594-be6a-ca6f1b4e22d0/uploads/2b450af2-a6aa-4143-9990-54cddfa80106?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [
      {
        createdAt: '2024-02-27T19:16:32.850Z',
        eTag: 'MjAyNC0wMi0yN1QxOToxNjozMi44NTAwNTda',
        id: '123',
        moveCode: 'WWYFP6',
        mtoShipments: [
          {
            createdAt: '2024-02-27T19:27:39.150Z',
            customerRemarks: '',
            destinationAddress: {
              city: 'Flagstaff',
              country: 'United States',
              id: '112e0d7b-90eb-44c4-80b1-5c1214fca0a7',
              postalCode: '86003',
              state: 'AZ',
              streetAddress1: 'N/A',
            },
            eTag: 'MjAyNC0wMi0yN1QxOToyNzozOS4xNTA3MjRa',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: 'f0082986-8e2f-443b-8411-191b3796e572',
            moveTaskOrderID: 'e23d629e-2a73-4b42-886b-fa60cb3db957',
            pickupAddress: {
              city: 'Tulsa',
              id: 'dac5e64d-1a69-4e83-a154-5fca04384544',
              postalCode: '74133',
              state: 'OK',
              streetAddress1: '8711 S 76th E Ave',
              streetAddress2: '',
            },
            requestedDeliveryDate: '2024-02-29',
            requestedPickupDate: '2024-03-01',
            shipmentType: 'HHG',
            status: 'SUBMITTED',
            updatedAt: '2024-02-27T19:27:39.150Z',
          },
        ],
        orders: {
          authorizedWeight: 8000,
          created_at: '2024-02-27T19:16:32.845Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_6',
          has_dependents: false,
          id: '786e60ec-1bbd-48dd-bc12-b6ffcf212c54',
          issue_date: '2024-02-29',
          new_duty_location: {
            address: {
              city: 'Flagstaff',
              country: 'United States',
              id: 'cd51f4db-6195-473a-86cd-c3e5e07640e4',
              postalCode: '86003',
              state: 'AZ',
              streetAddress1: 'n/a',
            },
            address_id: 'cd51f4db-6195-473a-86cd-c3e5e07640e4',
            affiliation: null,
            created_at: '2024-02-27T18:22:12.471Z',
            id: '6ea57f62-2995-4b0c-a0a8-f11782cc8a3b',
            name: 'Flagstaff, AZ 86003',
            updated_at: '2024-02-27T18:22:12.471Z',
          },
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          originDutyLocationGbloc: 'BGAC',
          origin_duty_location: {
            address: {
              city: 'Aberdeen Proving Ground',
              country: 'United States',
              id: 'b6ca003e-1528-4e7c-a43e-830222ca7fb3',
              postalCode: '21005',
              state: 'MD',
              streetAddress1: 'n/a',
            },
            address_id: 'b6ca003e-1528-4e7c-a43e-830222ca7fb3',
            affiliation: 'ARMY',
            created_at: '2024-02-27T18:22:12.471Z',
            id: '61e092c4-575c-458a-9c3f-b93ad373c454',
            name: 'Aberdeen Proving Ground, MD 21005',
            transportation_office: {
              address: {
                city: 'Aberdeen Proving Ground',
                country: 'United States',
                id: 'ac4dbfa5-3068-4f8f-99d1-3cd850412bb9',
                postalCode: '21005',
                state: 'MD',
                streetAddress1: '4305 Susquehanna Ave',
                streetAddress2: 'Room 307',
              },
              created_at: '2018-05-28T14:27:41.772Z',
              gbloc: 'BGAC',
              id: '6a27dfbd-2a49-485f-86dd-49475d5facef',
              name: 'PPPO Aberdeen Proving Ground - USA',
              phone_lines: [],
              updated_at: '2018-05-28T14:27:41.772Z',
            },
            transportation_office_id: '6a27dfbd-2a49-485f-86dd-49475d5facef',
            updated_at: '2024-02-27T18:22:12.471Z',
          },
          report_by_date: '2024-02-29',
          service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-27T19:16:32.845Z',
          uploaded_orders: {
            id: 'b392f96f-20a0-43d2-9ca3-643cfd3b4182',
            service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
            uploads: [
              {
                bytes: 1137126,
                contentType: 'image/png',
                createdAt: '2024-02-27T19:16:38.998Z',
                filename: 'Screenshot 2024-02-15 at 12.22.53 PM.png',
                id: 'bc6c0e2d-fbef-4c32-8487-92c14b613040',
                status: 'PROCESSING',
                updatedAt: '2024-02-27T19:16:38.998Z',
                url: '/storage/user/f94c8fa6-89de-4594-be6a-ca6f1b4e22d0/uploads/bc6c0e2d-fbef-4c32-8487-92c14b613040?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
      {
        createdAt: '2024-02-27T19:17:00.321Z',
        eTag: 'MjAyNC0wMi0yN1QxOToxNzowMC4zMjE3MzFa',
        id: 'testSubmittedMove',
        moveCode: 'PV96MG',
        orders: {
          authorizedWeight: 13000,
          created_at: '2024-02-27T19:17:00.311Z',
          entitlement: {
            proGear: 2000,
            proGearSpouse: 500,
          },
          grade: 'E_9',
          has_dependents: false,
          id: '6d1406d6-152e-475c-9365-2c20b6016541',
          issue_date: '2024-03-01',
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
            created_at: '2024-02-27T18:22:12.471Z',
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
            updated_at: '2024-02-27T18:22:12.471Z',
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
            created_at: '2024-02-27T18:22:12.471Z',
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
            updated_at: '2024-02-27T18:22:12.471Z',
          },
          report_by_date: '2024-03-01',
          service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
          spouse_has_pro_gear: false,
          status: 'DRAFT',
          updated_at: '2024-02-27T19:17:00.311Z',
          uploaded_orders: {
            id: 'f2a842f2-4708-442c-87cb-67ff388abf92',
            service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
            uploads: [
              {
                bytes: 1792102,
                contentType: 'image/png',
                createdAt: '2024-02-27T19:17:05.565Z',
                filename: 'Screenshot 2024-02-15 at 12.22.53 PM (3).png',
                id: '2b450af2-a6aa-4143-9990-54cddfa80106',
                status: 'PROCESSING',
                updatedAt: '2024-02-27T19:17:05.565Z',
                url: '/storage/user/f94c8fa6-89de-4594-be6a-ca6f1b4e22d0/uploads/2b450af2-a6aa-4143-9990-54cddfa80106?contentType=image%2Fpng',
              },
            ],
          },
        },
        status: 'SUBMITTED',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
  },
  currentOrders: {
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    has_dependents: false,
    issue_date: '2020-08-11',
    grade: 'E_1',
    moves: ['123'],
    origin_duty_location: {
      name: 'Test Duty Location',
      address: {
        postalCode: '123456',
      },
    },
    new_duty_location: {
      name: 'New Test Duty Location',
      address: {
        postalCode: '123456',
      },
    },
    report_by_date: '2020-08-31',
    service_member_id: '666',
    spouse_has_pro_gear: false,
    status: MOVE_STATUSES.DRAFT,
    uploaded_orders: {
      uploads: [],
    },
  },
  currentMove: {
    id: '123',
    locator: 'CXVV3F',
    selected_move_type: 'HHG',
    service_member_id: '666',
    status: MOVE_STATUSES.DRAFT,
  },
  moveIsApproved: true,
  entitlement: {},
  mtoShipment: {
    id: 'testMtoShipment789',
    agents: [],
    customerRemarks: 'please be carefule',
    moveTaskOrderID: '123',
    pickupAddress: {
      city: 'Beverly Hills',
    },
    requestedDeliveryDate: '2020-08-31',
    requestedPickupDate: '2020-08-31',
    shipmentType: 'HHG',
    status: MOVE_STATUSES.SUBMITTED,
    updatedAt: '2020-09-02T21:08:38.392Z',
  },
  mtoShipments: [
    {
      id: 'testMtoShipment789',
      agents: [],
      customerRemarks: 'please be carefule',
      moveTaskOrderID: '123',
      pickupAddress: {
        city: 'Beverly Hills',
      },
      requestedDeliveryDate: '2020-08-31',
      requestedPickupDate: '2020-08-31',
      shipmentType: 'HHG',
      status: MOVE_STATUSES.SUBMITTED,
      updatedAt: '2020-09-02T21:08:38.392Z',
    },
  ],
  onDidMount: jest.fn(),
  showLoggedInUser: jest.fn(),
  updateShipmentList: jest.fn(),
  setMsg: jest.fn(),
  updateAllMoves: jest.fn(),
};

describe('Summary page', () => {
  describe('if the user can add another shipment', () => {
    selectCurrentMoveFromAllMoves.mockImplementation(() => testMove);
    it('displays the Add Another Shipment section', () => {
      renderWithRouterProp(<Summary {...testProps} />, {
        path: customerRoutes.MOVE_REVIEW_PATH,
        params: { moveId: '123' },
      });

      expect(screen.getByRole('link', { name: 'Add another shipment' })).toHaveAttribute(
        'href',
        '/moves/123/shipment-type',
      );
    });

    it('displays contact local PPPO office message', async () => {
      renderWithRouterProp(<Summary {...testProps} />, {
        path: customerRoutes.MOVE_REVIEW_PATH,
        params: { moveId: '123' },
      });
      expect(await screen.findByText(/\*To change these fields, contact your local PPPO office/)).toBeInTheDocument();
    });

    it('displays a button that opens a modal', async () => {
      renderWithRouterProp(<Summary {...testProps} />, {
        path: customerRoutes.MOVE_REVIEW_PATH,
        params: { moveId: '123' },
      });

      expect(
        screen.queryByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).not.toBeInTheDocument();

      expect(screen.getByTitle('Help with adding shipments')).toBeInTheDocument();
      await userEvent.click(screen.getByTitle('Help with adding shipments'));

      expect(
        screen.getByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).toBeInTheDocument();
    });

    it('add shipment modal displays default text, nothing is disabled', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

      renderWithRouterProp(<Summary {...testProps} />, {
        path: customerRoutes.MOVE_REVIEW_PATH,
        params: { moveId: '123' },
      });

      expect(
        screen.queryByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).not.toBeInTheDocument();

      expect(screen.getByTitle('Help with adding shipments')).toBeInTheDocument();
      await userEvent.click(screen.getByTitle('Help with adding shipments'));

      expect(
        screen.getByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).toBeInTheDocument();

      // verify it display default text in modal when nothing is disabled
      expect(await screen.findByText(/If none of these apply to you, you probably/)).toBeInTheDocument();

      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.PPM);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTS);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTSR);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.BOAT);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.MOBILE_HOME);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE);
    });

    it('add shipment modal displays still in dev mode', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));

      renderWithRouterProp(<Summary {...testProps} />, {
        path: customerRoutes.MOVE_REVIEW_PATH,
        params: { moveId: '123' },
      });

      expect(
        screen.queryByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).not.toBeInTheDocument();

      await userEvent.click(screen.getByTitle('Help with adding shipments'));

      expect(
        screen.getByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).toBeInTheDocument();

      // verify it display default text in modal feature flag is enabled. display under construction text
      expect(
        await screen.findByText(
          /Some shipment types are still being developed and will become available at a later date./,
        ),
      ).toBeInTheDocument();

      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.PPM);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTS);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.NTSR);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.BOAT);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.MOBILE_HOME);
      expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE);
    });
  });
  afterEach(jest.clearAllMocks);
});
