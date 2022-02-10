import React from 'react';
import { render, screen } from '@testing-library/react';

import NTSRShipmentInfoList from './NTSRShipmentInfoList';

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
