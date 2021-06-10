import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

const info = {
  agents: [
    {
      agentType: 'RELEASING_AGENT',
      firstName: 'Quinn',
      lastName: 'Ocampo',
      phone: '999-999-9999',
      email: 'quinnocampo@myemail.com',
    },
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  secondaryAddresses: {
    secondaryPickupAddress: {
      street_address_1: '444 S 131st St',
      city: 'San Antonio',
      state: 'TX',
      postal_code: '78234',
    },
    secondaryDeliveryAddress: {
      street_address_1: '7 Q St',
      city: 'Austin',
      state: 'TX',
      postal_code: '78722',
    },
  },
};

const headers = ['Releasing agent', 'Receiving agent', 'Secondary addresses', 'Pickup', 'Destination'];

describe('Shipment Details Sidebar', () => {
  it('renders all fields when provided', () => {
    render(<ShipmentDetailsSidebar {...info} />);

    headers.forEach((header) => {
      expect(screen.getByText(header)).toBeInTheDocument();
    });

    expect(screen.getByText(info.agents[0].email)).toBeInTheDocument();
    expect(screen.getByText(info.agents[1].email)).toBeInTheDocument();
    expect(
      screen.getByText(info.secondaryAddresses.secondaryPickupAddress.street_address_1, { exact: false }),
    ).toBeInTheDocument();
    expect(
      screen.getByText(info.secondaryAddresses.secondaryDeliveryAddress.street_address_1, { exact: false }),
    ).toBeInTheDocument();
  });

  it('renders nothing with no info passed in', () => {
    render(<ShipmentDetailsSidebar />);

    headers.forEach((header) => {
      expect(screen.queryByText(header)).toBeNull();
    });
  });
});
