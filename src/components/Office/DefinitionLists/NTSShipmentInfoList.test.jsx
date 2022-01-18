import React from 'react';
import { render, screen } from '@testing-library/react';
import { object, text } from '@storybook/addon-knobs';

import NTSShipmentInfoList from './NTSShipmentInfoList';

const showWhenCollapsed = ['counselorRemarks'];
const warnIfMissing = ['serviceOrderNumber', 'counselorRemarks', 'tacType', 'sacType'];
const errorIfMissing = ['storageFacility'];

const shipment = {
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
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryPickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  agents: [
    {
      agentType: 'RELEASING_AGENT',
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
  usesExternalVendor: true,
};

describe('NTS Shipment Info List renders all fields when provided and expanded', () => {
  it.each([
    ['usesExternalVendor', 'External vendor'],
    ['storageFacilityName', shipment.storageFacility.facilityName],
    ['serviceOrderNumber', shipment.serviceOrderNumber],
    ['storageFacilityAddress', shipment.storageFacility.address.streetAddress1],
    ['pickupAddress', shipment.pickupAddress.streetAddress1],
    ['secondaryPickupAddress', shipment.secondaryPickupAddress.streetAddress1],
    ['agent', shipment.agents[0].email, { exact: false }],
    ['counselorRemarks', shipment.counselorRemarks],
    ['customerRemarks', shipment.customerRemarks],
    ['tacType', '1234 (HHG)'],
    ['sacType', '1234123412 (NTS)'],
  ])('Verify Shipment field %s with value %s is present', async (shipmentField, shipmentFieldValue) => {
    render(<NTSShipmentInfoList isExpanded shipment={shipment} />);
    const shipmentFieldElement = screen.getByTestId(shipmentField);
    expect(shipmentFieldElement).toHaveTextContent(shipmentFieldValue);
  });
});

describe('NTS Shipment Info List renders missing non-required items correctly', () => {
  it.each(['counselorRemarks', 'tacType', 'sacType', 'serviceOrderNumber'])(
    'Verify Shipment field %s displays "—" with a warning class',
    async (shipmentField) => {
      render(
        <NTSShipmentInfoList
          isExpanded
          shipment={{
            requestedDeliveryDate: text('requestedDeliveryDate', shipment.requestedDeliveryDate),
            storageFacility: object('storageFacility', shipment.storageFacility),
            pickupAddress: object('pickupAddress', shipment.pickupAddress),
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

describe('NTS Shipment Info List renders missing required items correctly', () => {
  it.each(['storageFacilityName', 'storageFacilityAddress'])(
    'Verify Shipment field %s displays "Missing" with an error class',
    async (shipmentField) => {
      render(
        <NTSShipmentInfoList
          isExpanded
          shipment={{
            counselorRemarks: text('counselorRemarks', shipment.counselorRemarks),
            requestedPickupDate: text('requestedPickupDate', shipment.requestedPickupDate),
            pickupAddress: object('pickupAddress', shipment.pickupAddress),
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

describe('NTS Shipment Info List collapsed view', () => {
  it('hides fields when collapsed unless explicitly passed', () => {
    render(
      <NTSShipmentInfoList
        isExpanded={false}
        shipment={shipment}
        warnIfMissing={warnIfMissing}
        errorIfMissing={errorIfMissing}
        showWhenCollapsed={showWhenCollapsed}
      />,
    );

    expect(screen.queryByTestId('usesExternalVendor')).toBeNull();
    expect(screen.queryByTestId('storageFacility')).toBeNull();
    expect(screen.queryByTestId('serviceOrderNumber')).toBeNull();
    expect(screen.queryByTestId('secondaryPickupAddress')).toBeNull();
    expect(screen.queryByTestId('agents')).toBeNull();
    expect(screen.getByTestId('counselorRemarks')).toBeInTheDocument();
  });
});
