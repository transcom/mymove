import React from 'react';
import { createMocks } from 'react-idle-timer';
import { act } from 'react-dom/test-utils';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import LogoutOnInactivity from './LogoutOnInactivity';

import { MockProviders } from 'testUtils';

const mockState = (loggedIn) => {
  return {
    auth: {
      isLoggedIn: loggedIn,
    },
  };
};

// These tests assume that the idle time limit >= 1 second
const idleTimeout = 3000;
const warningTime = 1000;

const renderComponent = ({ loggedIn }) => {
  render(
    <MockProviders initialState={mockState(loggedIn)}>
      <LogoutOnInactivity idleTimeout={idleTimeout} warningTime={warningTime} />
    </MockProviders>,
  );
};

const sleep = (time) => {
  return new Promise((resolve) => {
    setTimeout(resolve, time);
  });
};

describe('LogoutOnInactivity', () => {
  describe('when user is logged in', () => {
    beforeAll(() => {
      createMocks();
      global.fetch = jest.fn();
    });

    beforeEach(async () => {
      await act(async () => renderComponent({ loggedIn: true }));
    });

    it('renders without crashing or erroring', async () => {
      const wrapper = screen.getByTestId('logoutOnInactivityWrapper');
      expect(wrapper).toBeInTheDocument();
      expect(screen.queryAllByText('Something')).toHaveLength(0);
    });

    it('becomes idle and triggers a warning', async () => {
      // alert is missing before the user is idle for too long
      expect(
        screen.queryByText('You have been inactive and will be logged out', { exact: false }),
      ).not.toBeInTheDocument();
      await act(async () => {
        return sleep(idleTimeout - warningTime);
      });

      expect(screen.getByText('You have been inactive and will be logged out', { exact: false })).toBeInTheDocument();
    });

    it('removes the idle alert if a user performs a click', async () => {
      // alert is missing before the user is idle for too long
      expect(
        screen.queryByText('You have been inactive and will be logged out', { exact: false }),
      ).not.toBeInTheDocument();
      await act(async () => {
        return sleep(idleTimeout - warningTime);
      });

      // alert is present after user is idle for too long
      expect(screen.getByText('You have been inactive and will be logged out', { exact: false })).toBeInTheDocument();

      const wrapper = screen.getByTestId('logoutOnInactivityWrapper');
      await userEvent.click(wrapper);

      // alert is not present after the click
      expect(
        screen.queryByText('You have been inactive and will be logged out', { exact: false }),
      ).not.toBeInTheDocument();
    });
  });

  describe('when user is not logged in', () => {
    it('does not render the LogoutOnInactivity component', async () => {
      await act(async () => renderComponent({ loggedIn: false }));
      const wrapper = screen.queryByTestId('logoutOnInactivityWrapper');
      expect(wrapper).not.toBeInTheDocument();
    });
  });
});
