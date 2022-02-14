import React from 'react';
import { render, screen } from '@testing-library/react';

import NTSShipmentInfoList from './NTSShipmentInfoList';

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
};

describe('NTS Shipment Info List renders all fields when provided and expanded', () => {
  describe('For external vendors', () => {
    it.each([
      ['usesExternalVendor', 'External vendor'],
      ['requestedPickupDate', shipment.requestedPickupDate],
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
      render(<NTSShipmentInfoList isExpanded shipment={{ ...shipment, usesExternalVendor: true }} />);
      const shipmentFieldElement = screen.getByTestId(shipmentField);
      expect(shipmentFieldElement).toHaveTextContent(shipmentFieldValue);
    });
  });

  describe('For GHC prime contractor', () => {
    it.each([
      ['usesExternalVendor', 'GHC prime contractor'],
      ['requestedPickupDate', shipment.requestedPickupDate],
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
});

describe('NTS Shipment Info List renders only the requested pickup date and requested pickup address by default when collapsed when there are no props to control otherwise', () => {
  it.each([
    ['requestedPickupDate', shipment.requestedPickupDate],
    ['pickupAddress', shipment.pickupAddress.streetAddress1],
  ])('Verify Shipment field %s with value %s is present', async (shipmentField, shipmentFieldValue) => {
    render(<NTSShipmentInfoList isExpanded shipment={{ ...shipment, usesExternalVendor: true }} />);
    const shipmentFieldElement = screen.getByTestId(shipmentField);
    expect(shipmentFieldElement).toHaveTextContent(shipmentFieldValue);
  });
});

describe.each([
  ['usesExternalVendor', 'usesExternalVendor'],
  ['storageFacility', 'storageFacilityName'],
  ['serviceOrderNumber', 'serviceOrderNumber'],
  ['storageFacility', 'storageFacilityAddress'],
  ['secondaryPickupAddress', 'secondaryPickupAddress'],
  ['counselorRemarks', 'counselorRemarks'],
  ['customerRemarks', 'customerRemarks'],
  ['tacType', 'tacType'],
  ['sacType', 'sacType'],
])('NTS Shipment Info List toggles the fields', (shipmentField, testId) => {
  it(`renders ${testId} as missing and always shows when collapsed`, () => {
    render(
      <NTSShipmentInfoList
        errorIfMissing={[shipmentField]}
        shipment={{
          pickupAddress: {
            streetAddress1: '441 SW Rio de la Plata Drive',
            city: 'Tacoma',
            state: 'WA',
            postalCode: '98421',
          },
        }}
      />,
    );

    const shipmentFieldElement = screen.getByTestId(testId);
    expect(shipmentFieldElement.parentElement).toHaveClass('missingInfoError');
  });

  it(`renders ${testId} with a warning and always shows when collapsed`, () => {
    render(
      <NTSShipmentInfoList
        warnIfMissing={[shipmentField]}
        shipment={{
          pickupAddress: {
            streetAddress1: '441 SW Rio de la Plata Drive',
            city: 'Tacoma',
            state: 'WA',
            postalCode: '98421',
          },
        }}
      />,
    );

    const shipmentFieldElement = screen.getByTestId(testId);
    expect(shipmentFieldElement.parentElement).toHaveClass('warning');
  });

  it(`hides ${testId} even if it has the content and the list is expanded`, () => {
    render(<NTSShipmentInfoList isExpanded neverShow={[shipmentField]} shipment={{ ...shipment }} />);

    expect(screen.queryByTestId(testId)).toBeNull();
  });

  it(`always shows ${testId} when it has been marked as show when collapsed and the shipment list is collapsed`, () => {
    render(<NTSShipmentInfoList showWhenCollapsed={[shipmentField]} shipment={{ ...shipment }} />);

    expect(screen.getByTestId(testId)).toBeInTheDocument();
  });
});
