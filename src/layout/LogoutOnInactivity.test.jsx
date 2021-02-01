import React from 'react';
import { mount } from 'enzyme';
import IdleTimer from 'react-idle-timer';

import { LogoutOnInactivity } from './LogoutOnInactivity';

function mountIdleTimer() {
  return mount(<LogoutOnInactivity isLoggedIn />);
}

function mountIdleTimerWithLoggedOutUser() {
  return mount(<LogoutOnInactivity />);
}

describe('LogoutOnInactivity', () => {
  describe('component', () => {
    let component;

    beforeEach(() => {
      jest.useFakeTimers();
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

      jest.advanceTimersByTime(14 * 60 * 1000);

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
      const component = mountIdleTimerWithLoggedOutUser();
      const timer = component.find('IdleTimer');
      const idleTimer = component.find(IdleTimer);

      expect(timer.exists()).toBe(false);
      expect(idleTimer.exists()).toBe(false);
    });
  });
});
