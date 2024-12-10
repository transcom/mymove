import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentInfoListSelector from './ShipmentInfoListSelector';

import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';

const info = {
  requestedPickupDate: '2020-03-26',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  secondaryPickupAddress: {
    streetAddress1: '800 S 2nd St',
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
  counselorRemarks: 'counselor approved',
  customerRemarks: 'customer requested',
};

const ppmInfo = {
  ppmShipment: {
    actualMoveDate: null,
    hasRequestedAdvance: true,
    advanceAmountRequested: 598700,
    approvedAt: null,
    createdAt: '2022-04-29T21:48:21.581Z',
    eTag: 'MjAyMi0wNC0yOVQyMTo0ODoyMS41ODE0MzFa',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: true,
    id: 'b6ec215c-2cef-45fe-8d4a-35f445cd4768',
    proGearWeight: 1987,
    reviewedAt: null,
    shipmentId: 'b5c2d9a1-d1e6-485d-9678-8b62deb0d801',
    spouseProGearWeight: 498,
    status: 'SUBMITTED',
    submittedAt: '2022-04-29T21:48:21.573Z',
  },
};

describe('Shipment Info List', () => {
  it.each([
    ['HHG', SHIPMENT_OPTIONS.HHG, 'shipment-info-list'],
    ['NTS-release', SHIPMENT_OPTIONS.NTSR, 'nts-release-shipment-info-list'],
    ['default', SHIPMENT_OPTIONS.HHG, 'shipment-info-list'],
  ])('when the shipment type is %s it selects the %s shipment', async (_, shipmentType, testId) => {
    render(<ShipmentInfoListSelector shipment={info} shipmentType={shipmentType} />);

    expect(await screen.findByTestId(testId)).toBeInTheDocument();
  });

  it('renders a PPM shipment list', async () => {
    render(<ShipmentInfoListSelector shipment={ppmInfo} shipmentType={SHIPMENT_OPTIONS.PPM} />);

    expect(await screen.findByTestId('ppm-shipment-info-list')).toBeInTheDocument();
  });

  it('renders a NTS shipment list', async () => {
    render(<ShipmentInfoListSelector shipment={info} shipmentType={SHIPMENT_OPTIONS.NTS} />);

    expect(await screen.findByTestId('nts-shipment-info-list')).toBeInTheDocument();
  });

  it('renders a Mobile Home shipment list', async () => {
    render(<ShipmentInfoListSelector shipment={info} shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME} />);

    expect(await screen.findByText('Mobile home year')).toBeInTheDocument();
  });

  it('renders a Boat shipment list', async () => {
    render(<ShipmentInfoListSelector shipment={info} shipmentType={SHIPMENT_OPTIONS.BOAT} />);

    expect(await screen.findByText('Boat year')).toBeInTheDocument();
  });

  it('renders a Boat Haul Away shipment list', async () => {
    render(<ShipmentInfoListSelector shipment={info} shipmentType={SHIPMENT_TYPES.BOAT_HAUL_AWAY} />);

    expect(await screen.findByText('Trailer')).toBeInTheDocument();
  });

  it('renders a Boat Tow Away shipment list', async () => {
    render(<ShipmentInfoListSelector shipment={info} shipmentType={SHIPMENT_TYPES.BOAT_TOW_AWAY} />);

    expect(await screen.findByText('Boat make')).toBeInTheDocument();
  });
});
