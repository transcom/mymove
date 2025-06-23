import React from 'react';
import { render, screen, queryByTestId, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { Provider } from 'react-redux';
import { v4 } from 'uuid';

import ServicesCounselingMoveInfo from './ServicesCounselingMoveInfo';

import { mockPage, ReactQueryWrapper } from 'testUtils';
import { roleTypes } from 'constants/userRoles';
import { configureStore } from 'shared/store';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';

const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: 'PPM',
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    weightTickets: [],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

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

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('hooks/queries'),
  useTXOMoveInfoQueries: () => {
    return {
      customerData: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
      move: {
        lockedByOfficeUserID: '2744435d-7ba8-4cc5-bae5-f302c72c966e',
        lockExpiresAt: '2099-10-15T23:48:35.420Z',
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
        },
        report_by_date: '2018-08-01',
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
  usePPMShipmentAndDocsOnlyQueries: jest.fn(),
}));

mockPage('pages/Office/ServicesCounselingMoveDetails/ServicesCounselingMoveDetails');
mockPage('pages/Office/PPM/ReviewDocuments/ReviewDocuments');
mockPage('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment');
mockPage('pages/Office/CustomerSupportRemarks/CustomerSupportRemarks');
mockPage('pages/Office/MoveTaskOrder/MoveTaskOrder');
mockPage('pages/Office/MoveHistory/MoveHistory');
mockPage('pages/Office/SupportingDocuments/SupportingDocuments');
mockPage('pages/Office/ServicesCounselingMoveDocumentWrapper/ServicesCounselingMoveDocumentWrapper');
mockPage('pages/Office/CustomerInfo/CustomerInfo');
mockPage('pages/Office/ServicesCounselingEditShipmentDetails/ServicesCounselingEditShipmentDetails');
mockPage('pages/Office/ServicesCounselingReviewShipmentWeights/ServicesCounselingReviewShipmentWeights');

const renderSCMoveInfo = (nestedPath = 'details', state = {}) => {
  const mockStore = configureStore({
    ...loggedInTIOState,
    ...state,
  });

  // Render the SC Move Info page with redux and routing setup.
  // Nestes the SC Move Info under /counseling/moves/:moveCode/* as done in the app since the SC Move Info component uses nested pathing.
  return render(
    <MemoryRouter initialEntries={[`/counseling/moves/${testMoveCode}/${nestedPath}`]}>
      <Provider store={mockStore.store}>
        <ReactQueryWrapper>
          <Routes>
            <Route
              key="scMoveInfoRoute"
              path="/counseling/moves/:moveCode/*"
              element={<ServicesCounselingMoveInfo />}
            />
          </Routes>
        </ReactQueryWrapper>
      </Provider>
    </MemoryRouter>,
  );
};

describe('Services Counseling Move Info Container', () => {
  describe('Basic rendering', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      renderSCMoveInfo();

      const h2 = screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('should render the tab container with two tabs, move details and move history', async () => {
      renderSCMoveInfo();

      expect(screen.getByTestId('MoveDetails-Tab')).toBeInTheDocument();
      expect(screen.getByTestId('MoveTaskOrder-Tab')).toBeInTheDocument();
      expect(screen.getByTestId('MoveHistory-Tab')).toBeInTheDocument();
    });

    it('should render the customer header', async () => {
      renderSCMoveInfo();

      expect(screen.getByRole('heading', { name: 'Kerry, Smith', level: 2 })).toBeInTheDocument();
    });

    it('should render the system error when there is an error', () => {
      renderSCMoveInfo('details', { interceptor: { hasRecentError: true, traceId: 'some-trace-id' } });

      expect(screen.getByText('Technical Help Desk').closest('a')).toHaveAttribute(
        'href',
        'mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil',
      );
      expect(screen.getByTestId('system-error').textContent).toEqual(
        "Something isn't working, but we're not sure what. Wait a minute and try again.If that doesn't fix it, contact the Technical Help Desk (usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil) and give them this code: some-trace-id",
      );
    });

    it('should not render system error when there is not an error', () => {
      renderSCMoveInfo('details', { interceptor: { hasRecentError: false, traceId: '' } });

      expect(queryByTestId(document.documentElement, 'system-error')).not.toBeInTheDocument();
    });
  });

  describe('routing', () => {
    it.each([
      ['Services Counseling Move Details', '/'],
      ['Services Counseling Move Details', 'details'],
      ['Review Documents', 'shipments/SHIP123/document-review'],
      ['Services Counseling Add Shipment', 'new-shipment/hhg'],
      ['Move Task Order', 'mto'],
      ['Customer Support Remarks', 'customer-support-remarks'],
      ['Move History', 'history'],
      ['Services Counseling Move Document Wrapper', 'allowances'],
      ['Services Counseling Move Document Wrapper', 'orders'],
      ['Customer Info', 'customer'],
      ['Services Counseling Edit Shipment Details', 'shipments/SHIP123'],
      ['Services Counseling Edit Shipment Details', 'shipments/SHIP123/advance'],
      ['Review Documents', 'shipments/:shipmentId/document-review'],
      ['Review Documents', 'shipments/:shipmentId/view-document'],
      ['Services Counseling Review Shipment Weights', 'review-shipment-weights'],
    ])(
      'should render the %s component when at the route: /counseling/moves/:moveCode/%s',
      async (componentName, route) => {
        // Render the component at the route
        renderSCMoveInfo(route);

        // Wait for loading to finish
        await waitFor(() => expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument());

        // Assert that the mock component is rendered
        await expect(screen.getByText(`Mock ${componentName} Component`)).toBeInTheDocument();
      },
    );
  });
  it('renders a lock icon', async () => {
    renderSCMoveInfo();

    await waitFor(() => {
      const lockIcon = screen.queryByTestId('locked-move-banner');
      expect(lockIcon).toBeInTheDocument();
    });
  });

  it('should render the Supporting Documents component', async () => {
    const componentName = 'Supporting Documents';
    const nestedPath = 'supporting-documents';

    renderSCMoveInfo(nestedPath);

    // Wait for loading to finish
    await waitFor(() => expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument());

    // Assert that the mock component is rendered
    await waitFor(() => {
      expect(screen.getByText(`Mock ${componentName} Component`)).toBeInTheDocument();
    });
  });

  it('should render the PPM closeout About page component', async () => {
    const componentName = 'About your PPM';
    const nestedPath = 'shipments/:shipmentId/about';

    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      error: null,
    });

    renderSCMoveInfo(nestedPath);

    // Wait for loading to finish
    await waitFor(() => expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument());

    // Assert that the mock component is rendered
    await waitFor(() => {
      expect(screen.getByText(`${componentName}`)).toBeInTheDocument();
    });
  });
});
