import React from 'react';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import routeData from 'react-router-dom';
import { shallow } from 'enzyme';

import SignIn from './SignIn';

import { MockRouterProvider, renderWithProviders } from 'testUtils';
import { isBooleanFlagEnabledUnauthenticated } from 'utils/featureFlags';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabledUnauthenticated: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

afterEach(() => {
  jest.resetAllMocks();
  cleanup();
  jest.spyOn(routeData, 'useLocation').mockReturnValue({
    pathname: '/',
    search: '',
    state: null,
  });
});

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'sign-in',
  }),
}));

describe('SignIn tests', () => {
  it('renders without crashing', () => {
    render(
      <MockRouterProvider>
        <SignIn />
      </MockRouterProvider>,
    );
  });

  it('does not render content of error parameter', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '?error=SOME_ERROR',
      state: null,
    });

    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });

    expect(screen.getByText('An error occurred')).toBeInTheDocument();
    expect(
      screen.getByText('There was an error during your last sign in attempt. Please try again.'),
    ).toBeInTheDocument();
    expect(screen.queryByText('SOME_ERROR')).not.toBeInTheDocument();
  });

  it('shows the EULA when the signin button is clicked and hides the EULA when cancel is clicked', () => {
    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });

    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
    screen.getByTestId('signin').click();
    expect(screen.getByTestId('modal')).toBeInTheDocument();
    screen.getByLabelText('Cancel').click();
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('show logout message when hasLoggedOut state is true', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { hasLoggedOut: true },
    });

    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });

    expect(screen.getByText('You have signed out of MilMove')).toBeInTheDocument();
  });

  it('does not show logout message when hasLoggedOut state is false', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { hasLoggedOut: false },
    });

    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });

    expect(screen.queryByText('You have signed out of MilMove')).not.toBeInTheDocument();
  });

  it('show logout message when timedout state is true', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { timedout: true },
    });

    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });
    expect(screen.getByText('You have been logged out due to inactivity.')).toBeInTheDocument();
  });

  it('does not show logout message when timedout state is false', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { timedout: false },
    });

    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });
    expect(screen.queryByText('You have been logged out due to inactivity.')).not.toBeInTheDocument();
  });
  it('renders with the correct page title', () => {
    const div = document.createElement('div');
    shallow(<SignIn />, div);
    expect(document.title).toContain('Sign In');
  });

  it('renders red warning text', () => {
    const context = { siteName: 'TestMove', showLoginWarning: true };
    renderWithProviders(<SignIn />, {
      context,
    });
    expect(
      screen.getByText('Use of this system is by invitation only, following mandatory screening for'),
    ).toBeInTheDocument();
    expect(screen.getByText('DO NOT PROCEED if you have not gone through that')).toBeInTheDocument();
    expect(screen.queryAllByText('Failure to do so will likely result in you having to resubmit your shipment in the'));
    expect(screen.queryAllByText('and could cause a delay in your shipment being moved.'));
  });

  it('shows the EULA when the create account button is clicked and hides the EULA when cancel is clicked', async () => {
    isBooleanFlagEnabledUnauthenticated.mockImplementation(() => Promise.resolve(true));
    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });

    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
    await waitFor(() => {
      screen.getByTestId('createAccount').click();
    });
    expect(screen.getByTestId('modal')).toBeInTheDocument();
    screen.getByLabelText('Cancel').click();
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('renders non-cac error', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { noValidCAC: true },
    });

    const context = { siteName: 'TestMove' };
    renderWithProviders(<SignIn />, {
      context,
    });

    expect(screen.queryByText('If you do not have a CAC do not request your account here')).not.toBeInTheDocument();
  });
});
