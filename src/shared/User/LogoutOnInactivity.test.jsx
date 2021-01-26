/* eslint-disable no-only-tests/no-only-tests */
/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';
import { LogoutOnInactivity } from './LogoutOnInactivity';

describe('LogoutOnInactivity', () => {
  const defaultProps = {
    isLoggedIn: true,
  };

  describe('component', () => {
    let wrapper;

    beforeEach(() => {
      wrapper = shallow(<LogoutOnInactivity {...defaultProps} />);
    });

    it('renders without crashing or erroring', () => {
      const timerWrapper = wrapper.find('IdleTimer');
      expect(timerWrapper).toBeDefined();
      expect(wrapper.find('SomethingWentWrong')).toHaveLength(0);
    });

    describe('when user is idle', () => {
      it('renders the inactivity alert', () => {
        wrapper.setState({ isIdle: true });
        expect(wrapper.find('Alert').exists()).toBe(true);
      });
    });

    describe('when user is not idle', () => {
      it('does not render the inactivity alert', () => {
        wrapper.setState({ isIdle: false });
        expect(wrapper.find('Alert').exists()).toBe(false);
      });
    });
  });
});
