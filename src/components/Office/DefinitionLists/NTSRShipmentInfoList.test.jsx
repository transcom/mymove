import React from 'react';
import { render, screen } from '@testing-library/react';
import { object, text } from '@storybook/addon-knobs';

import NTSRShipmentInfoList from './NTSRShipmentInfoList';

const showWhenCollapsed = ['counselorRemarks'];
const warnIfMissing = ['ntsRecordedWeight', 'serviceOrderNumber', 'counselorRemarks', 'tacType', 'sacType'];
const errorIfMissing = ['storageFacility'];

const shipment = {
  ntsRecordedWeight: 2000,
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

describe('NTSR Shipment Info List renders all fields when provided and expanded', () => {
  it.each([
    ['ntsRecordedWeight', '2,000 lbs'],
    ['storageFacilityName', shipment.storageFacility.facilityName],
    ['serviceOrderNumber', shipment.serviceOrderNumber],
    ['storageFacilityAddress', shipment.storageFacility.address.streetAddress1],
    ['destinationAddress', shipment.destinationAddress.streetAddress1],
    ['secondaryDeliveryAddress', shipment.secondaryDeliveryAddress.streetAddress1],
    ['agent', shipment.agents[0].email, { exact: false }],
    ['counselorRemarks', shipment.counselorRemarks],
    ['customerRemarks', shipment.customerRemarks],
    ['tacType', '1234 (HHG)'],
    ['sacType', '1234123412 (NTS)'],
  ])('Verify Shipment field %s with value %s is present', async (shipmentField, shipmentFieldValue) => {
    render(<NTSRShipmentInfoList isExpanded shipment={shipment} />);
    const shipmentFieldElement = screen.getByTestId(shipmentField);
    expect(shipmentFieldElement).toHaveTextContent(shipmentFieldValue);
  });
});

describe('NTSR Shipment Info List renders missing non-required items correctly', () => {
  it.each(['counselorRemarks', 'tacType', 'sacType', 'ntsRecordedWeight', 'serviceOrderNumber'])(
    'Verify Shipment field %s displays "—" with a warning class',
    async (shipmentField) => {
      render(
        <NTSRShipmentInfoList
          isExpanded
          shipment={{
            requestedDeliveryDate: text('requestedDeliveryDate', shipment.requestedDeliveryDate),
            storageFacility: object('storageFacility', shipment.storageFacility),
            destinationAddress: object('destinationAddress', shipment.destinationAddress),
          }}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
        />,
      );
      const shipmentFieldElement = screen.getByTestId(shipmentField);
      expect(shipmentFieldElement).toHaveTextContent('—');
      expect(shipmentFieldElement.parentElement).toHaveClass('warning');
    },
  );
});

describe('NTSR Shipment Info List renders missing required items correctly', () => {
  it.each(['storageFacilityName', 'storageFacilityAddress'])(
    'Verify Shipment field %s displays "Missing" with an error class',
    async (shipmentField) => {
      render(
        <NTSRShipmentInfoList
          shipment={{
            counselorRemarks: text('counselorRemarks', shipment.counselorRemarks),
            requestedDeliveryDate: text('requestedDeliveryDate', shipment.requestedDeliveryDate),
            destinationAddress: object('destinationAddress', shipment.destinationAddress),
          }}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
        />,
      );
      const shipmentFieldElement = screen.getByTestId(shipmentField);
      expect(shipmentFieldElement).toHaveTextContent('Missing');
      expect(shipmentFieldElement.parentElement).toHaveClass('missingInfoError');
    },
  );
});

describe('NTSR Shipment Info List collapsed view', () => {
  it('hides fields when collapsed unless explicitly passed', () => {
    render(
      <NTSRShipmentInfoList
        isExpanded={false}
        shipment={shipment}
        warnIfMissing={warnIfMissing}
        errorIfMissing={errorIfMissing}
        showWhenCollapsed={showWhenCollapsed}
      />,
    );

    expect(screen.queryByTestId('ntsRecordedWeight')).toBeNull();
    expect(screen.queryByTestId('storageFacility')).toBeNull();
    expect(screen.queryByTestId('serviceOrderNumber')).toBeNull();
    expect(screen.queryByTestId('secondaryDeliveryAddress')).toBeNull();
    expect(screen.queryByTestId('agents')).toBeNull();
    expect(screen.getByTestId('counselorRemarks')).toBeInTheDocument();
  });
});
