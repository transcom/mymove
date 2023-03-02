/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';
import { Provider } from 'react-redux';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';

import { OfficeApp } from './index';

import { roleTypes } from 'constants/userRoles';
import { configureStore } from 'shared/store';

const mockPage = (path, name) => {
  return jest.mock(path, () => {
    // Create component name from path, if not provided (e.g. 'MoveQueue' -> 'Move Queue')
    const componentName =
      name ||
      path
        .substring(path.lastIndexOf('/') + 1)
        .replace(/([A-Z])/g, ' $1')
        .trim();

    return () => <div>{`Mock ${componentName} Component`}</div>;
  });
};

// Mock the components that are routed to from the index, ordered the same as the routes in the index file
mockPage('pages/SignIn/SignIn');
mockPage('pages/InvalidPermissions/InvalidPermissions');
mockPage('pages/Office/MoveQueue/MoveQueue');
mockPage('pages/Office/PaymentRequestQueue/PaymentRequestQueue');
mockPage('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment');
mockPage('pages/Office/ServicesCounselingQueue/ServicesCounselingQueue');
mockPage('pages/Office/ServicesCounselingMoveInfo/ServicesCounselingMoveInfo');
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
mockPage('pages/PrimeUI/CreateServiceItem/CreateServiceItem', 'Prime Simulator Create Service Item');
mockPage('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateReweigh', 'Prime Simulator Shipment Update Reweigh');
mockPage('pages/Office/QAECSRMoveSearch/QAECSRMoveSearch', 'QAE CSR Move Search');
mockPage('pages/Office/TXOMoveInfo/TXOMoveInfo', 'TXO Move Info');
mockPage('pages/PrimeUI/AvailableMoves/AvailableMovesQueue', 'Prime Simulator Available Moves Queue');
mockPage('components/NotFound/NotFound');

afterEach(() => {
  cleanup();
  jest.resetAllMocks();
});

const loggedInTOOState = {
  auth: {
    activeRole: roleTypes.TOO,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId123: {
        id: 'userId123',
        roles: [{ roleType: roleTypes.TOO }],
      },
    },
  },
};

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

const loggedInSCState = {
  auth: {
    activeRole: roleTypes.SERVICES_COUNSELOR,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId345: {
        id: 'userId345',
        roles: [{ roleType: roleTypes.SERVICES_COUNSELOR }],
      },
    },
  },
};

const loggedInPrimeState = {
  auth: {
    activeRole: roleTypes.PRIME_SIMULATOR,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId456: {
        id: 'userId456',
        roles: [{ roleType: roleTypes.PRIME_SIMULATOR }],
      },
    },
  },
};

const loggedInQAEState = {
  auth: {
    activeRole: roleTypes.QAE_CSR,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId567: {
        id: 'userId567',
        roles: [{ roleType: roleTypes.QAE_CSR }],
      },
    },
  },
};

const loggedOutState = {
  auth: {
    activeRole: null,
    isLoading: false,
    isLoggedIn: false,
  },
};

