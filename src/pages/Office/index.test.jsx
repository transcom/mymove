/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';

import OfficeApp from './index';

import { configureStore } from 'shared/store';
import { mockPage } from 'testUtils';
import { loadPublicSchema, loadInternalSchema } from 'shared/Swagger/ducks';
import { loadUser } from 'store/auth/actions';

afterEach(() => {
  cleanup();
});

let mockPath = '/';

jest.mock('react-router-dom', () => {
  const actual = jest.requireActual('react-router-dom');
  return {
    ...actual,
    useLocation: () => ({ pathname: mockPath }),
  };
});

// Mock Redux actions to prevent actual API calls
jest.mock('store/auth/actions', () => ({
  loadUser: jest.fn(() => async () => {}),
}));

jest.mock('store/onboarding/actions', () => ({
  initOnboarding: jest.fn(() => async () => {}),
}));

jest.mock('shared/Swagger/ducks', () => ({
  loadInternalSchema: jest.fn(() => async () => {}),
  loadPublicSchema: jest.fn(() => async () => {}),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
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
mockPage('pages/Office/QAECSRMoveSearch/QAECSRMoveSearch', 'QAE CSR Move Search');
mockPage('pages/Office/TXOMoveInfo/TXOMoveInfo', 'TXO Move Info');
mockPage('pages/PrimeUI/AvailableMoves/AvailableMovesQueue', 'Prime Simulator Available Moves Queue');
mockPage('components/NotFound/NotFound');

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
    user: {
      userId123: {
        permissions: [],
        privileges: [],
        roles: [],
      },
    },
  },
};

const loggedInState = {
  ...defaultState,
  auth: {
    ...defaultState.auth,
    isLoggedIn: true,
    activeRole: 'TOO',
    hasError: true,
  },
  entities: {
    user: {
      testUser: {
        id: 'testUser',
        roles: [{ roleType: 'TOO' }],
      },
    },
  },
};

const renderWithState = (state, path) => {
  mockPath = path;
  const mockStore = configureStore({ ...state });

  const minProps = {
    initOnboarding: jest.fn(),
    loadInternalSchema: jest.fn(),
    loadUser: jest.fn(),
  };

  return render(
    <MemoryRouter initialEntries={[path]}>
      <Provider store={mockStore.store}>
        <OfficeApp {...minProps} />
      </Provider>
    </MemoryRouter>,
  );
};

describe('Office App', () => {
  it('renders Sign In page when user is logged out', async () => {
    renderWithState(defaultState, '/sign-in');
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
    renderWithState(updatedState, '/');

    await waitFor(() =>
      expect(screen.getByText(/This system is currently undergoing maintenance/i)).toBeInTheDocument(),
    );
  });
  it('shows loading spinner when showLoadingSpinner is true', async () => {
    const updatedState = {
      ...loggedInState,
      generalState: {
        ...loggedInState.generalState,
        showLoadingSpinner: true,
        loadingSpinnerMessage: 'Loading...',
      },
    };

    renderWithState(updatedState, '/');
    await waitFor(() => {
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
      expect(screen.getByText(/Loading.../i)).toBeInTheDocument();
    });
  });
  it('handles the SignIn URL for not logged in user and fetches inital data', async () => {
    renderWithState(defaultState, '/');

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
  it('renders SystemError when there is a recent error on root path', async () => {
    const updatedState = {
      ...defaultState,
      auth: {
        ...defaultState.auth,
        isLoggedIn: true,
        activeRole: null,
      },
      interceptor: {
        hasRecentError: true,
        traceId: 'trace-id-123',
      },
    };

    renderWithState(updatedState, '/');
    await waitFor(() => {
      expect(screen.getByText(/trace-id-123/)).toBeInTheDocument();
      expect(screen.getByText(/something isn't working, but we're not sure what/i)).toBeInTheDocument();
    });
  });
  it('handles the Invalid Permissions URL for logged in user', async () => {
    renderWithState(loggedInState, '/invalid-permissions');

    expect(screen.getByText('Skip to content')).toBeInTheDocument();
    expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument();

    expect(screen.getByText(/You do not have permission to access this site/i)).toBeInTheDocument();
  });
});
