import React from 'react';
import { render, screen } from '@testing-library/react';
import { object, text } from '@storybook/addon-knobs';

import NTSRShipmentInfoList from './NTSRShipmentInfoList';

const showWhenCollapsed = ['counselorRemarks'];
const warnIfMissing = ['primeActualWeight', 'serviceOrderNumber', 'counselorRemarks', 'tacType', 'sacType'];
const errorIfMissing = ['storageFacility'];

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

describe('NTSR Shipment Info List', () => {
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

    const receivingAgent = screen.getByTestId('agent');
    expect(receivingAgent).toHaveTextContent(info.agents[0].email, { exact: false });

    const counselorRemarks = screen.getByTestId('counselorRemarks');
    expect(counselorRemarks).toHaveTextContent(info.counselorRemarks);

    const customerRemarks = screen.getByTestId('customerRemarks');
    expect(customerRemarks).toHaveTextContent(info.customerRemarks);

    const tacType = screen.getByTestId('tacType');
    expect(tacType).toHaveTextContent('1234 (HHG)');

    const sacType = screen.getByTestId('sacType');
    expect(sacType).toHaveTextContent('1234123412 (NTS)');
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
        warnIfMissing={warnIfMissing}
        errorIfMissing={errorIfMissing}
        showWhenCollapsed={showWhenCollapsed}
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
        warnIfMissing={warnIfMissing}
        errorIfMissing={errorIfMissing}
        showWhenCollapsed={showWhenCollapsed}
      />,
    );

    const storageFacility = screen.getByTestId('storageFacilityName');
    expect(storageFacility).toHaveTextContent('Missing');
    expect(storageFacility.parentElement).toHaveClass('missingInfoError');

    const storageFacilityAddress = screen.getByTestId('storageFacilityAddress');
    expect(storageFacilityAddress).toHaveTextContent('Missing');
    expect(storageFacilityAddress.parentElement).toHaveClass('missingInfoError');
  });

  it('hides fields when collapsed unless explicitly passed', () => {
    render(
      <NTSRShipmentInfoList
        isExpanded={false}
        shipment={info}
        warnIfMissing={warnIfMissing}
        errorIfMissing={errorIfMissing}
        showWhenCollapsed={showWhenCollapsed}
      />,
    );

    expect(screen.queryByTestId('primeActualWeight')).toBeNull();
    expect(screen.queryByTestId('storageFacility')).toBeNull();
    expect(screen.queryByTestId('serviceOrderNumber')).toBeNull();
    expect(screen.queryByTestId('secondaryDeliveryAddress')).toBeNull();
    expect(screen.queryByTestId('agents')).toBeNull();
    expect(screen.getByTestId('counselorRemarks')).toBeInTheDocument();
  });
});
