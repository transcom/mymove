import React from 'react';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

export default {
  title: 'TOO/TIO Components|Shipment Service Items Table',
  component: ShipmentServiceItemsTable,
};

export const HHGServiceItems = () => <ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.HHG} />;

export const NTSServiceItems = () => <ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.NTS} />;
