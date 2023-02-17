import React from 'react';
import { cleanup, render, screen } from '@testing-library/react';
import routeData from 'react-router-dom';

import SignIn from './SignIn';

import { MockRouting } from 'testUtils';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => jest.fn(),
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

describe('SignIn tests', () => {
  it('renders without crashing', () => {
    render(
      <MockRouting>
        <SignIn />
      </MockRouting>,
    );
  });

  it('does not render content of error parameter', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '?error=SOME_ERROR',
      state: null,
    });

    const context = { siteName: 'TestMove' };
    render(
      <MockRouting>
        <SignIn context={context} />
      </MockRouting>,
    );

    expect(screen.getByText('An error occurred')).toBeInTheDocument();
    expect(
      screen.getByText('There was an error during your last sign in attempt. Please try again.'),
    ).toBeInTheDocument();
    expect(screen.queryByText('SOME_ERROR')).not.toBeInTheDocument();
  });

  it('shows the EULA when the signin button is clicked and hides the EULA when cancel is clicked', () => {
    const context = { siteName: 'TestMove' };
    render(
      <MockRouting>
        <SignIn context={context} />
      </MockRouting>,
    );

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
    render(
      <MockRouting>
        <SignIn context={context} />
      </MockRouting>,
    );

    expect(screen.getByText('You have signed out of MilMove')).toBeInTheDocument();
  });

  it('does not show logout message when hasLoggedOut state is false', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { hasLoggedOut: false },
    });

    const context = { siteName: 'TestMove' };
    render(
      <MockRouting>
        <SignIn context={context} />
      </MockRouting>,
    );

    expect(screen.queryByText('You have signed out of MilMove')).not.toBeInTheDocument();
  });

  it('show logout message when timedout state is true', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { timedout: true },
    });

    const context = { siteName: 'TestMove' };
    render(
      <MockRouting>
        <SignIn context={context} />
      </MockRouting>,
    );
    expect(screen.getByText('You have been logged out due to inactivity.')).toBeInTheDocument();
  });

  it('does not show logout message when timedout state is false', () => {
    jest.spyOn(routeData, 'useLocation').mockReturnValue({
      pathname: '/sign-in',
      search: '',
      state: { timedout: false },
    });

    const context = { siteName: 'TestMove' };
    render(
      <MockRouting>
        <SignIn context={context} />
      </MockRouting>,
    );
    expect(screen.queryByText('You have been logged out due to inactivity.')).not.toBeInTheDocument();
  });
});
