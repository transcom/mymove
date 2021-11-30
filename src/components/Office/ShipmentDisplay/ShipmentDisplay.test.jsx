import React from 'react';
import { mount, shallow } from 'enzyme';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentDisplay from './ShipmentDisplay';

const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
}));

const info = {
  heading: 'HHG',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '987 Fairway Dr',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  counselorRemarks: 'counselor approved',
};

const secondaryPickupAddressInfo = {
  secondaryPickupAddress: {
    streetAddress1: '800 S 2nd St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  ...info,
};

const postalOnly = {
  heading: 'HHG',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    postalCode: '98421',
  },
};

const diversion = {
  heading: 'HHG',
  isDiversion: true,
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  counselorRemarks: 'counselor approved',
};

const cancelled = {
  heading: 'HHG',
  isDiversion: false,
  shipmentStatus: 'CANCELED',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
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
  it('renders with edit button', async () => {
    render(<ShipmentDisplay shipmentId="1" displayInfo={info} onChange={jest.fn()} isSubmitted={false} editURL="/" />);

    const button = screen.getByRole('button', { name: 'Edit shipment' });
    expect(button).toBeInTheDocument();
    userEvent.click(button);
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/');
    });
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
