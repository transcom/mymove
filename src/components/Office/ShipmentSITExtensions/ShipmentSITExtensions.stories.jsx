import React from 'react';

import ShipmentSITExtensions from './ShipmentSITExtensions';
import {
  SITExtensions,
  SITStatusOrigin,
  SITStatusDestination,
  SITShipment,
  SITStatusWithPastSITOriginServiceItem,
  SITStatusWithPastSITServiceItems,
} from './ShipmentSITExtensionsTestParams';

export default {
  title: 'Office Components/Shipment SIT Extensions',
};

export const AtOriginNoPreviousSIT = () => <ShipmentSITExtensions sitStatus={SITStatusOrigin} shipment={SITShipment} />;

export const AtDestinationNoPreviousSIT = () => (
  <ShipmentSITExtensions sitStatus={SITStatusDestination} shipment={SITShipment} />
);

export const AtDestinationPreviousSITAtOrigin = () => (
  <ShipmentSITExtensions sitStatus={SITStatusWithPastSITOriginServiceItem} shipment={SITShipment} />
);
export const AtDestinationPreviousMulitpleSIT = () => (
  <ShipmentSITExtensions sitStatus={SITStatusWithPastSITServiceItems} shipment={SITShipment} />
);
export const AtDestinationPreviousSITAndExtension = () => (
  <ShipmentSITExtensions
    sitExtensions={SITExtensions}
    sitStatus={SITStatusWithPastSITServiceItems}
    shipment={SITShipment}
  />
);
