import React from 'react';
import { mount } from 'enzyme';
import { queryByTestId, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { Provider } from 'react-redux';

import TXOMoveInfo from './TXOMoveInfo';

import { mockPage, MockProviders } from 'testUtils';
import { useTXOMoveInfoQueries, useUserQueries } from 'hooks/queries';
import { tooRoutes } from 'constants/routes';
import { roleTypes } from 'constants/userRoles';
import { configureStore } from 'shared/store';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

mockPage('pages/Office/MoveDetails/MoveDetails');
mockPage('pages/Office/MoveDocumentWrapper/MoveDocumentWrapper');
mockPage('pages/Office/MoveTaskOrder/MoveTaskOrder');
mockPage('pages/Office/PaymentRequestReview/PaymentRequestReview');
mockPage('pages/Office/ReviewBillableWeight/ReviewBillableWeight');
mockPage('pages/Office/CustomerSupportRemarks/CustomerSupportRemarks');
mockPage('pages/Office/EvaluationReports/EvaluationReports');
mockPage('pages/Office/EvaluationReport/EvaluationReport');
mockPage('pages/Office/EvaluationViolations/EvaluationViolations');
mockPage('pages/Office/MoveHistory/MoveHistory');
mockPage('pages/Office/MovePaymentRequests/MovePaymentRequests');
mockPage('pages/Office/CustomerInfo/CustomerInfo');
mockPage('pages/Office/Forbidden/Forbidden');

const testMoveCode = '1A5PM3';
const loggedInTIOState = {
  auth: {
    activeRole: roleTypes.TIO,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId234: {
        id: 'userId234',
        roles: [{ roleType: roleTypes.TIO }],
      },
    },
  },
};

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('hooks/queries'),
  useTXOMoveInfoQueries: jest.fn(),
  useUserQueries: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const basicUseTXOMoveInfoQueriesValue = {
  customerData: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
  move: {
    lockedByOfficeUserID: '2744435d-7ba8-4cc5-bae5-f302c72c966e',
  },
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
      address: {
        postalCode: '94535',
      },
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

const user = {
  isLoading: false,
  isError: false,
  data: {
    office_user: { id: '2744435d-7ba8-4cc5-bae5-f302c72c9632' },
  },
};

// Render the TXO Move Info page with redux and routing setup.
// Nestes the TXOMoveInfo under /moves/:moveCode/* as done in the app since the TXOMoveInfo component uses nested pathing.
const renderTXOMoveInfo = (nestedPath = 'details', state = {}) => {
  const mockStore = configureStore({
    ...loggedInTIOState,
    ...state,
  });

  return render(
    <MemoryRouter initialEntries={[`/moves/${testMoveCode}/${nestedPath}`]}>
      <Provider store={mockStore.store}>
        <Routes>
          <Route key="txoMoveInfoRoute" path="/moves/:moveCode/*" element={<TXOMoveInfo />} />
        </Routes>
      </Provider>
    </MemoryRouter>,
  );
};

beforeEach(() => {
  useTXOMoveInfoQueries.mockReturnValue(basicUseTXOMoveInfoQueriesValue);
  useUserQueries.mockReturnValue(user);
});

describe('TXO Move Info Container', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useTXOMoveInfoQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders path={tooRoutes.BASE_MOVE_VIEW_PATH} params={{ moveCode: testMoveCode }}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useTXOMoveInfoQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders path={tooRoutes.BASE_MOVE_VIEW_PATH} params={{ moveCode: testMoveCode }}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('should render the move tab container', () => {
      useTXOMoveInfoQueries.mockReturnValue(basicUseTXOMoveInfoQueriesValue);
      const wrapper = mount(
        <MockProviders path={tooRoutes.BASE_MOVE_VIEW_PATH} params={{ moveCode: testMoveCode }}>
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
      renderTXOMoveInfo('', { interceptor: { hasRecentError: true, traceId: 'some-trace-id' } });

      expect(screen.getByText('Technical Help Desk').closest('a')).toHaveAttribute(
        'href',
        'mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil',
      );
      expect(screen.getByTestId('system-error').textContent).toEqual(
        "Something isn't working, but we're not sure what. Wait a minute and try again.If that doesn't fix it, contact the Technical Help Desk (usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil) and give them this code: some-trace-id",
      );
    });

    it('should not render system error when there is not an error', () => {
      renderTXOMoveInfo('', { interceptor: { hasRecentError: false, traceId: '' } });

      expect(queryByTestId(document.documentElement, 'system-error')).not.toBeInTheDocument();
    });

    it('renders a lock icon when move lock flag is on', async () => {
      isBooleanFlagEnabled.mockResolvedValue(true);
      useTXOMoveInfoQueries.mockReturnValue(basicUseTXOMoveInfoQueriesValue);

      render(
        <MockProviders path={tooRoutes.BASE_MOVE_VIEW_PATH} params={{ moveCode: testMoveCode }}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      await waitFor(() => {
        const banner = screen.queryByTestId('locked-move-banner');
        expect(banner).toBeInTheDocument();
      });
    });
    it('does NOT render a lock icon when move lock flag is off', async () => {
      isBooleanFlagEnabled.mockResolvedValue(false);
      useTXOMoveInfoQueries.mockReturnValue(basicUseTXOMoveInfoQueriesValue);

      render(
        <MockProviders path={tooRoutes.BASE_MOVE_VIEW_PATH} params={{ moveCode: testMoveCode }}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      await waitFor(() => {
        const banner = screen.queryByTestId('locked-move-banner');
        expect(banner).not.toBeInTheDocument();
      });
    });
  });

  describe('routing', () => {
    it.each([
      ['Move Details', '/'],
      ['Move Details', 'details'],
      ['Move Document Wrapper', 'allowances'],
      ['Move Document Wrapper', 'orders'],
      ['Move Task Order', 'mto'],
      ['Payment Request Review', 'payment-requests/REQ123'],
      ['Move Payment Requests', 'payment-requests'],
      ['Review Billable Weight', 'billable-weight'],
      ['Customer Support Remarks', 'customer-support-remarks'],
      ['Evaluation Reports', 'evaluation-reports'],
      ['Move History', 'history'],
      ['Customer Info', 'customer'],
      ['Forbidden', 'evaluation-reports/123'], // Permission restricted
      ['Forbidden', 'evaluation-reports/report123/violations'], // Permission restricted
    ])('should render the %s component when at the route: /moves/:moveCode/%s', async (componentName, nestedPath) => {
      renderTXOMoveInfo(nestedPath);

      // Wait for loading to finish
      await waitFor(() => expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument());

      // Assert that the mock component is rendered
      await expect(screen.getByText(`Mock ${componentName} Component`)).toBeInTheDocument();
    });
  });
});
