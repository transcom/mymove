import React from 'react';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

export default {
  title: 'Office Components/Shipment Service Items Table',
  component: ShipmentServiceItemsTable,
};

const destZip3 = '112';
const sameDestZip3 = '902';
const pickupZip3 = '902';

export const HHGLonghaulServiceItems = () => (
  <ShipmentServiceItemsTable destinationZip3={destZip3} pickupZip3={pickupZip3} shipmentType={SHIPMENT_OPTIONS.HHG} />
);

export const HHGShorthaulServiceItems = () => (
  <ShipmentServiceItemsTable
    destinationZip3={sameDestZip3}
    pickupZip3={pickupZip3}
    shipmentType={SHIPMENT_OPTIONS.HHG}
  />
);

export const NTSServiceItems = () => (
  <ShipmentServiceItemsTable destinationZip3={destZip3} pickupZip3={pickupZip3} shipmentType={SHIPMENT_OPTIONS.NTS} />
);

export const NTSRServiceItems = () => (
  <ShipmentServiceItemsTable destinationZip3={destZip3} pickupZip3={pickupZip3} shipmentType={SHIPMENT_OPTIONS.NTSR} />
);
