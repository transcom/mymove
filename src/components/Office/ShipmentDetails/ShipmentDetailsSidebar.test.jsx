import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

const shipment = {
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
      streetAddress1: '444 S 131st St',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78234',
    },
    secondaryDeliveryAddress: {
      streetAddress1: '7 Q St',
      city: 'Austin',
      state: 'TX',
      postalCode: '78722',
    },
  },
};

const headers = ['Releasing agent', 'Receiving agent', 'Secondary addresses', 'Pickup', 'Destination'];

describe('Shipment Details Sidebar', () => {
  it('renders all fields when provided', () => {
    render(<ShipmentDetailsSidebar shipment={shipment} />);

    headers.forEach((header) => {
      expect(screen.getByText(header)).toBeInTheDocument();
    });

    expect(screen.getByText(shipment.agents[0].email)).toBeInTheDocument();
    expect(screen.getByText(shipment.agents[1].email)).toBeInTheDocument();
    expect(
      screen.getByText(shipment.secondaryAddresses.secondaryPickupAddress.streetAddress1, { exact: false }),
    ).toBeInTheDocument();
    expect(
      screen.getByText(shipment.secondaryAddresses.secondaryDeliveryAddress.streetAddress1, { exact: false }),
    ).toBeInTheDocument();
  });

  it('renders nothing with no info passed in', () => {
    render(<ShipmentDetailsSidebar />);

    headers.forEach((header) => {
      expect(screen.queryByText(header)).toBeNull();
    });
  });
});
