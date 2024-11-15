import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { v4 } from 'uuid';

import '@testing-library/jest-dom/extend-expect';

import MultiMovesLandingPage from './MultiMovesLandingPage';

import { MockProviders } from 'testUtils';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { ORDERS_TYPE } from 'constants/orders';

// Mock external dependencies
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

jest.mock('utils/featureFlags', () => ({
  detectFlags: jest.fn(() => ({ multiMove: true })),
}));

jest.mock('store/auth/actions', () => ({
  loadUser: jest.fn(),
}));

jest.mock('store/entities/actions', () => ({
  updateAllMoves: jest.fn(),
}));

jest.mock('store/general/actions', () => ({
  setCanAddOrders: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('store/onboarding/actions', () => ({
  initOnboarding: jest.fn(),
}));

jest.mock('shared/Swagger/ducks', () => ({
  loadInternalSchema: jest.fn(),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectAllMoves: jest.fn(),
  selectServiceMemberFromLoggedInUser: jest.fn(),
  selectCanAddOrders: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const defaultProps = {
  showLoggedInUser: jest.fn(),
  updateAllMoves: jest.fn(),
  setCanAddOrders: jest.fn(),
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  serviceMember: {
    affiliation: 'COAST_GUARD',
    backup_contacts: ['bc0c2ec7-252f-41f6-b1ff-4c9bb270ef41'],
    backup_mailing_address: {
      city: 'Beverly Hills',
      id: 'b1adf427-7743-4fbd-950c-d0fcc25168b9',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
    },
    created_at: '2024-02-15T14:43:31.492Z',
    edipi: '8362534852',
    email_is_preferred: true,
    first_name: 'Jim',
    id: v4(),
    is_profile_complete: true,
    last_name: 'Bean',
    orders: [
      '444de44f-608e-4b99-b66b-dc1fce8a12fd',
      'c1786dd4-771c-4b66-bdec-39960f57f890',
      'a6ca098a-effd-492e-bb1c-edd76568c66b',
    ],
    personal_email: 'multiplemoves@PPM.com',
    residential_address: {
      city: 'Beverly Hills',
      id: '8ace1b49-a1ea-4dd0-aa66-e786b2d220f9',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
    },
    telephone: '212-123-4567',
    updated_at: '2024-02-16T20:41:19.454Z',
    user_id: '68f8baa7-ed00-4ad9-ad3c-a849688cb537',
  },
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-01-31T16:29:53.290Z',
        eTag: 'MjAyNC0wMS0zMVQxNjoyOTo1My4yOTA0OTRa',
        id: '9211d4e2-5b92-42bb-9758-7ac1f329a8d6',
        moveCode: 'YJ9M34',
        orders: {
          id: '40475a80-5340-4722-88d1-3cc9764414d6',
          created_at: '2024-01-31T16:29:53.285657Z',
          updated_at: '2024-01-31T16:29:53.285657Z',
          service_member_id: '6686d242-e7af-4a06-afd7-7be423bfca2d',
          issue_date: '2024-01-31T00:00:00Z',
          report_by_date: '2024-02-09T00:00:00Z',
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          orders_type_detail: null,
          has_dependents: false,
          spouse_has_pro_gear: false,
          origin_duty_location: {
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Tinker AFB, OK 73145',
            affiliation: 'AIR_FORCE',
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            address: {
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Tinker AFB',
              state: 'OK',
              postal_code: '73145',
              country: 'United States',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            TransportationOffice: {
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              shipping_office_id: 'c2c440ae-5394-4483-84fb-f872e32126bb',
              ShippingOffice: null,
              name: 'PPPO Tinker AFB - USAF',
              Address: {
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                created_at: '2018-05-28T14:27:40.597383Z',
                updated_at: '2018-05-28T14:27:40.597383Z',
                street_address_1: '7330 Century Blvd',
                street_address_2: 'Bldg 469',
                street_address_3: null,
                city: 'Tinker AFB',
                state: 'OK',
                postal_code: '73145',
                country: 'United States',
              },
              address_id: '410b18bc-b270-4b52-9211-532fffc6f59e',
              latitude: 35.429035,
              longitude: -97.39955,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday – Friday: 0715 – 1600; Limited Service from 1130-1230',
              services: 'Walk-In Help; Briefings; Appointments; QA Inspections',
              note: null,
              gbloc: 'HAFC',
              created_at: '2018-05-28T14:27:40.605679Z',
              updated_at: '2018-05-28T14:27:40.60568Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          origin_duty_location_id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
          new_duty_location_id: '5c182566-0e6e-46f2-9eef-f07963783575',
          new_duty_location: {
            id: '5c182566-0e6e-46f2-9eef-f07963783575',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Fort Sill, OK 73503',
            affiliation: 'ARMY',
            address_id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
            address: {
              id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Fort Sill',
              state: 'OK',
              postal_code: '73503',
              country: 'United States',
            },
            transportation_office_id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
            TransportationOffice: {
              id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
              shipping_office_id: '5a3388e1-6d46-4639-ac8f-a8937dc26938',
              ShippingOffice: null,
              name: 'PPPO Fort Sill - USA',
              Address: {
                id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
                created_at: '2018-05-28T14:27:35.538742Z',
                updated_at: '2018-05-28T14:27:35.538743Z',
                street_address_1: '4700 Mow Way Rd',
                street_address_2: 'Room 110',
                street_address_3: null,
                city: 'Fort Sill',
                state: 'OK',
                postal_code: '73503',
                country: 'United States',
              },
              address_id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
              latitude: 34.647964,
              longitude: -98.41231,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday - Friday 0830-1530; Sat/Sun/Federal Holidays closed',
              services: 'Walk-In Help; Appointments; QA Inspections; Appointments 06 and above',
              note: null,
              gbloc: 'JEAT',
              created_at: '2018-05-28T14:27:35.547257Z',
              updated_at: '2018-05-28T14:27:35.547257Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          uploaded_orders_id: 'f779f6a2-48e2-47fe-87be-d93e8aa711fe',
          status: 'DRAFT',
          grade: 'E_7',
          Entitlement: null,
          entitlement_id: 'a1bf0035-4f28-45b8-af1a-556848d29e44',
          UploadedAmendedOrders: null,
          uploaded_amended_orders_id: null,
          amended_orders_acknowledged_at: null,
          origin_duty_location_gbloc: 'HAFC',
        },
        status: 'DRAFT',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [],
  },
};

const defaultPropsNoMoves = {
  showLoggedInUser: jest.fn(),
  updateAllMoves: jest.fn(),
  setCanAddOrders: jest.fn(),
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  serviceMember: {
    affiliation: 'COAST_GUARD',
    backup_contacts: ['bc0c2ec7-252f-41f6-b1ff-4c9bb270ef41'],
    backup_mailing_address: {
      city: 'Beverly Hills',
      id: 'b1adf427-7743-4fbd-950c-d0fcc25168b9',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
    },
    created_at: '2024-02-15T14:43:31.492Z',
    edipi: '8362534852',
    email_is_preferred: true,
    first_name: 'Jim',
    id: v4(),
    is_profile_complete: true,
    last_name: 'Bean',
    orders: [
      '444de44f-608e-4b99-b66b-dc1fce8a12fd',
      'c1786dd4-771c-4b66-bdec-39960f57f890',
      'a6ca098a-effd-492e-bb1c-edd76568c66b',
    ],
    personal_email: 'multiplemoves@PPM.com',
    residential_address: {
      city: 'Beverly Hills',
      id: '8ace1b49-a1ea-4dd0-aa66-e786b2d220f9',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
    },
    telephone: '212-123-4567',
    updated_at: '2024-02-16T20:41:19.454Z',
    user_id: '68f8baa7-ed00-4ad9-ad3c-a849688cb537',
  },
  serviceMemberMoves: {
    currentMove: [],
    previousMoves: [],
  },
};

const defaultPropsMultipleMove = {
  showLoggedInUser: jest.fn(),
  updateAllMoves: jest.fn(),
  setCanAddOrders: jest.fn(),
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  serviceMember: {
    affiliation: 'COAST_GUARD',
    backup_contacts: ['bc0c2ec7-252f-41f6-b1ff-4c9bb270ef41'],
    backup_mailing_address: {
      city: 'Beverly Hills',
      id: 'b1adf427-7743-4fbd-950c-d0fcc25168b9',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
    },
    created_at: '2024-02-15T14:43:31.492Z',
    edipi: '8362534852',
    email_is_preferred: true,
    first_name: 'Jim',
    id: v4(),
    is_profile_complete: true,
    last_name: 'Bean',
    orders: [
      '444de44f-608e-4b99-b66b-dc1fce8a12fd',
      'c1786dd4-771c-4b66-bdec-39960f57f890',
      'a6ca098a-effd-492e-bb1c-edd76568c66b',
    ],
    personal_email: 'multiplemoves@PPM.com',
    residential_address: {
      city: 'Beverly Hills',
      id: '8ace1b49-a1ea-4dd0-aa66-e786b2d220f9',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
    },
    telephone: '212-123-4567',
    updated_at: '2024-02-16T20:41:19.454Z',
    user_id: '68f8baa7-ed00-4ad9-ad3c-a849688cb537',
  },
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-01-31T16:29:53.290Z',
        eTag: 'MjAyNC0wMS0zMVQxNjoyOTo1My4yOTA0OTRa',
        id: '9211d4e2-5b92-42bb-9758-7ac1f329a8d6',
        moveCode: 'YJ9M34',
        orders: {
          id: '40475a80-5340-4722-88d1-3cc9764414d6',
          created_at: '2024-01-31T16:29:53.285657Z',
          updated_at: '2024-01-31T16:29:53.285657Z',
          service_member_id: '6686d242-e7af-4a06-afd7-7be423bfca2d',
          issue_date: '2024-01-31T00:00:00Z',
          report_by_date: '2024-02-09T00:00:00Z',
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          orders_type_detail: null,
          has_dependents: false,
          spouse_has_pro_gear: false,
          origin_duty_location: {
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Tinker AFB, OK 73145',
            affiliation: 'AIR_FORCE',
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            address: {
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Tinker AFB',
              state: 'OK',
              postal_code: '73145',
              country: 'United States',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            TransportationOffice: {
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              shipping_office_id: 'c2c440ae-5394-4483-84fb-f872e32126bb',
              ShippingOffice: null,
              name: 'PPPO Tinker AFB - USAF',
              Address: {
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                created_at: '2018-05-28T14:27:40.597383Z',
                updated_at: '2018-05-28T14:27:40.597383Z',
                street_address_1: '7330 Century Blvd',
                street_address_2: 'Bldg 469',
                street_address_3: null,
                city: 'Tinker AFB',
                state: 'OK',
                postal_code: '73145',
                country: 'United States',
              },
              address_id: '410b18bc-b270-4b52-9211-532fffc6f59e',
              latitude: 35.429035,
              longitude: -97.39955,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday – Friday: 0715 – 1600; Limited Service from 1130-1230',
              services: 'Walk-In Help; Briefings; Appointments; QA Inspections',
              note: null,
              gbloc: 'HAFC',
              created_at: '2018-05-28T14:27:40.605679Z',
              updated_at: '2018-05-28T14:27:40.60568Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          origin_duty_location_id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
          new_duty_location_id: '5c182566-0e6e-46f2-9eef-f07963783575',
          new_duty_location: {
            id: '5c182566-0e6e-46f2-9eef-f07963783575',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Fort Sill, OK 73503',
            affiliation: 'ARMY',
            address_id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
            address: {
              id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Fort Sill',
              state: 'OK',
              postal_code: '73503',
              country: 'United States',
            },
            transportation_office_id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
            TransportationOffice: {
              id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
              shipping_office_id: '5a3388e1-6d46-4639-ac8f-a8937dc26938',
              ShippingOffice: null,
              name: 'PPPO Fort Sill - USA',
              Address: {
                id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
                created_at: '2018-05-28T14:27:35.538742Z',
                updated_at: '2018-05-28T14:27:35.538743Z',
                street_address_1: '4700 Mow Way Rd',
                street_address_2: 'Room 110',
                street_address_3: null,
                city: 'Fort Sill',
                state: 'OK',
                postal_code: '73503',
                country: 'United States',
              },
              address_id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
              latitude: 34.647964,
              longitude: -98.41231,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday - Friday 0830-1530; Sat/Sun/Federal Holidays closed',
              services: 'Walk-In Help; Appointments; QA Inspections; Appointments 06 and above',
              note: null,
              gbloc: 'JEAT',
              created_at: '2018-05-28T14:27:35.547257Z',
              updated_at: '2018-05-28T14:27:35.547257Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          uploaded_orders_id: 'f779f6a2-48e2-47fe-87be-d93e8aa711fe',
          status: 'DRAFT',
          grade: 'E_7',
          Entitlement: null,
          entitlement_id: 'a1bf0035-4f28-45b8-af1a-556848d29e44',
          UploadedAmendedOrders: null,
          uploaded_amended_orders_id: null,
          amended_orders_acknowledged_at: null,
          origin_duty_location_gbloc: 'HAFC',
        },
        status: 'DRAFT',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [
      {
        createdAt: '2024-01-31T16:29:53.290Z',
        eTag: 'MjAyNC0wMS0zMVQxNjoyOTo1My4yOTA0OTRa',
        id: '9211d4e2-5b92-42bb-9758-7ac1f329a8d6',
        moveCode: 'ABC123',
        orders: {
          id: '40475a80-5340-4722-88d1-3cc9764414d6',
          created_at: '2024-01-31T16:29:53.285657Z',
          updated_at: '2024-01-31T16:29:53.285657Z',
          service_member_id: '6686d242-e7af-4a06-afd7-7be423bfca2d',
          issue_date: '2024-01-31T00:00:00Z',
          report_by_date: '2024-02-09T00:00:00Z',
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          orders_type_detail: null,
          has_dependents: false,
          spouse_has_pro_gear: false,
          origin_duty_location: {
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Tinker AFB, OK 73145',
            affiliation: 'AIR_FORCE',
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            address: {
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Tinker AFB',
              state: 'OK',
              postal_code: '73145',
              country: 'United States',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            TransportationOffice: {
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              shipping_office_id: 'c2c440ae-5394-4483-84fb-f872e32126bb',
              ShippingOffice: null,
              name: 'PPPO Tinker AFB - USAF',
              Address: {
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                created_at: '2018-05-28T14:27:40.597383Z',
                updated_at: '2018-05-28T14:27:40.597383Z',
                street_address_1: '7330 Century Blvd',
                street_address_2: 'Bldg 469',
                street_address_3: null,
                city: 'Tinker AFB',
                state: 'OK',
                postal_code: '73145',
                country: 'United States',
              },
              address_id: '410b18bc-b270-4b52-9211-532fffc6f59e',
              latitude: 35.429035,
              longitude: -97.39955,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday – Friday: 0715 – 1600; Limited Service from 1130-1230',
              services: 'Walk-In Help; Briefings; Appointments; QA Inspections',
              note: null,
              gbloc: 'HAFC',
              created_at: '2018-05-28T14:27:40.605679Z',
              updated_at: '2018-05-28T14:27:40.60568Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          origin_duty_location_id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
          new_duty_location_id: '5c182566-0e6e-46f2-9eef-f07963783575',
          new_duty_location: {
            id: '5c182566-0e6e-46f2-9eef-f07963783575',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Fort Sill, OK 73503',
            affiliation: 'ARMY',
            address_id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
            address: {
              id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Fort Sill',
              state: 'OK',
              postal_code: '73503',
              country: 'United States',
            },
            transportation_office_id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
            TransportationOffice: {
              id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
              shipping_office_id: '5a3388e1-6d46-4639-ac8f-a8937dc26938',
              ShippingOffice: null,
              name: 'PPPO Fort Sill - USA',
              Address: {
                id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
                created_at: '2018-05-28T14:27:35.538742Z',
                updated_at: '2018-05-28T14:27:35.538743Z',
                street_address_1: '4700 Mow Way Rd',
                street_address_2: 'Room 110',
                street_address_3: null,
                city: 'Fort Sill',
                state: 'OK',
                postal_code: '73503',
                country: 'United States',
              },
              address_id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
              latitude: 34.647964,
              longitude: -98.41231,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday - Friday 0830-1530; Sat/Sun/Federal Holidays closed',
              services: 'Walk-In Help; Appointments; QA Inspections; Appointments 06 and above',
              note: null,
              gbloc: 'JEAT',
              created_at: '2018-05-28T14:27:35.547257Z',
              updated_at: '2018-05-28T14:27:35.547257Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          uploaded_orders_id: 'f779f6a2-48e2-47fe-87be-d93e8aa711fe',
          status: 'DRAFT',
          grade: 'E_7',
          Entitlement: null,
          entitlement_id: 'a1bf0035-4f28-45b8-af1a-556848d29e44',
          UploadedAmendedOrders: null,
          uploaded_amended_orders_id: null,
          amended_orders_acknowledged_at: null,
          origin_duty_location_gbloc: 'HAFC',
        },
        status: 'APPROVED',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
      {
        createdAt: '2024-02-18T16:29:53.290Z',
        eTag: 'MjAyNC0wMS0zMVQxNjoyOTo1My4yOTA0OTRb',
        id: '9211d4e2-5b92-42bb-9758-7ac1f329a8d7',
        moveCode: 'DEF456',
        orders: {
          id: '40475a80-5340-4722-88d1-3cc9764414d7',
          created_at: '2024-01-31T16:29:53.285657Z',
          updated_at: '2024-01-31T16:29:53.285657Z',
          service_member_id: '6686d242-e7af-4a06-afd7-7be423bfca2d',
          issue_date: '2024-01-31T00:00:00Z',
          report_by_date: '2024-02-09T00:00:00Z',
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          orders_type_detail: null,
          has_dependents: false,
          spouse_has_pro_gear: false,
          origin_duty_location: {
            id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Tinker AFB, OK 73145',
            affiliation: 'AIR_FORCE',
            address_id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
            address: {
              id: '7e3ea97c-da9f-4fa1-8a11-87063c857635',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Tinker AFB',
              state: 'OK',
              postal_code: '73145',
              country: 'United States',
            },
            transportation_office_id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
            TransportationOffice: {
              id: '7876373d-57e4-4cde-b11f-c26a8feee9e8',
              shipping_office_id: 'c2c440ae-5394-4483-84fb-f872e32126bb',
              ShippingOffice: null,
              name: 'PPPO Tinker AFB - USAF',
              Address: {
                id: '410b18bc-b270-4b52-9211-532fffc6f59e',
                created_at: '2018-05-28T14:27:40.597383Z',
                updated_at: '2018-05-28T14:27:40.597383Z',
                street_address_1: '7330 Century Blvd',
                street_address_2: 'Bldg 469',
                street_address_3: null,
                city: 'Tinker AFB',
                state: 'OK',
                postal_code: '73145',
                country: 'United States',
              },
              address_id: '410b18bc-b270-4b52-9211-532fffc6f59e',
              latitude: 35.429035,
              longitude: -97.39955,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday – Friday: 0715 – 1600; Limited Service from 1130-1230',
              services: 'Walk-In Help; Briefings; Appointments; QA Inspections',
              note: null,
              gbloc: 'HAFC',
              created_at: '2018-05-28T14:27:40.605679Z',
              updated_at: '2018-05-28T14:27:40.60568Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          origin_duty_location_id: '2d6eab7d-1a21-4f29-933e-ee8fa7dbc314',
          new_duty_location_id: '5c182566-0e6e-46f2-9eef-f07963783575',
          new_duty_location: {
            id: '5c182566-0e6e-46f2-9eef-f07963783575',
            created_at: '2024-01-26T16:46:34.047004Z',
            updated_at: '2024-01-26T16:46:34.047004Z',
            name: 'Fort Sill, OK 73503',
            affiliation: 'ARMY',
            address_id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
            address: {
              id: 'ed62ba0b-a3cb-47ac-81ae-5b27ade4592b',
              created_at: '2024-01-26T16:46:34.047004Z',
              updated_at: '2024-01-26T16:46:34.047004Z',
              street_address_1: 'n/a',
              street_address_2: null,
              street_address_3: null,
              city: 'Fort Sill',
              state: 'OK',
              postal_code: '73503',
              country: 'United States',
            },
            transportation_office_id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
            TransportationOffice: {
              id: '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9',
              shipping_office_id: '5a3388e1-6d46-4639-ac8f-a8937dc26938',
              ShippingOffice: null,
              name: 'PPPO Fort Sill - USA',
              Address: {
                id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
                created_at: '2018-05-28T14:27:35.538742Z',
                updated_at: '2018-05-28T14:27:35.538743Z',
                street_address_1: '4700 Mow Way Rd',
                street_address_2: 'Room 110',
                street_address_3: null,
                city: 'Fort Sill',
                state: 'OK',
                postal_code: '73503',
                country: 'United States',
              },
              address_id: 'abbc0af9-2394-4e36-be84-811ad8f6060b',
              latitude: 34.647964,
              longitude: -98.41231,
              PhoneLines: null,
              Emails: null,
              hours: 'Monday - Friday 0830-1530; Sat/Sun/Federal Holidays closed',
              services: 'Walk-In Help; Appointments; QA Inspections; Appointments 06 and above',
              note: null,
              gbloc: 'JEAT',
              created_at: '2018-05-28T14:27:35.547257Z',
              updated_at: '2018-05-28T14:27:35.547257Z',
              provides_ppm_closeout: true,
            },
            provides_services_counseling: true,
          },
          uploaded_orders_id: 'f779f6a2-48e2-47fe-87be-d93e8aa711fe',
          status: 'DRAFT',
          grade: 'E_7',
          Entitlement: null,
          entitlement_id: 'a1bf0035-4f28-45b8-af1a-556848d29e44',
          UploadedAmendedOrders: null,
          uploaded_amended_orders_id: null,
          amended_orders_acknowledged_at: null,
          origin_duty_location_gbloc: 'HAFC',
        },
        status: 'DRAFT',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
  },
};

describe('MultiMovesLandingPage', () => {
  it('renders the component with moves', () => {
    render(
      <MockProviders>
        <MultiMovesLandingPage {...defaultProps} />
      </MockProviders>,
    );

    // Check for specific elements
    expect(screen.getByTestId('customerHeader')).toBeInTheDocument();
    expect(screen.getByTestId('welcomeHeaderPrevMoves')).toBeInTheDocument();
    expect(screen.getByText('Welcome to MilMove!')).toBeInTheDocument();
    expect(screen.getByText('Create a Move')).toBeInTheDocument();

    // Assuming there are two move headers and corresponding move containers
    expect(screen.getAllByText('Current Move')).toHaveLength(1);
    expect(screen.getAllByText('Previous Moves')).toHaveLength(1);
  });

  it('renders move data correctly if one move', () => {
    render(
      <MockProviders>
        <MultiMovesLandingPage {...defaultProps} />
      </MockProviders>,
    );

    expect(screen.getByText('Jim Bean')).toBeInTheDocument();
    expect(screen.getByText('#YJ9M34')).toBeInTheDocument();
    expect(screen.getByTestId('welcomeHeaderPrevMoves')).toBeInTheDocument();
    expect(screen.getByTestId('createMoveBtn')).toBeInTheDocument();
    expect(screen.getByTestId('currentMoveHeader')).toBeInTheDocument();
    expect(screen.getByTestId('currentMoveContainer')).toBeInTheDocument();
    expect(screen.getByTestId('prevMovesHeader')).toBeInTheDocument();
    expect(screen.getByText('You have no previous moves.')).toBeInTheDocument();
  });

  it('renders move data correctly if multiple moves', () => {
    render(
      <MockProviders>
        <MultiMovesLandingPage {...defaultPropsMultipleMove} />
      </MockProviders>,
    );

    expect(screen.getByText('Jim Bean')).toBeInTheDocument();
    expect(screen.getByText('#YJ9M34')).toBeInTheDocument();
    expect(screen.getByTestId('welcomeHeaderPrevMoves')).toBeInTheDocument();
    expect(screen.getByTestId('createMoveBtn')).toBeInTheDocument();
    expect(screen.getByTestId('currentMoveHeader')).toBeInTheDocument();
    expect(screen.getByTestId('currentMoveContainer')).toBeInTheDocument();
    expect(screen.getByTestId('prevMovesHeader')).toBeInTheDocument();
    expect(screen.getByText('#ABC123')).toBeInTheDocument();
    expect(screen.getByText('#DEF456')).toBeInTheDocument();
  });

  it('navigates the user when create move button is clicked', () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => defaultPropsMultipleMove.serviceMember);
    render(
      <MockProviders>
        <MultiMovesLandingPage {...defaultPropsMultipleMove} />
      </MockProviders>,
    );

    expect(screen.getByTestId('createMoveBtn')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('createMoveBtn'));
    expect(defaultPropsMultipleMove.setCanAddOrders).toHaveBeenCalledWith(true);
    expect(mockNavigate).toHaveBeenCalled();
  });

  it('renders move data correctly if no moves', () => {
    render(
      <MockProviders>
        <MultiMovesLandingPage {...defaultPropsNoMoves} />
      </MockProviders>,
    );

    expect(screen.getByText('Jim Bean')).toBeInTheDocument();
    expect(screen.getByTestId('currentMoveHeader')).toBeInTheDocument();
    expect(screen.getByText('You do not have a current move.')).toBeInTheDocument();
    expect(screen.getByTestId('prevMovesHeader')).toBeInTheDocument();
    expect(screen.getByText('You have no previous moves.')).toBeInTheDocument();
  });
});
