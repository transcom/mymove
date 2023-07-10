import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const destZip3 = '112';
const sameDestZip3 = '902';
const pickupZip3 = '902';

describe('Shipment Service Items Table', () => {
  describe('renders the hhg longhaul shipment type with service items', () => {
    it.each([
      ['Domestic linehaul'],
      ['Fuel surcharge'],
      ['Domestic origin price'],
      ['Domestic destination price'],
      ['Domestic packing'],
      ['Domestic unpacking'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(
        <ShipmentServiceItemsTable
          destinationZip3={destZip3}
          pickupZip3={pickupZip3}
          shipmentType={SHIPMENT_OPTIONS.HHG}
        />,
      );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 6 items', level: 4 }),
      ).toBeInTheDocument();
      // expect(await screen.findByRole('heading', { level: 4 })).toHaveTextContent(
      //   /Service items for this shipment 6 items/,
      // );
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });

  describe('renders the hhg shorthaul shipment type with service items', () => {
    it.each([
      ['Domestic shorthaul'],
      ['Fuel surcharge'],
      ['Domestic origin price'],
      ['Domestic destination price'],
      ['Domestic packing'],
      ['Domestic unpacking'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(
        <ShipmentServiceItemsTable
          destinationZip3={sameDestZip3}
          pickupZip3={pickupZip3}
          shipmentType={SHIPMENT_OPTIONS.HHG}
        />,
      );
      // expect(await screen.findByRole('heading', { level: 4 })).toHaveTextContent(
      //   /Service items for this shipment 6 items/,
      // );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 6 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });

  describe('renders the nts shipment type with service items', () => {
    it.each([
      ['Domestic linehaul'],
      ['Fuel surcharge'],
      ['Domestic origin price'],
      ['Domestic destination price'],
      ['Domestic NTS packing'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(
        <ShipmentServiceItemsTable
          destinationZip3={destZip3}
          pickupZip3={pickupZip3}
          shipmentType={SHIPMENT_OPTIONS.NTS}
        />,
      );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 5 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });

  describe('renders the nts release shipment type with service items', () => {
    it.each([
      ['Domestic linehaul'],
      ['Fuel surcharge'],
      ['Domestic origin price'],
      ['Domestic destination price'],
      ['Domestic unpacking'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(
        <ShipmentServiceItemsTable
          destinationZip3={destZip3}
          pickupZip3={pickupZip3}
          shipmentType={SHIPMENT_OPTIONS.NTSR}
        />,
      );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 5 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });
});
