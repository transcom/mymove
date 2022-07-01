import React from 'react';
import { createMocks } from 'react-idle-timer';
import { act } from 'react-dom/test-utils';
import { render, screen } from '@testing-library/react';
import { userEvent } from '@storybook/testing-library';

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
const idleTimeLimitSeconds = 1;
const warningTimeLimitSeconds = 1;

const renderComponent = ({ loggedIn }) => {
  render(
    <MockProviders initialState={mockState(loggedIn)}>
      <LogoutOnInactivity
        maxIdleTimeInSeconds={idleTimeLimitSeconds}
        maxWarningTimeInSeconds={warningTimeLimitSeconds}
      />
    </MockProviders>,
  );
};

const sleep = (seconds) => {
  return new Promise((resolve) => {
    setTimeout(resolve, seconds * 1000);
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
      // alert is missing before the user is idle for the timeout duration
      expect(screen.queryByTestId('logoutAlert')).not.toBeInTheDocument();
      await act(async () => {
        return sleep(idleTimeLimitSeconds);
      });

      expect(screen.queryByTestId('logoutAlert')).toBeInTheDocument();
    });

    it('removes the idle alert if a user performs a click', async () => {
      // alert is missing before the user is idle for the timeout duration
      expect(screen.queryByTestId('logoutAlert')).not.toBeInTheDocument();
      await act(async () => {
        return sleep(idleTimeLimitSeconds);
      });

      const wrapper = screen.getByTestId('logoutOnInactivityWrapper');
      userEvent.click(wrapper);

      // alert is not present after the click
      expect(screen.queryByTestId('logoutAlert')).not.toBeInTheDocument();
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
