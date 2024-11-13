import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const destZip3 = '112';
const sameDestZip3 = '902';
const pickupZip3 = '902';
const internationalMarket = 'i';
const domesticMarket = 'd';
const isOconus = true;

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
          marketCode={domesticMarket}
          isOconus={!isOconus}
        />,
      );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 6 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });

  describe('renders the international Conus hhg shipment type with POE service items', () => {
    it.each([
      ['International Shipping & Linehaul'],
      ['International POE Fuel Surcharge'],
      ['International HHG pack'],
      ['International HHG unpack'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(
        <ShipmentServiceItemsTable
          destinationZip3={destZip3}
          pickupZip3={pickupZip3}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          marketCode={internationalMarket}
          isOconus
        />,
      );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 4 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });

  describe('renders the international Oconus hhg shipment type with POD service items', () => {
    it.each([
      ['International Shipping & Linehaul'],
      ['International POD Fuel Surcharge'],
      ['International HHG pack'],
      ['International HHG unpack'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(
        <ShipmentServiceItemsTable
          destinationZip3={destZip3}
          pickupZip3={pickupZip3}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          marketCode={internationalMarket}
          isOconus={!isOconus}
        />,
      );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 4 items', level: 4 }),
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
      render(
        <ShipmentServiceItemsTable
          destinationZip3={sameDestZip3}
          pickupZip3={pickupZip3}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          marketCode={domesticMarket}
          isOconus={!isOconus}
        />,
      );
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
          marketCode={domesticMarket}
          isOconus
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
          marketCode={domesticMarket}
          isOconus
        />,
      );
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 5 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });
});
