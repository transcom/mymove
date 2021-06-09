import React from 'react';
import { render, screen, within } from '@testing-library/react';

import ShipmentInfoList from './ShipmentInfoList';

const info = {
  requestedMoveDate: '2020-03-26',
  originAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  secondPickupAddress: {
    street_address_1: '800 S 2nd St',
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
  secondDestinationAddress: {
    street_address_1: '987 Fairway Dr',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
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

const labels = {
  requestedMoveDate: 'Requested move date',
  originAddress: 'Origin address',
  secondPickupAddress: 'Second pickup address',
  destinationAddress: 'Destination address',
  secondDestinationAddress: 'Second destination address',
  agents: ['Releasing agent', 'Receiving agent'],
  counselorRemarks: 'Counselor remarks',
  customerRemarks: 'Customer remarks',
};

describe('Shipment Info List', () => {
  it('renders all fields when provided', () => {
    render(<ShipmentInfoList {...info} />);

    const requestedMoveDate = screen.getByText(labels.requestedMoveDate);
    expect(within(requestedMoveDate.parentElement).getByText('26 Mar 2020')).toBeTruthy();

    const originAddress = screen.getByText(labels.originAddress);
    expect(
      within(originAddress.parentElement).getByText(info.originAddress.street_address_1, { exact: false }),
    ).toBeTruthy();

    const secondPickupAddress = screen.getByText(labels.secondPickupAddress);
    expect(
      within(secondPickupAddress.parentElement).getByText(info.secondPickupAddress.street_address_1, { exact: false }),
    ).toBeTruthy();

    const destinationAddress = screen.getByText(labels.destinationAddress);
    expect(
      within(destinationAddress.parentElement).getByText(info.destinationAddress.street_address_1, { exact: false }),
    ).toBeTruthy();

    const secondDestinationAddress = screen.getByText(labels.secondDestinationAddress);
    expect(
      within(secondDestinationAddress.parentElement).getByText(info.secondDestinationAddress.street_address_1, {
        exact: false,
      }),
    ).toBeTruthy();

    const releasingAgent = screen.getByText(labels.agents[0]);
    expect(within(releasingAgent.parentElement).getByText(info.agents[0].email, { exact: false })).toBeTruthy();

    const receivingAgent = screen.getByText(labels.agents[1]);
    expect(within(receivingAgent.parentElement).getByText(info.agents[1].email, { exact: false })).toBeTruthy();

    const counselorRemarks = screen.getByText(labels.counselorRemarks);
    expect(within(counselorRemarks.parentElement).getByText(info.counselorRemarks)).toBeTruthy();

    const customerRemarks = screen.getByText(labels.customerRemarks);
    expect(within(customerRemarks.parentElement).getByText(info.customerRemarks)).toBeTruthy();
  });

  it('does not render secondary addresses or agents when not provided', () => {
    render(
      <ShipmentInfoList
        requestedMoveDate={info.requestedMoveDate}
        originAddress={info.originAddress}
        destinationAddress={info.destinationAddress}
      />,
    );

    expect(screen.queryByText(labels.secondPickupAddress)).toBeFalsy();
    expect(screen.queryByText(labels.secondDestinationAddress)).toBeFalsy();
    expect(screen.queryByText(labels.agents[0])).toBeFalsy();
    expect(screen.queryByText(labels.agents[1])).toBeFalsy();
  });
});
