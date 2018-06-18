import 'raf/polyfill';
import React from 'react';
import { shallow } from 'enzyme';
import { Landing } from '.';

describe('HomePage tests', () => {
  let wrapper;
  describe('when not loggedIn', () => {
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={false} />, div);
      expect(wrapper.find('.usa-grid').length).toEqual(1);
    });
  });
  describe('when  loggedIn', () => {
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={true} />, div);
      expect(wrapper.find('.usa-grid').length).toEqual(1);
    });
  });
});
