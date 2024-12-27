import React from 'react';

import { SHIPMENT_OPTIONS, MARKET_CODES } from '../../../shared/constants';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

export default {
  title: 'Office Components/Shipment Service Items Table',
  component: ShipmentServiceItemsTable,
};

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

export const HHGLonghaulServiceItems = () => <ShipmentServiceItemsTable shipment={domesticHhgShipment} />;

export const HHGShorthaulServiceItems = () => <ShipmentServiceItemsTable shipment={domesticHhgShipmentSameZip3} />;

export const NTSServiceItems = () => <ShipmentServiceItemsTable shipment={domesticNtsShipment} />;

export const NTSRServiceItems = () => <ShipmentServiceItemsTable shipment={domesticNtsrShipment} />;
