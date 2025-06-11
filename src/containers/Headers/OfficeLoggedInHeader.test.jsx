import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedOfficeLoggedInHeader from './OfficeLoggedInHeader';

import { LogoutUser } from 'utils/api';
import { logOut } from 'store/auth/actions';
import { MockProviders } from 'testUtils';
import { roleTypes } from 'constants/userRoles';
import { checkForLockedMovesAndUnlock } from 'services/ghcApi';

jest.mock('store/auth/actions', () => ({
  ...jest.requireActual('store/auth/actions'),
  logOut: jest.fn().mockImplementation(() => ({ type: '' })),
}));

jest.mock('utils/api', () => ({
  LogoutUser: jest.fn(() => ({ then: () => {} })),
}));

jest.mock('services/ghcApi', () => ({
  checkForLockedMovesAndUnlock: jest.fn(() => Promise.resolve()),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: jest.fn().mockReturnValue({ pathname: '/' }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const localStorageMock = (() => {
  let store = {};

  return {
    getItem(key) {
      return store[key] || null;
    },
    setItem(key, value) {
      store[key] = value.toString();
    },
    removeItem(key) {
      delete store[key];
    },
    clear() {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'sessionStorage', {
  value: localStorageMock,
});

describe('OfficeLoggedInHeader', () => {
  it('renders the office logged in header', () => {
    render(
      <MockProviders>
        <ConnectedOfficeLoggedInHeader />
      </MockProviders>,
    );

    const homeLink = screen.getByTitle('Home');
    expect(homeLink).toBeInstanceOf(HTMLAnchorElement);

    const signInButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signInButton).toBeInstanceOf(HTMLButtonElement);
  });

  it('shows the correct queue link for the TIO', () => {
    const testState = {
      auth: {
        activeRole: roleTypes.TIO,
        isLoading: false,
        isLoggedIn: true,
      },
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            roles: [{ roleType: roleTypes.TIO }],
            office_user: {
              first_name: 'Amanda',
              last_name: 'Gorman',
              transportation_office: {
                gbloc: 'ABCD',
              },
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={testState}>
        <ConnectedOfficeLoggedInHeader />
      </MockProviders>,
    );

    const queueLink = screen.getByText('ABCD payment requests');
    expect(queueLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it('shows the correct queue link for the TOO', () => {
    const testState = {
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
            office_user: {
              first_name: 'Amanda',
              last_name: 'Gorman',
              transportation_office: {
                gbloc: 'ABCD',
              },
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={testState}>
        <ConnectedOfficeLoggedInHeader />
      </MockProviders>,
    );

    const queueLink = screen.getByText('ABCD moves');
    expect(queueLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it('shows the correct queue link for the services counselor', () => {
    const testState = {
      auth: {
        activeRole: roleTypes.SERVICES_COUNSELOR,
        isLoading: false,
        isLoggedIn: true,
      },
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            roles: [{ roleType: roleTypes.SERVICES_COUNSELOR }],
            office_user: {
              first_name: 'Amanda',
              last_name: 'Gorman',
              transportation_office: {
                gbloc: 'ABCD',
              },
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={testState}>
        <ConnectedOfficeLoggedInHeader />
      </MockProviders>,
    );

    const queueLink = screen.getByText('ABCD');
    expect(queueLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it('signs out the user when sign out is clicked', async () => {
    const sessionStorageClearSpy = jest.spyOn(window.sessionStorage, 'clear');
    render(
      <MockProviders>
        <ConnectedOfficeLoggedInHeader />
      </MockProviders>,
    );

    const signOutButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signOutButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signOutButton);

    expect(logOut).toHaveBeenCalled();
    expect(LogoutUser).toHaveBeenCalled();
    expect(sessionStorageClearSpy).toHaveBeenCalled();
  });

  it('renders the GBLOC switcher when the current user is signed in with the HQ role', async () => {
    const testState = {
      auth: {
        activeRole: roleTypes.HQ,
        isLoading: false,
        isLoggedIn: true,
      },
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            roles: [{ roleType: roleTypes.HQ }],
            office_user: {
              first_name: 'Amanda',
              last_name: 'Gorman',
              transportation_office: {
                gbloc: 'KKFA',
              },
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={testState}>
        <ConnectedOfficeLoggedInHeader />
      </MockProviders>,
    );

    const gblocSwitcher = screen.getByTestId('gbloc_switcher');
    expect(gblocSwitcher).toBeInstanceOf(HTMLDivElement);
    expect(gblocSwitcher.firstChild).toBeInstanceOf(HTMLSelectElement);
    expect(gblocSwitcher.firstChild.firstChild).toBeInstanceOf(HTMLOptionElement);
  });

  it('unlocks all moves when the milmove logo is clicked', async () => {
    const testState = {
      auth: {
        activeRole: roleTypes.QAE,
        isLoading: false,
        isLoggedIn: true,
      },
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            roles: [{ roleType: roleTypes.QAE }],
            office_user: {
              first_name: 'Amanda',
              last_name: 'Gorman',
              transportation_office: {
                gbloc: 'KKFA',
              },
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={testState}>
        <ConnectedOfficeLoggedInHeader />
      </MockProviders>,
    );

    expect(checkForLockedMovesAndUnlock).toHaveBeenCalledTimes(1);
  });
});
