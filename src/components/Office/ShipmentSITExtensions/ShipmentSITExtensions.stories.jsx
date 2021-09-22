import React from 'react';

import ShipmentSITExtensions from './ShipmentSITExtensions';
import {
  SITExtensions,
  SITStatus,
  SITShipment,
  SITStatusWithPastSITServiceItems,
} from './ShipmentSITExtensionsTestParams';

export default {
  title: 'Office Components/Shipment SIT Extensions',
};

export const Default = () => (
  <ShipmentSITExtensions sitExtensions={SITExtensions} sitStatus={SITStatus} shipment={SITShipment} />
);

export const WithPastSIT = () => (
  <ShipmentSITExtensions
    sitExtensions={SITExtensions}
    sitStatus={SITStatusWithPastSITServiceItems}
    shipment={SITShipment}
  />
);
