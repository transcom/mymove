import React from 'react';
import { render, screen, within } from '@testing-library/react';

import ShipmentInfoList from './ShipmentInfoList';

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

const labels = {
  requestedPickupDate: 'Requested pickup date',
  pickupAddress: 'Origin address',
  secondaryPickupAddress: 'Second pickup address',
  destinationAddress: 'Destination address',
  secondaryDeliveryAddress: 'Second destination address',
  agents: ['Releasing agent', 'Receiving agent'],
  counselorRemarks: 'Counselor remarks',
  customerRemarks: 'Customer remarks',
};

describe('Shipment Info List', () => {
  it('renders all fields when provided', () => {
    render(<ShipmentInfoList shipment={info} />);

    const requestedPickupDate = screen.getByText(labels.requestedPickupDate);
    expect(within(requestedPickupDate.parentElement).getByText('26 Mar 2020')).toBeInTheDocument();

    const pickupAddress = screen.getByText(labels.pickupAddress);
    expect(
      within(pickupAddress.parentElement).getByText(info.pickupAddress.streetAddress1, { exact: false }),
    ).toBeInTheDocument();

    const secondaryPickupAddress = screen.getByText(labels.secondaryPickupAddress);
    expect(
      within(secondaryPickupAddress.parentElement).getByText(info.secondaryPickupAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const destinationAddress = screen.getByText(labels.destinationAddress);
    expect(
      within(destinationAddress.parentElement).getByText(info.destinationAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const secondaryDeliveryAddress = screen.getByText(labels.secondaryDeliveryAddress);
    expect(
      within(secondaryDeliveryAddress.parentElement).getByText(info.secondaryDeliveryAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const releasingAgent = screen.getByText(labels.agents[0]);
    expect(within(releasingAgent.parentElement).getByText(info.agents[0].email, { exact: false })).toBeInTheDocument();

    const receivingAgent = screen.getByText(labels.agents[1]);
    expect(within(receivingAgent.parentElement).getByText(info.agents[1].email, { exact: false })).toBeInTheDocument();

    const counselorRemarks = screen.getByText(labels.counselorRemarks);
    expect(within(counselorRemarks.parentElement).getByText(info.counselorRemarks)).toBeInTheDocument();

    const customerRemarks = screen.getByText(labels.customerRemarks);
    expect(within(customerRemarks.parentElement).getByText(info.customerRemarks)).toBeInTheDocument();
  });

  it('does not render secondary addresses or agents when not provided', () => {
    render(
      <ShipmentInfoList
        shipment={{
          requestedPickupDate: info.requestedPickupDate,
          pickupAddress: info.pickupAddress,
          destinationAddress: info.destinationAddress,
        }}
      />,
    );

    expect(screen.queryByText(labels.secondaryPickupAddress)).toBeNull();
    expect(screen.queryByText(labels.secondaryDeliveryAddress)).toBeNull();
    expect(screen.queryByText(labels.agents[0])).toBeNull();
    expect(screen.queryByText(labels.agents[1])).toBeNull();
  });
});
