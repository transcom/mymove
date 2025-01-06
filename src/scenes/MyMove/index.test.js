import React, { Suspense } from 'react';
import { shallow } from 'enzyme';
import { Provider } from 'react-redux';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';

import { CustomerApp } from './index';

import Footer from 'components/Customer/Footer';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { configureStore } from 'shared/store';
import { roleTypes } from 'constants/userRoles';
import { mockPage } from 'testUtils';
import OktaErrorBanner from 'components/OktaErrorBanner/OktaErrorBanner';

// Mock lazy loaded pages
mockPage('pages/SignIn/SignIn');
mockPage('pages/InvalidPermissions/InvalidPermissions');

// Mock imported components
jest.mock('components/Statements/AccessibilityStatement', () => () => (
  <div>Mock Accessibility Statement Component</div>
));
jest.mock('components/Statements/PrivacyAndPolicyStatement', () => () => (
  <div>Mock Privacy And Policy Statement Component</div>
));

const loggedOutState = {
  auth: {
    activeRole: null,
    isLoading: false,
    isLoggedIn: false,
  },
};

const loggedInCustomerState = {
  auth: {
    activeRole: roleTypes.CUSTOMER,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId567: {
        id: 'userId567',
        roles: [{ roleType: roleTypes.CUSTOMER }],
      },
    },
  },
};

const renderAtRoute = (path = '/', state = {}) => {
  const mockStore = configureStore(state);

  const minProps = {
    initOnboarding: jest.fn(),
    loadInternalSchema: jest.fn(),
    loadUser: jest.fn(),
    context: {
      flags: {
        hhgFlow: false,
      },
    },
  };
  render(
    <MemoryRouter initialEntries={[path]}>
      <Provider store={mockStore.store}>
        <Suspense fallback={<div>Loading...</div>}>
          <CustomerApp {...minProps} />
        </Suspense>
      </Provider>
    </MemoryRouter>,
  );
};

afterEach(() => {
  jest.resetAllMocks();
});

describe('CustomerApp tests', () => {
  describe('Customer App logged out routing', () => {
    it.each([
      ['Sign In', '/sign-in'],
      ['Privacy And Policy Statement', '/privacy-and-security-policy'],
      ['Accessibility Statement', '/accessibility'],
      ['Forbidden', '/forbidden', 'You are forbidden to use this endpoint'],
      ['Server Error', '/server_error', 'We are experiencing an internal server error'],
      ['Invalid Permissions', '/invalid-permissions'],
      ['Sign In', '/badurl'], // Bad urls redirect to sign in when not logged in
    ])('renders the %s component at URL: %s', async (component, path, expected = `Mock ${component} Component`) => {
      renderAtRoute(path, loggedOutState);

      // Header content should be rendered
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('An official website of the United States government')).toBeInTheDocument(); // GovBanner
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText(expected)));
    });
  });

  describe('Customer App logged in routing', () => {
    it.each([
      ['Sign In', '/sign-in'],
      ['Privacy And Policy Statement', '/privacy-and-security-policy'],
      ['Accessibility Statement', '/accessibility'],
      ['Forbidden', '/forbidden', 'You are forbidden to use this endpoint'],
      ['Server Error', '/server_error', 'We are experiencing an internal server error'],
      ['Invalid Permissions', '/invalid-permissions'],
      ['Sign In', '/badurl'], // Bad urls redirect to sign in when not logged in
    ])('renders the %s component at URL: %s', async (component, path, expected = `Mock ${component} Component`) => {
      renderAtRoute(path, loggedInCustomerState);

      // Header content should be rendered
      expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument(); // CUIHeader
      expect(screen.getByText('Skip to content')).toBeInTheDocument(); // BypassBlock
      expect(screen.getByText('An official website of the United States government')).toBeInTheDocument(); // GovBanner
      expect(screen.getByTestId('signin')).toBeInTheDocument(); // Sign In button

      // Wait for and lazy load, validate correct component was rendered
      await waitFor(() => expect(screen.getByText(expected)));
    });
  });

  describe('with GHC/HHG feature flags turned off', () => {
    it('renders without crashing or erroring', () => {
      const mockStore = configureStore(loggedInCustomerState);

      const noFlagProps = {
        initOnboarding: jest.fn(),
        loadInternalSchema: jest.fn(),
        loadUser: jest.fn(),
        context: {
          flags: {
            hhgFlow: false,
            ghcFlow: false,
          },
        },
      };

      render(
        <MemoryRouter initialEntries={['/']}>
          <Provider store={mockStore.store}>
            <Suspense fallback={<div>Loading...</div>}>
              <CustomerApp {...noFlagProps} />
            </Suspense>
          </Provider>
        </MemoryRouter>,
      );
    });

    expect(screen.queryByText('Missing Context')).not.toBeInTheDocument();
    expect(screen.queryByText('Error')).not.toBeInTheDocument();
  });

  describe('with GHC/HHG feature flags turned on', () => {
    it('renders without crashing or erroring', () => {
      const mockStore = configureStore(loggedInCustomerState);

      const flagsOnProps = {
        initOnboarding: jest.fn(),
        loadInternalSchema: jest.fn(),
        loadUser: jest.fn(),
        context: {
          flags: {
            hhgFlow: true,
            ghcFlow: true,
          },
        },
      };

      render(
        <MemoryRouter initialEntries={['/']}>
          <Provider store={mockStore.store}>
            <Suspense fallback={<div>Loading...</div>}>
              <CustomerApp {...flagsOnProps} />
            </Suspense>
          </Provider>
        </MemoryRouter>,
      );
    });

    expect(screen.queryByText('Missing Context')).not.toBeInTheDocument();
    expect(screen.queryByText('Error')).not.toBeInTheDocument();
  });

  describe('Page components', () => {
    let wrapper;

    const minProps = {
      initOnboarding: jest.fn(),
      loadInternalSchema: jest.fn(),
      loadUser: jest.fn(),
      context: {
        flags: {
          hhgFlow: false,
        },
      },
    };

    beforeEach(() => {
      wrapper = shallow(<CustomerApp {...minProps} />);
    });

    it('renders without crashing or erroring', () => {
      const appWrapper = wrapper.find('div');
      expect(appWrapper).toBeDefined();
      expect(wrapper.find(SomethingWentWrong)).toHaveLength(0);
    });

    it('renders LoggedOutHeader component by default', () => {
      expect(wrapper.find('LoggedOutHeader')).toHaveLength(1);
    });

    it('renders Footer component', () => {
      expect(wrapper.find(Footer)).toHaveLength(1);
    });

    it('fetches initial data', () => {
      expect(minProps.loadUser).toHaveBeenCalled();
      expect(minProps.loadInternalSchema).toHaveBeenCalled();
    });

    it('renders the fail whale', () => {
      wrapper.setState({ hasError: true });
      expect(wrapper.find(SomethingWentWrong)).toHaveLength(1);
    });
    it('renders the error banner when redirected from Okta', () => {
      wrapper.setState({ oktaErrorBanner: true });
      expect(wrapper.find(OktaErrorBanner)).toHaveLength(1);
    });
  });

  it('renders the Maintenance page flag is true', async () => {
    const minPropsWithMaintenance = {
      initOnboarding: jest.fn(),
      loadInternalSchema: jest.fn(),
      loadUser: jest.fn(),
      underMaintenance: true,
      context: {
        flags: {
          hhgFlow: false,
        },
      },
    };

    const wrapper = shallow(<CustomerApp {...minPropsWithMaintenance} />);

    // maintenance page should be rendered
    expect(wrapper.find('MaintenancePage')).toHaveLength(1);
  });
});
