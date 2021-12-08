import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

import { SHIPMENT_OPTIONS } from 'shared/constants';

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
      render(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC} />);
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 6 items', level: 4 }),
      ).toBeInTheDocument();
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
      render(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC} />);
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 6 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });

  describe('renders the nts shipment type with service items', () => {
    it.each([['Domestic linehaul'], ['Fuel surcharge'], ['Domestic origin price'], ['Domestic NTS packing']])(
      'expects %s to be in the document',
      async (serviceItem) => {
        render(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.NTS} />);
        expect(
          await screen.findByRole('heading', { name: 'Service items for this shipment 4 items', level: 4 }),
        ).toBeInTheDocument();
        expect(screen.getByText(serviceItem)).toBeInTheDocument();
      },
    );
  });

  describe('renders the nts release shipment type with service items', () => {
    it.each([
      ['Domestic linehaul'],
      ['Fuel surcharge'],
      ['Domestic origin price'],
      ['Domestic destination price'],
      ['Domestic unpacking'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.NTSR} />);
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 5 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });
});
