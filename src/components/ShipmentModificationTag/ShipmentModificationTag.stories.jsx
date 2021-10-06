import React from 'react';

import ShipmentModificationTag from './ShipmentModificationTag';

import { shipmentModificationTypes } from 'constants/shipments';

export default {
  title: 'Components/Shipment Modification Tag',
  component: ShipmentModificationTag,
};

export const CANCELED = () => <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.CANCELED} />;
export const DIVERSION = () => (
  <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.DIVERSION} />
);
