import React from 'react';
import { mount } from 'enzyme';
import { render, screen, queryByTestId } from '@testing-library/react';

import ServicesCounselingMoveInfo from './ServicesCounselingMoveInfo';

import { MockProviders } from 'testUtils';
import { useMoveDetailsQueries } from 'hooks/queries';
import { ORDERS_TYPE, ORDERS_TYPE_DETAILS } from 'constants/orders';
import MOVE_STATUSES from 'constants/moves';

const testMoveCode = '1A5PM3';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: '1A5PM3' }),
}));

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('hooks/queries'),
  useTXOMoveInfoQueries: () => {
    return {
      customerData: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
      order: {
        id: '4321',
        customerID: '2468',
        uploaded_order_id: '2',
        departmentIndicator: 'Navy',
        grade: 'E-6',
        originDutyLocation: {
          name: 'JBSA Lackland',
        },
        destinationDutyLocation: {
          name: 'JB Lewis-McChord',
        },
        report_by_date: '2018-08-01',
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
  useMoveDetailsQueries: jest.fn(),
}));

const newMoveDetailsQuery = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
  },
  order: {
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
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    order_type_detail: ORDERS_TYPE_DETAILS.HHG_PERMITTED,
    tac: '9999',
  },
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
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
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
  mtoServiceItems: [],
  mtoAgents: [],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const loadingReturnValue = {
  ...newMoveDetailsQuery,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...newMoveDetailsQuery,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('Services Counseling Move Info Container', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useMoveDetailsQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}/details`]}>
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMoveDetailsQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}/details`]}>
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);
    it('should render the move tab container', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}/details`]}>
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );

      expect(wrapper.find('CustomerHeader').exists()).toBe(true);
    });
    it('should render the system error when there is an error', () => {
      render(
        <MockProviders
          initialState={{ interceptor: { hasRecentError: true, traceId: 'some-trace-id' } }}
          initialEntries={[`/counseling/moves/${testMoveCode}/details`]}
        >
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );
      expect(screen.getByText('Technical Help Desk').closest('a')).toHaveAttribute(
        'href',
        'https://move.mil/customer-service#technical-help-desk',
      );
      expect(screen.getByTestId('system-error').textContent).toEqual(
        "Something isn't working, but we're not sure what. Wait a minute and try again.If that doesn't fix it, contact the Technical Help Desk and give them this code: some-trace-id",
      );
    });
    it('should not render system error when there is not an error', () => {
      render(
        <MockProviders
          initialState={{ interceptor: { hasRecentError: false, traceId: '' } }}
          initialEntries={[`/counseling/moves/${testMoveCode}/details`]}
        >
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );
      expect(queryByTestId(document.documentElement, 'system-error')).not.toBeInTheDocument();
    });
  });
  describe('routing', () => {
    useMoveDetailsQueries.mockReturnValue(newMoveDetailsQuery);
    it('should handle the Services Counseling Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}/details`]}>
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );

      expect(wrapper.find('ServicesCounselingMoveDetails')).toHaveLength(1);
    });

    it('should redirect from move info root to the Services Counseling Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}`]}>
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('ServicesCounselingMoveDetails');
      expect(renderedRoute).toHaveLength(1);
    });
  });
});
