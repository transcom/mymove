import React from 'react';
import { mount } from 'enzyme';
import { LocationsDisplay } from './Locations';

const defaultProps = {
  shipment: {
    delivery_address: {},
    pickup_address: {},
    has_delivery_address: true,
    has_secondary_pickup_address: false,
    secondary_pickup_address: {},
    service_member: { current_station: { address: {} } },
  },
};
describe('Locations component test', () => {
  describe('Locations display', () => {
    const shipment = {
      ...defaultProps.shipment,
      pickup_address: {
        street_address_1: '123 Disney Rd',
        city: 'Los Angeles',
        state: 'CA',
        postal_code: '90210',
      },
    };
    const wrapper = mount(<LocationsDisplay shipment={shipment} />);
    it('should render 2 headers', () => {
      const headers = wrapper.find('.column-subhead');
      expect(headers.length).toBe(2);
    });
    it('should render the Pickup header', () => {
      const Pickup = wrapper.find('.column-subhead').getElements()[0];
      const { className, children } = Pickup.props;
      expect(className).toBe('column-subhead');
      expect(children).toBe('Pickup');
    });
    it('should render the Delivery header', () => {
      const Delivery = wrapper.find('.column-subhead').getElements()[1];
      const { className, children } = Delivery.props;
      expect(className).toBe('column-subhead');
      expect(children).toBe('Delivery');
    });
    it('should render primary pickup address', () => {
      console.log('what', wrapper.debug());
    });
  });
});
