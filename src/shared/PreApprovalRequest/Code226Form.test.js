import React from 'react';
import { shallow } from 'enzyme';
import { Code226Form } from './Code226Form';

let wrapper;
describe('code 35A details component', () => {
  describe('renders', () => {
    wrapper = shallow(<Code226Form ship_line_item_schema={{}} />);

    it('without crashing', () => {
      expect(wrapper.exists('SwaggerField')).toBe(true);
    });
  });
});
