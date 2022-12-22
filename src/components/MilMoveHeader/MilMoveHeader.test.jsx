import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import LoggedOutUserInfo from './LoggedOutUserInfo';
import CustomerUserInfo from './CustomerUserInfo';
import OfficeUserInfo from './OfficeUserInfo';

import MilMoveHeader from './index';

import { MockProviders } from 'testUtils';

describe('MilMoveHeader and User Infos', () => {
  it('renders the base header with nothing in it', () => {
    render(<MilMoveHeader />);

    const homeLink = screen.getByTitle('Home');
    expect(homeLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it('renders the Logged Out User Info', async () => {
    const signInHandler = jest.fn();
    render(
      <MilMoveHeader>
        <LoggedOutUserInfo handleLogin={signInHandler} />
      </MilMoveHeader>,
    );
    const signInButton = screen.getByRole('button', { name: 'Sign In' });
    expect(signInButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signInButton);

    await waitFor(() => {
      expect(signInHandler).toHaveBeenCalled();
    });
  });

  it('renders the Customer User Info without the profile icon', async () => {
    const signOutHandler = jest.fn();
    render(
      <MilMoveHeader>
        <CustomerUserInfo handleLogout={signOutHandler} showProfileLink={false} />
      </MilMoveHeader>,
    );

    const signOutButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signOutButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signOutButton);

    await waitFor(() => {
      expect(signOutHandler).toHaveBeenCalled();
    });

    expect(screen.queryByTitle('profile-link')).not.toBeInTheDocument();
  });

  it('renders the Customer User Info with the profile icon', async () => {
    const signOutHandler = jest.fn();
    render(
      <MockProviders>
        <MilMoveHeader>
          <CustomerUserInfo handleLogout={signOutHandler} showProfileLink />
        </MilMoveHeader>
      </MockProviders>,
    );

    const signOutButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signOutButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signOutButton);

    await waitFor(() => {
      expect(signOutHandler).toHaveBeenCalled();
    });

    const profileLink = screen.getByTitle('profile-link');
    expect(profileLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it('renders the Office User Info', async () => {
    const testProps = {
      lastName: 'Baker',
      firstName: 'Riley',
      handleLogout: jest.fn(),
    };

    render(
      <MilMoveHeader>
        <OfficeUserInfo {...testProps} />
      </MilMoveHeader>,
    );

    const signOutButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signOutButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signOutButton);

    await waitFor(() => {
      expect(testProps.handleLogout).toHaveBeenCalled();
    });

    expect(screen.getByText('Baker, Riley'));
  });
});
