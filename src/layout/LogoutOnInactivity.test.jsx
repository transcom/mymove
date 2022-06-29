import React from 'react';
import { createMocks } from 'react-idle-timer';
import { prettyDOM, render, screen } from '@testing-library/react';

import LogoutOnInactivity from './LogoutOnInactivity';

import { MockProviders } from 'testUtils';

const history = {
  goBack: jest.fn(),
  push: jest.fn(),
};

const mockState = {
  auth: {
    isLoggedIn: true,
  },
};

function mountIdleTimer() {
  return render(
    <MockProviders initialState={mockState}>
      <LogoutOnInactivity history={history} />
    </MockProviders>,
  );
}

// function mountIdleTimerWithLoggedOutUser() {
//   return render(
//     <MockProviders initialState={mockState}>
//       <LogoutOnInactivity history={history} />
//     </MockProviders>,
//   );
// }

describe('LogoutOnInactivity', () => {
  describe('component', () => {
    // let component;

    beforeAll(() => {
      createMocks();
    });

    beforeEach(() => {
      jest.useFakeTimers();
      mountIdleTimer();
    });

    it('renders without crashing or erroring', () => {
      const wrapper = screen.getByTestId('logoutOnInactivityWrapper');

      expect(wrapper).toBeDefined();
      expect(screen.queryAllByText('Something')).toHaveLength(0);
    });

    it('becomes idle after 14 minutes', () => {
      const wrapper = screen.getByTestId('logoutOnInactivityWrapper');

      expect(screen.queryByText('You have been inactive')).not.toBeInTheDocument();

      // This line is not working; the state of the document printed immediately afterward shows that the user is not idle
      jest.advanceTimersByTime(14 * 60 * 1000);
      // eslint-disable-next-line no-console
      console.log(prettyDOM(wrapper));

      expect(screen.queryByText('You have been inactive')).toBeInTheDocument();
    });

    // describe('when user is idle', () => {
    //   it('renders the inactivity alert', () => {
    //     component.setState({ isIdle: true });

    //     expect(component.find('Alert').exists()).toBe(true);
    //   });
    // });

    // describe('when user is not idle', () => {
    //   it('does not render the inactivity alert', () => {
    //     component.setState({ isIdle: false });

    //     expect(component.find('Alert').exists()).toBe(false);
    //   });
    // });
  });

  // describe('when user is not logged in', () => {
  //   it('does not render the IdleTimer component', () => {
  //     const component = mountIdleTimerWithLoggedOutUser();
  //     const timer = component.find('IdleTimer');
  //     const idleTimer = component.find(IdleTimer);

  //     expect(timer.exists()).toBe(false);
  //     expect(idleTimer.exists()).toBe(false);
  //   });
  // });
});
