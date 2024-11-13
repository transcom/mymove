import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

import { SHIPMENT_OPTIONS, MARKET_CODES } from 'shared/constants';

const destinationAddress = {
  postalCode: '11234',
  isOconus: false,
};

const destinationAddressSameZip3 = {
  postalCode: '90299',
  isOconus: false,
};

const pickupAddress = {
  postalCode: '90210',
  isOconus: false,
};

const domesticHhgShipment = {
  shipmentType: SHIPMENT_OPTIONS.HHG,
  marketCode: MARKET_CODES.DOMESTIC,
  pickupAddress,
  destinationAddress,
};

const domesticNtsShipment = {
  shipmentType: SHIPMENT_OPTIONS.NTS,
  marketCode: MARKET_CODES.DOMESTIC,
  pickupAddress,
  destinationAddress,
};

const domesticNtsrShipment = {
  shipmentType: SHIPMENT_OPTIONS.NTSR,
  marketCode: MARKET_CODES.DOMESTIC,
  pickupAddress,
  destinationAddress,
};

const domesticHhgShipmentSameZip3 = {
  shipmentType: SHIPMENT_OPTIONS.HHG,
  marketCode: MARKET_CODES.DOMESTIC,
  pickupAddress,
  destinationAddress: destinationAddressSameZip3,
};

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
      render(<ShipmentServiceItemsTable shipment={domesticHhgShipment} />);
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
      render(<ShipmentServiceItemsTable shipment={domesticHhgShipmentSameZip3} />);
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
      render(<ShipmentServiceItemsTable shipment={domesticNtsShipment} />);
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
      render(<ShipmentServiceItemsTable shipment={domesticNtsrShipment} />);
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 5 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });
});
