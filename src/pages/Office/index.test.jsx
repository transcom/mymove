/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import OfficeApp from './index';

import { officeRoles, roleTypes } from 'constants/userRoles';
import { configureStore } from 'shared/store';
import { mockPage } from 'testUtils';
import { loadPublicSchema, loadInternalSchema } from 'shared/Swagger/ducks';
import { loadUser } from 'store/auth/actions';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

// Mock Redux actions to prevent actual API calls
let mockPath = '/';
jest.mock('react-router-dom', () => {
  const actual = jest.requireActual('react-router-dom');
  return {
    ...actual,
    useLocation: () => ({ pathname: mockPath }),
    useNavigate: () => jest.fn(),
  };
});

jest.mock('store/auth/actions', () => ({
  loadUser: jest.fn(() => async () => {}),
  setActiveRole: jest.fn().mockImplementation(() => ({ type: '' })),
  logOut: jest.fn().mockImplementation(() => ({ type: '' })),
}));

jest.mock('shared/Swagger/ducks', () => ({
  loadInternalSchema: jest.fn(() => async () => {}),
  loadPublicSchema: jest.fn(() => async () => {}),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

jest.mock('services/ghcApi', () => ({
  checkForLockedMovesAndUnlock: jest.fn(() => Promise.resolve()),
}));

// Mock the components that are routed to from the index, ordered the same as the routes in the index file
mockPage('pages/SignIn/SignIn');
mockPage('pages/InvalidPermissions/InvalidPermissions');
mockPage('pages/Office/MoveQueue/MoveQueue');
mockPage('pages/Office/HeadquartersQueues/HeadquartersQueues', 'Headquarters Queues');
mockPage('pages/Office/PaymentRequestQueue/PaymentRequestQueue');
mockPage('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment');
mockPage('pages/Office/ServicesCounselingQueue/ServicesCounselingQueue');
mockPage('pages/Office/ServicesCounselingMoveInfo/ServicesCounselingMoveInfo');
mockPage('pages/Office/AddShipment/AddShipment');
mockPage('pages/Office/EditShipmentDetails/EditShipmentDetails');
mockPage('pages/PrimeUI/MoveTaskOrder/MoveDetails', 'Prime Simulator Move Details');
mockPage('pages/PrimeUI/Shipment/PrimeUIShipmentCreate', 'Prime Simulator Shipment Create');
mockPage('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateAddress', 'Prime Simulator Shipment Update Address');
mockPage('pages/PrimeUI/Shipment/PrimeUIShipmentUpdate', 'Prime Simulator Shipment Update');
mockPage('pages/PrimeUI/CreatePaymentRequest/CreatePaymentRequest', 'Prime Simulator Create Payment Request');
mockPage(
  'pages/PrimeUI/UploadPaymentRequestDocuments/UploadPaymentRequestDocuments',
  'Prime Simulator Upload Payment Request Documents',
);
mockPage(
  'pages/PrimeUI/UploadServiceRequestDocuments/UploadServiceRequestDocuments',
  'Prime Simulator Upload Service Request Documents',
);
mockPage('pages/PrimeUI/CreateServiceItem/CreateServiceItem', 'Prime Simulator Create Service Item');
mockPage('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateReweigh', 'Prime Simulator Shipment Update Reweigh');
mockPage('pages/Office/MoveSearch/MoveSearch', 'Move Search');
mockPage('pages/Office/TXOMoveInfo/TXOMoveInfo', 'TXO Move Info');
mockPage('pages/PrimeUI/AvailableMoves/AvailableMovesQueue', 'Prime Simulator Available Moves Queue');
mockPage('components/NotFound/NotFound');

const { persistor } = configureStore();

afterEach(async () => {
  cleanup();
  await persistor.purge();
});

const defaultState = {
  auth: {
    activeRole: null,
    hasErrored: false,
    hasSucceeded: true,
    isLoading: false,
    isLoggedIn: false,
    underMaintenance: false,
  },
  swaggerInternal: {
    hasErrored: false,
    hasSucceeded: true,
    isLoading: false,
  },
  interceptor: {
    hasRecentError: false,
    traceId: 'mock-trace-id',
  },
  generalState: {
    showLoadingSpinner: false,
    loadingSpinnerMessage: null,
  },
  entities: {
    user: [
      {
        permissions: [],
        privileges: [],
        inactiveRoles: [],
      },
    ],
  },
};

// Render the OfficeApp component with routing and Redux setup for the provided route and role
const renderOfficeApp = ({ state = defaultState, path = '/', role = null } = {}) => {
  mockPath = path;

  const finalState = role
    ? {
        ...defaultState,
        auth: {
          activeRole: role,
          isLoading: false,
          isLoggedIn: true,
        },
        entities: {
          user: [
            {
              id: 'testUser',
              office_user: { id: 'testUser', first_name: 'bob', last_name: 'robert' },
              inactiveRoles: [],
            },
          ],
        },
      }
    : state;

  const mockStore = configureStore({ ...finalState });
  const queryClient = new QueryClient();
  isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
  const userRoles = role ? [{ roleType: role }] : [];
  render(
    <MemoryRouter initialEntries={[path]}>
      <Provider store={mockStore.store} initialState={finalState}>
        <QueryClientProvider client={queryClient}>
          <OfficeApp
            loadInternalSchema={jest.fn()}
            loadPublicSchema={jest.fn()}
            loadUser={jest.fn()}
            hasRecentError={false}
            activeRole={role || null}
            userRoles={userRoles}
            traceId=""
            loginIsLoading={!!role}
            userIsLoggedIn={!!role}
            hqRoleFlag
            gsrRoleFlag
          />
        </QueryClientProvider>
      </Provider>
    </MemoryRouter>,
  );
};

describe('Office App', () => {
  it('renders Sign In page when user is logged out', async () => {
    renderOfficeApp({ path: '/sign-in' });
    await waitFor(() => expect(screen.getByText(/sign in/i)).toBeInTheDocument());
  });
  it('displays Maintenance page when under maintenance is true', async () => {
    const updatedState = {
      ...defaultState,
      auth: {
        ...defaultState.auth,
        underMaintenance: true,
      },
    };
    renderOfficeApp({ state: updatedState });

    await waitFor(() =>
      expect(screen.getByText(/This system is currently undergoing maintenance/i)).toBeInTheDocument(),
    );
  });
  it('shows loading spinner when showLoadingSpinner is true', async () => {
    const updatedState = {
      ...defaultState,
      auth: {
        ...defaultState.auth,
        isLoggedIn: true,
        activeRole: officeRoles.task_ordering_officer,
      },
      generalState: {
        ...defaultState.generalState,
        showLoadingSpinner: true,
        loadingSpinnerMessage: 'Loading...',
      },
      entities: {
        user: [
          {
            id: 'testUser',
            office_user: { id: 'testUser', first_name: 'bob', last_name: 'robert' },
            inactiveRoles: [],
          },
        ],
      },
    };

    renderOfficeApp({ state: updatedState });
    await waitFor(() => {
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
      expect(screen.getByText(/Loading.../i)).toBeInTheDocument();
    });
  });
  it('handles the SignIn URL for not logged in user and fetches inital data', async () => {
    renderOfficeApp();

    expect(screen.getByText('Skip to content')).toBeInTheDocument();
    expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument();
    expect(screen.getByTestId('signin')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByText('Mock Sign In Component')));
    await waitFor(() => {
      expect(loadUser).toHaveBeenCalled();
      expect(loadInternalSchema).toHaveBeenCalled();
      expect(loadPublicSchema).toHaveBeenCalled();
      expect(screen.getByText('Mock Sign In Component'));
    });
  });
  it('renders SystemError when there is a recent error on root paths', async () => {
    const updatedState = {
      ...defaultState,
      auth: {
        ...defaultState.auth,
        isLoggedIn: true,
        activeRole: roleTypes.SERVICES_COUNSELOR,
      },
      entities: {
        user: [
          {
            id: 'testUser',
            office_user: { id: 'testUser', first_name: 'bob', last_name: 'robert' },
            inactiveRoles: [],
          },
        ],
      },
      interceptor: {
        hasRecentError: true,
        traceId: 'trace-id-123',
      },
    };

    renderOfficeApp({ state: updatedState });
    await waitFor(() => {
      expect(screen.getByText(/trace-id-123/)).toBeInTheDocument();
      expect(screen.getByText(/something isn't working, but we're not sure what/i)).toBeInTheDocument();
    });
  });
  describe('logged out routing', () => {
    it('handles the SignIn URL for not logged in user', async () => {
      renderOfficeApp({ path: '/sign-in' });

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Sign In Component')));
    });
    it('handles the Invalid Permissions URL for not logged in user', async () => {
      renderOfficeApp({ path: '/invalid-permissions' });

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });

    it('handles a bad URL for not logged in user', async () => {
      renderOfficeApp({ path: '/bad-path' });

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait to be redirected to the Sign In page
      await waitFor(() => expect(screen.getByText('Mock Sign In Component')));
    });
  });

  describe('logged in routing', () => {
    it('handles the Invalid Permissions URL', async () => {
      renderOfficeApp({ path: '/invalid-permissions', role: roleTypes.TOO });

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader

      await waitFor(() => expect(screen.getByText('Sign out')).toBeInTheDocument()); // Sign Out button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });

    it('renders the 404 component when the route is not found', async () => {
      renderOfficeApp({ path: '/not-a-real-route', role: roleTypes.QAE });

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button
      expect(screen.getByText('Error - 404')).toBeInTheDocument();
      expect(screen.getByText("We can't find the page you're looking for")).toBeInTheDocument();
    });

    it.each([
      ['Move Queue', '/moves/queue', roleTypes.TOO],
      ['Headquarters Queues', '/hq/queues', roleTypes.HQ],
      ['Payment Request Queue', '/invoicing/queue', roleTypes.TIO],
      ['Services Counseling Add Shipment', '/new-shipment/PPM', roleTypes.SERVICES_COUNSELOR],
      ['Services Counseling Queue', '/counseling', roleTypes.SERVICES_COUNSELOR],
      ['Services Counseling Queue', '/PPM-closeout', roleTypes.SERVICES_COUNSELOR],
      ['Services Counseling Move Info', '/counseling/moves/test123/', roleTypes.SERVICES_COUNSELOR],
      ['Edit Shipment Details', '/moves/test123/shipments/ship123', roleTypes.TOO],
      ['Prime Simulator Move Details', '/simulator/moves/test123/details', roleTypes.PRIME_SIMULATOR],
      ['Prime Simulator Shipment Create', '/simulator/moves/test123/shipments/new', roleTypes.PRIME_SIMULATOR],
      [
        'Prime Simulator Shipment Update Address',
        '/simulator/moves/test123/shipments/ship123/addresses/update',
        roleTypes.PRIME_SIMULATOR,
      ],
      ['Prime Simulator Shipment Update', '/simulator/moves/test123/shipments/ship123', roleTypes.PRIME_SIMULATOR],
      [
        'Prime Simulator Create Payment Request',
        '/simulator/moves/test123/payment-requests/new',
        roleTypes.PRIME_SIMULATOR,
      ],
      [
        'Prime Simulator Upload Payment Request Documents',
        '/simulator/moves/test123/payment-requests/req123/upload',
        roleTypes.PRIME_SIMULATOR,
      ],
      [
        'Prime Simulator Create Service Item',
        '/simulator/moves/test123/shipments/ship123/service-items/new',
        roleTypes.PRIME_SIMULATOR,
      ],
      [
        'Prime Simulator Shipment Update Reweigh',
        '/simulator/moves/test123/shipments/ship123/reweigh/req123/update',
        roleTypes.PRIME_SIMULATOR,
      ],
      ['Move Search', '/', roleTypes.QAE],
      ['Move Search', '/qaecsr/search', roleTypes.QAE],
      ['Move Search', '/', roleTypes.GSR, true],
      ['Move Search', '/qaecsr/search', roleTypes.GSR, true],
      ['TXO Move Info', '/moves/move123', roleTypes.TIO],
      ['Payment Request Queue', '/', roleTypes.TIO],
      ['Move Queue', '/', roleTypes.TOO],
      ['Headquarters Queues', '/', roleTypes.HQ],
      ['Services Counseling Queue', '/', roleTypes.SERVICES_COUNSELOR],
      ['Prime Simulator Available Moves Queue', '/', roleTypes.PRIME_SIMULATOR],
      ['Services Counseling Move Info', '/moves/move123/shipments/:shipmentId/advance', roleTypes.TOO],
    ])('renders the %s component at %s as a %s with sufficient permissions', async (component, path, role) => {
      renderOfficeApp({ path, role });

      const delay = new Promise((resolve) => {
        setTimeout(() => resolve(), 100);
      });

      const afterDelay = delay.then(async () => {
        expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
        expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
        expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button
        screen.debug();
        // Wait for lazy load, validate correct component was rendered
        await waitFor(() => expect(screen.getByText(`Mock ${component} Component`)));
      });

      await afterDelay;
    });

    it.each([
      ['Move Queue', '/moves/queue', roleTypes.PRIME_SIMULATOR],
      ['Payment Request Queue', '/invoicing/queue', roleTypes.PRIME_SIMULATOR],
      ['Services Counseling Add Shipment', '/new-shipment/PPM', roleTypes.PRIME_SIMULATOR],
      ['Services Counseling Move Info', '/counseling/moves/test123/', roleTypes.QAE],
      ['Edit Shipment Details', '/moves/test123/shipments/ship123', roleTypes.QAE],
      ['Prime Simulator Move Details', '/simulator/moves/test123/details', roleTypes.QAE],
      ['Prime Simulator Shipment Create', '/simulator/moves/test123/shipments/new', roleTypes.QAE],
      [
        'Prime Simulator Shipment Update Address as QAE',
        '/simulator/moves/test123/shipments/ship123/addresses/update',
        roleTypes.QAE,
      ],
      ['Prime Simulator Shipment Update', '/simulator/moves/test123/shipments/ship123', roleTypes.QAE],
      ['Prime Simulator Create Payment Request', '/simulator/moves/test123/payment-requests/new', roleTypes.QAE],
      ['Prime Simulator Create Payment Request as QAE', '/simulator/moves/test123/payment-requests/new', roleTypes.QAE],
      [
        'Prime Simulator Upload Payment Request Documents as QAE',
        '/simulator/moves/test123/payment-requests/req123/upload',
        roleTypes.QAE,
      ],
      [
        'Prime Simulator Create Service Item as QAE',
        '/simulator/moves/test123/shipments/ship123/service-items/new',
        roleTypes.QAE,
      ],
      [
        'Prime Simulator Shipment Update Reweigh as QAE',
        '/simulator/moves/test123/shipments/ship123/reweigh/re123/update',
        roleTypes.QAE,
      ],
      ['Services Counseling Move Info as CSR', '/counseling/moves/test123/', roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE],
      ['Edit Shipment Details as CSR', '/moves/test123/shipments/ship123', roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE],
      [
        'Prime Simulator Move Details as CSR',
        '/simulator/moves/test123/details',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      [
        'Prime Simulator Shipment Create as CSR',
        '/simulator/moves/test123/shipments/new',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      [
        'Prime Simulator Shipment Update Address as CSR',
        '/simulator/moves/test123/shipments/ship123/addresses/update',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      [
        'Prime Simulator Shipment Update as CSR',
        '/simulator/moves/test123/shipments/ship123',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      [
        'Prime Simulator Create Payment Request as CSR',
        '/simulator/moves/test123/payment-requests/new',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      [
        'Prime Simulator Upload Payment Request Documents as CSR',
        '/simulator/moves/test123/payment-requests/req123/upload',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      [
        'Prime Simulator Create Service Item as CSR',
        '/simulator/moves/test123/shipments/ship123/service-items/new',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      [
        'Prime Simulator Shipment Update Reweigh as CSR',
        '/simulator/moves/test123/shipments/ship123/reweigh/re123/update',
        roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
      ],
      ['Move Search', '/qaecsr/search', roleTypes.TIO],
      ['TXO Move Info', '/moves/move123', roleTypes.PRIME_SIMULATOR],
    ])('denies access to %s when user has insufficient permission', async (component, path, role) => {
      renderOfficeApp({ path, role });

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button

      // Wait for lazy load, validate invalid permissions component was rendered
      await waitFor(() => {
        expect(screen.queryByText(`Mock ${component} Component`)).not.toBeInTheDocument();
        expect(screen.getByText('Mock Invalid Permissions Component')).toBeInTheDocument();
      });
    });
  });
});
