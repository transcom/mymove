import React from 'react';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

export default {
  title: 'Office Components/Shipment Service Items Table',
  component: ShipmentServiceItemsTable,
};

export const HHGServiceItems = () => <ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.HHG} />;

export const NTSRServiceItems = () => <ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.NTSR} />;
