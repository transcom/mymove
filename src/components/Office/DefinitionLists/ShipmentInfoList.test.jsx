import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';

import ShipmentInfoList from './ShipmentInfoList';

import { ADDRESS_UPDATE_STATUS } from 'constants/shipments';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

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
  tertiaryPickupAddress: {
    streetAddress1: '654 S 3rd Ave',
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
  tertiaryDeliveryAddress: {
    streetAddress1: '235 Fairview Dr',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  mtoAgents: [
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
  tertiaryPickupAddress: 'Third pickup address',
  destinationAddress: 'Delivery Address',
  secondaryDeliveryAddress: 'Second delivery address',
  tertiaryDeliveryAddress: 'Third delivery address',
  mtoAgents: ['Releasing agent', 'Receiving agent'],
  counselorRemarks: 'Counselor remarks',
  customerRemarks: 'Customer remarks',
};

describe('Shipment Info List', () => {
  it('renders all fields when provided and expanded', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    render(<ShipmentInfoList isExpanded shipment={info} />);

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

    await waitFor(() => {
      const tertiaryPickupAddress = screen.getByText(labels.tertiaryPickupAddress);
      expect(
        within(tertiaryPickupAddress.parentElement).getByText(info.tertiaryPickupAddress.streetAddress1, {
          exact: false,
        }),
      ).toBeInTheDocument();
    });

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

    const tertiaryDeliveryAddress = screen.getByText(labels.tertiaryDeliveryAddress);
    expect(
      within(tertiaryDeliveryAddress.parentElement).getByText(info.tertiaryDeliveryAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const releasingAgent = screen.getByText(labels.mtoAgents[0]);
    expect(
      within(releasingAgent.parentElement).getByText(info.mtoAgents[0].email, { exact: false }),
    ).toBeInTheDocument();

    const receivingAgent = screen.getByText(labels.mtoAgents[1]);
    expect(
      within(receivingAgent.parentElement).getByText(info.mtoAgents[1].email, { exact: false }),
    ).toBeInTheDocument();

    const counselorRemarks = screen.getByText(labels.counselorRemarks);
    expect(within(counselorRemarks.parentElement).getByText(info.counselorRemarks)).toBeInTheDocument();

    const customerRemarks = screen.getByText(labels.customerRemarks);
    expect(within(customerRemarks.parentElement).getByText(info.customerRemarks)).toBeInTheDocument();
  });

  it('does not render mtoAgents when not provided', async () => {
    render(
      <ShipmentInfoList
        shipment={{
          requestedPickupDate: info.requestedPickupDate,
          pickupAddress: info.pickupAddress,
          destinationAddress: info.destinationAddress,
        }}
      />,
    );

    expect(await screen.queryByText(labels.secondaryPickupAddress)).toBeInTheDocument();
    expect(await screen.queryByText(labels.secondaryDeliveryAddress)).toBeInTheDocument();
    expect(await screen.queryByText(labels.mtoAgents[0])).not.toBeInTheDocument();
    expect(await screen.queryByText(labels.mtoAgents[1])).not.toBeInTheDocument();
  });

  it('renders appropriate fields when provided and collapsed', async () => {
    render(<ShipmentInfoList isExpanded={false} shipment={info} />);

    const requestedPickupDate = screen.getByText(labels.requestedPickupDate);
    expect(within(requestedPickupDate.parentElement).getByText('26 Mar 2020')).toBeInTheDocument();

    const pickupAddress = screen.getByText(labels.pickupAddress);
    expect(
      within(pickupAddress.parentElement).getByText(info.pickupAddress.streetAddress1, { exact: false }),
    ).toBeInTheDocument();

    expect(screen.queryByText(labels.secondaryPickupAddress)).toBeInTheDocument();

    const destinationAddress = screen.getByText(labels.destinationAddress);
    expect(
      within(destinationAddress.parentElement).getByText(info.destinationAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    expect(screen.queryByText(labels.secondaryDeliveryAddress)).toBeInTheDocument();

    expect(await screen.queryByText(labels.mtoAgents[0])).not.toBeInTheDocument();
    expect(await screen.queryByText(labels.mtoAgents[1])).not.toBeInTheDocument();

    const counselorRemarks = screen.getByText(labels.counselorRemarks);
    expect(within(counselorRemarks.parentElement).getByText(info.counselorRemarks)).toBeInTheDocument();

    const customerRemarks = screen.getByText(labels.customerRemarks);
    expect(within(customerRemarks.parentElement).getByText(info.customerRemarks)).toBeInTheDocument();
  });

  it('renders Review required instead of delivery address when the Prime has submitted a delivery address change', async () => {
    render(
      <ShipmentInfoList
        shipment={{
          requestedPickupDate: info.requestedPickupDate,
          pickupAddress: info.pickupAddress,
          destinationAddress: info.destinationAddress,
          deliveryAddressUpdate: { status: ADDRESS_UPDATE_STATUS.REQUESTED },
        }}
        errorIfMissing={[
          {
            fieldName: 'destinationAddress',
            condition: (shipment) => shipment.deliveryAddressUpdate?.status === ADDRESS_UPDATE_STATUS.REQUESTED,
            optional: true,
          },
        ]}
      />,
    );

    const destinationAddress = screen.getByText(labels.destinationAddress);
    // The delivery address will not render the address field
    // when the Prime requests a dest add update
    expect(
      within(destinationAddress.parentElement).queryByText(info.destinationAddress.streetAddress1, {
        exact: false,
      }),
    ).not.toBeInTheDocument();

    expect(within(destinationAddress.parentElement).getByText('Review required')).toBeInTheDocument();
  });
});
