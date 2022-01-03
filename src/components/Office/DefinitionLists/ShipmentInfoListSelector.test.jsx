import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentInfoListSelector from './ShipmentInfoListSelector';

import { SHIPMENT_OPTIONS } from 'shared/constants';

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

describe('Shipment Info List', () => {
  it.each([
    ['HHG', SHIPMENT_OPTIONS.HHG, 'shipment-info-list'],
    ['NTS-release', SHIPMENT_OPTIONS.NTSR, 'nts-release-shipment-info-list'],
    ['default', SHIPMENT_OPTIONS.HHG, 'shipment-info-list'],
  ])('when the shipment type is %s it selects the %s shipment', async (_, shipmentType, testId) => {
    render(<ShipmentInfoListSelector shipment={info} shipmentType={shipmentType} />);

    expect(await screen.findByTestId(testId)).toBeInTheDocument();
  });
});
