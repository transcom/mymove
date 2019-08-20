import React from 'react';
import { shallow } from 'enzyme';
import { Code35FormAlert } from './Code35FormAlert';

describe('code 35A Alert component', () => {
  describe('renders an alert', () => {
    let wrapper = shallow(<Code35FormAlert showAlert={true} />);

    it('without crashing', () => {
      expect(wrapper.exists('Alert')).toBe(true);
    });
  });

  describe('does not render an alert', () => {
    let wrapper = shallow(<Code35FormAlert showAlert={false} />);

    it('without crashing', () => {
      expect(wrapper.exists('Alert')).toBe(false);
    });
  });
});
