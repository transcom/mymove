import React from 'react';
import { shallow } from 'enzyme';
import ShipmentDisplay from './ShipmentDisplay';

const info = {
  heading: 'HHG',
  requestedMoveDate: '26 Mar 2020',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
};

describe('Shipment Container', () => {
  it('renders the container successfully', () => {
    const wrapper = shallow(<ShipmentDisplay displayInfo={info} />);
    expect(wrapper.find('div[data-cy="shipment-display"]').exists()).toBe(true);
  });
  it('renders a checkbox with id passed in', () => {
    const wrapper = shallow(<ShipmentDisplay checkboxId="TESTING123" displayInfo={info} />);
    expect(wrapper.find('input[data-cy="shipment-display-checkbox"]').props().id).toBe('TESTING123');
  });
});
