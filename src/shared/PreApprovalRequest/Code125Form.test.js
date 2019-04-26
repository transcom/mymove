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
      expect(wrapper.exists('SwaggerField')).toBe(true);
    });

    it('contains address with correct props', () => {
      expect(wrapper.exists('AddressElementEdit[fieldName="address"]')).toBe(true);
      expect(wrapper.find('AddressElementEdit').prop('zipPattern')).toEqual('USA');
      expect(wrapper.find('AddressElementEdit').prop('title')).toEqual('Truck-to-truck transfer location');
    });

    it('date is required', () => {
      expect(wrapper.find('SwaggerField[fieldName="date"]').prop('required')).toBe(true);
    });

    it('time is optional', () => {
      expect(wrapper.find('SwaggerField[fieldName="time"]').prop('required')).toBeFalsy();
    });
  });
});
