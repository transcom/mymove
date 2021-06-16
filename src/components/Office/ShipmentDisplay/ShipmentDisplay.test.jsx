import React from 'react';
import { mount, shallow } from 'enzyme';
import { render, screen } from '@testing-library/react';

import ShipmentDisplay from './ShipmentDisplay';

const info = {
  heading: 'HHG',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
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
  secondaryDeliveryAddress: {
    street_address_1: '987 Fairway Dr',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
  counselorRemarks: 'counselor approved',
};

const secondaryPickupAddressInfo = {
  secondaryPickupAddress: {
    street_address_1: '800 S 2nd St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  ...info,
};

const postalOnly = {
  heading: 'HHG',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    postal_code: '98421',
  },
};

const diversion = {
  heading: 'HHG',
  isDiversion: true,
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
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
  counselorRemarks: 'counselor approved',
};

const cancelled = {
  heading: 'HHG',
  isDiversion: false,
  isCancelled: true,
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
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
  counselorRemarks: 'counselor approved',
};

describe('Shipment Container', () => {
  it('renders the container successfully', () => {
    const wrapper = shallow(
      <ShipmentDisplay shipmentId="1" displayInfo={info} onChange={jest.fn()} isSubmitted={false} />,
    );
    expect(wrapper.find('div[data-testid="shipment-display"]').exists()).toBe(true);
  });
  it('renders secondary address info when present', () => {
    render(
      <ShipmentDisplay
        shipmentId="1"
        displayInfo={secondaryPickupAddressInfo}
        onChange={jest.fn()}
        isSubmitted={false}
      />,
    );
    expect(screen.getByText('Second pickup address')).toBeInTheDocument();
  });
  it('renders with postal only address', () => {
    const wrapper = mount(
      <ShipmentDisplay shipmentId="1" displayInfo={postalOnly} onChange={jest.fn()} isSubmitted={false} />,
    );
    expect(wrapper.find('div[data-testid="shipment-display"]').exists()).toBe(true);
  });
  it('renders with comments', () => {
    render(<ShipmentDisplay shipmentId="1" displayInfo={info} onChange={jest.fn()} isSubmitted={false} />);
    expect(screen.getByText('Counselor remarks')).toBeInTheDocument();
  });
  it('renders with edit button', () => {
    render(<ShipmentDisplay shipmentId="1" displayInfo={info} onChange={jest.fn()} isSubmitted={false} editURL="/" />);
    expect(screen.getByRole('button', { name: 'Edit shipment' })).toBeInTheDocument();
  });
  it('renders without edit button', () => {
    render(<ShipmentDisplay shipmentId="1" displayInfo={info} onChange={jest.fn()} isSubmitted={false} />);
    expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
  });
  it('renders with diversion tag', () => {
    render(<ShipmentDisplay shipmentId="1" displayInfo={diversion} onChange={jest.fn()} isSubmitted={false} />);
    expect(screen.getByText('diversion')).toBeInTheDocument();
  });
  it('renders with cancelled tag', () => {
    render(<ShipmentDisplay shipmentId="1" displayInfo={cancelled} onChange={jest.fn()} isSubmitted={false} />);
    expect(screen.getByText('cancelled')).toBeInTheDocument();
  });
});
