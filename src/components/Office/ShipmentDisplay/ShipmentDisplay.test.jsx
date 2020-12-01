import React from 'react';
import { shallow, mount } from 'enzyme';

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

const postalOnly = {
  heading: 'HHG',
  requestedMoveDate: '26 Mar 2020',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    postal_code: '98421',
  },
};

describe('Shipment Container', () => {
  it('renders the container successfully', () => {
    const wrapper = shallow(
      <ShipmentDisplay shipmentId="1" displayInfo={info} onChange={jest.fn()} isSubmitted={false} />,
    );
    expect(wrapper.find('div[data-testid="shipment-display"]').exists()).toBe(true);
  });
  it('renders with postal only address', () => {
    const wrapper = mount(
      <ShipmentDisplay shipmentId="1" displayInfo={postalOnly} onChange={jest.fn()} isSubmitted={false} />,
    );
    expect(wrapper.find('div[data-testid="shipment-display"]').exists()).toBe(true);
  });
});
