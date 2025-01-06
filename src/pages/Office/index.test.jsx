/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';
import { Provider } from 'react-redux';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';

import { OfficeApp } from './index';

import { roleTypes } from 'constants/userRoles';
import { configureStore } from 'shared/store';
import { mockPage } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

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

afterEach(() => {
  cleanup();
  jest.resetAllMocks();
});

const createMockStore = (role) => {
  if (!role) {
    // If no role provided, use logged out state
    const loggedOutState = {
      auth: {
        activeRole: null,
        isLoading: false,
        isLoggedIn: false,
      },
    };

    return configureStore(loggedOutState);
  }

  // Otherwise, use logged in state with the provided role
  const state = {
    auth: {
      activeRole: role,
      isLoading: false,
      isLoggedIn: true,
    },
    entities: {
      user: {
        userId123: {
          id: 'userId123',
          roles: [{ roleType: role }],
        },
      },
    },
  };

  return configureStore(state);
};

// Render the OfficeApp component with routing and Redux setup for the provided route and role
const renderOfficeAppAtRoute = (route, role) => {
  isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
  const mockStore = createMockStore(role);
  const userRoles = role ? [{ roleType: role }] : [];
  render(
    <MemoryRouter initialEntries={[route]}>
      <Provider store={mockStore.store}>
        <OfficeApp
          router={{ location: { pathname: route } }}
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
      </Provider>
    </MemoryRouter>,
  );
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
    describe('if redirected from okta', () => {
      it('renders the okta banner following log out', () => {
        wrapper.setState({ oktaLoggedOut: true });
        expect(wrapper.find('OktaLoggedOutBanner')).toHaveLength(1);
      });
    });
    describe('if redirected from okta', () => {
      it('renders the okta banner when not logged out', () => {
        wrapper.setState({ oktaNeedsLoggedOut: true });
        expect(wrapper.find('OktaNeedsLoggedOutBanner')).toHaveLength(1);
      });
    });
  });

  describe('logged out routing', () => {
    it('handles the SignIn URL for not logged in user', async () => {
      renderOfficeAppAtRoute('/sign-in');

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Sign In Component')));
    });

    it('handles the Invalid Permissions URL for not logged in user', async () => {
      renderOfficeAppAtRoute('/invalid-permissions');

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });

    it('handles a bad URL for not logged in user', async () => {
      renderOfficeAppAtRoute('/bad-path');

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
      renderOfficeAppAtRoute('/invalid-permissions', roleTypes.TOO);

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });

    it('renders the 404 component when the route is not found', async () => {
      renderOfficeAppAtRoute('/not-a-real-route', roleTypes.QAE);

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button
      await expect(screen.getByText('Error - 404')).toBeInTheDocument();
      await expect(screen.getByText("We can't find the page you're looking for")).toBeInTheDocument();
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
      ['QAE CSR Move Search', '/', roleTypes.QAE],
      ['QAE CSR Move Search', '/qaecsr/search', roleTypes.QAE],
      ['QAE CSR Move Search', '/', roleTypes.GSR, true],
      ['QAE CSR Move Search', '/qaecsr/search', roleTypes.GSR, true],
      ['TXO Move Info', '/moves/move123', roleTypes.TIO],
      ['Payment Request Queue', '/', roleTypes.TIO],
      ['Move Queue', '/', roleTypes.TOO],
      ['Headquarters Queues', '/', roleTypes.HQ],
      ['Services Counseling Queue', '/', roleTypes.SERVICES_COUNSELOR],
      ['Prime Simulator Available Moves Queue', '/', roleTypes.PRIME_SIMULATOR],
      ['Services Counseling Move Info', '/moves/move123/shipments/:shipmentId/advance', roleTypes.TOO],
    ])('renders the %s component at %s as a %s with sufficient permissions', async (component, path, role) => {
      renderOfficeAppAtRoute(path, role);

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button

      // Wait for lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText(`Mock ${component} Component`)));
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
      ['QAE CSR Move Search', '/qaecsr/search', roleTypes.TIO],
      ['TXO Move Info', '/moves/move123', roleTypes.PRIME_SIMULATOR],
    ])('denies access to %s when user has insufficient permission', async (component, path, role) => {
      renderOfficeAppAtRoute(path, role);

      // Header content should be rendered
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Sign out')).toBeInTheDocument(); // Sign Out button

      // Wait for lazy load, validate invalid permissions component was rendered
      await waitFor(() => expect(screen.getByText('Mock Invalid Permissions Component')));
    });
  });

  it('renders the Maintenance page flag is true', async () => {
    const mockMaintenanceOfficeProps = {
      loadUser: jest.fn(),
      loadInternalSchema: jest.fn(),
      loadPublicSchema: jest.fn(),
      logOut: jest.fn(),
      underMaintenance: true,
      hasRecentError: false,
      traceId: '',
    };

    const wrapper = shallow(<OfficeApp {...mockMaintenanceOfficeProps} router={{ location: { pathname: '/' } }} />);

    // maintenance page should be rendered
    expect(wrapper.find('MaintenancePage')).toHaveLength(1);
  });
});
