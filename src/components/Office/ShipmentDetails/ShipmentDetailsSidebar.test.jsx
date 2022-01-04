import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

import { MockProviders } from 'testUtils';

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
  storageFacility: {
    facilityName: 'Some storage facility',
    streetAddress1: '456 S 131st St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78212',
  },
  serviceOrderNumber: '1234',
};

const headers = [
  'Releasing agent',
  'Receiving agent',
  'Service order number',
  'Facility info and address',
  'Secondary addresses',
  'Pickup',
  'Destination',
];

describe('Shipment Details Sidebar', () => {
  it('renders all fields when provided', () => {
    render(
      <MockProviders>
        <ShipmentDetailsSidebar shipment={shipment} />
      </MockProviders>,
    );

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
    expect(screen.getByText(shipment.storageFacility.streetAddress1, { exact: false })).toBeInTheDocument();
    expect(screen.getByText(shipment.serviceOrderNumber)).toBeInTheDocument();
  });

  it('renders nothing with no info passed in', () => {
    render(
      <MockProviders>
        <ShipmentDetailsSidebar />
      </MockProviders>,
    );

    headers.forEach((header) => {
      expect(screen.queryByText(header)).toBeNull();
    });
  });
});
