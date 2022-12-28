import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedOfficeLoggedInHeader from './OfficeLoggedInHeader';

import { LogoutUser } from 'utils/api';
import { logOut } from 'store/auth/actions';
import { MockProviders } from 'testUtils';
import { roleTypes } from 'constants/userRoles';

jest.mock('store/auth/actions', () => ({
  ...jest.requireActual('store/auth/actions'),
  logOut: jest.fn().mockImplementation(() => ({ type: '' })),
}));

jest.mock('utils/api', () => ({
  LogoutUser: jest.fn(() => ({ then: () => {} })),
}));

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
  });
});
