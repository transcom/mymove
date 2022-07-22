import React from 'react';
import { mount } from 'enzyme';
import { queryByTestId, render, screen } from '@testing-library/react';

import TXOMoveInfo from './TXOMoveInfo';

import { MockProviders } from 'testUtils';
import { useTXOMoveInfoQueries } from 'hooks/queries';

const testMoveCode = '1A5PM3';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: '1A5PM3' }),
}));

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('hooks/queries'),
  useTXOMoveInfoQueries: jest.fn(),
}));

const basicUseTXOMoveInfoQueriesValue = {
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

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('TXO Move Info Container', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useTXOMoveInfoQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/details`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useTXOMoveInfoQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/details`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('should render the move tab container', () => {
      useTXOMoveInfoQueries.mockReturnValueOnce(basicUseTXOMoveInfoQueriesValue);
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/details`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      expect(wrapper.find('CustomerHeader').exists()).toBe(true);
      expect(wrapper.find('header.nav-header').exists()).toBe(true);
      expect(wrapper.find('nav.tabNav').exists()).toBe(true);
      expect(wrapper.find('li.tabItem').length).toEqual(6);

      expect(wrapper.find('span.tab-title').at(0).text()).toContain('Move details');
      expect(wrapper.find('span.tab-title + span').at(0).exists()).toBe(false);
      expect(wrapper.find('span.tab-title').at(1).text()).toContain('Move task order');
      expect(wrapper.find('span.tab-title').at(2).text()).toContain('Payment requests');
      expect(wrapper.find('span.tab-title').at(3).text()).toContain('Customer support remarks');
      expect(wrapper.find('span.tab-title').at(4).text()).toContain('Quality assurance');
      expect(wrapper.find('span.tab-title').at(5).text()).toContain('Move history');

      expect(wrapper.find('li.tabItem a').at(0).prop('href')).toEqual(`/moves/${testMoveCode}/details`);
      expect(wrapper.find('li.tabItem a').at(1).prop('href')).toEqual(`/moves/${testMoveCode}/mto`);
      expect(wrapper.find('li.tabItem a').at(2).prop('href')).toEqual(`/moves/${testMoveCode}/payment-requests`);
      expect(wrapper.find('li.tabItem a').at(3).prop('href')).toEqual(
        `/moves/${testMoveCode}/customer-support-remarks`,
      );
      expect(wrapper.find('li.tabItem a').at(4).prop('href')).toEqual(`/moves/${testMoveCode}/evaluation-reports`);
      expect(wrapper.find('li.tabItem a').at(5).prop('href')).toEqual(`/moves/${testMoveCode}/history`);
    });
    it('should render the system error when there is an error', () => {
      useTXOMoveInfoQueries.mockReturnValueOnce(basicUseTXOMoveInfoQueriesValue);

      render(
        <MockProviders
          initialState={{ interceptor: { hasRecentError: true, traceId: 'some-trace-id' } }}
          initialEntries={[`/moves/${testMoveCode}/details`]}
        >
          <TXOMoveInfo />
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
      useTXOMoveInfoQueries.mockReturnValueOnce(basicUseTXOMoveInfoQueriesValue);
      render(
        <MockProviders
          initialState={{ interceptor: { hasRecentError: false, traceId: '' } }}
          initialEntries={[`/moves/${testMoveCode}/details`]}
        >
          <TXOMoveInfo />
        </MockProviders>,
      );
      expect(queryByTestId(document.documentElement, 'system-error')).not.toBeInTheDocument();
    });
  });

  describe('routing', () => {
    beforeAll(() => {
      useTXOMoveInfoQueries.mockReturnValue(basicUseTXOMoveInfoQueriesValue);
    });

    it('should handle the Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/details`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      expect(wrapper.find('MoveDetails')).toHaveLength(1);
    });

    it('should redirect from move info root to the Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('MoveDetails');
      expect(renderedRoute).toHaveLength(1);
    });

    it('should handle the Move Orders route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/orders`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual(['/moves/:moveCode/allowances', '/moves/:moveCode/orders']);
    });

    it('should handle the Allowances route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/allowances`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual(['/moves/:moveCode/allowances', '/moves/:moveCode/orders']);
    });

    it('should handle the Move Task Order route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/mto`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/moves/:moveCode/mto');
    });

    it('should handle the Move Payment Requests route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/payment-requests`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/moves/:moveCode/payment-requests');
    });

    it('should handle the Billable Weight route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/billable-weight`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/moves/:moveCode/billable-weight');
    });

    it('should handle the Move History route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveCode}/history`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/moves/:moveCode/history');
    });
  });
});
