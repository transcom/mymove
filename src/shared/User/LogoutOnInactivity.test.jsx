/* eslint-disable no-only-tests/no-only-tests */
/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import sinon from 'sinon';
import IdleTimer from 'react-idle-timer';
import { LogoutOnInactivity } from './LogoutOnInactivity';

const defaultProps = {
  isLoggedIn: true,
};

function mountIdleTimer(props = {}) {
  return mount(<LogoutOnInactivity {...defaultProps} />);
}

const loggedOutProps = {
  isLoggedIn: false,
};

function mountIdleTimerWithLoggedOutUser(props = {}) {
  return mount(<LogoutOnInactivity {...loggedOutProps} />);
}

describe('LogoutOnInactivity', () => {
  describe('component', () => {
    let component;
    let clock;

    beforeEach(() => {
      clock = sinon.useFakeTimers();
      component = mountIdleTimer();
    });

    it('renders without crashing or erroring', () => {
      const timer = component.find('IdleTimer');

      expect(timer).toBeDefined();
      expect(component.find('SomethingWentWrong')).toHaveLength(0);
    });

    it('becomes idle after 14 minutes', () => {
      const idleTimer = component.find(IdleTimer);

      expect(idleTimer.state().idle).toBe(false);

      clock.tick(14 * 1000 * 60);

      expect(component.state().isIdle).toBe(true);
      expect(idleTimer.state().idle).toBe(true);
    });

    describe('when user is idle', () => {
      it('renders the inactivity alert', () => {
        component.setState({ isIdle: true });

        expect(component.find('Alert').exists()).toBe(true);
      });
    });

    describe('when user is not idle', () => {
      it('does not render the inactivity alert', () => {
        component.setState({ isIdle: false });

        expect(component.find('Alert').exists()).toBe(false);
      });
    });
  });

  describe('when user is not logged in', () => {
    it('does not render the IdleTimer component', () => {
      let component;
      component = mountIdleTimerWithLoggedOutUser();
      const timer = component.find('IdleTimer');

      expect(timer.exists()).toBe(false);
    });
  });
});
