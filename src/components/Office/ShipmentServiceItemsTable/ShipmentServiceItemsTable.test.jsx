import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

import * as api from 'services/ghcApi';
import { SHIPMENT_OPTIONS, MARKET_CODES } from 'shared/constants';

const reServiceItemResponse = [
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'POEFSC',
    serviceName: 'International POE Fuel Surcharge',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'PODFSC',
    serviceName: 'International POD Fuel Surcharge',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'UBP',
    serviceName: 'International UB',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'IUBPK',
    serviceName: 'International UB pack',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'IUBUPK',
    serviceName: 'International UB unpack',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'POEFSC',
    serviceName: 'International POE Fuel Surcharge',
    shipmentType: 'HHG',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'PODFSC',
    serviceName: 'International POD Fuel Surcharge',
    shipmentType: 'HHG',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'ISLH',
    serviceName: 'International Shipping & Linehaul',
    shipmentType: 'HHG',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'IHPK',
    serviceName: 'International HHG pack',
    shipmentType: 'HHG',
  },
  {
    isAutoApproved: true,
    marketCode: 'i',
    serviceCode: 'IHUPK',
    serviceName: 'International HHG unpack',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'ICRT',
    serviceName: 'International crating',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDASIT',
    serviceName: "International destination add'l day SIT",
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDDSIT',
    serviceName: 'International destination SIT delivery',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDFSIT',
    serviceName: 'International destination 1st day SIT',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDSHUT',
    serviceName: 'International destination shuttle service',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOASIT',
    serviceName: "International origin add'l day SIT",
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOFSIT',
    serviceName: 'International origin 1st day SIT',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOPSIT',
    serviceName: 'International origin SIT pickup',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOSHUT',
    serviceName: 'International origin shuttle service',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IUCRT',
    serviceName: 'International uncrating',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'ICRT',
    serviceName: 'International crating',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDASIT',
    serviceName: "International destination add'l day SIT",
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDDSIT',
    serviceName: 'International destination SIT delivery',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDFSIT',
    serviceName: 'International destination 1st day SIT',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDSHUT',
    serviceName: 'International destination shuttle service',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOASIT',
    serviceName: "International origin add'l day SIT",
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOFSIT',
    serviceName: 'International origin 1st day SIT',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOPSIT',
    serviceName: 'International origin SIT pickup',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOSHUT',
    serviceName: 'International origin shuttle service',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IUCRT',
    serviceName: 'International uncrating',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOSFSC',
    serviceName: 'International Origin SIT Fuel Surcharge',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDSFSC',
    serviceName: 'International Destination SIT Fuel Surcharge',
    shipmentType: 'HHG',
  },
  {
    marketCode: 'i',
    serviceCode: 'IOSFSC',
    serviceName: 'International Origin SIT Fuel Surcharge',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
  {
    marketCode: 'i',
    serviceCode: 'IDSFSC',
    serviceName: 'International Destination SIT Fuel Surcharge',
    shipmentType: 'UNACCOMPANIED_BAGGAGE',
  },
];

jest
  .spyOn(api, 'getAllReServiceItems')
  .mockImplementation(() => Promise.resolve({ data: JSON.stringify(reServiceItemResponse) }));

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

const oconusPickupAddress = {
  postalCode: '90210',
  isOconus: true,
};

const oconusDestinationAddress = {
  postalCode: '90210',
  isOconus: true,
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

const intlUbConusToOconusShipment = {
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
  marketCode: MARKET_CODES.INTERNATIONAL,
  pickupAddress,
  destinationAddress: oconusDestinationAddress,
};

const intlUbOconusToConusShipment = {
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
  marketCode: MARKET_CODES.INTERNATIONAL,
  pickupAddress: oconusPickupAddress,
  destinationAddress,
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

  describe('renders the intl UB shipment type (CONUS -> OCONUS) with service items', () => {
    it.each([
      ['International UB'],
      ['International POE Fuel Surcharge'],
      ['International UB pack'],
      ['International UB unpack'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(<ShipmentServiceItemsTable shipment={intlUbConusToOconusShipment} />);
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 4 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });

  describe('renders the intl UB shipment type (OCONUS -> CONUS) with service items', () => {
    it.each([
      ['International UB'],
      ['International POD Fuel Surcharge'],
      ['International UB pack'],
      ['International UB unpack'],
    ])('expects %s to be in the document', async (serviceItem) => {
      render(<ShipmentServiceItemsTable shipment={intlUbOconusToConusShipment} />);
      expect(
        await screen.findByRole('heading', { name: 'Service items for this shipment 4 items', level: 4 }),
      ).toBeInTheDocument();
      expect(screen.getByText(serviceItem)).toBeInTheDocument();
    });
  });
});
