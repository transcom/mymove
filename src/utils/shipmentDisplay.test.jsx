import { mount } from 'enzyme';

import { formatAddress, formatCustomerDestination } from './shipmentDisplay';

describe('shipmentDisplay utils', () => {
  describe('formatAddress', () => {
    describe('all address parts provided', () => {
      const shipmentAddress = {
        street_address_1: '555 Main Street',
        city: 'Celebration',
        state: 'FL',
        postal_code: '34747',
      };
      const component = mount(formatAddress(shipmentAddress));
      it('includes full address with comma seperator', () => {
        expect(component.at(0).text()).toEqual('555 Main Street');
        // Must use the character code for nbsp
        expect(component.at(1).text()).toEqual(',\xa0');
        expect(component.at(2).text()).toEqual('Celebration, FL 34747');
      });
    });
    describe('street address is missing', () => {
      const shipmentAddress = {
        city: 'Celebration',
        state: 'FL',
        postal_code: '34747',
      };
      const component = mount(formatAddress(shipmentAddress));
      it('formats as single line', () => {
        expect(component.text()).toEqual('Celebration, FL 34747');
      });
    });
    describe('postal code only', () => {
      const shipmentAddress = {
        postal_code: '34747',
      };
      const component = mount(formatAddress(shipmentAddress));
      it('omits city and state', () => {
        expect(component.text()).toEqual('34747');
      });
    });
  });
  describe('formatDestination', () => {
    it('shows entire address', () => {});

    it('shows postalCode if address is not provided', () => {
      expect(formatCustomerDestination(null, '11111')).toBe('11111');
    });
  });
});
