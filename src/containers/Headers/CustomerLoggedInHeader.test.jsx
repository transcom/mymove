import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedCustomerLoggedInHeader from './CustomerLoggedInHeader';

import { LogoutUser } from 'utils/api';
import { logOut } from 'store/auth/actions';
import { selectIsProfileComplete } from 'store/entities/selectors';
import { MockProviders } from 'testUtils';

jest.mock('store/auth/actions', () => ({
  ...jest.requireActual('store/auth/actions'),
  logOut: jest.fn().mockImplementation(() => ({ type: '' })),
}));

jest.mock('utils/api', () => ({
  LogoutUser: jest.fn(() => ({ then: () => {} })),
}));

jest.mock('store/entities/selectors', () => ({
  selectIsProfileComplete: jest.fn(),
}));

describe('CustomerLoggedInHeader', () => {
  it('renders the customer logged in header', () => {
    render(
      <MockProviders>
        <ConnectedCustomerLoggedInHeader />
      </MockProviders>,
    );

    const homeLink = screen.getByTitle('Home');
    expect(homeLink).toBeInstanceOf(HTMLAnchorElement);

    const signOutButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signOutButton).toBeInstanceOf(HTMLButtonElement);
  });

  it('does not show the profile icon if the profile is not complete', async () => {
    selectIsProfileComplete.mockImplementation(() => false);

    render(
      <MockProviders>
        <ConnectedCustomerLoggedInHeader />
      </MockProviders>,
    );

    expect(screen.queryByTitle('profile-link')).not.toBeInTheDocument();
  });

  it('shows the profile icon if the profile is complete', () => {
    selectIsProfileComplete.mockImplementation(() => true);

    render(
      <MockProviders>
        <ConnectedCustomerLoggedInHeader />
      </MockProviders>,
    );

    const profileLink = screen.getByTitle('profile-link');
    expect(profileLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it('signs out the user when sign out is clicked', async () => {
    render(
      <MockProviders>
        <ConnectedCustomerLoggedInHeader />
      </MockProviders>,
    );

    const signOutButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signOutButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signOutButton);

    expect(logOut).toHaveBeenCalled();
    expect(LogoutUser).toHaveBeenCalled();
  });
});
