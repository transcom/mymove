import { shallow } from 'enzyme/build';
import React from 'react';
import { Code125Form } from './Code125Form';

const simpleSchema = {
  properties: {
    address: {},
  },
};

let wrapper;
describe('code 125 details component', () => {
  describe('renders', () => {
    wrapper = shallow(<Code125Form ship_line_item_schema={simpleSchema} />);

    it('without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('SwaggerField')).toBe(true);
    });

    it('contains address', () => {
      expect(wrapper.exists('AddressElementEdit[fieldName="address"]')).toBe(true);
    });
  });
});
