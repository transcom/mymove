import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { object, text } from '@storybook/addon-knobs';

import NTSRShipmentInfoList from './NTSRShipmentInfoList';

import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const showWhenCollapsed = ['counselorRemarks'];
const warnIfMissing = [
  { fieldName: 'ntsRecordedWeight' },
  { fieldName: 'serviceOrderNumber' },
  { fieldName: 'counselorRemarks' },
  { fieldName: 'tacType' },
  { fieldName: 'sacType' },
];
const errorIfMissing = [{ fieldName: 'storageFacility' }];

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
  requestedPickupDate: '24 Mar 2020',
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
  tertiaryDeliveryAddress: {
    streetAddress1: '813 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  mtoAgents: [
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

const shipmentWithDeliveryAddressUpdate = {
  ...shipment,
  deliveryAddressUpdate: {
    contractorRemarks: 'dddd',
    id: 'a5e4fbfb-c8b3-4c6a-99d1-a3bec5acd52e',
    newAddress: {
      city: 'Fairfield',
      county: 'TULSA',
      eTag: 'MjAyNC0wNy0yNlQyMTozMDo1NS43OTg4NDVa',
      id: '8e553ed8-7311-479e-8f4e-d8ab974d1f6a',
      postalCode: '74133',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    newSitDistanceBetween: 1712,
    oldSitDistanceBetween: 0,
    originalAddress: {
      city: 'Fairfield',
      country: 'US',
      county: 'SOLANO',
      eTag: 'MjAyNC0wNy0yNlQxNDowNzozMS45NjY2Nlo=',
      id: 'd788f33a-f78f-48bb-a095-531555f124fc',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    shipmentID: 'c7493a05-ff3b-45ba-8420-de348c8ac63d',
    sitOriginalAddress: {
      city: 'Fairfield',
      country: 'US',
      county: 'SOLANO',
      eTag: 'MjAyNC0wNy0yNlQxNDowNzozMS45NjY2Nlo=',
      id: 'd788f33a-f78f-48bb-a095-531555f124fc',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    status: 'REQUESTED',
  },
};

describe('NTSR Shipment Info', () => {
  describe('NTSR Shipment Info List renders all fields when provided and expanded', () => {
    it.each([
      ['ntsRecordedWeight', '2,000 lbs'],
      ['storageFacilityName', shipment.storageFacility.facilityName],
      ['serviceOrderNumber', shipment.serviceOrderNumber],
      ['storageFacilityAddress', shipment.storageFacility.address.streetAddress1],
      ['requestedPickupDate', shipment.requestedPickupDate],
      ['requestedDeliveryDate', shipment.requestedDeliveryDate],
      ['destinationAddress', shipment.destinationAddress.streetAddress1],
      ['secondaryDeliveryAddress', shipment.secondaryDeliveryAddress.streetAddress1],
      ['tertiaryDeliveryAddress', shipment.tertiaryDeliveryAddress.streetAddress1],
      ['receivingAgent', shipment.mtoAgents[0].email, { exact: false }],
      ['counselorRemarks', shipment.counselorRemarks],
      ['customerRemarks', shipment.customerRemarks],
      ['tacType', '1234 (HHG)'],
      ['sacType', '1234123412 (NTS)'],
    ])('Verify Shipment field %s with value %s is present', async (shipmentField, shipmentFieldValue) => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <NTSRShipmentInfoList isExpanded shipment={shipment} />
        </MockProviders>,
      );
      await waitFor(() => {
        const shipmentFieldElement = screen.getByTestId(shipmentField);
        expect(shipmentFieldElement).toHaveTextContent(shipmentFieldValue);
      });
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
      expect(screen.queryByTestId('tertiaryDeliveryAddress')).toBeNull();
      expect(screen.queryByTestId('receivingAgent')).toBeNull();
      expect(screen.getByTestId('counselorRemarks')).toBeInTheDocument();
    });
  });

  describe('NTSR Shipment Info List Delivery Address Request', () => {
    it('renders Review required instead of delivery address when the Prime has submitted a delivery address change', async () => {
      render(
        <NTSRShipmentInfoList
          isExpanded
          shipment={shipmentWithDeliveryAddressUpdate}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
        />,
      );

      const destinationAddress = screen.getByTestId('destinationAddress');
      expect(destinationAddress).toBeInTheDocument();
      expect(destinationAddress).toHaveTextContent('Review required');
    });
  });
});
