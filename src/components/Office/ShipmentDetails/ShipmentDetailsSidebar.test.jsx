import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

import { MockProviders } from 'testUtils';
import { LOA_TYPE } from 'shared/constants';
import { formatAccountingCode } from 'utils/shipmentDisplay';

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
  tacType: LOA_TYPE.HHG,
  sacType: LOA_TYPE.NTS,
};

const ordersLOA = {
  tac: '1234',
  sac: '567',
  ntsTac: '8912',
  ntsSac: '345',
};

const headers = [
  'Releasing agent',
  'Receiving agent',
  'Service order number',
  'Facility info and address',
  'Secondary addresses',
  'Pickup',
  'Destination',
  'Accounting codes',
];

describe('Shipment Details Sidebar', () => {
  it('renders all fields when provided', () => {
    render(
      <MockProviders>
        <ShipmentDetailsSidebar shipment={shipment} ordersLOA={ordersLOA} />
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
    expect(screen.getByText(formatAccountingCode(ordersLOA.tac, shipment.tacType))).toBeInTheDocument();
    expect(screen.getByText(formatAccountingCode(ordersLOA.ntsSac, shipment.sacType))).toBeInTheDocument();
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