describe('Office App', () => {
  const mockOfficeProps = {
    loadUser: jest.fn(),
    loadInternalSchema: jest.fn(),
    loadPublicSchema: jest.fn(),
    logOut: jest.fn(),
    hasRecentError: false,
    traceId: '',
  };

  describe('component', () => {
    let wrapper;

    beforeEach(() => {
      wrapper = shallow(<OfficeApp {...mockOfficeProps} router={{ location: { pathname: '/' } }} />);
    });

    it('renders without crashing or erroring', () => {
      const officeWrapper = wrapper.find('div');
      expect(officeWrapper).toBeDefined();
      expect(wrapper.find('SomethingWentWrong')).toHaveLength(0);
    });

    it('renders the logged out header by default', () => {
      expect(wrapper.find('LoggedOutHeader')).toHaveLength(1);
    });

    it('fetches initial data', () => {
      expect(mockOfficeProps.loadUser).toHaveBeenCalled();
      expect(mockOfficeProps.loadInternalSchema).toHaveBeenCalled();
      expect(mockOfficeProps.loadPublicSchema).toHaveBeenCalled();
    });

    describe('if an error occurs', () => {
      it('renders the fail whale', () => {
        wrapper.setState({ hasError: true });
        expect(wrapper.find('SomethingWentWrong')).toHaveLength(1);
      });
    });
  });

  describe('logged out routing', () => {
    it('handles the SignIn URL for not logged in user', async () => {
      const mockStore = configureStore(loggedOutState);
      render(
        <MemoryRouter initialEntries={['/sign-in']}>
          <Provider store={mockStore.store}>
            <OfficeApp
              router={{ location: { pathname: '/sign-in' } }}
              loadInternalSchema={jest.fn()}
              loadPublicSchema={jest.fn()}
              loadUser={jest.fn()}
              hasRecentError={false}
              traceId=""
            />
          </Provider>
        </MemoryRouter>,
      );

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Sign In Component')));
    });

    it('handles the Invalid Permissions URL for not logged in user', async () => {
      const mockStore = configureStore(loggedOutState);
      render(
        <MemoryRouter initialEntries={['/invalid-permissions']}>
          <Provider store={mockStore.store}>
            <OfficeApp
              router={{ location: { pathname: '/invalid-permissions' } }}
              loadInternalSchema={jest.fn()}
              loadPublicSchema={jest.fn()}
              loadUser={jest.fn()}
              hasRecentError={false}
              traceId=""
            />
          </Provider>
        </MemoryRouter>,
      );

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });

    it('handles a bad URL for not logged in user', async () => {
      const mockStore = configureStore(loggedOutState);
      render(
        <MemoryRouter initialEntries={['/bad-path']}>
          <Provider store={mockStore.store}>
            <OfficeApp
              router={{ location: { pathname: '/bad-path' } }}
              loadInternalSchema={jest.fn()}
              loadPublicSchema={jest.fn()}
              loadUser={jest.fn()}
              hasRecentError={false}
              userRoles={[]}
              traceId=""
              loginIsLoading={false}
              userIsLoggedIn={false}
            />
          </Provider>
        </MemoryRouter>,
      );

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
      const mockStore = configureStore(loggedInTOOState);
      render(
        <MemoryRouter initialEntries={['/invalid-permissions']}>
          <Provider store={mockStore.store}>
            <OfficeApp
              router={{ location: { pathname: '/invalid-permissions' } }}
              loadInternalSchema={jest.fn()}
              loadPublicSchema={jest.fn()}
              loadUser={jest.fn()}
              hasRecentError={false}
              userRoles={[{ roleType: loggedInTOOState.auth.activeRole }]}
              traceId=""
              loginIsLoading={false}
              userIsLoggedIn
            />
          </Provider>
        </MemoryRouter>,
      );

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });

    it.each([
      ['Move Queue', '/moves/queue', loggedInTOOState],
      ['Payment Request Queue', '/invoicing/queue', loggedInTIOState],
      ['Services Counseling Add Shipment', '/new-PPM', loggedInSCState],
      ['Services Counseling Queue', '/counseling', loggedInSCState],
      ['Services Counseling Queue', '/PPM-closeout', loggedInSCState],
      ['Services Counseling Move Info', '/counseling/moves/test123/', loggedInSCState],
      ['Edit Shipment Details', '/moves/test123/shipments/ship123', loggedInTOOState],
      ['Prime Simulator Move Details', '/simulator/moves/test123/details', loggedInPrimeState],
      ['Prime Simulator Shipment Create', '/simulator/moves/test123/shipments/new', loggedInPrimeState],
      [
        'Prime Simulator Shipment Update Address',
        '/simulator/moves/test123/shipments/ship123/addresses/update',
        loggedInPrimeState,
      ],
      ['Prime Simulator Shipment Update', '/simulator/moves/test123/shipments/ship123', loggedInPrimeState],
      ['Prime Simulator Create Payment Request', '/simulator/moves/test123/payment-requests/new', loggedInPrimeState],
      [
        'Prime Simulator Upload Payment Request Documents',
        '/simulator/moves/test123/payment-requests/req123/upload',
        loggedInPrimeState,
      ],
      [
        'Prime Simulator Create Service Item',
        '/simulator/moves/test123/shipments/ship123/service-items/new',
        loggedInPrimeState,
      ],
      [
        'Prime Simulator Shipment Update Reweigh',
        '/simulator/moves/test123/shipments/ship123/reweigh/re123/update',
        loggedInPrimeState,
      ],
      ['QAE CSR Move Search', '/qaecsr/search', loggedInQAEState],
      ['TXO Move Info', '/moves/move123', loggedInTIOState],
      ['Payment Request Queue', '/', loggedInTIOState],
      ['Move Queue', '/', loggedInTOOState],
      ['Services Counseling Queue', '/', loggedInSCState],
      ['QAE CSR Move Search', '/', loggedInQAEState],
      ['Prime Simulator Available Moves Queue', '/', loggedInPrimeState],
      // ['Not Found', '/this/is/a/bad/path', loggedInTOOState],
    ])(
      'correctly displays the %s component at %s as a user with sufficient permissions',
      async (component, path, initialState) => {
        const mockStore = configureStore(initialState);

        render(
          <MemoryRouter initialEntries={[path]}>
            <Provider store={mockStore.store}>
              <OfficeApp
                router={{ location: { pathname: path } }}
                loadInternalSchema={jest.fn()}
                loadPublicSchema={jest.fn()}
                loadUser={jest.fn()}
                hasRecentError={false}
                activeRole={initialState.auth.activeRole}
                userRoles={[{ roleType: initialState.auth.activeRole }]}
                traceId=""
                loginIsLoading={false}
                userIsLoggedIn
              />
            </Provider>
          </MemoryRouter>,
        );

        // Header content should be rendered
        expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
        expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
        expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button

        // Wait for lazy load, validate correct component was rendered
        await waitFor(() => expect(screen.getByText(`Mock ${component} Component`)));
      },
    );

    it('handles the Move Queue URL with insufficient permission', async () => {
      const mockStore = configureStore(loggedInTIOState);
      render(
        <MemoryRouter initialEntries={['/moves/queue']}>
          <Provider store={mockStore.store}>
            <OfficeApp
              router={{ location: { pathname: '/moves/queue' } }}
              loadInternalSchema={jest.fn()}
              loadPublicSchema={jest.fn()}
              loadUser={jest.fn()}
              hasRecentError={false}
              userRoles={[{ roleType: loggedInTIOState.auth.activeRole }]}
              traceId=""
              loginIsLoading={false}
              userIsLoggedIn
            />
          </Provider>
        </MemoryRouter>,
      );

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button

      // Wait for and lazy load, validate redirected to invalid permissions page
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });
  });
});
