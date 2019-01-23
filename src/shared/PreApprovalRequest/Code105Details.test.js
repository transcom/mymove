import React from 'react';
import { shallow } from 'enzyme';
import { Code105Details } from './Code105Details';

let wrapper;
describe('Renders without crashing', () => {
  describe('Code 105B/E details component', () => {
    wrapper = shallow(<Code105Details />);

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('div')).toBe(true);
    });
  });
});
