import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { object, text } from '@storybook/addon-knobs';

import NTSRShipmentInfoList from './NTSRShipmentInfoList';

const info = {
  primeActualWeight: 2000,
  storageFacility: {
    address: {
      city: 'Anytown',
      country: 'USA',
      postalCode: '90210',
      state: 'OK',
      streetAddress1: '555 Main Ave',
      streetAddress2: 'Apartment 900',
    },
    facilityName: 'my storage',
    lotNumber: '2222',
  },
  serviceOrderNumber: '12341234',
  requestedDeliveryDate: '26 Mar 2020',
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  agents: [
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  counselorRemarks:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ' +
    'Morbi porta nibh nibh, ac malesuada tortor egestas.',
  customerRemarks: 'Ut enim ad minima veniam',
  tacType: 'HHG',
  sacType: 'NTS',
  tac: '1234',
  sac: '1234123412',
};

const labels = {
  primeActualWeight: 'Shipment weight',
  storageFacility: 'Storage facility info',
  serviceOrderNumber: 'Service order #',
  storageFacilityAddress: 'Storage facility address',
  requestedDeliveryDate: 'Preferred delivery date',
  destinationAddress: 'Delivery address',
  secondaryDeliveryAddress: 'Second delivery address',
  agents: 'Receiving agent',
  customerRemarks: 'Customer remarks',
  counselorRemarks: 'Counselor remarks',
  tacType: 'TAC',
  sacType: 'SAC',
};

describe('Shipment Info List', () => {
  it('renders all fields when provided and expanded', () => {
    render(<NTSRShipmentInfoList isExpanded shipment={info} />);
    const primeActualWeight = screen.getByTestId('primeActualWeight');
    expect(primeActualWeight).toHaveTextContent('2000 lbs');

    const storageFacility = screen.getByTestId('storageFacilityName');
    expect(storageFacility).toHaveTextContent(info.storageFacility.facilityName);

    const serviceOrderNumber = screen.getByTestId('serviceOrderNumber');
    expect(serviceOrderNumber).toHaveTextContent(info.serviceOrderNumber);

    const storageFacilityAddress = screen.getByTestId('storageFacilityAddress');
    expect(storageFacilityAddress).toHaveTextContent(info.storageFacility.address.streetAddress1);

    const destinationAddress = screen.getByTestId('destinationAddress');
    expect(destinationAddress).toHaveTextContent(info.destinationAddress.streetAddress1);

    const secondaryDeliveryAddress = screen.getByTestId('secondaryDeliveryAddress');
    expect(secondaryDeliveryAddress).toHaveTextContent(info.secondaryDeliveryAddress.streetAddress1);

    const receivingAgent = screen.getByText(labels.agents);
    expect(within(receivingAgent.parentElement).getByText(info.agents[0].email, { exact: false })).toBeInTheDocument();

    const counselorRemarks = screen.getByText(labels.counselorRemarks);
    expect(within(counselorRemarks.parentElement).getByText(info.counselorRemarks)).toBeInTheDocument();

    const customerRemarks = screen.getByText(labels.customerRemarks);
    expect(within(customerRemarks.parentElement).getByText(info.customerRemarks)).toBeInTheDocument();

    const tacType = screen.getByText(labels.tacType);
    expect(within(tacType.parentElement).getByText('1234 (HHG)')).toBeInTheDocument();

    const sacType = screen.getByText(labels.sacType);
    expect(within(sacType.parentElement).getByText('1234123412 (NTS)')).toBeInTheDocument();
  });

  it('renders a dash and adds a warning class for non-required missing items', () => {
    render(
      <NTSRShipmentInfoList
        isExpanded
        shipment={{
          requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
          storageFacility: object('storageFacility', info.storageFacility),
          destinationAddress: object('destinationAddress', info.destinationAddress),
        }}
      />,
    );

    const counselorRemarks = screen.getByTestId('counselorRemarks');
    expect(counselorRemarks).toHaveTextContent('—');
    expect(counselorRemarks.parentElement).toHaveClass('warning');

    const tacType = screen.getByTestId('tacType');
    expect(tacType).toHaveTextContent('—');
    expect(tacType.parentElement).toHaveClass('warning');

    const sacType = screen.getByTestId('sacType');
    expect(sacType).toHaveTextContent('—');
    expect(sacType.parentElement).toHaveClass('warning');

    const primeActualWeight = screen.getByTestId('primeActualWeight');
    expect(primeActualWeight).toHaveTextContent('—');
    expect(primeActualWeight.parentElement).toHaveClass('warning');

    const serviceOrderNumber = screen.getByTestId('serviceOrderNumber');
    expect(serviceOrderNumber).toHaveTextContent('—');
    expect(serviceOrderNumber.parentElement).toHaveClass('warning');
  });

  it('shows Missing and adds missing class for required missing items', () => {
    render(
      <NTSRShipmentInfoList
        shipment={{
          counselorRemarks: text('counselorRemarks', info.counselorRemarks),
          requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
          destinationAddress: object('destinationAddress', info.destinationAddress),
        }}
      />,
    );

    const storageFacility = screen.getByTestId('storageFacilityName');
    expect(storageFacility).toHaveTextContent('Missing');
    expect(storageFacility.parentElement).toHaveClass('missingInfoError');

    const storageFacilityAddress = screen.getByTestId('storageFacilityAddress');
    expect(storageFacilityAddress).toHaveTextContent('Missing');
    expect(storageFacilityAddress.parentElement).toHaveClass('missingInfoError');
  });

  it('hides fields when collapsed', () => {
    render(<NTSRShipmentInfoList isExpanded={false} shipment={info} />);

    expect(screen.queryByTestId('primeActualWeight')).toBeNull();
    expect(screen.queryByTestId('storageFacility')).toBeNull();
    expect(screen.queryByTestId('serviceOrderNumber')).toBeNull();
    expect(screen.queryByTestId('secondaryDeliveryAddress')).toBeNull();
    expect(screen.queryByTestId('agents')).toBeNull();
    expect(screen.queryByTestId('customerRemarks')).toBeNull();
  });
});
