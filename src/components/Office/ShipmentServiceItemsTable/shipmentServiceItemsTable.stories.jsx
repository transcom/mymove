import React from 'react';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

export default {
  title: 'TOO/TIO Components|Shipment Service Items Table',
  component: ShipmentServiceItemsTable,
};

export const HHGServiceItems = () => <ShipmentServiceItemsTable shipmentType="hhg" />;

export const NTSServiceItems = () => <ShipmentServiceItemsTable shipmentType="nts" />;
